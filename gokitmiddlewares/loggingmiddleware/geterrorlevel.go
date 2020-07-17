package loggingmiddleware

import (
	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog"
)

func getErrorLevel(c Config, err error) zerolog.Level {
	lvl, ok := getLevelByErrorCode(c.ErrorCodeMapLevel, err)
	if !ok {
		lvl, ok = getLevelByErrorSeverity(c.SeverityMapLevel, err)
		if !ok {
			lvl = c.DefaultErrorLevel
		}
	}
	return lvl
}

func getLevelByErrorCode(m map[errors.Code]zerolog.Level, err error) (zerolog.Level, bool) {
	if len(m) == 0 {
		return zerolog.DebugLevel, false
	}
	lvl, ok := m[errors.GetCode(err)]
	if !ok {
		return zerolog.DebugLevel, false
	}
	return lvl, ok
}

func getLevelByErrorSeverity(m map[errors.Severity]zerolog.Level, err error) (zerolog.Level, bool) {
	if len(m) == 0 {
		return zerolog.DebugLevel, false
	}
	lvl, ok := m[errors.GetSeverity(err)]
	if !ok {
		return zerolog.DebugLevel, false
	}
	return lvl, ok
}
