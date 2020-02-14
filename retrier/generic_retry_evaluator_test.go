package retrier

import (
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewGenericRetryEvaluator_DefaultValues(t *testing.T) {
	evaluator := NewGenericRetryEvaluator(GenericRetryEvaluatorSettings{})
	assert.Equal(t, 5, evaluator.MaxAttempts, "max attempts")
	assert.Equal(t, EvaluationPolicyBlacklist, evaluator.ErrorsCodesPolicy, "errors codes evaluation policy")
	assert.Len(t, evaluator.ErrorsCodes, 0, "errors codes list")
	assert.Equal(t, EvaluationPolicyBlacklist, evaluator.ErrorsSeveritiesPolicy, "errors severity evaluation policy")
	assert.Len(t, evaluator.ErrorsSeverities, 0, "errors severity")
}

func TestGenericRetryEvaluator_IsRetryable(t *testing.T) {
	someErrCode := errors.Code("SOME_CODE")
	otherErrCode := errors.Code("OTHER_CODE")

	tests := []struct {
		name                string
		settings            GenericRetryEvaluatorSettings
		attemptNumber       int
		attemptError        error
		expectedIsRetryable bool
	}{
		{
			name: "Attempt less than max attempts is retryable",
			settings: GenericRetryEvaluatorSettings{
				MaxAttempts: 5,
			},
			attemptNumber:       4,
			attemptError:        errors.E("any error"),
			expectedIsRetryable: true,
		},
		{
			name: "Attempt same as max attempts is retryable",
			settings: GenericRetryEvaluatorSettings{
				MaxAttempts: 5,
			},
			attemptNumber:       5,
			attemptError:        errors.E("any error"),
			expectedIsRetryable: true,
		},
		{
			name: "Attempt more than max attempts is not retryable",
			settings: GenericRetryEvaluatorSettings{
				MaxAttempts: 5,
			},
			attemptNumber:       6,
			attemptError:        errors.E("any error"),
			expectedIsRetryable: false,
		},
		{
			name: "Error code in black list is not retryable",
			settings: GenericRetryEvaluatorSettings{
				ErrorsCodesPolicy: EvaluationPolicyBlacklist,
				ErrorsCodes:       []errors.Code{someErrCode},
			},
			attemptError:        errors.E(someErrCode, "some error"),
			expectedIsRetryable: false,
		},
		{
			name: "Error code not in white list is not retryable",
			settings: GenericRetryEvaluatorSettings{
				ErrorsCodesPolicy: EvaluationPolicyWhitelist,
				ErrorsCodes:       []errors.Code{someErrCode},
			},
			attemptError:        errors.E(otherErrCode, "other error"),
			expectedIsRetryable: false,
		},
		{
			name: "Error severity in black list is not retryable",
			settings: GenericRetryEvaluatorSettings{
				ErrorsSeveritiesPolicy: EvaluationPolicyBlacklist,
				ErrorsSeverities:       []errors.Severity{errors.SeverityFatal},
			},
			attemptError:        errors.E(errors.SeverityFatal, "fatal error"),
			expectedIsRetryable: false,
		},
		{
			name: "Error severity not in white list is not retryable",
			settings: GenericRetryEvaluatorSettings{
				ErrorsSeveritiesPolicy: EvaluationPolicyWhitelist,
				ErrorsSeverities:       []errors.Severity{errors.SeverityRuntime},
			},
			attemptError:        errors.E(errors.SeverityInput, "input error"),
			expectedIsRetryable: false,
		},
		{
			name: "Error is retryable when it passes both code and severity conditions",
			settings: GenericRetryEvaluatorSettings{
				ErrorsCodesPolicy:      EvaluationPolicyWhitelist,
				ErrorsCodes:            []errors.Code{someErrCode},
				ErrorsSeveritiesPolicy: EvaluationPolicyBlacklist,
				ErrorsSeverities:       []errors.Severity{errors.SeverityFatal},
			},
			attemptError:        errors.E(someErrCode, errors.SeverityRuntime, "informative error"),
			expectedIsRetryable: true,
		},
	}

	for _, test := range tests {
		evaluator := NewGenericRetryEvaluator(test.settings)
		isRetryable := evaluator.IsRetryable(test.attemptNumber, test.attemptError)

		assert.Equal(t, test.expectedIsRetryable, isRetryable, "[%s] is retryable", test.name)
	}
}
