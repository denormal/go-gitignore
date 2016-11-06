package gitignore

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const File = ".gitignore"

type repository struct {
	ignore
	_cache Cache
	_file  string
} // repository{}

func NewRepository(base, file string) (GitIgnore, error) {
	return NewRepositoryWithCache(base, file, nil)
} // NewRepository()

func NewRepositoryWithCache(base, file string, cache Cache) (GitIgnore, error) {
	// if the cache is not given, then use the default global cache
	if cache == nil {
		cache = global
	}

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

// Absolute attempts to match an absolute path against this repository. If the
// path is not located under the base directory of this repository, or is not
// matched by this repository, nil is returned.
func (p *repository) Absolute(path string, isdir bool) Match {
	// does the file share the same directory as this ignore file?
	if !strings.HasPrefix(path, p.Base()) {
		return nil
	}

	// extract the relative path of this file
	_prefix := len(p.Base()) + 1
	_rel := string(path[_prefix:])
	return p.Relative(_rel, isdir)
} // Absolute()

// Relative attempts to match a path relative to the repository base directory.
// If the path is not matched by the repository, nil is returned.
func (p *repository) Relative(path string, isdir bool) Match {
	// if we are on Windows, then translate the path to Unix form
	_rel := path
	if runtime.GOOS == "windows" {
		_rel = filepath.ToSlash(_rel)
	}

	// repository matching:
	//		- a child path cannot be considered if its parent is ignored
	//		- a .gitignore in a lower directory overrides a .gitignore in a
	//		  higher directory

	// matching algorithm:
	//		- descend from the repository base to the path parent attempting
	//		  to match the descendant path (e.g. a, a/b, a/b/c, ...)
	//		- if the descendant path is ignored, then the path is ignored
	//		- otherwise, attempt to match from the path tail (e.g. the file
	//		  name), up to the repository base (i.e. the full, relative path),
	//		  and if the path is matched, return the match

	// extract the directory components of the path
	_path := strings.Split(_rel, string(_SEPARATOR))
	_length := len(_path)
	if _length > 1 {
		// we have at least one directory component, so attempt to match
		// the ancestral path
		_parent := _path[:_length-1]
		_match := p.down(p._base, _parent)

		// if the parent directory is ignored, then the path is ignored
		if _match != nil {
			if _match.Ignore() {
				return _match
			}
		}
	}

	// otherwise, we attempt to match from the file, up the parent directories
	return p.up(_path[_length-1], _path[:_length-1], isdir)
} // Relative()

func (p *repository) down(path string, remaining []string) Match {
	// if we have no remaining path elements, we cannot descend further
	if len(remaining) == 0 {
		return nil
	}

	// attempt to load the .gitignore in the parent directory
	//		- the parent is given relative to the base
	_file := filepath.Join(path, p._file)
	_ignore, _err := NewWithCache(_file, p._cache)
	if _err != nil {
		if !os.IsNotExist(_err) {
			// TODO: can we do better?
			panic(_err)
		}
	} else if _ignore != nil {
		// does the remaining path match?
		//		- we are only matching directories
		//		- we match iteratively, starting with the first remaining
		//		  component, and then adding others, to mimic traversing
		//		  down the remaining path
		for _i := 1; _i <= len(remaining); _i++ {
			_remaining := filepath.Join(remaining[:_i]...)
			_match := _ignore.Relative(_remaining, true)
			if _match != nil {
				if _match.Ignore() {
					return _match
				}
			}
		}
	}

	// descend the directory tree
	_path := filepath.Join(path, remaining[0])
	return p.down(_path, remaining[1:])
} // down()

func (p *repository) up(path string, remaining []string, isdir bool) Match {
	// attempt to load the .gitignore in the parent directory
	//		- the parent is given relative to the base
	_remaining := filepath.Join(remaining...)
	_file := filepath.Join(p._base, _remaining, p._file)
	_ignore, _err := NewWithCache(_file, p._cache)
	if _err != nil {
		if !os.IsNotExist(_err) {
			// TODO: can we do better?
			panic(_err)
		}
	} else if _ignore != nil {
		// does this path match?
		_match := _ignore.Relative(path, isdir)
		if _match != nil {
			return _match
		}
	}

	// if we have no remaining path elements, we cannot ascend
	if len(remaining) == 0 {
		return nil
	}

	// ascend the directory tree
	_last := len(remaining) - 1
	_path := filepath.Join(remaining[_last], path)
	return p.up(_path, remaining[:_last], isdir)
} // up()
