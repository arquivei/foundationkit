package avroutil

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myStruct struct {
	X int
}

func TestMockAddBinding(t *testing.T) {
	decoder := new(MockDecoder)
	decoder.AddBinding([]byte("int"), 1, nil)
	decoder.AddBinding([]byte("string"), "str", nil)
	decoder.AddBinding([]byte("struct"), myStruct{X: 1}, nil)
	decoder.AddBinding([]byte("error"), nil, errors.New("err"))
	var outputInt int
	var outputString string
	var outputStruct myStruct

	_ = decoder.Decode(context.Background(), []byte("int"), &outputInt)
	_ = decoder.Decode(context.Background(), []byte("string"), &outputString)
	_ = decoder.Decode(context.Background(), []byte("struct"), &outputStruct)

	assert.Equal(t, 1, outputInt)
	assert.Equal(t, "str", outputString)
	assert.Equal(t, myStruct{X: 1}, outputStruct)

	assert.Error(t, decoder.Decode(context.Background(), []byte("error"), &outputInt))

	decoder.AssertExpectations(t)
}

func TestMockAddJSON(t *testing.T) {
	decoder := new(MockDecoder)
	decoder.AddJSONBinding([]byte("int"), "1", nil)
	decoder.AddJSONBinding([]byte("string"), "\"str\"", nil)
	decoder.AddJSONBinding([]byte("struct"), "{\"X\": 1}", nil)
	decoder.AddJSONBinding([]byte("error"), "", errors.New("err"))
	var outputInt int
	var outputString string
	var outputStruct myStruct

	_ = decoder.Decode(context.Background(), []byte("int"), &outputInt)
	_ = decoder.Decode(context.Background(), []byte("string"), &outputString)
	_ = decoder.Decode(context.Background(), []byte("struct"), &outputStruct)

	assert.Equal(t, 1, outputInt)
	assert.Equal(t, "str", outputString)
	assert.Equal(t, myStruct{X: 1}, outputStruct)

	assert.Error(t, decoder.Decode(context.Background(), []byte("error"), &outputInt))

	decoder.AssertExpectations(t)
}
