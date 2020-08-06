package metrifier

import (
	"time"

	"github.com/arquivei/foundationkit/errors"
)

// Span is used to metrify a piece of computation.
type Span struct {
	begin  time.Time
	labels map[string]string
	m      *Metrifier
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

// End ends a span and calculate and publish the metrics.
func (s *Span) End(err error) {
	labels := make([]string, 0, 2*len(s.m.labels))
	labels = append(labels, labelErrorCode, getErrorCode(err))
	for k, v := range s.labels {
		labels = append(labels, k, v)
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
