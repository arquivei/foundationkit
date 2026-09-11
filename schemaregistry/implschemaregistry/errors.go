package implschemaregistry

import "github.com/arquivei/foundationkit/errors"

const (
	ErrCodeNetworkError                = errors.Code("ErrCodeNetworkError")
	ErrCodeDecodeResponseError         = errors.Code("ErrCodeDecodeResponseError")
	ErrCodeParseSchemaError            = errors.Code("ErrCodeParseSchemaError")
	ErrCodeMakeRequestBodyError        = errors.Code("ErrCodeMakeRequestBodyError")
	ErrCodeNotFound                    = errors.Code("ErrCodeNotFound")
	ErrCodeServerError                 = errors.Code("ErrCodeServerError")
	ErrCodeUnexpectedResponseCodeError = errors.Code("ErrCodeUnexpectedResponseCodeError")
)
