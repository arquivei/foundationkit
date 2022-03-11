package metricsmiddleware

import "context"

// LabelsDecoder defines an interface to decode labels for the internal metrifier.
type LabelsDecoder interface {
	// Labels return the complete list of all available labels that will be returned by the Decoder.
	Labels() []string
	// Decode extracts a map of labels considering the request, response and error.
	Decode(ctx context.Context, req, resp interface{}, err error) map[string]string
}
