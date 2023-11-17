package trace

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// MuxHTTPMiddleware sets up a handler to start tracing the incoming
// requests. The service parameter should describe the name of the
// (virtual) server handling the request.
// This will be set in the tag 'net.host.name'. It defaults to the
// ip.
func MuxHTTPMiddleware(service string) mux.MiddlewareFunc {
	return otelmux.Middleware(service)
}

// SetTraceInRequest will put the trace in @r headers
func SetTraceInRequest(r *http.Request) {
	otel.GetTextMapPropagator().Inject(
		r.Context(),
		propagation.HeaderCarrier(r.Header),
	)
}

// SetTraceInResponse will put the trace in @r headers
func SetTraceInResponse(ctx context.Context, r http.ResponseWriter) {
	otel.GetTextMapPropagator().Inject(
		ctx,
		propagation.HeaderCarrier(r.Header()),
	)
}
