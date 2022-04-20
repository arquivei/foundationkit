package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Due to the non-determinism nature functions , it is really hard to produce
// more complex tests, as the fields are not guaranteed to be in any order
func TestFlatten(t *testing.T) {
	structA := struct {
		FooA string
	}{
		FooA: "abC%#1s",
	}
	expectedFlattenA := `FooA: 'abC%#1s'`
	assert.Equal(t, expectedFlattenA, Flatten(structA), "test failed for struct A")

	structB := struct {
		FooB int
	}{
		FooB: 2,
	}
	expectedFlattenB := `FooB: '2'`
	assert.Equal(t, expectedFlattenB, Flatten(structB), "test failed for struct B")

	structC := struct {
		FooC bool
	}{
		FooC: false,
	}
	expectedFlattenC := `FooC: 'false'`
	assert.Equal(t, expectedFlattenC, Flatten(structC), "test failed for struct C")

	structD := struct {
		FooD interface{}
	}{
		FooD: "5",
	}
	expectedFlattenD := `FooD: '5'`
	assert.Equal(t, expectedFlattenD, Flatten(structD), "test failed for struct D")

	structE := struct {
		FooE1 interface{}
		FooE2 interface{} `secret:"true"`
	}{
		FooE1: "5",
		FooE2: "q1w2e3",
	}
	expectedFlattenE := `FooE1: '5'`
	assert.Equal(t, expectedFlattenE, Flatten(structE), "test failed for struct E")

	assert.NotPanics(t, func() {
		Flatten(nil)
	})

	assert.Equal(t, "", Flatten(nil))
}
