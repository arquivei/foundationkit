package cuf

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbededObject(cuf).Msg("Some message")
func (c CUF) MarshalZerologObject(e *zerolog.Event) {
	var result string
	if c.initialized {
		result = c.String()
	}
	e.Str("cuf", result)
}
