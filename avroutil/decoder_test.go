package avroutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/schemaregistry/implschemaregistry"

	"github.com/stretchr/testify/assert"
)

const (
	tagsSchemaID schemaregistry.ID = 1
)

var (
	tagsSchemaStr = string(loadFile("schemas/Tags.avsc"))
)

type tagsType struct {
	Tags []string `avro:"Tags"`
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name           string
		input          []byte
		expectedResult tagsType
		expectedError  string
	}{
		{
			name:  "Success",
			input: []byte{0, 0, 0, 0, 1, 5, 12, 2, 97, 2, 98, 2, 99, 0},
			expectedResult: tagsType{
				Tags: []string{"a", "b", "c"},
			},
			expectedError: "",
		},
		{
			name:          "Error - Failed to split",
			input:         []byte{0},
			expectedError: "invalid message length",
		},
		{
			name:          "Error - Schema ID not found",
			input:         []byte{0, 0, 0, 0, 99, 5, 12, 2, 97, 2, 98, 2, 99, 0},
			expectedError: "could not find schema [id=99]",
		},
	}
	schemaRegistry := implschemaregistry.MustNewMock(map[schemaregistry.ID]string{
		tagsSchemaID: tagsSchemaStr,
	})
	decoder := NewDecoder(schemaRegistry)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tags tagsType
			err := decoder.Decode(context.Background(), test.input, &tags)
			if test.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, tags)
			} else {
				assert.EqualError(t, errors.GetRootErrorWithKV(err), test.expectedError)
			}
		})
	}
}

func loadFile(name string) []byte {
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		panic(err)
	}
	return content
}

func TestSplitAvroWireFormatMessage(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		expectedID   schemaregistry.ID
		expectedData []byte
	}{
		{
			"zeros",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			schemaregistry.ID(0),
			[]byte{0x00},
		},
		{
			"one",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0xFF},
			schemaregistry.ID(1),
			[]byte{0xFF},
		},
		{
			"bignumber",
			[]byte{0x00, 0x0F, 0x00, 0x00, 0x00, 0xAB},
			schemaregistry.ID(0x0F000000),
			[]byte{0xAB},
		},
		{
			"empty",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x01},
			schemaregistry.ID(0x01),
			[]byte{},
		},
	}
	for _, test := range tests {
		id, data, err := SplitAvroWireFormatMessage(test.data)
		assert.NoError(t, err, test.name)
		assert.Equal(t, test.expectedID, id, test.name)
		assert.Equal(t, test.expectedData, data, test.name)
	}
}

func TestSplitAvroWireFormatMessageError(t *testing.T) {
	_, _, err := SplitAvroWireFormatMessage([]byte{0x00, 0x00, 0x00, 0x00})
	assert.Error(t, err)

	_, _, err = SplitAvroWireFormatMessage([]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00})
	assert.Error(t, err)
}
