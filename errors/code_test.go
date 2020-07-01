package errors

import (
	"errors"
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

func TestEqualsCode(t *testing.T) {
	assert.True(t, EqualsCode(Code("NO_FOOD"), Code("NO_FOOD")), "same error code")
	assert.False(t, EqualsCode(Code("RESIDENT_EVIL"), Code("VERONICA")), "different error code")
	assert.True(t, EqualsCode(ErrorCodeEmpty, ErrorCodeEmpty), "empty error code is the same")
}

func TestSameCode(t *testing.T) {
	errWithCodeTalker := E("Metal Gear Solid V: the Phantom Pain", Code("TALKER"))
	anotherErrWithCodeTalker := E("Metal Gear Solid V: the Phantom Pain (goty)", Code("TALKER"))
	errWithCodeVeronica := E("Resident Evil", Code("VERONICA"))

	assert.True(t, SameCode(errWithCodeTalker, anotherErrWithCodeTalker), "same code")
	assert.False(t, SameCode(errWithCodeTalker, errWithCodeVeronica), "different codes")
	assert.False(t, SameCode(errWithCodeTalker, errors.New("konami")), "right error has no code")
	assert.False(t, SameCode(errors.New("capcom"), errWithCodeVeronica), "left error has no code")
}
