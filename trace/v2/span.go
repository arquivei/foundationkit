package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Span is the individual component of a trace. It represents a single named
// and timed operation of a workflow that is traced
type Span struct {
	span trace.Span
}

// End stops an span sampling
func (s Span) End() {
	s.span.End()
}

// Start creates a new span. If other spans were created using @ctx, this method
// will bind them all
func Start(ctx context.Context, name string) (context.Context, Span) {
	ctx, span := otel.Tracer("").Start(ctx, name)
	return withTraceInfo(ctx, span), Span{span}
}
