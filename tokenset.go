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

// tokenset represents an ordered list of Tokens
type tokenset []*Token

// String() returns a concatenated string of all runes represented by the
// list of tokens.
func (t tokenset) String() string {
	// concatenate the tokens into a single string
	_rtn := ""
	for _, _t := range []*Token(t) {
		_rtn = _rtn + _t.Token()
	}
	return _rtn
} // String()
