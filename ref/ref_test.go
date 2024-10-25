package ref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStr(t *testing.T) {
	value := "abcde"
	assert.Equal(t, Str(value), &value)
}

func TestBool(t *testing.T) {
	valueTrue := true
	assert.Equal(t, Bool(valueTrue), &valueTrue)

	valueFalse := false
	assert.Equal(t, Bool(valueFalse), &valueFalse)
}

func TestFloat64(t *testing.T) {
	value := 1.2
	assert.Equal(t, Float64(value), &value)
}

func TestAsStr(t *testing.T) {
	value := "bla"
	assert.Equal(t, AsStr(&value), value)
	assert.Equal(t, AsStr(nil), "")
}

func TestTime(t *testing.T) {
	value := time.Now()
	assert.Equal(t, Time(value), &value)
}

func TestUInt64(t *testing.T) {
	value := uint64(1)
	assert.Equal(t, UInt64(value), &value)
}

func TestOf_Str(t *testing.T) {
	value := "abcde"
	assert.Equal(t, Of(value), &value)
}

func TestOf_Bool(t *testing.T) {
	valueTrue := true
	assert.Equal(t, Of(valueTrue), &valueTrue)

	valueFalse := false
	assert.Equal(t, Of(valueFalse), &valueFalse)
}

func TestOf_Float64(t *testing.T) {
	value := 1.2
	assert.Equal(t, Of(value), &value)
}

func TestValue(t *testing.T) {
	assert.Equal(t, Value(Of("")), "")
	assert.Equal(t, Value(Of("bla")), "bla")
	assert.Equal(t, Value[string](nil), "")
}

func TestOf_Time(t *testing.T) {
	value := time.Now()
	assert.Equal(t, Of(value), &value)
}

func TestOf_UInt64(t *testing.T) {
	value := uint64(1)
	assert.Equal(t, Of(value), &value)
}
