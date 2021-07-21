package implschemaregistry

import (
	"context"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/stretchr/testify/assert"
)

func TestMockRepositoryGetIDBySchema(t *testing.T) {
	mock := MustNewMock(map[schemaregistry.ID]string{
		1: tagsSchemaStr,
	})

	id, schema, err := mock.GetIDBySchema(context.Background(), "mysubject", tagsSchemaStr)
	assert.NoError(t, err)
	assert.Equal(t, schemaregistry.ID(1), id)
	assert.Equal(t, tagsSchemaStr, schema.String())

	_, _, err = mock.GetIDBySchema(context.Background(), "mysubject", "wrongschema")
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "avro: unknown type: wrongschema")

	_, _, err = mock.GetIDBySchema(context.Background(), "mysubject", `{"name":"b","type":"string"}`)
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "could not find schema [subject=mysubject]")
}

func TestMockRepositoryGetSchemaByID(t *testing.T) {
	mock := MustNewMock(map[schemaregistry.ID]string{
		1: tagsSchemaStr,
	})

	schema, err := mock.GetSchemaByID(context.Background(), schemaregistry.ID(1))
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())

	_, err = mock.GetSchemaByID(context.Background(), schemaregistry.ID(2))
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "could not find schema [id=2]")
}

func TestMustNewMock_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustNewMock(map[schemaregistry.ID]string{1: "ops"})
	})
}
