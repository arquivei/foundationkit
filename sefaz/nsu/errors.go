package nsu

import "github.com/arquivei/foundationkit/errors"

var (
	// ErrCannotParse is returned when a NSU cannot be parsed from a string
	ErrCannotParse = errors.New("failed to parse nsu from string")
)
