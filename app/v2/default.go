package app

import (
	"context"

	"github.com/arquivei/foundationkit/app/v2/logger"

	"github.com/rs/zerolog/log"
)

var (
	// This is the default app.
	defaultApp *App
)

type AppConfig interface {
	GetAppConfig() Config
}

// Bootstrap initializes the config structure, the log and creates a new app internally.
func Bootstrap(appVersion string, config AppConfig) {
	SetupConfig(config)
	appConfig := config.GetAppConfig()
	logger.Setup(appConfig.App.Log, appVersion)

	log.Info().Str("config", logger.Flatten(config)).Msg("[app] Configuration loaded and global logger configured.")
	defaultApp = New(appConfig)
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown.
// This is expected to be called only once and will panic if called a second time.
func RunAndWait(f MainLoopFunc) {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RunAndWait(f)
}

// Shutdown calls all shutdown methods ordered by priority.
// Handlers are processed from higher priority to lower priority.
func Shutdown(ctx context.Context) error {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return defaultApp.Shutdown(ctx)
}

// RegisterShutdownHandler adds a shutdown handler to the app. Shutdown Handlers are executed
// one at a time from the highest priority to the lowest priority. Shutdown handlers of the same
// priority are normally executed in the added order but this is not guaranteed.
func RegisterShutdownHandler(sh *ShutdownHandler) {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RegisterShutdownHandler(sh)
}

// ReadinessProbeGoup is a collection of readiness probes.
func ReadinessProbeGoup() *ProbeGroup {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return &defaultApp.Ready
}

// HealthinessProbeGroup is a colection of healthiness probes.
func HealthinessProbeGroup() *ProbeGroup {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return &defaultApp.Healthy
}
