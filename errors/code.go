package errors

// Code is the error code
type Code string

// String returns the code as a string
func (c Code) String() string {
	return string(c)
}

const (
	// ErrorCodeEmpty is an empty error code
	ErrorCodeEmpty = Code("")
)

// GetErrorCode returns the error code. If the error doen't contains an error code, returns ErrorCodeEmpty.
func GetErrorCode(err error) Code {
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
