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

type Error interface {
	error

	// Position returns the position of the error within the .gitignore file
	// (if any)
	Position() Position

	// Underlying returns the underlying error, permitting direct comparison
	// against the wrapped error.
	Underlying() error
}

type err struct {
	error
	_position Position
} // err()

// NewError returns a new Error instance for the given error e and position p.
func NewError(e error, p Position) Error {
	return &err{error: e, _position: p}
} // NewError()

func (e *err) Position() Position { return e._position }

func (e *err) Underlying() error { return e.error }

// ensure err satisfies the Error interface
var _ Error = &err{}
