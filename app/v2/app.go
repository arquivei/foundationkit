package app

import (
	"container/heap"
	"context"
	"net/http"
	"net/http/pprof" // Sadly, this also changes the DefaultMux to have the pprof URLs
	"os/signal"
	"syscall"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/trace"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// MainLoopFunc is the functions runned by app. If it finishes, it will trigger a shutdown.
// The context will be canceled when application is gracefully shutting down.
type MainLoopFunc func(context.Context) error

// App represents an application with a main loop and a shutdown routine
type App struct {
	Ready   ProbeGroup
	Healthy ProbeGroup

	mainLoopCtx       context.Context
	cancelMainLoopCtx func()

	shutdownHandlers shutdownHeap
	gracePeriod      time.Duration
	shutdownTimeout  time.Duration

	mainReadinessProbe  Probe
	mainHealthnessProbe Probe
}

// New returns a new App.
// If ctx contains a zerolog logger it is used for logging.
// adminPort must be a valid port number or it will fail silently.
func New(c Config) *App {
	log.Trace().Msg("[app] Creating new app")

	readinessProbeGroup := NewProbeGroup("readiness")
	healthnessProbeGroup := NewProbeGroup("healthness")

	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		Ready:               readinessProbeGroup,
		Healthy:             healthnessProbeGroup,
		mainLoopCtx:         ctx,
		cancelMainLoopCtx:   cancel,
		gracePeriod:         c.App.Shutdown.GracePeriod,
		shutdownTimeout:     c.App.Shutdown.Timeout,
		mainReadinessProbe:  readinessProbeGroup.MustNewProbe("fkit/app", false),
		mainHealthnessProbe: healthnessProbeGroup.MustNewProbe("fkit/app", true),
	}

	app.startAdminServer(c)

	if c.App.Trace.Exporter != "" {
		trace.SetupTrace(c.App.Trace)
	}

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
			log.Fatal().Err(err).Msg("Failed to start admin server.")
		}
	}()
}

// Shutdown calls all shutdown methods, in order they were added.
func (app *App) Shutdown(ctx context.Context) (err error) {
	log.Trace().Int("shutdown_handlers", app.shutdownHandlers.Len()).Msg("[app] Starting graceful shutdown.")

	app.cancelMainLoopCtx()

	const op = errors.Op("app.App.Shutdown")

	if app.shutdownTimeout > 0 {
		log.Trace().Dur("shutdown_timeout", app.shutdownTimeout).Msg("[app] Configuring a timeout for the shutdown.")
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, app.shutdownTimeout)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		err = errors.E(op, "shutdown deadline has been reached")
	case err = <-app.shutdownAllHandlers(ctx):
	}
	if err != nil {
		log.Trace().Err(err).Msg("[app] Graceful shutdown failed.")
		return errors.E(op, err)
	}

	log.Trace().Msg("[app] Graceful shutdown finished successfully.")
	return nil
}

func (app *App) shutdownAllHandlers(ctx context.Context) chan error {
	const op = errors.Op("shutdownAllHandlers")
	done := make(chan error, 1)
	go func() {
		defer close(done)
		for app.shutdownHandlers.Len() > 0 {
			h := heap.Pop(&app.shutdownHandlers).(*ShutdownHandler)
			if ctx.Err() != nil {
				done <- errors.E(op, "shutdown deadline has been reached")
			}

			logger := log.With().
				Str("shutdown_handler_name", h.Name).
				Uint8("shutdown_handler_priority", uint8(h.Priority)).
				Dur("shutdown_handler_timeout", h.Timeout).
				Str("shutdown_handler_policy", ErrorPolicyString(h.Policy)).Logger()

			logger.Trace().Msg("[app] Executing shutdown handler.")
			if err := h.Execute(ctx); err != nil {
				logger.Trace().Msg("[app] Shutdown handler failed.")
				done <- errors.E(op, err)
			}
			logger.Trace().Msg("[app] Shutdown handler finished.")
		}
	}()
	return done
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown
func (app *App) RunAndWait(mainLoop MainLoopFunc) {
	log.Trace().Msg("[app] Starting run and wait.")

	errs := make(chan error, 1)

	go app.runMainLoop(mainLoop, errs)

	notifyCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Await for OS signal or main loop to finishes by itself
	select {
	case err := <-errs:
		if err != nil {
			log.Error().Err(err).Msg("[app] Main Loop finished by itself with error.")
		} else {
			log.Warn().Msg("[app] Main Loop finished by itself without error. Ideally the main loop should be finished by a graceful shutdown handler.")
		}
	case <-notifyCtx.Done():
		log.Info().
			Dur("grace_period", app.gracePeriod).
			Msg("[app] Graceful shutdown signal received.")
	}

	app.mainReadinessProbe.SetNotOk()
	log.Info().
		Dur("grace_period", app.gracePeriod).
		Msg("[app] Awaiting for grace period to end.")
	time.Sleep(app.gracePeriod)
	log.Info().Msg("[app] Grace period is over, initiating shutdown procedures...")

	app.logAppTerminated(app.Shutdown(context.Background()))

	select {
	case <-app.mainLoopCtx.Done():
	case <-time.After(10 * time.Second):
		log.Error().Msg("[app] Main loop didn't finished by itself.")
	}

	// This forces kubernetes kills the pod if some other code is holding the main func.
	app.mainHealthnessProbe.SetNotOk()
}

func (app *App) runMainLoop(mainLoop MainLoopFunc, errs chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			errs <- errors.NewFromRecover(r)
		}
	}()

	if mainLoop == nil {
		errs <- errors.New("main loop is nil")
		return
	}
	log.Info().Msg("[app] Application main loop starting now!")
	app.mainReadinessProbe.SetOk()

	errs <- mainLoop(app.mainLoopCtx)
}

func (app *App) logAppTerminated(err error) {
	if err == nil {
		log.Info().Msg("[app] App gracefully terminated.")
	} else {
		log.Error().Err(err).Msg("[app] App terminated with error.")
	}
}

// RegisterShutdownHandler adds a shutdown handler to the app. Shutdown Handlers are executed
// one at a time from the highest priority to the lowest priority. Shutdown handlers of the same
// priority are normally executed in the added order but this is not guaranteed.
func (app *App) RegisterShutdownHandler(sh *ShutdownHandler) {
	if sh.Name == "" {
		panic("shutdown handler name must not be an empty string")
	}
	if len(app.shutdownHandlers) == 0 {
		heap.Init(&app.shutdownHandlers)
	}
	heap.Push(&app.shutdownHandlers, sh)

	log.Trace().
		Str("shutdown_handler_name", sh.Name).
		Uint8("shutdown_handler_priority", uint8(sh.Priority)).
		Dur("shutdown_handler_timeout", sh.Timeout).
		Str("shutdown_handler_policy", ErrorPolicyString(sh.Policy)).
		Msg("[app] Shutdown handler registered")
}
