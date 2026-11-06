package trace

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/oauth2/google"
)

const (
	gcpTelemetryEndpoint = "telemetry.googleapis.com"
	gcpTraceScope        = "https://www.googleapis.com/auth/trace.append"
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
		ctx := context.Background()
		// Cloud Trace is reached through the Telemetry OTLP endpoint, authenticated
		// with Application Default Credentials.
		client, err := google.DefaultClient(ctx, gcpTraceScope)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create stackdriver trace http client")
		}
		gcpExporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(gcpTelemetryEndpoint),
			otlptracehttp.WithHTTPClient(client),
			otlptracehttp.WithHeaders(map[string]string{"x-goog-user-project": c.Stackdriver.ProjectID}),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create stackdriver trace exporter")
		}
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(gcpExporter),
			sdktrace.WithResource(resource.NewSchemaless(
				attribute.String("gcp.project_id", c.Stackdriver.ProjectID),
			)),
		)
		opencensus.InstallTraceBridge(opencensus.WithTracerProvider(tp))
		log.Info().Dur("took", time.Since(start)).Msg("Stackdriver loaded")
	case "":
	default:
		log.Fatal().Str("exporter", exporter).Msg("This exporter is not supported")
	}
	defaultProbabilitySample = c.ProbabilitySample
}
