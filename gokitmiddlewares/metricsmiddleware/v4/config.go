package metricsmiddleware

import (
	"github.com/arquivei/foundationkit/metrifier"
)

// Config is used to configure a metrics middleware.
type Config struct {
	// Metrifier is the metrifier configuration
	Metrifier metrifier.Config

	// LabelsDecoder extracts labels from the request, response or error.
	// This is optional, can be nil
	LabelsDecoder LabelsDecoder

	// ExternalMetrics is executed after the main metrifier is called.
	// This is intended to calculate custom metrics.
	// This is optional, can be nil.
	ExternalMetrics ExternalMetrics
}

// WithLabelsDecoder adds a LabelsDecoder to the metrics middleware.
func (c Config) WithLabelsDecoder(d LabelsDecoder) Config {
	c.LabelsDecoder = d
	c.Metrifier.ExtraLabels = d.Labels()
	return c
}

// WithExternalMetrics adds ExternalMetrics to the metrics middleware.
func (c Config) WithExternalMetrics(m ExternalMetrics) Config {
	c.ExternalMetrics = m
	return c
}

// NewDefaultConfig returns a new Config with sane defaults.
func NewDefaultConfig(endpoint string) Config {
	config := Config{
		Metrifier: metrifier.NewDefaultConfig("endpoint", ""),
	}
	config.Metrifier.ConstLabels = map[string]string{
		"fikt_endpoint": endpoint,
	}
	return config
}
