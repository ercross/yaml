package yaml

type parserStack struct {
	elements []parserStackFrame
}

func newStack(maxStackDepth int) *parserStack {
	return &parserStack{

		// allow at most 50 levels of stack depth
		elements: make([]parserStackFrame, 0, maxStackDepth),
	}
}

func (s *parserStack) push(node parserStackFrame) {
	s.elements = append(s.elements, node)
}
func (s *parserStack) pop() parserStackFrame {
	if len(s.elements) == 0 {
		panic("cannot pop from empty stack")
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element
}

func (s *parserStack) peek() parserStackFrame {
	if len(s.elements) == 0 {
		panic("cannot peek from empty stack")
	}
	return s.elements[len(s.elements)-1]
}
func (s *parserStack) isEmpty() bool {
	return len(s.elements) == 0
}
func (s *parserStack) size() int {
	return len(s.elements)
}
