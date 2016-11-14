package gitignore_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/denormal/go-gitignore"
)

type repositorytest struct {
	directory string
	cache     gitignore.Cache
	cached    bool
	instance  func(string) (gitignore.GitIgnore, error)
} // repostorytest{}

func TestRepository(t *testing.T) {
	_test := &repositorytest{}
	_test.instance = func(path string) (gitignore.GitIgnore, error) {
		return gitignore.NewRepository(path, "")
	}

	// perform the repository tests
	repository(t, _test)

	// remove the temporary directory used for this test
	defer os.RemoveAll(_test.directory)
} // TestRepository()

func TestRepositoryWithCache(t *testing.T) {
	_test := &repositorytest{}
	_test.cache = gitignore.NewCache()
	_test.cached = true
	_test.instance = func(path string) (gitignore.GitIgnore, error) {
		return gitignore.NewRepositoryWithCache(path, "", _test.cache)
	}

	// perform the repository tests
	repository(t, _test)

	// remove the temporary directory used for this test
	_err := os.RemoveAll(_test.directory)
	if _err != nil {
		t.Fatalf(
			"unable to remove temporary directory %s: %s",
			_test.directory, _err.Error(),
		)
	}

	// recreate the temporary directory
	//		- this remove & recreate gives us an empty directory for the
	//		  repository test
	//		- this lets us test the caching
	_err = os.MkdirAll(_test.directory, _GITMASK)
	if _err != nil {
		t.Fatalf(
			"unable to recreate temporary directory %s: %s",
			_test.directory, _err.Error(),
		)
	}
	defer os.RemoveAll(_test.directory)

	// repeat the repository tests
	//		- these should succeed using just the cache data
	repository(t, _test)
} // TestRepositoryWithCache()

func TestInvalidRepository(t *testing.T) {
	_test := &repositorytest{}
	_test.instance = func(path string) (gitignore.GitIgnore, error) {
		return gitignore.NewRepository(path, "")
	}

	// perform the invalid repository tests
	invalid(t, _test)
} // TestInvalidRepository()

func TestInvalidRepositoryWithCache(t *testing.T) {
	_test := &repositorytest{}
	_test.cache = gitignore.NewCache()
	_test.cached = true
	_test.instance = func(path string) (gitignore.GitIgnore, error) {
		return gitignore.NewRepositoryWithCache(path, "", _test.cache)
	}

	// perform the invalid repository tests
	invalid(t, _test)

	// repeat the tests using a default cache
	_test.cache = nil
	invalid(t, _test)
} // TestInvalidRepositoryWithCache()

func repository(t *testing.T, test *repositorytest) {
	// if the test has no configured directory, then create a new
	// directory with the required .gitignore files
	if test.directory == "" {
		// create a temporary directory populated with sample .gitignore files
		//		- first, augment the test data to include file names
		_map := make(map[string]string)
		for _k, _content := range _GITREPOSITORY {
			_map[_k+"/"+gitignore.File] = _content
		}
		_dir, _err := dir(_map)
		if _err != nil {
			t.Fatalf("unable to create temporary directory: %s", _err.Error())
		}
		test.directory = _dir
	}

	// create the repository
	_repository, _err := test.instance(test.directory)
	if _err != nil {
		t.Fatalf("unable to create repository: %s", _err.Error())
	}

	// ensure we have a non-nill repository returned
	if _repository == nil {
		t.Error("expected non-nill GitIgnore repository instance; nil found")
	}

	// ensure the base of the repository is correct
	if _repository.Base() != test.directory {
		t.Errorf(
			"repository.Base() mismatch; expected %q, got %q",
			test.directory, _repository.Base(),
		)
	}

	// perform the repository matching using absolute paths
	_cb := func(path string, isdir bool) gitignore.Match {
		_path := filepath.Join(_repository.Base(), path)
		return _repository.Absolute(_path, isdir)
	}
	for _, _test := range _REPOSITORYMATCHES {
		do(t, _cb, _test)
	}

	// repeat the tests using relative paths
	_cb = func(path string, isdir bool) gitignore.Match {
		return _repository.Relative(path, isdir)
	}
	for _, _test := range _REPOSITORYMATCHES {
		do(t, _cb, _test)
	}

	// perform absolute path tests with paths not under the same repository
	_map := make(map[string]string)
	for _, _test := range _REPOSITORYMATCHES {
		_map[_test.Path] = " "
	}
	_new, _err := dir(_map)
	if _err != nil {
		t.Fatalf("unable to create temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_new)

	// first, perform Match() tests
	for _, _test := range _REPOSITORYMATCHES {
		_path := filepath.Join(_new, _test.Local())
		_match, _err := _repository.Match(_path)
		if _err != nil {
			t.Fatalf("unexpected Match() error: %s", _err.Error())
		} else if _match != nil {
			t.Fatalf("unexpected match; expected nil, got %v", _match)
		}
	}

	// next, perform Absolute() tests
	for _, _test := range _REPOSITORYMATCHES {
		// build the absolute path
		_path := filepath.Join(_new, _test.Local())

		// we don't expect to match paths not under this repository
		_match := _repository.Absolute(_path, _test.IsDir())
		if _match != nil {
			t.Fatalf("unexpected match; expected nil, got %v", _match)
		}
	}

	// now, repeat the Match() test after having first removed the
	// temporary directory
	//		- we are testing correct handling of missing files
	_err = os.RemoveAll(_new)
	if _err != nil {
		t.Fatalf(
			"unable to remove temporary directory %s: %s",
			_new, _err.Error(),
		)
	}
	for _, _test := range _REPOSITORYMATCHES {
		_path := filepath.Join(_new, _test.Local())
		_, _err := _repository.Match(_path)
		if _err == nil {
			t.Fatalf("expected Match() error; non found for %s", _path)
		} else if !os.IsNotExist(_err) {
			t.Fatalf("unexpected Match() error: %s", _err.Error())
		}
	}

	// ensure Match() behaves as expected if the absolute path cannot
	// be determined
	//		- we do this by choosing as our working directory a path
	//		  that this process does not have permission to
	_dir, _err := dir(nil)
	if _err != nil {
		t.Fatalf("unable to create temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf("unable to retrieve working directory: %s", _err.Error())
	}
	_err = os.Chdir(_dir)
	if _err != nil {
		t.Fatalf("unable to chdir into temporary directory: %s", _err.Error())
	}
	defer os.Chdir(_cwd)

	// remove permission from the temporary directory
	_err = os.Chmod(_dir, 0)
	if _err != nil {
		t.Fatalf(
			"unable to remove temporary directory %s: %s",
			_dir, _err.Error(),
		)
	}

	// perform the repository tests
	for _, _test := range _REPOSITORYMATCHES {
		_match, _err := _repository.Match(_test.Local())
		if _err == nil {
			t.Fatalf("expected Match() error; non found for %s", _test.Path)
		} else if _match != nil {
			t.Fatalf("unexpected match; expected nil, not %v", _match)
		}
	}
} // repository()

func invalid(t *testing.T, test *repositorytest) {
	// create a temporary file to use as the repository
	_file, _err := file("")
	if _err != nil {
		t.Fatalf("unable to create temporary file: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	// test repository instance creation against a file
	_repository, _err := test.instance(_file.Name())
	if _err == nil {
		t.Errorf(
			"invalid repository error; expected %q, got nil",
			gitignore.InvalidDirectoryError.Error(),
		)
	} else if _err != gitignore.InvalidDirectoryError {
		t.Errorf(
			"invalid repository mismatch; expected %q, got %q",
			gitignore.InvalidDirectoryError.Error(), _err.Error(),
		)
	}

	// ensure no repository is returned
	if _repository != nil {
		t.Error(
			"invalid repository; expected nil, got %v",
			_repository,
		)
	}

	// now, remove the temporary file and repeat the tests
	_err = os.Remove(_file.Name())
	if _err != nil {
		t.Fatalf(
			"unable to remove temporary file %s: %s",
			_file.Name(), _err.Error(),
		)
	}

	// test repository instance creating against a missing file
	_repository, _err = test.instance(_file.Name())
	if _err == nil {
		t.Errorf(
			"invalid repository error; expected %q, got nil",
			gitignore.InvalidDirectoryError.Error(),
		)
	} else if !os.IsNotExist(_err) {
		t.Errorf(
			"invalid repository mismatch; "+
				"expected no such file or directory, got %q",
			_err.Error(),
		)
	}

	// ensure no repository is returned
	if _repository != nil {
		t.Errorf(
			"invalid repository; expected nil, got %v",
			_repository,
		)
	}

	// ensure we can't create a repository instance where the absolute path
	// of the repository cannot be determined
	//		- we do this by choosing a working directory this process does
	//		  not have access to and using a relative path
	_map := map[string]string{gitignore.File: _GITIGNORE}
	_dir, _err := dir(_map)
	if _err != nil {
		t.Fatalf("unable to create a temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	// now change the working directory
	_cwd, _err := os.Getwd()
	if _err != nil {
		t.Fatalf("unable to retrieve working directory: %s", _err.Error())
	}
	_err = os.Chdir(_dir)
	if _err != nil {
		t.Fatalf("unable to chdir into temporary directory: %s", _err.Error())
	}
	defer os.Chdir(_cwd)

	// remove permissions from the working directory
	_err = os.Chmod(_dir, 0)
	if _err != nil {
		t.Fatalf("unable remove temporary directory permissions: %s: %s",
			_dir, _err.Error(),
		)
	}

	// test repository instance creating against a relative path
	//		- the relative path exists
	_repository, _err = test.instance(gitignore.File)
	if _err == nil {
		t.Errorf("expected repository error, got nil")
	} else if os.IsNotExist(_err) {
		t.Errorf(
			"unexpected repository error; file exists, but %q returned",
			_err.Error(),
		)
	}

	// next, create a repository where we do not have read permission
	// to a .gitignore file within the repository
	//		- this should trigger a panic() when attempting a file match
	for _, _test := range _REPOSITORYMATCHES {
		_map[_test.Path] = " "
	}
	_dir, _err = dir(_map)
	if _err != nil {
		t.Fatalf("unable to create a temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	_git := filepath.Join(_dir, gitignore.File)
	_err = os.Chmod(_git, 0)
	if _err != nil {
		t.Fatalf("unable remove temporary .gitignore permissions: %s: %s",
			_git, _err.Error(),
		)
	}

	// attempt to match a path in this repository
	//		- it can be anything, so we just use the .gitignore itself
	//		- between each test we recreate the repository instance to
	//		  remove the effect of any caching
	_instance := func() gitignore.GitIgnore {
		if test.cached {
			if test.cache != nil {
				test.cache = gitignore.NewCache()
			}
		}

		// create the new repository
		_repository, _err := test.instance(_dir)
		if _err != nil {
			t.Fatalf("unable to create repository: %s", _err.Error())
		}

		// return the repository
		return _repository
	}
	for _, _match := range _REPOSITORYMATCHES {
		_test := _match
		_path := filepath.Join(_dir, _test.Local())

		// try Match() with an absolute path
		panics(t, "Match()", func() (gitignore.Match, error) {
			return _instance().Match(_path)
		})
		// try Absolute() matching
		panics(t, "Absolute()", func() (gitignore.Match, error) {
			return _instance().Absolute(_path, _test.IsDir()), nil
		})
		// try Relative() matching
		panics(t, "Relative()", func() (gitignore.Match, error) {
			return _instance().Relative(_test.Local(), _test.IsDir()), nil
		})
	}
} // invalid()

func panics(t *testing.T, tag string, action func() (gitignore.Match, error)) {
	var _wg sync.WaitGroup

	_wg.Add(1)
	go func() {
		defer func() {
			_recover := recover()
			if _recover != nil {
				_, _ok := _recover.(error)
				if _ok {
					return
				} else {
					t.Fatalf(
						"%s: panic expected with error, non-error found: %v",
						tag, _recover,
					)
				}
			} else {
				t.Fatalf("%s: panic expected, none found", tag)
			}
		}()
		defer _wg.Done()

		// attempt to match against the absolute path
		_match, _err := action()
		if _err != nil {
			t.Fatalf("%s: unexpected error from match: %s", tag, _err.Error())
		} else if _match != nil {
			t.Fatalf("%s: unexpected match: %v", tag, _match)
		}
	}()

	// wait for this test to complete
	_wg.Wait()
} // panics()
