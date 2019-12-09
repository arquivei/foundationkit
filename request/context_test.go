package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDContextOperations(t *testing.T) {
	ctx := context.Background()
	assert.Empty(t, GetRequestIDFromContext(ctx))

	ctx = WithRequestID(ctx)
	assert.NotEmpty(t, GetRequestIDFromContext(ctx))
}
