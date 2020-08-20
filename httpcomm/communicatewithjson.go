package httpcomm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/trace"
)

// CommunicateWithJSON uses a given @httpClient to communicate with a HTTP
// server at @fullURL using the specified @httpMethod. This function will
// send @requestData marshalled as a JSON content and wait for the server
// response. The response will be written into the given @outResponse (it
// will mutate the parameter contents) using the json unmarshaling rules set
// by @outResponse itself.
//
// To reduce damage in case of bugs or attacks, a @maxAcceptedBodySize must be
// passed, and bodies contents larger than this will be considered noxious and
// return an error.
//
// A value of @maxErrBodySize must be passed to indicate how much of response
// contents can be added into the error message.
func CommunicateWithJSON(
	ctx context.Context,
	httpClient http.Client,
	httpMethod HTTPMethod,
	fullURL string,
	requestData interface{},
	maxAcceptedBodySize int64,
	maxErrBodySize int,

	// output response, on success, will be overwritten. This is not
	// a nice design, but it allows the ErrCodeDecodeError well unmarshalling
	// json to stay inside this function, and avoids the need of a cast by the caller.
	outResponse interface{},
) (trace.Trace, error) {
	const op = errors.Op("httpcomm.CommunicateWithJSON")

	if ctx.Err() != nil {
		return trace.Trace{}, errors.E(
			op,
			ErrCodeExpiredContext,
			errors.SeverityRuntime,
			"refusing request due to expired context",
			errors.KV("CONTEXT_ERROR", ctx.Err().Error()),
		)
	}

	requestDataBody, err := json.Marshal(requestData)
	if err != nil {
		return trace.Trace{}, errors.E(op, ErrCodeRequestError, errors.SeverityFatal, err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, httpMethod, fullURL, bytes.NewReader(requestDataBody))
	if err != nil {
		return trace.Trace{}, errors.E(op, ErrCodeRequestError, errors.SeverityFatal, err)
	}

	if trace.IDIsEmpty(trace.GetIDFromContext(ctx)) {
		// Force a random trace into request if no trace in context.
		// FIXME : trace.GetIDFromContext(ctx) results in a warning if there's no context in the trace. Implement
		//         a trace.ExistsInContext(ctx) to allow checking.
		ctx = trace.WithTrace(ctx, trace.Trace{})
	}
	trace.SetTraceInHTTPRequest(ctx, httpRequest)

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		if isHTTPTimeoutError(err) {
			return trace.Trace{}, errors.E(op, ErrCodeTimeout, errors.SeverityRuntime, err)
		}

		return trace.Trace{}, errors.E(
			op,
			ErrCodeRequestError,
			errors.SeverityRuntime,
			err,
		)
	}

	defer httpResponse.Body.Close()
	responseTrace := trace.GetTraceFromHTTPResponse(httpResponse)

	limitedReader := io.LimitReader(httpResponse.Body, maxAcceptedBodySize+1)
	contents, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return trace.Trace{}, errors.E(
			op,
			ErrCodeDecodeError,
			errors.SeverityRuntime,
			err,
		)
	}

	if int64(len(contents)) > maxAcceptedBodySize {
		return responseTrace, errors.E(
			op,
			ErrCodeResponseTooLong,
			errors.SeverityRuntime,
			errors.Errorf("received contents longer than the allowed %d bytes", maxAcceptedBodySize),
		)
	}

	if err := json.Unmarshal(contents, outResponse); err != nil {
		contentsStr := string(contents)
		if len(contentsStr) > maxErrBodySize {
			contentsStr = string(contents[0:maxErrBodySize-5]) + "(...)"
		}

		return responseTrace, errors.E(
			op,
			ErrCodeDecodeError,
			errors.SeverityRuntime,
			fmt.Errorf("failed to decode received response: %v", err),
			errors.KV("HTTP", httpResponse.StatusCode),
			errors.KV("BODY", contentsStr),
		)
	}

	return responseTrace, nil
}

func isHTTPTimeoutError(httpError error) bool {
	err, ok := httpError.(*url.Error)
	if !ok {
		return false
	}

	return err.Timeout()
}
