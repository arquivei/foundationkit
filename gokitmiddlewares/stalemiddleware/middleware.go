package stalemiddleware

import (
	"context"
	"time"

	"github.com/arquivei/foundationkit/app"
	"github.com/go-kit/kit/endpoint"
)

// New returns a new stale middleware
// This middleware marks the application as
// unhealthy if no request is received in
// the allotted time.
//
// Unhealthy is a final state for now, so the checker
// stops after reaching unhealthy.
//
// With a better health handler we could try
// to recover from it.
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
		if time.Since(time.Unix(0, *lastKnownRequestTime)) > c.MaxTimeBetweenRequests {
			logAndSetUnhealthy(c)
			return
		}
	}
}

func logAndSetUnhealthy(c Config) {
	if c.Logger != nil {
		c.Logger.Warn().
			Dur("timeout", c.MaxTimeBetweenRequests).
			Msg("Endpoint didn't receive any request and it's stale")
	}
	app.SetUnhealthy()
}
