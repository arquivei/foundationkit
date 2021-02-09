package splitio

import (
	"io"
	"strings"

	"github.com/arquivei/foundationkit/errors"
	"github.com/rs/zerolog/log"
)

type zerologLogger struct{}

// NewZerologLogger returns a new logger for Split IO
func NewZerologLogger() io.Writer {
	return &zerologLogger{}
}

func (e *zerologLogger) Write(p []byte) (int, error) {
	errStr := strings.TrimPrefix(string(p), "ERROR - ")
	errStr = strings.TrimSuffix(errStr, "\n")

	log.Logger.Warn().
		Err(errors.E(errors.Op("splitio"), errStr)).
		Msg("Failed to check if Feature is enabled in Split IO")

	return len(p), nil
}
