package apiutil

import (
	"context"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
)

// BaseEndpointResponse should be used to extend responses from endpoint
type BaseEndpointResponse struct {
	// Example: "51e8cb62dbc57774a7720ef828e96c34"
	TraceID   trace.ID    `json:"trace_id,omitempty"` // [Legacy] Already exists in Trace.
	Trace     trace.Trace `json:"trace,omitempty"`
	RequestID request.ID  `json:"request_id,omitempty"`
}

// CreateBaseEndpointResponse creates and returns a struct of BaseEndpointResponse filled with data from context and error.
func CreateBaseEndpointResponse(ctx context.Context) BaseEndpointResponse {
	return BaseEndpointResponse{
		TraceID:   trace.GetIDFromContext(ctx), // [Legacy] Already exists in Trace.
		Trace:     trace.GetFromContext(ctx),
		RequestID: request.GetIDFromContext(ctx),
	}
}
