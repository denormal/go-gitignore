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

package gitignore

import (
	"os"
	"path/filepath"
)

// exclude attempts to return the GitIgnore instance for the
// $GIT_DIR/info/exclude from the working copy to which path belongs.
func exclude(path string) (GitIgnore, error) {
	// attempt to locate GIT_DIR
	_gitdir := os.Getenv("GIT_DIR")
	if _gitdir == "" {
		_gitdir = filepath.Join(path, ".git")
	}
	_info, _err := os.Stat(_gitdir)
	if _err != nil {
		if os.IsNotExist(_err) {
			return nil, nil
		} else {
			return nil, _err
		}
	} else if !_info.IsDir() {
		return nil, nil
	}

	// is there an info/exclude file within this directory?
	_file := filepath.Join(_gitdir, "info", "exclude")
	_, _err = os.Stat(_file)
	if _err != nil {
		if os.IsNotExist(_err) {
			return nil, nil
		} else {
			return nil, _err
		}
	}

	// attempt to load the exclude file
	return NewFromFile(_file)
} // exclude()
