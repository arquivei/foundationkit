package cuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUFSuccess(t *testing.T) {
	tests := []struct {
		str     string
		isValid bool
	}{
		{str: "11", isValid: true}, {str: "12", isValid: true},
		{str: "13", isValid: true}, {str: "14", isValid: true},
		{str: "15", isValid: true}, {str: "16", isValid: true},
		{str: "17", isValid: true}, {str: "21", isValid: true},
		{str: "22", isValid: true}, {str: "23", isValid: true},
		{str: "24", isValid: true}, {str: "25", isValid: true},
		{str: "26", isValid: true}, {str: "27", isValid: true},
		{str: "28", isValid: true}, {str: "29", isValid: true},
		{str: "31", isValid: true}, {str: "32", isValid: true},
		{str: "33", isValid: true}, {str: "35", isValid: true},
		{str: "41", isValid: true}, {str: "42", isValid: true},
		{str: "43", isValid: true}, {str: "50", isValid: true},
		{str: "51", isValid: true}, {str: "52", isValid: true},
		{str: "53", isValid: true},

		{str: "", isValid: false}, {str: "000", isValid: false},
		{str: "0", isValid: false}, {str: "00", isValid: false},
		{str: "1", isValid: false}, {str: "9", isValid: false},
		{str: "10", isValid: false}, {str: "18", isValid: false},
		{str: "19", isValid: false}, {str: "20", isValid: false},
		{str: "30", isValid: false}, {str: "34", isValid: false},
		{str: "36", isValid: false}, {str: "40", isValid: false},
		{str: "45", isValid: false}, {str: "49", isValid: false},
		{str: "54", isValid: false}, {str: "55", isValid: false},
		{str: "56", isValid: false}, {str: "57", isValid: false},
		{str: "foo", isValid: false}, {str: "bar", isValid: false},
	}
	for _, test := range tests {
		assert.Equalf(t, test.isValid, isValidUF(test.str), "UF validation failed for uf %s", test.str)
	}
}

func TestCUFString(t *testing.T) {
	cuf := CUF("35")
	assert.Equal(t, "35", cuf.String())
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected CUF
		err      string
	}{
		{
			name:     "Valid cUF",
			input:    "35",
			expected: CUF("35"),
			err:      "",
		},
		{
			name:     "Invalid cUF",
			input:    "00",
			expected: CUF(""),
			err:      "invalid cUF",
		},
		{
			name:     "Missing cUF",
			input:    "",
			expected: CUF(""),
			err:      "missing cUF",
		},
	}

	for _, test := range tests {
		cuf, err := New(test.input)
		if test.err != "" {
			assert.EqualError(t, err, test.err, "[%s] Error mismatch", test.name)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expected, cuf, "[%s] cUF mismatch", test.name)
	}
}

func TestCUF_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		data       []byte
		isValid    bool
		expectedUF CUF
	}{
		{
			data:       []byte(`"31"`),
			expectedUF: CUF("31"),
			isValid:    true,
		},
		{
			data:       []byte(`"banana"`),
			expectedUF: CUF(""),
			isValid:    false,
		},
		{
			data:       []byte(`"30"`),
			expectedUF: CUF(""),
			isValid:    false,
		},
		{
			data:       []byte(`""`),
			expectedUF: CUF(""),
			isValid:    false,
		},
	}
	for _, test := range tests {
		var actual CUF
		err := actual.UnmarshalJSON(test.data)
		if test.isValid {
			assert.NoError(t, err, "Valid cUF produced error")
		} else {
			assert.Error(t, err, "Invalid cUF didnt produce error")
		}
		assert.Equal(t, actual, test.expectedUF, "Expected UF mismatch")
	}
}
