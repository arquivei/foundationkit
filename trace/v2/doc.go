/*
# Trace v2

This guide will help you to setup and use the trace/v2 lib. We have a usage example in `examples` folder.

## Configuring

This lib provides a `Config` struct and this section will explain how to properly configure.

All configuration is done via environment variables for the OpenTelemetry SDK. Please refer to
https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/ for more information.

## Setting up and using the lib

# Call trace.Setup to initialize the OpenTemetry SDK

```golang

traceShutdown := trace.Setup(ctx)
app.RegisterShutdownHandler(

		&app.ShutdownHandler{
			Name:     "opentelemetry_trace",
			Priority: app.ShutdownPriority(shutdownPriorityTrace),
			Handler:  traceShutdown,
			Policy:   app.ErrorPolicyAbort,
	})

```

# Setting the trace up AFTER create a new App in main.go

```golang

	func main() {
		app.SetupConfig(&config)

	    // (...)

		if err := app.NewDefaultApp(ctx); err != nil {
			log.Ctx(ctx).Fatal().Err(err).Msg("Failed to create app")
		}

		setupTrace()

	    // (...)
	}

```

# Starting a span

```golang
ctx, span := trace.Start(ctx, "SPAN-NAME")
defer span.End()
```

# Recovering trace information and logging it

```golang

	type TraceInfo struct {
		ID        string
		IsSampled bool
	}

```

```golang
t := trace.GetTraceInfoFromContext(ctx)
log.Ctx(ctx).Info().EmbedObject(t).Msg("Hello")
```

Refer to `examples/` directory for a working example.

## Propagating the trace

### API

# Using in transport layer

Use the middleware to automatically handle traces from HTTP Requests.

```golang

	func MakeHTTPHandler(e endpoint.Endpoint) http.Handler {
	    // (...)

		r := mux.NewRouter()
		r.Use(trace.MuxHTTPMiddleware("SERVER-NAME"))

	    //(...)
		return r
	}

```

The "SERVER-NAME" is optional and OpenTelemetry will default to the host's IP.

# Using with go-kit endpoints

```golang

	e := endpoint.Chain(
		trace.EndpointMiddleware("my-endpoint"),
		loggingmiddleware.MustNew(loggingConfig),
	)(srv))

```

# Using in a HTTP request

```golang
request, err := http.NewRequestWithContext(

	ctx,
	"POST",
	url,
	bytes.NewReader(body),

)

	if err != nil {
	    return "", errors.E(op, err)
	}

trace.SetTraceInRequest(request)
```

### Workers

For exporting the trace you can use `trace.ToMap` and `trace.FromMap` to export the necessary information to a `map[string]string` that could be marshaled into messages or used as message metadata if you broker supports it.
*/
package trace // import "github.com/arquivei/foundationkit/trace/v2"
