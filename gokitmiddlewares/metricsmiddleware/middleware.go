package metricsmiddleware

import (
	"context"

	"github.com/arquivei/foundationkit/gokitmiddlewares"
	"github.com/arquivei/foundationkit/metrifier"
	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
)

// LabelsDecoder defines a functions that generates a label value by processing the
// endpoint's request and response.
type LabelsDecoder func(ctx context.Context, req, resp interface{}, err error) map[string]string

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

	shouldDecodeLabels := c.LabelsDecoder != nil

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			s := m.Begin()
			defer func() {
				defer func() {
					if r := recover(); r != nil {
						log.Ctx(ctx).Error().Msgf("Metrics middleware panicked! Please check your code and configuration: %v", r)
					}
				}()
				if shouldDecodeLabels {
					s = s.WithLabels(c.LabelsDecoder(ctx, req, resp, err))
				}
				s.End(err)
			}()
			return next(ctx, req)
		}
	}, nil
}
