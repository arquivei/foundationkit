// This example ilustrates how shutdowns are handled.
//
// Shutdown is executed from High to Low priority.
//
// Graceful shutdown is triggered by either the main loop finishing by itself or by receiving
// a SIGINT signal (ctrl+c on the terminal).
//
// Run this example and try both scenarios. Let it finishes by itself or kill it with ctrl+c.
//
// To run this code:
// go run -ldflags="-X main.version=v0.0.1" ./app/v2/examples/shutdown-handlers/ -app-log-human -app-log-level=trace
package main

import (
	"context"
	"errors"
	"time"

	"github.com/arquivei/foundationkit/app/v2"
	"github.com/rs/zerolog/log"
)

type config struct {
	app.Config
}

var (
	cfg     config
	version = "development"
)

func main() {
	app.Bootstrap(version, &cfg)

	app.RegisterShutdownHandler(&app.ShutdownHandler{
		Name:    "first",
		Timeout: time.Second,
		Handler: func(ctx context.Context) error {
			log.Info().Msg("First shutdown handler executed.")
			return nil
		},
		Policy:   app.ErrorPolicyWarn,
		Priority: app.ShutdownPriority(30),
	})
	app.RegisterShutdownHandler(&app.ShutdownHandler{
		Name:    "second",
		Timeout: time.Second,
		Handler: func(ctx context.Context) error {
			log.Info().Msg("Second handler will fail but will only cause a warn.")
			return errors.New("some error")
		},
		Policy:   app.ErrorPolicyWarn,
		Priority: app.ShutdownPriority(20),
	})

	app.RegisterShutdownHandler(&app.ShutdownHandler{
		Name:    "third",
		Timeout: time.Second,
		Handler: func(ctx context.Context) error {
			log.Info().Msg("Third handler executed.")
			return nil
		},
		Policy:   app.ErrorPolicyWarn,
		Priority: app.ShutdownPriority(10),
	})

	app.RunAndWait(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			return errors.New("app timeout reached")
		}
	})

}
