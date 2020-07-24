package stakeholder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeText(t *testing.T) {
	tests := []struct {
		name     string
		input    Type
		expected string
	}{
		{
			name:     "CPF",
			input:    TypePerson,
			expected: "cpf",
		},
		{
			name:     "CNPJ",
			input:    TypeCompany,
			expected: "cnpj",
		},
		{
			name:     "Unknown",
			input:    TypeUnknown,
			expected: "unknown",
		},
		{
			name:     "Unmapped value",
			input:    7,
			expected: "unexpected: 7",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, TypeText(test.input), "[%s] Failed to check type text", test.name)
	}
}
