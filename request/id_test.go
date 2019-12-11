package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDString(t *testing.T) {
	tests := []struct {
		name      string
		timestamp uint64
		randomID  string
		expected  string
	}{
		{
			name:      "With timestamp and random ID",
			timestamp: 19,
			randomID:  "random",
			expected:  "19-random",
		},
		{
			name:     "Without timestamp and with random ID",
			randomID: "random",
			expected: "0-random",
		},
		{
			name:      "With timestamp and without random ID",
			timestamp: 19,
			expected:  "19-",
		},
		{
			name:     "Without timestamp and random ID",
			expected: "",
		},
	}

	for _, test := range tests {
		id := ID{test.timestamp, test.randomID}
		assert.Equal(t, test.expected, id.String(), test.name)
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		timestamp uint64
		randomID  string
		expected  bool
	}{
		{
			name:      "With timestamp and random ID",
			timestamp: 19,
			randomID:  "random",
			expected:  false,
		},
		{
			name:     "Without timestamp and with random ID",
			randomID: "random",
			expected: false,
		},
		{
			name:      "With timestamp and without random ID",
			timestamp: 19,
			expected:  false,
		},
		{
			name:     "Without timestamp and random ID",
			expected: true,
		},
	}

	for _, test := range tests {
		id := ID{test.timestamp, test.randomID}
		assert.Equal(t, test.expected, id.IsEmpty(), test.name)
	}
}

func TestIDMarshall(t *testing.T) {
	id := ID{timestamp: 1576072698019, randomID: "01DVTM1P53ZVBRCCM4F9SCRK09"}

	var serializedMessage struct {
		ID *ID
	}
	serializedMessage.ID = &id

	buffer, err := json.Marshal(&serializedMessage)
	assert.NoError(t, err)
	assert.Equal(t, `{"ID":"1576072698019-01DVTM1P53ZVBRCCM4F9SCRK09"}`, string(buffer))
}

func TestIDUnmarshall(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expected      ID
		expectedError string
	}{
		{
			name:     "Success",
			json:     `{"ID":"1576072698019-01DVTM1P53ZVBRCCM4F9SCRK09"}`,
			expected: ID{timestamp: 1576072698019, randomID: "01DVTM1P53ZVBRCCM4F9SCRK09"},
		},
		{
			name:          "Wrong format",
			json:          `{"ID":"1-2-3"}`,
			expectedError: "wrong format for request id",
		},
		{
			name:          "Atoi failed",
			json:          `{"ID":"a-01DVTM1P53ZVBRCCM4F9SCRK09"}`,
			expectedError: `strconv.Atoi: parsing "a": invalid syntax`,
		},
	}

	for _, test := range tests {
		var serializedMessage struct {
			ID *ID
		}
		err := json.Unmarshal([]byte(test.json), &serializedMessage)
		if test.expectedError != "" {
			assert.EqualError(t, err, test.expectedError, test.name)
		} else {
			assert.NoError(t, err, test.name)
			assert.Equal(t, test.expected, *serializedMessage.ID, test.name)
		}
	}
}
