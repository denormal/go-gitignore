package gitignore_test

import (
	"strings"
	"testing"

	"github.com/denormal/go-gitignore"
)

func TestMatchRelative(t *testing.T) {
	// create a temporary .gitignore
	_buffer, _err := buffer(_GITMATCH)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}

	// ensure we can run NewGitIgnore()
	//		- ensure we encounter no errors
	_position := []gitignore.Position{}
	_error := func(e gitignore.Error) bool {
		_position = append(_position, e.Position())
		return true
	}

	// ensure we have a non-nil GitIgnore instance
	_ignore := gitignore.NewGitIgnore(_buffer, "", _error)
	if _ignore == nil {
		t.Error("expected non-nil GitIgnore instance; nil found")
	}

	// ensure we encountered the right number of errors
	if len(_position) != _GITBADMATCHPATTERNS {
		t.Errorf(
			"match error mismatch; expected %d errors, got %d",
			_GITBADMATCHPATTERNS, len(_position),
		)
	}

	// test each of the defined matches
	for _, _test := range _GITMATCHES {
		// does this test represent a directory?
		//		- it does if the path ends in /
		_path := _test.Path
		_isdir := false
		if strings.HasSuffix(_path, "/") {
			_path = strings.TrimSuffix(_path, "/")
			_isdir = true
		}

		// attempt to match this path
		_match := _ignore.Relative(_path, _isdir)
		if _match == nil {
			// we have no match, is this expected?
			//		- a test that matches will list the expected pattern
			if _test.Pattern != "" {
				t.Errorf(
					"failed match; expected match for %q by %q",
					_test.Path, _test.Pattern,
				)
				continue
			}

			// since we have no match, ensure this test path is not ignored
			if _test.Ignore {
				t.Errorf(
					"failed ignore; no match for %q but expected to be ignored",
					_test.Path,
				)
			}
		} else {
			// we have a match, is this expected?
			//		- a test that matches will list the expected pattern
			if _test.Pattern == "" {
				t.Errorf(
					"unexpected match by %q; expected no match for %q",
					_match, _test.Path,
				)
				continue
			} else if _test.Pattern != _match.String() {
				t.Errorf(
					"mismatch for %q; expected match pattern %q, got %q",
					_test.Path, _test.Pattern, _match.String(),
				)
				continue
			}

			// since we have a match, are we expected to ignore this file?
			if _test.Ignore != _match.Ignore() {
				t.Errorf(
					"ignore mismatch; expected %v for %q Ignore(), "+
						"got %v from pattern %q",
					_test.Ignore, _test.Path, _match.Ignore(), _match,
				)
			}
		}
	}
} // TestMatchRelative()
