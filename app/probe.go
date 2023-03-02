package app

import (
	"regexp"
	"strings"

	"github.com/arquivei/foundationkit/errors"
)

// Probe stores the state of a probe (`true` or `false` for ok and not ok respectively).
type Probe struct {
	ok *bool
}

// Set changes the state of a probe. Use `true` for ok and `false` for not ok.
func (p *Probe) Set(ok bool) {
	*p.ok = ok
}

// SetOk sets the probe as ok. Same as `Set(true)`.
func (p *Probe) SetOk() {
	*p.ok = true
}

// SetNotOk sets the probe as not ok. Same as `Set(false)`.
func (p *Probe) SetNotOk() {
	*p.ok = false
}

// IsOk returns the state of the probe  (`true` or `false` for ok and not ok respectively).
func (p *Probe) IsOk() bool {
	return *p.ok
}

// ProbeGroup aggregates and manages probes.
type ProbeGroup struct {
	probes map[string]*bool
}

// NewProbeGroup returns a new ProbeGroup.
func NewProbeGroup() ProbeGroup {
	return ProbeGroup{
		probes: make(map[string]*bool),
	}
}

// NewProbe returns a new Probe with the given name.
func (m *ProbeGroup) NewProbe(name string, ok bool) (Probe, error) {
	if _, ok := m.probes[name]; ok {
		return Probe{}, errors.Errorf("probe '%s' already registered", name)
	}

	if err := m.checkName(name); err != nil {
		return Probe{}, err
	}

	s := ok
	m.probes[name] = &s
	return Probe{&s}, nil
}

// MustNewProbe returns a new Probe with the given name and panics in case of error.
func (m *ProbeGroup) MustNewProbe(name string, ok bool) Probe {
	p, err := m.NewProbe(name, ok)
	if err != nil {
		panic(err)
	}
	return p
}

var reIsValidProbeName = regexp.MustCompile("[a-zA-Z0-9_/-]{3,}")

func (ProbeGroup) checkName(name string) error {
	if !reIsValidProbeName.MatchString(name) {
		return errors.Errorf("name '%s' doesn't conform to '[a-zA-Z0-9_-]{3,}'", name)
	}
	return nil
}

// CheckProbes range through the probes and returns the state of the group.
// If any probe is not ok (false) the group state is not ok (false).
// If the group is not ok, it's also returned the cause in the second return parameter.
// If more than one probe is not ok, the causes are concatenated by a comma.
func (m *ProbeGroup) CheckProbes() (bool, string) {
	ok := true
	cause := strings.Builder{}

	for name, probeOk := range m.probes {
		if !*probeOk {
			ok = false
			if cause.Len() > 0 {
				cause.WriteString(",")
			}
			cause.WriteString(name)
		}
	}

	return ok, cause.String()
}
