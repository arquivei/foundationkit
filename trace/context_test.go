package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceIDContextOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetTraceIDFromContext(ctx)))

	id := NewTraceID()
	ctx = WithTraceID(ctx, id)
	assert.Equal(t, id.String(), GetTraceIDFromContext(ctx).String())
}
