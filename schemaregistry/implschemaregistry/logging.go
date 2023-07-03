package implschemaregistry

import (
	"context"

	"github.com/hamba/avro/v2"
	"github.com/rs/zerolog/log"

	"github.com/arquivei/foundationkit/schemaregistry"
)

type loggingRepository struct {
	next schemaregistry.Repository
}

// WrapWithLogging wraps @next with a logging layer
func WrapWithLogging(next schemaregistry.Repository) schemaregistry.Repository {
	return &loggingRepository{
		next: next,
	}
}

func (r *loggingRepository) GetSchemaByID(ctx context.Context, id schemaregistry.ID) (_ avro.Schema, err error) {
	defer func() {
		logger := log.Ctx(ctx)
		if err != nil {
			logger.Error().
				Err(err).
				EmbedObject(id).
				Msg("GetSchemaByID returned an error")
		} else {
			logger.Debug().
				EmbedObject(id).
				Msg("GetSchemaByID returned successfully")
		}
	}()
	return r.next.GetSchemaByID(ctx, id)
}

func (r *loggingRepository) GetIDBySchema(
	ctx context.Context,
	subject schemaregistry.Subject,
	schema string,
) (id schemaregistry.ID, _ avro.Schema, err error) {
	defer func() {
		logger := log.Ctx(ctx)
		if err != nil {
			logger.Error().
				Err(err).
				Str("subject", string(subject)).
				Msg("GetIDBySchema returned an error")
		} else {
			logger.Debug().
				EmbedObject(id).
				Str("subject", string(subject)).
				Msg("GetIDBySchema returned successfully")
		}
	}()
	return r.next.GetIDBySchema(ctx, subject, schema)
}
