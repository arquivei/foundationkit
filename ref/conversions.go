package ref

import (
	"time"
)

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T {
	return &v
}

// PtrNil returns a pointer to v but if v is zero, returns nil.
func PtrNil[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// Val returns the value of a pointer. If the pointer is nil returns the zero value.
func Val[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

// Str gets the pointer of the @input
// Deprecated: use Ptr instead.
func Str(input string) *string {
	return &input
}

// AsStr gets the value of @input. It returns zero value if @input is nil.
// Deprecated: use Ptr instead.
func AsStr(input *string) string {
	if input != nil {
		return *input
	}
	return ""
}

// Bool gets the pointer of the @input
// Deprecated: use Ptr instead.
func Bool(input bool) *bool {
	return &input
}

// UInt64 gets the pointer of the @input
// Deprecated: use Ptr instead.
func UInt64(input uint64) *uint64 {
	return &input
}

// Time gets the pointer of the @input
// Deprecated: use Ptr instead.
func Time(input time.Time) *time.Time {
	return &input
}

// Float64 gets the pointer of the @input
// Deprecated: use Ptr instead.
func Float64(input float64) *float64 {
	return &input
}
