package trace

import (
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/rs/zerolog/log"
	"go.opencensus.io/trace"
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
		stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID:            c.Stackdriver.ProjectID,
			BundleCountThreshold: 100,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create stackdriver trace exporter")
		}
		trace.RegisterExporter(stackdriverExporter)
		log.Info().Dur("took", time.Since(start)).Msg("Stackdriver loaded")
	case "":
	default:
		log.Fatal().Str("exporter", exporter).Msg("This exporter is not supported")
	}
	defaultProbabilitySample = c.ProbabilitySample
}
