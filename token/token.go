package token

import (
	"fmt"
)

type Type int8

const (
	TypeUnknown Type = iota
	// TypeData could refer to a key or value of any data type
	TypeData
	TypeColon
	TypeDocumentStart
	TypeDocumentEnd

	// TypeIndentation indicates Type is an indentation.
	// An indentation is zero or more spaces preceding a newline rune
	TypeIndentation
	TypeNewline
	TypePipe
	TypeComma
	TypeGreaterThan
	TypeQuestionMark
	TypeExclamationMark
	TypeAmpersand
	TypeAsterisk
	TypeComment
	TypeOpeningSquareBracket
	TypeClosingSquareBracket
	TypeOpeningCurlyBrace
	TypeClosingCurlyBrace
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

var symbolToType map[rune]Type = map[rune]Type{
	newline:              TypeNewline,
	colon:                TypeColon,
	pipe:                 TypePipe,
	comma:                TypeComma,
	greaterThan:          TypeGreaterThan,
	questionMark:         TypeQuestionMark,
	exclamationMark:      TypeExclamationMark,
	ampersand:            TypeAmpersand,
	asterisk:             TypeAsterisk,
	openingSquareBracket: TypeOpeningSquareBracket,
	closingSquareBracket: TypeClosingSquareBracket,
	openingCurlyBrace:    TypeOpeningCurlyBrace,
	closingCurlyBrace:    TypeClosingCurlyBrace,
}

type Token struct {
	Type     Type
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

func NewToken(typ Type, value string, line, column int) Token {
	return Token{
		Type:  typ,
		Value: value,
		Position: Location{
			line:   line,
			column: column,
		},
	}
}
