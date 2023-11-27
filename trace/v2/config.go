package trace

import (
	"context"
	"os"

	"github.com/arquivei/foundationkit/app"

	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Setup use Config to setup an trace exporter and returns a shutdown handler
func Setup(ctx context.Context) app.ShutdownFunc {
	lintOtelEnvVariables()

	// Set the OpenTelemetry to use foundation's kit default logger.
	otel.SetLogger(zerologr.New(&log.Logger))

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Warn().Err(err).Msg("[foundationkit:trace/v2] OpenTelemetry raised an error!")
	}))

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("[foundationkit:trace/v2] Failed to create trace exporter")
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exporter))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(newPropagator())

	return tp.ForceFlush
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(exporter trace.SpanExporter) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
}

// lintOtelEnvVariables logs a warning if an important variable is empty
// If any of these variables are empty, the OpenTelemetry SDK may not
// work as expected. This will not break the code, but could lead
// to loses of traces.
//
// Linted variables:
//   - OTEL_SERVICE_NAME: Because defaults to 'unkown_service' and this says nothing about the service.
//   - OTEL_EXPORTER_OTLP_ENDPOINT: because it defaults to localhost and could lose exported traces.
//
// Variables that are probably OK being empty:
//   - OTEL_TRACES_SAMPLER_ARG - It defaults to "1.0"
//   - OTEL_TRACES_SAMPLER - Ut defaults to "parentbased_always_on"
//
// Other variables don't seems to make too much of a difference.
func lintOtelEnvVariables() {
	for _, env := range []string{
		"OTEL_SERVICE_NAME",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	} {
		lintEnvVariable(env)
	}
}

func lintEnvVariable(env string) {
	if os.Getenv(env) == "" {
		log.Warn().Str("env", env).Msg("[foundationkit:trace/v2] An OpenTelemetry environment variable is empty but it probably shouldn't be. OpenTelemetry will use it's default, but this can lead to traces not being exported or recorded with a wrong name. Please read the docs at https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/.")
	}
}
