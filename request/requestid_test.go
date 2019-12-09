package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, ID(""), GetRequestIDFromContext(ctx))

	ctx = WithRequestID(ctx)
	assert.NotEmpty(t, GetRequestIDFromContext(ctx))
}
