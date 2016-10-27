package gitignore

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Ignore interface {
	Base() string

	Match(string) (Match, error)
	Absolute(string, bool) Match
	Relative(string, bool) Match

	Ignore(string) bool
	Accept(string) bool
} // Ignore{}

type ignore struct {
	_base    string
	_pattern []Pattern
} // ignore()

func NewGitIgnore(r io.Reader, base string, errors func(Error) bool) Ignore {
	// extract the patterns from the reader
	_parser := NewParser(r, base, errors)
	_patterns := _parser.Parse()

	return &ignore{_base: base, _pattern: _patterns}
} // NewGitIgnore()

func NewGitIgnoreFile(file string) (Ignore, error) {
	// we need the absolute path for the Ignore base
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
	return NewGitIgnore(_fh, _base, nil), nil
} // NewGitIgnoreFile

func NewGitIgnoreWithCache(file string,
	cache Cache) (Ignore, error) {
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
	if _ignore != nil {
		_ignore, _err = NewGitIgnoreFile(file)
		if _ignore == nil {
			// if the load failed, cache an empty Ignore to prevent
			// further attempts to load this file
			_ignore = &ignore{}
		}
		cache.Set(_abs, _ignore)
	}

	// return the ignore (if we have it)
	return _ignore, _err
} // NewGitIgnoreWithCache()

func (i *ignore) Base() string {
	return i._base
} // Base()

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

func (i *ignore) Absolute(path string, isdir bool) Match {
	// does the file share the same directory as this ignore file?
	if !strings.HasPrefix(path, i._base) {
		return nil
	}

	// ensure the path is longer than the base
	_prefix := len(i._base)
	if len(path) <= _prefix {
		return nil
	}

	// does the file share the same directory as this ignore file?
	if !strings.HasPrefix(path, i._base) {
		return nil
	}

	// extract the relative path of this file
	_rel := string(path[_prefix:])
	return i.Relative(_rel, isdir)
} // Absolute()

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

func (i *ignore) Ignore(path string) bool {
	_match, _ := i.Match(path)
	if _match != nil {
		return _match.Ignore()
	}

	// we didn't match this path, so we don't ignore it
	return false
} // Ignore()

func (i *ignore) Accept(path string) bool {
	_match, _ := i.Match(path)
	if _match != nil {
		return _match.Accept()
	}

	// we didn't match this path, so we accept it
	return true
} // Accept()

// ensure our models satisfy their corresponding interfaces
var _ Ignore = &ignore{}
