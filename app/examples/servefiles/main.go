package main

// This is an example program on how to use the app package. It uses to FileServer
// handler to serve all files on a given directory.
//
// To run this program in the commandline you could use:
// go run ./app/examples/servefiles/ -log-human -log-level=debug

import (
	"context"
	"net/http"
	"time"

	"github.com/arquivei/foundationkit/app"
	"github.com/arquivei/foundationkit/app/appoptions"
	"github.com/arquivei/foundationkit/log"
)

var version = "development" // Should be replaced on compilation phase

var config struct {
	Log  log.Config
	HTTP struct {
		Port string `default:"8000"`
	}
	Dir string `default:"."`
}

func main() {
	app.SetupConfig(&config)
	ctx := log.SetupLoggerWithContext(context.Background(), config.Log, version)

	// New app. Passing a write timeout configuration
	// as an example of the Option pattern usage. The second
	// paramater onwards can be omitted, and default values will be
	// used.
	app.NewDefaultApp(
		ctx,
		appoptions.WithWriteTimeout(10*time.Second),
	)

	// Some initialization, could take a while.
	// It's a good practice to initialize everything before calling RunAndWait because
	// readiness probe is already up and reporting the app is not ready yet.
	httpServer := &http.Server{Addr: ":" + config.HTTP.Port, Handler: http.FileServer(http.Dir(config.Dir))}

	// You can register the shutdown handlers at any order, but do it before starting the app
	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "http_server",
			Priority: app.ShutdownPriority(100),
			Handler:  httpServer.Shutdown,
			Policy:   app.ErrorPolicyAbort,
		},
	)

	// Run the main loop until it finishes or receives termination signal
	// On this point the readiness probe starts returning success.
	app.RunAndWait(func() error {
		return httpServer.ListenAndServe()
	})
}
