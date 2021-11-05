package trace

import (
	"strings"

	stackdriverexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/arquivei/foundationkit/app"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Config represents all trace's configuration
type Config struct {
	Exporter          string  `default:""`
	ProbabilitySample float64 `default:"0"`
	Stackdriver       struct {
		ProjectID string
	}
}

// Setup use @c to setup an trace exporter and returns a shutdown handler
func Setup(c Config) app.ShutdownFunc {
	var exporter trace.SpanExporter

	switch e := strings.ToLower(c.Exporter); e {
	case "stackdriver":
		var err error
		exporter, err = stackdriverexporter.New(
			stackdriverexporter.WithProjectID(c.Stackdriver.ProjectID),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create exporter")
		}
	case "":
		// No trace will be exported, but will be created
	default:
		log.Fatal().Str("exporter", e).Msg("This exporter is not supported")
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(
			trace.ParentBased(trace.TraceIDRatioBased(c.ProbabilitySample)),
		),
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.ForceFlush
}
