/*
Package trace provides distributed tracing.

Initialization

In the config variable, add the "Config" struct:

	var config struct {
		Trace trace.Config
	}

The "Config" has the following values:

	type Config struct {
		Exporter          string  `default:""` // empty string or "stackdriver"
		ProbabilitySample float64 `default:"0"` // [0, 1]
		Stackdriver       struct {
			ProjectID string
		}
	}

Initialize your config using the "app.SetupConfig"

	app.SetupConfig(&config)

Now, initialize your trace exporter using "trace.SetupTrace":

	trace.SetupTrace(config.Trace)

Service

It will be used the following service as example:

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

And will be used the following job as example of job/event:

	type Job interface {
		//(...)
		trace.Trace
	}

HTTP Layer

Retrieve the Trace from the HTTP request and pass it along:

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

Encode the trace in the HTTP respose:

	func EncodeResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
		response := r.(Response)
		trace.SetInHTTPResponse(response.Trace, w)
		// (...)
	}
Using spans

Inside a service function:

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

Logging

Use the function "GetIDFromContext" to log the Trace ID:

	func (l *logging) Do(ctx ontext.Context, req Request) (response Response, err error) {
		logger := log.Logger.With().
			EmbedObject(trace.GetIDFromContext(ctx)).
			Logger()
		// (...)
	}


*/
package trace
