package trackingmiddleware

import (
	"net/http"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
)

// New instantiates a new tracking middleware wrapping the @next handler.
func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = request.WithNewRequestID(ctx)
		ctx = trace.WithTrace(ctx, trace.GetTraceFromHTTPRequest(r))

		request.SetInHTTPResponse(request.GetRequestIDFromContext(ctx), w)
		trace.SetTraceInHTTPResponse(trace.GetTraceFromContext(ctx), w)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
