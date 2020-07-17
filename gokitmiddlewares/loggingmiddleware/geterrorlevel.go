package loggingmiddleware

import (
	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog"
)

func getErrorLevel(c Config, err error) zerolog.Level {
	if lvl, ok := getLevelByErrorCode(c.ErrorCodeMapLevel, err); ok {
		return lvl
	}

	if lvl, ok := getLevelByErrorSeverity(c.SeverityMapLevel, err); ok {
		return lvl
	}

	return c.DefaultErrorLevel
}

func getLevelByErrorCode(m map[errors.Code]zerolog.Level, err error) (zerolog.Level, bool) {
	if len(m) == 0 {
		return zerolog.DebugLevel, false
	}

	if lvl, ok := m[errors.GetCode(err)]; ok {
		return lvl, true
	}

	return zerolog.DebugLevel, false
}

func getLevelByErrorSeverity(m map[errors.Severity]zerolog.Level, err error) (zerolog.Level, bool) {
	if len(m) == 0 {
		return zerolog.DebugLevel, false
	}

	if lvl, ok := m[errors.GetSeverity(err)]; ok {
		return lvl, true
	}

	return zerolog.DebugLevel, false
}
