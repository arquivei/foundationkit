package trace

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arquivei/foundationkit/app"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Config represents all trace's configuration
type Config struct {
	ServiceName       string
	ServiceVersion    string
	Exporter          string
	ProbabilitySample float64
	Stackdriver       StackdriverConfig
	OTLP              OTLPConfig
}

// Setup use Config to setup an trace exporter and returns a shutdown handler
func Setup(c Config) app.ShutdownFunc {
	exporter, err := newExporter(c)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create trace exporter")
	}

	res, err := newResource(c.ServiceName, c.ServiceVersion)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create trace resource")
	}

	tp := newTraceProvider(res, c.ProbabilitySample, exporter)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(newPropagator())

	return tp.ForceFlush
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(res *resource.Resource, prob float64, exporter trace.SpanExporter) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSampler(
			trace.ParentBased(trace.TraceIDRatioBased(prob)),
		),
		trace.WithBatcher(exporter),
	)
}

func newExporter(c Config) (trace.SpanExporter, error) {
	var exporter trace.SpanExporter
	var err error

	switch exporterName := strings.ToLower(c.Exporter); exporterName {
	case "stackdriver":
		exporter, err = newStackdriverExporter(c.Stackdriver)
	case "otlp":
		exporter, err = newOTLPExporter(c.OTLP)
	case "":
		// No trace will be exported, but will be created
	default:
		err = errors.New("invalid exporter name: " + exporterName)
	}

	if err != nil {
		return nil, fmt.Errorf("creating trace exporter: %w", err)
	}

	return exporter, nil
}
