package app

import (
	"context"
	"time"
)

var (
	// DefaultGracePeriod is the default value for the grace period.
	// During normal shutdown procedures, the shutdown function will wait
	// this amount of time before actually starting calling the shutdown handlers.
	DefaultGracePeriod = 3 * time.Second
	// DefaultShutdownTimeout is the default value for the timeout during shutdown.
	DefaultShutdownTimeout = 5 * time.Second
	// DefaultAdminPort is the default port the app will bind the admin HTTP interface.
	DefaultAdminPort = "9000"
	// This is the default app.
	defaultApp *App
)

// NewDefaultApp creates and sets the default app. The default app is controlled by
// public functions in app package
func NewDefaultApp(ctx context.Context) (err error) {
	defaultApp, err = New(ctx, DefaultAdminPort)
	if err != nil {
		return err
	}
	defaultApp.GracePeriod = DefaultGracePeriod
	defaultApp.ShutdownTimeout = DefaultShutdownTimeout
	return nil
}

// RunAndWait calls the RunAndWait of the default app
func RunAndWait(f MainLoopFunc) {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RunAndWait(f)
}

// Shutdown calls the Shutdown of the default app
func Shutdown(ctx context.Context) error {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return defaultApp.Shutdown(ctx)
}

// RegisterShutdownHandler calls the RegisterShutdownHandler from the default app
func RegisterShutdownHandler(sh *ShutdownHandler) {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RegisterShutdownHandler(sh)
}

// ReadinessProbeGoup TODO
func ReadinessProbeGoup() *ProbeGroup {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return &defaultApp.Ready
}

// HealthinessProbeGroup TODO
func HealthinessProbeGroup() *ProbeGroup {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return &defaultApp.Healthy
}
