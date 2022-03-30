package metricsmiddleware

import "context"

// ExternalMetrics is called after the internal metrifier is called.
// This functions should compute other metrics that are not computed by
// the internal metrifier (request latency and count).
type ExternalMetrics func(ctx context.Context, req, resp interface{}, err error)
