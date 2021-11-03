package trace

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged
// using: log.With().EmbededObject(t).Msg("Some message")
func (t TraceInfo) MarshalZerologObject(e *zerolog.Event) {
	e.Str("trace_id", t.ID)
	e.Bool("trace_sampled", t.IsSampled)
}
