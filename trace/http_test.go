package trace

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFromHTTPRequest(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")
	r.Header.Add(headerProbabilitySample, "0.5")

	trace := GetFromHTTPRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Equal(t, 0.5, *trace.ProbabilitySample)
}

func TestGetFromHTTPRequest_ErrorParseProbabilitySample(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")
	r.Header.Add(headerProbabilitySample, "0.5a")

	trace := GetFromHTTPRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetFromHTTPRequest_WithoutHeader(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	trace := GetFromHTTPRequest(r)

	assert.True(t, IDIsEmpty(trace.ID))
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetFromHTTPRequest_WithoutProbabilitySample(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerTraceID, "00000000000000000000000000000019")

	trace := GetFromHTTPRequest(r)

	assert.Equal(t, "00000000000000000000000000000019", trace.ID.String())
	assert.Nil(t, trace.ProbabilitySample)
}

func TestGetFromHTTPRequest_WithoutTraceID(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(headerProbabilitySample, "0.5")

	trace := GetFromHTTPRequest(r)

	assert.True(t, IDIsEmpty(trace.ID))
	assert.Equal(t, 0.5, *trace.ProbabilitySample)
}

func TestGetFromHTTPResponse(t *testing.T) {
	ps := 0.75
	trace := Trace{
		ID:                decode([]byte("000000000000000000000000bebacafe")),
		ProbabilitySample: &ps,
	}

	response := http.Response{}
	response.Header = make(http.Header)
	response.Header.Set(headerTraceID, trace.ID.String())
	response.Header.Set(headerProbabilitySample, fmt.Sprintf("%f", *trace.ProbabilitySample))

	assert.Equal(t, "000000000000000000000000bebacafe", response.Header.Get(headerTraceID))
	assert.Equal(t, "0.750000", fmt.Sprintf("%f", *trace.ProbabilitySample))
}

func TestSetInHTTPResponse(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ID:                decode([]byte("00000000000000000000000000000019")),
		ProbabilitySample: &ps,
	}

	r := httptest.NewRecorder()
	SetInHTTPResponse(trace, r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetInHTTPResponse_EmptyTrace(t *testing.T) {
	defaultProbabilitySample = 0.5
	r := httptest.NewRecorder()
	SetInHTTPResponse(Trace{}, r)

	assert.NotEmpty(t, r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetInHTTPResponse_EmptyProbabilitySample(t *testing.T) {
	defaultProbabilitySample = 0.5
	trace := Trace{
		ID: decode([]byte("00000000000000000000000000000019")),
	}

	r := httptest.NewRecorder()
	SetInHTTPResponse(trace, r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetInHTTPResponse_EmptyTraceID(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ProbabilitySample: &ps,
	}

	r := httptest.NewRecorder()
	SetInHTTPResponse(trace, r)

	assert.NotEmpty(t, r.Header().Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header().Get(headerProbabilitySample))
}

func TestSetInHTTPRequest(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ID:                decode([]byte("00000000000000000000000000000019")),
		ProbabilitySample: &ps,
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetInHTTPRequest_EmptyTrace(t *testing.T) {
	defaultProbabilitySample = 0.5
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetInHTTPRequest(WithTrace(context.Background(), Trace{}), r)

	assert.NotEmpty(t, r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetInHTTPRequest_EmptyProbabilitySample(t *testing.T) {
	defaultProbabilitySample = 0.5
	trace := Trace{
		ID: decode([]byte("00000000000000000000000000000019")),
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.Equal(t, "00000000000000000000000000000019", r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}

func TestSetInHTTPRequest_EmptyTraceID(t *testing.T) {
	ps := 0.5
	trace := Trace{
		ProbabilitySample: &ps,
	}

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)
	SetInHTTPRequest(WithTrace(context.Background(), trace), r)

	assert.NotEmpty(t, r.Header.Get(headerTraceID))
	assert.Equal(t, "0.500000", r.Header.Get(headerProbabilitySample))
}
