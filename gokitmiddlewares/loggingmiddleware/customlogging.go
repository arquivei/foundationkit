package loggingmiddleware

import (
	"context"

	"github.com/rs/zerolog"
)

// LoggableEndpointRequest, if implemented, by the endpoint request type, will
// be used to enrich the log along with the EnrichLogWithRequest config.
type LoggableEndpointRequest interface {
	// EnrichLog enriches the current log context @zctx with the request data.
	EnrichLog(
		ctx context.Context,
		zctx zerolog.Context,
	) (context.Context, zerolog.Context)
}

// LoggableEndpointResponse, if implemented, by the endpoint response type, will
// be used to enrich the log along with the EnrichLogWithResponse config.
type LoggableEndpointResponse interface {
	// EnrichLog enriches the current log context @zctx with the response data.
	EnrichLog(
		ctx context.Context,
		zctx zerolog.Context,
	) zerolog.Context
}
