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
// this middleware is safe for using with middlewares from the transport layer that may have placed a
// trace or a request ID in the context.
//
// For the trace to work properly, it's expected that the trace.Setup() functions was called previously.
func New() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, r interface{}) (interface{}, error) {
			// Ensures that we have a request ID in the context
			ctx = request.WithNewRequestID(ctx)

			// Get trace from request if there is one or else
			// creates a new trace.
			if t, ok := r.(Traceable); ok {
				ctx = trace.WithTrace(ctx, t.Trace())
			} else {
				ctx = trace.WithNewTrace(ctx)
			}

			return next(ctx, r)
		}
	}
}
