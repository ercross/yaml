package parser

import (
	"github.com/ercross/yaml/test/data"
	"testing"
)

func TestScalarFrameParser(t *testing.T) {
	defaultIndentationLevel := 0
	nst := newNodeSyntaxTraverser(scalarNodeSyntax().head)
	for _, scalarNodeTokens := range testdata.ExpectedScalarTokens {
		frame := newScalarFrame(defaultIndentationLevel, nst)
		if err := frame.Build(scalarNodeTokens); err != nil {
			t.Error(err)
		}
	}
}
