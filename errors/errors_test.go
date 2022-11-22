package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("error")
	assert.EqualError(t, err, "error")
}

func TestOpAppending(t *testing.T) {
	err := E(SeverityFatal, "1st")
	err = E(err, Op("2nd"))
	err = E(Op("3rd"), err)
	err = E(err, Op("4th"))
	assert.Equal(t, "4th: 3rd: 2nd: 1st", err.Error())
}

func TestErrorf(t *testing.T) {
	err := Errorf("the following treta occurs: [%d] %s", 19, "TRETA")
	if !assert.EqualError(t, err, "the following treta occurs: [19] TRETA") {
		return
	}
}

func TestGetRootError(t *testing.T) {
	err := E("a", fmt.Errorf("err a"))
	err = E(Op("b"), err)
	err = E(Op("c"), err)
	err = E(Op("d"), err)
	assert.Equal(t, New("err a"), GetRootError(err))
}

func TestGetRootErrorNormalError(t *testing.T) {
	err := fmt.Errorf("err a")
	assert.Equal(t, New("err a"), GetRootError(err))
}

func TestConcatErrorsMessage(t *testing.T) {
	errs := ConcatErrorsMessage(New("a"), New("b"), New("c"), New("d"), New("e"), New("f"))
	assert.Equal(t, "a: b: c: d: e: f", errs)
}

func TestConcatErrors(t *testing.T) {
	errs := ConcatErrors(New("a"), New("b"), New("c"), New("d"), New("e"), New("f"))
	assert.EqualError(t, errs, "a: b: c: d: e: f")
}

func TestE_EmptyArgs(t *testing.T) {
	assert.True(t, strings.HasPrefix(E().Error(), "errors.E called with 0 args"))
}

func TestE_NoError(t *testing.T) {
	assert.NoError(t, E(Code("no error")))
}

func TestErrorString_NoError(t *testing.T) {
	var stringError string
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "root error is nil", r)
			assert.Equal(t, "", stringError)
		}
	}()
	err := Error{}
	stringError = err.Error()
	assert.FailNow(t, "panic didn't occurred")
}

func TestErrorWrap(t *testing.T) {
	rootError := testError("root error")

	err := E(Op("a"), rootError)
	assert.Equal(t, rootError, errors.Unwrap(err))
	err = E(Op("b"), err)
	err = E(Op("c"), err)
	err = E(Op("d"), err)

	assert.True(t, errors.Is(err, rootError))

	var destError testError
	assert.True(t, errors.As(err, &destError))
	assert.Equal(t, rootError, destError)
}

type testError string

func (e testError) Error() string {
	return string(e)
}

func TestErrorWrapMixed(t *testing.T) {
	err := E("root error", Code("CODE"), SeverityFatal)

	err = fmt.Errorf("wrapped: %w", err)
	err = E(Op("a"), err)

	assert.Equal(t, Code("CODE"), GetCode(err))
	assert.Equal(t, SeverityFatal, GetSeverity(err))
}
