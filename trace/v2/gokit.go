package trace

import (
	"context"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/request"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// EndpointMiddleware returns a new gokit endpoint.Middleware that wraps
// the next Endpoint with a span with the given name. The RequestID, if
// present, is injected as a tag and if the next Endpoint returns an
// error, it is registered in the span.
func EndpointMiddleware[Request any, Response any](name string) endpoint.Middleware[Request, Response] {
	return func(next endpoint.Endpoint[Request, Response]) endpoint.Endpoint[Request, Response] {
		return func(ctx context.Context, req Request) (resp Response, err error) {
			ctx, span := Start(ctx, name)

			requestID := request.GetIDFromContext(ctx)
			if !requestID.IsEmpty() {
				span.SetAttributes(attribute.String("request.id", requestID.String()))
			}

			setAttributes(span, req)

			resp, err = next(ctx, req)
			if err != nil {
				span.RecordError(err)
			}

			setAttributes(span, resp)

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

func setAttributes(span trace.Span, v any) {
	if ta, ok := v.(TraceAttributer); ok {
		span.SetAttributes(ta.TraceAttributes()...)
	}
}
