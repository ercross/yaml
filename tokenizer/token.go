package tokenizer

import (
	"fmt"
)

type tokenType int8

const (
	TokenTypeUnknown tokenType = iota
	// TokenTypeData could refer to a key or value of any data type
	TokenTypeData
	TokenTypeColon
	TokenTypeDocumentStart
	TokenTypeDocumentEnd

	// TokenTypeIndentation indicates tokenType is an indentation.
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

var symbolToTokenType map[rune]tokenType = map[rune]tokenType{
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

type token struct {
	Type     tokenType
	Value    string
	position location
}

// location represent the location of the token within the document
type location struct {
	line   int
	column int
}

func (l location) String() string {
	return fmt.Sprintf("line(%d): column(%d)", l.line, l.column)
}

func (t token) String() string {
	return fmt.Sprintf("token{Type: %v, Value: %v, Position: %v}", t.Type, t.Value, t.position)
}
