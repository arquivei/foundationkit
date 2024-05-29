package dontpanicmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog/log"
)

// New returns a new Don't Panic middleware.
func New[Request any, Response any]() endpoint.Middleware[Request, Response] {
	return func(next endpoint.Endpoint[Request, Response]) endpoint.Endpoint[Request, Response] {
		log.Debug().Msg("New dontpanic endpoint middleware")
		return func(ctx context.Context, req Request) (resp Response, err error) {
			panicErr := errors.DontPanic(func() {
				resp, err = next(ctx, req)
			})
			if panicErr != nil {
				return resp, panicErr
			}
			return
		}
	}
}
