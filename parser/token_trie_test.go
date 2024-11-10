package parser

import (
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/test/data"
	"testing"
)

func TestScalarNodeFinder(t *testing.T) {
	initTokenTrie()

	finder := newNodeTypeFinder()

	for _, sampleScalar := range data.ExpectedScalarTokens {
		finder.match(sampleScalar)
		if finder.done && finder.nodeType() != yaml.NodeTypeScalar {
			t.Errorf("expected node type %d, got %d", yaml.NodeTypeScalar, finder.nodeType())
		}
	}

}
