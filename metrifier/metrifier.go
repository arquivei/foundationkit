package metrifier

import (
	"time"

	"github.com/arquivei/foundationkit/errors"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	labelErrorCode = "error_code"
)

// Metrifier is a struct that helps metrify any computation
//
// It uses two Prometheus metrics: <system>_<subsystem>_execution_count and <system>_<subsystem>_execution_latency_seconds.
// Only one metrifier per <system>_<subsystem> is allowed, Prometheus will panic if
// it tries ro register the same metrics twice.
type Metrifier struct {
	count   *kitprometheus.Counter
	latency *kitprometheus.Summary
	labels  []string
}

// Begin creates and returns a Span.
func (m *Metrifier) Begin() Span {
	return Span{
		begin: time.Now(),
		m:     m,
	}
}

// MustNew returns a new Metrifier but panics in case of error.
func MustNew(c Config) Metrifier {
	m, err := New(c)
	if err != nil {
		panic(err)
	}
	return m
}

// New returns a new Metrifier.
func New(c Config) (Metrifier, error) {
	if c.System == "" {
		return Metrifier{}, errors.New("System is empty")
	}
	if c.Subsystem == "" {
		return Metrifier{}, errors.New("Subsystem is empty")
	}

	labelKeys := []string{labelErrorCode}
	if len(c.ExtraLabels) > 0 {
		labelKeys = append(labelKeys, c.ExtraLabels...)
	}

	return Metrifier{
		count: kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace:   c.System,
			Subsystem:   c.Subsystem,
			Name:        "execution_count",
			Help:        "Total amount of executions.",
			ConstLabels: c.ConstLabels,
		}, labelKeys),
		latency: kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace:   c.System,
			Subsystem:   c.Subsystem,
			Name:        "execution_latency_seconds",
			Help:        "Total duration of execution in seconds.",
			ConstLabels: c.ConstLabels,
			Objectives:  c.Summary.Objectives,
			MaxAge:      c.Summary.MaxAge,
			AgeBuckets:  c.Summary.AgeBuckets,
			BufCap:      c.Summary.BufCap,
		}, labelKeys),
	}, nil
}
