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
type ShutdownPriority int

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

// ShutdownFunc is a shutdown function that will be executed when the app is shutting down.
type ShutdownFunc func(context.Context) error

type shutdownHandler struct {
	Name     string
	Timeout  time.Duration
	Policy   ErrorPolicy
	Priority ShutdownPriority

	Shutdown ShutdownFunc
	executed bool
	err      error
	mu       sync.Mutex
	index    int
}

// NewShutdownHandler returns a new
func NewShutdownHandler(name string, fn ShutdownFunc, options ...interface{}) *shutdownHandler {
	if name == "" {
		panic("shutdown handler must have a name")
	}
	sh := &shutdownHandler{
		Name:     name,
		Shutdown: fn,
	}
	for _, o := range options {
		switch opt := o.(type) {
		case time.Duration:
			sh.Timeout = opt
		case ErrorPolicy:
			sh.Policy = opt
		case ShutdownPriority:
			sh.Priority = opt
		default:
			panic(errors.Errorf("invalid shutdown handler option: [%T]%v", o, o))
		}
	}
	return sh
}

// Execute runs the shutdown functions and handles timeout and error policy
func (sh *shutdownHandler) Execute(ctx context.Context) error {
	const op = errors.Op("app.shutdownHandler.Execute")

	sh.mu.Lock()
	defer sh.mu.Unlock()

	// The shutdown should run only once
	// Future calls will return the result of the first call
	if sh.executed {
		return sh.err
	}
	sh.executed = true

	// Avoid runnign if the context is already closed
	if ctx.Err() != nil {
		sh.err = ctx.Err()
		return sh.err
	}

	// Set the configured timeout, if any
	if sh.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, sh.Timeout)
		defer cancel()
	}

	// Execute the shutdown function and process the result
	err := sh.Shutdown(ctx)
	if err != nil {
		err = errors.E(op, errors.E(errors.Op(sh.Name), err))
		switch sh.Policy {
		case ErrorPolicyWarn:
			log.Ctx(ctx).Warn().Err(err).Msg("Shutdown handler failed")
		case ErrorPolicyAbort:
			sh.err = err
			// No need for logging here, this will happen latter
		case ErrorPolicyFatal:
			log.Ctx(ctx).Fatal().Err(err).Msg("Shutdown handler failed")
		case ErrorPolicyPanic:
			panic(err)
		default:
			panic(errors.Errorf("invalid error policy: %v", sh.Policy))
		}
	}

	log.Ctx(ctx).Info().
		Str("handler", sh.Name).
		Int("shutdown_priority", int(sh.Priority)).
		Msg("Shutdown successfull")

	return sh.err
}

// shutdownHeap is a heap implementation for the *shutdownHandler type
type shutdownHeap []*shutdownHandler

func (sq shutdownHeap) Len() int {
	return len(sq)
}

func (sq shutdownHeap) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return sq[i].Priority > sq[j].Priority
}

func (sq shutdownHeap) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
	sq[i].index, sq[j].index = i, j
}

func (sq *shutdownHeap) Push(x interface{}) {
	sh := x.(*shutdownHandler)
	sh.index = len(*sq)
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
