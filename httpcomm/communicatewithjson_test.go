package httpcomm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

func Test_communicateWithJSON_Success(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	type RequestType struct {
		Integer int
	}

	type ResponseType struct {
		Integer int
		String  string
	}

	ctx := context.Background()
	ctx = trace.WithNewTrace(ctx)
	requestTrace := trace.GetFromContext(ctx)

	var serverRequestID request.ID

	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			receivedTrace := trace.GetFromHTTPRequest(r)

			// Casting traces to String so assertion message is friendlier to human eyes
			assert.Equal(t, requestTrace.ID.String(), receivedTrace.ID.String(), "server trace id")
			assert.Equal(t, requestTrace.ProbabilitySample, receivedTrace.ProbabilitySample, "server trace probability sample")

			// checking headers
			assert.Equal(t, []string{"1", "2"}, r.Header.Values("HEADER_TEST"), "request headers")

			// TODO : need a request.NewID() method to avoid using context like this in tests
			serverRequestID = request.GetIDFromContext(request.WithNewID(context.Background()))
			trace.SetInHTTPResponse(receivedTrace, w)
			request.SetInHTTPResponse(serverRequestID, w)

			// NOTE : theres no implementation to send request id over http response yet

			w.Header().Add("some-header", "some-value")
			w.WriteHeader(http.StatusForbidden) // Should always be the last header

			var request RequestType
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				assert.FailNow(t, "Server received bad request", "Error: %v", err)
			}

			if request.Integer != 123 {
				assert.FailNow(t, "Server received request integer does not match", "Integer: expected %d, received %d", 123, request.Integer)
			}

			response := ResponseType{
				Integer: 456,
				String:  "Not stringer",
			}

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				assert.FailNow(t, "Server should not fail to encode response", "Error: %v", err)
			}
		},
	))

	var response ResponseType

	details, err := communicateWithJSON(
		ctx,
		http.Client{},
		http.MethodPost,
		testServer.URL,
		RequestType{Integer: 123},
		map[string][]string{"HEADER_TEST": {"1", "2"}},
		100,
		20,
		/*out*/ &response,
	)

	if !assert.NoError(t, err, "Should not return error") {
		return
	}

	assert.Equal(t, 456, response.Integer, "Response integer")
	assert.Equal(t, "Not stringer", response.String)

	// Casting traces to String so assertion message is friedlier to human eyes
	assert.Equal(t, requestTrace.ID.String(), details.Trace.ID.String(), "trace id")
	assert.Equal(t, requestTrace.ProbabilitySample, details.Trace.ProbabilitySample, "probability sample")
	assert.Equal(t, serverRequestID.String(), details.RequestID.String(), "request ID")
	assert.Equal(t, http.StatusForbidden, details.StatusCode, "http status code")

	if assert.True(t, len(details.Header) > 0, "should have headers") {
		assert.Equal(t, "some-value", details.Header.Get("some-header"), "custom header value")
	}
}

func Test_communicateWithJSON_CommunicationErrors(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	type RequestType struct {
		Integer int
	}

	type ResponseType struct {
		Integer int
		String  string
	}

	tests := []struct {
		name                   string
		requestInteger         int
		serverResponseStatus   int
		serverResponseContents string
		expectedCode           errors.Code
		expectedSeverity       errors.Severity
		expectedMessage        string
	}{
		{
			name:                   "Unexpected HTML response",
			requestInteger:         1,
			serverResponseStatus:   http.StatusForbidden,
			serverResponseContents: `<html><head><title>403</title></head><body>Ah ah ah, you didn't say the magic word!</body></html>`,
			expectedCode:           ErrCodeDecodeError,
			expectedSeverity:       errors.SeverityRuntime,
			expectedMessage:        "failed to decode received response: invalid character '<' looking for beginning of value [HTTP=403,BODY=<html><head><ti(...)]",
		},
		{
			name:                   "Unexpected HTML response",
			requestInteger:         1,
			serverResponseStatus:   http.StatusForbidden,
			serverResponseContents: `<html><head><title>Bacon ipsum</title></head><body>Bacon ipsum dolor amet venison andouille buffalo short ribs tenderloin ground round</body></html>`,
			expectedCode:           ErrCodeResponseTooLong,
			expectedSeverity:       errors.SeverityRuntime,
			expectedMessage:        "received contents longer than the allowed 100 bytes",
		},
	}

	for _, tc := range tests {
		testCase := tc

		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			testServer := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					var request RequestType
					err := json.NewDecoder(r.Body).Decode(&request)
					if err != nil {
						assert.FailNow(tt, "Received bad request", "Error: %v", err)
					}

					if request.Integer != testCase.requestInteger {
						assert.FailNow(tt, "Received request integer does not match", "Integer: expected %d, received %d", testCase.requestInteger, request.Integer)
					}

					w.WriteHeader(tc.serverResponseStatus)
					fmt.Fprintln(w, testCase.serverResponseContents)
				},
			))

			var response ResponseType

			_, err := communicateWithJSON(
				context.Background(),
				http.Client{},
				http.MethodPost,
				testServer.URL,
				RequestType{Integer: testCase.requestInteger},
				nil,
				100,
				20,
				/*out*/ response,
			)

			assert.Error(tt, err, "Error expected")
			assert.Equal(tt, testCase.expectedCode, errors.GetCode(err), "Error code")
			assert.Equal(tt, testCase.expectedSeverity, errors.GetSeverity(err), "Error severity")
			assert.EqualError(tt, errors.GetRootErrorWithKV(err), testCase.expectedMessage, "Error message")
		})
	}
}

func Test_communicateWithJSON_Timeout(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
		},
	))

	var request struct{}
	var response struct{}

	_, err := communicateWithJSON(
		context.Background(),
		http.Client{Timeout: 1},
		http.MethodPost,
		testServer.URL,
		request,
		nil,
		2*(1<<10), // 2KB
		20,
		/*out*/ response,
	)

	assert.Error(t, err, "Error expected")
	assert.Equal(t, ErrCodeTimeout, errors.GetCode(err), "Error code")
	assert.Equal(t, errors.SeverityRuntime, errors.GetSeverity(err), "Error severity")
	// Not checking message as it changes on every test
}

func Test_communicateWithJSON_ExpiredContext(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
		},
	))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var request struct{}
	var response struct{}

	_, err := communicateWithJSON(
		ctx,
		http.Client{Timeout: 1},
		http.MethodPost,
		testServer.URL,
		request,
		nil,
		2*(1<<10), // 2KB
		20,
		/*out*/ response,
	)

	assert.Error(t, err, "Error expected")
	assert.Equal(t, ErrCodeExpiredContext, errors.GetCode(err), "Error code")
	assert.Equal(t, errors.SeverityRuntime, errors.GetSeverity(err), "Error severity")
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "refusing request due to expired context [CONTEXT_ERROR=context canceled]", "Error message")
}
