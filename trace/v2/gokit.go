package trace

import (
	"context"

	"github.com/arquivei/foundationkit/request"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/attribute"
)

// EndpointMiddleware returns a new gokit endpoint.Middleware that wraps
// the next Endpoint with a span with the given name. The RequestID, if
// present, is injected as a tag and if the next Endpoint returns an
// error, it is registered in the span.
func EndpointMiddleware(name string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			ctx, span := Start(ctx, name)

			requestID := request.GetIDFromContext(ctx)
			if !requestID.IsEmpty() {
				span.SetAttributes(attribute.String("request.id", requestID.String()))
			}

			if ta, ok := req.(TraceAttributer); ok {
				span.SetAttributes(ta.TraceAttributes()...)
			}

			resp, err = next(ctx, req)
			if err != nil {
				span.RecordError(err)
			}

			if ta, ok := resp.(TraceAttributer); ok {
				span.SetAttributes(ta.TraceAttributes()...)
			}

			span.End()

			return resp, err
		}
	}
}

// TraceAttributer is a interface for returning OpenTelemetry span attributes.
// If the Request or Response (or both) implement this interface the EndpointMiddleware
// will set theses attributes on the span.
type TraceAttributer interface {
	TraceAttributes() []attribute.KeyValue
}
