package gitignore

import (
	"errors"
)

var (
	// define the standard lexer errors
	CarriageReturnError = errors.New("unexpected carriage return '\\r'")

	// define the standard parser errors
	InvalidPatternError = errors.New("invalid pattern")

	// define the standard repository errors
	InvalidDirectoryError = errors.New("invalid directory")
)
