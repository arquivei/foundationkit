package errors

// Op is the operation that encapsulated the error
type Op string

func (o Op) String() string {
	return string(o)
}
