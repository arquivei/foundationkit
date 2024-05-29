// This is a simple example that shows how to setup the metrics middleware.
// This can be run with `go run main.go`
package main

import (
	"context"
	"os"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/gokitmiddlewares/metricsmiddleware/v4"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Request is a Request to the greeter endpoint
type Request struct {
	Name string
}

// Response is the Response of the greeter endpoint
type Response struct {
	Message string
}

// greeter is an endpoint that takes a name and says hello.
func greeter(_ context.Context, req Request) (Response, error) {
	resp := Response{
		Message: "Hello " + req.Name + "!",
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
	if req.(Request).Name == "" {
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
		count.Add((float64(len(req.(Request).Name))))
	}
}

func main() {
	// Just some basic logger initialization
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Create some endpoint
	e := greeter

	// Chain metrics middleware
	e = endpoint.Chain(
		metricsmiddleware.MustNew[Request, Response](
			metricsmiddleware.NewDefaultConfig("endpointTest").
				WithLabelsDecoder(labelsDecoder{}).
				WithExternalMetrics(newExternalMetrics("system", "subsystem")),
		),
	)(e)

	// Let's just run the example for fun.
	ctx := context.Background()
	req := Request{Name: "World"}
	resp, err := e(ctx, req)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	log.Info().Msg(resp.Message)
}
