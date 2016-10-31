package gitignore

import (
	"fmt"
)

type TokenType int

const (
	// this must be the first token type
	ILLEGAL TokenType = iota

	EOF
	EOL
	WHITESPACE

	COMMENT

	SEPARATOR

	NEGATION

	PATTERN

	ANY

	// this must be the last token type
	BAD
)

// Token represents a parsed token from a .gitignore stream, encapsulating the
// token type, the runes comprising the token, and the position within the
// stream of the first rune of the token.
type Token struct {
	Type TokenType
	Word []rune
	Position
} // Token{}

// NewToken returns a Token instance of the given type_, represented by the
// word runes, at the stream position pos. If the token type is not know, the
// returned instance will have type BAD.
func NewToken(type_ TokenType, word []rune, pos Position) *Token {
	// ensure the type is valid
	if type_ < ILLEGAL || type_ > BAD {
		type_ = BAD
	}

	// return the token
	return &Token{Type: type_, Word: word, Position: pos}
} // NewToken()

// Name returns a string representation of the Token type.
func (t *Token) Name() string {
	switch t.Type {
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
