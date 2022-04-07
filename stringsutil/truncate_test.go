package stringsutil

import (
	"strings"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestTruncate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		s        string
		i        int
		expected string
	}{
		{
			name:     "Size bigger than string",
			s:        "pudim",
			i:        100,
			expected: "pudim",
		},
		{
			name:     "Size same as string",
			s:        "flan",
			i:        4,
			expected: "flan",
		},
		{
			name:     "Size less than string",
			s:        "marshmallow",
			i:        3,
			expected: "mar",
		},
		{
			name:     "Size less than string (with UTF-8)",
			s:        "mãrshmallow",
			i:        3,
			expected: "mãr",
		},
		{
			name:     "Size is zero",
			s:        "anything goes",
			i:        0,
			expected: "",
		},
		{
			name:     "Size is 3000",
			s:        strings.Repeat("a", 4000),
			i:        3000,
			expected: strings.Repeat("a", 3000),
		},
		{
			name:     "String with UTF-8 that has a bigger length than @i but has less elements",
			s:        "00000000000000000000000000맱00000",
			i:        33,
			expected: "00000000000000000000000000맱00000",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.NotPanics(t, func() {
				assert.Equal(t, test.expected, Truncate(test.s, test.i))
			})
		})
	}
}

func FuzzTruncate(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, i int) {
		_ = Truncate(s, i)
	})
}
