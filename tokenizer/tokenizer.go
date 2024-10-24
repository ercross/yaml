package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenizer struct {
	tokens                []token
	startLine             int
	startColumn           int
	multilineTokenBuilder strings.Builder
	endBuildOnNext        rune

	// indentationCharacter must be consistent once set
	indentationCharacter rune
}

func newTokenizer() *tokenizer {
	return &tokenizer{}
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

func (t *tokenizer) endBuild() {
	t.multilineTokenBuilder.Reset()
	t.endBuildOnNext = 0
	t.startLine = 0
	t.startColumn = 0
}

func (t *tokenizer) startBuilding(breakOn rune, lineNumber int, column int) {
	t.endBuildOnNext = breakOn
	t.startLine = lineNumber
	t.startColumn = column
}

func (t *tokenizer) tokenizeLine(line string, lineNumber int) {
	if len(line) == 0 {
		return
	}

	column := 0

	rawLine := []byte(line)
	for len(rawLine) > 0 {

		r, runeSize := utf8.DecodeRune(rawLine)

		if column == 0 {
			if isWhiteSpaceCharacter(r) {
				if t.indentationCharacter == 0 {
					t.indentationCharacter = r
				}
				if t.indentationCharacter != r {
					panic(fmt.Sprintf("inconsistent indentation character found at %d:%d", lineNumber, column))
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
					t.tokens = append(t.tokens, newToken(TokenTypeIndentation, b.String(), lineNumber, column))
				}
				continue
			}

			if r == period || r == dash {
				allowedLength := 3
				if len(rawLine) != allowedLength {
					panic(fmt.Sprintf("document start or end tokens must be on a separate line"))
				}

				t.startBuilding(r, lineNumber, column)
				for len(rawLine) > 0 {
					if r != t.endBuildOnNext {
						panic(fmt.Sprintf("unexpected rune %v", string(r)))
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
				t.tokens = append(t.tokens, newToken(tt, t.multilineTokenBuilder.String(), lineNumber, column))
				t.endBuild()
				return
			}
		}

		if (r == doubleQuote || r == singleQuote) && !t.isEscapeSequence(r) {
			rawLine = rawLine[runeSize:]
			if t.isBuilding() && t.endBuildOnNext == r {
				t.tokens = append(t.tokens, newToken(TokenTypeData, t.multilineTokenBuilder.String(), t.startLine, t.startColumn))
				t.endBuild()
				column++
				continue
			}
			if !t.isBuilding() && t.endBuildOnNext != r {
				t.startBuilding(r, lineNumber, column)
				continue
			}
			panic("unknown token build state")
		}

		if t.isBuilding() {
			// build data token
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
			t.tokens = append(t.tokens, newToken(TokenTypeComment, extractComment(line), lineNumber, column))
			return
		}

		// check for YAML-meaningful symbol
		if tt, ok := symbolToTokenType[r]; ok {
			rawLine = rawLine[runeSize:]
			t.tokens = append(t.tokens, newToken(tt, "", lineNumber, column))
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

			t.tokens = append(t.tokens, newToken(TokenTypeData, b.String(), lineNumber, startColumn))
			continue
		}

		panic(fmt.Sprintf("unknown token %v on %d:%d", string(r), lineNumber, column))
	}
}

func (t *tokenizer) isBuilding() bool {
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

func (t *tokenizer) isEscapeSequence(r rune) bool {
	if !t.isBuilding() {
		return false
	}

	return r == doubleQuote && strings.HasSuffix(t.multilineTokenBuilder.String(), "\\")
}

func isWhiteSpaceCharacter(r rune) bool {
	return r == whitespace || r == tab
}
