package splitio

import "golang.org/x/net/context"

// MockFeaturesToContext writes the mocked behaviours in @features to the context
func MockFeaturesToContext(ctx context.Context, features map[Feature]bool) context.Context {
	for k, v := range features {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// MockFeatureListToContext writes the mocked behaviours in @features to the context
func MockFeatureListToContext(ctx context.Context, features []Feature) context.Context {
	for _, feature := range features {
		ctx = context.WithValue(ctx, feature, true)
	}
	return ctx
}
