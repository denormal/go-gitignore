package gitignore

import (
	"bufio"
	"io"
)

//
// inspired by https://blog.gopheracademy.com/advent-2014/parsers-lexers/
//

// lexer is the implementation of the .gitignore lexical analyser
type lexer struct {
	_r        *bufio.Reader
	_unread   []rune
	_offset   int
	_line     int
	_column   int
	_previous int
} // lexer{}

// Lexer is the interface to the lexical analyser for .gitignore files
type Lexer interface {
	Next() (*Token, Error)

	Position() Position
	String() string
} // Lexer{}

// NewLexer returns a Lexer instance for the io.Reader r.
func NewLexer(r io.Reader) Lexer {
	return &lexer{_r: bufio.NewReader(r), _line: 1, _column: 1}
} // NewLexer()

// Next returns the next Token from the Lexer reader. If an error is
// encountered, it will be returned as an Error instance, detailing the error
// and its position within the stream.
func (l *lexer) Next() (*Token, Error) {
	// read the next rune
	_r, _err := l.read()
	if _err != nil {
		return nil, _err
	}

	switch _r {
	// end of file
	case _EOF:
		return l.token(EOF, nil), nil

	// whitespace ' ', '\t'
	case _SPACE:
		fallthrough
	case _TAB:
		l.unread(_r)
		_rtn, _err := l.whitespace()
		return l.token(WHITESPACE, _rtn), _err

	// end of line '\n' or '\r\n'
	case _CR:
		fallthrough
	case _NEWLINE:
		l.unread(_r)
		_rtn, _err := l.eol()
		return l.token(EOL, _rtn), _err

	// comment '#'
	case _COMMENT:
		l.unread(_r)
		_rtn, _err := l.eol()
		return l.token(COMMENT, _rtn), _err

	// separator '/'
	case _SEPARATOR:
		return l.token(SEPARATOR, []rune{_r}), nil

	// negation '!'
	case _NEGATION:
		return l.token(NEGATION, []rune{_r}), nil

	// any '**'
	case _WILDCARD:
		// is the wildcard followed by another wildcard?
		//      - does this represent the "any" token (i.e. "**")
		_next, _err := l.read()
		if _err != nil {
			return nil, _err
		} else if _next == _WILDCARD {
			return l.token(ANY, []rune{_WILDCARD, _WILDCARD}), nil
		}

		// otherwise we have a single wildcard, so we treat it as
		// part of a pattern
		l.unread(_next)
		fallthrough

	// pattern
	default:
		l.unread(_r)
		_rtn, _err := l.pattern()
		// if we have an empty pattern, then skip to the next token
		//      - this can happen when we encounter a line continuation
		//        at the end of a pattern, followed immediately by a
		//        blank line
		if len(_rtn) == 0 && _err == nil {
			return l.Next()
		} else {
			return l.token(PATTERN, _rtn), _err
		}
	}
} // Next()

// Position returns the current position of the Lexer.
func (l *lexer) Position() Position {
	return NewPosition(l._line, l._column, l._offset)
} // Position()

// String returns the string representation of the current position of the
// Lexer.
func (l *lexer) String() string {
	return l.Position().String()
} // String()

//
// private methods
//

// read the next rune from the stream. Return an Error if there is a problem
// reading from the stream. If the end of stream is reached, return the EOF
// Token.
func (l *lexer) read() (rune, Error) {
	var _r rune
	var _err error

	// do we have any unread runes to read?
	_length := len(l._unread)
	if _length > 0 {
		_r = l._unread[_length-1]
		l._unread = l._unread[:_length-1]

		// otherwise, attempt to read a new rune
	} else {
		_r, _, _err = l._r.ReadRune()
		if _err == io.EOF {
			return _EOF, nil
		}
	}

	// increment the offset and column counts
	l._offset++
	l._column++

	// return the rune
	return _r, l.err(_err)
} // read()

// unread returns the given runes to the stream, making them eligible to be
// read again. The runes are returned in the order given, so the last rune
// specified will be the next rune read from the stream.
func (l *lexer) unread(r ...rune) {
	// do we have an runes to "unread"
	_length := len(r)
	if _length == 0 {
		return
	}

	// initialise the unread rune list if necessary
	if l._unread == nil {
		l._unread = make([]rune, 0)
	}
	l._unread = append(l._unread, r...)

	// decrement the offset and column counts
	//      - we have to take care of column being 0
	//      - NOTE: this won't unwind indefinitely
	for ; _length > 0; _length-- {
		l._offset--
		if l._column == 1 {
			l._column = l._previous
			if l._line != 1 {
				l._line--
			}
		} else {
			l._column--
		}
	}
} // unread()

// peek returns the next rune in the stream without consuming it (i.e. it will
// be returned by the next call to read or peek). peek will return an error if
// there is a problem reading from the stream.
func (l *lexer) peek() (rune, Error) {
	// read the next rune
	_r, _err := l.read()
	if _err != nil {
		return _r, _err
	}

	// unread & return the rune
	l.unread(_r)
	return _r, _err
} // peek()

// newline adjusts the positional counters when an end of line is reached
func (l *lexer) newline() {
	// adjust the counters for the new line
	l._previous = l._column
	l._column = 1
	l._line++
} // newline()

// escape attempts to read an escape sequence (e.g. '\ ') form the input
// stream. An error will be returned if there is an error reading from the
// stream. escape returns just the escape rune if the following rune is either
// end of line or end of file (since .gitignore files do not support line
// continuations).
func (l *lexer) escape() ([]rune, Error) {
	// attempt to process the escape sequence
	_peek, _err := l.peek()
	if _err != nil {
		return nil, _err
	}

	// what is the next rune after the escape?
	switch _peek {
	// are we at the end of the line or file?
	//      - we return just the escape rune
	case _CR:
		fallthrough
	case _NEWLINE:
		fallthrough
	case _EOF:
		return []rune{_ESCAPE}, nil
	}

	// otherwise, return the escape and the next rune
	//      - we know read() will succeed here since we used peek() above
	l.read()
	return []rune{_ESCAPE, _peek}, nil
} // escape()

// eol returns all runes from the current position to the end of the line. An
// error is returned if there is a problem reading from the stream, or if a
// carriage return character '\r' is encountered that is not followed by a
// newline '\n'.
func (l *lexer) eol() ([]rune, Error) {
	// read the to the end of the line
	//      - we should only be called here when we encounter an end of line
	//        sequence or a comment
	_line := make([]rune, 0, 1)

	// loop until there's nothing more to do
	for {
		_next, _err := l.read()
		if _err != nil {
			return _line, _err
		}

		// read until we have a newline or we're at end of file
		switch _next {
		// end of file
		case _EOF:
			return _line, nil

		// carriage return - we expect to see a newline next
		case _CR:
			_peek, _err := l.peek()
			if _err != nil {
				return _line, _err
			} else if _peek != _NEWLINE {
				return _line, l.err(CarriageReturnError)
			}

		// newline
		case _NEWLINE:
			_line = append(_line, _next)
			return _line, nil
		}

		// otherwise, add this rune to the line
		_line = append(_line, _next)
	}
} // eol()

// whitespace returns all whitespace (i.e. ' ' and '\t') runes in a sequence,
// or an error if there is a problem reading the next runes.
func (l *lexer) whitespace() ([]rune, Error) {
	// read until we hit the first non-whitespace rune
	_ws := make([]rune, 0, 1)

	// loop until there's nothing more to do
	for {
		_next, _err := l.read()
		if _err != nil {
			return _ws, _err
		}

		// what is this next rune?
		switch _next {
		// space or tab is consumed
		case _SPACE:
			fallthrough
		case _TAB:
			break

		// non-whitespace rune
		default:
			// return the rune to the buffer and we're done
			l.unread(_next)
			return _ws, nil
		}

		// add this rune to the whitespace
		_ws = append(_ws, _next)
	}
} // whitespace()

// pattern returns all runes representing a file or path pattern, delimited
// either by unescaped whitespace, a path separator '/' or enf of file. An
// error is returned if a problem is encountered reading from the stream.
func (l *lexer) pattern() ([]rune, Error) {
	// read until we hit the first whitespace/end of line/eof rune
	_pattern := make([]rune, 0, 1)

	// loop until there's nothing more to do
	for {
		_next, _err := l.read()
		if _err != nil {
			return _pattern, _err
		}

		// what is the next rune?
		switch _next {
		// whitespace, newline, end of file, separator
		case _SPACE:
			fallthrough
		case _TAB:
			fallthrough
		case _CR:
			fallthrough
		case _NEWLINE:
			fallthrough
		case _SEPARATOR:
			fallthrough
		case _EOF:
			// return this rune to the lexer
			l.unread(_next)

			// return what we have
			return _pattern, nil

		// escape sequence - consume the next rune
		case _ESCAPE:
			_escape, _err := l.escape()
			if _err != nil {
				return _pattern, _err

				// if we have an empty escape sequence, then we're at
				// the end of the line, so we can ignore the continuation
			} else if _escape == nil {
				l.newline()
				continue
			}

			// otherwise, the escape sequence is part of the pattern
			_pattern = append(_pattern, _escape...)

		// comment character - this is an error
		case _COMMENT:
			return _pattern, l.err(IllegalRuneError)

		// any other character, we add to the pattern
		default:
			_pattern = append(_pattern, _next)
		}
	}
} // pattern()

// token returns a Token instance of the given type_ represented by word runes.
func (l *lexer) token(type_ TokenType, word []rune) *Token {
	// extract the lexer position
	//      - the column is taken from the current column position
	//        minus the length of the consumed "word"
	_word := len(word)
	_column := l._column - _word
	_offset := l._offset - _word
	position := NewPosition(l._line, _column, _offset)

	// if this is a newline or comment token, we adjust the line & column counts
	switch type_ {
	case EOL:
		fallthrough
	case COMMENT:
		l.newline()
	}

	// return the Token
	return NewToken(type_, word, position)
} // token()

// err returns an Error encapsulating the error e and the current Lexer
// position.
func (l *lexer) err(e error) Error {
	// do we have an error?
	if e == nil {
		return nil
	} else {
		return NewError(e, l.Position())
	}
} // err()

// ensure the lexer conforms to the lexer interface
var _ Lexer = &lexer{}
