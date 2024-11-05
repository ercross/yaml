package parser

import "github.com/ercross/yaml/token"

type indentationRelationship int8

const (
	indentationRelationshipUnknown indentationRelationship = iota
	indentationRelationshipChild
	indentationRelationshipParentLevel
	indentationRelationSibling
)

// indentationManager manages nested yaml.Node indentation,
// ensuring that each yaml.Node nested under another
// has the correct indentation with respect to previous yaml.Node indentation levels
type indentationManager struct {
	stack                          []int
	indentationLevelMultipleFactor int
}

func newIndentationManager() *indentationManager {
	documentNodeIndentationLevel := 0
	return &indentationManager{
		stack:                          []int{documentNodeIndentationLevel},
		indentationLevelMultipleFactor: documentNodeIndentationLevel,
	}
}

func (m *indentationManager) pop() {
	if len(m.stack) == 2 {
		m.indentationLevelMultipleFactor = m.stack[0]
	}

	if len(m.stack) == 1 {
		// do not pop the last item on stack
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
}

func (m *indentationManager) peek() int {
	// indentationManager has been initialized with a default indentation level
	// and pop can not remove the default indentation level,
	// so this operation is safe
	return m.stack[len(m.stack)-1]
}

// push a newIndentation onto indentationManager.stack
//
// Check indentationManager.canPush for push rules
func (m *indentationManager) push(newIndentation int) {
	if len(m.stack) == 1 {
		m.indentationLevelMultipleFactor = newIndentation - m.stack[0]
	}
	if len(m.stack) > 1 && newIndentation%m.indentationLevelMultipleFactor != 0 {
		panic("can not push inconsistent newIndentation. " +
			"New newIndentation should be a multiple or sub-multiple of indentationLevelMultipleFactor")
	}

	if !m.canPush(newIndentation) {
		// if parser stack contains say [0, 2, 4, 6] and @newIndentation = 4,
		// then there's a need to unwind the stack by 2 levels such that the yaml.Node
		// currently built in stack frame (F4) (i.e., newIndentation 6) is popped
		// and added as child to F3 (i.e., frame with newIndentation 4).
		// Then F3 (now containing f4) is popped and added as child to F2.
		panic("can not push a lower newIndentation. unwind stack and retry push")
	}
	m.stack = append(m.stack, newIndentation)
}

// canPush check that newIndentationLevel can be pushed onto indentationManager
func (m *indentationManager) canPush(newIndentationLevel int) bool {
	if newIndentationLevel < 0 || m.peek() > newIndentationLevel {
		return false
	}

	if len(m.stack) > 1 {
		return newIndentationLevel%m.indentationLevelMultipleFactor == 0
	}

	return true
}

// determineRelationship finds the hierarchical relationship between newIndentationLevel and existing indentations.
// If indentationRelationship is indentationRelationshipParentLevel, ancestorPathLength is the distance between
// newIndentationLevel and its direct parent or parent sibling.
//
// Note that if indentationRelationship is successfully determined (i.e., not indentationRelationshipUnknown),
// that doesn't imply that indentationManager would allow pushing @newIndentationLevel.
func (m *indentationManager) determineRelationship(newIndentationLevel int) (relationship indentationRelationship, ancestorPathLength int) {

	if m.peek() < newIndentationLevel {
		return indentationRelationshipChild, -1
	}

	if m.peek() == newIndentationLevel {
		return indentationRelationSibling, -1
	}

	depth := 0
	for i := len(m.stack) - 1; i >= 0; i-- {
		if m.stack[i] == newIndentationLevel {
			return indentationRelationshipParentLevel, depth
		}
		depth++
	}

	return indentationRelationshipUnknown, -1
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
			relationship, _ = m.determineRelationship(len(tokens[i].Value))
			return relationship, len(tokens[i].Value)
		}

		break
	}

	return indentationRelationshipUnknown, 0
}
