package request

import "context"

type contextKeyType int

const (
	contextKeyID contextKeyType = iota
)

// WithID returns a context with the given ID
//
// If there is already an ID in the context, the context is returned unchanged.
func WithID(ctx context.Context, id ID) context.Context {
	if v := ctx.Value(contextKeyID); v == nil {
		return context.WithValue(ctx, contextKeyID, id)
	}
	return ctx
}

// WithNewID calls WithID with a new ID
func WithNewID(ctx context.Context) context.Context {
	return WithID(ctx, newID())
}

// GetIDFromContext returns the request ID in the context.
// Will return a empty ID if it is not set
func GetIDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(contextKeyID).(ID); ok {
		return id
	}
	return ID{}
}
