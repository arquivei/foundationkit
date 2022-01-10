package gzip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGunzip(t *testing.T) {
	validGzip := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xca, 0xcf, 0x4, 0x4, 0x0, 0x0, 0xff, 0xff, 0x6b, 0xbc, 0xd2, 0x97, 0x2, 0x0, 0x0, 0x0}
	output, err := Decompress(validGzip)
	assert.NoError(t, err)
	assert.Equal(t, []byte("oi"), output)

	invalidGzip := []byte("invalid")
	_, err = Decompress(invalidGzip)
	assert.Error(t, err)

	assert.NotPanics(t, func() {
		output := MustDecompress(validGzip)
		assert.Equal(t, []byte("oi"), output)
	})

	assert.Panics(t, func() {
		MustDecompress(invalidGzip)
	})
}

func TestGzip(t *testing.T) {
	expectedGzip := []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xca, 0xcf, 0x4, 0x4, 0x0, 0x0, 0xff, 0xff, 0x6b, 0xbc, 0xd2, 0x97, 0x2, 0x0, 0x0, 0x0}
	validInput := []byte("oi")
	output, err := Compress(validInput)
	assert.NoError(t, err)
	assert.Equal(t, expectedGzip, output)

	assert.NotPanics(t, func() {
		output := MustCompress(validInput)
		assert.Equal(t, expectedGzip, output)
	})
}
