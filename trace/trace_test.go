package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type traceTest struct {
	name                      string
	trace                     Trace
	defaultProbabilitySample  float64
	expectedProbabilitySample float64
	expectedID                *ID
}

func TestEnsureTraceNotEmpty(t *testing.T) {
	ps := 1.9
	newID := NewTraceID()

	tests := []traceTest{
		{
			name:                      "All Empty",
			trace:                     Trace{},
			defaultProbabilitySample:  0.5,
			expectedProbabilitySample: 0.5,
		},
		{
			name: "ID Empty",
			trace: Trace{
				ProbabilitySample: &ps,
			},
			defaultProbabilitySample:  0.5,
			expectedProbabilitySample: ps,
		},
		{
			name: "Probability Sample Empty",
			trace: Trace{
				ID: newID,
			},
			defaultProbabilitySample:  0.5,
			expectedProbabilitySample: 0.5,
		},
		{
			name: "None Empty",
			trace: Trace{
				ID:                newID,
				ProbabilitySample: &ps,
			},
			defaultProbabilitySample:  0.5,
			expectedProbabilitySample: ps,
			expectedID:                &newID,
		},
	}

	for _, test := range tests {
		defaultProbabilitySample = test.defaultProbabilitySample
		test.trace = ensureTraceNotEmpty(test.trace)
		ctx := WithTrace(context.Background(), test.trace)
		assertTrace(t, test, "Trace from Test")

		traceFromCtx := GetTraceFromContext(ctx)
		test.trace = traceFromCtx
		assertTrace(t, test, "Trace from Context")
	}
}

func assertTrace(t *testing.T, test traceTest, traceOrigin string) {
	if !assert.NotNil(t, test.trace, fmt.Sprintf("trace should not be nil [%s][%s]", traceOrigin, test.name)) {
		return
	}
	if !assert.NotNil(t, test.trace.ProbabilitySample, fmt.Sprintf("probability sample should not be nil [%s][%s]", traceOrigin, test.name)) {
		return
	}
	if !assert.False(t, IDIsEmpty(test.trace.ID), fmt.Sprintf("trace ID should not be empty [%s][%s]", traceOrigin, test.name)) {
		return
	}

	assert.Equal(t, test.expectedProbabilitySample, *test.trace.ProbabilitySample, fmt.Sprintf("probability sample should be equal [%s][%s]", traceOrigin, test.name))
	if test.expectedID != nil {
		assert.Equal(t, test.expectedID.String(), test.trace.ID.String(), fmt.Sprintf("trace ID should be equal [%s][%s]", traceOrigin, test.name))
	} else {
		assert.NotNil(t, test.trace.ID, fmt.Sprintf("trace ID should not be nil [%s][%s]", traceOrigin, test.name))
	}
}
