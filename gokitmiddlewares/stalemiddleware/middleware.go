package stalemiddleware

import (
	"context"
	"time"

	"github.com/arquivei/foundationkit/app"
	"github.com/go-kit/kit/endpoint"
)

func New(c Config) endpoint.Middleware {
	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				if !app.IsHealthy() {
					app.SetHealthy()
				}
			case <-time.After(c.Timeout):
				if app.IsHealthy() {
					if c.Logger != nil {
						c.Logger.Warn().
							Dur("timeout", c.Timeout).
							Msg("Endpoint didn't receive any request and it's stale")
					}
					app.SetUnhealthy()
				}
			}
		}
	}()

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			ch <- struct{}{}
			return next(ctx, request)
		}
	}
}
