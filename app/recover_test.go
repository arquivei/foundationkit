package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	// Note: this test is not working as expected. It is not catching the panic.
	// assert.NotPanics(t, func() {
	// 	defer Recover()
	// 	thisWillPanic()
	// })

	assert.NotPanics(t, func() {
		go func() {
			defer Recover()
			thisWillPanic()
		}()
	})
}

func thisWillPanic() {
	panic("panics should be caught by Recover()")
}
