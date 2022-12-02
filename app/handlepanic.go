package app

import (
	"runtime/debug"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/stringsutil"
	"github.com/rs/zerolog/log"
)

// HandlePanic executes the function and, if the function panics, recovers from the panic
// and if panic persists, logs it.
func HandlePanic(f func()) {
	panicErr := errors.DontPanic(f)
	if panicErr != nil {
		log.Fatal().
			Err(panicErr).
			Str("panic_stack", stringsutil.Truncate(string(debug.Stack()), 1024)).
			Msg("App panicked!")
	}
}
