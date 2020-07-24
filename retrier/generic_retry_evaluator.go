package retrier

import "github.com/arquivei/foundationkit/errors"

// EvaluationPolicy allows indicating to the generic retry evaluator if it
// should work with blacklists or whitelists
type EvaluationPolicy int

const (
	// EvaluationPolicyBlacklist indicates that an evaluation will only succeed if
	// the information being verified is NOT present in a given list
	EvaluationPolicyBlacklist EvaluationPolicy = iota
	// EvaluationPolicyWhitelist indicates that an evaluation will only succeed if
	// the information being verified is present in a given list
	EvaluationPolicyWhitelist
)

// GenericRetryEvaluator is a retry evaluator that works with max attempts and
// both error codes and errors severity. Error codes and errors severity are both
// based on whitelist or blacklists, and are evaluated individually. If a match
// happens on blacklists, the error is assumed as non retryable. If a match happens
// on a whitelist, the errors is assumed as retryable. Both errors codes and errors
// severity must pass the list tests in order of the error to be retryable.
type GenericRetryEvaluator struct {
	MaxAttempts            int
	ErrorsCodesPolicy      EvaluationPolicy
	ErrorsCodes            []errors.Code
	ErrorsSeveritiesPolicy EvaluationPolicy
	ErrorsSeverities       []errors.Severity
}

// GenericRetryEvaluatorSettings is used to construct GenericRetryEvaluator instances.
// All fields with zero value with receive default values
type GenericRetryEvaluatorSettings struct {
	// MaxAttempts is how many attempts to allow, starting at 1. Defaults to 5.
	MaxAttempts int
	// ErrorCodesPolicy indicates if error codes should be held with
	// blacklists or whitelists. Defaults to blacklist.
	ErrorsCodesPolicy EvaluationPolicy
	// ErrorsCodesList is the list of errors codes to use as a base. Defaults to empty.
	ErrorsCodes []errors.Code
	// ErrorSeveritiesPolicy indicates if error severity should be held with
	// blacklists or whitelists. Defaults to blacklist.
	ErrorsSeveritiesPolicy EvaluationPolicy
	// ErrorsSeveritiesList is the list of errors severity to use as a base. Defaults to empty.
	ErrorsSeverities []errors.Severity
}

// NewGenericRetryEvaluator will return an instance of generic retry evaluator
func NewGenericRetryEvaluator(settings GenericRetryEvaluatorSettings) *GenericRetryEvaluator {
	if settings.MaxAttempts == 0 {
		settings.MaxAttempts = 5
	}

	return &GenericRetryEvaluator{
		MaxAttempts:            settings.MaxAttempts,
		ErrorsCodesPolicy:      settings.ErrorsCodesPolicy,
		ErrorsCodes:            settings.ErrorsCodes,
		ErrorsSeveritiesPolicy: settings.ErrorsSeveritiesPolicy,
		ErrorsSeverities:       settings.ErrorsSeverities,
	}
}

// IsRetryable will return true if the @err in the given @attempt can be retried,
// or false otherwise.
//
// This function works with preset black or white lists to error codes and errors
// severity. See type definition for more information.
func (e *GenericRetryEvaluator) IsRetryable(attempt int, attemptError error) bool {
	const op = errors.Op("retrier.GenericRetryEvaluator.IsRetryable")

	if attempt > e.MaxAttempts {
		return false
	}

	canRetryOnErrorCode, err := isErrorCodeRetryable(errors.GetCode(attemptError), e.ErrorsCodesPolicy, e.ErrorsCodes)
	if err != nil {
		panic(errors.E(op, err))
	}

	canRetryOnErrorSeverity, err := isErrorSeverityRetryable(errors.GetSeverity(attemptError), e.ErrorsSeveritiesPolicy, e.ErrorsSeverities)
	if err != nil {
		panic(errors.E(op, err))
	}

	return canRetryOnErrorCode && canRetryOnErrorSeverity
}

func isErrorCodeRetryable(
	errCode errors.Code,
	evaluationPolicy EvaluationPolicy,
	evaluationCodes []errors.Code,
) (bool, error) {
	const op = errors.Op("isErrorCodeRetryable")

	errCodeFound := false

	for _, listCode := range evaluationCodes {
		if errCode == listCode {
			errCodeFound = true
			break
		}
	}

	switch evaluationPolicy {
	case EvaluationPolicyBlacklist:
		if errCodeFound {
			return false, nil
		}
	case EvaluationPolicyWhitelist:
		if !errCodeFound {
			return false, nil
		}
	default:
		return false, errors.E(op, "bad evaluation policy")
	}

	return true, nil
}

func isErrorSeverityRetryable(
	errSeverity errors.Severity,
	evaluationPolicy EvaluationPolicy,
	evaluationSeverities []errors.Severity,
) (bool, error) {
	const op = errors.Op("isErrorSeverityRetryable")

	errSeverityFound := false

	for _, listSeverity := range evaluationSeverities {
		if errSeverity == listSeverity {
			errSeverityFound = true
			break
		}
	}

	switch evaluationPolicy {
	case EvaluationPolicyBlacklist:
		if errSeverityFound {
			return false, nil
		}
	case EvaluationPolicyWhitelist:
		if !errSeverityFound {
			return false, nil
		}
	default:
		return false, errors.E(op, "bad evaluation policy")
	}

	return true, nil
}
