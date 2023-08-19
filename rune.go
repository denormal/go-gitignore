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

const (
	// define the sentinel runes of the lexer
	_EOF       = rune(0)
	_CR        = rune('\r')
	_NEWLINE   = rune('\n')
	_COMMENT   = rune('#')
	_SEPARATOR = rune('/')
	_ESCAPE    = rune('\\')
	_SPACE     = rune(' ')
	_TAB       = rune('\t')
	_NEGATION  = rune('!')
	_WILDCARD  = rune('*')
)
