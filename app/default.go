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
	// This is the default app.
	defaultApp *App
)

// NewDefaultApp creates and sets the defualt app. The default app is controlled by
// public functions in app package
func NewDefaultApp(ctx context.Context, mainLoop MainLoopFunc) (err error) {
	defaultApp, err = New(ctx, mainLoop)
	if err != nil {
		return err
	}
	defaultApp.GracePeriod = DefaultGracePeriod
	defaultApp.ShutdownTimeout = DefaultShutdownTimeout
	return nil
}

// RunAndWait calls the RunAndWait of the default app
func RunAndWait() error {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return defaultApp.RunAndWait()
}

// Shutdown calls the Shutdown of the default app
func Shutdown(ctx context.Context) error {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	return defaultApp.Shutdown(ctx)
}

// RegisterShutdownHandler calls the RegisterShutdownHandler from the default app
func RegisterShutdownHandler(name string, fn ShutdownFunc, options ...interface{}) {
	if defaultApp == nil {
		panic("default app not initialized")
	}
	defaultApp.RegisterShutdownHandler(name, fn, options...)
}
