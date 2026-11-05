package implschemaregistry

import (
	"context"
	"strconv"
	"sync"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"golang.org/x/sync/singleflight"

	"github.com/arquivei/avro/v2"
)

type cacheRepository struct {
	next            schemaregistry.Repository
	schemaByIDCache map[schemaregistry.ID]avro.Schema
	idBySchemaCache map[string]schemaregistry.ID
	lock            sync.RWMutex
	schemaByIdGroup singleflight.Group
	idBySchemaGroup singleflight.Group
}

type idBySchemaResult struct {
	id     schemaregistry.ID
	schema avro.Schema
}

// WrapWithCache wraps @next with a schemaByIDCache layer that stores the result indefinitely
func WrapWithCache(next schemaregistry.Repository) schemaregistry.Repository {
	return &cacheRepository{
		next:            next,
		lock:            sync.RWMutex{},
		schemaByIDCache: map[schemaregistry.ID]avro.Schema{},
		idBySchemaCache: map[string]schemaregistry.ID{},
	}
}

func (r *cacheRepository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (avro.Schema, error) {
	const op = errors.Op("implschemaregistry.cacheRepository.GetSchemaById")

	if schema, ok := r.tryGetSchemaFromCache(id); ok {
		return schema, nil
	}

	v, err, _ := r.schemaByIdGroup.Do(strconv.FormatUint(uint64(id), 10), func() (any, error) {
		if schema, ok := r.tryGetSchemaFromCache(id); ok {
			return schema, nil
		}

		schema, err := r.next.GetSchemaByID(context.WithoutCancel(ctx), id)
		if err != nil {
			return nil, err
		}

		r.storeSchemaByID(id, schema)

		return schema, nil
	})

	if err != nil {
		return nil, errors.E(op, err)
	}

	schema, _ := v.(avro.Schema)

	return schema, nil
}

func (r *cacheRepository) tryGetSchemaFromCache(id schemaregistry.ID) (avro.Schema, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	schema, ok := r.schemaByIDCache[id]
	return schema, ok
}

func (r *cacheRepository) tryGetIDFromSchemaCache(schema string) (schemaregistry.ID, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	id, ok := r.idBySchemaCache[schema]
	return id, ok
}

func (r *cacheRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (schemaregistry.ID, avro.Schema, error) {
	const op = errors.Op("implschemaregistry.cacheRepository.GetIDBySchema")

	if id, avroSchema, ok := r.tryGetIDAndSchemaFromCache(schema); ok {
		return id, avroSchema, nil
	}

	v, err, _ := r.idBySchemaGroup.Do(string(subject)+"\x00"+schema, func() (any, error) {
		if id, avroSchema, ok := r.tryGetIDAndSchemaFromCache(schema); ok {
			return idBySchemaResult{id: id, schema: avroSchema}, nil
		}

		id, avroSchema, err := r.next.GetIDBySchema(context.WithoutCancel(ctx), subject, schema)
		if err != nil {
			return idBySchemaResult{id: id}, err
		}

		r.storeIDBySchema(id, schema)
		r.storeSchemaByID(id, avroSchema)

		return idBySchemaResult{id: id, schema: avroSchema}, nil
	})

	result, _ := v.(idBySchemaResult)

	if err != nil {
		return result.id, nil, errors.E(op, err)
	}

	return result.id, result.schema, nil
}

func (r *cacheRepository) tryGetIDAndSchemaFromCache(schema string) (schemaregistry.ID, avro.Schema, bool) {
	id, ok := r.tryGetIDFromSchemaCache(schema)
	if !ok {
		return 0, nil, false
	}

	avroSchema, ok := r.tryGetSchemaFromCache(id)

	return id, avroSchema, ok
}

func (r *cacheRepository) storeSchemaByID(id schemaregistry.ID, schema avro.Schema) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.schemaByIDCache[id] = schema
}

func (r *cacheRepository) storeIDBySchema(id schemaregistry.ID, schema string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.idBySchemaCache[schema] = id
}
