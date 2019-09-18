package trace

import (
	"net/http"
)

// GetTraceIDFromHTTPRequest attempts to return a trace ID read from the @r
// http request by obtaining the value in the `X-TRACEID` http header field.
func GetTraceIDFromHTTPRequest(r *http.Request) ID {
	s := r.Header.Get("X-TRACEID")
	return Decode([]byte(s))
}
