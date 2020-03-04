package message

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/arquivei/foundationkit/errors"
	"github.com/oklog/ulid"
)

func getTypeName(data interface{}) string {
	valueOf := reflect.ValueOf(data)

	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	}
	return valueOf.Type().Name()
}

// ParseTypeAndDataVersion returns the message type and data version for a given data
// The type a hyphen separated string
// The data should be a struct whose name is like (all those results in the same type):
// MyMessageV1 -> "my-message", 1
// myMessageV1 -> "my-message", 1
// MYVersionV1 -> "my-message", 1
func ParseTypeAndDataVersion(data interface{}) (Type, DataVersion, error) {
	const op = errors.Op("message.ParseTypeAndDataVersion")

	typeName := getTypeName(data)
	vIdx := strings.LastIndex(typeName, "V")

	if vIdx < 1 || vIdx == len(typeName)-1 {
		return "", 0, errors.E(op, errors.Errorf("invalid type name, expected '<type>V<version>' but got '%s'", typeName))
	}

	v, err := strconv.Atoi(typeName[vIdx+1:])
	if err != nil {
		return "", 0, errors.E(op, err)
	}

	runes := []rune(typeName[0:vIdx])
	hasNeighborCaseChanged := func(i int) bool {
		if i == 0 || i == vIdx-1 { // case is considered not changing in the extremes of the slice
			return false
		}
		return unicode.IsLower(runes[i+1]) || unicode.IsLower(runes[i-1])
	}

	t := strings.Builder{}
	t.Grow(vIdx)
	for i := 0; i < vIdx; i++ {
		r := runes[i]
		if unicode.IsUpper(r) && hasNeighborCaseChanged(i) {
			t.WriteRune('-')
		}
		t.WriteRune(unicode.ToLower(r))
	}
	return Type(t.String()), DataVersion(v), nil
}

// newMessageID creates a new ULID ID for messages
func newMessageID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

// New creates a new message with auto discovering of message type and version
// The type name of the data param must be in the format <type>V<version>, like nsuReceivedV3
// The data must be json marshable
func New(ctx context.Context, source Source, data interface{}) (Message, error) {
	const op = errors.Op("message.New")
	messageType, messageDataVersion, err := ParseTypeAndDataVersion(data)
	if err != nil {
		return Message{}, errors.E(op, err)
	}
	d, err := json.Marshal(data)
	if err != nil {
		return Message{}, errors.E(op, err)
	}
	return Message{
		SchemaVersion: SchemaVersion3,
		ID:            newMessageID(),
		Source:        source,
		Type:          messageType,
		CreatedAt:     time.Now().Format(time.RFC3339),
		DataVersion:   messageDataVersion,
		Data:          json.RawMessage(d),
		Context:       ctx,
	}, nil
}
