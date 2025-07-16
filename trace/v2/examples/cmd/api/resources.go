package main

import (
	"context"
	"net/http"
	"time"

	"github.com/arquivei/foundationkit/app"
	"github.com/arquivei/foundationkit/gokitmiddlewares/loggingmiddleware"
	"github.com/arquivei/foundationkit/httpmiddlewares/enrichloggingmiddleware"
	"github.com/arquivei/foundationkit/request"
	tracev1 "github.com/arquivei/foundationkit/trace"
	"github.com/arquivei/foundationkit/trace/v2"
	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping"
	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping/apiping"
	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping/implping"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

func setupTrace() {
	// This is only to show that v1 and v2 can coexist
	// This is not recommended in production.
	tracev1.SetupTrace(config.TraceV1)

	traceShutdown := trace.Setup(context.Background())
	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "opentelemetry_trace",
			Priority: shutdownPriorityTrace,
			Handler:  traceShutdown,
			Policy:   app.ErrorPolicyAbort,
		})
}

func getEndpoint() endpoint.Endpoint {
	loggingConfig := loggingmiddleware.NewDefaultConfig("ping")

	pongAdapter := implping.NewHTTPPongAdapter(
		&http.Client{Timeout: config.Pong.HTTP.Timeout},
		config.Pong.HTTP.URL,
	)

	pingEndpoint := endpoint.Chain(
		trace.EndpointMiddleware("ping-pong-endpoint"),
		loggingmiddleware.MustNew(loggingConfig),
	)(apiping.MakeAPIPingEndpoint(
		ping.NewService(pongAdapter),
	))

	return pingEndpoint
}

func getHTTPServer() *http.Server {
	r := mux.NewRouter()

	r.Use(
		// This is deprecated. Used when we can't ditch trace v1.
		// trackingmiddleware.New,
		// This is the preferred way
		trace.MuxHTTPMiddleware(""),
		request.HTTPMiddleware,
		enrichloggingmiddleware.New,
	)

	r.PathPrefix("/ping/").Handler(apiping.MakeHTTPHandler(getEndpoint()))

	httpAddr := ":" + config.HTTP.Port
	httpServer := &http.Server{
		Addr:              httpAddr,
		ReadHeaderTimeout: time.Second,
		Handler:           r}

	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "http_server",
			Priority: shutdownPriorityHTTP,
			Handler:  httpServer.Shutdown,
			Policy:   app.ErrorPolicyAbort,
		})

	return httpServer
}
