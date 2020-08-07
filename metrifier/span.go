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
func (s Span) WithLabels(l map[string]string) Span {
	if len(l) == 0 {
		return s
	}

	return Span{
		begin:  s.begin,
		m:      s.m,
		labels: mergeMaps(s.labels, l),
	}
}

func mergeMaps(maps ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}

// End ends a span and calculate and publish the metrics.
func (s Span) End(err error) {
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
