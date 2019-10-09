package stakeholder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDigitOnly(t *testing.T) {
	tests := []struct {
		str     string
		isValid bool
	}{
		{
			str:     "12345678904567890123",
			isValid: true,
		},
		{
			str:     "11111111111111111111111111111",
			isValid: true,
		},
		{
			str:     "0125855001-239",
			isValid: false,
		},
		{
			str:     "1a3456789012345",
			isValid: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.isValid, isDigitOnly(test.str), "Test case failed [%s]", test.str)
	}
}

// The following CPF values were automatically generated using an external tool.
func TestCheckCPF(t *testing.T) {
	tests := []struct {
		name          string
		cpf           string
		expectedError error
	}{
		{
			name:          "valid CPF",
			cpf:           "64287245881",
			expectedError: nil,
		},
		{
			name:          "invalid CPF format",
			cpf:           "153.806.664-50",
			expectedError: ErrCPFNotDigitsOnly,
		},
		{
			name:          "invalid CPF validation digit",
			cpf:           "64287245883",
			expectedError: ErrCPFValidationDigit,
		},
		{
			name:          "invalid CPF length",
			cpf:           "6428724588",
			expectedError: ErrCPFInvalidLength,
		},
		{
			name:          "empty CPF",
			cpf:           "",
			expectedError: ErrCPFEmpty,
		},
	}
	for _, test := range tests {
		err := checkCPF(test.cpf)
		assert.Equalf(t, test.expectedError, err, "Test [%s] failed.", test.name)
	}
}

func TestCheckCNPJ(t *testing.T) {
	tests := []struct {
		cnpj          string
		expectedError error
	}{
		{cnpj: "", expectedError: ErrCNPJEmpty},
		{cnpj: "1234567890123", expectedError: ErrCNPJInvalidLength},
		{cnpj: "123456789012345", expectedError: ErrCNPJInvalidLength},
		{cnpj: "1234/678901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "12345-78901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "a2345678901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "1#345678901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "12@45678901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "123^5678901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "123456 8901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "1234567{901234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "12345678$01234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "123456789~1234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "1234567890_234", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "12345678901m34", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "123456789012a4", expectedError: ErrCNPJNotDigitsOnly},
		{cnpj: "1234567890123=", expectedError: ErrCNPJNotDigitsOnly},

		// The following CNPJ's values were automatically generated using an external tool.
		{cnpj: "99859557000140", expectedError: nil},
		{cnpj: "99859557000149", expectedError: ErrCNPJValidationDigit},
		{cnpj: "31140936000141", expectedError: nil},
		{cnpj: "31140936000161", expectedError: ErrCNPJValidationDigit},
		{cnpj: "79867125000173", expectedError: nil},
		{cnpj: "79867125000111", expectedError: ErrCNPJValidationDigit},
		{cnpj: "42885777000120", expectedError: nil},
		{cnpj: "42885777000199", expectedError: ErrCNPJValidationDigit},
		{cnpj: "76623396000195", expectedError: nil},
		{cnpj: "76623396000105", expectedError: ErrCNPJValidationDigit},
		{cnpj: "49819843000103", expectedError: nil},
		{cnpj: "49819843000100", expectedError: ErrCNPJValidationDigit},
		{cnpj: "11664002000100", expectedError: nil},
		{cnpj: "11664002000111", expectedError: ErrCNPJValidationDigit},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedError, checkCNPJ(test.cnpj), "Test failed for input case [%s]", test.cnpj)
	}
}
