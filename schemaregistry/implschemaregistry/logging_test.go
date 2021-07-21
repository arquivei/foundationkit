package implschemaregistry

import (
	"context"
	"errors"
	"testing"

	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/stretchr/testify/assert"
)

func TestLoggingGetSchemaByID(t *testing.T) {
	schemaregistryID := schemaregistry.ID(10)

	mock := new(mockTestRepository)
	mock.
		On("GetSchemaByID", schemaregistryID).
		Return(tagsSchema, nil).
		Once()

	repository := WrapWithLogging(mock)

	schema, err := repository.GetSchemaByID(
		context.Background(),
		schemaregistryID,
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchema.String(), schema.String())
}

func TestLoggingGetSchemaByID_Error(t *testing.T) {
	schemaregistryID := schemaregistry.ID(10)

	mock := new(mockTestRepository)
	mock.
		On("GetSchemaByID", schemaregistryID).
		Return(tagsSchema, errors.New("error")).
		Once()

	repository := WrapWithLogging(mock)

	schema, err := repository.GetSchemaByID(
		context.Background(),
		schemaregistryID,
	)
	assert.EqualError(t, err, "error")
	assert.Equal(t, tagsSchema.String(), schema.String())
}

func TestLoggingGetIDBySchema(t *testing.T) {
	schemaregistrySubject := schemaregistry.Subject("mysubject")
	schemaregistryID := schemaregistry.ID(10)

	mock := new(mockTestRepository)
	mock.
		On("GetIDBySchema", schemaregistrySubject, tagsSchema.String()).
		Return(schemaregistryID, tagsSchema, nil).
		Once()

	repository := WrapWithLogging(mock)

	id, schema, err := repository.GetIDBySchema(
		context.Background(),
		schemaregistrySubject,
		tagsSchema.String(),
	)
	assert.NoError(t, err)
	assert.Equal(t, tagsSchema.String(), schema.String())
	assert.Equal(t, schemaregistryID, id)
}

func TestLoggingGetIDBySchema_Error(t *testing.T) {
	schemaregistrySubject := schemaregistry.Subject("mysubject")
	schemaregistryID := schemaregistry.ID(10)

	mock := new(mockTestRepository)
	mock.
		On("GetIDBySchema", schemaregistrySubject, tagsSchema.String()).
		Return(schemaregistryID, tagsSchema, errors.New("error")).
		Once()

	repository := WrapWithLogging(mock)

	id, schema, err := repository.GetIDBySchema(
		context.Background(),
		schemaregistrySubject,
		tagsSchema.String(),
	)
	assert.EqualError(t, err, "error")
	assert.Equal(t, tagsSchema.String(), schema.String())
	assert.Equal(t, schemaregistryID, id)
}
