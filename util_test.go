package gitignore_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/denormal/go-gitignore"
)

func create(content string) (*os.File, error) {
	// create a temporary file
	_file, _err := ioutil.TempFile("", "gitignore")
	if _err != nil {
		return nil, _err
	}

	// populate this file with the example .gitignore
	_, _err = io.WriteString(_file, content)
	if _err != nil {
		defer os.Remove(_file.Name())
		return nil, _err
	}
	_, _err = _file.Seek(0, io.SeekStart)
	if _err != nil {
		defer os.Remove(_file.Name())
		return nil, _err
	}

	// we have a temporary file containing the .gitignore
	return _file, nil
} // create()

func coincident(a, b gitignore.Position) bool {
	return a.Line == b.Line && a.Column == b.Column && a.Offset == b.Offset
} // coincident()

func position(p gitignore.Position) string {
	return fmt.Sprintf("%d:%d [%d]", p.Line, p.Column, p.Offset)
} // position()
