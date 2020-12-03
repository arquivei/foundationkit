package errors

import (
	"runtime/debug"
	"strings"
)

// CodePanic represents the error code for panic
var CodePanic Code = "PANIC"

// DontPanic executes f and, if f panics, recovers from the panic
// and returns the panic wrapped as an Error.
func DontPanic(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewFromRecover(r)
			err = E(newOpFromPanicStack(), err)
		}
	}()

	f()
	return
}

func newOpFromPanicStack() (op Op) {
	defer func() {
		if r := recover(); r != nil {
			op = ""
		}
	}()

	stack := string(debug.Stack())

	runtimePanicSplit := strings.SplitAfterN(stack, "runtime/panic.go:", 2)
	panicFuncSplit := strings.SplitAfterN(runtimePanicSplit[1], "\n\t", 2)
	lastFileSplit := strings.SplitN(panicFuncSplit[1], " +", 2)
	pathSplit := strings.Split(lastFileSplit[0], "/")
	packageName := pathSplit[len(pathSplit)-2]
	fileAndLine := pathSplit[len(pathSplit)-1]

	return Op(packageName + "/" + fileAndLine)
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
