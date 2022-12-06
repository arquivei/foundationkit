package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		shoulPanic bool
	}{
		{
			name:       "Recover should recover from panic and log it",
			shoulPanic: true,
		},

		{
			name:       "Recover should not recover if there is no panic",
			shoulPanic: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Panics(t, thisWillPanic)
			assert.NotPanics(t, func() {
				go func() {
					defer Recover()
					if test.shoulPanic {
						thisWillPanic()
					}
				}()
			})
		})
	}
}

func thisWillPanic() {
	panic("panics should be caught by Recover()")
}
