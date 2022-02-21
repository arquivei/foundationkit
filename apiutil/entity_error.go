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
	return ErrorDescription{
		Code:    getErrorCode(err).String(),
		Message: errors.GetRootErrorWithKV(err).Error(),
	}
}

func getErrorCode(err error) errors.Code {
	switch code := errors.GetCode(err); code {
	case errors.CodeEmpty:
		return getErrorCodeBasedOnSeverity(errors.GetSeverity(err))
	default:
		return code
	}
}

func getErrorCodeBasedOnSeverity(code errors.Severity) errors.Code {
	switch code {
	case errors.SeverityInput:
		return ErrCodeBadRequest
	default:
		return ErrCodeInternal
	}
}
