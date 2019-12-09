package request

import (
	"crypto/rand"
	"fmt"

	"github.com/oklog/ulid"
)

// ID is an identifier to uniquely identify a request
type ID string

func (i ID) String() string {
	return string(i)
}

func newID() ID {
	now := ulid.Now()
	return ID(fmt.Sprintf("%d-%s", now, ulid.MustNew(now, rand.Reader).String()))
}
