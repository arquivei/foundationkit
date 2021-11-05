package trace

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type contextKeyType int

const (
	contextKeyTraceInfo contextKeyType = iota
)

// TraceInfo carries the trace informations
type TraceInfo struct {
	ID        string
	IsSampled bool
}

// GetTraceInfoFromContext returns a TraceInfo from the context @ctx or logs
// if there is no TraceInfo in context
func GetTraceInfoFromContext(ctx context.Context) TraceInfo {
	if t, ok := ctx.Value(contextKeyTraceInfo).(TraceInfo); ok {
		return t
	}
	log.Warn().
		Str("method", "trace.GetTraceInfoFromContext").
		Msg("[FoundationKit] There is no Trace Info in context.")
	return TraceInfo{}
}

func withTraceInfo(ctx context.Context, s trace.Span) context.Context {
	if v := ctx.Value(contextKeyTraceInfo); v != nil {
		return ctx
	}

	traceInfo := TraceInfo{
		ID:        s.SpanContext().TraceID().String(),
		IsSampled: s.SpanContext().IsSampled(),
	}

	return context.WithValue(ctx, contextKeyTraceInfo, traceInfo)
}
