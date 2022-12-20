package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error is the error struct that should be returned in all functions. It
// is ready to hold useful data about the error and has methods that make
// building and extracting information very easy.
// It also implements the Go's error interface
type Error struct {
	// Severity contains information about the nature of the error so you can
	// decide what to do
	Severity Severity
	// Err is the previous error, usually given by a lower level layer to be
	// wrapped
	Err error
	// Code is an application friendly way to describe the error
	Code Code
	// Op is the operation that was happening when the error occurred. In other
	// words, is where the error happened
	Op Op
	// KVs is a map os values you can use to enrich your error with relevant
	// information
	KVs []KeyValue
}

// Error formats the error information into a string. By implementing this
// method, we implement Go's error interface
func (e Error) Error() string {
	const sep = ": "
	s := strings.Builder{}
	kvb := strings.Builder{}
	if e.Op != "" {
		s.WriteString(string(e.Op))
		s.WriteString(sep)
	}
	if len(e.KVs) > 0 {
		for _, kv := range e.KVs {
			if kvb.Len() > 0 {
				kvb.WriteString(",")
			}
			kvb.WriteString(kv.String())
		}
	}
	var err = e.Err
	for {
		if err == nil {
			panic("root error is nil")
		}
		if innerErr, ok := err.(Error); ok {
			if innerErr.Op != "" {
				s.WriteString(string(innerErr.Op))
				s.WriteString(sep)
			}
			if len(innerErr.KVs) > 0 {
				for _, kv := range innerErr.KVs {
					if kvb.Len() > 0 {
						kvb.WriteString(",")
					}
					kvb.WriteString(kv.String())
				}
			}
			err = innerErr.Err
		} else {
			s.WriteString(err.Error())
			break
		}
	}
	if kvb.Len() > 0 {
		s.Grow(kvb.Len() + 2)
		s.WriteString(" [")
		s.WriteString(kvb.String())
		s.WriteString("]")
	}

	return s.String()
}

// String is required to implement the stringer interface
func (e Error) String() string {
	return e.Error()
}

// Unwrap returns the previous error
func (e Error) Unwrap() error {
	return e.Err
}

// E is a helper function for building errors.
//
// If called with no arguments, it returns an error solely containing a message
// informing that.
//
// This method can be called with any set of parameters of any type, but it
// requires either an error or a string at least, otherwise it will return nil.
//
// Parameters can be passed in any order and even be of repeated types, although
// only the last value of each type will be considered, except for KeyValue and
// []KeyValue, which will be concatenated to the struct's KVs value.
//
// Types other than string, Code, Severity, error, Op, KeyValue, or []KeyValue
// will simply be ignored.
func E(args ...interface{}) error {
	e := Error{}
	if len(args) == 0 {
		msg := "errors.E called with 0 args"
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = fmt.Sprintf("%v - %v:%v", msg, file, line)
		}
		e.Err = errors.New(msg)
	}

	for _, arg := range args {
		switch a := arg.(type) {
		case Code:
			e.Code = a
		case Severity:
			e.Severity = a
		case error:
			e.Err = a
		case string:
			e.Err = New(a)
		case Op:
			e.Op = a
		case KeyValue:
			e.KVs = append(e.KVs, a)
		case []KeyValue:
			e.KVs = append(e.KVs, a...)
		}
	}

	// If no error was provided, assume there was no error
	if e.Err == nil {
		return nil
	}

	return e
}

// New returns a new error. It is a wrap of Go's errors.New method that when
// given an empty string, returns nil
func New(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

// Errorf returns a error based on given params. It is a wrap of the fmt.Errorf
// method that returns nil when given an empty string
func Errorf(s string, params ...interface{}) error {
	if s == "" {
		return nil
	}
	return fmt.Errorf(s, params...)
}

// GetRootError returns the deepest error in the Err stack. That is, while an
// Error has a previous Error, keep getting the previous and returns when
// previous no longer has an Err
func GetRootError(err error) error {
	for {
		if myErr, ok := err.(Error); ok && myErr.Err != nil {
			err = myErr.Err
			continue
		}
		break
	}
	return err
}

// ConcatErrorsMessage concatenates all the error messages from the given errors
func ConcatErrorsMessage(errs ...error) string {
	s := strings.Builder{}
	for _, e := range errs {
		if s.Len() > 0 {
			s.WriteString(": ")
		}
		s.WriteString(e.Error())
	}
	return s.String()
}

// ConcatErrors returns an error with a message that is the concatenation of all
// messages of the given errors
func ConcatErrors(errs ...error) error {
	return New(ConcatErrorsMessage(errs...))
}
