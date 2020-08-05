package metricsmiddleware

import (
	"context"
	"errors"

	"github.com/arquivei/foundationkit/gokitmiddlewares"
	"github.com/arquivei/foundationkit/metrifier"
	"github.com/go-kit/kit/endpoint"
)

// LabelDecoder defines a functions that generates a label value by processing the
// endpoint's request and response.
type LabelDecoder func(ctx context.Context, req, resp interface{}, err error) string

// MustNew returns a new metrics middleware but panics in case of error.
func MustNew(c Config) endpoint.Middleware {
	return gokitmiddlewares.Must(New(c))
}

// New returns a new metrics middleware.
func New(c Config) (endpoint.Middleware, error) {
	m, err := metrifier.New(c.Metrifier)
	if err != nil {
		return nil, err
	}

	if len(c.Metrifier.ExtraLabels) == 0 && len(c.LabelDecoders) > 0 {

	}

	if len(c.Metrifier.ExtraLabels) != len(c.LabelDecoders) {
		return nil, errors.New("wrong config: Metrifier.ExtraLabels should have the same keys as LabelDecoders")
	}
	for _, l := range c.Metrifier.ExtraLabels {
		if _, ok := c.LabelDecoders[l]; !ok {
			return nil, errors.New("label declared but missing decoder: " + l)
		}
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			s := m.Begin()
			defer func() {
				var labels map[string]string
				if len(c.LabelDecoders) > 0 {
					labels := make(map[string]string)
					for label, decoder := range c.LabelDecoders {
						labels[label] = decoder(ctx, req, resp, err)
					}
				}
				s.WithLabels(labels).End(err)
			}()
			return next(ctx, req)
		}
	}, nil
}
