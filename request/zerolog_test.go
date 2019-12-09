package request

import (
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNotPanicMarshalZerologObject(t *testing.T) {
	tests := []struct {
		name string
		id   ID
	}{
		{
			name: "New ID",
			id:   newID(),
		},
		{
			name: "Empty ID",
			id:   ID{},
		},
	}

	for _, test := range tests {
		assert.NotPanics(t, func() { log.Debug().EmbedObject(test.id) }, test.name)
	}
}
