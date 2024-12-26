package request

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/oklog/ulid/v2"
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

// UnmarshalJSON parses an ID from a json
func (i *ID) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)

	id, err := Parse(string(b))
	if err != nil {
		return err
	}

	*i = id
	return nil
}

// MarshalJSON converts ID to a string.
func (i ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// Parse read the given @id into an request ID.
func Parse(id string) (ID, error) {
	slices := strings.Split(id, ("-"))

	if len(slices) != 2 {
		return ID{}, errors.New("wrong format for request id")
	}

	timestamp, err := strconv.Atoi(slices[0])
	if err != nil {
		return ID{}, err
	}

	return ID{
		timestamp: uint64(math.Abs(float64(timestamp))),
		randomID:  slices[1],
	}, nil
}

// IsEmpty returns true if an request ID is in it's zero-value format
func IsEmpty(id ID) bool {
	return (id == ID{})
}

func newID() ID {
	now := ulid.Now()
	return ID{
		timestamp: now,
		randomID:  ulid.MustNew(now, rand.Reader).String(),
	}
}
