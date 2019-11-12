package cuf

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUF(t *testing.T) {
	tests := []struct {
		str    string
		errMsg string
	}{
		{str: "11"}, {str: "12"}, {str: "13"}, {str: "14"}, {str: "15"},
		{str: "16"}, {str: "17"}, {str: "21"}, {str: "22"}, {str: "23"},
		{str: "24"}, {str: "25"}, {str: "26"}, {str: "27"}, {str: "28"},
		{str: "29"}, {str: "31"}, {str: "32"}, {str: "33"}, {str: "35"},
		{str: "41"}, {str: "42"}, {str: "43"}, {str: "50"}, {str: "51"},
		{str: "52"}, {str: "53"},
		{str: "", errMsg: "input cUF should have 2 digits: "},
		{str: "000", errMsg: "input cUF should have 2 digits: 000"},
		{str: "035", errMsg: "input cUF should have 2 digits: 035"},
		{str: "0", errMsg: "input cUF should have 2 digits: 0"},
		{str: "00", errMsg: "invalid cUF code: 00"},
		{str: "1", errMsg: "input cUF should have 2 digits: 1"},
		{str: "9", errMsg: "input cUF should have 2 digits: 9"},
		{str: "10", errMsg: "invalid cUF code: 10"},
		{str: "18", errMsg: "invalid cUF code: 18"},
		{str: "19", errMsg: "invalid cUF code: 19"},
		{str: "20", errMsg: "invalid cUF code: 20"},
		{str: "30", errMsg: "invalid cUF code: 30"},
		{str: "34", errMsg: "invalid cUF code: 34"},
		{str: "36", errMsg: "invalid cUF code: 36"},
		{str: "40", errMsg: "invalid cUF code: 40"},
		{str: "45", errMsg: "invalid cUF code: 45"},
		{str: "49", errMsg: "invalid cUF code: 49"},
		{str: "54", errMsg: "invalid cUF code: 54"},
		{str: "55", errMsg: "invalid cUF code: 55"},
		{str: "56", errMsg: "invalid cUF code: 56"},
		{str: "57", errMsg: "invalid cUF code: 57"},
		{str: "foo", errMsg: "input cUF should have 2 digits: foo"},
		{str: "bar", errMsg: "input cUF should have 2 digits: bar"},
	}
	for _, test := range tests {
		uf, err := parseUF(test.str)
		if test.errMsg == "" {
			assert.NoError(t, err)
			assert.Equal(t, test.str, strconv.Itoa(int(uf)))

		} else {
			assert.EqualError(t, err, test.errMsg)
		}

	}
}

func TestCUFString(t *testing.T) {
	cuf, err := New("35")
	assert.NoError(t, err)
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
			expected: MustNew("35"),
			err:      "",
		},
		{
			name:     "Wrong UF code",
			input:    "00",
			expected: CUF{},
			err:      "cuf.New: invalid cUF code: 00",
		},
		{
			name:     "Not number cUF string",
			input:    "vc",
			expected: CUF{},
			err:      "cuf.New: cUF could not be converted to integer: vc",
		},
		{
			name:     "Not number cUF string",
			input:    "bar",
			expected: CUF{},
			err:      "cuf.New: input cUF should have 2 digits: bar",
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
			expectedUF: MustNew("31"),
			isValid:    true,
		},
		{
			data:       []byte(`"banana"`),
			expectedUF: CUF{},
			isValid:    false,
		},
		{
			data:       []byte(`"30"`),
			expectedUF: CUF{},
			isValid:    false,
		},
		{
			data:       []byte(`""`),
			expectedUF: CUF{},
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

func TestIsValid(t *testing.T) {
	cuf := MustNew("35")
	assert.True(t, IsValid(cuf))
	assert.False(t, IsValid(CUF{}))
}

func TestMustNew(t *testing.T) {
	assert.Panics(t, func() { MustNew("blabla") })
	assert.Equal(t, MustNew("35"), CUF{true, 35})
}

func TestString(t *testing.T) {
	assert.Panics(t, func() {
		d := CUF{}
		_ = d.String()
	}, "CUF not initialized")
}

func TestMarshalJSON_WithError(t *testing.T) {
	cuf := CUF{}
	_, err := cuf.MarshalJSON()
	assert.EqualError(t, err, "CUF.MarshalJSON: CUF not initialized")
}

func TestMarshalJSON(t *testing.T) {
	cuf := MustNew("35")
	r, _ := cuf.MarshalJSON()
	assert.Equal(t, []byte(`"35"`), r)
}
