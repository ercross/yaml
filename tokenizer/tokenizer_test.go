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

		if len(tokens) != len(testdata.ExpectedScalarTokens[i]) {
			t.Errorf("expected token length %d but got %d on line %d", len(testdata.ExpectedScalarTokens[i]), len(tokens), i+1)
			continue
		}

		for j, tk := range tokens {
			if testdata.ExpectedScalarTokens[i][j] != tk {
				t.Errorf("Token at positon %d: expected %s, got %s", tk.Position, testdata.ExpectedScalarTokens[i][j], tk)
			}
		}
	}
}
