package retrier

import (
	"time"

	"github.com/arquivei/foundationkit/errors"
)

// Retrier is a configurable retry helper based on the use of
// the Strategy pattern.
type Retrier struct {
	RetryEvaluator    RetryEvaluator
	BackoffCalculator BackoffCalculator
	ErrorWrapper      ErrorWrapper
}

// RetryEvaluator defines the logic to decide if an operation should be
// retried or not. Returns bool if the operation should be retried.
type RetryEvaluator interface {
	IsRetryable(attempt int, err error) bool
}

// BackoffCalculator defines how much time to backoff based on the
// current attempt number.
type BackoffCalculator interface {
	CalculateBackoff(attempt int) time.Duration
}

// ErrorWrapper defines how the retry middleware should wrap errors on each successive call.
// Note that implementation of this interface are likely to be NOT-REENTRANT.
type ErrorWrapper interface {
	WrapError(attempt int, err error) error
}

// Settings holds the settings to instantiate a new Retrier. If a field has zero-value,
// a sane default is assumed.
type Settings struct {
	// Evaluator is the ruleset that decides if an error can be retried or not.
	// Defaults to GenericRetryEvaluator with retry only on SeverityRuntime and
	// 5 attempts.
	RetryEvaluator RetryEvaluator
	// BackoffCalculator is how long to wait between each attempt. Defaults
	// to ExponentialBackoffCalculator with 100ms +-20ms and multiplier of 2
	BackoffCalculator BackoffCalculator
	// ErrorWrapper is how to handle the error returned when all attempts have failed.
	// Defaults to LastErrorWrapper
	ErrorWrapper ErrorWrapper
}

// NewRetrier returns a new instance of Retrier, configured
// with the strategies passed by parameter by @settings
func NewRetrier(settings Settings) *Retrier {
	if settings.RetryEvaluator == nil {
		settings.RetryEvaluator = NewGenericRetryEvaluator(GenericRetryEvaluatorSettings{
			MaxAttempts:            5,
			ErrorsSeveritiesPolicy: EvaluationPolicyWhitelist,
			ErrorsSeverities:       []errors.Severity{errors.SeverityRuntime},
		})
	}

	if settings.BackoffCalculator == nil {
		settings.BackoffCalculator = NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
			BaseBackoff:        100 * time.Millisecond,
			RandomExtraBackoff: 20 * time.Millisecond,
			Multiplier:         2.0,
		})
	}

	if settings.ErrorWrapper == nil {
		settings.ErrorWrapper = NewLastErrorWrapper()
	}

	return &Retrier{
		RetryEvaluator:    settings.RetryEvaluator,
		BackoffCalculator: settings.BackoffCalculator,
		ErrorWrapper:      settings.ErrorWrapper,
	}
}

// ExecuteOperation runs the @operation with retry logic. It return the operation
// error according to the set errors wrapper strategy.
func (r Retrier) ExecuteOperation(operation func() error) error {
	for attempt := 1; ; attempt++ {
		err := operation()

		if err == nil {
			return nil
		}

		isRetryable := r.RetryEvaluator.IsRetryable(attempt+1, err)
		if !isRetryable {
			return r.ErrorWrapper.WrapError(attempt, err)
		}

		time.Sleep(r.BackoffCalculator.CalculateBackoff(attempt))
	}
}
