# Foundation Kit

[![PkgGoDev](https://pkg.go.dev/badge/github.com/arquivei/foundationkit)](https://pkg.go.dev/github.com/arquivei/foundationkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/arquivei/foundationkit)](https://goreportcard.com/report/github.com/arquivei/foundationkit)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

This project contains very opinionated packages for common operations.

## Request

### Usage

It will be used the following service as example:

```golang
type Service interface {
    Do(context.Context, Request) (Response, error)
}

type Response struct {
    //(...)
    RequestID request.ID
}
```

#### HTTP Layer

Use the method `request.WithID` to create and put a Request ID in context

```golang
func MakeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(Request)

		ctx = request.WithID(ctx)

		response, err := s.Do(ctx, req)

		return response, err
	}
}
```

#### Logging

Use the method `request.GetIDFromContext` to log the Request ID

```golang
func (l *logging) Do(ctx ontext.Context, req Request) (response Response, err error) {
	logger := log.Logger.With().
		EmbedObject(request.GetIDFromContext(ctx)).
		Logger()
    // (...)
}
```

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

### Service

It will be used the following service as example:

```golang
type Service interface {
    Do(context.Context, Request) (Response, error)
}

type Request struct {
    // (...)
    trace.Trace
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

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
    var req Request

    // (...)
    req.Trace = trace.GetTraceFromHTTRequest(r)

    return req, nil
}


func MakeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(Request)

		ctx = trace.WithTraceAndLabels(ctx, req.Trace, getLabelsFromRequest(req))

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
    trace.SetInHTTPRequest(ctx, req)

    // Create a job/event to send to a queue
    // or
    // Create the Service Response
    response.Trace = trace.GetFromContext(ctx)
    newJob.Trace = trace.GetFromContext(ctx)
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
		EmbedObject(trace.GetIDFromContext(ctx)).
		Logger()
    // (...)
}
```

#### Encoding

```golang
func EncodeResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
    response := r.(Response)
    trace.SetInHTTPResponse(response.Trace, w)
    // (...)
}
```


## SplitIO

### Using the client directly
```golang
client := splitio.MustNewClient(config)
attributes := splitIO.Attributes{
    "age": 37,
}
if client.IsFeatureEnabled("MY_FEATURE_FLAG", attributes) {
    // do stuff
}
```

### Using the middleware
```golang
// in the service
MY_FF := splitio.Feature("MY_FEATURE_FLAG")
Features := []splitio.Feature{
    MY_FF,
}

// In the transport layer
func GetUserFromRequest(ctx context.Context, request interface{}) map[User]Attributes {
    // Extract user and attributes from the endpoint request
}

// in the main package
client := splitio.MustNewClient(config)
middlewareConfig := splitio.DefaultFFMidlewareConfig()
middlewareConfig.MultiUserDecodeFn = myapi.GetUserFromRequest
middlewareConfig.Features = myservice.Features

middleware := NewFeatureFlagMiddleware(client, middlewareConfig)
myEndpoint := endpoint.Chain(
    // ...
    middleware,
    // ...
)(myEndpoint)

// In your code
if splitio.IsFeatureEnabled(ctx, MY_FF) {
    // do stuff
}
```

Although the initial setup is more complex, it has the advantage of setting up everything only once, and then integrating seamlessly with new feature flags.

Each feature is checked once per request, for the users and attributes specified in the `MultiUserDecodeFn`. The behaviour is stored in the context, so that it is possible to check anywhere in the service if the feature is enabled without caring about which user or attributes should be used.