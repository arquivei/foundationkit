package errors

// Op is the operation that encapsulated the error
type Op string

func (op Op) String() string {
	return string(op)
}

func (op Op) Apply(err *Error) {
	err.Op = op
}
