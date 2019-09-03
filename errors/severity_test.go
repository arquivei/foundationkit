package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSeverity(t *testing.T) {
	// Basic error
	var err error = Error{
		Severity: SeverityInput,
		Err:      New("my error"),
	}

	assert.Equal(t, SeverityInput, GetSeverity(err))
	assert.EqualError(t, err, "my error")

	// Wraps error with runtime severity
	err = Error{
		Severity: SeverityRuntime,
		Err:      err,
	}

	assert.Equal(t, SeverityRuntime, GetSeverity(err))
	assert.EqualError(t, err, "my error")

	// Wraps error with an error with without severity
	// Keeps previous severity
	err = Error{
		Err: err,
	}

	assert.Equal(t, SeverityRuntime, GetSeverity(err))
	assert.EqualError(t, err, "my error")
}

func TestUnsetSeverity(t *testing.T) {
	// Basic error
	var err error = Error{
		Err: New("my error"),
	}

	assert.Equal(t, SeverityUnset, GetSeverity(err))
	assert.EqualError(t, err, "my error")
}

func TestSeverityText(t *testing.T) {
	assert.Equal(t, "runtime", SeverityRuntime.String())
	assert.Equal(t, "input", SeverityInput.String())
	assert.Equal(t, "fatal", SeverityFatal.String())
	assert.Equal(t, "", SeverityUnset.String())
	assert.Equal(t, "other", Severity("other").String())
}
