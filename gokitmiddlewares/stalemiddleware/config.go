package stalemiddleware

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Logger                 *zerolog.Logger
	MaxTimeBetweenRequests time.Duration
	StartCheckAfter        time.Duration
}

var DefaultConfig = Config{
	MaxTimeBetweenRequests: time.Minute,
	Logger:                 &log.Logger,
	StartCheckAfter:        10 * time.Second,
}
