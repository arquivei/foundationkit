package apiping

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/arquivei/foundationkit/errors"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHTTPHandler returns a new http handler for endpoint
func MakeHTTPHandler(e endpoint.Endpoint) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	httpHandler := kithttp.NewServer(
		e,
		decodeRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/ping/v1", httpHandler).Methods("POST")

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body Request

	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&body)

	return body, nil
}

func encodeResponse(
	ctx context.Context,
	w http.ResponseWriter,
	r interface{},
) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(r)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	r := ResponseError{
		Error: err.Error(),
	}
	w.WriteHeader(getHTTPStatus(err))
	json.NewEncoder(w).Encode(r)
}

func getHTTPStatus(err error) (s int) {
	switch errors.GetSeverity(err) {
	case errors.SeverityInput:
		return http.StatusBadRequest
	}

	// If we don't know what happened...
	return http.StatusInternalServerError
}
