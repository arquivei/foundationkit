package ref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOfStr(t *testing.T) {
	value := "abcde"
	assert.Equal(t, Of(value), &value)
}

func TestOfBool(t *testing.T) {
	valueTrue := true
	assert.Equal(t, Of(valueTrue), &valueTrue)

	valueFalse := false
	assert.Equal(t, Of(valueFalse), &valueFalse)
}

func TestOfFloat64(t *testing.T) {
	value := 1.2
	assert.Equal(t, Of(value), &value)
}

func TestValue(t *testing.T) {
	assert.Equal(t, Value(Of("")), "")
	assert.Equal(t, Value(Of("bla")), "bla")
	assert.Equal(t, Value[string](nil), "")
}

func TestOfTime(t *testing.T) {
	value := time.Now()
	assert.Equal(t, Of(value), &value)
}

func TestOfUInt64(t *testing.T) {
	value := uint64(1)
	assert.Equal(t, Of(value), &value)
}
