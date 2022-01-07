package apiutil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
	"github.com/rs/zerolog/log"
)

// GetHTTPStatusFunc is a type of function that should
// take an error and return an HTTP status code.
type GetHTTPStatusFunc func(error) int

// ParseErrorFunc is a type of function that should
// take an error parses in a struct.
type ParseErrorFunc func(context.Context, error) interface{}

// NewHTTPErrorJSONEncoder returns a new
func NewHTTPErrorJSONEncoder(
	getHTTPStatus GetHTTPStatusFunc,
	parseErrorFunc ParseErrorFunc,
) func(context.Context, error, http.ResponseWriter) {
	if getHTTPStatus == nil {
		panic("getHTTPStatus is nil")
	}

	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		resp := parseErrorFunc(ctx, err)

		w.WriteHeader(getHTTPStatus(err))
		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			log.Error().Err(errors.E(err, errors.KV("encode", encodeErr))).
				EmbedObject(trace.GetFromContext(ctx)).
				EmbedObject(request.GetIDFromContext(ctx)).
				Msg("Failed to write endpoint error response")
		}
	}
}

// GetDefaultErrorHTTPStatusCode returns an HTTP status code base on an error.
//
// This is a default implementation that sets the satus code base on the error
// severity.
func GetDefaultErrorHTTPStatusCode(err error) (s int) {
	switch errors.GetSeverity(err) {
	case errors.SeverityInput:
		return http.StatusBadRequest
	}

	switch errors.GetCode(err) {
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	}

	// If we don't know what happend...
	return http.StatusInternalServerError
}
