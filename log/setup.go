package log

import (
	"context"
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
	Human bool   `default:"false"`
	Hook  struct {
		Stackdriver bool `default:"false"`
	}
}

// SetupLogger sets the global logger by configuring the global zerolog.Log and
// also the go's log package.
func SetupLogger(config Config, version string, extraLogWriters ...io.Writer) {
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

	// Adds stackdriver severity hook
	if config.Hook.Stackdriver {
		log.Logger = log.Logger.Hook(stackdriverSeverityHook{})
	}

	// Replace standard go logger with zerolog
	hooked := log.Hook(noLevelWarnHook{})
	stdlog.SetFlags(0)
	stdlog.SetOutput(hooked)
}

// SetupLoggerWithContext returns a context enriched with a logger. The logger
// is created using SetupLogger, what implies that it will be also available
// globally.
func SetupLoggerWithContext(ctx context.Context, config Config, version string,
	extraLogWriters ...io.Writer) context.Context {
	SetupLogger(config, version, extraLogWriters...)
	return log.Logger.WithContext(ctx)
}

// MustParseLevel transforms a string in a zerolog level
func MustParseLevel(l string) zerolog.Level {
	zl, err := zerolog.ParseLevel(strings.ToLower(l))
	if err != nil {
		panic(err)
	}

	return zl
}
