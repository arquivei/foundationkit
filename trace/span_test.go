package trace

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSpan struct {
	name                     string
	defaultProbabilitySample float64
	parent                   Trace
	expectedIsSample         bool
}

func TestSpan(t *testing.T) {
	// PS = probability sample
	// DPS = default probability sample (context's probability sample)

	zero := 0.0
	one := 1.0

	tests := []testSpan{
		{
			name:                     "With Parent with PS = 1 and DPS = 0; Should Sample",
			defaultProbabilitySample: 0,
			parent: Trace{
				ID:                NewTraceID(),
				ProbabilitySample: &one,
			},
			expectedIsSample: true,
		},
		{
			name:                     "With Parent with PS = 0 and DPS = 1; Should not Sample",
			defaultProbabilitySample: 1,
			parent: Trace{
				ID:                NewTraceID(),
				ProbabilitySample: &zero,
			},
			expectedIsSample: false,
		},
		{
			name:                     "With Parent with PS = 1 and DPS = nil; Should Sample",
			defaultProbabilitySample: -1,
			parent: Trace{
				ID:                NewTraceID(),
				ProbabilitySample: &one,
			},
			expectedIsSample: true,
		},
		{
			name:                     "With Parent with PS = 0 and DPS = nil; Should not Sample",
			defaultProbabilitySample: -1,
			parent: Trace{
				ID:                NewTraceID(),
				ProbabilitySample: &zero,
			},
			expectedIsSample: false,
		},
		{
			name:                     "Without Parent; DPS = 0; Should not Sample",
			defaultProbabilitySample: 0,
			parent:                   Trace{},
			expectedIsSample:         false,
		},
		{
			name:                     "Without Parent; DPS = 1; Should Sample",
			defaultProbabilitySample: 1,
			parent:                   Trace{},
			expectedIsSample:         true,
		},
	}

	for _, test := range tests {
		defaultProbabilitySample = test.defaultProbabilitySample
		ctx := WithTrace(context.Background(), test.parent)

		var s Span
		ctx, s = StartSpanWithParent(ctx, "test")
		assertResponse(ctx, s, test, t)

		ctx, s = StartSpan(ctx, "test")
		assertResponse(ctx, s, test, t)
	}
}

func assertResponse(ctx context.Context, s Span, test testSpan, t *testing.T) {
	assert.Equal(t, test.expectedIsSample, s.span.SpanContext().IsSampled(), fmt.Sprintf("isSample not equal [%s]", test.name))
	if !IDIsEmpty(test.parent.ID) {
		assert.Equal(t, test.parent.ID.String(), s.span.SpanContext().TraceID.String(), fmt.Sprintf("trace ID not equal [%s]", test.name))
	}

	trace := GetTraceFromContext(ctx)
	if !assert.NotNil(t, trace, fmt.Sprintf("trace should not be nil [%s]", test.name)) {
		return
	}
	if !assert.NotNil(t, trace.ProbabilitySample, fmt.Sprintf("probability sample should not be nil [%s]", test.name)) {
		return
	}
	if !assert.False(t, IDIsEmpty(trace.ID), fmt.Sprintf("trace ID should not be empty [%s]", test.name)) {
		return
	}

	if test.parent.ProbabilitySample != nil {
		assert.Equal(t, *test.parent.ProbabilitySample, *trace.ProbabilitySample, fmt.Sprintf("probability sample should be equal [%s]", test.name))
	} else {
		assert.Equal(t, test.defaultProbabilitySample, *trace.ProbabilitySample, fmt.Sprintf("probability sample should be equal (DPS) [%s]", test.name))
	}

	if !IDIsEmpty(test.parent.ID) {
		assert.Equal(t, test.parent.ID.String(), trace.ID.String(), fmt.Sprintf("trace ID should be equal [%s]", test.name))
	}

	ctx = withLabels(ctx, map[string]string{
		"key": "value",
	})

	assert.NotPanics(t, func() { s.End(errors.New("error label")) })
	assert.NotPanics(t, func() { setSpanLabels(ctx, s.span) })
}

func TestCreateSampleContext(t *testing.T) {
	traceID := NewTraceID()
	span0 := createSpanContext(traceID.String(), 0)
	span1 := createSpanContext(traceID.String(), 1)

	assert.Equal(t, traceID.String(), span0.TraceID.String())
	assert.Equal(t, traceID.String(), span1.TraceID.String())
	assert.False(t, span0.IsSampled())
	assert.True(t, span1.IsSampled())
}

func TestSpanName(t *testing.T) {
	assert.Equal(t, "bla1-bla2-bla3", spanName([]string{"bla1", "bla2", "bla3"}))
}

func TestSetSpanLabels(t *testing.T) {
	ctx := withLabels(context.Background(), map[string]string{
		"key": "value",
	})
	var s Span
	ctx, s = StartSpan(ctx, "test")
	assert.NotPanics(t, func() { setSpanLabels(ctx, s.span) })
}

func TestEnd(t *testing.T) {
	ctx := withLabels(context.Background(), map[string]string{
		"key": "value",
	})
	var s Span
	_, s = StartSpan(ctx, "test")
	assert.NotPanics(t, func() { s.End(errors.New("error label")) })
}
