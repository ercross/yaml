package yaml

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type (
	tokenType int8
	dataType  int8
)

type TokenType int

const (
	TokenTypeUnknown tokenType = iota
	// TokenTypeData could refer to a key or value of any data type
	TokenTypeData
	TokenTypeColon
	TokenTypeDash
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
	TokenTypeSingleQuote
	TokenTypeDoubleQuote
	TokenTypePeriod
	TokenTypeOpeningSquareBracket
	TokenTypeClosingSquareBracket
	TokenTypeOpeningCurlyBrace
	TokenTypeClosingCurlyBrace
)

const (
	dash                 rune = '-'
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
	dash:                 TokenTypeDash,
	newline:              TokenTypeNewline,
	colon:                TokenTypeColon,
	pipe:                 TokenTypePipe,
	comma:                TokenTypeComma,
	greaterThan:          TokenTypeGreaterThan,
	questionMark:         TokenTypeQuestionMark,
	exclamationMark:      TokenTypeExclamationMark,
	ampersand:            TokenTypeAmpersand,
	asterisk:             TokenTypeAsterisk,
	singleQuote:          TokenTypeSingleQuote,
	doubleQuote:          TokenTypeDoubleQuote,
	period:               TokenTypePeriod,
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

func newToken(typ tokenType, value string, line, column int) token {
	return token{
		Type:  typ,
		Value: value,
		position: location{
			line:   line,
			column: column,
		},
	}
}

func tokenizeLine(line string, lineNumber int) ([]token, error) {
	var tokens []token
	if len(line) == 0 {
		return tokens, nil
	}

	column := 0

	rawLine := []byte(line)
	for len(rawLine) > 0 {

		// check for indentation
		r, runeSize := utf8.DecodeRune(rawLine)
		if column == 0 && unicode.IsSpace(r) {
			var b strings.Builder

			for unicode.IsSpace(r) {
				rawLine = rawLine[runeSize:]
				b.WriteString(" ")
				column++
				r, runeSize = utf8.DecodeRune(rawLine)
			}
			if b.Len() > 0 {
				tokens = append(tokens, newToken(TokenTypeIndentation, b.String(), lineNumber, column))
			}
			continue
		}

		if r == commentStarter {
			return []token{newToken(TokenTypeComment, extractComment(line), lineNumber, column)}, nil
		}

		// check for symbol tokens
		if tt, ok := symbolToTokenType[r]; ok {
			rawLine = rawLine[runeSize:]
			tokens = append(tokens, newToken(tt, "", lineNumber, column))
			column++
			continue
		}

		// check for data
		if isData(r) {
			startColumn := column
			var b strings.Builder
			for !isYAMLValidSymbol(r) {
				b.WriteRune(r)
				rawLine = rawLine[runeSize:]
				column++
				r, runeSize = utf8.DecodeRune(rawLine)
				if r == utf8.RuneError && runeSize == 1 {
					return nil, fmt.Errorf("invalid character on line %d; column %d", lineNumber, column)
				}
			}

			tokens = append(tokens, newToken(TokenTypeData, b.String(), lineNumber, startColumn))
			continue
		}

		tokens = append(tokens, newToken(TokenTypeUnknown, string(r), lineNumber, column))
		column++
	}

	return tokens, nil
}

func extractComment(line string) string {
	parts := strings.Split("#", line)
	return parts[1]
}

func isData(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isYAMLValidSymbol(r rune) bool {
	_, ok := symbolToTokenType[r]
	return ok
}
