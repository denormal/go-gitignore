package gitignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/denormal/go-gitignore"
)

func TestRepository(t *testing.T) {
	// create a temporary directory populated with sample .gitignore files
	//		- first, augment the test data to include file names
	_map := make(map[string]string)
	for _k, _content := range _GITPROJECT {
		_name := _k + "/" + gitignore.File
		_map[_name] = _content
	}
	_dir, _err := dir(_map)
	if _err != nil {
		t.Fatalf("unable to create temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	// create the project
	_project, _err := gitignore.NewRepository(_dir, "")
	if _err != nil {
		t.Fatalf("unable to create gitignore project: %s", _err.Error())
	}

	// ensure we have a non-nill project returned
	if _project == nil {
		t.Error("expected non-nill GitIgnore project instance; nil found")
	}

	// ensure the base of the project is correct
	if _project.Base() != _dir {
		t.Errorf(
			"project.Base() mismatch; expected %q, got %q",
			_dir, _project.Base(),
		)
	}

	// perform the project matching using absolute paths
	for _, _test := range _PROJECTMATCHES {
		match(t, _project, _project.Base(), _test)
	}

	// repeat the tests using relative paths
	for _, _test := range _PROJECTMATCHES {
		match(t, _project, "", _test)
	}
} // TestRepository()

func TestRepositoryWithCache(t *testing.T) {
	// create a temporary directory for this test
	_dir, _err := dir(nil)
	if _err != nil {
		t.Fatalf("unable to create temporary directory: %s", _err.Error())
	}
	defer os.RemoveAll(_dir)

	// create the project cache from the test data
	_cache := gitignore.NewCache()
	for _path, _content := range _GITPROJECT {
		_buffer, _err := buffer(_content)
		if _err != nil {
			t.Fatalf("unable to create io.Reader buffer: %s", _err.Error())
		}

		// create the GitIgnore instance
		_ignore := gitignore.New(_buffer, _dir, nil)

		// store the GitIgnore against the absolute path
		//		- TODO: handle Windows paths
		_abs := filepath.Join(_dir, _path, gitignore.File)
		_cache.Set(_abs, _ignore)
	}

	// create the git project
	_project, _err := gitignore.NewRepositoryWithCache(_dir, "", _cache)
	if _err != nil {
		t.Fatalf("unable to create cached project: %s", _err.Error())
	}

	// ensure we have a non-nill project returned
	if _project == nil {
		t.Error("expected non-nill GitIgnore project instance; nil found")
	}

	// ensure the base of the project is correct
	if _project.Base() != _dir {
		t.Errorf(
			"project.Base() mismatch; expected %q, got %q",
			_GITBASE, _project.Base(),
		)
	}

	// perform the project matching using absolute paths
	for _, _test := range _PROJECTMATCHES {
		match(t, _project, _project.Base(), _test)
	}

	// repeat the tests using relative paths
	for _, _test := range _PROJECTMATCHES {
		match(t, _project, "", _test)
	}
} // TestRepositoryWithCache()
