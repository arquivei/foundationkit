package accesskey

import (
	"github.com/arquivei/foundationkit/errors"
)

var (
	// ErrCodeInvalidAccessKey is an error code used to imply that an access key provided was not valid
	ErrCodeInvalidAccessKey = errors.Code("INVALID_ACCESS_KEY")
	// ErrCodeEmptyAccessKey is an error code used to imply that an access key provided was empty
	ErrCodeEmptyAccessKey = errors.Code("EMPTY_ACCESS_KEY")

	// ErrEmptyAccessKey is returned when the provided access key is an empty string
	ErrEmptyAccessKey = errors.New("access key is empty")
	// ErrInvalidLenght the access key does not contains 44 digits
	ErrInvalidLenght = errors.New("access key does not have 44 characters")
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
)
