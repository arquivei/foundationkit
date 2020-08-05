package metricsmiddleware

import "github.com/arquivei/foundationkit/metrifier"

// Config is used to configure a metrics middleware.
type Config struct {
	// Metrifier is the metrifier configuration
	Metrifier metrifier.Config

	// LabelDecoders defines decoder for labels. The Labels must also be declared
	// in the Metrifier.ExtraLabels.
	LabelDecoders map[string]LabelDecoder
}

// NewDefaultConfig returns a new Config with sane defaults.
func NewDefaultConfig(system, subsystem string) Config {
	return Config{
		Metrifier: metrifier.NewDefaultConfig(system, subsystem),
	}
}
