package app

import (
	"context"
	"testing"
	"time"

	"github.com/arquivei/foundationkit/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

func newAppTestingConfig() Config {
	cfg := Config{}
	cfg.App.Log = log.Config{
		Level: "disabled",
	}
	cfg.App.AdminServer.Enabled = false
	cfg.App.Shutdown.GracePeriod = 3 * time.Second
	cfg.App.Shutdown.Timeout = 5 * time.Second
	return cfg
}

func TestRunAndWait(t *testing.T) {
	assert.Panics(t, func() {
		a := App{}
		main := func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		}
		a.RunAndWait(main)
		a.RunAndWait(main)
	}, "Panics if RunAndWait is called more than once.")

	assert.NotPanics(t, func() {
		a := New(newAppTestingConfig())
		a.RunAndWait(func(ctx context.Context) error {
			return nil
		})

	}, "Calling RunAndWait once should not Panic.")
}

func TestAppShutdown(t *testing.T) {
	assert.NotPanics(t, func() {
		var shutdownHandlerCalled bool
		a := New(Config{})
		a.RegisterShutdownHandler(&ShutdownHandler{
			Name: "testing_handler",
			Handler: func(ctx context.Context) error {
				shutdownHandlerCalled = true
				return nil
			},
		})
		err := a.Shutdown(context.Background())
		assert.NoError(t, err, "Shutdown should not fail.")
		assert.True(t, shutdownHandlerCalled, "Shutdown handler should be executed during shutdown.")
	}, "Calling RunAndWait once should not Panic")
}
