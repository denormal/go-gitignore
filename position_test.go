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
	"testing"

	"github.com/ianlewis/go-gitignore"
)

func TestPosition(t *testing.T) {
	// test the conversion of Positions to strings
	for _, _p := range _POSITIONS {
		_position := gitignore.Position{
			File:   _p.File,
			Line:   _p.Line,
			Column: _p.Column,
			Offset: _p.Offset,
		}

		// ensure the string representation of the Position is as expected
		_rtn := _position.String()
		if _rtn != _p.String {
			t.Errorf(
				"position mismatch; expected %q, got %q",
				_p.String, _rtn,
			)
		}
	}
} // TestPosition()
