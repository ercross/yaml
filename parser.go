package yaml

import "fmt"

var ErrUnexpectedTokenType = fmt.Errorf("unexpected token type")

type parser struct {
	stack        *parserStack
	anchors      map[string]Node
	ast          *abstractSyntaxTree
	unknownFrame *unknownFrame
}

func newParser(maxNestingLevel int) *parser {
	return &parser{
		stack:   newStack(maxNestingLevel),
		anchors: make(map[string]Node),
	}
}

func (p *parser) parseTokens(tokens []token, poppedFrame parserStackFrame) error {
	if poppedFrame != nil {
		p.handlePoppedFrame(poppedFrame)
	}

	if p.stack.isEmpty() && p.unknownFrame.containsData() {
		frame, createNewDocument, remaining := p.unknownFrame.attemptToDetermineFrameType(tokens)
		if frame != nil {
			p.stack.push(frame)
			p.unknownFrame.clearAccumulatedData()
		}
		if createNewDocument {
			p.ast.documents = append(p.ast.documents, newDocumentNode())
		}
		if remaining != nil {
			return p.parseTokens(remaining, nil)
		}
		return nil
	}

	frame := p.stack.peek()
	remaining, childFrame, err := frame.Parse(tokens)
	if err != nil {
		return fmt.Errorf("error parsing tokens: %w", err)
	}
	if childFrame != nil {
		p.stack.push(childFrame)
	}
	if len(remaining) > 0 {
		return p.parseTokens(remaining, p.stack.pop())
	}

	return nil
}

func (p *parser) handlePoppedFrame(poppedFrame parserStackFrame) {

	// if poppedFrame is an anchor, add its Node to parser anchors
	if poppedFrame.Type() == frameTypeAnchorParser {
		p.anchors[poppedFrame.Builder().ToNode().Key()] = poppedFrame.Builder().ToNode()
		return
	}

	// if stack is empty, node is an independent entry of the parser ast
	if p.stack.isEmpty() {
		p.ast.addChild(poppedFrame.Builder().ToNode())
		return
	}

	if p.stack.peek().Type() == frameTypeScalarParser {
		panic("scalar frame can not handle child")
	}

	// if stack is not empty, add node to next frame child
	p.stack.peek().Builder().AddChild(poppedFrame.Builder().ToNode())
}
