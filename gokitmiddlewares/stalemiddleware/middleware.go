package stalemiddleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// New returns a new stale middleware
// This middleware marks the application as
// unhealthy if no request is received in
// the allotted time.
//
// If a request is received after being unhealthy,
// it's not considered stale anymore abd becomes
// healthy again.
func New(c Config) endpoint.Middleware {
	var lastKnownRequestTime int64

	go backgroundStaleCheck(c, &lastKnownRequestTime)

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			lastKnownRequestTime = time.Now().UnixNano()
			return next(ctx, request)
		}
	}
}

func backgroundStaleCheck(c Config, lastKnownRequestTime *int64) {
	time.Sleep(c.StartCheckAfter)
	*lastKnownRequestTime = time.Now().UnixNano()
	for {
		time.Sleep(c.MaxTimeBetweenRequests)

		if isUnhealthy(c, *lastKnownRequestTime) {
			if c.HealthinessPobe.IsOk() {
				logAndSetUnhealthy(c)
			}
			continue
		} else if !c.HealthinessPobe.IsOk() {
			logAndSetHealthy(c)
		}
	}
}

func logAndSetUnhealthy(c Config) {
	if c.Logger != nil {
		c.Logger.Warn().
			Dur("timeout", c.MaxTimeBetweenRequests).
			Msg("Endpoint didn't receive any request and it's stale")
	}
	c.HealthinessPobe.SetNotOk()
}

func logAndSetHealthy(c Config) {
	if c.Logger != nil {
		c.Logger.Info().
			Dur("timeout", c.MaxTimeBetweenRequests).
			Msg("Endpoint received a request and is no longer stale")
	}
	c.HealthinessPobe.SetOk()
}

func isUnhealthy(c Config, last int64) bool {
	return time.Since(time.Unix(0, last)) > c.MaxTimeBetweenRequests
}
