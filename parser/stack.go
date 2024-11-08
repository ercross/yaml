package parser

// stack serves as the building stage for yaml.Node using Frame
type stack struct {
	elements           []Frame
	indentationManager *indentationManager
}

func newStack() *stack {
	return &stack{
		elements:           make([]Frame, 0),
		indentationManager: newIndentationManager(),
	}
}

func (s *stack) push(frame Frame) {
	s.elements = append(s.elements, frame)
	s.indentationManager.push(frame.AllowedIndentationLevel(), frame.NodeType())
}
func (s *stack) pop() Frame {
	if len(s.elements) == 0 {
		panic("cannot pop from empty stack")
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	s.indentationManager.pop()
	return element
}

func (s *stack) peek() Frame {
	if len(s.elements) == 0 {
		panic("cannot peek from empty stack")
	}
	return s.elements[len(s.elements)-1]
}
func (s *stack) isEmpty() bool {
	return len(s.elements) == 0
}
func (s *stack) size() int {
	return len(s.elements)
}

func (s *stack) clear() {
	s.elements = []Frame{}
	s.indentationManager = newIndentationManager()
}
