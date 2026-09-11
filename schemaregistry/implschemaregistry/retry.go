package implschemaregistry

import (
	"context"

	"github.com/arquivei/avro/v2"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/retrier"
	"github.com/arquivei/foundationkit/schemaregistry"
)

type retryRepository struct {
	next    schemaregistry.Repository
	retrier *retrier.Retrier
}

// WrapWithRetry wraps @next with a retry layer, configured according to @settings.
// See retrier.Settings for the defaults assumed on zero-valued fields.
func WrapWithRetry(next schemaregistry.Repository, settings retrier.Settings) schemaregistry.Repository {
	return &retryRepository{
		next:    next,
		retrier: retrier.NewRetrier(settings),
	}
}

func (r *retryRepository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (schema avro.Schema, err error) {
	const op = errors.Op("implschemaregistry.retryRepository.GetSchemaByID")

	err = r.retrier.ExecuteOperation(func() error {
		var innerErr error
		schema, innerErr = r.next.GetSchemaByID(ctx, id)
		return innerErr
	})
	if err != nil {
		return nil, errors.E(op, err)
	}

	return schema, nil
}

func (r *retryRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (id schemaregistry.ID, avroSchema avro.Schema, err error) {
	const op = errors.Op("implschemaregistry.retryRepository.GetIDBySchema")

	err = r.retrier.ExecuteOperation(func() error {
		var innerErr error
		id, avroSchema, innerErr = r.next.GetIDBySchema(ctx, subject, schema)
		return innerErr
	})
	if err != nil {
		return 0, nil, errors.E(op, err)
	}

	return id, avroSchema, nil
}
