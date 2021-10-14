package implschemaregistry

import (
	"context"
	"sync"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro"
)

type cacheRepository struct {
	next                schemaregistry.Repository
	schemaByIDCache     map[schemaregistry.ID]avro.Schema
	idBySchemaCache     map[string]schemaregistry.ID
	lockSchemaByIDCache sync.RWMutex
	lockIDBySchemaCache sync.RWMutex
}

// WrapWithCache wraps @next with a schemaByIDCache layer that stores the result indefinitely
func WrapWithCache(next schemaregistry.Repository) schemaregistry.Repository {
	return &cacheRepository{
		next:                next,
		lockSchemaByIDCache: sync.RWMutex{},
		lockIDBySchemaCache: sync.RWMutex{},
		schemaByIDCache:     map[schemaregistry.ID]avro.Schema{},
		idBySchemaCache:     map[string]schemaregistry.ID{},
	}
}

func (r *cacheRepository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (avro.Schema, error) {
	const op = errors.Op("implschemaregistry.cacheRepository.GetSchemaById")

	if schema, ok := r.tryGetSchemaFromCache(id); ok {
		return schema, nil
	}

	schema, err := r.next.GetSchemaByID(ctx, id)
	if err != nil {
		return nil, errors.E(op, err)
	}

	r.storeSchemaByID(id, schema)

	return schema, nil
}

func (r *cacheRepository) tryGetSchemaFromCache(id schemaregistry.ID) (avro.Schema, bool) {
	r.lockSchemaByIDCache.RLock()
	defer r.lockSchemaByIDCache.RUnlock()

	schema, ok := r.schemaByIDCache[id]
	return schema, ok
}
func (r *cacheRepository) tryGetIDFromSchemaCache(schema string) (schemaregistry.ID, bool) {
	r.lockSchemaByIDCache.RLock()
	defer r.lockSchemaByIDCache.RUnlock()

	id, ok := r.idBySchemaCache[schema]
	return id, ok
}

func (r *cacheRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	const op = errors.Op("implschemaregistry.cacheRepository.GetIDBySchema")

	if id, ok := r.tryGetIDFromSchemaCache(schema); ok {
		if avroSchema, ok := r.tryGetSchemaFromCache(id); ok {
			return id, avroSchema, nil
		}
	}

	id, avroSchema, err := r.next.GetIDBySchema(ctx, subject, schema)
	if err != nil {
		return 0, nil, errors.E(op, err)
	}

	r.storeIDBySchemaCache(id, schema)
	r.storeSchemaByID(id, avroSchema)

	return id, avroSchema, nil
}

func (r *cacheRepository) storeSchemaByID(id schemaregistry.ID, schema avro.Schema) {
	r.lockIDBySchemaCache.Lock()
	defer r.lockIDBySchemaCache.Unlock()

	r.schemaByIDCache[id] = schema
}

func (r *cacheRepository) storeIDBySchemaCache(id schemaregistry.ID, schema string) {
	r.lockIDBySchemaCache.Lock()
	defer r.lockIDBySchemaCache.Unlock()

	r.idBySchemaCache[schema] = id
}
