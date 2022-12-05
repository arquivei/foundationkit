package app

import (
	"runtime/debug"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/stringsutil"
	"github.com/rs/zerolog/log"
)

// PanicHandler recovers from the panic and if panic persists, logs it.
func PanicHandler() {
	if r := recover(); r != nil {
		log.Fatal().
			Err(errors.NewFromRecover(r)).
			Str("panic_stack", stringsutil.Truncate(string(debug.Stack()), 1024)).
			Msg("App panicked!")
	}
}
