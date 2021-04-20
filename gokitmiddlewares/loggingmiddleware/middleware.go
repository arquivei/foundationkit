package loggingmiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/gokitmiddlewares"
	logutil "github.com/arquivei/foundationkit/log"

	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MustNew calls New and panics in case of error.
func MustNew(c Config) endpoint.Middleware {
	return gokitmiddlewares.Must(New(c))
}

// New returns a new go-kit logging middleware with the given name and configuration.
//
// Fields Config.Name and Config.Logger are mandatory.
// Considering that this middleware puts a logger inside the context, this should always
// be the outter middleware when using endpoint.Chain.
func New(c Config) (endpoint.Middleware, error) {
	if c.Name == "" {
		return nil, errors.New("endpoint name is empty")
	}

	if c.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		log.Debug().Str("config", logutil.Flatten(c)).Msg("New logging endpoint middleware")

		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			begin := time.Now()

			l, ctx := initLoggerContext(ctx, *c.Logger)

			enrichLoggerContext(ctx, l, c, req)
			ctx = doCustomEnrichRequest(ctx, c, l, req)

			defer func() {
				var r interface{}
				// Panics are handled as errors and re-raised
				if r = recover(); r != nil {
					err = errors.NewFromRecover(r)
					log.Ctx(ctx).Warn().Err(err).
						Msg("Logging endpoint middleware is handling an uncaught a panic. Please fix it!")
				}
				enrichLoggerAfterResponse(l, c, begin, resp)
				doCustomEnrichResponse(ctx, c, l, resp, err)

				doLogging(l, c, err)
				if r != nil {
					panic(r)
				}
			}()

			return next(ctx, req)
		}
	}, nil
}

func doLogging(l *zerolog.Logger, c Config, err error) {
	if err != nil {
		l.WithLevel(getErrorLevel(c, err)).
			Err(err).
			EmbedObject(errors.GetCode(err)).
			EmbedObject(errors.GetSeverity(err)).
			Msg("Request failed")
		return
	}

	l.WithLevel(c.SuccessLevel).Msg("Request successful")
}

func toString(i interface{}, n int) string {
	s := fmt.Sprintf("%#v", i)
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n]
}

func doCustomEnrichRequest(
	ctx context.Context,
	config Config,
	logger *zerolog.Logger,
	request interface{},
) context.Context {
	if typedReq, ok := request.(LoggableEndpointRequest); ok {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			ctx, zctx = typedReq.EnrichLog(ctx, zctx)
			return zctx
		})
	} else if config.EnrichLogWithRequest != nil {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			ctx, zctx = config.EnrichLogWithRequest(ctx, zctx, request)
			return zctx
		})
	}
	return ctx
}

func doCustomEnrichResponse(
	ctx context.Context,
	config Config,
	logger *zerolog.Logger,
	response interface{},
	err error,
) {
	if typedReq, ok := response.(LoggableEndpointResponse); ok {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			zctx = typedReq.EnrichLog(ctx, zctx)
			return zctx
		})
	} else if config.EnrichLogWithResponse != nil {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			zctx = config.EnrichLogWithResponse(ctx, zctx, response, err)
			return zctx
		})
	}
}
