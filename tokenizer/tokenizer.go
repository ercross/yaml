package tokenizer

import (
	"fmt"
	"github.com/ercross/yaml/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	complexTokenBuilder *complexTokenBuilder

	// indentationCharacter must be consistent once set
	indentationCharacter rune
}

// manages build process for complex tokens (e.g., quoted strings)
type complexTokenBuilder struct {
	startLine                     int
	startColumn                   int
	builder                       strings.Builder
	endBuildOnNext                rune
	startCharacterOccurrenceCount int
}

var symbolToTokenType map[rune]token.Type = map[rune]token.Type{
	token.CharNewline:              token.TypeNewline,
	token.CharColon:                token.TypeColon,
	token.CharPipe:                 token.TypePipe,
	token.CharComma:                token.TypeComma,
	token.CharGreaterThan:          token.TypeGreaterThan,
	token.CharQuestionMark:         token.TypeQuestionMark,
	token.CharExclamationMark:      token.TypeExclamationMark,
	token.CharAmpersand:            token.TypeAmpersand,
	token.CharAsterisk:             token.TypeAsterisk,
	token.CharOpeningSquareBracket: token.TypeOpeningSquareBracket,
	token.CharClosingSquareBracket: token.TypeClosingSquareBracket,
	token.CharOpeningCurlyBrace:    token.TypeOpeningCurlyBrace,
	token.CharClosingCurlyBrace:    token.TypeClosingCurlyBrace,
}

func New() *Tokenizer {
	return &Tokenizer{
		complexTokenBuilder: &complexTokenBuilder{},
	}
}

func (t *complexTokenBuilder) endBuild() {
	t.builder.Reset()
	t.endBuildOnNext = 0
	t.startLine = 0
	t.startColumn = 0
}

func (t *complexTokenBuilder) startBuilding(breakOn rune, lineNumber int, column int) {
	t.endBuildOnNext = breakOn
	t.startLine = lineNumber
	t.startColumn = column
}

func (t *Tokenizer) Run(in <-chan string, out chan<- []token.Token) error {

	lineNumber := 0
	for line := range in {
		lineNumber++
		tokens, err := t.tokenize(line, lineNumber)
		if err != nil {
			return err
		}
		out <- tokens
	}
	close(out)
	return nil
}

func (t *Tokenizer) tokenize(line string, lineNumber int) (tokens []token.Token, err error) {
	if len(line) == 0 {
		return tokens, nil
	}

	column := 1
	rawLine := []byte(line)

	for len(rawLine) > 0 {

		r, runeSize := utf8.DecodeRune(rawLine)

		if column == 1 {
			if isWhiteSpaceCharacter(r) {
				if err = t.handleWhitespace(tokens, rawLine, &column, lineNumber); err != nil {
					return tokens, err
				}
				continue
			}

			if (r == token.CharPeriod || r == token.CharDash) && len(tokens) == 0 {
				return t.handleDocumentStarters(rawLine, lineNumber)
			}
		}

		if (r == token.CharDoubleQuote || r == token.CharSingleQuote) && !t.complexTokenBuilder.isEscapeSequence(r) {
			t.complexTokenBuilder.startCharacterOccurrenceCount++
			rawLine = rawLine[runeSize:]
			if t.complexTokenBuilder.isBuilding() && t.complexTokenBuilder.endBuildOnNext == r {
				if t.complexTokenBuilder.canEndBuilding(rawLine) {
					tokens = append(tokens, token.New(token.TypeData, t.complexTokenBuilder.builder.String(), t.complexTokenBuilder.startLine, t.complexTokenBuilder.startColumn))
					t.complexTokenBuilder.endBuild()
					column++
					continue
				}

				t.complexTokenBuilder.builder.WriteString(string(r))
				continue
			}

			t.complexTokenBuilder.startBuilding(r, lineNumber, column)
			continue

		}

		if t.complexTokenBuilder.isBuilding() {
			// build data Token
			t.complexTokenBuilder.builder.WriteString(string(r))
			rawLine = rawLine[runeSize:]
			column++
			continue
		}

		if isWhiteSpaceCharacter(r) {

			// strip out whitespaces
			rawLine = rawLine[runeSize:]
			column++
			continue
		}

		if r == token.CharCommentStarter {
			tokens = append(tokens, token.New(token.TypeComment, extractComment(line), lineNumber, column))
			return
		}

		// check for YAML-meaningful symbol
		if tt, ok := symbolToTokenType[r]; ok {
			rawLine = rawLine[runeSize:]
			tokens = append(tokens, token.New(tt, "", lineNumber, column))
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

			tokens = append(tokens, token.New(token.TypeData, b.String(), lineNumber, startColumn))
			continue
		}

		return nil, fmt.Errorf("unknown Token %v on %d:%d", string(r), lineNumber, column)
	}

	return tokens, nil
}

func (t complexTokenBuilder) isBuilding() bool {
	return t.startColumn > 0 && t.endBuildOnNext != 0
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

func (t complexTokenBuilder) isEscapeSequence(r rune) bool {
	if !t.isBuilding() {
		return false
	}

	return r == token.CharDoubleQuote && strings.HasSuffix(t.builder.String(), `\`)
}

func isWhiteSpaceCharacter(r rune) bool {
	return r == token.CharWhitespace || r == token.CharTab
}

func (t complexTokenBuilder) canEndBuilding(rawLine []byte) bool {

	nextCharacter, size := utf8.DecodeRune(rawLine)
	if size == utf8.RuneError {
		return false
	}

	if (nextCharacter == token.CharWhitespace || nextCharacter == token.CharNewline) && t.startCharacterOccurrenceCount%2 == 0 {
		return true
	}
	return false
}

func (t *Tokenizer) handleDocumentStarters(rawLine []byte, lineNumber int) ([]token.Token, error) {
	column := 1
	allowedLength := 3
	if len(rawLine) != allowedLength {
		return nil, fmt.Errorf("document start [---] or end [...] tokens must be alone on a separate line")
	}
	r, runeSize := utf8.DecodeRune(rawLine)

	t.complexTokenBuilder.startBuilding(r, lineNumber, column)
	for len(rawLine) > 0 {
		if r != t.complexTokenBuilder.endBuildOnNext {
			return nil, fmt.Errorf("unexpected rune %v", string(r))
		}
		t.complexTokenBuilder.builder.WriteString(string(r))
		rawLine = rawLine[runeSize:]
		r, runeSize = utf8.DecodeRune(rawLine)
		column++
	}

	tt := token.TypeDocumentEnd
	if r == token.CharDash {
		tt = token.TypeDocumentStart
	}
	tokens := []token.Token{token.New(tt, t.complexTokenBuilder.builder.String(), lineNumber, column)}
	t.complexTokenBuilder.endBuild()
	return tokens, nil
}

func (t *Tokenizer) handleWhitespace(tokens []token.Token, rawLine []byte, column *int, lineNumber int) error {
	r, runeSize := utf8.DecodeRune(rawLine)
	if t.indentationCharacter == 0 {
		t.indentationCharacter = r
	}
	if t.indentationCharacter != r {
		return fmt.Errorf("inconsistent indentation character found at %d:%d", lineNumber, *column)
	}

	// build indentation
	var b strings.Builder
	for r == token.CharWhitespace {
		rawLine = rawLine[runeSize:]
		b.WriteString(" ")
		*column++
		r, runeSize = utf8.DecodeRune(rawLine)
	}
	if b.Len() > 0 {
		tokens = append(tokens, token.New(token.TypeIndentation, b.String(), lineNumber, *column))
	}

	return nil
}
