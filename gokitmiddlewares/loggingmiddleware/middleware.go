package loggingmiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/gokitmiddlewares"
	logutil "github.com/arquivei/foundationkit/log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MustNew calls New and panics in case of error.
func MustNew[Request any, Response any](c Config) endpoint.Middleware[Request, Response] {
	return gokitmiddlewares.Must(New[Request, Response](c))
}

// New returns a new go-kit logging middleware with the given name and configuration.
//
// Fields Config.Name and Config.Logger are mandatory.
// Considering that this middleware puts a logger inside the context, this should always
// be the outter middleware when using endpoint.Chain.
func New[Request any, Response any](c Config) (endpoint.Middleware[Request, Response], error) {
	if c.Name == "" {
		return nil, errors.New("endpoint name is empty")
	}

	if c.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	return func(next endpoint.Endpoint[Request, Response]) endpoint.Endpoint[Request, Response] {
		log.Debug().Str("config", logutil.Flatten(c)).Msg("New logging endpoint middleware")

		return func(ctx context.Context, req Request) (resp Response, err error) {
			begin := time.Now()

			ctx = initLoggerContext(ctx, *c.Logger)
			l := log.Ctx(ctx)

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
	}
	if config.EnrichLogWithRequest != nil {
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
	}
	if config.EnrichLogWithResponse != nil {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			zctx = config.EnrichLogWithResponse(ctx, zctx, response, err)
			return zctx
		})
	}
}
