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
	c, err := checkAndFixConfig(c)
	if err != nil {
		return nil, err
	}

	m, err := metrifier.New(c.Metrifier)
	if err != nil {
		return nil, err
	}

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			s := m.Begin()
			defer func() {
				s.WithLabels(makeLabels(ctx, c, req, resp, err)).End(err)
			}()
			return next(ctx, req)
		}
	}, nil
}

// makeLabels returns an map of labels for the request.
func makeLabels(ctx context.Context, c Config, req, resp interface{}, err error) map[string]string {
	if len(c.LabelDecoders) == 0 {
		return nil
	}

	var labels map[string]string
	if len(c.LabelDecoders) > 0 {
		labels = make(map[string]string)
		for label, decoder := range c.LabelDecoders {
			labels[label] = decoder(ctx, req, resp, err)
		}
	}
	return labels
}

// checkAndFixConfig checks if ExtraLabels and LabelsDecoder are consistent.
// If ExtraLabels is empty, it fixes the config by coping the LabelDecoders keys.
func checkAndFixConfig(c Config) (Config, error) {
	if len(c.Metrifier.ExtraLabels) == 0 && len(c.LabelDecoders) > 0 {
		return fixExtraLabels(c), nil
	}

	if len(c.Metrifier.ExtraLabels) != len(c.LabelDecoders) {
		return Config{}, errors.New("wrong config: Metrifier.ExtraLabels should have the same keys as LabelDecoders")
	}

	for _, l := range c.Metrifier.ExtraLabels {
		if _, ok := c.LabelDecoders[l]; !ok {
			return Config{}, errors.New("label declared but missing decoder: " + l)
		}
	}

	return c, nil
}

// fixExtraLabels copies the LabelDecoder keys to the ExtraLables slice.
func fixExtraLabels(c Config) Config {
	c.Metrifier.ExtraLabels = make([]string, 0, len(c.LabelDecoders))
	for k := range c.LabelDecoders {
		c.Metrifier.ExtraLabels = append(c.Metrifier.ExtraLabels, k)
	}
	return c
}
