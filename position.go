package gitignore

import (
	"fmt"
)

// Position records the position of the .gitignore parser.
type Position struct {
	Line   int
	Column int
	Offset int
} // Position{}

// NewPosition returns the Position instance for the given line, column, and
// rune offset.
func NewPosition(line int, column int, offset int) Position {
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
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
