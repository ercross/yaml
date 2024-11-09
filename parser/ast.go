package parser

import "github.com/ercross/yaml"

type AbstractSyntaxTree struct {
	documents []*yaml.DocumentNode
}

func newAbstractSyntaxTree() *AbstractSyntaxTree {
	return &AbstractSyntaxTree{
		documents: []*yaml.DocumentNode{
			yaml.NewDocumentNode(),
		},
	}
}

// addChild adds n to the yaml document currently being built
func (ast *AbstractSyntaxTree) addChild(n yaml.Node) {
	currentlyParsedDocumentIndex := len(ast.documents) - 1
	ast.documents[currentlyParsedDocumentIndex].AddChild(n)
}

func (ast *AbstractSyntaxTree) startAnotherDocument() {
	ast.documents = append(ast.documents, yaml.NewDocumentNode())
}
