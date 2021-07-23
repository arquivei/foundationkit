package implschemaregistry

import (
	"context"
	"sync"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro"
)

type cacheRepository struct {
	next  schemaregistry.Repository
	cache map[schemaregistry.ID]avro.Schema
	lock  sync.RWMutex
}

// WrapWithCache wraps @next with a cache layer that stores the result indefinitely
func WrapWithCache(next schemaregistry.Repository) schemaregistry.Repository {
	return &cacheRepository{
		next:  next,
		lock:  sync.RWMutex{},
		cache: map[schemaregistry.ID]avro.Schema{},
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

	r.storeResult(id, schema)

	return schema, nil
}

func (r *cacheRepository) tryGetSchemaFromCache(id schemaregistry.ID) (avro.Schema, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	schema, ok := r.cache[id]
	return schema, ok
}

func (r *cacheRepository) storeResult(id schemaregistry.ID, schema avro.Schema) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.cache[id] = schema
}

func (r *cacheRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	return r.next.GetIDBySchema(ctx, subject, schema)
}
