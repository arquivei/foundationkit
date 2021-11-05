package backoffmiddleware

import (
	"context"
	"math/rand"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/go-kit/kit/endpoint"
)

// Config contains the config for the exponential backoff retrier
type Config struct {
	// InitialDelay represents the delay after the first error, before adding
	// the spread
	InitialDelay time.Duration

	// MaxDelay represents the max delay after an error, before adding the
	// spread
	MaxDelay time.Duration

	// Spread is the percentage of the current delay that can be added as a
	// random term. For example, with a delay of 10s and 20% spread, the
	// calculated delay will be between 10s and 12s.
	Spread float64

	// Factor represents how bigger the next delay wil be in comparison to the
	// current one
	Factor float64

	// MaxRetries indicates how many times this middleware should retry.
	// SeverityRuntime errors are always retried and don't count . SeverityInput errors are never retried.
	MaxRetries int
}

// MaxRetriesInfinite constant indicate that the middleware should never give up retrying
const MaxRetriesInfinite = -1

// NewDefaultConfig returns a ready-to-use Config with sane defaults
func NewDefaultConfig() Config {
	return Config{
		InitialDelay: 200 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Spread:       0.2,
		Factor:       1.5,
		MaxRetries:   MaxRetriesInfinite,
	}
}

// New tries to execute @next.Process() until it succeeds. Each failure is
// followed by an exponentially increasing delay.
func New(config Config) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			return runWithBackoff(ctx, config, next, request)
		}
	}
}

func runWithBackoff(
	ctx context.Context,
	config Config,
	next endpoint.Endpoint,
	request interface{},
) (interface{}, error) {
	delay := config.InitialDelay
	retries := 0

	response, err := next(ctx, request)
	for err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		switch errors.GetSeverity(err) {
		case errors.SeverityInput:
			return nil, err
		case errors.SeverityFatal:
			return nil, err
		case errors.SeverityRuntime:
			// always retry
		default:
			retries++
			if config.MaxRetries != MaxRetriesInfinite &&
				retries > config.MaxRetries {
				return nil, err
			}
		}

		amountToSleep := addSpread(delay, config.Spread)

		waitCtx, cancelFn := context.WithTimeout(context.Background(), amountToSleep)
		defer cancelFn()

		select {
		case <-waitCtx.Done():
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		delay = time.Duration(float64(delay) * config.Factor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		response, err = next(ctx, request)
	}
	return response, nil
}

func addSpread(delay time.Duration, spread float64) time.Duration {
	spreadRange := int64(float64(delay.Nanoseconds()) * spread)
	return delay + time.Duration(rand.Int63n(spreadRange))*time.Nanosecond //nolint:gosec
}
