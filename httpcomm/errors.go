package httpcomm

import "github.com/arquivei/foundationkit/errors"

const (
	// ErrCodeExpiredContext is returned when an operation won't be
	// executed due to it's context being expired before the operation
	// start.
	ErrCodeExpiredContext errors.Code = "EXPIRED_CONTEXT"

	// ErrCodeRequestError is returned when a request could not be generated
	// for some reason or other.
	ErrCodeRequestError errors.Code = "REQUEST_ERROR"

	// ErrCodeDecodeError is returned when a received response body could not be
	// decoded. This usually means transport layer errors, such as HTTP or DNS.
	ErrCodeDecodeError errors.Code = "DECODE_ERROR"

	// ErrCodeResponseTooLong is returned when a received response body is too
	// long and can be assumed as an serious malfunction or an attack.
	ErrCodeResponseTooLong errors.Code = "RESPONSE_TOO_LONG"

	// ErrCodeTimeout is returned when a client side timeout on the HTTP Client
	// is detected. Note that as there are various ways of a timeout occurring, this
	// error code might not be returned in every type of timeout error happening.
	ErrCodeTimeout errors.Code = "TIMEOUT"

	// ErrCodeMissing is returned when a received response has an error without
	// code. This should never happen, indicating unexpected behavior in the
	// HTTP Server.
	ErrCodeMissing errors.Code = "CODE_MISSING"
)
