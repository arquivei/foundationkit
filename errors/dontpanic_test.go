package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDontPanic(t *testing.T) {
	tests := []struct {
		name      string
		expected  string
		panicFunc func()
	}{
		{
			name:     "panic()",
			expected: "errors/dontpanic_test.go:20: panic: panic message aaaaaah!",
			panicFunc: func() {
				panic("panic message aaaaaah!")
			},
		},
		{
			name:     "index out of range",
			expected: "errors/dontpanic_test.go:28: panic: runtime error: index out of range [1] with length 0",
			panicFunc: func() {
				var s []string
				ss := s[1]
				fmt.Println(ss)
			},
		},
		{
			name:     "invalid memory address or nil pointer dereference",
			expected: "errors/dontpanic_test.go:40: panic: runtime error: invalid memory address or nil pointer dereference",
			panicFunc: func() {
				type a struct {
					b string
				}
				var aa *a
				fmt.Println(aa.b)
			},
		},
		{
			name:     "panic() inside a func inside a func",
			expected: "errors/dontpanic_test.go:48: panic: panic message aaaaaah!",
			panicFunc: func() {
				anotherFunc1 := func() {
					panic("panic message aaaaaah!")
				}
				anotherFunc2 := func() {
					anotherFunc1()
				}
				anotherFunc2()
			},
		},
		{
			name:     "integer divide by zero",
			expected: "errors/dontpanic_test.go:63: panic: runtime error: integer divide by zero",
			panicFunc: func() {
				d := func(i int) {
					fmt.Println(2 / i)
				}
				d(0)
			},
		},
		{
			name:     "send on closed channel",
			expected: "errors/dontpanic_test.go:72: panic: send on closed channel",
			panicFunc: func() {
				c := make(chan int)
				close(c)
				c <- 1
			},
		},
	}

	//nolint:gocritic
	func4 := func(f func()) error {
		return DontPanic(f)
	}
	//nolint:gocritic
	func3 := func(f func()) error {
		return func4(f)
	}
	//nolint:gocritic
	func2 := func(f func()) error {
		return func3(f)
	}
	//nolint:gocritic
	func1 := func(f func()) error {
		return func2(f)
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actual := func1(test.panicFunc)
			assert.EqualError(t, actual, test.expected)
		})
	}
}

func TestNewOpFromPanicStack_Recover(t *testing.T) {
	op := newOpFromPanicStack()
	assert.Equal(t, Op(""), op)

	err := E(op, "new error without op")
	assert.EqualError(t, err, "new error without op")
}

func TestNewFromRecover_UsingError(t *testing.T) {
	err := NewFromRecover(E(
		Op("TestNewFromRecover"),
		"new error",
		SeverityInput,
		Code("CODE"),
	))
	assert.EqualError(t, err, "TestNewFromRecover: new error")
	assert.Equal(t, CodePanic, GetCode(err))
	assert.Equal(t, SeverityFatal, GetSeverity(err))
}
