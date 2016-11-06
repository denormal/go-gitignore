package gitignore_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/denormal/go-gitignore"
)

func TestMatchAbsolute(t *testing.T) {
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
	_ignore := gitignore.NewGitIgnore(_buffer, _GITBASE, _error)
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

	// perform the absolute path matching
	for _, _test := range _GITMATCHES {
		match(t, _ignore, _GITBASE, _test)
	}
} // TestMatchAbsolute()

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
	_ignore := gitignore.NewGitIgnore(_buffer, _GITBASE, _error)
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

	// perform the absolute path matching
	for _, _test := range _GITMATCHES {
		match(t, _ignore, "", _test)
	}
} // TestMatchRelative()

func match(t *testing.T, i gitignore.GitIgnore, base string, test_ test) {
	var _match gitignore.Match

	// does this test represent a directory?
	//		- it does if the path ends in /
	_path := test_.Path
	_isdir := false
	if strings.HasSuffix(_path, "/") {
		_path = strings.TrimSuffix(_path, "/")
		_isdir = true
	}

	// are we matching a relative or absolute path?
	if base == "" {
		_match = i.Relative(_path, _isdir)
	} else {
		_path = filepath.Join(base, _path)
		_match = i.Absolute(_path, _isdir)
	}

	// attempt to match this path
	if _match == nil {
		// we have no match, is this expected?
		//		- a test that matches will list the expected pattern
		if test_.Pattern != "" {
			t.Errorf(
				"failed match; expected match for %q by %q",
				test_.Path, test_.Pattern,
			)
			return
		}

		// since we have no match, ensure this test path is not ignored
		if test_.Ignore {
			t.Errorf(
				"failed ignore; no match for %q but expected to be ignored",
				test_.Path,
			)
		}
	} else {
		// we have a match, is this expected?
		//		- a test that matches will list the expected pattern
		if test_.Pattern == "" {
			t.Errorf(
				"unexpected match by %q; expected no match for %q",
				_match, test_.Path,
			)
			return
		} else if test_.Pattern != _match.String() {
			t.Errorf(
				"mismatch for %q; expected match pattern %q, got %q",
				test_.Path, test_.Pattern, _match.String(),
			)
			return
		}

		// since we have a match, are we expected to ignore this file?
		if test_.Ignore != _match.Ignore() {
			t.Errorf(
				"ignore mismatch; expected %v for %q Ignore(), "+
					"got %v from pattern %q",
				test_.Ignore, test_.Path, _match.Ignore(), _match,
			)
		}
	}
} // match()
