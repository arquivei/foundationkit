package timeoutmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/errors"
)

// New returns a new timeout middleware.
//
// After timeout is reached, if the middleware is configured to wait,
// it will just cancel the context and wait for next endpoint to return.
// But if the middleware is configured to not wait, it will run the next endpoint
// inside a go-routine and return error as soon as the context is canceled.
func New[Request any, Response any](c Config) (endpoint.Middleware[Request, Response], error) {
	return func(next endpoint.Endpoint[Request, Response]) endpoint.Endpoint[Request, Response] {
		// Timeout is disabled
		if c.Timeout <= 0 {
			return next
		}

		return func(ctx context.Context, request Request) (response Response, err error) {
			ctx, cancel := context.WithTimeout(ctx, c.Timeout)
			defer cancel()

			// Override error code and severity based on the context
			defer func() {
				if err != nil && ctx.Err() != nil {
					err = errors.E(err, c.ErrorSeverity, c.ErrorCode)
				}
			}()

			if c.Wait {
				return next(ctx, request)
			}
			return nextNoWait(ctx, next, request)
		}
	}, nil
}

// nextNoWait runs next but don't wait for a response in case of canceled context
func nextNoWait[Request any, Response any](ctx context.Context, next endpoint.Endpoint[Request, Response], request Request) (Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-runNextAsync(ctx, next, request):
		return r.response, r.err
	}
}

type asyncResult struct {
	response interface{}
	err      error
}

// runNextAsync executes next inside a go-routine and returns the result in a channel.
func runNextAsync[Request any, Response any](ctx context.Context, next endpoint.Endpoint[Request, Response], request Request) <-chan asyncResult {
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
