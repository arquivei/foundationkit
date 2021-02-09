package contextmap

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// NewEndpointMiddleware ensures that there is a ContextMap in the context of an endpoint.
func NewEndpointMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			ctx = New().WithCtx(ctx)
			return next(ctx, req)
		}
	}
}
