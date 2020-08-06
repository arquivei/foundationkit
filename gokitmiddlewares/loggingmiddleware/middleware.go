package loggingmiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/gokitmiddlewares"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog"
)

// MustNew calls New and panics in case of error.
func MustNew(c Config) endpoint.Middleware {
	return gokitmiddlewares.Must(New(c))
}

// New returns a new go-kit logging middleware with the given name and configuration.
//
// Fields Config.Name and Config.Logger are mandatory.
func New(c Config) (endpoint.Middleware, error) {
	if c.Name == "" {
		return nil, errors.New("endpoint name is empty")
	}

	if c.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	shouldEnrichLogWithRequest := c.EnrichLogWithRequest != nil
	shouldEnrichLogWithResponse := c.EnrichLogWithResponse != nil

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			begin := time.Now()

			l, ctx := initLoggerContext(ctx, *c.Logger)

			enrichLoggerContext(ctx, l, c, req)

			if shouldEnrichLogWithRequest {
				l.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
					ctx, zctx = c.EnrichLogWithRequest(ctx, zctx, req)
					return zctx
				})
			}

			defer func() {
				enrichLoggerAfterResponse(l, c, begin, resp)

				if shouldEnrichLogWithResponse {
					l.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
						return c.EnrichLogWithResponse(ctx, zctx, resp, err)
					})
				}

				doLogging(l, c, err)
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
