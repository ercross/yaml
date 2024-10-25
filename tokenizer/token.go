package tokenizer

import (
	"fmt"
)

type TokenType int8

const (
	TokenTypeUnknown TokenType = iota
	// TokenTypeData could refer to a key or value of any data type
	TokenTypeData
	TokenTypeColon
	TokenTypeDocumentStart
	TokenTypeDocumentEnd

	// TokenTypeIndentation indicates TokenType is an indentation.
	// An indentation is zero or more spaces preceding a newline rune
	TokenTypeIndentation
	TokenTypeNewline
	TokenTypePipe
	TokenTypeComma
	TokenTypeGreaterThan
	TokenTypeQuestionMark
	TokenTypeExclamationMark
	TokenTypeAmpersand
	TokenTypeAsterisk
	TokenTypeComment
	TokenTypeOpeningSquareBracket
	TokenTypeClosingSquareBracket
	TokenTypeOpeningCurlyBrace
	TokenTypeClosingCurlyBrace
)

const (
	dash                 rune = '-'
	whitespace                = ' '
	tab                       = '\t'
	commentStarter            = '#'
	newline                   = '\n'
	colon                     = ':'
	pipe                      = '|'
	comma                     = ','
	greaterThan               = '>'
	questionMark              = '?'
	exclamationMark           = '!'
	ampersand                 = '&'
	asterisk                  = '*'
	singleQuote               = '\''
	doubleQuote               = '"'
	period                    = '.'
	openingSquareBracket      = '['
	closingSquareBracket      = ']'
	openingCurlyBrace         = '{'
	closingCurlyBrace         = '}'
)

var symbolToTokenType map[rune]TokenType = map[rune]TokenType{
	newline:              TokenTypeNewline,
	colon:                TokenTypeColon,
	pipe:                 TokenTypePipe,
	comma:                TokenTypeComma,
	greaterThan:          TokenTypeGreaterThan,
	questionMark:         TokenTypeQuestionMark,
	exclamationMark:      TokenTypeExclamationMark,
	ampersand:            TokenTypeAmpersand,
	asterisk:             TokenTypeAsterisk,
	openingSquareBracket: TokenTypeOpeningSquareBracket,
	closingSquareBracket: TokenTypeClosingSquareBracket,
	openingCurlyBrace:    TokenTypeOpeningCurlyBrace,
	closingCurlyBrace:    TokenTypeClosingCurlyBrace,
}

type Token struct {
	Type     TokenType
	Value    string
	Position Location
}

// Location represent the Location of the Token within the document
type Location struct {
	line   int
	column int
}

func (l Location) String() string {
	return fmt.Sprintf("line(%d): column(%d)", l.line, l.column)
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Value: %v, Position: %v}", t.Type, t.Value, t.Position)
}

func NewToken(typ TokenType, value string, line, column int) Token {
	return Token{
		Type:  typ,
		Value: value,
		Position: Location{
			line:   line,
			column: column,
		},
	}
}
