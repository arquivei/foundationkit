package loggingmiddleware

import (
	"github.com/arquivei/foundationkit/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is a struct that configures the behavior of a
// loggingmiddleware. All fields should be set unless marked
// as optional.
type Config struct {
	// TruncRequestAt will truncate the serialized request if
	// the string length is above this value.
	// Setting to zero disables truncating.
	TruncRequestAt int

	// TruncResponseAt will truncate the serialized request if
	// the string length is above this value.
	// Setting to zero disables truncating.
	TruncResponseAt int

	// LogRequestIfLevel will log the request if this
	// log level is enabled in the logger.
	LogRequestIfLevel zerolog.Level

	// LogResponseIfLevel will log the response if this
	// log level is enabled in the logger.
	LogResponseIfLevel zerolog.Level

	// SuccessLevel is the log level of requests that are processed without error.
	SuccessLevel zerolog.Level

	// DefaultErrorLevel is the log level of requests that finished with some error.
	// May be override by ErrorCodeMapLevel or SeverityMapLevel
	DefaultErrorLevel zerolog.Level

	// ErrorCodeMapLevel is used to override the error level using the error code.
	// This is optional
	ErrorCodeMapLevel map[errors.Code]zerolog.Level

	// SeverityMapLevel is used to override the error level using the error severity.
	// This is optional
	SeverityMapLevel map[errors.Severity]zerolog.Level

	// Logger is the default logger to be put in the context.
	// If there is already a logger in the context, the context
	// is not updated and the existing logger is used.
	Logger *zerolog.Logger

	// Meta are extra keys logged under 'endpoint_meta' key.
	Meta Meta
}

// DefaultConfig contains sane defaults to configure a new logging middleware.
var DefaultConfig = Config{
	Logger: &log.Logger,

	LogRequestIfLevel:  zerolog.DebugLevel,
	TruncRequestAt:     200,
	LogResponseIfLevel: zerolog.DebugLevel,
	TruncResponseAt:    200,

	SuccessLevel:      zerolog.InfoLevel,
	DefaultErrorLevel: zerolog.ErrorLevel,

	SeverityMapLevel: map[errors.Severity]zerolog.Level{
		errors.SeverityInput:   zerolog.InfoLevel,
		errors.SeverityRuntime: zerolog.WarnLevel,
		errors.SeverityFatal:   zerolog.ErrorLevel,
	},
}
