package trace

import (
	stackdriverexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/sdk/trace"
)

// StackdriverConfig contains all the configuration of the Starkdriver exporter.
type StackdriverConfig struct {
	ProjectID string
}

func newStackdriverExporter(c StackdriverConfig) (trace.SpanExporter, error) {
	return stackdriverexporter.New(
		stackdriverexporter.WithProjectID(c.ProjectID),
	)
}
