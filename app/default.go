package app

import (
	"context"
	"time"
)

var (
	// DefaultGracePeriod is the default value for the grace period.
	// During normal shutdown procedures, the shutdown function will wait
	// this amount of time before actually starting calling the shurdown handlers.
	DefaultGracePeriod = 3 * time.Second
	// DefaultShutdownTimeout is the defualt value for the timeout during shutdown.
	DefaultShutdownTimeout = 5 * time.Second
	// DefaultAdminPort is the default port the app will bind the admin HTTP interface.
	DefaultAdminPort = "9000"
	// This is the default app.
	defaultApp *App
)

// NewDefaultApp creates and sets the defualt app. The default app is controlled by
// public functions in app package
func NewDefaultApp(ctx context.Context, mainLoop MainLoopFunc) (err error) {
	defaultApp, err = New(ctx, DefaultAdminPort, mainLoop)
	if err != nil {
		return err
	}
	defaultApp.GracePeriod = DefaultGracePeriod
	defaultApp.ShutdownTimeout = DefaultShutdownTimeout
	return nil
}

// RunAndWait calls the RunAndWait of the default app
func RunAndWait() {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RunAndWait()
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

// IsReady returns if the default app is ready (to be used by kubernetes readyness probe)
func IsReady() bool {
	return defaultApp.Ready
}

// SetReady sets the app to ready state
func SetReady() {
	defaultApp.Ready = true
}

// SetUnready sets the app to unready state
func SetUnready() {
	defaultApp.Ready = false
}

// IsHealthy returns if the app is healthy. Unhealthy apps are killed by the kubernetes.
func IsHealthy() bool {
	return defaultApp.Healthy
}

// SetHealthy sets the app to an healthy state
func SetHealthy() {
	defaultApp.Healthy = true
}

// SetUnhealthy sets the app to an unhealthy state
func SetUnhealthy() {
	defaultApp.Healthy = false
}
