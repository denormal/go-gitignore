package gitignore

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// use an empty GitIgnore for cached lookups
var empty = &ignore{}

// GitIgnore is the interface to .gitignore files. It provides methods for
// testing files for matching the .gitignore file, and then determining
// whether a file should be ignored or included.
type GitIgnore interface {
	Base() string

	Match(string) (Match, error)
	Absolute(string, bool) Match
	Relative(string, bool) Match

	Ignore(string) bool
	Include(string) bool
} // GitIgnore{}

// ignore is the implementation of the .gitignore file, containing the base
// path f the file, as well as the ordered list of patterns contained in the
// file
type ignore struct {
	_base    string
	_pattern []Pattern
} // ignore()

// NewGitIgnore creates a new GitIgnore instance from the patterns listed in t,
// representing a .gitignore file in the base directory. If errors is given, it
// will be invoked for every error encountered when parsing the .gitignore
// patterns. Parsing will terminate if errors is called and returns false,
// otherwise, parsing will continue until end of file has been reached.
func New(r io.Reader, base string, errors func(Error) bool) GitIgnore {
	// extract the patterns from the reader
	_parser := NewParser(r, errors)
	_patterns := _parser.Parse()

	return &ignore{_base: base, _pattern: _patterns}
} // New()

// NewGitIgnoreFile creates a GitIgnore instance from the given file. An error
// will be returned if file cannot be opened or its absolute path determined.
func NewFromFile(file string) (GitIgnore, error) {
	// we need the absolute path for the GitIgnore base
	_file, _err := filepath.Abs(file)
	if _err != nil {
		return nil, _err
	}
	_base := filepath.Dir(_file)

	// attempt to open the ignore file to create the io.Reader
	_fh, _err := os.Open(_file)
	if _err != nil {
		return nil, _err
	}
	return New(_fh, _base, nil), nil
} // NewFromFile()

// NewGitIgnoreCached returns a GitIgnore instance (using NewGetIgnoreFile)
// for the given file. If the file has been loaded before, its GitIgnore
// instance will be returned from the cache rather than being reloaded. If
// cache is not defined, NewGitIgnoreCached will use a default cache.
//
// If NewGitIgnoreFile returns an error, NewGitIgnoreCached will store an empty
// GitIgnore (i.e. no patterns) against the file to prevent repeated parse
// attempts on subsequent requests for the same file. Subsequent calls to
// NewGitIgnoreCached for a file that could not be loaded due to an error will
// return nil.
func NewWithCache(file string, cache Cache) (GitIgnore, error) {
	// if we haven't been given a cache, use the default cache
	if cache == nil {
		cache = global
	}

	// use the file absolute path as its key into the cache
	_abs, _err := filepath.Abs(file)
	if _err != nil {
		return nil, _err
	}

	_ignore := cache.Get(_abs)
	if _ignore == nil {
		_ignore, _err = NewFromFile(file)
		if _ignore == nil {
			// if the load failed, cache an empty GitIgnore to prevent
			// further attempts to load this file
			_ignore = empty
		}
		cache.Set(_abs, _ignore)
	}

	// return the ignore (if we have it)
	if _ignore == empty {
		return nil, _err
	} else {
		return _ignore, _err
	}
} // NewWithCache()

// Base returns the directory containing the .gitignore file for this GitIgnore.
func (i *ignore) Base() string {
	return i._base
} // Base()

// Match attempts to match the path against this GitIgnore. If the path is
// matched by a GitIgnore pattern, its Match will be returned. Match will
// return an error if its not possible to determine the absolute path of the
// given path, or if its not possible to determine if the path represents a
// file or a directory.
func (i *ignore) Match(path string) (Match, error) {
	// ensure we have the absolute path for the given file
	_path, _err := filepath.Abs(path)
	if _err != nil {
		return nil, _err
	}

	// is the path a file or a directory?
	_info, _err := os.Stat(_path)
	if _err != nil {
		return nil, _err
	}
	_isdir := _info.IsDir()

	// attempt to match the absolute path
	return i.Absolute(_path, _isdir), nil
} // Match()

// Absolute attempts to match an absolute path against this GitIgnore. If the
// path is not located under the base directory of this GitIgnore, or is not
// matched by this GitIgnore, nil is returned.
func (i *ignore) Absolute(path string, isdir bool) Match {
	// does the file share the same directory as this ignore file?
	if !strings.HasPrefix(path, i._base) {
		return nil
	}

	// extract the relative path of this file
	_prefix := len(i._base) + 1
	_rel := string(path[_prefix:])
	return i.Relative(_rel, isdir)
} // Absolute()

// Relative attempts to match a path relative to the GitIgnore base directory.
// If the path is not matched by the GitIgnore, nil is returned.
func (i *ignore) Relative(path string, isdir bool) Match {
	// if we are on Windows, then translate the path to Unix form
	_rel := path
	if runtime.GOOS == "windows" {
		_rel = filepath.ToSlash(_rel)
	}

	// iterate over the patterns for this ignore file
	//      - iterate in reverse, since later patterns overwrite earlier
	for _i := len(i._pattern) - 1; _i >= 0; _i-- {
		_pattern := i._pattern[_i]
		if _pattern.Match(_rel, isdir) {
			return _pattern
		}
	}

	// we don't match this file
	return nil
} // Relative()

// Ignore returns true if the path is ignored by this GitIgnore. Paths that are
// not matched by this GitIgnore are not ignored.
func (i *ignore) Ignore(path string) bool {
	_match, _ := i.Match(path)
	if _match != nil {
		return _match.Ignore()
	}

	// we didn't match this path, so we don't ignore it
	return false
} // Ignore()

// Include returns true if the path is included by this GitIgnore. Paths that
// are not matched by this GitIgnore are always included.
func (i *ignore) Include(path string) bool {
	_match, _ := i.Match(path)
	if _match != nil {
		return _match.Include()
	}

	// we didn't match this path, so we include it
	return true
} // Include()

// ensure ignore satisfies the GitIgnore interface
var _ GitIgnore = &ignore{}
