package metricsmiddleware

import (
	"github.com/arquivei/foundationkit/metrifier"
)

// Config is used to configure a metrics middleware.
type Config struct {
	// Metrifier is the metrifier configuration
	Metrifier metrifier.Config

	// LabelDecoder defines decoder for labels. The Labels must also be declared
	// in the Metrifier.ExtraLabels.
	LabelsDecoder LabelsDecoder
}

// NewDefaultConfig returns a new Config with sane defaults.
func NewDefaultConfig(system, subsystem string) Config {
	return Config{
		Metrifier: metrifier.NewDefaultConfig(system, subsystem),
	}
}
