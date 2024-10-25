package ref

import "time"

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

// Deprecated: Str is deprecated. Use `Of` instead.
func Str(input string) *string {
	return &input
}

// Deprecated: AsStr is deprecated. Use `Value` instead.
func AsStr(input *string) string {
	if input != nil {
		return *input
	}
	return ""
}

// Deprecated: Bool is deprecated. Use `Of` instead.
func Bool(input bool) *bool {
	return &input
}

// Deprecated: UInt64 is deprecated. Use `Of` instead.
func UInt64(input uint64) *uint64 {
	return &input
}

// Deprecated: Time is deprecated. Use `Of` instead.
func Time(input time.Time) *time.Time {
	return &input
}

// Deprecated: Float64 is deprecated. Use `Of` instead.
func Float64(input float64) *float64 {
	return &input
}
