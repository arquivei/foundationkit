package main

// This example shows how panics are captured and gracefully printed
// by the app.Recover() function.
//
// To run this program in the commandline you could use:
// go run ./app/examples/panic/ -log-human -log-level=trace

import (
	"context"

	"github.com/arquivei/foundationkit/app"
	"github.com/arquivei/foundationkit/log"
)

var version = "development"

var config struct {
	Log log.Config
}

func main() {
	defer app.Recover()

	app.SetupConfig(&config)
	ctx := log.SetupLoggerWithContext(context.Background(), config.Log, version)

	// New app
	if err := app.NewDefaultApp(ctx); err != nil {
		panic(err)
	}

	// Comment this next line to see the other panic
	// thisWillPanic()

	app.RunAndWait(func() error {
		panic("panics inside run and wait will trigger a graceful shutdown")
	})
}

// nolint: unused
func thisWillPanic() {
	panic("panics outside RunAndWait should be caught by app.Recover()")
}
