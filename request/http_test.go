package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFromHTTPRequest(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	r.Header.Add(HTTPHeaderID, "1598556827687-01EGRPJW17DMNTVAK3C6F5WQMW")
	assert.Equal(t, "1598556827687-01EGRPJW17DMNTVAK3C6F5WQMW", GetFromHTTPRequest(r).String())
}

func TestGetFromHTTPRequest_WithoutHeader(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	assert.True(t, IsEmpty(GetFromHTTPRequest(r)))
}

func TestSetInHTTPRequest(t *testing.T) {
	id := newID()

	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	SetInHTTPRequest(WithID(context.Background(), id), r)

	assert.Equal(t, id.String(), r.Header.Get(HTTPHeaderID))
}

func TestSetInHTTPRequest_EmptyID(t *testing.T) {
	r, err := http.NewRequestWithContext(context.Background(), "POST", "URL", nil)
	assert.NoError(t, err)

	SetInHTTPRequest(WithID(context.Background(), ID{}), r)

	assert.Empty(t, r.Header.Get(HTTPHeaderID))
}

func TestGetFromHTTPResponse(t *testing.T) {
	id := newID()

	response := http.Response{}
	response.Header = make(http.Header)
	response.Header.Set(HTTPHeaderID, id.String())

	assert.Equal(t, id.String(), GetFromHTTPResponse(&response).String())
}

func TestSetInHTTPResponse(t *testing.T) {
	id := newID()

	r := httptest.NewRecorder()
	SetInHTTPResponse(id, r)

	assert.Equal(t, id.String(), r.Header().Get(HTTPHeaderID))
}

func TestSetInHTTPResponse_EmptyID(t *testing.T) {
	r := httptest.NewRecorder()
	SetInHTTPResponse(ID{}, r)

	assert.Empty(t, r.Header().Get(HTTPHeaderID))
}
