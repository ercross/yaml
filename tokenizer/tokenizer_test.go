package tokenizer

import (
	"testing"
)

func TestTokenizeLine(t *testing.T) {
	tkn := NewTokenizer()
	var (
		tokens []Token
		err    error
	)

	for i, line := range scalarLines {
		tokens, err = tkn.Tokenize(line, i+1)
		if err != nil {
			t.Errorf("tokenizer error: %v", err)
		}

		if len(tokens) != len(expectedScalarTokens[i]) {
			t.Errorf("expected token length %d but got %d on line %d", len(expectedScalarTokens[i]), len(tokens), i+1)
			continue
		}

		for j, token := range tokens {
			if expectedScalarTokens[i][j] != token {
				t.Errorf("Token at positon %d: expected %s, got %s", token.Position, expectedScalarTokens[i][j], token)
			}
		}
	}
}
