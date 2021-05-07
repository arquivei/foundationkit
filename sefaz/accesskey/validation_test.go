package accesskey

import (
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestCheck_ShouldReturnNilForValidAccessKeys(t *testing.T) {
	key := AccessKey("35190901180169000133550010000297701022024512")

	err := Check(key)

	assert.Nil(t, err)

	key = AccessKey("35210501180169000133551210501100013022024517")

	err = CheckNFF(key)

	assert.Nil(t, err)
}

func TestCheck_ShouldReturnValidationCodesForNFe(t *testing.T) {
	tests := map[string]struct {
		key          AccessKey
		expectedCode errors.Code
	}{
		"Empty access key": {
			key:          "",
			expectedCode: ErrCodeEmptyAccessKey,
		},
		"Invalid length": {
			key:          "351909011801690001335500100002",
			expectedCode: ErrCodeInvalidLength,
		},
		"Non numeric key": {
			key:          "3519O901180169O00133550010000297701022024512",
			expectedCode: ErrCodeInvalidCharacter,
		},
		"Invalid UF": {
			key:          "95190901180169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidUF,
		},
		"Invalid month": {
			key:          "35194201180169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidMonth,
		},
		"Invalid CPF": {
			key:          "42160300005475399984558910004353741212308338",
			expectedCode: ErrCodeInvalidCPFCNPJ,
		},
		"Invalid CNPJ": {
			key:          "35190901429169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidCPFCNPJ,
		},
		"Invalid model": {
			key:          "35190901180169000133420010000297701022024512",
			expectedCode: ErrCodeInvalidModel,
		},
		"Invalid digit": {
			key:          "35190901180169000133550010000297701022024513",
			expectedCode: ErrCodeInvalidDigit,
		},
	}

	for testName, testCase := range tests {
		err := Check(testCase.key)
		code := errors.GetCode(err)

		if code != testCase.expectedCode {
			t.Errorf("Wrong validation for key '%s', expected an '%s' error", testCase.key, testName)
		}
	}
}

func TestCheck_ShouldReturnValidationCodesForNFF(t *testing.T) {
	tests := map[string]struct {
		key          AccessKey
		expectedCode errors.Code
	}{
		"Empty access key": {
			key:          "",
			expectedCode: ErrCodeEmptyAccessKey,
		},
		"Invalid length": {
			key:          "351909011801690001335500100002",
			expectedCode: ErrCodeInvalidLength,
		},
		"Non numeric key": {
			key:          "3519O901180169O00133550010000297701022024512",
			expectedCode: ErrCodeInvalidCharacter,
		},
		"Invalid UF": {
			key:          "95190901180169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidUF,
		},
		"Invalid month": {
			key:          "35194201180169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidMonth,
		},
		"Invalid CPF": {
			key:          "42160300005475399984558910004353741212308338",
			expectedCode: ErrCodeInvalidCPFCNPJ,
		},
		"Invalid CNPJ": {
			key:          "35190901429169000133550010000297701022024512",
			expectedCode: ErrCodeInvalidCPFCNPJ,
		},
		"Invalid model": {
			key:          "35190901180169000133420010000297701022024512",
			expectedCode: ErrCodeInvalidModel,
		},
		"Invalid digit": {
			key:          "35190901180169000133550010000297701022024513",
			expectedCode: ErrCodeInvalidDigit,
		},
		"Invalid NFF - Invalid serie": {
			key:          "35210501180169000133550210501100013022024512",
			expectedCode: ErrCodeInvalidSerieForNFF,
		},
		"Invalid NFF - Invalid numero: emission date": {
			key:          "35210501180169000133551210532100013022024515",
			expectedCode: ErrCodeInvalidNumeroForNFF,
		},
		"Invalid NFF - Invalid numero: emitter cpf or cnpj": {
			key:          "35210501180169000133551210501300013022024514",
			expectedCode: ErrCodeInvalidNumeroForNFF,
		},
		"Invalid NFF - Invalid numero: invalid cnpj": {
			key:          "35210500091387151630551210501100013022024511",
			expectedCode: ErrCodeInvalidCNPJForNFF,
		},
		"Invalid NFF - Invalid numero: invalid cpf": {
			key:          "35210501180169000133551210501200013022024510",
			expectedCode: ErrCodeInvalidCPFForNFF,
		},
	}

	for testName, testCase := range tests {
		err := CheckNFF(testCase.key)
		code := errors.GetCode(err)

		if code != testCase.expectedCode {
			t.Errorf("Wrong validation for key '%s', expected an '%s' error", testCase.key, testName)
		}
	}
}
