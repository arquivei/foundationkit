package stringutils

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestTruncate(t *testing.T) {
	type input struct {
		str  string
		size int
	}
	tests := []struct {
		name     string
		input    input
		expected string
	}{
		{
			name:     "Size bigger than string",
			input:    input{str: "pudim", size: 100},
			expected: "pudim",
		},
		{
			name:     "Size same as string",
			input:    input{str: "flan", size: 4},
			expected: "flan",
		},
		{
			name:     "Size less than string",
			input:    input{str: "marshmallow", size: 3},
			expected: "mar",
		},
		{
			name:     "Size less than string (with UTF-8)",
			input:    input{str: "mãrshmallow", size: 3},
			expected: "mãr",
		},
		{
			name:     "Size is zero",
			input:    input{str: "anything goes", size: 0},
			expected: "",
		},
	}

	for _, test := range tests {
		actual := Truncate(test.input.str, test.input.size)
		assert.Equal(t, test.expected, actual, "[%s] Unexpected value", test.name)

	}
}
