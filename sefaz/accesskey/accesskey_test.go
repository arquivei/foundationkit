package accesskey

import (
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestAccessKeyValidator(t *testing.T) {
	tests := []struct {
		name              string
		accessKey         AccessKey
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
