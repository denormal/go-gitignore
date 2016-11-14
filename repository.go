package gitignore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const File = ".gitignore"

type repository struct {
	ignore
	_cache Cache
	_file  string
} // repository{}

func NewRepository(base, file string) (GitIgnore, error) {
	return NewRepositoryWithCache(base, file, NewCache())
} // NewRepository()

func NewRepositoryWithCache(base, file string, cache Cache) (GitIgnore, error) {
	// extract the absolute path of the base directory
	_base, _err := filepath.Abs(base)
	if _err != nil {
		return nil, _err
	}

	// ensure the given base is a directory
	_info, _err := os.Stat(_base)
	if _err != nil {
		return nil, _err
	} else if !_info.IsDir() {
		return nil, InvalidDirectoryError
	}

	// if we haven't been given a base file name, use the default
	if file == "" {
		file = File
	}

	// return the repository instance
	_ignore := ignore{_base: _base}
	return &repository{ignore: _ignore, _cache: cache, _file: file}, nil
} // NewRepositoryWithCache()

// Match attempts to match the path against this repository. If the path is
// matched by a repository pattern, its Match will be returned. Match will
// return an error if its not possible to determine the absolute path of the
// given path, or if its not possible to determine if the path represents a
// file or a directory.
func (r *repository) Match(path string) (Match, error) {
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
	return r.Absolute(_path, _isdir), nil
} // Match()

// Absolute attempts to match an absolute path against this repository. If the
// path is not located under the base directory of this repository, or is not
// matched by this repository, nil is returned.
func (r *repository) Absolute(path string, isdir bool) Match {
	// does the file share the same directory as this ignore file?
	if !strings.HasPrefix(path, r.Base()) {
		return nil
	}

	// extract the relative path of this file
	_prefix := len(r.Base()) + 1
	_rel := string(path[_prefix:])
	return r.Relative(_rel, isdir)
} // Absolute()

// Relative attempts to match a path relative to the repository base directory.
// If the path is not matched by the repository, nil is returned.
func (r *repository) Relative(path string, isdir bool) Match {
	// if there's no path, then there's nothing to match
	_path := filepath.Clean(path)
	if _path == "." {
		return nil
	}

	// repository matching:
	//		- a child path cannot be considered if its parent is ignored
	//		- a .gitignore in a lower directory overrides a .gitignore in a
	//		  higher directory

	// first, is the parent directory ignored?
	//		- extract the parent directory from the current path
	_parent, _local := filepath.Split(_path)
	_match := r.Relative(_parent, true)
	if _match != nil {
		if _match.Ignore() {
			return _match
		}
	}

	// the parent directory isn't ignored, so we now look at the original path
	//		- we consider .gitignore files in the current directory first, then
	//		  move up the path hierarchy
	var _last string
	for {
		_file := filepath.Join(r._base, _parent, r._file)
		_ignore, _err := NewWithCache(_file, r._cache, nil)
		if _err != nil {
			if !os.IsNotExist(_err) {
				// TODO: can we do better?
				panic(errors.New(_file + ": " + _err.Error()))
			}
		} else if _ignore != nil {
			_match := _ignore.Relative(_local, isdir)
			if _match != nil {
				return _match
			}
		}

		// if there's no parent, then we're done
		//		- since we use filepath.Clean() we look for "."
		if _parent == "." {
			return nil
		}

		// we don't have a match for this file, so we progress up the
		// path hierarchy
		//		- we are manually building _local using the .gitignore
		//		  separator "/", which is how we handle operating system
		//		  file system differences
		_parent, _last = filepath.Split(_parent)
		_parent = filepath.Clean(_parent)
		_local = _last + string(_SEPARATOR) + _local
	}
} // Relative()
