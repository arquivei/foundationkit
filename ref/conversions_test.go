package ref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	strPtr := Ptr("some string")

	assert.NotNil(t, strPtr)
	assert.Equal(t, "some string", *strPtr)
}

func TestPtrNil(t *testing.T) {
	nilStr := PtrNil("")
	assert.Nil(t, nilStr)

	someString := PtrNil("some string")
	assert.NotNil(t, someString)
	assert.Equal(t, "some string", *someString)
}

func TestVal(t *testing.T) {
	var nilStr *string
	assert.Equal(t, "", nilStr)

	someString := "some string"
	assert.Equal(t, "some string", Val(&someString))
}

// Deprecated tests

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
