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

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	EOL
	WHITESPACE
	COMMENT
	SEPARATOR
	NEGATION
	PATTERN
	ANY
	BAD
)

// String returns a string representation of the Token type.
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case EOL:
		return "EOL"
	case WHITESPACE:
		return "WHITESPACE"
	case COMMENT:
		return "COMMENT"
	case SEPARATOR:
		return "SEPARATOR"
	case NEGATION:
		return "NEGATION"
	case PATTERN:
		return "PATTERN"
	case ANY:
		return "ANY"
	default:
		return "BAD TOKEN"
	}
} // String()
