package gitignore

import (
	"fmt"
)

type Position struct {
	Offset int
	Column int
	Line   int
} // Position{}

func NewPosition(line int, column int, offset int) Position {
	return Position{Offset: offset, Column: column, Line: line}
} // NewPosition()

func (p Position) String() string {
	if p.Line == 0 {
		return fmt.Sprintf("+%d", p.Offset)
	} else if p.Column == 0 {
		return fmt.Sprintf("%d", p.Line)
	} else {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
} // String()
