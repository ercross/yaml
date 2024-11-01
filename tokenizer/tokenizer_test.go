package tokenizer

import (
	"github.com/ercross/yaml/test/data"
	"github.com/ercross/yaml/token"
	"testing"
)

func TestTokenizeLine(t *testing.T) {
	tkn := New()
	var (
		tokens []token.Token
		err    error
	)

	for i, line := range data.ScalarLines {
		tokens, err = tkn.Tokenize(line, i+1)
		if err != nil {
			t.Errorf("tokenizer error: %v", err)
		}

		if len(tokens) != len(data.ExpectedScalarTokens[i]) {
			t.Errorf("expected token length %d but got %d on line %d", len(data.ExpectedScalarTokens[i]), len(tokens), i+1)
			continue
		}

		for j, tk := range tokens {
			if data.ExpectedScalarTokens[i][j] != tk {
				t.Errorf("Token at positon %d: expected %s, got %s", tk.Position, data.ExpectedScalarTokens[i][j], tk)
			}
		}
	}
}
