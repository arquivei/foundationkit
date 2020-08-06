package metrifier

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Config is used to configure a Metrifier.
type Config struct {
	// System identifies the application being instrumented.
	//
	// This field is mandatory.
	//
	// If you have multiple deployments consider adding the environment here.
	// Examples: myserviceprod, myservicetesting
	System string

	// Subsystem identifies which part of the code is being instrumented.
	//
	// This field is mandatory.
	//
	// Examples: myserviceendpoint, myentityrepository, myjobqueue
	Subsystem string

	// ExtraLabels are extra prometheus labels that will be attached to the metric.
	//
	// The Metrifier will always create and manage the "error_code" label
	// automatically. The ExtraLabels field is for any extra label that
	// will be (and must be) passed by calling the Span.WithLabels.
	ExtraLabels []string

	// The bellow fields are passed directly to Prometheus.

	// ConstLabels are used to attach fixed labels to this metric. Metrics
	// with the same fully-qualified name must have the same label names in
	// their ConstLabels.
	//
	// Due to the way a Summary is represented in the Prometheus text format
	// and how it is handled by the Prometheus server internally, “quantile”
	// is an illegal label name. Construction of a Summary or SummaryVec
	// will panic if this label name is used in ConstLabels.
	//
	// ConstLabels are only used rarely. In particular, do not use them to
	// attach the same labels to all your metrics. Those use cases are
	// better covered by target labels set by the scraping Prometheus
	// server, or by one specific metric (e.g. a build_info or a
	// machine_role metric). See also
	// https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels,-not-static-scraped-labels
	ConstLabels prometheus.Labels

	// Summary contains only Summary specific configurations
	Summary ConfigSummary
}

// ConfigSummary contains some Summary specific configuration
type ConfigSummary struct {
	// Objectives defines the quantile rank estimates with their respective
	// absolute error. If Objectives[q] = e, then the value reported for q
	// will be the φ-quantile value for some φ between q-e and q+e.  The
	// default value is an empty map, resulting in a summary without
	// quantiles.
	Objectives map[float64]float64

	// MaxAge defines the duration for which an observation stays relevant
	// for the summary. Must be positive. The default value is DefMaxAge.
	MaxAge time.Duration

	// AgeBuckets is the number of buckets used to exclude observations that
	// are older than MaxAge from the summary. A higher number has a
	// resource penalty, so only increase it if the higher resolution is
	// really required. For very high observation rates, you might want to
	// reduce the number of age buckets. With only one age bucket, you will
	// effectively see a complete reset of the summary each time MaxAge has
	// passed. The default value is DefAgeBuckets.
	AgeBuckets uint32

	// BufCap defines the default sample stream buffer size.  The default
	// value of DefBufCap should suffice for most uses. If there is a need
	// to increase the value, a multiple of 500 is recommended (because that
	// is the internal buffer size of the underlying package
	// "github.com/bmizerany/perks/quantile").
	BufCap uint32
}

// NewDefaultConfig returns the Config struct filled with sane defaults.
func NewDefaultConfig(system, subsystem string) Config {
	return Config{
		System:    system,
		Subsystem: subsystem,
		Summary: ConfigSummary{
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	}
}
