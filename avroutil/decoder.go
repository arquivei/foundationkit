package avroutil

import (
	"context"
	"encoding/binary"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
	"github.com/arquivei/foundationkit/trace"

	"github.com/hamba/avro"
)

// Decoder is able to transform avro's wire format data into golang's
// concrete types
type Decoder interface {
	// Decode decodes @data, in the wire format, into @output.
	Decode(ctx context.Context, data []byte, output interface{}) error
}

type implDecoder struct {
	schemaRepository schemaregistry.Repository
}

// NewDecoder returns a concrete implementation of Decoder, that
// fetches schemas in schema registry
func NewDecoder(schemaRepository schemaregistry.Repository) Decoder {
	return &implDecoder{
		schemaRepository: schemaRepository,
	}
}

func (i *implDecoder) Decode(ctx context.Context, data []byte, output interface{}) error {
	const op = errors.Op("avroutil.implDecoder.Decode")
	ctx, span := trace.StartSpan(ctx, "AvroDecode")
	defer span.End(nil)

	schemaID, data, err := SplitAvroWireFormatMessage(data)
	if err != nil {
		return errors.E(op, err, errors.SeverityInput)
	}

	schema, err := i.schemaRepository.GetSchemaByID(ctx, schemaID)
	if err != nil {
		return errors.E(op, err)
	}

	err = avro.Unmarshal(schema, data, output)
	return errors.E(op, err, errors.SeverityInput)
}

// DecodeWireFormatMessage decodes any @msg data encoded with the wire format
// into the @output variable. The @output parameter must be defined as specified by the
// library: github.com/hamba/avro
// Deprecated: Use avroutil.Decoder instead
func DecodeWireFormatMessage(
	ctx context.Context,
	msg []byte,
	schemaRegistry schemaregistry.Repository,
	output interface{},
) error {
	const op = errors.Op("avroutil.DecodeWireFormatMessage")
	err := NewDecoder(schemaRegistry).Decode(ctx, msg, output)
	return errors.E(op, err)
}

// SplitAvroWireFormatMessage extracts the schema ID and the data from a
// message written with the wire format.
// The header is a 5-byte slice, where the first byte is equal to x00, and
// the last four represent the ID in a 32-bit big endian integer encoding.
// Deprecated: This is a low level implementation detail of avro decoding. Use
// the high-level avroutil.Decoder instead
func SplitAvroWireFormatMessage(msg []byte) (schemaregistry.ID, []byte, error) {
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
