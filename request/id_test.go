package request

import (
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
