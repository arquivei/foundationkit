package loggingmiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
)

// New returns a new go-kit logging middleware with the given name and configuration.
// It will panic if the name is empty.
func New(name string, c Config) endpoint.Middleware {
	if name == "" {
		panic("endpoint name is empty")
	}

	if c.Logger == nil {
		panic("logger is nil")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			begin := time.Now()

			ctx = initLoggerContext(ctx, *c.Logger)

			enrichLoggerContext(ctx, name, c, req)

			defer func() {
				doLogging(ctx, c, begin, resp, err)
			}()

			return next(ctx, req)
		}
	}
}

func doLogging(ctx context.Context, c Config, begin time.Time, resp interface{}, err error) {
	enrichLoggerAfterResponse(ctx, c, begin, resp)

	l := log.Ctx(ctx)

	if err != nil {
		e := l.WithLevel(getErrorLevel(c, err)).Err(err)
		if code := errors.GetCode(err); code.String() != "" {
			e = e.Str("error_code", code.String())
		}
		if s := errors.GetSeverity(err); s.String() != "" {
			e = e.Str("error_severity", s.String())
		}
		e.Msg("Request failed")
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
