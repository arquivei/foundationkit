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

// Metrifier is a struct that helps metrify any computation using
// two Promethus metrics: execution_count and execution_latency_seconds.
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

// Span is used to metrify a piece of computation.
type Span struct {
	begin  time.Time
	labels map[string]string
	m      *Metrifier
}

// End ends a span and calculate and publish the metrics.
func (s *Span) End(err error) {
	labels := make([]string, 0, 2*len(s.m.labels))
	labels = append(labels, labelErrorCode, getErrorCode(err))

	if len(s.labels) > 0 {
		for k, v := range s.labels {
			labels = append(labels, k, v)
		}
	}

	s.m.latency.With(labels...).Observe(time.Since(s.begin).Seconds())
	s.m.count.With(labels...).Add(1)
}

func getErrorCode(err error) string {
	if err == nil {
		return ""
	}
	c := errors.GetCode(err).String()
	if c == "" {
		return "UNKNOWN"
	}
	return c
}

// WithLabels appends the given labels on the Prometheus metrics.
func (s *Span) WithLabels(l map[string]string) *Span {
	if s.labels != nil {
		for k, v := range l {
			s.labels[k] = v
		}
	} else {
		s.labels = l
	}
	return s
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
			Objectives:  c.Objectives,
			MaxAge:      c.MaxAge,
			AgeBuckets:  c.AgeBuckets,
			BufCap:      c.BufCap,
		}, labelKeys),
	}, nil
}
