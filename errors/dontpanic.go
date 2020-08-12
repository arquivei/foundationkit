package errors

var CodePanic Code = "PANIC"

// DontPanic executes f and, if f panics, recovers from the panic
// and returns the panic wrapped as an Error.
func DontPanic(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewFromRecover(r)
		}
	}()

	f()
	return
}

// NewFromRecover returns a new Error created from the result of a recover.
// If r is an Error, this will be used so Op and KV are preserved
func NewFromRecover(r interface{}) Error {
	var err error

	if rr, ok := r.(Error); ok {
		err = rr
	} else {
		err = Errorf("panic: %v", r)
	}

	return Error{
		Err:      err,
		Severity: SeverityFatal,
		Code:     CodePanic,
	}
}
