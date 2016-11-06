package gitignore_test

import (
	"testing"

	"os"
	"path/filepath"

	"github.com/denormal/go-gitignore"
)

func TestNewFromFile(t *testing.T) {
	// create a temporary .gitignore
	_file, _err := file(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	// ensure we can run NewFromFile()
	_ignore, _err := gitignore.NewFromFile(_file.Name())
	if _err != nil {
		t.Fatalf("unable to open temporary .gitignore: %s", _err.Error())
	}

	// ensure we have a non-nil GitIgnore instance
	if _ignore == nil {
		t.Error("expected non-nil GitIgnore instance; nil found")
	}

	// ensure the base of the ignore is the directory of the temporary file
	_dir := filepath.Dir(_file.Name())
	if _ignore.Base() != _dir {
		t.Errorf(
			"gitignore.Base() mismatch; expected %q, got %q",
			_dir, _ignore.Base(),
		)
	}
} // TestNewFromFile()

func TestNewWithCache(t *testing.T) {
	// create a cache for this test
	_cache := gitignore.NewCache()

	// create a temporary .gitignore
	_file, _err := file(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	//  ensure we can run NewFromFile()
	_ignore, _err := gitignore.NewFromFile(_file.Name())
	if _err != nil {
		t.Fatalf("unable to open temporary .gitignore: %s", _err.Error())
	}

	// ensure we have a non-nil GitIgnore instance
	if _ignore == nil {
		t.Error("expected non-nil GitIgnore instance; nil found")
	}

	// store the GitIgnore in the cache
	_cache.Set(_file.Name(), _ignore)

	// remove the .gitignore file from disk
	_err = os.Remove(_file.Name())
	if _err != nil {
		t.Errorf("unable to remove temporary .gitignore: %s", _err.Error())
	}

	// attempt to create this .gitignore again, but this time using the cache
	//		- ensure the retrieved GitIgnore matches the stored instance
	_new, _err := gitignore.NewWithCache(_file.Name(), _cache)
	if _err != nil {
		t.Errorf(
			"unexpected error retrieving cached .gitignore: %s", _err.Error(),
		)
		return
	} else if _new != _ignore {
		t.Errorf(
			"gitignore.NewWithCache() mismatch; expected %v, got %v",
			_ignore, _new,
		)
		return
	}

	// ensure the base of the ignore is the directory of the temporary file
	_dir := filepath.Dir(_file.Name())
	if _new.Base() != _dir {
		t.Errorf(
			"gitignore.Base() mismatch; expected %q, got %q",
			_dir, _new.Base(),
		)
	}
} // TestNewWithCache()

func TestNew(t *testing.T) {
	// create a temporary .gitignore
	_file, _err := file(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	// ensure we can run NewGitIgnore()
	//		- ensure we encounter 2 errors
	_position := []gitignore.Position{}
	_error := func(e gitignore.Error) bool {
		_position = append(_position, e.Position())
		return true
	}

	_dir := filepath.Dir(_file.Name())
	_ignore := gitignore.New(_file, _dir, _error)

	// ensure we have a non-nil GitIgnore instance
	if _ignore == nil {
		t.Error("expected non-nil GitIgnore instance; nil found")
	}

	// ensure the base of the ignore is the directory of the temporary file
	if _ignore.Base() != _dir {
		t.Errorf(
			"gitignore.Base() mismatch; expected %q, got %q",
			_dir, _ignore.Base(),
		)
	}

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

			// ensure the positions are correct
			if !coincident(_got, _expected) {
				t.Errorf("bad pattern position mismatch; expected %q, got %q",
					position(_expected), position(_got),
				)
			}
		}
	}
} // TestNew()
