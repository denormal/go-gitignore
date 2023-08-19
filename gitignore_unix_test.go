// Copyright 2016 Denormal Limited
// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !windows

package gitignore_test

import (
	"os"
	"testing"

	"github.com/ianlewis/go-gitignore"
)

// TODO(#18): Re-enable TestNewFromFile on Windows.
func TestNewFromFile(t *testing.T) {
	_test := &gitignoretest{}
	_test.instance = func(file string) (gitignore.GitIgnore, error) {
		return gitignore.NewFromFile(file)
	}

	// perform the gitignore test
	withfile(t, _test, _GITIGNORE)
} // TestNewFromFile()

// TODO(#18): Re-enable TestNewFromWhitespaceFile on Windows.
func TestNewFromWhitespaceFile(t *testing.T) {
	_test := &gitignoretest{}
	_test.instance = func(file string) (gitignore.GitIgnore, error) {
		return gitignore.NewFromFile(file)
	}

	// perform the gitignore test
	withfile(t, _test, _GITIGNORE_WHITESPACE)
} // TestNewFromWhitespaceFile()

// TODO(#18): Re-enable TestNewFromEmptyFile on Windows.
func TestNewFromEmptyFile(t *testing.T) {
	_test := &gitignoretest{}
	_test.instance = func(file string) (gitignore.GitIgnore, error) {
		return gitignore.NewFromFile(file)
	}

	// perform the gitignore test
	withfile(t, _test, "")
} // TestNewFromEmptyFile()

// TODO(#18): Re-enable TestNewWithErrors on Windows.
func TestNewWithErrors(t *testing.T) {
	_test := &gitignoretest{}
	_test.error = func(e gitignore.Error) bool {
		_test.errors = append(_test.errors, e)
		return true
	}
	_test.instance = func(file string) (gitignore.GitIgnore, error) {
		// reset the error slice
		_test.errors = make([]gitignore.Error, 0)

		// attempt to create the GitIgnore instance
		_ignore := gitignore.NewWithErrors(file, _test.error)

		// if we encountered errors, and the first error has a zero position
		// then it represents a file access error
		//		- extract the error and return it
		//		- remove it from the list of errors
		var _err error
		if len(_test.errors) > 0 {
			if _test.errors[0].Position().Zero() {
				_err = _test.errors[0].Underlying()
				_test.errors = _test.errors[1:]
			}
		}

		// return the GitIgnore instance
		return _ignore, _err
	}

	// perform the gitignore test
	withfile(t, _test, _GITIGNORE)

	_test.error = nil
	withfile(t, _test, _GITIGNORE)
} // TestNewWithErrors()

// TODO(#18): Re-enable TestNewWithCache on Windows.
func TestNewWithCache(t *testing.T) {
	// perform the gitignore test with a custom cache
	_test := &gitignoretest{}
	_test.cached = true
	_test.cache = gitignore.NewCache()
	_test.instance = func(file string) (gitignore.GitIgnore, error) {
		// reset the error slice
		_test.errors = make([]gitignore.Error, 0)

		// attempt to create the GitIgnore instance
		_ignore := gitignore.NewWithCache(file, _test.cache, _test.error)

		// if we encountered errors, and the first error has a zero position
		// then it represents a file access error
		//		- extract the error and return it
		//		- remove it from the list of errors
		var _err error
		if len(_test.errors) > 0 {
			if _test.errors[0].Position().Zero() {
				_err = _test.errors[0].Underlying()
				_test.errors = _test.errors[1:]
			}
		}

		// return the GitIgnore instance
		return _ignore, _err
	}

	// perform the gitignore test
	withfile(t, _test, _GITIGNORE)

	// repeat the tests while accumulating errors
	_test.error = func(e gitignore.Error) bool {
		_test.errors = append(_test.errors, e)
		return true
	}
	withfile(t, _test, _GITIGNORE)

	// create a temporary .gitignore
	_file, _err := file(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	// attempt to load the .gitignore file
	_ignore, _err := _test.instance(_file.Name())
	if _err != nil {
		t.Fatalf("unable to open temporary .gitignore: %s", _err.Error())
	}

	// remove the .gitignore and try again
	os.Remove(_file.Name())

	// ensure the retrieved GitIgnore matches the stored instance
	_new, _err := _test.instance(_file.Name())
	if _err != nil {
		t.Fatalf(
			"unexpected error retrieving cached .gitignore: %s", _err.Error(),
		)
	} else if _new != _ignore {
		t.Fatalf(
			"gitignore.NewWithCache() mismatch; expected %v, got %v",
			_ignore, _new,
		)
	}
} // TestNewWithCache()
