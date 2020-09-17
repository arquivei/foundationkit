package trace

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceIDGeneration(t *testing.T) {
	// Should generate random non-empty traceIDs
	id1 := NewID()
	id2 := NewID()
	assert.False(t, IDIsEmpty(id1))
	assert.False(t, IDIsEmpty(id2))
	assert.NotEqual(t, id1, id2)
}

func TestTraceIDMarshall(t *testing.T) {
	serialized := `{"ID":"5be27e0caf6c932033e5525a0979b074"}`
	id := ID{91, 226, 126, 12, 175, 108, 147, 32, 51, 229, 82, 90, 9, 121, 176, 116}
	var serializedMsg struct {
		ID *ID
	}

	serializedMsg.ID = &id

	buffer, err := json.Marshal(&serializedMsg)
	assert.NoError(t, err)
	assert.Equal(t, serialized, string(buffer))
}

func TestTraceIDUnmarshall(t *testing.T) {
	serialized := `{"ID":"5be27e0caf6c932033e5525a0979b074"}`
	id := ID{91, 226, 126, 12, 175, 108, 147, 32, 51, 229, 82, 90, 9, 121, 176, 116}
	var serializedMsg struct {
		ID *ID
	}
	err := json.Unmarshal([]byte(serialized), &serializedMsg)
	assert.NoError(t, err)
	assert.Equal(t, id, *serializedMsg.ID)
}

func TestDecode(t *testing.T) {
	tests := []struct {
		description string
		input       []byte
		expectedID  ID
	}{
		{
			description: "decoding a correct input",
			input:       []byte("12345678901234567890123456789012"),
			expectedID:  ID{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12},
		},
		{
			description: "decode is case insensitive",
			input:       []byte("aAbBcCdDeEfFaaAAbbBBccCCddDDeeEE"),
			expectedID:  ID{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0xaa, 0xaa, 0xbb, 0xbb, 0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee},
		},
		{
			description: "smaller input are suffixed with zeros",
			input:       []byte("1234567890123456789012345678"),
			expectedID:  ID{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00},
		},
		{
			description: "larger values are truncated",
			input:       []byte("12345678901234567890123456789012bebacafe"),
			expectedID:  ID{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12},
		},
		{
			description: "decodes goes only up to the first invalid character, after it comes 0's suffix",
			input:       []byte("1234567890_123456789012345678901"),
			expectedID:  ID{0x12, 0x34, 0x56, 0x78, 0x90, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, test := range tests {
		id := decode(test.input)
		assert.Equal(t, test.expectedID, id, "[%s] id differs", test.description)
	}
}

func TestTraceIDSerialization_Invalid(t *testing.T) {
	buffer := []byte("{\"ID\":\"kartoffelpüree mit käse und remoulade\"}")
	var serializedMsg struct {
		ID *ID
	}
	err := json.Unmarshal(buffer, &serializedMsg)
	assert.Nil(t, err)
	assert.True(t, IDIsEmpty(*serializedMsg.ID))
}

func TestEnsureIDNotEmpty(t *testing.T) {
	var id ID
	assert.True(t, IDIsEmpty(id))
	newID := EnsureIDNotEmpty(id)
	assert.False(t, IDIsEmpty(newID))
	sameID := EnsureIDNotEmpty(newID)
	assert.Equal(t, newID.String(), sameID.String())
}
