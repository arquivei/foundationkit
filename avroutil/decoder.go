package avroutil

import (
	"context"
	"encoding/binary"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/trace"

	"github.com/hamba/avro/v2"
)

// Decoder is able to transform avro's wire format data into golang's
// concrete types
type Decoder interface {
	// Decode decodes @data, in the wire format, into @output.
	Decode(ctx context.Context, data []byte, output interface{}) error
}

type implDecoder struct {
	schemaRepository schemaregistry.Repository
	avroAPI          avro.API
}

// NewDecoder returns a concrete implementation of Decoder, that
// fetches schemas in schema registry
func NewDecoder(schemaRepository schemaregistry.Repository, options ...option) Decoder {
	return &implDecoder{
		schemaRepository: schemaRepository,
		avroAPI:          newConfig(options...).Freeze(),
	}
}

func (i *implDecoder) Decode(
	ctx context.Context,
	data []byte,
	output interface{},
) error {
	const op = errors.Op("avroutil.implDecoder.Decode")
	ctx, span := trace.StartSpan(ctx, "AvroDecode")
	defer span.End(nil)

	schemaID, data, err := splitAvroWireFormatMessage(data)
	if err != nil {
		return errors.E(op, err, errors.SeverityInput)
	}

	schema, err := i.schemaRepository.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return errors.E(op, err)
	}

	err = i.avroAPI.Unmarshal(schema, data, output)
	return errors.E(op, err, errors.SeverityInput)
}

// splitAvroWireFormatMessage extracts the schema ID and the data from a
// message written with the wire format.
// The header is a 5-byte slice, where the first byte is equal to x00, and
// the last four represent the ID in a 32-bit big endian integer encoding.
func splitAvroWireFormatMessage(msg []byte) (schemaregistry.ID, []byte, error) {
	const op = errors.Op("avroutil.SplitAvroWireFormatMessage")
	if len(msg) < 5 {
		return 0, nil, errors.E(op, "invalid message length")
	}
	if msg[0] != 0x00 {
		return 0, nil, errors.E(op, "invalid magic byte")
	}
	schemaID := schemaregistry.ID(binary.BigEndian.Uint32(msg[1:5]))
	data := msg[5:]
	return schemaID, data, nil
}
