package gitignore_test

import (
	"testing"

	"os"

	"github.com/denormal/go-gitignore"
)

// TestLexer tests the behaviour of gitignore.Lexer
func TestLexer(t *testing.T) {
	// create a temporary .gitignore
	_file, _err := create(_GITIGNORE)
	if _err != nil {
		t.Fatalf("unable to create temporary .gitignore: %s", _err.Error())
	}
	defer os.Remove(_file.Name())

	// ensure we have a non-nil Lexer instance
	_lexer := gitignore.NewLexer(_file)
	if _lexer == nil {
		t.Error("expected non-nil Lexer instance; nil found")
	}

	// ensure the stream of tokens is as we expect
	for _, _expected := range _GITTOKENS {
		// extract the next token from the lexer
		_got, _err := _lexer.Next()

		// ensure we did not receive an error
		if _err != nil {
			t.Errorf(
				"unable to retrieve expected token; %s at %s",
				_err.Error(), position(_err.Position()),
			)
		}

		// ensure the token is as we expect
		if _got.Type != _expected.Type {
			t.Errorf(
				"token type mismatch; expected type %d, got %d [%s]",
				_expected.Type, _got.Type, _got,
			)
		}

		// ensure the token has the correct name
		if _got.Name() != _expected.Name {
			t.Errorf(
				"token name mismatch; expected name %q, got %q [%s]",
				_expected.Name, _got.Name(), _got,
			)
		}

	}

	// ensure there are no more tokens
	_next, _err := _lexer.Next()
	if _err != nil {
		t.Errorf("unexpected error on end of token test: %s", _err.Error())
	} else if _next == nil {
		t.Errorf("unexpected nil token at end of test")
	} else if _next.Type != gitignore.EOF {
		t.Errorf(
			"token type mismatch; expected type %d, got %d [%s]",
			gitignore.EOF, _next.Type, _next,
		)
	}
} // TestLexer()
