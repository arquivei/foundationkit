package logger

import "github.com/rs/zerolog"

type noLevelWarnHook struct{}

func (h noLevelWarnHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.NoLevel {
		e.Str("level", zerolog.WarnLevel.String())
	}
}
