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

// Match represents the interface of successful matches against a .gitignore
// pattern set. A Match can be queried to determine whether the matched path
// should be ignored or included (i.e. was the path matched by a negated
// pattern), and to extract the position of the pattern within the .gitignore,
// and a string representation of the pattern.
type Match interface {
	// Ignore returns true if the match pattern describes files or paths that
	// should be ignored.
	Ignore() bool

	// Include returns true if the match pattern describes files or paths that
	// should be included.
	Include() bool

	// String returns a string representation of the matched pattern.
	String() string

	// Position returns the position in the .gitignore file at which the
	// matching pattern was defined.
	Position() Position
}
