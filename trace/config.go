package trace

import (
	"strings"
	"time"

	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/bridge/opencensus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var defaultProbabilitySample float64

// Config represents the informations that must be
// set to configure the trace
type Config struct {
	Exporter          string  `default:""`
	ProbabilitySample float64 `default:"0"`
	Stackdriver       struct {
		ProjectID string
	}
}

// SetupTrace configure a trace exporter, defined in @c
func SetupTrace(c Config) {
	switch exporter := strings.ToLower(c.Exporter); exporter {
	case "stackdriver":
		start := time.Now()
		gcpExporter, err := gcptrace.New(gcptrace.WithProjectID(c.Stackdriver.ProjectID))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create stackdriver trace exporter")
		}
		tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(gcpExporter))
		opencensus.InstallTraceBridge(opencensus.WithTracerProvider(tp))
		log.Info().Dur("took", time.Since(start)).Msg("Stackdriver loaded")
	case "":
	default:
		log.Fatal().Str("exporter", exporter).Msg("This exporter is not supported")
	}
	defaultProbabilitySample = c.ProbabilitySample
}
