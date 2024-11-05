package parser

import (
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/token"
)

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

	// stack is designed to hold only exclusively nested nodes indentation,
	// one child at a time, never multiple children indentation at a time
	stack []indentation

	// indentationLevelModuloFactor is a number N such that
	// stack[i].level % N = 0 || N % stack[i].level = 0.
	//
	// It helps to manage consistency of the indentation.level, such
	// that if earlier indentations have been different by two whitespaces,
	// an indentation.level with 3 whitespace will not be allowed on the stack.
	//
	// Once indentationLevelModuloFactor is set, it should stay consistent throughout
	// current yaml.DocumentNode and should not be reset
	indentationLevelModuloFactor int
}

type indentation struct {
	level    int
	nodeType yaml.NodeType
}

func newIndentation(level int, nodeType yaml.NodeType) indentation {
	return indentation{level: level, nodeType: nodeType}
}

func newIndentationManager() *indentationManager {

	documentNodeIndentation := newIndentation(0, yaml.NodeTypeDocument)
	return &indentationManager{
		stack:                        []indentation{documentNodeIndentation},
		indentationLevelModuloFactor: documentNodeIndentation.level,
	}
}

func (m *indentationManager) pop() {
	if len(m.stack) == 1 {
		// do not pop the last item on stack
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
}

// peek returns the top indentation.level without removing it from indentationManager.stack
func (m *indentationManager) peek() int {
	// indentationManager has been initialized with a default indentation level
	// and pop can not remove the default indentation level,
	// so this operation is safe
	return m.stack[len(m.stack)-1].level
}

// push a newIndentation onto indentationManager.stack
//
// Check indentationManager.canPush for push rules
func (m *indentationManager) push(newIndentationLevel int, nodeType yaml.NodeType) {

	// set indentationManager.indentationLevelModuloFactor
	if len(m.stack) == 1 {
		m.indentationLevelModuloFactor = newIndentationLevel - m.stack[0].level
	}
	if len(m.stack) > 1 && newIndentationLevel%m.indentationLevelModuloFactor != 0 {
		panic("can not push inconsistent newIndentationLevel. " +
			"New newIndentationLevel should be a multiple or sub-multiple of indentationLevelModuloFactor")
	}

	nin := newIndentation(newIndentationLevel, nodeType)
	if !m.canPush(nin) {
		// if parser stack contains say [0, 2, 4, 6] and @newIndentationLevel = 4,
		// then there's a need to unwind the stack by 2 levels such that the yaml.Node
		// currently built in stack frame (F4) (i.e., newIndentationLevel 6) is popped
		// and added as child to F3 (i.e., frame with newIndentationLevel 4).
		// Then F3 (now containing f4) is popped and added as child to F2.
		panic("can not push a lower newIndentationLevel. unwind stack and retry push")
	}
	m.stack = append(m.stack, nin)
}

// canPush check that newIndentationLevel can be pushed onto indentationManager
func (m *indentationManager) canPush(newIndentation indentation) bool {

	// existing indentation.level can not be greater than newIndentation.level
	// because indentationManager.stack is expected to only hold exclusively nested nodes indentation,
	// one child at a time, never multiple children on the stack at a time
	if newIndentation.level < 0 || m.peek() > newIndentation.level {
		// invalid indentation.level
		return false
	}

	if len(m.stack) == 1 && newIndentation.level != m.peek() {
		// can not push a different indentation.level directly on yaml.DocumentNodeType.
		// Indentations pushed directly on the document node must be at the same indentation.level
		// with the document node
		return false
	}

	if len(m.stack) > 1 {
		if newIndentation.level%m.indentationLevelModuloFactor != 0 {
			return false
		}

		// if indentationLevel is different, then stack top indentation.nodeType must be nestable
		if newIndentation.level != m.peek() && !m.stack[len(m.stack)-1].nodeType.IsNestable() {
			return false
		}
		return true
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
		if m.stack[i].level == newIndentationLevel {
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
