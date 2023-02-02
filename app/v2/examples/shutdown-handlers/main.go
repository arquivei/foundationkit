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
			time.Sleep(time.Second)
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
			time.Sleep(time.Second)
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
			time.Sleep(time.Second)
			log.Info().Msg("Third handler executed.")
			return nil
		},
		Policy:   app.ErrorPolicyWarn,
		Priority: app.ShutdownPriority(10),
	})

	// We will trigger a second call to Shutdown just to show how it is handled.
	// You should not do this in production code.
	secondaryShutdown, secondaryShutdownCancel := context.WithCancel(context.Background())
	defer secondaryShutdownCancel()

	app.RunAndWait(func(ctx context.Context) error {
		// triggering the second shutdown after the functions ends and on a goroutine.
		defer func() {
			go func() {
				log.Info().Msg("Triggering another Shutdown. This should emit a warning but it will wait and yield the same result as the original Shutdown call.")
				defer secondaryShutdownCancel()
				err := app.Shutdown(context.Background())
				log.Warn().Err(err).Msg("Second Shutdown finished.")
			}()
		}()

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			return errors.New("app timeout reached")
		}
	})
	// Just wait for the secondary shutdown to finish so we can read it in the logs.
	<-secondaryShutdown.Done()
}
