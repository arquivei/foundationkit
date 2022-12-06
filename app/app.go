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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MainLoopFunc is the functions runned by app. If it finishes, it will trigger a shutdown
type MainLoopFunc func() error

// App represents an application with a main loop and a shutdown routine
type App struct {
	logger *zerolog.Logger

	Ready   ProbeGroup
	Healthy ProbeGroup

	shutdownHandlers shutdownHeap
	GracePeriod      time.Duration
	ShutdownTimeout  time.Duration

	mainReadinessProbe  Probe
	mainHealthnessProbe Probe
}

// New returns a new App.
// If ctx contains a zerolog logger it is used for logging.
// adminPort must be a valid port number or it will fail silently.
func New(ctx context.Context, adminPort string) (*App, error) {
	log.Trace().Msg("[app] Creating new app")

	app := &App{
		logger:  log.Ctx(ctx),
		Ready:   NewProbeGroup(),
		Healthy: NewProbeGroup(),
	}

	mainReadinessProbe, err := app.Ready.NewProbe("fkit/app", false)
	if err != nil {
		return nil, err
	}

	mainHealthnessProbe, err := app.Healthy.NewProbe("fkit/app", true)
	if err != nil {
		return nil, err
	}

	app.mainReadinessProbe = mainReadinessProbe
	app.mainHealthnessProbe = mainHealthnessProbe

	{ // This spwans an admin HTTP server for this
		mux := http.NewServeMux()

		mux.Handle("/metrics", promhttp.Handler())

		mux.HandleFunc("/healthy", func(w http.ResponseWriter, _ *http.Request) {
			isHealthy, cause := app.Healthy.CheckProbes()
			if isHealthy {
				w.WriteHeader(http.StatusOK)
				//nolint:errcheck
				w.Write([]byte("OK"))
				log.Trace().Msg("[app] Healthiness probe replied: I'm healthy!")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				//nolint:errcheck
				w.Write([]byte(cause))
				log.Trace().Str("cause", cause).Msg("[app] Healthiness probe replied: I'm unhealthy.")
			}
		})

		mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
			isReady, cause := app.Ready.CheckProbes()
			if isReady {
				w.WriteHeader(http.StatusOK)
				//nolint:errcheck
				w.Write([]byte("OK"))
				log.Trace().Msg("[app] Readiness probe replied: I'm ready!")
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				//nolint:errcheck
				w.Write([]byte(cause))
				log.Trace().Str("cause", cause).Msg("[app] Readiness probe replied: I'm unready.")
			}
		})

		mux.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
		mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		mux.HandleFunc("/debug/dump/goroutines", dumpGoroutines)
		mux.HandleFunc("/debug/dump/memory", dumpMemProfile)
		mux.HandleFunc("/debug/dump/memstats", dumpMemStats)

		server := http.Server{
			Addr:              ":" + adminPort,
			Handler:           mux,
			ReadHeaderTimeout: 60 * time.Second,
		}
		//nolint:errcheck
		go server.ListenAndServe()
	}

	return app, nil
}

// MustNew returns a new App, but panics if there is an error.
func MustNew(ctx context.Context, adminPort string) *App {
	app, err := New(ctx, adminPort)
	if err != nil {
		panic(err)
	}
	return app
}

// Shutdown calls all shutdown methods, in order they were added.
func (a *App) Shutdown(ctx context.Context) (err error) {
	log.Trace().Msg("[app] Starting graceful shutdown.")

	const op = errors.Op("app.App.Shutdown")

	if a.ShutdownTimeout > 0 {
		log.Trace().Dur("shutdown_timeout", a.ShutdownTimeout).Msg("[app] Configuring a timeout for the shutdown.")
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, a.ShutdownTimeout)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		err = errors.E(op, "shutdown deadline has been reached")
	case err = <-a.shutdownAllHandlers(ctx):
	}
	if err != nil {
		log.Trace().Err(err).Msg("[app] Graceful shutdown failed.")
		return errors.E(op, err)
	}

	log.Trace().Msg("[app] Graceful shutdown finished successfully.")
	return nil
}

func (a *App) shutdownAllHandlers(ctx context.Context) chan error {
	const op = errors.Op("shutdownAllHandlers")
	done := make(chan error)
	go func() {
		defer close(done)
		for a.shutdownHandlers.Len() > 0 {
			h := heap.Pop(&a.shutdownHandlers).(*ShutdownHandler)
			if ctx.Err() != nil {
				done <- errors.E(op, "shutdow deadline has been reached")
			}
			trace := log.Trace().
				Str("shutdown_handler_name", h.Name).
				Uint8("shutdown_handler_priority", uint8(h.Priority)).
				Dur("shutdown_handler_timeout", h.Timeout).
				Str("shutdown_handler_policy", ErrorPolicyString(h.Policy))

			trace.Msg("[app] Executing shutdown handler.")
			if err := h.Execute(ctx); err != nil {
				trace.Msg("[app] Shutdown handler failed.")
				done <- errors.E(op, err)
			}
			trace.Msg("[app] Shutdown handler finished.")
		}
	}()
	return done
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown
func (a *App) RunAndWait(mainLoop MainLoopFunc) {
	log.Trace().Msg("[app] Starting run and wait.")

	errs := make(chan error)

	go func() {
		defer Recover()

		a.logger.Info().Msg("Application main loop starting now!")
		if mainLoop == nil {
			errs <- errors.New("main loop is nil")
			return
		}
		a.mainReadinessProbe.SetOk()
		errs <- mainLoop()
	}()

	notifyCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var err error
	ctx := a.logger.WithContext(context.Background())
	select {
	case <-notifyCtx.Done():
		a.mainReadinessProbe.SetNotOk()
		a.logger.Info().
			Dur("grace_period", a.GracePeriod).
			Msg("Graceful shutdown signal received! Awaiting for grace period to end.")
		time.Sleep(a.GracePeriod)
		a.logger.Info().Msg("Grace period is over, initiating shutdown procedures...")
		err = a.Shutdown(ctx)
	case err = <-errs:
		a.mainReadinessProbe.SetNotOk()
		a.logger.Info().Err(err).Msg("Main Loop finished by itself, initiating shutdown procedures...")
		err = a.Shutdown(ctx)
	}
	if err == nil {
		a.logger.Info().Msg("App gracefully terminated.")
	} else {
		a.logger.Error().Err(err).Msg("App terminated with error.")
	}

	// This forces kubernetes kills the pod if some other code is holding the main func.
	a.mainHealthnessProbe.SetNotOk()
}

// RegisterShutdownHandler adds a handler in the end of the list. During shutdown all handlers are executed in the order they were added
func (a *App) RegisterShutdownHandler(sh *ShutdownHandler) {
	if sh.Name == "" {
		panic("Shutdown handler name must not be an empty string")
	}
	if len(a.shutdownHandlers) == 0 {
		heap.Init(&a.shutdownHandlers)
	}
	heap.Push(&a.shutdownHandlers, sh)

	log.Trace().
		Str("shutdown_handler_name", sh.Name).
		Uint8("shutdown_handler_priority", uint8(sh.Priority)).
		Dur("shutdown_handler_timeout", sh.Timeout).
		Str("shutdown_handler_policy", ErrorPolicyString(sh.Policy)).
		Msg("[app] Shutdown handler registered")
}
