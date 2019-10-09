package stakeholder

import "github.com/arquivei/foundationkit/errors"

var (
	// ErrInvalidType is returned when the a Stakeholder has invalid Type value
	ErrInvalidType = errors.New("invalid stakeholder type")
	// ErrCPFInvalidLength is returned when a CPF has the wrong size
	ErrCPFInvalidLength = errors.New("cpf has invalid length")
	// ErrCPFEmpty is returned when CPF is empty
	ErrCPFEmpty = errors.New("cpf has empty value")
	// ErrCPFValidationDigit is returned when either of the validation digit is not valid
	ErrCPFValidationDigit = errors.New("cpf has wrong validation digit")
	// ErrCPFNotDigitsOnly is returned when CPF has any char that is not a number
	ErrCPFNotDigitsOnly = errors.New("cpf must be numbers only")

	// ErrCNPJInvalidLength is returned when a CNPJ has the wrong size
	ErrCNPJInvalidLength = errors.New("cnpj has invalid length")
	// ErrCNPJEmpty is returned when CPF is empty
	ErrCNPJEmpty = errors.New("cnpj has empty value")
	// ErrCNPJValidationDigit is returned when either of the validation digit is not valid
	ErrCNPJValidationDigit = errors.New("cnpj has wrong validation digit")
	// ErrCNPJNotDigitsOnly is returned when CNPJ has any char that is not a number
	ErrCNPJNotDigitsOnly = errors.New("cnpj must be numbers only")
)

func newInvalidStakeholderTypeError(t Type) error {
	return errors.Errorf("invalid stakeholder type: %v", TypeText(t))
}
