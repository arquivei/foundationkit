package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPStatus(t *testing.T) {
	expectedStatus := http.StatusOK

	var err error = Error{
		HTTPStatus: 200,
		Err:        New("some error"),
	}

	s := GetHTTPStatus(err)

	assert.Equal(t, expectedStatus, s.Int())
	assert.Equal(t, http.StatusText(expectedStatus), s.String())
	assert.EqualError(t, err, "some error")
}
