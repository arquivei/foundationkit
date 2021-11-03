package apiping

import (
	"context"

	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping"

	"github.com/go-kit/kit/endpoint"
)

// MakeAPIPingEndpoint returns an edpoint
func MakeAPIPingEndpoint(service ping.Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(Request)
		return service.Ping(ctx, ping.Request{
			Num:   req.Num,
			Sleep: req.Sleep,
		})
	}
}
