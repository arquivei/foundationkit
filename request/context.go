package request

import "context"

type contextKeyType int

const (
	contextKeyID contextKeyType = iota
)

// WithRequestID checks the context if it already has a ID. If not,
// creates a new one and returns a new context with it.
func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyID, newID())
}

// GetRequestIDFromContext returns the request ID in the context.
// Will return a empty ID if it is not set
func GetRequestIDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(contextKeyID).(ID); ok {
		return id
	}
	return ""
}
