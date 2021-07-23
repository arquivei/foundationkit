package avroutil

import (
	"context"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/schemaregistry/implschemaregistry"

	"github.com/stretchr/testify/assert"
)

func TestEncoder(t *testing.T) {
	tests := []struct {
		name                string
		input               tagsType
		writerSchemaStr     string
		schemaRepository    schemaregistry.Repository
		expectedResult      []byte
		expectedNewError    string
		expectedEncodeError string
	}{
		{
			name: "Success",
			input: tagsType{
				Tags: []string{"a", "b", "c"},
			},
			writerSchemaStr: tagsSchemaStr,
			schemaRepository: implschemaregistry.MustNewMock(map[schemaregistry.ID]string{
				tagsSchemaID: tagsSchemaStr,
			}),
			expectedResult: []byte{0, 0, 0, 0, 1, 5, 12, 2, 97, 2, 98, 2, 99, 0},
		},
		{
			name: "Error - failed to get schema ID",
			input: tagsType{
				Tags: []string{"a", "b", "c"},
			},
			writerSchemaStr:  tagsSchemaStr,
			schemaRepository: implschemaregistry.MustNewMock(map[schemaregistry.ID]string{}),
			expectedNewError: "could not find schema [subject=subject]",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoder, err := NewEncoder(context.Background(), test.schemaRepository, "subject", test.writerSchemaStr)
			if test.expectedNewError != "" {
				assert.EqualError(t, errors.GetRootErrorWithKV(err), test.expectedNewError)
				return
			}
			assert.NoError(t, err)
			result, err := encoder.Encode(test.input)
			if test.expectedEncodeError == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			} else {
				assert.EqualError(t, errors.GetRootErrorWithKV(err), test.expectedEncodeError)
			}
		})
	}
}
