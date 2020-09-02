package httpcomm

import (
	"net/http"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
)

// ResponseDetails holds details of the received response, such as HTTP Status
// and headers, if available.
type ResponseDetails struct {
	Trace      trace.Trace
	RequestID  request.ID
	StatusCode int
	Header     http.Header
}
