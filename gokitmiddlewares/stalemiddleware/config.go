package stalemiddleware

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Timeout time.Duration
	Logger  *zerolog.Logger
}

var DefaultConfig = Config{
	Timeout: 60 * time.Second,
	Logger:  &log.Logger,
}
