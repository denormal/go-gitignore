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
	"fmt"
)

// Token represents a parsed token from a .gitignore stream, encapsulating the
// token type, the runes comprising the token, and the position within the
// stream of the first rune of the token.
type Token struct {
	Type TokenType
	Word []rune
	Position
}

// NewToken returns a Token instance of the given t, represented by the
// word runes, at the stream position pos. If the token type is not know, the
// returned instance will have type BAD.
func NewToken(t TokenType, word []rune, pos Position) *Token {
	// ensure the type is valid
	if t < ILLEGAL || t > BAD {
		t = BAD
	}

	// return the token
	return &Token{Type: t, Word: word, Position: pos}
} // NewToken()

// Name returns a string representation of the Token type.
func (t *Token) Name() string {
	return t.Type.String()
} // Name()

// Token returns the string representation of the Token word.
func (t *Token) Token() string {
	return string(t.Word)
} // Token()

// String returns a string representation of the Token, encapsulating its
// position in the input stream, its name (i.e. type), and its runes.
func (t *Token) String() string {
	return fmt.Sprintf("%s: %s %q", t.Position.String(), t.Name(), t.Token())
} // String()
