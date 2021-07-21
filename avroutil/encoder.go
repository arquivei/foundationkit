package avroutil

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"

	"github.com/hamba/avro"
)

// Encoder is able to transform golang's concrete types into avro's wire format
type Encoder interface {
	// Encode encodes @input in the wire format
	Encode(input interface{}) ([]byte, error)
}

type implEncoder struct {
	writerSchemaID schemaregistry.ID
	writerSchema   avro.Schema
}

// NewEncoder returns a concrete implementation of Decoder, that
// fetches schemas in schema registry
// Parameters:
// - @schemaRepository: repository for avro schemas.
// - @writerSchemaStr: avro schema, in the avsc format, used to marshall the
//   objects. This schema must be previusly registered in the schema registry
//   exactly as provided.
func NewEncoder(
	ctx context.Context,
	schemaRepository schemaregistry.Repository,
	subject schemaregistry.Subject,
	writerSchemaStr string,
) (Encoder, error) {
	const op = errors.Op("avroutil.NewEncoder")
	schemaID, parsedSchema, err := schemaRepository.GetIDBySchema(
		ctx,
		subject,
		writerSchemaStr,
	)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return &implEncoder{
		writerSchemaID: schemaID,
		writerSchema:   parsedSchema,
	}, nil
}

func (e *implEncoder) Encode(input interface{}) ([]byte, error) {
	const op = errors.Op("avroutil.implEncoder.Encode")

	avroData, err := avro.Marshal(e.writerSchema, input)
	if err != nil {
		return nil, errors.E(op, err)
	}

	wireFormat, err := joinAvroWireFormatMessage(e.writerSchemaID, avroData)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return wireFormat, nil
}

// joinAvroWireFormatMessage uses the schema ID and the data from a message
// to write it as an Avro Wire Format message.
// The header is a 5-byte slice, where the first byte is equal to x00, and
// the last four represent the ID in a 32-bit big endian integer encoding.
func joinAvroWireFormatMessage(
	schemaID schemaregistry.ID,
	msg []byte,
) ([]byte, error) {
	const op = errors.Op("joinAvroWireFormatMessage")

	buf := new(bytes.Buffer)

	if err := buf.WriteByte(0x00); err != nil {
		return nil, errors.E(op, err)
	}

	if err := binary.Write(buf, binary.BigEndian, int32(schemaID)); err != nil {
		return nil, errors.E(op, err)
	}

	if _, err := buf.Write(msg); err != nil {
		return nil, errors.E(op, err)
	}

	return buf.Bytes(), nil
}
