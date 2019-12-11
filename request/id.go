package request

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

// UnmarshalJSON parses an ID from a json
func (i *ID) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	slices := strings.Split(string(b), ("-"))

	if len(slices) != 2 {
		return errors.New("wrong format for request id")
	}

	timestamp, err := strconv.Atoi(slices[0])
	if err != nil {
		return err
	}

	i.timestamp = uint64(timestamp)
	i.randomID = slices[1]
	return nil
}

// MarshalJSON converts ID to a string.
func (i ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func newID() ID {
	now := ulid.Now()
	return ID{
		timestamp: now,
		randomID:  ulid.MustNew(now, rand.Reader).String(),
	}
}
