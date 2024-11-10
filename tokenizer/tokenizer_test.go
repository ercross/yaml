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

	for i, line := range testdata.ScalarLines {
		tokens, err = tkn.Tokenize(line, i+1)
		if err != nil {
			t.Errorf("tokenizer error: %v", err)
		}

		if len(tokens) != len(testdata.ScalarTokens[i]) {
			t.Errorf("expected token length %d but got %d on line %d", len(testdata.ScalarTokens[i]), len(tokens), i+1)
			continue
		}

		for j, tk := range tokens {
			if testdata.ScalarTokens[i][j] != tk {
				t.Errorf("Token at positon %d: expected %s, got %s", tk.Position, testdata.ScalarTokens[i][j], tk)
			}
		}
	}
}
