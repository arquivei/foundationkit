package kitlog

import (
	"github.com/arquivei/foundationkit/errors"

	"github.com/go-kit/log"
	"github.com/rs/zerolog"
)

type zerologWrapper struct {
	logger zerolog.Logger
}

// NewKitLogger returns a new zerolog wrapper
func NewKitLogger(logger zerolog.Logger) log.Logger {
	return &zerologWrapper{
		logger: logger,
	}
}

func (l *zerologWrapper) Log(args ...interface{}) error {
	var hasError bool
	var err error
	ctx := l.logger.With()
	logger := ctx.Logger()
	logPointer := &logger

	if len(args)%2 != 0 {
		args = append(args, "")
	}

	for i := 0; i < len(args); i += 2 {
		if args[i] == "err" || args[i] == "error" {
			hasError = true
			err = errors.Errorf("%v", args[i+1])
		} else {
			ctx = ctx.Interface(args[i].(string), args[i+1])
		}
	}
	if hasError && err.Error() != "" {
		logPointer.Error().Err(err).Msg("Logging from go-kit")
	} else {
		logPointer.Info().Msg("Logging from go-kit")
	}
	return nil
}
