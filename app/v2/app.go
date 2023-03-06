package app

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof" // Sadly, this also changes the DefaultMux to have the pprof URLs
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// MainLoopFunc is the function runned by app. If it finishes, it will trigger a shutdown.
// The context will be canceled when application is gracefully shutting down.
type MainLoopFunc func(context.Context) error

// App represents an application with a main loop and a shutdown routine
type App struct {
	Ready   ProbeGroup
	Healthy ProbeGroup

	mainLoopCtx       context.Context
	cancelMainLoopCtx func()

	readinessProbe   Probe
	healthinessProbe Probe

	isRunning atomic.Bool
	shutdown  struct {
		handlers    shutdownHeap
		timeout     time.Duration
		gracePeriod time.Duration
		err         error
		once        sync.Once
	}
}

// New returns a new App using the app.Config struct for configuration.
func New(c Config) *App {
	log.Trace().Msg("[app] Creating new app")

	readinessProbeGroup := NewProbeGroup("readiness")
	healthinessProbeGroup := NewProbeGroup("healthiness")

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		Ready:             readinessProbeGroup,
		Healthy:           healthinessProbeGroup,
		mainLoopCtx:       ctx,
		cancelMainLoopCtx: cancel,
		shutdown: struct {
			handlers    shutdownHeap
			timeout     time.Duration
			gracePeriod time.Duration
			err         error
			once        sync.Once
		}{
			gracePeriod: c.App.Shutdown.GracePeriod,
			timeout:     c.App.Shutdown.Timeout,
		},
		readinessProbe:   readinessProbeGroup.MustNewProbe("fkit/app", false),
		healthinessProbe: healthinessProbeGroup.MustNewProbe("fkit/app", true),
	}

	app.startAdminServer(c)

	return app
}

func (app *App) startAdminServer(c Config) {
	if !c.App.AdminServer.Enabled {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/healthy", &app.Healthy)
	mux.Handle("/ready", &app.Ready)

	if c.App.AdminServer.With.DebugURLs {
		mux.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
		mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		mux.HandleFunc("/debug/dump/goroutines", dumpGoroutines)
		mux.HandleFunc("/debug/dump/memory", dumpMemProfile)
		mux.HandleFunc("/debug/dump/memstats", dumpMemStats)
	}

	server := http.Server{
		Addr:              ":" + c.App.AdminServer.Port,
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err).Msg("[app] Failed to start admin server.")
		}
	}()
}

// Shutdown calls all shutdown methods ordered by priority.
// Handlers are processed from higher priority to lower priority.
func (app *App) Shutdown(ctx context.Context) error {
	app.shutdown.once.Do(func() {
		log.Trace().Int("shutdown_handlers", app.shutdown.handlers.Len()).Msg("[app] Starting graceful shutdown.")
		app.doShutdown(ctx)
		if app.shutdown.err != nil {
			app.shutdown.err = fmt.Errorf("app.App.Shutdown: %w", app.shutdown.err)
			log.Error().Err(app.shutdown.err).Msg("[app] Graceful shutdown failed.")
			return
		}
		log.Info().Msg("[app] Graceful shutdown successful.")
	})

	return app.shutdown.err
}

func (app *App) doShutdown(ctx context.Context) {
	app.cancelMainLoopCtx()

	if app.shutdown.timeout > 0 {
		log.Trace().Dur("shutdown_timeout", app.shutdown.timeout).Msg("[app] Configuring a timeout for the shutdown.")
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, app.shutdown.timeout)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		app.shutdown.err = ctx.Err()
	case app.shutdown.err = <-app.shutdownAllHandlers(ctx):
	}
}

func (app *App) shutdownAllHandlers(ctx context.Context) chan error {
	done := make(chan error, 1)
	go func() {
		defer close(done)
		for app.shutdown.handlers.Len() > 0 {
			h := heap.Pop(&app.shutdown.handlers).(*ShutdownHandler)
			if ctx.Err() != nil {
				done <- fmt.Errorf("shutdownAllHandlers: %w", ctx.Err())
				return
			}

			logger := log.With().
				Str("shutdown_handler_name", h.Name).
				Uint8("shutdown_handler_priority", uint8(h.Priority)).
				Dur("shutdown_handler_timeout", h.Timeout).
				Str("shutdown_handler_policy", ErrorPolicyString(h.Policy)).Logger()

			logger.Trace().Msg("[app] Executing shutdown handler.")
			err := h.Execute(ctx)
			logger.Trace().Err(err).Msg("[app] Shutdown handler finished.")
			if err != nil {
				done <- fmt.Errorf("shutdownAllHandlers: %w", err)
			}
		}
	}()
	return done
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown.
// This is expected to be called only once and will panic if called a second time.
func (app *App) RunAndWait(mainLoop MainLoopFunc) {
	if alreadyRunning := app.isRunning.Swap(true); alreadyRunning {
		panic("[app] RunAndWait called more than once")
	}

	log.Trace().Msg("[app] Starting run and wait.")

	// Run main loop on a go-routine
	errs := make(chan error, 1)
	go app.runMainLoop(mainLoop, errs)
	app.waitMainLoopOrSignal(errs)

	// App is shutting down...
	app.readinessProbe.SetNotOk()
	app.waitGracePeriod()
	_ = app.Shutdown(context.Background())
	app.waitMainLoopFinish(10 * time.Second)

	// This forces kubernetes kills the pod if some other code is holding the main func.
	app.healthinessProbe.SetNotOk()
}

func (app *App) runMainLoop(mainLoop MainLoopFunc, errs chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			errs <- recoverErr(r)
		}
	}()

	if mainLoop == nil {
		errs <- errors.New("main loop is nil")
		return
	}
	log.Info().Msg("[app] Application main loop starting now!")
	app.readinessProbe.SetOk()

	errs <- mainLoop(app.mainLoopCtx)
}

func (app *App) waitMainLoopOrSignal(errs <-chan error) {
	gracefulShutdownCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	select {
	case err := <-errs:
		if err != nil {
			log.Error().Err(err).Msg("[app] Main Loop finished by itself with error.")
		} else {
			log.Warn().Msg("[app] Main Loop finished by itself without error. Ideally the main loop should be finished by a graceful shutdown handler.")
		}
	case <-gracefulShutdownCtx.Done():
		log.Info().
			Dur("grace_period", app.shutdown.gracePeriod).
			Msg("[app] Graceful shutdown signal received.")
	}
}

func (app *App) waitGracePeriod() {
	if app.shutdown.gracePeriod <= 0 {
		return
	}

	log.Info().
		Dur("grace_period", app.shutdown.gracePeriod).
		Msg("[app] Awaiting for grace period to end.")
	time.Sleep(app.shutdown.gracePeriod)
	log.Info().Msg("[app] Grace period is over, initiating shutdown procedures...")
}

func (app *App) waitMainLoopFinish(timeout time.Duration) {
	select {
	case <-app.mainLoopCtx.Done():
	case <-time.After(timeout):
		log.Error().Msg("[app] Main loop is shutting down but the main loop didn't finish.")
	}
}

// RegisterShutdownHandler adds a shutdown handler to the app. Shutdown Handlers are executed
// one at a time from the highest priority to the lowest priority. Shutdown handlers of the same
// priority are normally executed in the added order but this is not guaranteed.
func (app *App) RegisterShutdownHandler(sh *ShutdownHandler) {
	if sh.Name == "" {
		panic("shutdown handler name must not be an empty string")
	}
	if len(app.shutdown.handlers) == 0 {
		heap.Init(&app.shutdown.handlers)
	}
	heap.Push(&app.shutdown.handlers, sh)

	log.Trace().
		Str("shutdown_handler_name", sh.Name).
		Uint8("shutdown_handler_priority", uint8(sh.Priority)).
		Dur("shutdown_handler_timeout", sh.Timeout).
		Str("shutdown_handler_policy", ErrorPolicyString(sh.Policy)).
		Msg("[app] Shutdown handler registered.")
}
