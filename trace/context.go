package trace

import "context"

type contextKeyType int

const (
	contextKeyTraceID contextKeyType = iota
)

// WithTraceID instantiates a new child context from @parent with the
// given @traceID value set
func WithTraceID(parent context.Context, traceID ID) context.Context {
	return context.WithValue(parent, contextKeyTraceID, traceID)
}

// GetTraceIDFromContext returns the trace ID set in the context, if any,
// or an empty trace id if none is set
func GetTraceIDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(contextKeyTraceID).(ID); ok {
		return id
	}
	return ID{}
}
