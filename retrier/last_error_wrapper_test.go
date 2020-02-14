package retrier

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLastErrorWrapper(t *testing.T) {
	wrapper := NewLastErrorWrapper()

	assert.NoError(t, wrapper.WrapError(0, nil), "should not return error on nil error parameter")

	err := errors.New("my error")
	assert.Equal(t, err, wrapper.WrapError(0, err), "when receiving error, should return the same error")
}
