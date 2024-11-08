package parser

import (
	"errors"
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/test"
	"github.com/ercross/yaml/token"
	"testing"
)

func TestPush(t *testing.T) {
	t.Parallel()

	t.Run("Test document node indentation level", func(t *testing.T) {
		TestDocumentNodeIndentationLevel(t)
	})

	t.Run("Test push child node on non-nestable node", func(t *testing.T) {
		TestPushChildOnNonNestableNode(t)
	})

	t.Run("Test push child node on nestable node", func(t *testing.T) {
		TestPushChildOnNestableNode(t)
	})

	t.Run("Test push parent-level or ancestor-level node on stack", func(t *testing.T) {
		TestPushParentLevelNode(t)
	})

	t.Run("Test disallow inconsistent indentation", func(t *testing.T) {
		TestPushModuloIncompatibleIndentation(t)
	})
}

func TestDocumentNodeIndentationLevel(t *testing.T) {
	m := newIndentationManager()
	test.AssertEqualInt(t, 0, m.peek().level, "Initial indentation level should be 0")
}

func TestPushChildOnNonNestableNode(t *testing.T) {
	m := newIndentationManager()

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("push child node on document node should panic")
		}
		if !errors.Is(r.(error), errChildNodeOnNonNestableNode) {
			t.Errorf("unexpected error: error is not child node on non-nestable node")
		}
	}()
	m.push(2, yaml.NodeTypeScalar)
}

func TestPushSiblingOnDocumentNode(t *testing.T) {
	m := newIndentationManager()

	m.push(0, yaml.NodeTypeSequenceBlockStyle)
	test.AssertEqualInt(t, 0, m.peek().level, "top indentation level should be 0")
}

func TestPushChildOnNestableNode(t *testing.T) {
	m := newIndentationManager()

	m.push(0, yaml.NodeTypeMappingBlockStyle)
	m.push(2, yaml.NodeTypeScalar)
	test.AssertEqualInt(t, 2, m.peek().level, "top indentation level should be 2")
	test.AssertEqualInt(t, 2, *m.indentationLevelModuloFactor, "indentation level modulo should be 2")
}

func TestPushParentLevelNode(t *testing.T) {
	m := newIndentationManager()
	m.push(0, yaml.NodeTypeSequenceBlockStyle)
	m.push(2, yaml.NodeTypeMappingBlockStyle)
	m.push(4, yaml.NodeTypeMappingBlockStyle)
	m.push(6, yaml.NodeTypeScalar)

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("push parent level on stack top-element should panic")
		}
		if !errors.Is(r.(error), errParentLevelIndentation) {
			t.Errorf("unexpected error: error is not parent level indentation")
		}
	}()
	m.push(4, yaml.NodeTypeScalar)
}

func TestPushModuloIncompatibleIndentation(t *testing.T) {
	m := newIndentationManager()
	m.push(0, yaml.NodeTypeSequenceBlockStyle)
	m.push(2, yaml.NodeTypeMappingBlockStyle)

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("push parent level on stack top-element should panic")
		}
		if !errors.Is(r.(error), errModuloFactorIncompatibleIndentation) {
			t.Errorf("unexpected error: error is not modulo factor incompatible indentation")
		}
	}()
	m.push(5, yaml.NodeTypeMappingBlockStyle)
}

func TestDetermineRelationship(t *testing.T) {
	m := newIndentationManager()
	actualRelationship, _ := m.determineRelationship(0)
	test.AssertEqualInt(t, indentationRelationSibling, actualRelationship, "should be sibling")

	actualRelationship, _ = m.determineRelationship(2)
	test.AssertEqualInt(t, indentationRelationshipUnknown, actualRelationship, "should be unknown")

	m.push(0, yaml.NodeTypeSequenceBlockStyle)
	actualRelationship, _ = m.determineRelationship(2)
	test.AssertEqualInt(t, indentationRelationshipChild, actualRelationship, "should be child")

	m.push(2, yaml.NodeTypeScalar)
	actualRelationship, _ = m.determineRelationship(0)
	test.AssertEqualInt(t, indentationRelationshipParentLevel, actualRelationship, "should be parent level")
}

func TestIndentationManager_FindIndentation(t *testing.T) {
	m := newIndentationManager()
	m.push(0, yaml.NodeTypeSequenceBlockStyle)
	tokens := []token.Token{
		{Type: token.TypeIndentation, Value: "    "}, // Indentation level 4
		{Type: token.TypeData, Value: "data"},
	}

	relationship, indentCount := m.findIndentation(tokens)
	test.AssertEqualInt(t, indentationRelationshipChild, relationship, "Initial node with no prior indentation should be a sibling")
	test.AssertEqualInt(t, 4, indentCount, "Initial indentation count should be 0")

	m.push(4, yaml.NodeTypeSequenceBlockStyle)
	tokens = []token.Token{
		{Type: token.TypeIndentation, Value: "    "}, // Same indentation level
		{Type: token.TypeData, Value: "data"},
	}

	relationship, indentCount = m.findIndentation(tokens)
	test.AssertEqualInt(t, indentationRelationSibling, relationship, "Indentation level 4 should be treated as a sibling when the stack level is also 4")
	test.AssertEqualInt(t, 4, indentCount, "Indentation count should match token length")

	tokens = []token.Token{
		{Type: token.TypeIndentation, Value: "        "}, // Indentation level 8
		{Type: token.TypeData, Value: "data"},
	}

	relationship, indentCount = m.findIndentation(tokens)
	test.AssertEqualInt(t, indentationRelationshipChild, relationship, "Indentation level 8 should be treated as a child of level 4")
	test.AssertEqualInt(t, 8, indentCount, "Indentation count should match token length")
}
