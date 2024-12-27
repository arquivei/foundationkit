// This is a simple example that shows how to setup the metrics middleware.
// This can be run with `go run main.go`
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/arquivei/foundationkit/app"
	"github.com/arquivei/foundationkit/gokitmiddlewares/metricsmiddleware/v4"
	flog "github.com/arquivei/foundationkit/log"

	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version = "v0.0.0-dev"
	cfg     struct {
		// Programs can have any configuration the want.

		HTTP struct {
			Addr string `default:"localhost:8000"`
		}
	}
)

// request is a request to the greeter endpoint
type request struct {
	Name string
}

// response is the response of the greeter endpoint
type response struct {
	Message string
}

// greeter is an endpoint that takes a name and says hello.
func greeter(_ context.Context, req interface{}) (interface{}, error) {
	resp := response{
		Message: "Hello " + req.(request).Name + "!",
	}

	return resp, nil
}

// labelsDecoder is an example that creates labels for the greeter endpoint. To
// avoid overwriting of your labels by the metrics collector, you may want to
// prefix them with a system-specific name
type labelsDecoder struct{}

func (labelsDecoder) Labels() []string {
	return []string{"greeter_empty_name"}
}

func (labelsDecoder) Decode(ctx context.Context, req, resp interface{}, err error) map[string]string {
	if req.(request).Name == "" {
		return map[string]string{"greeter_empty_name": "true"}
	}
	return map[string]string{"greeter_empty_name": "false"}
}

// newExternalMetrics is an example on how to implement external metrics.
func newExternalMetrics(system, subsystem string) func(ctx context.Context, req, resp interface{}, err error) {
	count := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: system,
		Subsystem: subsystem,
		Name:      "letters",
		Help:      "Total amount letters.",
	}, nil)

	return func(ctx context.Context, req, resp interface{}, err error) {
		count.Add((float64(len(req.(request).Name))))
	}
}

func main() {
	app.SetupConfig(&cfg)
	flog.SetupLogger(flog.Config{
		Level: "info",
		Human: true,
	}, version)

	// Just some basic logger initialization
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx := context.Background()
	app.NewDefaultApp(ctx)

	server := newHTTPServer(makeEndpoint())
	app.RunAndWait(func() error {
		log.Info().Msg("Server ready for receiving requests!")
		log.Info().Msg(`Try running 'curl "http://localhost:8000/greet/World"'`)
		log.Info().Msg(`And than running 'curl "http://localhost:9000/metrics"'`)
		return server.ListenAndServe()
	})
}

func makeEndpoint() endpoint.Endpoint {
	// Create some endpoint
	e := greeter

	// Chain metrics middleware
	e = endpoint.Chain(
		metricsmiddleware.MustNew(
			metricsmiddleware.NewDefaultConfig("endpointTest").
				WithLabelsDecoder(labelsDecoder{}).
				WithExternalMetrics(newExternalMetrics("system", "subsystem")),
		),
	)(e)

	return e
}

func newHTTPServer(e endpoint.Endpoint) *http.Server {
	httpServer := &http.Server{Addr: cfg.HTTP.Addr, Handler: makeHandler(e)}

	// You can register the shutdown handlers at any order, but do it before starting the app
	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "http_server",
			Priority: app.ShutdownPriority(100),
			Handler:  httpServer.Shutdown,
			Policy:   app.ErrorPolicyAbort,
		},
	)
	return httpServer
}

func makeHandler(e endpoint.Endpoint) http.Handler {
	handler := kithttp.NewServer(
		e,
		decodeRequest,
		encodeResponse,
	)

	r := http.NewServeMux()

	r.Handle("/greet/{name}", handler)

	return r
}

func decodeRequest(_ context.Context, r *http.Request) (any, error) {
	return request{
		Name: r.PathValue("name"),
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, r any) error {
	return json.NewEncoder(w).Encode(r)
}
