package apiping

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// Request Request
type Request struct {
	Num   int           `json:"num"`
	Sleep time.Duration `json:"sleep"`
}

// TraceAttributes will inject attributes num and sleep into the trace span.
func (r Request) TraceAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Int("request.num", r.Num),
		attribute.Float64("request.sleep.seconds", r.Sleep.Seconds()),
	}
}
