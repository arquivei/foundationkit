package loggingmiddleware

import (
	"context"
	"fmt"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog"
)

// MustNew calls New and panics in case of error.
func MustNew(name string, c Config) endpoint.Middleware {
	m, err := New(name, c)
	if err != nil {
		panic(err)
	}
	return m
}

// New returns a new go-kit logging middleware with the given name and configuration.
// It will panic if the name is empty.
func New(name string, c Config) (endpoint.Middleware, error) {
	if name == "" {
		return nil, errors.New("endpoint name is empty")
	}

	if c.Logger == nil {
		return nil, errors.New("logger is nil")
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			begin := time.Now()

			l, ctx := initLoggerContext(ctx, *c.Logger)

			enrichLoggerContext(ctx, l, name, c, req)
			defer func() {
				enrichLoggerAfterResponse(l, c, begin, resp)
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
