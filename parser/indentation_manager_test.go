package parser

import (
	"github.com/ercross/yaml/test"
	"testing"

	"github.com/ercross/yaml/token"
)

func TestIndentationManager_PushPopPeek(t *testing.T) {
	m := newIndentationManager()
	test.AssertEqualInt(t, 0, m.peek(), "Initial indentation level should be 0")

	m.push(2)
	test.AssertEqualInt(t, 2, m.peek(), "Top of stack should now be 2")

	m.push(4)
	test.AssertEqualInt(t, 4, m.peek(), "Top of stack should now be 4")

	m.pop()
	test.AssertEqualInt(t, 2, m.peek(), "After pop, top of stack should be back to 2")

	m.pop()
	test.AssertEqualInt(t, 0, m.peek(), "After popping to the base, indentation level should return to 0")

	// Ensure that we can't pop beyond the initial level
	m.pop()
	test.AssertEqualInt(t, 0, m.peek(), "Should not pop the last item on stack")
}

func TestIndentationManager_DetermineRelationship(t *testing.T) {
	m := newIndentationManager()

	m.push(2)
	m.push(4)

	actualRelationship, _ := m.determineRelationship(6)
	test.AssertEqualInt(t, indentationRelationshipChild, actualRelationship, "6 should be considered a child of 4")

	actualRelationship, _ = m.determineRelationship(4)
	test.AssertEqualInt(t, indentationRelationSibling, actualRelationship, "4 should be considered a sibling of 4")

	actualRelationship, _ = m.determineRelationship(2)
	test.AssertEqualInt(t, indentationRelationshipParentLevel, actualRelationship, "2 should be considered a parent of 4")

	m.push(6)
	actualRelationship, pathLength := m.determineRelationship(2)
	test.AssertEqualInt(t, indentationRelationshipParentLevel, actualRelationship, "2 should be considered a parent of 6")
	test.AssertEqualInt(t, 2, pathLength, "path length should be 2")

	actualRelationship, _ = m.determineRelationship(3)
	test.AssertEqualInt(t, indentationRelationshipUnknown, actualRelationship, "3 should be considered an unknown relationship")
}

func TestIndentationManager_FindIndentation(t *testing.T) {
	m := newIndentationManager()

	tokens := []token.Token{
		{Type: token.TypeIndentation, Value: "    "}, // Indentation level 4
		{Type: token.TypeData, Value: "data"},
	}

	relationship, indentCount := m.findIndentation(tokens)
	test.AssertEqualInt(t, indentationRelationshipChild, relationship, "Initial node with no prior indentation should be a sibling")
	test.AssertEqualInt(t, 4, indentCount, "Initial indentation count should be 0")

	m.push(4)

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
