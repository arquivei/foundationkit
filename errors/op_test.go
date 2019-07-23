package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpString(t *testing.T) {
	const op Op = "op"
	assert.Equal(t, "op", op.String())
}
