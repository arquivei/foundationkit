package gzip

import (
	"bytes"
	"compress/gzip"

	"github.com/arquivei/foundationkit/errors"
)

// Compress the @input using the gzip format
func Compress(input []byte) ([]byte, error) {
	const op = errors.Op("gzip.Compress")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(input)
	if err != nil {
		return nil, errors.E(err, op)
	}
	err = zw.Close()
	if err != nil {
		return nil, errors.E(err, op)
	}
	return buf.Bytes(), nil
}

// MustCompress compresses @input using gzip, and panic if it fails
func MustCompress(input []byte) []byte {
	const op = errors.Op("gzip.MustCompress")
	output, err := Compress(input)
	if err != nil {
		panic(errors.E(err, op))
	}
	return output
}

// Decompress @input in the gzip format
func Decompress(input []byte) ([]byte, error) {
	const op = errors.Op("gzip.Decompress")
	b := bytes.NewReader(input)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, errors.E(err, op, errors.KV("step", "gzip.NewReader"))
	}

	var responseBuffer bytes.Buffer
	_, err = responseBuffer.ReadFrom(r)
	if err != nil {
		return nil, errors.E(err, op, errors.KV("step", "responseBuffer.ReadFrom"))
	}

	return responseBuffer.Bytes(), nil
}

// MustDecompress decompresses @input in the gzip format, and panic if it fails
func MustDecompress(input []byte) []byte {
	const op = errors.Op("gzip.MustDecompress")
	output, err := Decompress(input)
	if err != nil {
		panic(errors.E(err, op))
	}
	return output
}
