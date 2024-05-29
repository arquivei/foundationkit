package apiping

import (
	"context"

	"github.com/arquivei/foundationkit/endpoint"
	"github.com/arquivei/foundationkit/trace/v2/examples/services/ping"
)

// MakeAPIPingEndpoint returns an edpoint
func MakeAPIPingEndpoint(service ping.Service) endpoint.Endpoint[Request, string] {
	return func(ctx context.Context, req Request) (string, error) {
		return service.Ping(ctx, ping.Request{
			Num:   req.Num,
			Sleep: req.Sleep,
		})
	}
}
