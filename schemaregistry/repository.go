package schemaregistry

import (
	"context"
	"strconv"

	"github.com/hamba/avro/v2"
	"github.com/rs/zerolog"
)

// ID is the schema registry's schema ID
type ID uint32

// MarshalZerologObject implements the zerolog marshaler so it can be logged
// using: log.With().EmbedObject(id).Msg("Some message")
func (i ID) MarshalZerologObject(e *zerolog.Event) {
	e.Str("schemaregistry_id", strconv.Itoa(int(i)))
}

// Subject is the schema registry's subject for schemas
type Subject string

// Repository is responsible for retrieving schemas for a given ID
type Repository interface {
	GetSchemaByID(
		ctx context.Context,
		id ID,
	) (avro.Schema, error)
	GetIDBySchema(
		ctx context.Context,
		subject Subject,
		schema string,
	) (ID, avro.Schema, error)
}
