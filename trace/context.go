package trace

import (
	"context"

	"github.com/rs/zerolog/log"
)

type contextKeyType int

const (
	contextKeyTrace contextKeyType = iota
	contextKeyLabels
)

// WithTrace if there is no trace in the context, returns the @ctx with the @trace
// else returns the @ctx unchanged
func WithTrace(ctx context.Context, trace Trace) context.Context {
	if v := ctx.Value(contextKeyTrace); v != nil {
		return ctx
	}
	trace = ensureTraceNotEmpty(trace)
	return context.WithValue(ctx, contextKeyTrace, trace)
}

// WithNewTrace returns the same as WithTrace, but instead of receiving
// a trace, it creates a new one.
func WithNewTrace(ctx context.Context) context.Context {
	return WithTrace(ctx, Trace{})
}

// GetFromContext returns the Trace saved in @ctx
func GetFromContext(ctx context.Context) Trace {
	if t, ok := ctx.Value(contextKeyTrace).(Trace); ok {
		return t
	}
	log.Warn().
		Str("method", "trace.GetFromContext").
		Msg("[FoundationKit] There is no Trace in context. Use trace.WithTrace(context.Context, trace.Trace)")
	return Trace{}
}

// WithLabels returns the @parent context with the labels @labels
// If there are already labels in the context, the new labels will be merged
// into the existing one. If a label with the same key already exists, it will
// be overwritten
func WithLabels(parent context.Context, labels map[string]string) context.Context {
	if currentLabels := getLabelsFromContext(parent); currentLabels != nil {
		for k, v := range labels {
			currentLabels[k] = v
		}
		return parent
	}
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
	return GetFromContext(ctx).ID
}

// WithTraceAndLabels returns the @parent context with the Trace @trace
// and the labels @labels
func WithTraceAndLabels(parent context.Context, trace Trace, labels map[string]string) context.Context {
	parent = WithTrace(parent, trace)
	return WithLabels(parent, labels)
}
