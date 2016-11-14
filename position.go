package gitignore

import (
	"fmt"
)

// Position records the position of the .gitignore parser.
type Position struct {
	File   string
	Line   int
	Column int
	Offset int
} // Position{}

// NewPosition returns the Position instance for the given line, column, and
// rune offset.
func NewPosition(file string, line int, column int, offset int) Position {
	return Position{
		File:   file,
		Line:   line,
		Column: column,
		Offset: offset,
	}
} // NewPosition()

// String returns a string representation of the current position.
func (p Position) String() string {
	_prefix := ""
	if p.File != "" {
		_prefix = p.File + ": "
	}

	if p.Line == 0 {
		return fmt.Sprintf("%s+%d", _prefix, p.Offset)
	} else if p.Column == 0 {
		return fmt.Sprintf("%s%d", _prefix, p.Line)
	} else {
		return fmt.Sprintf("%s%d:%d", _prefix, p.Line, p.Column)
	}
} // String()

// Zero returns true if the Position represents the zero Position
func (p Position) Zero() bool {
	return p.Line+p.Column+p.Offset == 0
} // Zero()
