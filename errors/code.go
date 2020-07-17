package errors

import "github.com/rs/zerolog"

// Code is the error code
type Code string

// String returns the code as a string
func (c Code) String() string {
	return string(c)
}

// MarshalZerologObject allows for zerolog to
// log the error code as 'error_code': '...'
func (c Code) MarshalZerologObject(e *zerolog.Event) {
	e.Str("error_code", string(c))
}

const (
	// ErrorCodeEmpty is an empty error code
	ErrorCodeEmpty = Code("")
)

// GetCode returns the error code. If the error doesn't contains
// an error code, returns ErrorCodeEmpty
func GetCode(err error) Code {
	for {
		e, ok := err.(Error)
		if !ok {
			break
		}
		if e.Code != ErrorCodeEmpty {
			return e.Code
		}
		err = e.Err
	}

	return ErrorCodeEmpty
}

// GetErrorCode returns the error code. If the error doesn't contains an error code, returns ErrorCodeEmpty.
//
// Deprecated: use GetCode instead.
func GetErrorCode(err error) Code {
	return GetCode(err)
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
	return (lCode == rCode && lCode != ErrorCodeEmpty)
}
