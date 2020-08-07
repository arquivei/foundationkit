package retrier

import (
	"math/rand"
	"time"
)

// ExponentialBackoffCalculator is a simple backoff calculator that multiplies a
// base backoff time n times against a multiplier. The multiplier value
// must be higher than 1.0.
// The multiplier is self-multiplied for each attempt over 2. For example,
// if multiplier is 2, the backoff for each attempt will be multiplied by:
// Attempt 1 - multiplied by 1
// Attempt 2 - multiplied by 2
// Attempt 3 - multiplied by 4
// Attempt 4 - multiplied by 8
type ExponentialBackoffCalculator struct {
	BaseBackoff        time.Duration
	RandomExtraBackoff time.Duration
	Multiplier         float64
}

// ExponentialBackoffCalculatorSettings holds information for how
// to calculate the backoff. Zero values will be turned into sane defaults
type ExponentialBackoffCalculatorSettings struct {
	// BaseBackoff is the base time of how much to wait between attempts. Defaults to
	// 50ms
	BaseBackoff time.Duration
	// RandomExtraBackoff is an amount of extra time between retries that is randomized up to
	// this value. This helps with avoid too many retries happening at once. Defaults to 0 (disabled)
	RandomExtraBackoff time.Duration
	// Multiplier is how much to multiply the backoff time each attempt. It multiplies itself for
	// each attempt above 2. Defaults to 2 if set to a value < 1
	Multiplier float64
}

// NewExponentialBackoffCalculator return a new instance of ExponentialBackoffCalculator configured
// with the given @settings
func NewExponentialBackoffCalculator(settings ExponentialBackoffCalculatorSettings) *ExponentialBackoffCalculator {
	if settings.BaseBackoff == time.Duration(0) {
		settings.BaseBackoff = 50 * time.Millisecond
	}

	if settings.Multiplier < 1 {
		settings.Multiplier = 2
	}

	return &ExponentialBackoffCalculator{
		BaseBackoff:        settings.BaseBackoff,
		RandomExtraBackoff: settings.RandomExtraBackoff,
		Multiplier:         settings.Multiplier,
	}
}

// CalculateBackoff returns a backoff calculator calculated by multiplying a
// base backoff time @attempt-1 times against a multiplier . The multiplier value
// must be higher than 1.0. If @attempt is 1, return the base backoff. See type
// definition for more information
//
// Disables gosec lint here because it complains about math/rand instead of crypto/rand, but here it's ok.
// nolint: gosec
func (c ExponentialBackoffCalculator) CalculateBackoff(attempt int) time.Duration {
	multiplier := float64(1)
	for i := 1; i < attempt; i++ {
		multiplier *= c.Multiplier
	}

	backoff := float64(c.BaseBackoff)
	if c.RandomExtraBackoff > 0 {
		backoff += rand.Float64() * float64(c.RandomExtraBackoff)
	}
	backoff *= multiplier

	return time.Duration(backoff)
}
