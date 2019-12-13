package trace

import (
	"context"
	"strings"

	"github.com/arquivei/foundationkit/errors"
	"go.opencensus.io/trace"
)

// Span represents a span of a trace
type Span struct {
	span *trace.Span
}

// End ends the span and sets a label for @err, if exists
func (s *Span) End(err error) {
	if s.span != nil {
		s.span.End()
		if err != nil {
			s.span.AddAttributes(trace.StringAttribute("error_code", errors.GetErrorCode(err).String()))
		}
	}
}

// StartSpanWithParent returns a context and a span with the name @spanNameArgs.
// If exists a Trace in @ctx, the method will return a span with it as parent.
// Otherwise, the method will create a new span and return it
func StartSpanWithParent(ctx context.Context, spanNameArgs ...string) (newCtx context.Context, s Span) {
	t := GetTraceFromContext(ctx)

	parent := createSpanContext(t.ID.String(), *t.ProbabilitySample)

	newCtx, s.span = trace.StartSpanWithRemoteParent(ctx, spanName(spanNameArgs), *parent)

	setSpanLabels(newCtx, s.span)

	return
}

// StartSpan starts a span from @ctx and return it with a new context. The span returned
// has some labels defined in method setSpanLabels
func StartSpan(ctx context.Context, spanNameArgs ...string) (newCtx context.Context, s Span) {
	newCtx, s.span = trace.StartSpan(ctx, spanName(spanNameArgs))
	setSpanLabels(newCtx, s.span)
	return
}

func createSpanContext(traceIDStr string, probabilitySample float64) *trace.SpanContext {
	traceID := trace.TraceID(decode([]byte(traceIDStr)))
	samplingDecision := trace.ProbabilitySampler(probabilitySample)(trace.SamplingParameters{
		ParentContext: trace.SpanContext{},
		TraceID:       traceID,
	})

	var traceOptions uint32
	if samplingDecision.Sample {
		traceOptions = 1
	}

	return &trace.SpanContext{
		TraceID:      traceID,
		TraceOptions: trace.TraceOptions(traceOptions),
	}
}

func spanName(names []string) string {
	return strings.Join(names, "-")
}

func setSpanLabels(ctx context.Context, s *trace.Span) {
	labels := getLabelsFromContext(ctx)
	for key, value := range labels {
		s.AddAttributes(trace.StringAttribute(key, value))
	}
}
