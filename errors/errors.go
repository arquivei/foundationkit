package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error represents the error struct that should be returned in all functions
// Error implements the Go's error interface
type Error struct {
	Severity Severity
	Err      error
	Code     Code
	Op       Op
	KVs      []KeyValue
}

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

func (e Error) String() string {
	return e.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

// E is a helper function for builder errors
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

// New returns a new error
func New(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

// Errorf returns a error based on given params
func Errorf(s string, params ...interface{}) error {
	if s == "" {
		return nil
	}
	return fmt.Errorf(s, params...)
}

// GetRootError returns the Err field of Error struct or the error itself if it is of another type
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

// ConcatErrors concatenates all errors
func ConcatErrors(errs ...error) error {
	return New(ConcatErrorsMessage(errs...))
}
