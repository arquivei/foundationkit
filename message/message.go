package message

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

// Source is the type used to represent messages source
type Source string

func (s Source) String() string {
	return string(s)
}

// Type is the type used to represent messages type
type Type string

func (t Type) String() string {
	return string(t)
}

// DataVersion is the type used to represent messages data version
type DataVersion int

// SchemaVersion is the type used to represent mesages schema version
type SchemaVersion int

const (
	// SchemaVersion3 is the only schema version available right now.
	// All messages are produced with this as default.
	SchemaVersion3 SchemaVersion = 3
)

// Message is the default event structure
type Message struct {
	SchemaVersion SchemaVersion
	ID            string `json:"Id"`
	Source        Source
	Type          Type
	CreatedAt     string
	DataVersion   DataVersion
	Data          json.RawMessage
}

// MarshalZerologObject implements the zerolog marshaler so it can be logged using:
// log.With().EmbedObject(stakeholder).Msg("Some message")
func (m Message) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("msg_id", m.ID).
		Str("msg_source", m.Source.String()).
		Str("msg_created_at", m.CreatedAt).
		Str("msg_type", m.Type.String()).
		Int("msg_data_version", int(m.DataVersion))
}
