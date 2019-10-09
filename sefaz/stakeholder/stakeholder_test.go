package stakeholder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStakeholderString(t *testing.T) {
	s := Stakeholder("15380666450")
	assert.Equal(t, "15380666450", s.String())
}

func TestStakeholderJSONMarshaler(t *testing.T) {
	s := Stakeholder("15380666450")
	j, err := s.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"15380666450"`), j)
}

func TestStakeholderJSONUnmarshaler(t *testing.T) {
	testcases := []struct {
		Test     string
		Input    []byte
		Expected Stakeholder
		Err      string
	}{
		{
			Test:     "Valid CNPJ stakeholder",
			Input:    []byte(`"49205898000123"`),
			Expected: Stakeholder("49205898000123"),
			Err:      "",
		},
		{
			Test:     "Valid CPJ stakeholder",
			Input:    []byte(`"15380666450"`),
			Expected: Stakeholder("15380666450"),
			Err:      "",
		},
		{
			Test:     "Invalid JSON",
			Input:    []byte(`"15380666450`),
			Expected: Stakeholder(""),
			Err:      "stakeholder.UnmarshalJSON: unexpected end of JSON input",
		},
		{
			Test:     "Invalid cnpj",
			Input:    []byte(`"abc"`),
			Expected: Stakeholder(""),
			Err:      "stakeholder.UnmarshalJSON: stakeholder.Parse: invalid stakeholder type: unkown",
		},
	}

	for _, testcase := range testcases {
		var output Stakeholder
		err := output.UnmarshalJSON(testcase.Input)
		assert.Equal(t, testcase.Expected, output, "[%s] Unexpected output", testcase.Test)
		if testcase.Err != "" {
			assert.EqualError(t, err, testcase.Err, "[%s] Unexpected error", testcase.Test)
		} else {
			assert.NoError(t, err, "[%s] Unexpected error", testcase.Test)
		}
	}

}

func TestNewCPF(t *testing.T) {
	sh, err := NewCPF("15380666450")
	assert.NoError(t, err, "success stakeholder")
	assert.Equal(t, TypePerson, GetType(sh), "success stakeholder type")

	sh, err = NewCPF("abc")
	assert.Error(t, err, "fail stakeholder")
	assert.Equal(t, TypeUnknown, GetType(sh), "fail stakeholder type")
}

func TestNewCNPJ(t *testing.T) {
	sh, err := NewCNPJ("49205898000123")
	assert.NoError(t, err, "success stakeholder")
	assert.Equal(t, TypeCompany, GetType(sh), "success stakeholder type")

	sh, err = NewCNPJ("abc")
	assert.Error(t, err, "fail stakeholder")
	assert.Equal(t, TypeUnknown, GetType(sh), "fail stakeholder type")
}

func TestGetCPFCNPJ(t *testing.T) {
	tests := []struct {
		name         string
		input        Stakeholder
		expectedCpf  string
		expectedCnpj string
	}{
		{
			name:         "Valid CPF",
			input:        Stakeholder("15380666450"),
			expectedCpf:  "15380666450",
			expectedCnpj: "",
		},
		{
			name:         "Valid CNPJ",
			input:        Stakeholder("49205898000123"),
			expectedCpf:  "",
			expectedCnpj: "49205898000123",
		},
		{
			name:         "Invalid stakeholder",
			input:        Stakeholder("abc"),
			expectedCpf:  "",
			expectedCnpj: "",
		},
	}

	for _, test := range tests {
		cpf, cnpj := GetCPFCNPJ(test.input)
		assert.Equal(t, test.expectedCpf, cpf, "[%s] CPF mismatch")
		assert.Equal(t, test.expectedCnpj, cnpj, "[%s] CNPJ mismatch")
	}
}
