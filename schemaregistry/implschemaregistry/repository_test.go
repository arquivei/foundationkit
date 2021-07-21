package implschemaregistry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro"
	"github.com/stretchr/testify/assert"
)

var (
	tagsSchemaStr = `{"name":"simple","type":"record","fields":[{"name":"a","type":"long"},{"name":"b","type":"string"}]}`
	tagsSchema    = avro.MustParse(tagsSchemaStr)
)

func TestGetSchemaByID(t *testing.T) {
	tests := []struct {
		name          string
		httpResponse  map[string]interface{}
		expectedError string
	}{
		{
			name: "Success",
			httpResponse: map[string]interface{}{
				"schema": tagsSchemaStr,
			},
		},
		{
			name: "Error - Decode Response",
			httpResponse: map[string]interface{}{
				"schema": 1,
			},
			expectedError: "implschemaregistry.repository.GetSchemaByID: json: cannot unmarshal number into Go struct field getSchemaByIDResponse.schema of type string",
		},
		{
			name: "Error - Empty Schema",
			httpResponse: map[string]interface{}{
				"schema": "",
			},
			expectedError: "implschemaregistry.repository.GetSchemaByID: avro: unknown type:  [schema=]",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, "/schemas/ids/10", req.URL.String())
				assert.Equal(t, http.MethodGet, req.Method)
				rw.WriteHeader(200)
				out, _ := json.Marshal(test.httpResponse)
				_, _ = rw.Write(out)
			}))
			defer server.Close()
			repository := MustNew(server.URL, nil)
			schema, err := repository.GetSchemaByID(context.Background(), schemaregistry.ID(10))
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tagsSchemaStr, schema.String())
			}
		})
	}
}

func TestGetIDBySchema(t *testing.T) {
	tests := []struct {
		name string

		schema   string
		schemaID schemaregistry.ID
		subject  schemaregistry.Subject

		httpResponse map[string]interface{}
		statusCode   int

		expectedError string
	}{
		{
			name:     "Success - 200",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id":      schemaregistry.ID(10),
				"subject": schemaregistry.Subject("mysubject"),
				"version": 1,
				"schema":  tagsSchemaStr,
			},
			statusCode: 200,
		},
		{
			name:     "Error - 404",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id":      schemaregistry.ID(10),
				"subject": schemaregistry.Subject("mysubject"),
				"version": 1,
				"schema":  tagsSchemaStr,
			},
			statusCode:    404,
			expectedError: "implschemaregistry.repository.GetIDBySchema: schema registry returned 404 - subject or schema not found",
		},
		{
			name:     "Error - 500",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id":      schemaregistry.ID(10),
				"subject": schemaregistry.Subject("mysubject"),
				"version": 1,
				"schema":  tagsSchemaStr,
			},
			statusCode:    500,
			expectedError: "implschemaregistry.repository.GetIDBySchema: internal server error",
		},
		{
			name:     "Error - 666",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id":      schemaregistry.ID(10),
				"subject": schemaregistry.Subject("mysubject"),
				"version": 1,
				"schema":  tagsSchemaStr,
			},
			statusCode:    666,
			expectedError: "implschemaregistry.repository.GetIDBySchema: unexpected status code returned [statusCode=666]",
		},
		{
			name:          "Error - Invalid Schema",
			schema:        "",
			expectedError: "implschemaregistry.repository.GetIDBySchema: makeGetIDBySchemaRequestBody: unexpected end of JSON input",
		},
		{
			name:     "Error - Response parser",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id": "bla",
			},
			statusCode:    200,
			expectedError: "implschemaregistry.repository.GetIDBySchema: json: cannot unmarshal string into Go struct field getIDFromSchemaResponse.id of type schemaregistry.ID",
		},
		{
			name:     "Error - Empty schema",
			schema:   tagsSchemaStr,
			schemaID: schemaregistry.ID(10),
			subject:  schemaregistry.Subject("mysubject"),
			httpResponse: map[string]interface{}{
				"id":      schemaregistry.ID(10),
				"subject": schemaregistry.Subject("mysubject"),
				"version": 1,
				"schema":  "",
			},
			statusCode:    200,
			expectedError: "implschemaregistry.repository.GetIDBySchema: avro: unknown type: ",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Equal(t, "/subjects/mysubject", req.URL.String())
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "application/vnd.schemaregistry+json", req.Header.Get("Content-Type"))

				var body map[string]interface{}
				err := json.NewDecoder(req.Body).Decode(&body)
				assert.NoError(t, err)
				assert.Equal(t, body["schema"], test.schema)
				rw.WriteHeader(test.statusCode)
				out, _ := json.Marshal(test.httpResponse)
				_, _ = rw.Write(out)
			}))
			defer server.Close()

			repository := MustNew(server.URL, server.Client())
			id, schema, err := repository.GetIDBySchema(context.Background(), test.subject, test.schema)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.schemaID, id)
				assert.Equal(t, test.schema, schema.String())
			}
		})
	}
}

func TestMustNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		MustNew("", nil)
	})
}
