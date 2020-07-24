package trace

import (
	"context"

	"github.com/rs/zerolog/log"
)

type contextKeyType int

const (
	contextKeyTrace contextKeyType = iota
	contextKeyLabels
	// Deprecated: Only used in Deprecated Methods
	contextKeyTraceID
)

// WithTrace returns the @parent context with the Trace @trace
func WithTrace(parent context.Context, trace Trace) context.Context {
	trace = ensureTraceNotEmpty(trace)
	return context.WithValue(parent, contextKeyTrace, trace)
}

// GetTraceFromContext returns the Trace saved in @ctx
func GetTraceFromContext(ctx context.Context) Trace {
	if t, ok := ctx.Value(contextKeyTrace).(Trace); ok {
		return t
	}
	log.Warn().
		Str("method", "trace.GetTraceFromContext").
		Msg("[FoundationKit] There is no Trace in context. Use trace.WithTrace(context.Context, trace.Trace)")
	return Trace{}
}

func withLabels(parent context.Context, labels map[string]string) context.Context {
	return context.WithValue(parent, contextKeyLabels, labels)
}

func getLabelsFromContext(ctx context.Context) map[string]string {
	if l, ok := ctx.Value(contextKeyLabels).(map[string]string); ok {
		return l
	}
	return nil
}

// GetIDFromContext returns the Trace ID in the context.
// Will return a empty ID if a Trace is not set in context
func GetIDFromContext(ctx context.Context) ID {
	return GetTraceFromContext(ctx).ID
}

// WithTraceAndLabels returns the @parent context with the Trace @trace
// and the labels @labels
func WithTraceAndLabels(parent context.Context, trace Trace, labels map[string]string) context.Context {
	parent = WithTrace(parent, trace)
	return withLabels(parent, labels)
}

// WithTraceID instantiates a new child context from @parent with the
// given @traceID value set
//
// Deprecated: Should use WithTrace instead
func WithTraceID(parent context.Context, traceID ID) context.Context {
	return context.WithValue(parent, contextKeyTraceID, traceID)
}

// GetTraceIDFromContext returns the trace ID set in the context, if any,
// or an empty trace id if none is set
//
// Deprecated: Should use GetTraceFromContext instead
func GetTraceIDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(contextKeyTraceID).(ID); ok {
		return id
	}
	return ID{}
}
