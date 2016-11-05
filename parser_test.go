package gitignore_test

import (
	"strings"
	"testing"

	"github.com/denormal/go-gitignore"
)

// TestParser tests the behaviour of gitignore.Parser
func TestParser(t *testing.T) {
	// create a temporary .gitignore
	_buffer, _err := buffer(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}

	// ensure we can run NewGitIgnore()
	//		- ensure we encounter 2 errors
	_position := []gitignore.Position{}
	_error := func(e gitignore.Error) bool {
		_position = append(_position, e.Position())
		return true
	}

	// ensure we have a non-nil Parser instance
	_parser := gitignore.NewParser(_buffer, _error)
	if _parser == nil {
		t.Error("expected non-nil Parser instance; nil found")
	}

	// attempt to parse the .gitignore
	_patterns := _parser.Parse()

	// ensure we encountered the right number of errors
	if len(_position) != _GITBADPATTERNS {
		t.Errorf(
			"parse error mismatch; expected %d errors, got %d",
			_GITBADPATTERNS, len(_position),
		)
	} else {
		// ensure the error positions are correct
		for _i := 0; _i < _GITBADPATTERNS; _i++ {
			_got := _position[_i]
			_expected := _GITBADPOSITION[_i]

			if !coincident(_got, _expected) {
				t.Errorf(
					"bad pattern position mismatch; expected %q, got %q",
					position(_expected), position(_got),
				)
			}
		}
	}

	// ensure we encountered the right number of patterns
	if len(_patterns) != _GITPATTERNS {
		t.Errorf(
			"parse pattern mismatch; expected %d patterns, got %d",
			_GITPATTERNS, len(_patterns),
		)
	} else {
		// ensure the pattern positions are correct
		for _i := 0; _i < len(_patterns); _i++ {
			_got := _patterns[_i].Position()
			_expected := _GITPOSITION[_i]

			if !coincident(_got, _expected) {
				t.Errorf(
					"pattern position mismatch; expected %q, got %q",
					position(_expected), position(_got),
				)
			}
		}

		// ensure the retrieved patterns are correct
		_lines := strings.Split(_GITIGNORE, "\n")
		for _i := 0; _i < len(_patterns); _i++ {
			_pattern := _patterns[_i]
			_got := _pattern.String()
			_line := _pattern.Position().Line
			_expected := _lines[_line-1]

			if _got != _expected {
				t.Errorf(
					"pattern mismatch; expected %q, got %q at %s",
					_expected, _got, position(_pattern.Position()),
				)
			}
		}
	}
} // TestParser()

// TestParserError tests the behaviour of the gitignore.Parser with an error
// handler that returns false on receiving and error
func TestParserError(t *testing.T) {
	// create a temporary .gitignore
	_buffer, _err := buffer(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}

	// ensure we can run NewGitIgnore()
	//		- ensure we encounter 2 errors
	_position := []gitignore.Position{}
	_error := func(e gitignore.Error) bool {
		_position = append(_position, e.Position())
		return false
	}

	// ensure we have a non-nil Parser instance
	_parser := gitignore.NewParser(_buffer, _error)
	if _parser == nil {
		t.Error("expected non-nil Parser instance; nil found")
	}

	// attempt to parse the .gitignore
	_patterns := _parser.Parse()

	// ensure we encountered only one bad pattern
	if len(_position) != _GITBADPATTERNSFALSE {
		t.Errorf(
			"parse error mismatch; expected %d errors, got %d",
			1, len(_position),
		)
	} else {
		// ensure the error positions are correct
		for _i := 0; _i < _GITBADPATTERNSFALSE; _i++ {
			_got := _position[_i]
			_expected := _GITBADPOSITION[_i]

			if !coincident(_got, _expected) {
				t.Errorf(
					"bad pattern position mismatch; expected %q, got %q",
					position(_expected), position(_got),
				)
			}
		}
	}

	// ensure we encountered the right number of patterns
	if len(_patterns) != _GITPATTERNSFALSE {
		t.Errorf(
			"parse pattern mismatch; expected %d patterns, got %d",
			_GITPATTERNS, len(_patterns),
		)
	} else {
		// ensure the pattern positions are correct
		for _i := 0; _i < len(_patterns); _i++ {
			_got := _patterns[_i].Position()
			_expected := _GITPOSITION[_i]

			if !coincident(_got, _expected) {
				t.Errorf(
					"pattern position mismatch; expected %q, got %q",
					position(_expected), position(_got),
				)
			}
		}

		// ensure the retrieved patterns are correct
		_lines := strings.Split(_GITIGNORE, "\n")
		for _i := 0; _i < len(_patterns); _i++ {
			_pattern := _patterns[_i]
			_got := _pattern.String()
			_line := _pattern.Position().Line
			_expected := _lines[_line-1]

			if _got != _expected {
				t.Errorf(
					"pattern mismatch; expected %q, got %q at %s",
					_expected, _got, position(_pattern.Position()),
				)
			}
		}
	}
} // TestParserError()
