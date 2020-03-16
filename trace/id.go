package trace

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"sync"
)

// Generator is a interface to generate trace ids
type Generator interface {
	NewTraceID() ID
}

var traceIDGen Generator

func init() {
	gen := &defaultIDGenerator{}
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &gen.traceIDAdd,
	} {
		if err := binary.Read(crand.Reader, binary.LittleEndian, p); err != nil {
			panic(err)
		}
	}
	gen.traceIDRand = rand.New(rand.NewSource(rngSeed))
	traceIDGen = gen
}

// ID is an unique identifier (trace id) that can be use to identify
// one or more requests between distinct systems. It is a random-generated
// 16 bytes word, encoded as hexadecimal characters when in string format.
type ID [16]byte

// String will return the ID as 32 hexadecimal characters string
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// Parse returns an ID instance from a serialized traceID
func Parse(input string) ID {
	return decode([]byte(input))
}

// UnmarshalJSON parses an ID from a json. The ID is expected to a
// 32 hexadecimal characters string. This operation is case-insensitive.
func (id *ID) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	*id = decode(b)

	return nil
}

func decode(b []byte) ID {
	size := 32
	if len(b) < size {
		size = len(b)
	}
	b = b[:size]

	var new ID
	hex.Decode(new[:], b)

	return new
}

// MarshalJSON converts ID to a 32 hexadecimal characters string.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// NewTraceID generates a new random trace id
func NewTraceID() ID {
	return traceIDGen.NewTraceID()
}

// EnsureIDNotEmpty checks if the ID is not empty and return it, else it returns NewTraceID().
// The empty check follows the same rules as the IDIsEmpty function.
func EnsureIDNotEmpty(id ID) ID {
	if IDIsEmpty(id) {
		return traceIDGen.NewTraceID()
	}
	return id
}

// IDIsEmpty returns true if the @id is nil, empty or composed solely by zeroes.
func IDIsEmpty(id ID) bool {
	for _, b := range id {
		if b != 0 {
			return false
		}
	}
	return true
}

type defaultIDGenerator struct {
	sync.Mutex

	traceIDAdd  [2]uint64
	traceIDRand *rand.Rand
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// This function is thread safe.
func (gen *defaultIDGenerator) NewTraceID() ID {
	var tid [16]byte
	// Construct the trace ID from two outputs of traceIDRand, with a constant
	// added to each half for additional entropy.
	gen.Lock()
	defer gen.Unlock()

	binary.LittleEndian.PutUint64(tid[0:8], gen.traceIDRand.Uint64()+gen.traceIDAdd[0])
	binary.LittleEndian.PutUint64(tid[8:16], gen.traceIDRand.Uint64()+gen.traceIDAdd[1])
	return ID(tid)
}
