package parser

import "github.com/ercross/yaml/token"

type indentationRelationship int8

const (
	indentationRelationshipUnknown indentationRelationship = iota
	indentationRelationshipChild
	indentationRelationshipParentLevel
	indentationRelationSibling
)

// indentationManager helps to keep track of indentation correctness in nested frames
type indentationManager struct {
	indentations []int
}

func newIndentationManager() *indentationManager {
	documentNodeIndentationLevel := 0
	return &indentationManager{[]int{documentNodeIndentationLevel}}
}

func (m *indentationManager) popStack() {
	if len(m.indentations) == 1 {
		// do not pop the last item on stack
		return
	}
	m.indentations = m.indentations[:len(m.indentations)-1]
}

func (m *indentationManager) peek() int {
	// indentationManager has been initialized with a default indentation level
	// and popStack can not remove the default indentation level,
	// so this operation is safe
	return m.indentations[len(m.indentations)-1]
}

func (m *indentationManager) pushStack(indentation int) {
	m.indentations = append(m.indentations, indentation)
}

func (m *indentationManager) determineRelationship(newIndentationLevel int) indentationRelationship {

	if m.peek() > newIndentationLevel && (m.peek()%newIndentationLevel) == 0 {
		return indentationRelationshipChild
	}

	if m.peek() == newIndentationLevel {
		return indentationRelationSibling
	}

	for i := len(m.indentations) - 1; i >= 0; i-- {
		if m.indentations[i] == newIndentationLevel {
			return indentationRelationshipParentLevel
		}
	}

	return indentationRelationshipUnknown
}

func (m *indentationManager) findIndentation(tokens []token.Token) (relationship indentationRelationship, indentationCount int) {

	if m.peek() == 0 && tokens[0].Type == token.TypeData {
		return indentationRelationSibling, 0
	}

	for i := 0; i < len(tokens); i++ {

		if tokens[i].Type == token.TypeNewline {
			continue
		}
		if tokens[i].Type == token.TypeIndentation {
			return m.determineRelationship(len(tokens[i].Value)), len(tokens[i].Value)
		}

		break
	}

	return indentationRelationshipUnknown, 0
}
