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
	"path/filepath"
	"testing"

	"github.com/ianlewis/go-gitignore"
)

// TODO(#17): Re-enable TestMatchAbsolute on windows.
func TestMatchAbsolute(t *testing.T) {
	// create a temporary .gitignore
	_buffer, _err := buffer(_GITMATCH)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}

	// ensure we can run New()
	//		- ensure we encounter no errors
	_position := []gitignore.Position{}
	_error := func(e gitignore.Error) bool {
		_position = append(_position, e.Position())
		return true
	}

	// ensure we have a non-nil GitIgnore instance
	_ignore := gitignore.New(_buffer, _GITBASE, _error)
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
	_cb := func(path string, isdir bool) gitignore.Match {
		_path := filepath.Join(_GITBASE, path)
		return _ignore.Absolute(_path, isdir)
	}
	for _, _test := range _GITMATCHES {
		do(t, _cb, _test)
	}

	// perform absolute path tests with paths not under the same root
	// directory as the GitIgnore we are testing
	_new, _ := directory(t)
	defer os.RemoveAll(_new)

	for _, _test := range _GITMATCHES {
		_path := filepath.Join(_new, _test.Local())
		_match := _ignore.Match(_path)
		if _match != nil {
			t.Fatalf("unexpected match; expected nil, got %v", _match)
		}
	}
} // TestMatchAbsolute()
