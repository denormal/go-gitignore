package gitignore

import (
	"errors"
)

var (
	// define the standard lexer errors
	CarriageReturnError = errors.New("unexpected carriage return '\\r'")
	EscapeError         = errors.New("unexpected escape '\\'")
	ContinuationError   = errors.New("unexpected EOF after continuation")
	IllegalRuneError    = errors.New("illegal character")

	// define the standard parser errors
	EOLError            = errors.New("unexpected end of line")
	EOFError            = errors.New("unexpected end of file")
	NegationError       = errors.New("unexpected negation '!'")
	CommentError        = errors.New("unexpected comment '#'")
	InvalidPatternError = errors.New("invalid pattern")
)
