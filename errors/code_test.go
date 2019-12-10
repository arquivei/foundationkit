package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCode(t *testing.T) {
	assert.Equal(t, "error code", GetErrorCode(E(New("err"), Code("error code"))).String())
}

func TestErrorCode_WithoutErrorCode(t *testing.T) {
	var err error
	assert.Equal(t, ErrorCodeEmpty, GetErrorCode(err))

	err = New("my error")
	assert.Equal(t, ErrorCodeEmpty, GetErrorCode(err))

	err = Error{
		Err: New("my error"),
	}
	assert.Equal(t, ErrorCodeEmpty, GetErrorCode(err))

	assert.Equal(t, ErrorCodeEmpty, GetErrorCode(nil))
}
