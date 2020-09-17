package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetIDFromContext(ctx)))

	defaultProbabilitySample = 0
	trace := newTrace()
	ctx = WithTrace(ctx, trace)

	assert.False(t, IDIsEmpty(GetIDFromContext(ctx)))
	assert.Equal(t, trace.ID.String(), GetIDFromContext(ctx).String())
}

func TestTraceOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetFromContext(ctx).ID))

	defaultProbabilitySample = 0.5
	ctx = WithTrace(ctx, newTrace())

	trace := GetFromContext(ctx)

	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, defaultProbabilitySample, *trace.ProbabilitySample)
}

func TestTraceAndLabelsOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetFromContext(ctx).ID))
	assert.Nil(t, getLabelsFromContext(ctx))

	defaultProbabilitySample = 0.5
	ctx = WithTraceAndLabels(ctx, newTrace(), map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

	trace := GetFromContext(ctx)
	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, defaultProbabilitySample, *trace.ProbabilitySample)

	labels := getLabelsFromContext(ctx)
	for key, value := range labels {
		switch key {
		case "k1":
			assert.Equal(t, "v1", value)
		case "k2":
			assert.Equal(t, "v2", value)
		default:
			assert.FailNow(t, "none key is valid")
		}
	}
}

func TestTraceAndLabelsOperations_LabelsNil(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetFromContext(ctx).ID))
	assert.Nil(t, getLabelsFromContext(ctx))

	defaultProbabilitySample = 0.5
	ctx = WithTraceAndLabels(ctx, newTrace(), nil)

	trace := GetFromContext(ctx)
	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, defaultProbabilitySample, *trace.ProbabilitySample)

	labels := getLabelsFromContext(ctx)
	assert.Nil(t, labels)
}

func TestLabelsOperations(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, getLabelsFromContext(ctx))

	ctx = WithLabels(ctx, map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

	labels := getLabelsFromContext(ctx)
	assert.Equal(t, "v1", labels["k1"])
	assert.Equal(t, "v2", labels["k2"])

	ctx = WithLabels(ctx, map[string]string{
		"k2": "v22",
		"k3": "v3",
	})

	labels = getLabelsFromContext(ctx)
	assert.Equal(t, "v1", labels["k1"])
	assert.Equal(t, "v22", labels["k2"])
	assert.Equal(t, "v3", labels["k3"])
}

// [DEPRECATED] Testing a Deprecated Methods
func TestTraceIDContextOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetTraceIDFromContext(ctx)))

	id := NewID()
	ctx = WithTraceID(ctx, id)
	assert.Equal(t, id.String(), GetTraceIDFromContext(ctx).String())
}
