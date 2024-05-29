package metricsmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/gokitmiddlewares"
	logutil "github.com/arquivei/foundationkit/log"
	"github.com/arquivei/foundationkit/metrifier"
	"github.com/rs/zerolog/log"
)

// MustNew returns a new metrics middleware but panics in case of error.
func MustNew[Request any, Response any](c Config) endpoint.Middleware[Request, Response] {
	return gokitmiddlewares.Must(New[Request, Response](c))
}

// New returns a new metrics middleware.
func New[Request any, Response any](c Config) (endpoint.Middleware[Request, Response], error) {
	m, err := metrifier.New(c.Metrifier)
	if err != nil {
		return nil, err
	}

	return func(next endpoint.Endpoint[Request, Response]) endpoint.Endpoint[Request, Response] {
		log.Debug().Str("config", logutil.Flatten(c)).Msg("[metricsmiddleware] New metrics endpoint middleware")

		return func(ctx context.Context, req Request) (resp Response, err error) {
			defer func(s metrifier.Span) {
				var r interface{}
				// Panics are handled as errors and re-raised
				if r = recover(); r != nil {
					err = errors.E(errors.NewFromRecover(r), errors.SeverityFatal, errors.CodePanic)
					log.Ctx(ctx).Warn().Err(err).
						Msg("[metricsmiddleware] Metrics endpoint middleware is handling an uncaught a panic. Please fix it!")
				}
				metrify(ctx, c.LabelsDecoder, s, req, resp, err)
				if panicErr := tryRunExternalMetrics(ctx, c.ExternalMetrics, req, resp, err); panicErr != nil {
					log.Ctx(ctx).Error().Err(panicErr).
						Msg("[metricsmiddleware] External Metrics panicked. Please check you ExternalMetrics function.")
				}

				// re-raise panic
				if r != nil {
					panic(r)
				}
			}(m.Begin())
			return next(ctx, req)
		}
	}, nil
}

func metrify(ctx context.Context, labelsDecoder LabelsDecoder, s metrifier.Span, req, resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Ctx(ctx).Error().
				Err(errors.NewFromRecover(r)).
				Msg("[metricsmiddleware] Metrics middleware panicked! Please check your code and configuration.")
		}
	}()
	if labelsDecoder != nil {
		s = s.WithLabels(labelsDecoder.Decode(ctx, req, resp, err))
	}
	s.End(err)
}

func tryRunExternalMetrics(ctx context.Context, externalMetrics ExternalMetrics, req, resp interface{}, err error) (panicErr error) {
	if externalMetrics == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr = errors.E(errors.NewFromRecover(r), errors.SeverityFatal, errors.CodePanic)
		}
	}()

	externalMetrics(ctx, req, resp, err)
	return
}
