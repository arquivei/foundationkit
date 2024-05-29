package endpoint

import "context"

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Endpoint[Request any, Response any] func(ctx context.Context, request Request) (response Response, err error)

// Middleware is a chainable behavior modifier for endpoints.
type Middleware[Request any, Response any] func(Endpoint[Request, Response]) Endpoint[Request, Response]

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
func Chain[Request any, Response any](outer Middleware[Request, Response], others ...Middleware[Request, Response]) Middleware[Request, Response] {
	return func(next Endpoint[Request, Response]) Endpoint[Request, Response] {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func Gokit[Request any, Response any](e Endpoint[Request, Response]) func(ctx context.Context, request interface{}) (response interface{}, err error) {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return e(ctx, request.(Request))
	}
}
