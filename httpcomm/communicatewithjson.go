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
	"github.com/arquivei/foundationkit/request"
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

	details, err := communicateWithJSON(
		ctx,
		httpClient,
		httpMethod,
		fullURL,
		requestData,
		nil,
		maxAcceptedBodySize,
		maxErrBodySize,
		outResponse,
	)
	if err != nil {
		return details.Trace, errors.E(op, err)
	}

	return details.Trace, nil
}

// CommunicateWithJSONDetailed uses a given @httpClient to communicate with
// a HTTP server at @fullURL using the specified @httpMethod. This function will
// send @requestData marshalled as a JSON content and wait for the server
// response. The response will be written into the given @outResponse (it
// will mutate the parameter contents) using the json unmarshaling rules set
// by @outResponse itself. This function returns detailed information of the
// response. It is possible that the detailed information might hold some data
// even if error is not nil, but the data may not be complete.
//
// To reduce damage in case of bugs or attacks, a @maxAcceptedBodySize must be
// passed, and bodies contents larger than this will be considered noxious and
// return an error.
//
// A value of @maxErrBodySize must be passed to indicate how much of response
// contents can be added into the error message.
func CommunicateWithJSONDetailed(
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
) (ResponseDetails, error) {
	const op = errors.Op("httpcomm.CommunicateWithJSONDetailed")

	details, err := communicateWithJSON(
		ctx,
		httpClient,
		httpMethod,
		fullURL,
		requestData,
		nil,
		maxAcceptedBodySize,
		maxErrBodySize,
		outResponse,
	)
	if err != nil {
		return details, errors.E(op, err)
	}

	return details, nil
}

// CommunicateWithJSONAndHeadersDetailed uses a given @httpClient to communicate
// with a HTTP server at @fullURL using the specified @httpMethod. This function
// will send @requestData marshalled as a JSON content, with @headers as headers
// and wait for the server response. The response will be written into the given
// @outResponse (it will mutate the parameter contents) using the json
// unmarshaling rules set by @outResponse itself. This function returns detailed
// information of the response. It is possible that the detailed information
// might hold some data even if error is not nil, but the data may not be
// complete.
//
// To reduce damage in case of bugs or attacks, a @maxAcceptedBodySize must be
// passed, and bodies contents larger than this will be considered noxious and
// return an error.
//
// A value of @maxErrBodySize must be passed to indicate how much of response
// contents can be added into the error message.
func CommunicateWithJSONAndHeadersDetailed(
	ctx context.Context,
	httpClient http.Client,
	httpMethod HTTPMethod,
	fullURL string,
	requestData interface{},
	headers map[string][]string,
	maxAcceptedBodySize int64,
	maxErrBodySize int,

	// output response, on success, will be overwritten. This is not
	// a nice design, but it allows the ErrCodeDecodeError well unmarshalling
	// json to stay inside this function, and avoids the need of a cast by the caller.
	outResponse interface{},
) (ResponseDetails, error) {
	const op = errors.Op("httpcomm.CommunicateWithJSONAndHeadersDetailed")

	details, err := communicateWithJSON(
		ctx,
		httpClient,
		httpMethod,
		fullURL,
		requestData,
		headers,
		maxAcceptedBodySize,
		maxErrBodySize,
		outResponse,
	)
	if err != nil {
		return details, errors.E(op, err)
	}

	return details, nil
}

func communicateWithJSON(
	ctx context.Context,
	httpClient http.Client,
	httpMethod HTTPMethod,
	fullURL string,
	requestData interface{},
	headers map[string][]string,
	maxAcceptedBodySize int64,
	maxErrBodySize int,
	outResponse interface{},
) (ResponseDetails, error) {
	if ctx.Err() != nil {
		return ResponseDetails{}, errors.E(
			ErrCodeExpiredContext,
			errors.SeverityRuntime,
			"refusing request due to expired context",
			errors.KV("CONTEXT_ERROR", ctx.Err().Error()),
		)
	}

	httpRequest, err := makeHTTPRequest(ctx, fullURL, httpMethod, requestData, headers)
	if err != nil {
		return ResponseDetails{}, err
	}

	details, contents, err := communicateWithHTTPRequest(httpClient, maxAcceptedBodySize, httpRequest)
	if err != nil {
		return details, err
	}

	if err := json.Unmarshal(contents, outResponse); err != nil {
		contentsStr := string(contents)
		if len(contentsStr) > maxErrBodySize {
			contentsStr = string(contents[0:maxErrBodySize-5]) + "(...)"
		}

		return details, errors.E(
			ErrCodeDecodeError,
			errors.SeverityRuntime,
			fmt.Errorf("failed to decode received response: %v", err),
			errors.KV("HTTP", details.StatusCode),
			errors.KV("BODY", contentsStr),
		)
	}

	return details, nil
}

func makeHTTPRequest(
	ctx context.Context,
	fullURL string,
	httpMethod HTTPMethod,
	requestData interface{},
	headers map[string][]string,
) (*http.Request, error) {
	requestDataBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, errors.E(ErrCodeRequestError, errors.SeverityFatal, err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, httpMethod, fullURL, bytes.NewReader(requestDataBody))
	if err != nil {
		return nil, errors.E(ErrCodeRequestError, errors.SeverityFatal, err)
	}

	if trace.IDIsEmpty(trace.GetIDFromContext(ctx)) {
		// Force a random trace into request if no trace in context.
		// FIXME : trace.GetIDFromContext(ctx) results in a warning if there's no context in the trace. Implement
		//         a trace.ExistsInContext(ctx) to allow checking.
		ctx = trace.WithTrace(ctx, trace.Trace{})
	}
	trace.SetInHTTPRequest(ctx, httpRequest)

	for header, values := range headers {
		for _, value := range values {
			httpRequest.Header.Add(header, value)
		}
	}

	// NOTE : it is proposed that RequestID can be sent over HTTP Requests so that the receiver
	// side (the server) can log it; but not use it. This can aid in backtracking request is
	// systems where a same trace ID may result in many request ID's over time, and only a handful
	// of requests fails. There's no formal discussion over this proposal yet.

	return httpRequest, nil
}

func communicateWithHTTPRequest(
	httpClient http.Client,
	maxAcceptedBodySize int64,
	httpRequest *http.Request,
) (ResponseDetails, []byte, error) {
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		if isHTTPTimeoutError(err) {
			return ResponseDetails{}, nil, errors.E(ErrCodeTimeout, errors.SeverityRuntime, err)
		}

		return ResponseDetails{}, nil, errors.E(
			ErrCodeRequestError,
			errors.SeverityRuntime,
			err,
		)
	}

	details := ResponseDetails{
		StatusCode: httpResponse.StatusCode,
		Header:     httpResponse.Header,
	}

	defer httpResponse.Body.Close()
	details.Trace = trace.GetFromHTTPResponse(httpResponse)
	details.RequestID = request.GetFromHTTPResponse(httpResponse)

	limitedReader := io.LimitReader(httpResponse.Body, maxAcceptedBodySize+1)
	contents, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		return details, nil, errors.E(
			ErrCodeDecodeError,
			errors.SeverityRuntime,
			err,
		)
	}

	if int64(len(contents)) > maxAcceptedBodySize {
		return details, nil, errors.E(
			ErrCodeResponseTooLong,
			errors.SeverityRuntime,
			errors.Errorf("received contents longer than the allowed %d bytes", maxAcceptedBodySize),
		)
	}

	return details, contents, nil
}

func isHTTPTimeoutError(httpError error) bool {
	err, ok := httpError.(*url.Error)
	if !ok {
		return false
	}

	return err.Timeout()
}
