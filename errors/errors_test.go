package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("error")
	assert.EqualError(t, err, "error")
}

func TestOpAppending(t *testing.T) {
	err := New("1st", SeverityFatal)
	err = E(err, Op("2nd"))
	err = E(err, Op("3rd"))
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
	err := E(fmt.Errorf("err a"))
	err = E(err, Op("b"))
	err = E(err, Op("c"))
	err = E(err, Op("d"))
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

func TestE_NoError(t *testing.T) {
	assert.NoError(t, E(nil, Code("no error")))
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

	err := E(rootError, Op("a"))
	assert.Equal(t, rootError, errors.Unwrap(err))
	err = E(err, Op("b"))
	err = E(err, Op("c"))
	err = E(err, Op("d"))

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
	err := New("root error", Code("CODE"), SeverityFatal)

	err = fmt.Errorf("wrapped: %w", err)
	err = E(err, Op("a"))

	assert.Equal(t, Code("CODE"), GetCode(err))
	assert.Equal(t, SeverityFatal, GetSeverity(err))
}

func ExampleE() {
	// Calling E with no arguments results in an error with a message saying
	// that the function was called with no arguments
	errNoArgs := E(nil)
	fmt.Println(errNoArgs)

	// Calling E without an error will result in a nil return
	errNil := E(nil, Op("Error Example"), Code("ERROR_EXAMPLE_NIL"))
	fmt.Println(errNil)

	// E requires either a string or an err to return an error
	withString := New("This string will be used to build an error")
	previous := errors.New("Previous error")
	withError := E(previous)
	fmt.Println(withString, withError)

	// We can pass parameters of the same type more than once, but only the last
	// one will be considered
	multiOp := New("Multi op", Op("Op 1"), Op("Op 2"), Op("Op 3"))
	fmt.Println(multiOp)

	// Except for KeyValue and []KeyValue which will be concatenated
	kv := KeyValue{Key: "key1", Value: "val1"}
	multiKv := New("Multi kv", kv, kv, kv)
	fmt.Println(multiKv)

	// Full example
	op := Op("errors.errorExample")
	code := Code("ERROR_EXAMPLE")
	sev := SeverityRuntime
	kv = KeyValue{Key: "key", Value: "val"}

	err := New("Error example", op, code, sev, kv)
	fmt.Println(err)
}
