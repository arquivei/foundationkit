package nsu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNSUStringMethod(t *testing.T) {
	assert.Equal(t, "000000000000000", NSU("0").String())
	assert.Equal(t, "000000000000123", NSU("123").String())
	assert.Equal(t, "123456789012345", NSU("123456789012345").String())
}

func TestNSUJSONMarshaler(t *testing.T) {
	n := NSU("123")
	j, err := n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"000000000000123"`), j)
}

func TestNSUJSONUnmarshaler(t *testing.T) {
	testcases := []struct {
		Test     string
		Input    []byte
		Expected NSU
		Err      string
	}{
		{
			Test:     "Valid NSU",
			Input:    []byte(`"123"`),
			Expected: NSU("000000000000123"),
			Err:      "",
		},
		{
			Test:     "Invalid JSON",
			Input:    []byte(`"123`),
			Expected: NSU(""),
			Err:      "nsu.UnmarshalJSON: unexpected end of JSON input",
		},
		{
			Test:     "Invalid NSU",
			Input:    []byte(`"abc"`),
			Expected: NSU(""),
			Err:      "nsu.UnmarshalJSON: nsu.Parse: failed to parse nsu from string",
		},
	}

	for _, testcase := range testcases {
		var output NSU
		err := output.UnmarshalJSON(testcase.Input)
		assert.Equal(t, testcase.Expected, output, "[%s] Unexpected output", testcase.Test)
		if testcase.Err != "" {
			assert.EqualError(t, err, testcase.Err, "[%s] Unexpected error", testcase.Test)
		} else {
			assert.NoError(t, err, "[%s] Unexpected error", testcase.Test)
		}
	}

}

func TestParse(t *testing.T) {
	testcases := []struct {
		Test     string
		Input    string
		Expected NSU
		Err      string
	}{
		{
			Test:     "Valid NSU",
			Input:    "123",
			Expected: NSU("000000000000123"),
			Err:      "",
		},
		{
			Test:     "Empty NSU",
			Input:    "",
			Expected: "",
			Err:      "nsu.Parse: nsu is empty",
		},
		{
			Test:     "NSU too long",
			Input:    "0000000000001234",
			Expected: "",
			Err:      "nsu.Parse: nsu has more than 15 digits",
		},
		{
			Test:     "Invalid NSU",
			Input:    "abc",
			Expected: "",
			Err:      "nsu.Parse: failed to parse nsu from string",
		},
	}

	for _, testcase := range testcases {
		output, err := Parse(testcase.Input)
		assert.Equal(t, testcase.Expected, output, "[%s] Unexpected output", testcase.Test)
		if testcase.Err != "" {
			assert.EqualError(t, err, testcase.Err, "[%s] Unexpected error", testcase.Test)
		} else {
			assert.NoError(t, err, "[%s] Unexpected error", testcase.Test)
		}
	}
}

func TestCompare(t *testing.T) {
	assert.Equal(t, 0, Compare(NSU("000000000000123"), NSU("000000000000123")))
	assert.Equal(t, 1, Compare(NSU("000000000000123"), NSU("000000000000122")))
}
