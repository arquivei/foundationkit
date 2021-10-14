package implschemaregistry

import (
	"context"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCacheGetSchemaByID(t *testing.T) {
	mock := new(mockTestRepository)
	mock.On("GetSchemaByID", schemaregistry.ID(10)).
		Return(tagsSchema, nil).
		Once()
	mock.On("GetSchemaByID", schemaregistry.ID(99)).
		Return(tagsSchema, errors.New("not found")).
		Once()

	repository := WrapWithCache(mock)
	// First time, fetch from repository
	schema, err := repository.GetSchemaByID(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())

	// Second time, get from schemaByIDCache
	schema, err = repository.GetSchemaByID(context.Background(), 10)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())

	// Repository returned error
	_, err = repository.GetSchemaByID(context.Background(), 99)
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "not found")
}

func TestCacheGetIDBySchema(t *testing.T) {
	mock := new(mockTestRepository)
	mock.
		On("GetIDBySchema", schemaregistry.Subject("mysubject"), tagsSchemaStr).
		Return(schemaregistry.ID(10), tagsSchema, nil).
		Twice()
	mock.
		On("GetIDBySchema", schemaregistry.Subject("mysubject"), "wrongschema").
		Return(schemaregistry.ID(10), tagsSchema, errors.New("not found")).
		Once()

	repository := WrapWithCache(mock)
	// First time, fetch from repository
	id, schema, err := repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		tagsSchemaStr,
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())
	assert.Equal(t, schemaregistry.ID(10), id)

	// Second time, fetch from repository. CACHE IS NOT IMPLEMENTED
	id, schema, err = repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		tagsSchemaStr,
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchemaStr, schema.String())
	assert.Equal(t, schemaregistry.ID(10), id)

	// Repository returned error
	_, _, err = repository.GetIDBySchema(
		context.Background(),
		schemaregistry.Subject("mysubject"),
		"wrongschema",
	)
	assert.EqualError(t, errors.GetRootErrorWithKV(err), "not found")
}

type mockTestRepository struct {
	mock.Mock
}

func (m *mockTestRepository) GetSchemaByID(
	ctx context.Context,
	id schemaregistry.ID,
) (avro.Schema, error) {
	args := m.Called(id)
	return args.Get(0).(avro.Schema), args.Error(1)
}
func (m *mockTestRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	args := m.Called(subject, schema)
	return args.Get(0).(schemaregistry.ID),
		args.Get(1).(avro.Schema), args.Error(2)
}
