package gitignore

import (
	"path/filepath"
	"strings"

	"github.com/danwakefield/fnmatch"
)

// Pattern represents per-line patterns within a .gitignore file
type Pattern interface {
	Match
	Match(string, bool) bool
} // Pattern{}

// pattern is the base implementation of a .gitignore pattern
type pattern struct {
	_negated   bool
	_anchored  bool
	_directory bool
	_string    string
	_fnmatch   string
	_position  Position
} // pattern()

// name represents patterns matching a file or path name (i.e. the last
// component of a path)
type name struct {
	pattern
} // name{}

// path represents a pattern that contains at least one path separator within
// the pattern (i.e. not at the start or end of the pattern)
type path struct {
	pattern
	_depth int
} // path{}

// wildcard represents a pattern that contains at least one "any" token "**"
// allowing for recursive matching.
type wildcard struct {
	pattern
	_tokens []*Token
} // wildcard{}

// NewPattern returns a Pattern from the ordered set of Tokens. The tokens are
// assumed to represent a well-formed .gitignore pattern. A Pattern may be
// negated, anchored to the start of the path (relative to the base directory
// of tie containing .gitignore), or match directories only.
func NewPattern(tokens []*Token) Pattern {
	// extract the pattern position from first token
	_position := tokens[0].Position
	_string := Tokens(tokens).String()

	// is this a negated pattern?
	_negated := false
	if tokens[0].Type == NEGATION {
		_negated = true
		tokens = tokens[1:]
	}

	// is this pattern anchored to the start of the path?
	_anchored := false
	if tokens[0].Type == SEPARATOR {
		_anchored = true
		tokens = tokens[1:]
	}

	// is this pattern for directories only?
	_directory := false
	_last := len(tokens) - 1
	if tokens[_last].Type == SEPARATOR {
		_directory = true
		tokens = tokens[:_last]
	}

	// build the pattern expression
	_fnmatch := Tokens(tokens).String()
	_pattern := &pattern{
		_negated:   _negated,
		_anchored:  _anchored,
		_position:  _position,
		_directory: _directory,
		_string:    _string,
		_fnmatch:   _fnmatch,
	}
	return _pattern.compile(tokens)
} // NewPattern()

// compile generates a specific Pattern (i.e. name, path or wildcard)
// represented by the list of tokens.
func (p *pattern) compile(tokens []*Token) Pattern {
	// what tokens do we have in this pattern?
	//      - ANY token means we can match to any depth
	//      - SEPARATOR means we have path rather than file matching
	_separator := false
	for _, _token := range tokens {
		switch _token.Type {
		case ANY:
			return p.any(tokens)
		case SEPARATOR:
			_separator = true
		}
	}

	// should we perform path or name/file matching?
	if _separator {
		return p.path(tokens)
	} else {
		return p.name(tokens)
	}
} // compile()

// Ignore returns true if the pattern describes files or paths that should be
// ignored.
func (p *pattern) Ignore() bool { return !p._negated }

// Include returns true if the pattern describes files or paths that should be
// included (i.e. not ignored)
func (p *pattern) Include() bool { return p._negated }

// Position returns the position of the first token of this pattern.
func (p *pattern) Position() Position { return p._position }

// String returns the string representation of the pattern.
func (p *pattern) String() string { return p._string }

//
// name patterns
//      - designed to match trailing file/directory names only
//

// name returns a Pattern designed to match file or directory names, with no
// path elements.
func (p *pattern) name(tokens []*Token) Pattern {
	return &name{*p}
} // name()

// Match returns true if the given path matches the name pattern. If the
// pattern is meant for directories only, and the path is not a directory,
// Match will return false. The matching is performed by fnmatch(). It
// is assumed path is relative to the base path of the owning GitIgnore.
func (n *name) Match(path string, isdir bool) bool {
	// are we expecting a directory?
	if n._directory && !isdir {
		return false
	}

	// should we match the whole path, or just the last component?
	if n._anchored {
		return fnmatch.Match(n._fnmatch, path, 0)
	} else {
		_, _base := filepath.Split(path)
		return fnmatch.Match(n._fnmatch, _base, 0)
	}
} // Match()

//
// path patterns
//      - designed to match complete or partial paths (not just filenames)
//

// path returns a Pattern designed to match paths that include at least one
// path separator '/' neither at the end nor the start of the pattern.
func (p *pattern) path(tokens []*Token) Pattern {
	// how many directory components are we expecting?
	_depth := 0
	for _, _token := range tokens {
		if _token.Type == SEPARATOR {
			_depth++
		}
	}

	// return the pattern instance
	return &path{pattern: *p, _depth: _depth}
} // path()

// Match returns true if the given path matches the path pattern. If the
// pattern is meant for directories only, and the path is not a directory,
// Match will return false. The matching is performed by fnmatch()
// with flags set to FNM_PATHNAME. It is assumed path is relative to the
// base path of the owning GitIgnore.
func (p *path) Match(path string, isdir bool) bool {
	// are we expecting a directory
	if p._directory && !isdir {
		return false
	}

	// should we match the whole path?
	if p._anchored {
		return fnmatch.Match(p._fnmatch, path, fnmatch.FNM_PATHNAME)
	}

	// attempt to extract the last N elements of the path to match
	// the expected "depth" of this pattern
	_depth := p._depth
	_index := len(path) - 1
	for ; _index >= 0; _index-- {
		// this is safe to do, since the separator is a single-byte rune
		if rune(path[_index]) == _SEPARATOR {
			_depth--
			if _depth < 0 {
				break
			}
		}
	}

	// if we don't have enough elements in the given path, then we can't match
	if _depth > 0 {
		return false
	}

	// otherwise, truncate the path
	_path := path
	if _index >= 0 {
		_path = path[_index+1:]
	}

	// match against the trailing path elements
	return fnmatch.Match(p._fnmatch, _path, fnmatch.FNM_PATHNAME)
} // Match()

//
// "any" patterns
//

// any returns a Pattern designed to match paths that include at least one
// wildcard pattern '**', specifying recursive matching.
func (p *pattern) any(tokens []*Token) Pattern {
	// consider only the non-SEPARATOR tokens, as these will be matched
	// against the path components
	_tokens := make([]*Token, 0)
	for _, _token := range tokens {
		if _token.Type != SEPARATOR {
			_tokens = append(_tokens, _token)
		}
	}

	// if the pattern is not anchored at the start, but does not start with a
	// wildcard token, then add a wildcard to the start of the set of tokens
	//
	// this simplifies the matching, since we can treat /fu/bar as **/fu/bar
	if !p._anchored {
		if tokens[0].Type != ANY {
			_any := NewToken(ANY, nil, Position{})
			_tokens = append([]*Token{_any}, _tokens...)
		}
	}

	// store the tokens
	return &wildcard{*p, _tokens}
} // any()

// Match returns true if the given path matches the wildcard pattern. If the
// pattern is meant for directories only, and the path is not a directory,
// Match will return false. The matching is performed by recursively applying
// fnmatch() with flags set to FNM_PATHNAME. It is assumed path is relative to
// the base path of the owning GitIgnore.
func (w *wildcard) Match(path string, isdir bool) bool {
	// are we expecting a directory?
	if w._directory && !isdir {
		return false
	}

	// split the path into components
	_parts := strings.Split(path, string(_SEPARATOR))

	// attempt to match the parts against the pattern tokens
	return w.match(_parts, w._tokens)
} // Match()

// match performs the recursive matching for wildcard patterns. A wildcard
// token '**' may match any path component, or no path component.
func (w *wildcard) match(path []string, tokens []*Token) bool {
	// if we have no more tokens, then we have matched this path
	// if there are also no more path elements, otherwise there's no match
	if len(tokens) == 0 {
		return len(path) == 0
	}

	// what token are we trying to match?
	_token := tokens[0]
	switch _token.Type {
	case ANY:
		// since we can match anything, whether we actually match is
		// dependent on the tokens that follow

		// do the remaining tokens match the existing path?
		if w.match(path, tokens[1:]) {
			return true

		} else if len(path) != 0 {
			// attempt to match the existing tokens against the
			// rest of the path
			return w.match(path[1:], tokens)
		}

	default:
		// if we have a non-ANY token, then we must have a non-empty path
		if len(path) != 0 {
			// if the current path element matches this token,
			// we match if the remainder of the path matches the
			// remaining tokens
			if fnmatch.Match(_token.Token(), path[0], fnmatch.FNM_PATHNAME) {
				return w.match(path[1:], tokens[1:])
			}
		}
	}

	// if we are here, then we have no match
	return false
} // match()

// ensure the patterns confirm to the Pattern interface
var _ Pattern = &name{}
var _ Pattern = &path{}
var _ Pattern = &wildcard{}
