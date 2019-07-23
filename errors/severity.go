package errors

// Severity is the error severity. It's used to classify errors in groups to be easily handled by the code. For examle,
// a retry layer should be only checking for Runtime erros to retry. Or in an HTTP layer, erros of input type are always
// returned a 400 status.
type Severity string

const (
	//SeverityUnset indicates the severity was not set
	SeverityUnset = Severity("")
	// SeverityRuntime indicates the error is returned for an operation that should/could be executed again. For example, timeouts.
	SeverityRuntime = Severity("runtime")
	// SeverityFatal indicates the error is unrecoverable and the execution should stop, or not being retried.
	SeverityFatal = Severity("fatal")
	// SeverityInput indicates  an expected, like a bad user input/request. For example, an invalid email.
	SeverityInput = Severity("input")
)

func (s Severity) String() string {
	return string(s)
}

// GetSeverity returns the error severity. If there is not severity, SeverityUnset is returned.
func GetSeverity(err error) Severity {
	for {
		e, ok := err.(Error)
		if !ok {
			break
		}
		if e.Severity != SeverityUnset {
			return e.Severity
		}
		err = e.Err
	}

	return SeverityUnset
}
