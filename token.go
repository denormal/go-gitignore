package gitignore

import(
    "fmt"
)


type TokenType int


const(
    // this must be the first token type
    ILLEGAL    TokenType    = iota

    EOF
    EOL
    WHITESPACE

    COMMENT

    SEPARATOR

    NEGATION

    PATTERN

    ANY

    // this must be the last token type
    BAD
)


type Token struct {
    Type          TokenType
    Word        []rune
    Position
} // Token{}


func NewToken( type_   TokenType ,
               word  []rune      ,
               pos     Position  ) *Token {
    // ensure the type is valid
    if type_ < ILLEGAL || type_ > BAD {
        type_   = BAD
    }

    // return the token
    return &Token{ Type     : type_ ,
                   Word     : word  ,
                   Position : pos   }
} // NewToken()


func ( t *Token ) Name() string {
    switch t.Type {
        case ILLEGAL:       return "ILLEGAL"
        case EOF:           return "EOF"
        case EOL:           return "EOL"
        case WHITESPACE:    return "WHITESPACE"
        case COMMENT:       return "COMMENT"
        case SEPARATOR:     return "SEPARATOR"
        case NEGATION:      return "NEGATION"
        case PATTERN:       return "PATTERN"
        case ANY:           return "ANY"
        default:            return "BAD TOKEN"
    }
} // Name()


func ( t *Token ) Token() string {
    return string( t.Word )
} // Token()


func ( t *Token ) String() string {
    return fmt.Sprintf( "%s: %s %q\n" , t.Position.String() ,
                                        t.Name()            ,
                                        t.Token()           )
} // String()
