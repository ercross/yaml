package tokenizer

import (
	reader "github.com/ercross/yaml/file_reader"
	testdata "github.com/ercross/yaml/test/data"
	"github.com/ercross/yaml/token"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"testing"
)

func TestTokenizer(t *testing.T) {

	t.Run("Test tokenizer on scalar nodes file", func(t *testing.T) {
		pwd, _ := os.Getwd()

		scalarYamlFile := filepath.Dir(pwd) + "/test/data/scalars.yaml"
		runTokenizerTest(t, scalarYamlFile, testdata.ScalarLineTokens)
	})
}

func runTokenizerTest(t *testing.T, filename string, expectedOutput [][]token.Token) {
	lineChan := make(chan string, 10)
	tokenChan := make(chan []token.Token, 100)
	var producedTokenLines [][]token.Token

	errGrp := new(errgroup.Group)

	errGrp.Go(func() error {
		return reader.ReadFromYamlFile(filename, lineChan)
	})
	errGrp.Go(func() error {
		return New().Run(lineChan, tokenChan)
	})

	errGrp.Go(func() error {
		for tokens := range tokenChan {
			producedTokenLines = append(producedTokenLines, tokens)
		}
		return nil
	})

	if err := errGrp.Wait(); err != nil {
		t.Errorf("error-group received error: %v", err)
	}

	compareActualToExpected(t, producedTokenLines, expectedOutput)
}

func compareActualToExpected(t *testing.T, actualOutput, expectedOutput [][]token.Token) {

	if len(actualOutput) != len(expectedOutput) {
		t.Fatalf("actual output (%d) does not match expected output (%d)",
			len(actualOutput), len(expectedOutput))

	}

	for i, line := range actualOutput {
		if len(line) != len(expectedOutput[i]) {
			t.Errorf("expected token length %d but got %d on line %d",
				len(expectedOutput[i]), len(line), i+1)
		}

		for j, tk := range line {
			if expectedOutput[i][j] != tk {
				t.Errorf("token at position %d: expected %s, got %s", tk.Position, expectedOutput[i][j], tk)
			}
		}
	}
}
