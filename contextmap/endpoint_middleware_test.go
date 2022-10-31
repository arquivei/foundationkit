package contextmap

import (
	"context"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/assert"
)

// TestEndpointMiddleware tests that a request context, after passing through
// the endpoint middleware, is enriched with a ContextMap.
func TestEndpointMiddleware(t *testing.T) {
	endpointFn := func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		Ctx(ctx).Set("key", "value")
		Ctx(ctx).Set("key2", "value2")
		assert.Equal(t, "value", Ctx(ctx).Get("key"))
		assert.Equal(t, "value2", Ctx(ctx).Get("key2"))
		return nil, nil
	}
	endpointFn = endpoint.Chain(NewEndpointMiddleware())(endpointFn)

	_, _ = endpointFn(context.Background(), nil)
}
