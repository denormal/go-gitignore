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

package gitignore_test

import (
	"fmt"

	"github.com/ianlewis/go-gitignore"
)

func ExampleNewFromFile() {
	ignore, err := gitignore.NewFromFile("/my/project/.gitignore")
	if err != nil {
		panic(err)
	}

	// attempt to match an absolute path
	match := ignore.Match("/my/project/src/file.go")
	if match != nil {
		if match.Ignore() {
			fmt.Println("ignore file.go")
		}
	}

	// attempt to match a relative path
	//		- this is equivalent to the call above
	match = ignore.Relative("src/file.go", false)
	if match != nil {
		if match.Include() {
			fmt.Println("include file.go")
		}
	}
} // ExampleNewFromFile()

func ExampleNewRepository() {
	ignore, err := gitignore.NewRepository("/my/project")
	if err != nil {
		panic(err)
	}

	// attempt to match a directory in the repository
	match := ignore.Relative("src/examples", true)
	if match != nil {
		if match.Ignore() {
			fmt.Printf(
				"ignore src/examples because of pattern %q at %s",
				match, match.Position(),
			)
		}
	}

	// if we have an absolute path, or a path relative to the current
	// working directory we can use the short-hand methods
	if ignore.Include("/my/project/etc/service.conf") {
		fmt.Println("include the service configuration")
	}
} // ExampleNewRepository()
