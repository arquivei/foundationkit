package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// ErrorPolicy specifies what should be done when a handler fails
type ErrorPolicy int

const (
	// ErrorPolicyWarn prints the error as a warning and continues to the next handler. This is the default.
	ErrorPolicyWarn ErrorPolicy = iota
	// ErrorPolicyAbort stops the shutdown process and returns an error
	ErrorPolicyAbort
	// ErrorPolicyFatal logs the error as Fatal, it means the application will close immediately
	ErrorPolicyFatal
	// ErrorPolicyPanic panics if there is an error
	ErrorPolicyPanic
)

// MainLoopFunc is the functions runned by app. If it finishes, it will trigger a shutdown
type MainLoopFunc func() error

// ShutdownFunc is a shutdown function that will be executed when the app is shutting down.
type ShutdownFunc func(context.Context) error

type shutdownHandler struct {
	Name     string
	Shutdown ShutdownFunc
	Timeout  time.Duration
	Policy   ErrorPolicy
}

// App represents an application with a main loop and a shutdown routine
type App struct {
	ctx context.Context

	Ready   bool
	Healthy bool

	mainLoop         MainLoopFunc
	shutdownHandlers []shutdownHandler
	GracePeriod      time.Duration
	ShutdownTimeout  time.Duration
}

// New returns a new App.
func New(ctx context.Context, adminPort string, mainLoop MainLoopFunc) (*App, error) {
	if mainLoop == nil {
		return nil, errors.New("main loop is nil")
	}

	app := &App{
		ctx:      ctx,
		mainLoop: mainLoop,
		Ready:    false,
		Healthy:  true,
	}

	{ // This spwans an admin HTTP server for this
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthy", func(w http.ResponseWriter, _ *http.Request) {
			if app.Healthy {
				w.WriteHeader(http.StatusOK)
				//nolint:errcheck
				w.Write([]byte("OK"))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				//nolint:errcheck
				w.Write([]byte("NOT OK"))
			}
		})

		mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
			if app.Ready {
				w.WriteHeader(http.StatusOK)
				//nolint:errcheck
				w.Write([]byte("OK"))
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				//nolint:errcheck
				w.Write([]byte("NOT OK"))
			}
		})

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
func MustNew(ctx context.Context, adminPort string, mainLoop MainLoopFunc) *App {
	app, err := New(ctx, adminPort, mainLoop)
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
		err = ctx.Err()
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
		for _, h := range a.shutdownHandlers {
			if ctx.Err() != nil {
				done <- errors.E(op, ctx.Err())
			}
			err := a.executeHandlerShutdown(ctx, h)
			if err != nil {
				done <- errors.E(op, err)
			}
		}
	}()
	return done
}

func (a *App) executeHandlerShutdown(ctx context.Context, h shutdownHandler) error {
	const op = errors.Op("executeHandlerShutdown")

	if h.Timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	err := h.Shutdown(ctx)
	if err != nil {
		err = errors.E(op, errors.E(errors.Op(h.Name), err))
		switch h.Policy {
		case ErrorPolicyWarn:
			log.Ctx(a.ctx).Warn().Err(err).Msg("Shutdown handler failed")
		case ErrorPolicyAbort:
			// No need for logging here, this will happen latter
			return err
		case ErrorPolicyFatal:
			log.Ctx(a.ctx).Fatal().Err(err).Msg("Shutdown handler failed")
		case ErrorPolicyPanic:
			panic(err)
		default:
			panic(errors.Errorf("invalid error policy: %v", h.Policy))
		}
	}
	log.Ctx(a.ctx).Info().Str("handler", h.Name).Msg("Shutdown successfull")
	return nil
}

// RunAndWait executes the main loop on a go-routine and listens to SIGINT and SIGKILL to start the shutdown
func (a *App) RunAndWait() {
	errs := make(chan error)

	go func() {
		log.Ctx(a.ctx).Info().Msg("Application main loop starting now!")
		a.Ready = true
		errs <- a.mainLoop()
	}()

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	var err error
	ctx := log.Ctx(a.ctx).WithContext(context.Background())
	select {
	case <-a.ctx.Done():
		a.Ready = false
		log.Ctx(a.ctx).Info().Msg("App context canceled, initialing shutdown procedures...")
		err = a.Shutdown(ctx)
	case s := <-signals:
		a.Ready = false
		log.Ctx(a.ctx).Info().
			Str("signal", s.String()).
			Dur("grace_period", a.GracePeriod).
			Msg("Signal received. Waiting grace period...")
		time.Sleep(a.GracePeriod)
		log.Ctx(a.ctx).Info().Msg("Grace period is over, initiating shutdown procedures...")
		err = a.Shutdown(ctx)
	case err = <-errs: // App finished by itself
		a.Ready = false
	}
	if err == nil {
		log.Ctx(a.ctx).Info().Msg("App exited")
	} else {
		log.Ctx(a.ctx).Error().Err(err).Msg("App exited with error")
	}
	// This forces kubernetes kills the pod if some other code is holding the main func.
	a.Healthy = false
}

// RegisterShutdownHandler adds a handler in the end of the list. During shutdown all handlers are executed in the order they were added
func (a *App) RegisterShutdownHandler(name string, fn ShutdownFunc, options ...interface{}) {
	h := shutdownHandler{
		Name:     name,
		Shutdown: fn,
	}
	for _, o := range options {
		switch opt := o.(type) {
		case time.Duration:
			h.Timeout = opt
		case ErrorPolicy:
			h.Policy = opt
		default:
			panic(errors.Errorf("invalid shutdown handler option: [%T]%v", o, o))
		}
	}
	a.shutdownHandlers = append(a.shutdownHandlers, h)
}
