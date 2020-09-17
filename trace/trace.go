package trace

// Trace represents the informations that should be
// passed through systems
type Trace struct {
	// ID represents the Trace ID used in logging and
	// trace views. It will be used as the main span in
	ID ID

	// ProbabilitySample represents if the span will be
	// sampled or not. The two possibles values are 0 and 1
	ProbabilitySample *float64
}

func newTrace() Trace {
	return Trace{
		ID:                NewID(),
		ProbabilitySample: &defaultProbabilitySample,
	}
}

func (t Trace) isEmpty() bool {
	return IDIsEmpty(t.ID) || t.ProbabilitySample == nil
}

func ensureTraceNotEmpty(t Trace) Trace {
	if IDIsEmpty(t.ID) {
		t.ID = NewID()
	}
	if t.ProbabilitySample == nil {
		t.ProbabilitySample = &defaultProbabilitySample
	}
	return t
}
