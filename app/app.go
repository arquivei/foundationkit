package app

import (
	"container/heap"
	"context"
	"net/http"
	"net/http/pprof" // Sadly, this also changes the DefaultMux to have the pprof URLs
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// MainLoopFunc is the functions runned by app. If it finishes, it will trigger a shutdown
type MainLoopFunc func() error

// App represents an application with a main loop and a shutdown routine
type App struct {
	ctx context.Context

	Ready   ProbeGroup
	Healthy ProbeGroup

	shutdownHandlers shutdownHeap
	GracePeriod      time.Duration
	ShutdownTimeout  time.Duration

	mainReadinessProbe  Probe
	mainHealthnessProbe Probe
}

// New returns a new App.
func New(ctx context.Context, adminPort string) (*App, error) {
	app := &App{
		ctx:     ctx,
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
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				//nolint:errcheck
				w.Write([]byte(cause))
			}
		})

		mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
			isReady, cause := app.Ready.CheckProbes()
			if isReady {
				w.WriteHeader(http.StatusOK)
				//nolint:errcheck
				w.Write([]byte("OK"))
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				//nolint:errcheck
				w.Write([]byte(cause))
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
			Addr:    ":" + adminPort,
			Handler: mux,
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
func (a *App) Shutdown(ctx context.Context) error {
	const op = errors.Op("app.App.Shutdown")
	if a.ShutdownTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, a.ShutdownTimeout)
		defer cancel()
	}
	var err error
	select {
	case <-ctx.Done():
		err = errors.E(op, "shutdown deadline has been reached")
	case err = <-a.shutdownAllHandlers(ctx):
	}
	if err != nil {
		return errors.E(op, err)
	}
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
			if err := h.Execute(ctx); err != nil {
				done <- errors.E(op, err)
			}
		}
	}()
	return done
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown
func (a *App) RunAndWait(mainLoop MainLoopFunc) {
	errs := make(chan error)

	go func() {
		log.Ctx(a.ctx).Info().Msg("Application main loop starting now!")
		if mainLoop == nil {
			errs <- errors.New("main loop is nil")
			return
		}
		a.mainReadinessProbe.SetOk()
		errs <- mainLoop()
	}()

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	var err error
	ctx := log.Ctx(a.ctx).WithContext(context.Background())
	select {
	case <-a.ctx.Done():
		a.mainReadinessProbe.SetNotOk()
		log.Ctx(a.ctx).Info().Msg("App context canceled, initialing shutdown procedures...")
		err = a.Shutdown(ctx)
	case s := <-signals:
		a.mainReadinessProbe.SetNotOk()
		log.Ctx(a.ctx).Info().
			Str("signal", s.String()).
			Dur("grace_period", a.GracePeriod).
			Msg("Signal received. Waiting grace period...")
		time.Sleep(a.GracePeriod)
		log.Ctx(a.ctx).Info().Msg("Grace period is over, initiating shutdown procedures...")
		err = a.Shutdown(ctx)
	case err = <-errs:
		a.mainReadinessProbe.SetNotOk()
		log.Ctx(a.ctx).Info().Err(err).Msg("App finished by itself, initialing shutdown procedures...")
		err = a.Shutdown(ctx)
	}
	if err == nil {
		log.Ctx(a.ctx).Info().Msg("App exited")
	} else {
		log.Ctx(a.ctx).Error().Err(err).Msg("App exited with error")
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
}
