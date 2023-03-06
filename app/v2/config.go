package app

import (
	"time"

	"github.com/arquivei/foundationkit/app/v2/logger"

	"github.com/omeid/uconfig"
	"github.com/omeid/uconfig/plugins/defaults"
	"github.com/omeid/uconfig/plugins/env"
	"github.com/omeid/uconfig/plugins/flag"
	"github.com/rs/zerolog/log"
)

type Config struct {
	App struct {
		Log         logger.Config
		AdminServer struct {
			// Enabled sets the admin server
			Enabled bool `default:"true"`
			// DefaultAdminPort is the default port the app will bind the admin HTTP interface.
			Port string `default:"9000"`
			With struct {
				// DebugURLs sets the debug URLs in the admin server. To disable them, set to false.
				DebugURLs bool `default:"true"`
				// Metrics
				Metrics bool `default:"true"`
				// Probes
				Probes bool `default:"true"`
			}
		}
		Shutdown struct {
			// DefaultGracePeriod is the default value for the grace period.
			// During normal shutdown procedures, the shutdown function will wait
			// this amount of time before actually starting calling the shutdown handlers.
			GracePeriod time.Duration `default:"3s"`
			// DefaultShutdownTimeout is the default value for the timeout during shutdown.
			Timeout time.Duration `default:"5s"`
		}
	}
}

func (c Config) GetAppConfig() Config {
	return c
}

// setupConfig loads the configuration in the given struct. In case of error, prints help and exit application.
func SetupConfig(config any) {
	c, err := uconfig.New(config, defaults.New(), env.New(), flag.Standard())
	if err == nil {
		err = c.Parse()
	}
	if err != nil {
		c.Usage()
		log.Fatal().Err(err).Msg("[app] Failed to setup config!")
	}
}
