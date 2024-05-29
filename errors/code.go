package errors

import (
	"errors"

	"github.com/rs/zerolog"
)

// Code is the error code
type Code string

// String returns the code as a string
func (code Code) String() string {
	return string(code)
}

// MarshalZerologObject allows for zerolog to
// log the error code as 'error_code': '...'
func (code Code) MarshalZerologObject(e *zerolog.Event) {
	e.Str("error_code", string(code))
}

func (code Code) Apply(err *Error) {
	err.Code = code
}

const (
	// CodeEmpty is the zero-value for error codes
	CodeEmpty = Code("")
)

// GetCode returns the error code. If the error doesn't contains
// an error code, returns ErrorCodeEmpty
func GetCode(err error) Code {
	for {
		var e Error

		ok := errors.As(err, &e)
		if !ok {
			break
		}
		if e.Code != CodeEmpty {
			return e.Code
		}
		err = e.Err
	}

	return CodeEmpty
}

// EqualsCode returns true if @lCode and @rCode holds the same value, and
// false otherwise
func EqualsCode(lCode, rCode Code) bool {
	return (lCode == rCode)
}

// SameCode returns true if @lError and @rError holds error codes with the
// same value, and false otherwise. If one or both errors have no code, SameCode
// will return false.
func SameCode(lError, rError error) bool {
	lCode := GetCode(lError)
	rCode := GetCode(rError)
	return (lCode == rCode && lCode != CodeEmpty)
}
