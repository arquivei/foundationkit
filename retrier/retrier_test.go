package retrier

import (
	"testing"
	"time"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewRetrier_DefaultValues(t *testing.T) {
	retrier := NewRetrier(Settings{})
	assert.NotNil(t, retrier.RetryEvaluator, "evaluator should not be nil")
	assert.NotNil(t, retrier.BackoffCalculator, "backoff calculator should not be nil")
	assert.NotNil(t, retrier.ErrorWrapper, "error wrapper should not be nil")

	evaluator, ok := retrier.RetryEvaluator.(*GenericRetryEvaluator)
	assert.True(t, ok, "evaluator cast")
	assert.Equal(t, 5, evaluator.MaxAttempts, "evaluator - max retry attempts")
	assert.Equal(t, EvaluationPolicyBlacklist, evaluator.ErrorsCodesPolicy, "evaluator - errors codes evaluation policy")
	assert.Len(t, evaluator.ErrorsCodes, 0, "evaluator - errors codes list must be empty")
	assert.Equal(t, EvaluationPolicyWhitelist, evaluator.ErrorsSeveritiesPolicy, "evaluator - errors severities policy")
	assert.Len(t, evaluator.ErrorsSeverities, 1, "evaluator - errors severities list must have one element")
	assert.Equal(t, errors.SeverityRuntime, evaluator.ErrorsSeverities[0], "evaluator - error severity unknown")

	backoffCalculator, ok := retrier.BackoffCalculator.(*ExponentialBackoffCalculator)
	assert.True(t, ok, "backoff calculator cast")
	assert.Equal(t, 100*time.Millisecond, backoffCalculator.BaseBackoff, "backoff calculator - base backoff")
	assert.Equal(t, 20*time.Millisecond, backoffCalculator.RandomExtraBackoff, "backoff calculator - random extra backoff")
	assert.Equal(t, 2.0, backoffCalculator.Multiplier, "backoff calculator - multiplier")

	_, ok = retrier.ErrorWrapper.(*LastErrorWrapper)
	assert.True(t, ok, "error wrapper cast")
}

func TestRetrier_ExecuteOperation_Success(t *testing.T) {
	retrier := NewRetrier(
		Settings{
			BackoffCalculator: NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
				BaseBackoff:        1 * time.Nanosecond,
				Multiplier:         1.0,
				RandomExtraBackoff: 1 * time.Nanosecond,
			}),
		},
	)

	calls := 0
	err := retrier.ExecuteOperation(func() error {
		calls++
		if calls <= 2 {
			return errors.New("some error", errors.SeverityRuntime)
		}

		return nil
	})

	assert.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, 3, calls, "Operation calls count")
}

func TestRetrier_ExecuteOperation_CorrectMaxRetriesAttempt(t *testing.T) {
	retrier := NewRetrier(
		Settings{
			BackoffCalculator: NewExponentialBackoffCalculator(ExponentialBackoffCalculatorSettings{
				BaseBackoff:        1 * time.Nanosecond,
				Multiplier:         1.0,
				RandomExtraBackoff: 1 * time.Nanosecond,
			}),
			RetryEvaluator: NewGenericRetryEvaluator(GenericRetryEvaluatorSettings{
				MaxAttempts: 5,
			}),
		},
	)

	calls := 0
	err := retrier.ExecuteOperation(func() error {
		calls++
		return errors.New("some error", errors.SeverityRuntime)
	})

	assert.Error(t, err, "Execute should return an error")
	assert.Equal(t, 5, calls, "Operation calls count")
}
