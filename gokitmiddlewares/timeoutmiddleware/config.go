package timeoutmiddleware

import (
	"time"

	"github.com/arquivei/foundationkit/errors"
)

// Config is the configuration of the timeout middleware.
type Config struct {
	// Timeout is the duration of the timeout
	Timeout time.Duration
	// Wait indicates if the middleware should wait for
	// the next middleware to finish or it should return an error
	// right away.
	Wait bool
	// ErrorSeverity is the severity of the error when the context
	// is canceled, probably due the timeout
	ErrorSeverity errors.Severity
}

func NewDefaultConfig() Config {
	return Config{
		Timeout:       30 * time.Second,
		Wait:          false,
		ErrorSeverity: errors.SeverityRuntime,
	}
}
