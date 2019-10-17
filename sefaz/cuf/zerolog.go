package cuf

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbededObject(cuf).Msg("Some message")
func (c CUF) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("cuf", c.String())
}
