package apiutil

import (
	"context"
	"encoding/json"
	"net/http"
)

// EncodeJSONResponse encodes an endpoint response into json
func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
