package trace

import (
	"errors"
	"strings"

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
	Stackdriver       StackdriverConfig
	OTLP              OTLPConfig
}

// Setup use Config to setup an trace exporter and returns a shutdown handler
func Setup(c Config) app.ShutdownFunc {
	var exporter trace.SpanExporter
	var err error

	exporterName := strings.ToLower(c.Exporter)
	switch exporterName {
	case "stackdriver":
		exporter, err = newStackdriverExporter(c.Stackdriver)
	case "otlp":
		exporter, err = newOTLPExporter(c.OTLP)
	case "":
		// No trace will be exported, but will be created
	default:
		err = errors.New("invalid exporter")
	}

	if err != nil {
		log.Fatal().Str("exporter", exporterName).Err(err).Msg("Failed to create trace exporter")
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
