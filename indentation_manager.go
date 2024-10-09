package yaml

type indentationRelationship int8

const (
	indentationRelationshipUnknown indentationRelationship = iota
	indentationRelationshipChild
	indentationRelationshipSibling
	indentationRelationshipParentSibling
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
	m.indentations = m.indentations[:len(m.indentations)-2]
}

func (m *indentationManager) pushStack(indentation int) {
	m.indentations = append(m.indentations, indentation)
}

func (m *indentationManager) determineRelationship(indentationLevel int) indentationRelationship {
	siblingIndentationLevel := m.indentations[len(m.indentations)-1]
	if siblingIndentationLevel == indentationLevel {
		return indentationRelationshipSibling
	}

	if siblingIndentationLevel > indentationLevel {
		return indentationRelationshipChild
	}

	if len(m.indentations) > 1 {
		parentIndentationLevel := m.indentations[len(m.indentations)-2]
		if parentIndentationLevel == indentationLevel {
			return indentationRelationshipParentSibling
		}
	}

	return indentationRelationshipUnknown
}
