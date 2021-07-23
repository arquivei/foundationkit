package implschemaregistry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/trace"

	"github.com/hamba/avro"
)

type repository struct {
	getSchemaByIDURL string
	getIDBySchemaURL string
}

type getSchemaByIDResponse struct {
	Schema string `json:"schema"`
}

// MustNew returns a new valid schemaregistry.Repository but panics in case on an error
func MustNew(url string, httpClient *http.Client) schemaregistry.Repository {
	const op = errors.Op("implschemaregistry.MustNew")

	if url == "" {
		panic(errors.E(op, "missing url"))
	}

	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return &repository{
		getSchemaByIDURL: url + "schemas/ids/%d",
		getIDBySchemaURL: url + "subjects/%s",
	}
}

// nolint:gosec, noctx
func (r repository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (avro.Schema, error) {
	const op = errors.Op("implschemaregistry.repository.GetSchemaByID")
	fullURL := fmt.Sprintf(r.getSchemaByIDURL, id)

	_, span := trace.StartSpan(ctx, "GetSchemaByID")
	defer span.End(nil)

	// gosec linter is given the error "Potential HTTP request made with
	// variable url" in this line, but the URL must be built from a config
	httpResponse, err := http.Get(fullURL)
	if err != nil {
		return nil, errors.E(op, err, errors.SeverityRuntime)
	}
	defer httpResponse.Body.Close()

	response := getSchemaByIDResponse{}
	decoder := json.NewDecoder(httpResponse.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return nil, errors.E(op, err, errors.SeverityInput)
	}
	schema, err := avro.Parse(response.Schema)
	if err != nil {
		return nil, errors.E(
			op,
			err,
			errors.SeverityInput,
			errors.KV("schema", truncateStr(response.Schema, 50)),
		)
	}
	return schema, nil
}

type getIDFromSchemaResponse struct {
	Subject string            `json:"subject"`
	ID      schemaregistry.ID `json:"id"`
	Version int               `json:"version"`
	Schema  string            `json:"schema"`
}

// GetIDBySchema returns the avro schema ID by using @subject and @schema.
// DO REALLY NOTE that schema is not avro.Schema, but a string instad. The reason
// is that avro.Schema.String() returns the schema in it's canonical form, which may
// unexpectedly not be recognized by the schema registry (code 40403 schema not found)
//
// nolint:gosec, noctx
func (r repository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	const op = errors.Op("implschemaregistry.repository.GetIDBySchema")

	requestBody, err := makeGetIDBySchemaRequestBody(schema)
	if err != nil {
		return 0, nil, errors.E(op, errors.SeverityFatal, err)
	}

	fullURL := fmt.Sprintf(r.getIDBySchemaURL, subject)

	// gosec linter is given the error "Potential HTTP request made with
	// variable url" in this line, but the URL must be built from a config
	httpResponse, err := http.Post(
		fullURL,
		"application/vnd.schemaregistry+json",
		strings.NewReader(requestBody),
	)
	if err != nil {
		return 0, nil, errors.E(op, errors.SeverityRuntime, err)
	}
	defer httpResponse.Body.Close()

	switch httpResponse.StatusCode {
	// This switch can be improved further by reading the "error_code"
	// field in the body, if this level of information is needed.
	case 200:
	case 404:
		return 0, nil, errors.E(
			op,
			errors.SeverityInput,
			"schema registry returned 404 - subject or schema not found",
		)
	case 500:
		return 0, nil, errors.E(
			op,
			errors.SeverityRuntime,
			"internal server error",
		)
	default:
		return 0, nil, errors.E(
			op,
			errors.SeverityRuntime,
			"unexpected status code returned",
			errors.KV("statusCode", httpResponse.StatusCode),
		)
	}

	var getResponse getIDFromSchemaResponse
	if err = json.NewDecoder(httpResponse.Body).Decode(&getResponse); err != nil {
		return 0, nil, errors.E(op, errors.SeverityInput, err)
	}

	returnedSchema, err := avro.Parse(getResponse.Schema)
	if err != nil {
		return 0, nil, errors.E(op, err)
	}

	return getResponse.ID, returnedSchema, nil
}

func makeGetIDBySchemaRequestBody(schema string) (string, error) {
	const op = errors.Op("makeGetIDBySchemaRequestBody")

	buf := new(bytes.Buffer)
	if err := json.Compact(buf, []byte(schema)); err != nil {
		return "", errors.E(op, errors.SeverityFatal, err)
	}

	body := struct {
		Schema string `json:"schema"`
	}{
		Schema: buf.String(),
	}

	marshaledBody, err := json.Marshal(&body)
	if err != nil {
		return "", errors.E(op, errors.SeverityFatal, err)
	}

	return string(marshaledBody), nil
}
