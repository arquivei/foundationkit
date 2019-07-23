package log

import (
	"cloud.google.com/go/logging"
	"github.com/rs/zerolog"
)

type stackdriverSeverityHook struct{}

func (h stackdriverSeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Str("severity", LevelToSeverity(level).String())
}

// LevelToSeverity converts a zerolog level to the stackdriver severity
// Stackdriver has more levels than zerolog so we skip some severities.
// By default we set info when no level is provided.
var LevelToSeverity = func(level zerolog.Level) logging.Severity {
	switch level {
	case zerolog.DebugLevel:
		return logging.Debug
	// Let info falls into the defualt
	case zerolog.WarnLevel:
		return logging.Warning
	case zerolog.ErrorLevel:
		return logging.Error
	case zerolog.FatalLevel:
		return logging.Alert
	case zerolog.PanicLevel:
		return logging.Emergency
	default:
		return logging.Info
	}
}
