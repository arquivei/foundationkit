package accesskey

import (
	"github.com/arquivei/foundationkit/errors"
)

var (
	// ErrCodeEmptyAccessKey is an error code used to imply that an access key provided was empty
	ErrCodeEmptyAccessKey = errors.Code("EMPTY_ACCESS_KEY")
	// ErrCodeInvalidLength is an error code used to imply that an access key does not contains 44 digits
	ErrCodeInvalidLength = errors.Code("INVALID_LENGTH")
	// ErrCodeInvalidCharacter is an error code used to imply that an access key has non-numeric character(s)
	ErrCodeInvalidCharacter = errors.Code("INVALID_CHARACTER")
	// ErrCodeInvalidUF is an error code used to imply that an access key does not contain a valid IBGE UF code
	ErrCodeInvalidUF = errors.Code("INVALID_UF")
	// ErrCodeInvalidMonth is an error code used to imply that an access key does not contain a month value between 01-12
	ErrCodeInvalidMonth = errors.Code("INVALID_MONTH")
	// ErrCodeInvalidCPFCNPJ is an error code used to imply that an access key does not contain a valid CNPJ
	ErrCodeInvalidCPFCNPJ = errors.Code("INVALID_CPF_CNPJ")
	// ErrCodeInvalidModel is an error code used to imply that an access key does not contain a valid SEFAZ model
	ErrCodeInvalidModel = errors.Code("INVALID_MODEL")
	// ErrCodeInvalidDigit is an error code used to imply that an access has a verification digit mismatch
	ErrCodeInvalidDigit = errors.Code("INVALID_DIGIT")

	// ErrCodeInvalidSerieForNFF
	ErrCodeInvalidSerieForNFF = errors.Code("INVALID_SERIE_FOR_NFF")
	// ErrCodeInvalidCodigoForNFF
	ErrCodeInvalidNumeroForNFF = errors.Code("INVALID_NUMERO_FOR_NFF")
	// ErrCodeInvalidCNPJForNFF
	ErrCodeInvalidCNPJForNFF = errors.Code("INVALID_CNPJ_FOR_NFF")
	// ErrCodeInvalidCPFForNFF
	ErrCodeInvalidCPFForNFF = errors.Code("INVALID_CPF_FOR_NFF")

	// ErrEmptyAccessKey is returned when the provided access key is an empty string
	ErrEmptyAccessKey = errors.New("access key is empty")
	// ErrInvalidLength the access key does not contains 44 digits
	ErrInvalidLength = errors.New("access key does not have 44 characters")
	// ErrInvalidCharacter the access key has non-numeric character(s)
	ErrInvalidCharacter = errors.New("access key contains non-number characters")
	// ErrInvalidUF the access key does not contain a valid IBGE UF code
	ErrInvalidUF = errors.New("access key has invalid UF value")
	// ErrInvalidMonth month not between 01-12
	ErrInvalidMonth = errors.New("access key has invalid month value")
	// ErrInvalidCPFCNPJ the access key does not contain a valid CNPJ
	ErrInvalidCPFCNPJ = errors.New("access key has invalid CPF or CNPJ")
	// ErrInvalidModel is returned when the model is not a valid SEFAZ model
	ErrInvalidModel = errors.New("access key has invalid model")
	// ErrInvalidDigit verification digit mismatch
	ErrInvalidDigit = errors.New("access key has invalid validation digit")

	// ErrInvalidSerieForNFF
	ErrInvalidSerieForNFF = errors.New("access key invalid serie for NFF")
	// ErrInvalidCodigoForNFF
	ErrInvalidNumeroForNFF = errors.New("access key invalid numero for NFF")
)
