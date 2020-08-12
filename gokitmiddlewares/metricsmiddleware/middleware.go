package metricsmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/gokitmiddlewares"
	logutil "github.com/arquivei/foundationkit/log"
	"github.com/arquivei/foundationkit/metrifier"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
)

// LabelsDecoder defines a functions that generates a label value by processing the
// endpoint's request and response.
type LabelsDecoder func(ctx context.Context, req, resp interface{}, err error) map[string]string

// MustNew returns a new metrics middleware but panics in case of error.
func MustNew(c Config) endpoint.Middleware {
	return gokitmiddlewares.Must(New(c))
}

// New returns a new metrics middleware.
func New(c Config) (endpoint.Middleware, error) {
	m, err := metrifier.New(c.Metrifier)
	if err != nil {
		return nil, err
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		log.Debug().Str("config", logutil.Flatten(c)).Msg("New metrics endpoint middleware")

		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			defer func(s metrifier.Span) {
				var r interface{}
				// Panics are handled as errors and re-raised
				if r = recover(); r != nil {
					err = errors.E(errors.NewFromRecover(r), errors.SeverityFatal, errors.CodePanic)
					log.Ctx(ctx).Warn().Err(err).
						Msg("Metrics endpoint middleware is handling an uncaught a panic. Please fix it!")
				}
				metrify(ctx, c.LabelsDecoder, s, req, resp, err)
				if r != nil {
					panic(r)
				}
			}(m.Begin())
			return next(ctx, req)
		}
	}, nil
}

func metrify(ctx context.Context, dec LabelsDecoder, s metrifier.Span, req, resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Ctx(ctx).Error().
				Err(errors.NewFromRecover(r)).
				Msg("Metrics middleware panicked! Please check your code and configuration.")
		}
	}()
	if dec != nil {
		s = s.WithLabels(dec(ctx, req, resp, err))
	}
	s.End(err)
}
