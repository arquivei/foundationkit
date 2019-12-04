package trace

import (
	"context"
)

// Trace represents the informations that should be
// passed through systems
type Trace struct {
	// ID represents the Trace ID used in logging and
	// trace views. It will be used as the main span in
	ID ID

	// ProbabilitySample represents if the span will be
	// sampled or not. The two possibles values are 0 and 1
	ProbabilitySample *float64
}

func newTrace(defaultProbabilitySample float64) Trace {
	return Trace{
		ID:                NewTraceID(),
		ProbabilitySample: &defaultProbabilitySample,
	}
}

func createTraceIfEmpty(ctx context.Context, t *Trace, defaultProbabilitySample float64) context.Context {
	if t == nil || IDIsEmpty(t.ID) || t.ProbabilitySample == nil {
		*t = newTrace(defaultProbabilitySample)
	}
	ctx = WithTrace(ctx, *t)
	return ctx
}
