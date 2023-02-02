package app

import (
	"context"
	"sync"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog/log"
)

// ErrorPolicy specifies what should be done when a handler fails
type ErrorPolicy int

// ShutdownPriority is used to guide the execution of the shutdown handlers
// during a graceful shutdown. The shutdown is performed from the higher to the lowest
// priority
type ShutdownPriority uint8

const (
	// ErrorPolicyWarn prints the error as a warning and continues to the next handler. This is the default.
	ErrorPolicyWarn ErrorPolicy = iota
	// ErrorPolicyAbort stops the shutdown process and returns an error
	ErrorPolicyAbort
	// ErrorPolicyFatal logs the error as Fatal, it means the application will close immediately
	ErrorPolicyFatal
	// ErrorPolicyPanic panics if there is an error
	ErrorPolicyPanic
)

// ErrorPolicyString returns a string representation of a ErrorPolicy. This was intended for logging purposes.
func ErrorPolicyString(p ErrorPolicy) string {
	switch p {
	case ErrorPolicyAbort:
		return "abort"
	case ErrorPolicyFatal:
		return "fatal"
	case ErrorPolicyPanic:
		return "panic"
	case ErrorPolicyWarn:
		return "warn"
	default:
		return ""
	}
}

// ShutdownFunc is a shutdown function that will be executed when the app is shutting down.
type ShutdownFunc func(context.Context) error

// ShutdownHandler is a shutdown structure that allows configuring
// and storing shutdown information of an orchestrated shutdown flow.
type ShutdownHandler struct {
	Name     string
	Handler  ShutdownFunc
	Policy   ErrorPolicy
	Priority ShutdownPriority
	Timeout  time.Duration

	index int
	order int
	err   error
	once  sync.Once
}

// Execute runs the shutdown functions and handles timeout and error policy
func (sh *ShutdownHandler) Execute(ctx context.Context) error {
	sh.once.Do(func() {
		sh.doExecute(ctx)
		if sh.err != nil {
		}
	})
	return sh.err
}

func (sh *ShutdownHandler) doExecute(ctx context.Context) {
	const op = errors.Op("app.shutdownHandler.Execute")

	// Avoid running if the context is already closed
	if err := ctx.Err(); err != nil {
		sh.err = errors.E(errors.Op(sh.Name), err)
		sh.err = errors.E(op, sh.err)
		return
	}

	// Set the configured timeout, if any
	if sh.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, sh.Timeout)
		defer cancel()
	}

	// Execute the shutdown function and process the result
	sh.err = sh.Handler(ctx)
	if sh.err != nil {
		sh.err = errors.E(errors.Op(sh.Name), sh.err)
		sh.err = errors.E(op, sh.err)
		sh.applyErrorPolicy()
		return
	}

	log.Info().
		Str("handler", sh.Name).
		Uint8("shutdown_priority", uint8(sh.Priority)).
		Msg("[app] Shutdown successful.")
}

func (sh *ShutdownHandler) applyErrorPolicy() {
	switch sh.Policy {
	case ErrorPolicyWarn:
		log.Warn().
			Err(sh.err).
			Str("handler", sh.Name).
			Uint8("shutdown_priority", uint8(sh.Priority)).
			Msg("[app] Shutdown handler failed but graceful shutdown will continue.")
		sh.err = nil
	case ErrorPolicyAbort:
		// No need for logging here, this will happen latter
	case ErrorPolicyFatal:
		// KLUDGE: golang-ci linter is complaining about the log.Fatal() causing
		// the `defer cancel()` not to run. But on this case this is fine.
		//nolint:gocritic
		log.Fatal().
			Err(sh.err).
			Str("handler", sh.Name).
			Uint8("shutdown_priority", uint8(sh.Priority)).
			Msg("[app] Shutdown handler failed.")
	case ErrorPolicyPanic:
		panic(sh.err)
	default:
		panic(errors.Errorf("invalid error policy: %v", sh.Policy))
	}
}

// shutdownHeap is a heap implementation for the *shutdownHandler type
type shutdownHeap []*ShutdownHandler

func (sq shutdownHeap) Len() int {
	return len(sq)
}

func (sq shutdownHeap) Less(i, j int) bool {
	// If two items have the same priority, we use the first one inserted
	if sq[i].Priority == sq[j].Priority {
		return sq[i].order < sq[j].order
	}
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return sq[i].Priority > sq[j].Priority
}

func (sq shutdownHeap) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
	sq[i].index, sq[j].index = i, j
}

func (sq *shutdownHeap) Push(x interface{}) {
	sh := x.(*ShutdownHandler)
	n := len(*sq)
	sh.order = n
	sh.index = n
	*sq = append(*sq, sh)
}

func (sq *shutdownHeap) Pop() interface{} {
	old := *sq
	n := len(old)
	sh := old[n-1]
	old[n-1] = nil
	sh.index = -1
	*sq = old[0 : n-1]
	return sh
}
