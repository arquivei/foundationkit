package trace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	headerTraceID           = "X-TRACEID"
	headerProbabilitySample = "X-PROBABILITYSAMPLE"
)

// GetTraceFromHTTRequest returns a Trace using the trace
// ID and the probability sample get from the header of @r
func GetTraceFromHTTRequest(r *http.Request) Trace {
	idStr := r.Header.Get(headerTraceID)
	id := decode([]byte(idStr))

	probabilitySampleStr := r.Header.Get(headerProbabilitySample)
	probabilitySample, err := strconv.ParseFloat(probabilitySampleStr, 64)

	var probabilitySamplePtr *float64
	if err == nil {
		probabilitySamplePtr = &probabilitySample
	}

	return Trace{
		ID:                id,
		ProbabilitySample: probabilitySamplePtr,
	}
}

// SetTraceInHTTPRequest sets the header of @request using the
// trace in the @ctx. If @trace is empty or @request is nil, nothing will happen
func SetTraceInHTTPRequest(ctx context.Context, request *http.Request) {
	trace := GetTraceFromContext(ctx)

	if request == nil || trace.ProbabilitySample == nil || IDIsEmpty(trace.ID) {
		log.Warn().
			Str("method", "trace.SetTraceInHTTPRequest").
			Msg("[FoundationKit] Trace is empty or request is nil. Nothing will happen")
		return
	}

	request.Header.Set(headerTraceID, trace.ID.String())
	request.Header.Set(headerProbabilitySample, fmt.Sprintf("%f", *trace.ProbabilitySample))

	return
}

// SetTraceInHTTPResponse sets the header of @response using @trace.
// If @trace is empty or @response is nil, nothing will happen
func SetTraceInHTTPResponse(trace Trace, response http.ResponseWriter) {
	if response == nil || trace.ProbabilitySample == nil || IDIsEmpty(trace.ID) {
		log.Warn().
			Str("method", "trace.SetTraceInHTTPResponse").
			Msg("[FoundationKit] Trace is empty or response is nil. Nothing will happen")
		return
	}

	response.Header().Set(headerTraceID, trace.ID.String())
	response.Header().Set(headerProbabilitySample, fmt.Sprintf("%f", *trace.ProbabilitySample))

	return
}

// GetTraceIDFromHTTPRequest attempts to return a trace ID read from the @r
// http request by obtaining the value in the `X-TRACEID` http header field.
// [DEPRECATED] Should use GetTraceFromHTTRequest instead
func GetTraceIDFromHTTPRequest(r *http.Request) ID {
	s := r.Header.Get("X-TRACEID")
	return decode([]byte(s))
}
