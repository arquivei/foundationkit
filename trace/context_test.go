package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetIDFromContext(ctx)))

	trace := newTrace(0)
	ctx = WithTrace(ctx, trace)

	assert.False(t, IDIsEmpty(GetIDFromContext(ctx)))
	assert.Equal(t, trace.ID.String(), GetIDFromContext(ctx).String())
}

func TestTraceOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetTraceFromContext(ctx).ID))

	ps := 0.5
	ctx = WithTrace(ctx, newTrace(ps))

	trace := GetTraceFromContext(ctx)

	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, *trace.ProbabilitySample, ps)
}

func TestTraceAndLabelsOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetTraceFromContext(ctx).ID))
	assert.Nil(t, nil, getLabelsFromContext(ctx))

	ps := 0.5
	ctx = WithTraceAndLabels(ctx, newTrace(ps), map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

	trace := GetTraceFromContext(ctx)
	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, *trace.ProbabilitySample, ps)

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
	assert.True(t, IDIsEmpty(GetTraceFromContext(ctx).ID))
	assert.Nil(t, nil, getLabelsFromContext(ctx))

	ps := 0.5
	ctx = WithTraceAndLabels(ctx, newTrace(ps), nil)

	trace := GetTraceFromContext(ctx)
	assert.False(t, IDIsEmpty(trace.ID))
	assert.Equal(t, *trace.ProbabilitySample, ps)

	labels := getLabelsFromContext(ctx)
	assert.Nil(t, labels)
}

func TestLabelsOperations(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, nil, getLabelsFromContext(ctx))

	ctx = withLabels(ctx, map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

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

//
// [DEPRECATED]
//

func TestTraceIDContextOperations(t *testing.T) {
	ctx := context.Background()
	assert.True(t, IDIsEmpty(GetTraceIDFromContext(ctx)))

	id := NewTraceID()
	ctx = WithTraceID(ctx, id)
	assert.Equal(t, id.String(), GetTraceIDFromContext(ctx).String())
}
