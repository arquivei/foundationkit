# Foundation Kit

[![Go Report Card](https://goreportcard.com/badge/github.com/arquivei/foundationkit)](https://goreportcard.com/report/github.com/arquivei/foundationkit)

This project contains very opinionated packages for common operations.

## Trace

### Config and Initialization

In the config variable, add the `trace.Config` struct

```golang
var config struct {
    Trace trace.Config
}
```

The `trace.Config` has the following values:

```golang
type Config struct {
	Exporter          string  `default:""` // empty string or "stackdriver"
	ProbabilitySample float64 `default:"0"` // [0, 1]
	Stackdriver       struct {
		ProjectID string
	}
}
```

Initialize your config using the `app.SetupConfig`

```golang
app.SetupConfig(&config)
```

Now, inicialize your trace exporter using `trace.SetupTrace`

```golang
trace.SetupTrace(config.Trace)
```

### Service Running

It will be used the following service as example:

```golang
type Service interface {
    Do(context.Context, Request) (Response, error)
}

type Response struct {
    // (...)
    trace.Trace
}
```

And will be used the following job as example of job/event:

```golang
type Job interface {
    //(...)
    trace.Trace
}
```

#### HTTP Layer

```golang
func MakeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(Request)

		ctx = trace.WithTraceFromHTTPRequest(ctx, req, getLabelsFromRequest(req))

		response, err := s.Do(ctx, req)

		return response, err
	}
}

func getLabelsFromRequest(req Request) map[string]string{
    //(...)
}
```

#### Method Do

```golang
func (s service) Do(ctx ontext.Context, req Request) (response Response, err error) {

    // When receive the trace through a job/event
    ctx = trace.WithTrace(ctx, job.Trace)

    // When is the first span in service
    span := trace.StartSpanWithParent(ctx)
    defer span.End(err)

    // When is not the first span in service
    span := trace.StartSpan(ctx)
    defer span.End(err)

    // Make a POST in another service
    var req *http.Request
    prepareRequest(req)
    trace.SetTraceInHTTPRequest(ctx, req)

    // Create a job/event to send to a queue
    // or
    // Create the Service Response
    response.Trace = trace.GetTraceFromContext(ctx)
    newJob.Trace = trace.GetTraceFromContext(ctx)
}

func prepareRequest(req *http.Request) {
    //(...)
}
```

#### Logging

Use the method `trace.GetIDFromContext` to log the Trace ID

```golang
func (l *logging) Do(ctx ontext.Context, req Request) (response Response, err error) {
	logger := log.Logger.With().
		Str("trace_id", trace.GetIDFromContext(ctx)).
		Logger()
    // (...)
}
```

#### Encoding

```golang
func EncodeResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
    response := r.(Response)
    trace.SetTraceInHTTPResponse(response.Trace, w)
    // (...)
}
```