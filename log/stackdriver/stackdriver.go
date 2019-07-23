package stackdriver

import (
	"encoding/json"

	"cloud.google.com/go/logging"
	"github.com/arquivei/foundationkit/log"
	"github.com/rs/zerolog"
)

type stackdriveWriter struct {
	logger *logging.Logger
}

// FlushFunc blocks and flushes the writer buffer. This should be called before the program exits to send all remaining entries to stackdriver
type FlushFunc func() error

// NewStackdriveLevelWriter returns a new writer using async logging from stackdriver. It also returns a flush function that should be called before the program exits to clear the internal buffers
func NewStackdriveLevelWriter(client *logging.Client, name string) (zerolog.LevelWriter, FlushFunc) {
	logger := client.Logger(name)
	return &stackdriveWriter{
		logger: client.Logger(name),
	}, logger.Flush
}

func (l *stackdriveWriter) Write(p []byte) (n int, err error) {
	return l.WriteLevel(zerolog.InfoLevel, p)
}

func (l *stackdriveWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	l.logger.Log(logging.Entry{Severity: log.LevelToSeverity(level), Payload: json.RawMessage(p)})
	return len(p), nil
}
