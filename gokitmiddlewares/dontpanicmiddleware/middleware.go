package dontpanicmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
)

// New returns a new Don't Panic middleware.
func New() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		log.Debug().Msg("New dontpanic endpoint middleware")
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			panicErr := errors.DontPanic(func() {
				resp, err = next(ctx, req)
			})
			if panicErr != nil {
				return nil, panicErr
			}
			return
		}
	}
}
