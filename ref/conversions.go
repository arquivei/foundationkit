package ref

import (
	"time"
)

// Str gets the pointer of the @input
func Str(input string) *string {
	return &input
}

// AsStr gets the value of @input. It returns zero value if @input is nil.
func AsStr(input *string) string {
	if input != nil {
		return *input
	}
	return ""
}

// Bool gets the pointer of the @input
func Bool(input bool) *bool {
	return &input
}

// UInt64 gets the pointer of the @input
func UInt64(input uint64) *uint64 {
	return &input
}

// Time gets the pointer of the @input
func Time(input time.Time) *time.Time {
	return &input
}

// Float64 gets the pointer of the @input
func Float64(input float64) *float64 {
	return &input
}
