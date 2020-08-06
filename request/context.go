package request

import "context"

type contextKeyType int

const (
	contextKeyID contextKeyType = iota
)

// WithRequestID returns a context with the given ID
//
// If there is already an ID in the context, the context is returned unchanged.
func WithRequestID(ctx context.Context, id ID) context.Context {
	if v := ctx.Value(contextKeyID); v == nil {
		return context.WithValue(ctx, contextKeyID, id)
	}
	return ctx
}

// WithNewRequestID calls WithRequestIDpassign a new ID
func WithNewRequestID(ctx context.Context) context.Context {
	return WithRequestID(ctx, newID())
}

// GetRequestIDFromContext returns the request ID in the context.
// Will return a empty ID if it is not set
func GetRequestIDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(contextKeyID).(ID); ok {
		return id
	}
	return ID{}
}
