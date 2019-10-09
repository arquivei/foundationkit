package nsu

import "github.com/arquivei/foundationkit/errors"

var (
	// ErrCodeNotFound should be returned by repository when the NSU was not found for the given stakeholder
	ErrCodeNotFound = errors.Code("NSU_NOT_FOUND")

	// ErrCannotParse is returned when a NSU cannot be parsed from a string
	ErrCannotParse = errors.New("failed to parse nsu from string")
)
