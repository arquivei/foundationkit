package nsu

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbededObject(nsu).Msg("Some message")
func (nsu NSU) MarshalZerologObject(e *zerolog.Event) {
	e.Str("nsu", nsu.String())
}
