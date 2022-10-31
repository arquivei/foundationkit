// This example shows how panics are captured and gracefuly printed.
//
// To run this example
// go run -ldflags="-X main.version=v0.0.1" ./app/v2/examples/panic/ -app-log-human -app-log-level=trace
package main

import (
	"context"

	"github.com/arquivei/foundationkit/app/v2"
)

var (
	version = "development"
	config  struct {
		app.Config
	}
)

func main() {
	defer app.Recover()

	app.Bootstrap("", &config)

	// Comment this next line to see the other panic
	thisWillPanic()

	app.RunAndWait(func(_ context.Context) error {
		panic("panics inside run and wait will trigger a graceful shutdown")
	})
}

func thisWillPanic() {
	panic("panics outside RunAndWait should be caught by  app.Recover()")
}
