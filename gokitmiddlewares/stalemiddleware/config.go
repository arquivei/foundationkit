package stalemiddleware

import (
	"time"

	"github.com/arquivei/foundationkit/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is a struct used to configure a new stale middleware.
type Config struct {
	Logger                 *zerolog.Logger
	MaxTimeBetweenRequests time.Duration
	StartCheckAfter        time.Duration
	HealthinessPobe        app.Probe
}

// NewDefaultConfig returns a new `Config` with all values filled with a sane default.
func NewDefaultConfig(pg *app.ProbeGroup) Config {
	probe, err := pg.NewProbe("fkit/stale", true)
	if err != nil {
		panic(err)
	}
	return Config{
		MaxTimeBetweenRequests: time.Minute,
		Logger:                 &log.Logger,
		StartCheckAfter:        10 * time.Second,
		HealthinessPobe:        probe,
	}
}
