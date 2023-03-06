package logger

import (
	"io"
	stdlog "log"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the log configuration struct
type Config struct {
	Level string `default:"info"`
	Human bool
}

// SetupLogger sets the global logger by configuring the global zerolog.Log and
// also the go's log package.
func Setup(config Config, version string, extraLogWriters ...io.Writer) {
	if config.Human {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	if len(extraLogWriters) > 0 {
		extraLogWriters = append(extraLogWriters, log.Logger)
		log.Logger = log.Output(zerolog.MultiLevelWriter(extraLogWriters...))
	}

	zerolog.SetGlobalLevel(MustParseLevel(config.Level))

	// Adds some global keys
	log.Logger = log.With().
		Str("version", version).
		Str("goversion", runtime.Version()).
		Logger()

	// Replace standard go logger with zerolog
	hooked := log.Hook(noLevelWarnHook{})
	stdlog.SetFlags(0)
	stdlog.SetOutput(hooked)
}

// MustParseLevel transforms a string in a zerolog level
func MustParseLevel(l string) zerolog.Level {
	zl, err := zerolog.ParseLevel(strings.ToLower(l))
	if err != nil {
		panic(err)
	}

	return zl
}
