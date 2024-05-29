package avroutil

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro/v2"
)

// Encoder is able to transform golang's concrete types into avro's wire format
type Encoder interface {
	// Encode encodes @input in the wire format
	Encode(input interface{}) ([]byte, error)
}

type implEncoder struct {
	wireFormatEncoder WireFormatEncoder
	writerSchema      avro.Schema
	avroAPI           avro.API
}

// NewEncoder returns a concrete implementation of Decoder, that
// fetches schemas in schema registry
// Parameters:
//   - @schemaRepository: repository for avro schemas.
//   - @writerSchemaStr: avro schema, in the AVSC format, used to marshall the
//     objects. This schema must be previously registered in the schema registry
//     exactly as provided.
func NewEncoder(
	ctx context.Context,
	schemaRepository schemaregistry.Repository,
	subject schemaregistry.Subject,
	writerSchemaStr string,
	options ...option,
) (Encoder, error) {
	const op = errors.Op("avroutil.NewEncoder")

	encoder, err := NewWireFormatEncoder(ctx, schemaRepository, subject, writerSchemaStr)
	if err != nil {
		return nil, errors.E(err, op)
	}

	parsedAvroSchema, err := avro.Parse(writerSchemaStr)
	if err != nil {
		return nil, errors.E(err, op)
	}

	return &implEncoder{
		wireFormatEncoder: encoder,
		writerSchema:      parsedAvroSchema,
		avroAPI:           newConfig(options...).Freeze(),
	}, nil
}

func (e *implEncoder) Encode(input interface{}) ([]byte, error) {
	const op = errors.Op("avroutil.implEncoder.Encode")

	avroData, err := e.avroAPI.Marshal(e.writerSchema, input)
	if err != nil {
		return nil, errors.E(err, op)
	}

	wireFormat, err := e.wireFormatEncoder.BinaryToWireFormat(avroData)
	if err != nil {
		return nil, errors.E(err, op)
	}

	return wireFormat, nil
}
