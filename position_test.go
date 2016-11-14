package gitignore_test

import (
	"github.com/denormal/go-gitignore"
	"testing"
)

func TestPosition(t *testing.T) {
	// test the conversion of Positions to strings
	for _, _p := range _POSITIONS {
		_position := gitignore.NewPosition(
			_p.File, _p.Line, _p.Column, _p.Offset,
		)

		// ensure the string representation of the Position is as expected
		_rtn := _position.String()
		if _rtn != _p.String {
			t.Error(
				"position mismatch; expected %q, got %q",
				_p.String, _rtn,
			)
		}
	}
} // TestPosition()
