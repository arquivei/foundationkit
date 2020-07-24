package trace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTraceFromHTTRequest(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")
	r.Header.Add(headerProbabilitySample, "0.5")

	trace := GetTraceFromHTTRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Equal(t, 0.5, *trace.ProbabilitySample)
}

func TestGetTraceFromHTTRequest_ErrorParseProbabilitySample(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")
	r.Header.Add(headerProbabilitySample, "0.5a")

	trace := GetTraceFromHTTRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetTraceFromHTTRequest_WithoutHeader(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	trace := GetTraceFromHTTRequest(r)

	assert.True(t, IDIsEmpty(trace.ID))
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetTraceFromHTTRequest_WithoutProbabilitySample(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")

	trace := GetTraceFromHTTRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetTraceFromHTTRequest_WithoutTraceID(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerProbabilitySample, "0.5")

	trace := GetTraceFromHTTRequest(r)

	assert.True(t, IDIsEmpty(trace.ID))
	assert.Equal(t, 0.5, *trace.ProbabilitySample)
}

func TestSetTraceInHTTPResponse(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ID:                decode([]byte("00000000000000000000000000000019")),
		ProbabilitySample: &ps,
	}

	r := httptest.NewRecorder()
	SetTraceInHTTPResponse(trace, r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetTraceInHTTPResponse_EmptyTrace(t *testing.T) {
	defaultProbabilitySample = 0.5
	r := httptest.NewRecorder()
	SetTraceInHTTPResponse(Trace{}, r)

	assert.NotEmpty(t, r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetTraceInHTTPResponse_EmptyProbabilitySample(t *testing.T) {
	defaultProbabilitySample = 0.5
	trace := Trace{
		ID: decode([]byte("00000000000000000000000000000019")),
	}

	r := httptest.NewRecorder()
	SetTraceInHTTPResponse(trace, r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetTraceInHTTPResponse_EmptyTraceID(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ProbabilitySample: &ps,
	}

	r := httptest.NewRecorder()
	SetTraceInHTTPResponse(trace, r)

	assert.NotEmpty(t, r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetTraceInHTTPRequest(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ID:                decode([]byte("00000000000000000000000000000019")),
		ProbabilitySample: &ps,
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetTraceInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetTraceInHTTPRequest_EmptyTrace(t *testing.T) {
	defaultProbabilitySample = 0.5
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetTraceInHTTPRequest(WithTrace(context.Background(), Trace{}), r)

	assert.NotEmpty(t, r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetTraceInHTTPRequest_EmptyProbabilitySample(t *testing.T) {
	defaultProbabilitySample = 0.5
	trace := Trace{
		ID: decode([]byte("00000000000000000000000000000019")),
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetTraceInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetTraceInHTTPRequest_EmptyTraceID(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ProbabilitySample: &ps,
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetTraceInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.NotEmpty(t, r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

// [DEPRECATED] Testing a Deprecated Methods
func TestGetTraceIDFromHTTRequest(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")

	id := GetTraceIDFromHTTPRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", id.String())
}
