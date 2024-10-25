package tokenizer

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Tokenizer struct {
	startLine             int
	startColumn           int
	multilineTokenBuilder strings.Builder
	endBuildOnNext        rune

	// indentationCharacter must be consistent once set
	indentationCharacter rune
}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) endBuild() {
	t.multilineTokenBuilder.Reset()
	t.endBuildOnNext = 0
	t.startLine = 0
	t.startColumn = 0
}

func (t *Tokenizer) startBuilding(breakOn rune, lineNumber int, column int) {
	t.endBuildOnNext = breakOn
	t.startLine = lineNumber
	t.startColumn = column
}

func (t *Tokenizer) Tokenize(line string, lineNumber int) (tokens []Token, err error) {
	if len(line) == 0 {
		return tokens, nil
	}

	column := 1

	rawLine := []byte(line)
	for len(rawLine) > 0 {

		r, runeSize := utf8.DecodeRune(rawLine)

		if column == 1 {
			if isWhiteSpaceCharacter(r) {
				if t.indentationCharacter == 0 {
					t.indentationCharacter = r
				}
				if t.indentationCharacter != r {
					return nil, fmt.Errorf("inconsistent indentation character found at %d:%d", lineNumber, column)
				}

				// build indentation
				var b strings.Builder
				for r == whitespace {
					rawLine = rawLine[runeSize:]
					b.WriteString(" ")
					column++
					r, runeSize = utf8.DecodeRune(rawLine)
				}
				if b.Len() > 0 {
					tokens = append(tokens, NewToken(TokenTypeIndentation, b.String(), lineNumber, column))
				}
				continue
			}

			if r == period || r == dash {
				allowedLength := 3
				if len(rawLine) != allowedLength {
					return nil, fmt.Errorf("document start or end tokens must be on a separate line")
				}

				t.startBuilding(r, lineNumber, column)
				for len(rawLine) > 0 {
					if r != t.endBuildOnNext {
						return nil, fmt.Errorf("unexpected rune %v", string(r))
					}
					t.multilineTokenBuilder.WriteString(string(r))
					rawLine = rawLine[runeSize:]
					r, runeSize = utf8.DecodeRune(rawLine)
					column++
				}

				tt := TokenTypeDocumentEnd
				if r == dash {
					tt = TokenTypeDocumentStart
				}
				tokens = append(tokens, NewToken(tt, t.multilineTokenBuilder.String(), lineNumber, column))
				t.endBuild()
				return
			}
		}

		if (r == doubleQuote || r == singleQuote) && !t.isEscapeSequence(r) {
			rawLine = rawLine[runeSize:]
			if t.isBuilding() && t.endBuildOnNext == r {
				tokens = append(tokens, NewToken(TokenTypeData, t.multilineTokenBuilder.String(), t.startLine, t.startColumn))
				t.endBuild()
				column++
				continue
			}
			if !t.isBuilding() && t.endBuildOnNext != r {
				t.startBuilding(r, lineNumber, column)
				continue
			}
			return nil, errors.New("unknown Token build state")
		}

		if t.isBuilding() {
			// build data Token
			t.multilineTokenBuilder.WriteString(string(r))
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

		if r == commentStarter {
			tokens = append(tokens, NewToken(TokenTypeComment, extractComment(line), lineNumber, column))
			return
		}

		// check for YAML-meaningful symbol
		if tt, ok := symbolToTokenType[r]; ok {
			rawLine = rawLine[runeSize:]
			tokens = append(tokens, NewToken(tt, "", lineNumber, column))
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
					panic(fmt.Errorf("invalid character on line %d; column %d", lineNumber, column))
				}
			}

			tokens = append(tokens, NewToken(TokenTypeData, b.String(), lineNumber, startColumn))
			continue
		}

		return nil, fmt.Errorf("unknown Token %v on %d:%d", string(r), lineNumber, column)
	}

	return tokens, nil
}

func (t *Tokenizer) isBuilding() bool {
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

func (t *Tokenizer) isEscapeSequence(r rune) bool {
	if !t.isBuilding() {
		return false
	}

	return r == doubleQuote && strings.HasSuffix(t.multilineTokenBuilder.String(), "\\")
}

func isWhiteSpaceCharacter(r rune) bool {
	return r == whitespace || r == tab
}
