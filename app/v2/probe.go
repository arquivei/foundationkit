package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
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

// ProbeGroup aggregates and manages probes. Probes inside a group must have a unique name.
type ProbeGroup struct {
	name   string
	lock   *sync.RWMutex
	probes map[string]*bool
}

// NewProbeGroup returns a new ProbeGroup.
func NewProbeGroup(name string) ProbeGroup {
	return ProbeGroup{
		name:   name,
		lock:   &sync.RWMutex{},
		probes: make(map[string]*bool),
	}
}

// NewProbe returns a new Probe with the given name.
// The name must satisfy the following regular expression: [a-zA-Z0-9_/-]{3,}
func (g *ProbeGroup) NewProbe(name string, ok bool) (Probe, error) {
	if err := g.checkName(name); err != nil {
		return Probe{}, err
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if err := g.checkProbeAlreadyExists(name); err != nil {
		return Probe{}, err
	}

	return g.newProbe(name, ok)
}

func (g *ProbeGroup) checkProbeAlreadyExists(name string) error {
	if _, ok := g.probes[name]; ok {
		return fmt.Errorf("probe '%s' already registered", name)
	}
	return nil
}

func (g *ProbeGroup) newProbe(name string, ok bool) (Probe, error) {
	g.probes[name] = &ok
	return Probe{&ok}, nil
}

// MustNewProbe returns a new Probe with the given name and panics in case of error.
func (g *ProbeGroup) MustNewProbe(name string, ok bool) Probe {
	p, err := g.NewProbe(name, ok)
	if err != nil {
		panic(err)
	}
	return p
}

// ServeHTTP serves the Probe Group as an HTTP handler.
func (g *ProbeGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ok, cause := g.CheckProbes()
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	//nolint:errcheck
	w.Write([]byte(cause))
	log.Trace().Str("probe_group", g.name).Bool("probe_is_ok", ok).Str("cause", cause).Msg("[app] App was probed.")
}

var reIsValidProbeName = regexp.MustCompile("[a-zA-Z0-9_/-]{3,}")

func (ProbeGroup) checkName(name string) error {
	if !reIsValidProbeName.MatchString(name) {
		return fmt.Errorf("name '%s' doesn't conform to '[a-zA-Z0-9_-]{3,}'", name)
	}
	return nil
}

// CheckProbes range through the probes and returns the state of the group.
// If any probe is not ok (false) the group state is not ok (false).
// If the group is not ok, it's also returned the cause in the second return parameter.
// If all probes are ok (true) the cause is returned as OK.
// If more than one probe is not ok, the causes are concatenated by a comma.
func (g *ProbeGroup) CheckProbes() (bool, string) {
	cause := strings.Builder{}

	cause.Grow(len(g.name) + 3)
	cause.WriteString(g.name)
	cause.WriteString(":")

	g.lock.RLock()
	defer g.lock.RUnlock()

	ok := true
	for name, probeOk := range g.probes {
		if !*probeOk {
			if !ok { // There is already a probe that is not ok
				cause.WriteString(",")
			}
			cause.WriteString(name)
			ok = false
		}
	}
	if ok {
		cause.WriteString("OK")
	}

	return ok, cause.String()
}
