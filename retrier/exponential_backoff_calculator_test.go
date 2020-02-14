package retrier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffCalculator_DefaultValues(t *testing.T) {
	calculator := NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{})
	assert.Equal(t, 50*time.Millisecond, calculator.BaseBackoff, "base backoff")
	assert.Equal(t, time.Duration(0), calculator.RandomExtraBackoff, "random extra backoff")
	assert.Equal(t, 2.0, calculator.Multiplier, "multiplier")
}

func TestExponentialBackoffCalculator(t *testing.T) {
	// Test backoff
	squareBackoffCalculator := NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
		BaseBackoff: 3 * time.Second,
		Multiplier:  2,
	})
	assert.Equal(t, 3*time.Second, squareBackoffCalculator.CalculateBackoff(1), "square backoff attempt 1")
	assert.Equal(t, 6*time.Second, squareBackoffCalculator.CalculateBackoff(2), "square backoff attempt 2")
	assert.Equal(t, 12*time.Second, squareBackoffCalculator.CalculateBackoff(3), "square backoff attempt 3")
	assert.Equal(t, 24*time.Second, squareBackoffCalculator.CalculateBackoff(4), "square backoff attempt 4")

	unusualBackoffCalculator := NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
		BaseBackoff: 10 * time.Second,
		Multiplier:  1.5,
	})
	assert.Equal(t, 10000*time.Millisecond, unusualBackoffCalculator.CalculateBackoff(1), "unusual backoff attempt 1")
	assert.Equal(t, 15000*time.Millisecond, unusualBackoffCalculator.CalculateBackoff(2), "unusual backoff attempt 2")
	assert.Equal(t, 22500*time.Millisecond, unusualBackoffCalculator.CalculateBackoff(3), "unusual backoff attempt 3")
	assert.Equal(t, 33750*time.Millisecond, unusualBackoffCalculator.CalculateBackoff(4), "unusual backoff attempt 4")

	// Test that random backoff is indeed random
	withRandomBackoffCalculator := NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
		BaseBackoff:        1 * time.Hour,
		RandomExtraBackoff: 30 * time.Second,
	})
	firstRandomAttempt := withRandomBackoffCalculator.CalculateBackoff(2)
	foundDifferent := false
	for i := 0; i < 100; i++ {
		if firstRandomAttempt != withRandomBackoffCalculator.CalculateBackoff(2) {
			foundDifferent = true
			break
		}
	}
	assert.True(t, foundDifferent, "random backoff generating random values")
}
