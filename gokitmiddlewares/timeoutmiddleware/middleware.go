package timeoutmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/go-kit/kit/endpoint"
)

// New returns a new timeout middleware.
//
// After timeout is reached, if the middleware is configured to wait,
// it will just cancel the context and wait for next endpoint to return.
// But if the middleware is configured to not wait, it will run the next endpoint
// inside a go-routine and return error as soon as the context is canceled.
func New(c Config) (endpoint.Middleware, error) {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		// Timeout is disabled
		if c.Timeout <= 0 {
			return next
		}

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx, cancel := context.WithTimeout(ctx, c.Timeout)
			defer cancel()

			defer func() {
				if err != nil && ctx.Err() != nil {
					err = errors.E(err, c.ErrorSeverity)
				}
			}()

			if c.Wait {
				return next(ctx, request)
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case r := <-runAsync(ctx, next, request):

				return r.response, r.err
			}
		}
	}, nil
}

type asyncResult struct {
	response interface{}
	err      error
}

// runAsync executes next inside a go-routine and returns the result in a channel.
func runAsync(ctx context.Context, next endpoint.Endpoint, request interface{}) <-chan asyncResult {
	c := make(chan asyncResult)

	go func() {
		defer close(c)

		// Panics in go-routines must be captured inside the go-routine
		err := errors.DontPanic(func() {
			response, err := next(ctx, request)
			c <- asyncResult{
				response: response,
				err:      err,
			}
		})
		if err != nil {
			c <- asyncResult{
				err: err,
			}
		}
	}()

	return c
}
