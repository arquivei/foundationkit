package app

import (
	"container/heap"
	"context"
	"testing"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestShutdownhandlerHeap(t *testing.T) {
	h := shutdownHeap{}
	heap.Init(&h)
	assert.Equal(t, h.Len(), 0, "a new heap must be empty")

	assert.Panics(t, func() {
		heap.Push(&h, nil)
	}, "Push should panic if value is nil")

	assert.Panics(t, func() {
		heap.Push(&h, ShutdownHandler{})
	}, "Push should panic if if value type is not *ShutdownHandler")

	sh1 := &ShutdownHandler{
		Name:     "sh1",
		Priority: ShutdownPriority(0),
	}

	sh2 := &ShutdownHandler{
		Name:     "sh2",
		Priority: ShutdownPriority(10),
	}

	sh3 := &ShutdownHandler{
		Name:     "sh3",
		Priority: ShutdownPriority(5),
	}

	sh4 := &ShutdownHandler{
		Name:     "sh4",
		Priority: ShutdownPriority(5),
	}

	heap.Push(&h, sh1)
	heap.Push(&h, sh2)
	heap.Push(&h, sh4)
	heap.Push(&h, sh3)
	assert.Equal(t, h.Len(), 4, "heap should have 4 elements")

	p1 := heap.Pop(&h)
	p2 := heap.Pop(&h)
	p3 := heap.Pop(&h)
	p4 := heap.Pop(&h)

	assert.Equal(t, sh1, p4, "sh1 has de lowest priority and must be poped last")
	assert.Equal(t, sh2, p1, "sh2 has the highest priority and must be poped first")
	assert.Equal(t, sh3, p3, "sh3 must be poped after sh4 and before sh1")
	assert.Equal(t, sh4, p2, "sh4 must be poped after sh2 and before sh3")

	assert.Equal(t, h.Len(), 0, "heap should be empty")
}

func TestShutdownHandlerExecute(t *testing.T) {
	assert.Panics(t, func() {
		sh := &ShutdownHandler{}
		sh.Execute(context.TODO())
	}, "should panic if Handler is not set")

	sh := &ShutdownHandler{
		Name: "my_shutdown_handler",
		Handler: func(context.Context) error {
			return nil
		},
	}

	assert.False(t, sh.executed)
	assert.NoError(t, sh.err)

	err := sh.Execute(context.TODO())
	assert.NoError(t, err)
	assert.NoError(t, sh.err)
	assert.True(t, sh.executed)

	sh = &ShutdownHandler{
		Name: "my_failed_shutdown_handler",
		Handler: func(context.Context) error {
			return errors.New("my error")
		},
		Policy: ErrorPolicyAbort,
	}

	err = sh.Execute(context.TODO())
	assert.EqualError(t, err, "app.shutdownHandler.Execute: my_failed_shutdown_handler: my error")

	err2 := sh.Execute(context.TODO())
	assert.Equal(t, err, err2, "a second execution of handler should return the first error")
}

func TestShutdownHandlerExecute_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sh := &ShutdownHandler{
		Name: "my_failed_shutdown_handler",
		Handler: func(context.Context) error {
			return errors.New("my error")
		},
		Policy: ErrorPolicyAbort,
	}

	err := sh.Execute(ctx)
	assert.True(t, sh.executed)
	assert.EqualError(t, err, "app.shutdownHandler.Execute: my_failed_shutdown_handler: context canceled")
}

func TestShutdownHandlerExecute_Timeout(t *testing.T) {
	sh := &ShutdownHandler{
		Name: "my_failed_shutdown_handler",
		Handler: func(ctx context.Context) error {
			select {
			case <-time.After(2 * time.Nanosecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		Policy:  ErrorPolicyAbort,
		Timeout: time.Nanosecond,
	}

	ctx := context.Background()
	err := sh.Execute(ctx)
	assert.True(t, sh.executed)
	assert.EqualError(t, err, "app.shutdownHandler.Execute: my_failed_shutdown_handler: context deadline exceeded")
}
