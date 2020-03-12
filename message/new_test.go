package message

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetTypeName(t *testing.T) {
	type test1 struct{}
	d := test1{}
	assert.Equal(t, getTypeName(d), "test1")
	dPtr := &test1{}

	assert.Equal(t, getTypeName(dPtr), "test1")
}

func TestParseTypeAndDataVersion(t *testing.T) {
	// good
	type myMessageV1 struct{}
	type myMessageV2 struct{}
	type MyMessageV1 struct{}
	type myHTTPMessageV1 struct{}
	type HTTPMessageV1 struct{}
	type myMessageHTTPV1 struct{}
	// bad
	type myMessagev1 struct{}
	type myMessage struct{}
	type myMessageV struct{}
	type myMessageWithouVersion struct{}
	type V struct{}
	type v struct{}
	type V1 struct{}

	tests := []struct {
		name            string
		input           interface{}
		expectedType    Type
		expectedVersion DataVersion
		expectedErr     string
	}{
		{
			name:            "Private type with version 1",
			input:           myMessageV1{},
			expectedType:    Type("my-message"),
			expectedVersion: DataVersion(1),
			expectedErr:     "",
		},
		{
			name:            "Private type with version 2",
			input:           myMessageV2{},
			expectedType:    Type("my-message"),
			expectedVersion: DataVersion(2),
			expectedErr:     "",
		},
		{
			name:            "Public type with version 1",
			input:           MyMessageV1{},
			expectedType:    Type("my-message"),
			expectedVersion: DataVersion(1),
			expectedErr:     "",
		},
		{
			name:            "Private type with upper inner word and version 1",
			input:           myHTTPMessageV1{},
			expectedType:    Type("my-http-message"),
			expectedVersion: DataVersion(1),
			expectedErr:     "",
		},
		{
			name:            "Private type with upper word at begin and version 1",
			input:           HTTPMessageV1{},
			expectedType:    Type("http-message"),
			expectedVersion: DataVersion(1),
			expectedErr:     "",
		},
		{
			name:            "Private type with upper word at end and version 1",
			input:           myMessageHTTPV1{},
			expectedType:    Type("my-message-http"),
			expectedVersion: DataVersion(1),
			expectedErr:     "",
		},
		{
			name:            "With no version",
			input:           myMessage{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'myMessage'",
		},
		{
			name:            "With empty version",
			input:           myMessageV{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'myMessageV'",
		},
		{
			name:            "With broken version",
			input:           myMessageWithouVersion{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: strconv.Atoi: parsing \"ersion\": invalid syntax",
		},
		{
			name:            "With lower v in version",
			input:           myMessagev1{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'myMessagev1'",
		},

		{
			name:            "With no type",
			input:           V1{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'V1'",
		},
		{
			name:            "Only V",
			input:           V{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'V'",
		},
		{
			name:            "Only v",
			input:           v{},
			expectedType:    Type(""),
			expectedVersion: DataVersion(0),
			expectedErr:     "message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'v'",
		},
	}
	for _, test := range tests {
		actualType, actualVersion, err := ParseTypeAndDataVersion(test.input)
		if test.expectedErr != "" {
			assert.EqualError(t, err, test.expectedErr, "[%s] Error mismatch", test.name)
		} else {
			assert.NoError(t, err, "[%s] Error not expected", test.name)
		}
		assert.Equal(t, test.expectedType, actualType, "[%s] Type mismatch", test.name)
		assert.Equal(t, test.expectedVersion, actualVersion, "[%s] Version mismatch", test.name)
	}
}

type myMessageMockV1 struct {
	Msg string
	err error
}

func (m myMessageMockV1) MarshalJSON() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return json.Marshal(m.Msg)
}

func TestNewWithAutoType(t *testing.T) {
	type myBadType struct{}

	tests := []struct {
		name        string
		input       interface{}
		expectedErr string
	}{
		{
			name:        "Good message",
			input:       myMessageMockV1{Msg: "hello"},
			expectedErr: "",
		},
		{
			name:        "Unmarshable message",
			input:       myMessageMockV1{Msg: "hello", err: errors.New("my err")},
			expectedErr: "message.New: json: error calling MarshalJSON for type message.myMessageMockV1: my err",
		},
		{
			name:        "With bad message type",
			input:       myBadType{},
			expectedErr: "message.New: message.ParseTypeAndDataVersion: invalid type name, expected '<type>V<version>' but got 'myBadType'",
		},
	}

	ctx := context.TODO()
	for _, test := range tests {
		actual, err := New(ctx, Source("foundationkit-test"), test.input)
		if test.expectedErr != "" {
			assert.EqualError(t, err, test.expectedErr, "[%s] Error mismatch", test.name)
		} else {
			assert.NoError(t, err, "[%s] Error not expected", test.name)
			assert.NotEmpty(t, actual, "[%s] Message is empty")
		}
	}
}
