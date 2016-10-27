package gitignore

import (
	"fmt"
)

// Position records the position of the .gitignore parser.
type Position struct {
	Offset int
	Column int
	Line   int
} // Position{}

// NewPosition returns the Position instance for the given line, column, and
// rune offset.
func NewPosition(line int, column int, offset int) Position {
	return Position{Offset: offset, Column: column, Line: line}
} // NewPosition()

// String returns a string representation of the current position.
func (p Position) String() string {
	if p.Line == 0 {
		return fmt.Sprintf("+%d", p.Offset)
	} else if p.Column == 0 {
		return fmt.Sprintf("%d", p.Line)
	} else {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
} // String()
