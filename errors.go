package gitignore

import(
    "errors"
)


// define the standard lexer errors
var CarriageReturnError = errors.New( "unexpected carriage return '\\r'" )
var EscapeError         = errors.New( "unexpected escape '\\'" )
var ContinuationError   = errors.New( "unexpected EOF after continuation" )
var IllegalRuneError    = errors.New( "illegal character" )


// define the standard parser errors
var EOLError            = errors.New( "unexpected end of line" )
var EOFError            = errors.New( "unexpected end of file" )
var NegationError       = errors.New( "unexpected negation '!'" )
var CommentError        = errors.New( "unexpected comment '#'" )
var InvalidPatternError = errors.New( "invalid pattern" )
