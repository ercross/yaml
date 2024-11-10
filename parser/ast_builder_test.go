package parser

import (
	testdata "github.com/ercross/yaml/test/data"
	"testing"
)

func TestAstBuilder_ParseScalarTokenSequence(t *testing.T) {
	astBuilder := NewAstBuilder()
	var err error
	for _, tokens := range testdata.ScalarTokens {
		err = astBuilder.Build(tokens)
		if err != nil {
			t.Fatal(err)
		}
	}
}
