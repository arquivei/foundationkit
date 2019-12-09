package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDString(t *testing.T) {
	id := ID{
		timestamp: 19,
		randomID:  "random",
	}
	assert.Equal(t, "19-random", id.String())
}
