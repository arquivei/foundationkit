package stakeholder

import "github.com/rs/zerolog"

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbededObject(stakeholder).Msg("Some message")
func (s Stakeholder) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("stakeholder", s.String()).
		Str("stakeholder_type", TypeText(GetType(s)))
}
