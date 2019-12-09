package request

import (
	"crypto/rand"
	"fmt"

	"github.com/oklog/ulid"
)

// ID is an identifier to uniquely identify a request
type ID struct {
	timestamp uint64
	randomID  string
}

func (i ID) String() string {
	return fmt.Sprintf("%d-%s", i.timestamp, i.randomID)
}

func newID() ID {
	now := ulid.Now()
	return ID{
		timestamp: now,
		randomID:  ulid.MustNew(now, rand.Reader).String(),
	}
}
