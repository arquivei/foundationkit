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
	var id string
	if !i.IsEmpty() {
		id = fmt.Sprintf("%d-%s", i.timestamp, i.randomID)
	}
	return id
}

// IsEmpty returns true if there is no
// timestamp or randomID in Request ID @i
func (i ID) IsEmpty() bool {
	return i.timestamp == 0 && i.randomID == ""
}

func newID() ID {
	now := ulid.Now()
	return ID{
		timestamp: now,
		randomID:  ulid.MustNew(now, rand.Reader).String(),
	}
}
