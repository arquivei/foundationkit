package trace

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbededObject(cuf).Msg("Some message")
func (i ID) MarshalZerologObject(e *zerolog.Event) {
	e.Str("trace_id", i.String())
}
