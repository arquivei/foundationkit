package request

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	// HTTPHeaderID is the HTTP Header to be used when sending ID in
	// HTTP Requests or Responses
	HTTPHeaderID = "X-REQUESTID"
)

// GetFromHTTPRequest attempts to read a Request ID from the given @r http request's
// header. If no ID is found, or the found ID is ill-formatted, an empty ID is returned.
func GetFromHTTPRequest(r *http.Request) ID {
	idStr := r.Header.Get(HTTPHeaderID)

	id, err := Parse(idStr)
	if err != nil {
		return ID{}
	}

	return id
}

// SetInHTTPRequest will add the ID registered in @ctx into the given
// @request as a HTTP Header. If @request is nil, an warning log will
// be emitted and nothing will be changed.
func SetInHTTPRequest(ctx context.Context, request *http.Request) {
	if request == nil {
		log.Warn().
			Str("method", "request.SetInHTTPRequest").
			Msg("[FoundationKit] Request is nil. Nothing will happen")
		return
	}

	id := GetIDFromContext(ctx)
	request.Header.Set(HTTPHeaderID, id.String())
}

// GetFromHTTPResponse attempts to read a Request ID from the given @r http response's
// header. If no ID is found, or the found ID is ill-formatted, an empty ID is returned.
func GetFromHTTPResponse(r *http.Response) ID {
	idStr := r.Header.Get(HTTPHeaderID)

	id, err := Parse(idStr)
	if err != nil {
		return ID{}
	}

	return id
}

// SetInHTTPResponse will add the ID registered in @ctx into the given
// @request as a HTTP Header. If @request is nil, an warning log will
// be emitted and nothing will be changed.
//
// nolint: interfacer
func SetInHTTPResponse(id ID, response http.ResponseWriter) {
	if response == nil {
		log.Warn().
			Str("method", "request.SetInHTTPResponse").
			Msg("[FoundationKit] Response is nil. Nothing will happen")
	}

	response.Header().Set(HTTPHeaderID, id.String())
}

// HTTPMiddleware returns an http middleware that adds a request id
// to the context of the request and the same id in the header of the
// http response. If there is an active trace span, the request id is
// also registered as an attribute 'request.id'. It's important that
// this middleware comes after the trace middleware.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := newID()

			ctx := r.Context()
			ctx = WithID(ctx, id)
			r = r.WithContext(ctx)

			SetInHTTPResponse(id, w)

			if span := trace.SpanFromContext(ctx); span.IsRecording() {
				span.SetAttributes(attribute.String("request.id", id.String()))
			}

			next.ServeHTTP(w, r)
		})
}
