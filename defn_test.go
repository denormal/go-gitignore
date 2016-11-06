package gitignore_test

import (
	"github.com/denormal/go-gitignore"
)

type test struct {
	Path    string // test path
	Pattern string // matching pattern (if any)
	Ignore  bool   // whether the path is ignored or included
} // test{}

// define the constants for the unit tests
const (
	// define the example .gitignore file contents
	_GITIGNORE = `
# example .gitignore contents

!*.go

*.o
*.a

/ignore/this/path/
/and/**/all/**/these/**
!/but/not/this\ 

we support   spaces

/**/this.is.not/a ** valid/pattern
/**/nor/is/***/this
/nor/is***this
and this is #3 failure

but \this\ is / valid\#
`

	// define the example .gitignore file contents for the Match tests
	// these tests have been taken from
	//		https://github.com/sdobz/backup/gitignore
	//
	// https://github.com/sdobz/backup/blob/master/gitignore/gitignore_test.go
	_GITMATCH = `
*.[oa]
*.html
*.min.js

!foo*.html
foo-excl.html

vmlinux*

\!important!.txt

log/*.log
!/log/foo.log

**/logdir/log
**/foodir/bar
exclude/**

!findthis*

**/hide/**
subdir/subdir2/

/rootsubdir/

dirpattern/

README.md

# arch/foo/kernel/.gitignore
!arch/foo/kernel/vmlinux*

# htmldoc/.gitignore
!htmldoc/*.html

# git-sample-3/.gitignore
git-sample-3/*
!git-sample-3/foo
git-sample-3/foo/*
!git-sample-3/foo/bar
`

	// define the number of good & bad patterns in the .gitignore above
	_GITPATTERNS    = 8
	_GITBADPATTERNS = 4

	// define the number of good & bad patterns in the match .gitignore above
	_GITMATCHPATTERNS    = 24
	_GITBADMATCHPATTERNS = 0

	// define the number of good and bad patterns returned when the
	// gitignore.Parser error handler returns false upon receiving an error
	_GITPATTERNSFALSE    = 7
	_GITBADPATTERNSFALSE = 1

	// define the base path for a git repository
	_GITBASE = "/my/git/repository"

	// define the directory mask for any directories created during testing
	_GITMASK = 0700
)

var (
	// define the positions of the bad patterns
	_GITBADPOSITION = []gitignore.Position{
		gitignore.NewPosition(15, 19, 148),
		gitignore.NewPosition(16, 14, 178),
		gitignore.NewPosition(17, 8, 192),
		gitignore.NewPosition(18, 13, 212),
	}

	// define the positions of the good patterns
	_GITPOSITION = []gitignore.Position{
		gitignore.NewPosition(4, 1, 32),
		gitignore.NewPosition(6, 1, 39),
		gitignore.NewPosition(7, 1, 43),
		gitignore.NewPosition(9, 1, 48),
		gitignore.NewPosition(10, 1, 67),
		gitignore.NewPosition(11, 1, 91),
		gitignore.NewPosition(13, 1, 109),
		gitignore.NewPosition(20, 1, 224),
	}

	// define the token stream for the _GITIGNORE .gitignore
	_GITTOKENS = []struct {
		Type  gitignore.TokenType
		Name  string
		Token string
	}{
		// 1:
		{gitignore.EOL, "EOL", "\n"},
		// 2: # example .gitignore contents
		{gitignore.COMMENT, "COMMENT", "# example .gitignore contents"},
		{gitignore.EOL, "EOL", "\n"},
		// 3:
		{gitignore.EOL, "EOL", "\n"},
		// 4: !*.go
		{gitignore.NEGATION, "NEGATION", "!"},
		{gitignore.WILDCARD, "WILDCARD", "*"},
		{gitignore.PATTERN, "PATTERN", ".go"},
		{gitignore.EOL, "EOL", "\n"},
		// 5:
		{gitignore.EOL, "EOL", "\n"},
		// 6: *.o
		{gitignore.WILDCARD, "WILDCARD", "*"},
		{gitignore.PATTERN, "PATTERN", ".o"},
		{gitignore.EOL, "EOL", "\n"},
		// 7: *.a
		{gitignore.WILDCARD, "WILDCARD", "*"},
		{gitignore.PATTERN, "PATTERN", ".a"},
		{gitignore.EOL, "EOL", "\n"},
		// 8:
		{gitignore.EOL, "EOL", "\n"},
		// 9: /ignore/this/path/
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "ignore"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "this"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "path"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.EOL, "EOL", "\n"},
		// 10: /and/**/all/**/these/**
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "and"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "all"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "these"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.EOL, "EOL", "\n"},
		// 11: !/but/not/this\
		{gitignore.NEGATION, "NEGATION", "!"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "but"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "not"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "this"},
		{gitignore.EOL, "EOL", "\n"},
		// 12:
		{gitignore.EOL, "EOL", "\n"},
		// 13: we support   spaces
		{gitignore.PATTERN, "PATTERN", "we"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "support"},
		{gitignore.WHITESPACE, "WHITESPACE", "   "},
		{gitignore.PATTERN, "PATTERN", "spaces"},
		{gitignore.EOL, "EOL", "\n"},
		// 14:
		{gitignore.EOL, "EOL", "\n"},
		// 15: /**/this.is.not/a ** valid/pattern
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "this.is.not"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "a"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "valid"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "pattern"},
		{gitignore.EOL, "EOL", "\n"},
		// 16: /**/nor/is/***/this
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "nor"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "is"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.WILDCARD, "WILDCARD", "*"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "this"},
		{gitignore.EOL, "EOL", "\n"},
		// 17: /nor/is***this
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "nor"},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.PATTERN, "PATTERN", "is"},
		{gitignore.ANY, "ANY", "**"},
		{gitignore.WILDCARD, "WILDCARD", "*"},
		{gitignore.PATTERN, "PATTERN", "this"},
		{gitignore.EOL, "EOL", "\n"},
		// 18: and this is #3 failure
		{gitignore.PATTERN, "PATTERN", "and"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "this"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "is"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.COMMENT, "COMMENT", "#3 failure"},
		{gitignore.EOL, "EOL", "\n"},
		// 19:
		{gitignore.EOL, "EOL", "\n"},
		// 20: but \this\ is / valid
		{gitignore.PATTERN, "PATTERN", "but"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "\\this\\ is"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.SEPARATOR, "SEPARATOR", "/"},
		{gitignore.WHITESPACE, "WHITESPACE", " "},
		{gitignore.PATTERN, "PATTERN", "valid"},
		{gitignore.EOL, "EOL", "\n"},

		{gitignore.EOF, "EOF", ""},
	}

	// define match tests and their expected results
	_GITMATCHES = []test{
		{"!important!.txt", "\\!important!.txt", true},
		{"arch/", "", false},
		{"arch/foo/", "", false},
		{"arch/foo/kernel/", "", false},
		{"arch/foo/kernel/vmlinux.lds.S", "!arch/foo/kernel/vmlinux*", false},
		{"arch/foo/vmlinux.lds.S", "vmlinux*", true},
		{"bar/", "", false},
		{"bar/testfile", "", false},
		{"dirpattern", "", false},
		{"Documentation/", "", false},
		{"Documentation/foo-excl.html", "foo-excl.html", true},
		{"Documentation/foo.html", "!foo*.html", false},
		{"Documentation/gitignore.html", "*.html", true},
		{"Documentation/test.a.html", "*.html", true},
		{"exclude", "exclude/**", true},
		{"exclude/dir1", "exclude/**", true},
		{"exclude/dir1/dir2", "exclude/**", true},
		{"exclude/dir1/dir2/dir3", "exclude/**", true},
		{"exclude/dir1/dir2/dir3/testfile", "exclude/**", true},
		{"file.o", "*.[oa]", true},
		{"foodir", "", false},
		{"foodir/bar/", "**/foodir/bar", true},
		{"foodir/bar/testfile", "", false},
		{"git-sample-3/", "", false},
		{"git-sample-3/foo/", "!git-sample-3/foo", false},
		{"git-sample-3/foo/bar/", "!git-sample-3/foo/bar", false},
		{"git-sample-3/foo/test/", "git-sample-3/foo/*", true},
		{"git-sample-3/test/", "git-sample-3/*", true},
		{"git-sample-3", "", false},
		{"htmldoc/", "", false},
		{"htmldoc/docs.html", "!htmldoc/*.html", false},
		{"htmldoc/jslib.min.js", "*.min.js", true},
		{"lib.a", "*.[oa]", true},
		{"log/", "", false},
		{"log/foo.log", "!/log/foo.log", false},
		{"log/test.log", "log/*.log", true},
		{"rootsubdir/", "/rootsubdir/", true},
		{"rootsubdir/foo", "", false},
		{"src/", "", false},
		{"src/findthis.o", "!findthis*", false},
		{"src/internal.o", "*.[oa]", true},
		{"subdir/", "", false},
		{"subdir/hide/", "**/hide/**", true},
		{"subdir/hide/foo", "**/hide/**", true},
		{"subdir/logdir/", "", false},
		{"subdir/logdir/log/", "**/logdir/log", true},
		{"subdir/logdir/log/findthis.log", "!findthis*", false},
		{"subdir/logdir/log/foo.log", "log/*.log", true},
		{"subdir/logdir/log/test.log", "log/*.log", true},
		{"subdir/rootsubdir/", "", false},
		{"subdir/rootsubdir/foo", "", false},
		{"subdir/subdir2/", "subdir/subdir2/", true},
		{"subdir/subdir2/bar", "", false},
		{"README.md", "README.md", true},
	}

	// define the cache tests
	_CACHETEST = map[string]gitignore.GitIgnore{
		"a":     null(),
		"a/b":   null(),
		"a/b/c": nil,
	}

	// define a set of cache keys known not to be in the cache tests above
	_CACHEUNKNOWN = []string{
		"b",
		"b/c",
	}

	// define the set of .gitignore files for a project
	_GITPROJECT = map[string]string{
		// define the top-level .gitignore file
		"": `
# ignore .bak files
*.bak
`,
		// define subdirectory .gitignore files
		"a": `
# ignore .go files
*.go

# ignore every c directory
#	- this should be the same as c/
**/c/
`,
		"a/b": `
# include .go files in this directory
!*.go
`,
	}

	// define project match tests and their expected results
	_PROJECTMATCHES = []test{
		{"include.go", "", false},
		{"ignore.go.bak", "*.bak", true},
		{"a/ignore.go", "*.go", true},
		{"a/ignore.go.bak", "*.bak", true},
		{"a/include.sh", "", false},
		{"a/c/ignore.go", "**/c/", true},
		{"a/c/ignore.go.bak", "**/c/", true},
		{"a/c/ignore.sh", "**/c/", true},
		{"a/c/", "**/c/", true},
		{"a/b/c/d/ignore.go", "**/c/", true},
		{"a/b/c/d/ignore.go.bak", "**/c/", true},
		{"a/b/c/d/ignore.sh", "**/c/", true},
		{"a/b/c/d/", "**/c/", true},
		{"a/b/c/", "**/c/", true},
		{"a/b/include.go", "!*.go", false},
		{"a/b/ignore.go.bak", "*.bak", true},
		{"a/b/include.sh", "", false},
	}
)
