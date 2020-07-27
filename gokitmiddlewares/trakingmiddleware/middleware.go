package trackingmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
	"github.com/go-kit/kit/endpoint"
)

// New returns a new endpoint middleware that ensures that the context has a request ID and a trace.
//
// The request and trace packages both check the context before adding the information in the context so
// this middleware is safe for usign with middlewares from the trasport layer that may have placed a
// trace or a request ID in the context.
//
// For the trace to work properly, it's expected that the trace.Setup() functions was called previously.
func New() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			ctx = request.WithRequestID(ctx)
			ctx = trace.WithNewTrace(ctx)
			return next(ctx, r)
		}
	}
}
