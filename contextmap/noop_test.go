package contextmap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoop(t *testing.T) {
	noopCtx := Ctx(context.Background())
	assert.Equal(t, "", noopCtx.String())
	assert.Equal(t, noopCtx, noopCtx.Set("key", "value"))
	assert.Equal(t, nil, noopCtx.Get("key"))
}
