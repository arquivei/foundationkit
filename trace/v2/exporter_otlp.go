package trace

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

// OTLPConfig allows for configuring the OpenTelemery Protocol.
type OTLPConfig struct {
	// Endpoint allows one to set the address of the collector
	// endpoint that the driver will use to send spans.
	// It is a string in the form <host>:<port>.
	// If unset, it will instead try to use the default endpoint
	// from package otlptracehttp (at the time of this writing
	// the default is localhost:4318).
	// Note that the endpoint must not contain any URL path.
	Endpoint string
	// Compression tells the driver to compress the sent data.
	// Possible values are:
	// - "" (empty string) or "none": No compression
	// - "gzip": GZIP compression.
	Compression string
	// Insecure tells the driver to connect to the collector using the
	// HTTP scheme, instead of HTTPS.
	Insecure bool
}

func (c *OTLPConfig) Options() []otlptracehttp.Option {
	opts := make([]otlptracehttp.Option, 0, 10)

	if c.Endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(c.Endpoint))
	}

	switch c.Compression {
	case "", "none":
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.NoCompression))
	case "gzip":
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	default:
		log.Fatal().Msgf("Invalid OpenTelemetry trace exporter compression: %s", c.Compression)
	}

	if c.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	return opts
}

func newOTLPExporter(c OTLPConfig) (trace.SpanExporter, error) {
	return otlptracehttp.New(context.Background(), c.Options()...)
}
