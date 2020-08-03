package loggingmiddleware

import (
	"context"

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

	// EnrichRequestFunc is used to enrich the logger context with the request.
	//
	// It is called by the middleware to enrich the logger with information from
	// inside the request just before calling the next endpoint. This functions should
	// interpret the request and extract information and return the new zerolog context.
	//
	// Remember that the request can be logged as a whole by using LogRequestIfLevel and
	// TruncRequestAt.
	//
	// If set to 'nil' it is disabled.
	EnrichLogWithRequest EnrichLogWithRequestFunc

	// EnrichLogWithResponse is used to enrich the logger context with the response and error.
	//
	// It is called by the middleware to enrich the logger with information from
	// inside the respose and error after calling the next endpoint. This functions should
	// interpret the response and error and extract information and return the new zerolog context.
	//
	// Remember that the response can be logged as a whole by using LogResponseIfLevel and
	// TruncResponseAt.
	//
	// If set to 'nil' it is disabled.
	EnrichLogWithResponse EnrichLogWithResponseFunc
}

// EnrichLogWithRequestFunc is a function that receives a zerolog Context and the request. It should parse and
// extract information and return a new enriched zerolog context.
type EnrichLogWithRequestFunc func(ctx context.Context, zctx zerolog.Context, request interface{}) (context.Context, zerolog.Context)

// EnrichLogWithResponseFunc is a function that receives a zerolog Context and the response and error. It should
// parse and extract information and return a new enriched zerolog context. The error is already handled my the logging
// middleware if it is an error.Error, but it is passed to this functions should you need to extract another kind of data
// from the error.
type EnrichLogWithResponseFunc func(ctx context.Context, zctx zerolog.Context, response interface{}, err error) zerolog.Context

// DefaultConfig contains sane defaults to configure a new logging middleware.
func NewDefaultConfig() Config {
	return Config{
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
}
