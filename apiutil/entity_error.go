package apiutil

import "github.com/arquivei/foundationkit/errors"

const (
	// ErrCodeInternal is returned when an internal error happens.
	ErrCodeInternal = errors.Code("INTERNAL_ERROR")

	// ErrCodeBadRequest is returned when an error happens due to request data.
	ErrCodeBadRequest = errors.Code("BAD_REQUEST")

	// ErrCodeTimeout is returned then the request returns an error due to a timeout.
	ErrCodeTimeout = errors.Code("REQUEST_TIMEOUT")
)

// ErrorDescription represents the detailed returned error in all APIs
type ErrorDescription struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ParseError parses @err in a ErrorDescription
func ParseError(err error) ErrorDescription {
	// This should never happen, but...
	if err == nil {
		return ErrorDescription{
			Code:    ErrCodeInternal.String(),
			Message: "trying to encode nil error",
		}
	}
	code := errors.GetCode(err)
	if code == errors.CodeEmpty {
		code = ErrCodeInternal
	}
	return ErrorDescription{
		Code:    code.String(),
		Message: errors.GetRootErrorWithKV(err).Error(),
	}
}
