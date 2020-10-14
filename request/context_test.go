package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDContextOperations(t *testing.T) {
	ctx := context.Background()
	assert.Empty(t, GetIDFromContext(ctx))

	ctx = WithNewID(ctx)
	assert.NotEmpty(t, GetIDFromContext(ctx))
}
