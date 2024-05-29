package implschemaregistry

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro/v2"
)

type mockRepository struct {
	idToSchemas  map[schemaregistry.ID]avro.Schema
	schemasToIDs map[string]schemaregistry.ID
}

// MustNewMock creates a mock schema registry.
func MustNewMock(schemas map[schemaregistry.ID]string) schemaregistry.Repository {
	const op = errors.Op("implschemaregistry.MustNewMock")
	r := mockRepository{
		idToSchemas:  make(map[schemaregistry.ID]avro.Schema),
		schemasToIDs: make(map[string]schemaregistry.ID),
	}

	for id, schemaStr := range schemas {
		schema, err := avro.Parse(schemaStr)
		if err != nil {
			panic(errors.E(
				err,
				op,
				errors.KV("schema", truncateStr(schemaStr, 50)),
			))
		}
		r.idToSchemas[id] = schema
		r.schemasToIDs[schema.String()] = id
	}

	return r
}

func (r mockRepository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (avro.Schema, error) {
	const op = errors.Op("implschemaregistry.mockRepository.GetSchemaByID")

	if schema, ok := r.idToSchemas[id]; ok {
		return schema, nil
	}

	return nil, errors.New("could not find schema", errors.KV("id", id), op)
}

func (r mockRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	const op = errors.Op("implschemaregistry.mockRepository.GetIDBySchema")

	avroSchema, err := avro.Parse(schema)
	if err != nil {
		return 0, nil, errors.E(err, op)
	}

	if id, ok := r.schemasToIDs[avroSchema.String()]; ok {
		return id, avroSchema, nil
	}

	return 0, nil, errors.New("could not find schema", errors.KV("subject", subject), op)
}

func truncateStr(str string, size int) string {
	if len(str) > size {
		return str[0:size]
	}
	return str
}
