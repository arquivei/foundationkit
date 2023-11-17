# Examples

## API

This example will create an API (ping) that configures the trace, creates a span and propagates to another API (pong). The ping API is configured to use itself as pong API.

The whole example runs in docker with a [jaeger](https://www.jaegertracing.io/) collecting the traces.

To run the example, run `make run` from this directory to create the API.

Wait for the API to start serving requests. It will print the message `Application main loop starting now!` on the logs. 

After that, send a request using `make send`. You can change the parameters of the request on the `docker-compose.yaml` file.

Them visit http://localhost:16686/ to see the traces.

## Example V1 compatibility

In the file `trace/v2/examples/cmd/api/resources.go` we apply the new middlewares. But if you want to test with the old `trackingmiddleware.New`, comment the trace and request middlewares and uncomment the trackingmiddleware. 

``` go
	r.Use(
		// This is deprecated. Used when we can't ditch trace v1.
		// trackingmiddleware.New,
		// This is the preferred way
		trace.MuxHTTPMiddleware(""),
		request.HTTPMiddleware,
		enrichloggingmiddleware.New,
	)
```

