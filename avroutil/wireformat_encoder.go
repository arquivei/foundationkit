package avroutil

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/schemaregistry"
)

// WireFormatEncoder uses the Schema ID provided in the constructor to encodes Avro Messages of the same type into its
// Wire Format form. The `WireFormatEncoder` is used undercover by `Encoder` to deal with the Wire Format transformation
// step.
type WireFormatEncoder interface {
	BinaryToWireFormat(avroInput []byte) ([]byte, error)
}

type wireFormatEncoder struct {
	writerSchemaID schemaregistry.ID
}

// NewWireFormatEncoder returns a concrete implementation of WireFormatEncoder, that
// fetches schemas in schema registry
// Parameters:
// - @schemaRepository: repository for avro schemas.
// - @writerSchemaStr: avro schema, in the avsc format, used to marshall the
//   objects. This schema must be previously registered in the schema registry
//   exactly as provided.
func NewWireFormatEncoder(
	ctx context.Context,
	schemaRepository schemaregistry.Repository,
	subject schemaregistry.Subject,
	writerSchemaStr string,
) (WireFormatEncoder, error) {
	const op = errors.Op("avroutil.NewWireFormatEncoder")
	schemaID, _, err := schemaRepository.GetIDBySchema(
		ctx,
		subject,
		writerSchemaStr,
	)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return &wireFormatEncoder{
		writerSchemaID: schemaID,
	}, nil
}

// BinaryToWireFormat uses the schema ID and the data from a message
// to write it as an Avro Wire Format message.
// The header is a 5-byte slice, where the first byte is equal to x00, and
// the last four represent the ID in a 32-bit big endian integer encoding.
func (e *wireFormatEncoder) BinaryToWireFormat(avroInput []byte) ([]byte, error) {
	const op = errors.Op("avroutil.wireFormatEncoder.BinaryToWireFormat")

	buf := new(bytes.Buffer)

	if err := buf.WriteByte(0x00); err != nil {
		return nil, errors.E(op, err)
	}

	if err := binary.Write(buf, binary.BigEndian, int32(e.writerSchemaID)); err != nil {
		return nil, errors.E(op, err)
	}

	if _, err := buf.Write(avroInput); err != nil {
		return nil, errors.E(op, err)
	}

	return buf.Bytes(), nil
}
