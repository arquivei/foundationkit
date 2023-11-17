package trace

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TraceInfo carries the trace informations
type TraceInfo struct {
	ID        string
	IsSampled bool
}

// GetTraceInfoFromContext returns a TraceInfo from the context @ctx or logs
// if there is no TraceInfo in context
func GetTraceInfoFromContext(ctx context.Context) TraceInfo {
	sc := trace.SpanContextFromContext(ctx)
	t := TraceInfo{
		ID:        sc.TraceID().String(),
		IsSampled: sc.IsSampled(),
	}

	if t.ID != "" {
		return t
	}
	log.Warn().
		Str("method", "trace.GetTraceInfoFromContext").
		Msg("[FoundationKit] There is no Trace Info in context.")
	return TraceInfo{}
}

// ToMap extracts the current trace from the context and put it on a map.
// This can be used to serialize the trace context on json messages.
func ToMap(ctx context.Context) map[string]string {
	m := map[string]string{}
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(m))
	return m
}

// FromMap injects the trace context from the map into the context. This is the
// inverse of ToMap and could be used to extract trace context form json messages.
func FromMap(ctx context.Context, m map[string]string) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m))
}
