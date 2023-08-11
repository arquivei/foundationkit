package avroutil

import "github.com/hamba/avro/v2"

type option func(*avro.Config)

func newConfig(options ...option) *avro.Config {
	config := &avro.Config{
		MaxByteSliceSize: 10 * 1024 * 1024, // 10 MiB
	}

	for _, opt := range options {
		opt(config)
	}

	return config
}

// WithMaxByteSliceSize is the maximum size of `bytes` or `string` types the Reader will create, defaulting to 10MiB.
// If this size is exceeded, the Reader returns an error. This can be disabled by setting a negative number.
func WithMaxByteSliceSize(maxSize int) option {
	return func(c *avro.Config) {
		c.MaxByteSliceSize = maxSize
	}
}
