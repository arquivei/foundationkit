package ref

// Of returns a reference pointer of the @input
func Of[T any](input T) *T {
	return &input
}

// Value returns a value of the @input type. It the @input is nil, it returns the zero value of the type
func Value[T any](value *T) T {
	if value != nil {
		return *value
	}

	var zero T
	return zero
}
