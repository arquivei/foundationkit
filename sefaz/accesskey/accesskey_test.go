package accesskey

import (
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestAccessKeyValidator(t *testing.T) {
	tests := []struct {
		name              string
		accessKey         Key
		expectedErrorCode errors.Code
	}{
		{
			name:              "Access Key with 43 characters",
			accessKey:         "1234567890123456789012345678901234567890123",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with 45 characters",
			accessKey:         "123456789012345678901234567890123456789012345",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Empty access key",
			accessKey:         "",
			expectedErrorCode: ErrCodeEmptyAccessKey,
		},
		{
			name:              "Access Key with invalid format",
			accessKey:         "1234567890123456789412A456789012345678901235",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with invalid UF",
			accessKey:         "01181010807374000258550010000066931939775149",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with invalid month",
			accessKey:         "31181310807374000258550010000066931939775149",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with invalid CNPJ",
			accessKey:         "31181274837284738271550010000066931939775149",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with invalid CPF",
			accessKey:         "31171200079713500071550010000066931939775145",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access Key with invalid model",
			accessKey:         "31181210807374000258540010000066931939775149",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "Access key with invalid validation digit",
			accessKey:         "17181011010096000195550010000143181002924980",
			expectedErrorCode: ErrCodeInvalidAccessKey,
		},
		{
			name:              "valid access key 1",
			accessKey:         "17181011010096000195550010000143181002924981",
			expectedErrorCode: "",
		},
		{
			name:              "valid access key 2",
			accessKey:         "31181010807374000258550010000066931939775149",
			expectedErrorCode: "",
		},
		{
			name:              "valid access key with contingency emission type (type 6)",
			accessKey:         "31190461186888010580550040011001536037955637",
			expectedErrorCode: "",
		},
		{
			name:              "valid access Key with CPF",
			accessKey:         "31181200079713500075550010000066931939775144",
			expectedErrorCode: "",
		},
	}

	for _, test := range tests {
		s := NewValidator()

		err := s.Check(test.accessKey)

		if test.expectedErrorCode != "" {
			assert.Equalf(t, test.expectedErrorCode, errors.GetErrorCode(err),
				"test [%s] failed to match expected error code", test.name)
		} else {
			assert.NoErrorf(t, err, "test [%s] generated an unexpected error", test.name)
		}
	}
}

func TestCPFValidator(t *testing.T) {
	tests := []struct {
		name               string
		cpf                string
		expectedValidation bool
	}{
		{
			name:               "Repeated digits (0)",
			cpf:                "00000000000000",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (1)",
			cpf:                "00011111111111",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (2)",
			cpf:                "00022222222222",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (3)",
			cpf:                "00033333333333",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (4)",
			cpf:                "00044444444444",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (5)",
			cpf:                "00055555555555",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (6)",
			cpf:                "00066666666666",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (7)",
			cpf:                "00077777777777",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (8)",
			cpf:                "00088888888888",
			expectedValidation: false,
		},
		{
			name:               "Repeated digits (9)",
			cpf:                "00099999999999",
			expectedValidation: false,
		},
		{
			name:               "More than 14 digits",
			cpf:                "000441998598459",
			expectedValidation: false,
		},
		{
			name:               "Less than 14 digits",
			cpf:                "0004419985984",
			expectedValidation: false,
		},
		{
			name:               "Not numerical only",
			cpf:                "0004B19985984",
			expectedValidation: false,
		},
		{
			name:               "Valid (1)",
			cpf:                "00007671368458",
			expectedValidation: true,
		},
		{
			name:               "Valid (2)",
			cpf:                "00008027956412",
			expectedValidation: true,
		},
		{
			name:               "Valid (3)",
			cpf:                "00038997169491",
			expectedValidation: true,
		},
		{
			name:               "Valid (4)",
			cpf:                "00005842364417",
			expectedValidation: true,
		},
		{
			name:               "Valid (5)",
			cpf:                "00003911279426",
			expectedValidation: true,
		},
		{
			name:               "Valid (6)",
			cpf:                "00003752934425",
			expectedValidation: true,
		},
		{
			name:               "Valid (7)",
			cpf:                "00003811403427",
			expectedValidation: true,
		},
		{
			name:               "Valid (8)",
			cpf:                "00007333974413",
			expectedValidation: true,
		},
		{
			name:               "Valid (9)",
			cpf:                "00007673414490",
			expectedValidation: true,
		},
		{
			name:               "Valid (10)",
			cpf:                "00009099983450",
			expectedValidation: true,
		},
		{
			name:               "Valid (11)",
			cpf:                "00009590535496",
			expectedValidation: true,
		},
		{
			name:               "Valid (12)",
			cpf:                "00038814561800",
			expectedValidation: true,
		},
		{
			name:               "Invalid (1)",
			cpf:                "00007671368451",
			expectedValidation: false,
		},
		{
			name:               "Invalid (2)",
			cpf:                "00008027956411",
			expectedValidation: false,
		},
		{
			name:               "Invalid (3)",
			cpf:                "00038997169492",
			expectedValidation: false,
		},
		{
			name:               "Invalid (4)",
			cpf:                "00005842364411",
			expectedValidation: false,
		},
		{
			name:               "Invalid (5)",
			cpf:                "00003911279421",
			expectedValidation: false,
		},
		{
			name:               "Invalid (6)",
			cpf:                "00003752934421",
			expectedValidation: false,
		},
		{
			name:               "Invalid (7)",
			cpf:                "00003811403421",
			expectedValidation: false,
		},
		{
			name:               "Invalid (8)",
			cpf:                "00007333974411",
			expectedValidation: false,
		},
		{
			name:               "Invalid (9)",
			cpf:                "00007673414491",
			expectedValidation: false,
		},
		{
			name:               "Invalid (10)",
			cpf:                "00009099983451",
			expectedValidation: false,
		},
		{
			name:               "Invalid (11)",
			cpf:                "00009590535491",
			expectedValidation: false,
		},
	}

	for _, test := range tests {

		isValid := isValidCPF(test.cpf)
		assert.Equal(t, test.expectedValidation, isValid, test.name)

	}
}

func TestCNPJCPFValidator(t *testing.T) {
	tests := []struct {
		name               string
		cpfcnpj            string
		expectedValidation bool
	}{
		{
			name:               "Invalid CNPJ with inner valid CPF",
			cpfcnpj:            "12344199859845",
			expectedValidation: false,
		},
		{
			name:               "Valid CNPJ and invalid CPF",
			cpfcnpj:            "00012345678978",
			expectedValidation: true,
		},
	}

	for _, test := range tests {

		isValid := isValidCPFCNPJ(test.cpfcnpj)
		assert.Equal(t, test.expectedValidation, isValid, test.name)

	}
}
