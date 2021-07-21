package avroutil

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/stretchr/testify/mock"
)

// MockDecoder is a special implementation of Decoder, useful for
// unit testing
type MockDecoder struct {
	mock.Mock
}

// Decode decodes @data, in the wire format, into @output.
func (i *MockDecoder) Decode(ctx context.Context, data []byte, output interface{}) error {
	args := i.Called(data)
	err := args.Error(1)
	if err == nil {
		i.handleOutput(args, output)
	}
	return err
}

func (i *MockDecoder) handleOutput(args mock.Arguments, output interface{}) {
	isJSON := args.Bool(2)
	if isJSON {
		jsonMock := args.Get(0).(string)

		err := json.Unmarshal([]byte(jsonMock), output)
		if err != nil {
			panic(err)
		}
	} else {
		i.setMockToOutput(output, args.Get(0))
	}
}

func (i *MockDecoder) setMockToOutput(output, mocked interface{}) {
	rv1 := reflect.ValueOf(output)
	if rv1.Kind() != reflect.Ptr || rv1.IsNil() {
		panic("invalid output type")
	}
	rv2 := reflect.ValueOf(mocked)
	rv1.Elem().Set(rv2)
}

// AddBinding mocks the decoding result of @data to the @output struct
func (i *MockDecoder) AddBinding(data []byte, output interface{}, err error) *mock.Call {
	return i.On("Decode", data).Return(output, err, false)
}

// AddJSONBinding mocks the decoding of @data with the JSON-encoded data in
// @output. Use this function when the output structure is hard to define in
// unit tests, ex: anonymous structs with lots of pointers.
func (i *MockDecoder) AddJSONBinding(data []byte, output string, err error) *mock.Call {
	return i.On("Decode", data).Return(output, err, true)
}
