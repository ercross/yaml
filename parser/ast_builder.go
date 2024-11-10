package parser

import (
	"errors"
	"fmt"
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/token"
)

type AstBuilder struct {
	stack          *stack
	ast            *AbstractSyntaxTree
	awaitingParse  []token.Token
	nodeTypeFinder *nodeTypeFinder
}

func NewAstBuilder() *AstBuilder {
	initTokenTrie()
	return &AstBuilder{
		stack:          newStack(),
		ast:            newAbstractSyntaxTree(),
		nodeTypeFinder: newNodeTypeFinder(),
	}
}

func (builder *AstBuilder) AbstractSyntaxTree() *AbstractSyntaxTree {
	return builder.ast
}

func (builder *AstBuilder) createNewDocument() {
	// don't create a new document if first document has not been used
	if len(builder.ast.documents[0].Children()) == 0 {
		return
	}

	builder.unwindStack()
	builder.ast.documents = append(builder.ast.documents, yaml.NewDocumentNode())
}

// unwindStack removes all frames on stack (in a stack LIFO manner)
// and adds them to the current AbstractSyntaxTree.documents being built
func (builder *AstBuilder) unwindStack() {

	if builder.stack.isEmpty() {
		return
	}

	frame := builder.stack.pop()
	for builder.stack.size() > 0 {
		builder.stack.peek().Builder().AddChild(frame.Builder().ToNode())
		frame = builder.stack.pop()
	}
	builder.ast.addChild(frame.Builder().ToNode())

	builder.stack.clear()
}

// Build parses tokens, builds yaml.Node, and inserts the built nodes to AstBuilder.AbstractSyntaxTree
//
// Build maintains an internal state, which enables it to continuously build over multiple invocations
func (builder *AstBuilder) Build(tokens []token.Token) error {
	if len(tokens) == 0 {
		return errors.New("can not parse empty tokens")
	}

	if tokens[0].Type == token.TypeDocumentStart || tokens[0].Type == token.TypeDocumentEnd {
		builder.createNewDocument()
		return nil
	}

	builder.nodeTypeFinder.match(tokens)
	if !builder.nodeTypeFinder.done {

		builder.awaitingParse = append(builder.awaitingParse, tokens...)

		// continue finding on getting the next set of tokens since nodeTypeFinder is not done
		return nil
	}

	if builder.nodeTypeFinder.position.nodeType == yaml.NodeTypeUnknown {
		return fmt.Errorf("can not determine node type on %d", builder.nodeTypeFinder.position.tokenType)
	}

	tokens = append(builder.awaitingParse, tokens...)

	relationship, indentationLength := builder.stack.indentationManager.findIndentation(tokens)
	if relationship == indentationRelationshipUnknown {
		return fmt.Errorf("can not find indentation")
	}

	frame, err := builder.createNewFrame(builder.nodeTypeFinder.position.nodeType, indentationLength)
	if err != nil {
		return fmt.Errorf("failed to create new %d frame: %w", builder.nodeTypeFinder.position.nodeType, err)
	}

	err = builder.pushOnStack(frame, relationship)
	if err != nil {
		return fmt.Errorf("failed to push frame on stack: %w", err)
	}

	err = builder.stack.peek().Build(tokens)
	if err != nil {
		return fmt.Errorf("error building %d frame near line %s",
			builder.stack.peek().NodeType(), builder.stack.peek().Builder().CurrentPosition())
	}

	builder.nodeTypeFinder.reset()
	builder.awaitingParse = []token.Token{}
	return nil
}

func (builder *AstBuilder) createNewFrame(nt yaml.NodeType, indentation int) (Frame, error) {
	var frame Frame
	switch nt {
	case yaml.NodeTypeScalar:
		frame = newScalarFrame(indentation, newNodeSyntaxTraverser(scalarNodeSyntax().head))

	default:
		return nil, fmt.Errorf("can not handle NodeType %d", nt)
	}
	return frame, nil
}

func (builder *AstBuilder) handlePoppedFrame(poppedFrame Frame) {

	// if stack is empty, node is an independent entry of the AstBuilder ast
	if builder.stack.isEmpty() {
		builder.ast.addChild(poppedFrame.Builder().ToNode())
		return
	} else { // frame is a child of current stack-top frame
		builder.stack.peek().Builder().AddChild(poppedFrame.Builder().ToNode())
	}
}

func (builder *AstBuilder) pushOnStack(frame Frame, relationshipWithLastFrame indentationRelationship) error {
	switch relationshipWithLastFrame {
	case indentationRelationshipChild:
		builder.stack.push(frame)

	case indentationRelationshipParentLevel:
		if builder.stack.size() < 2 {
			panic("incorrect indentation relationship. " +
				"frame can not be on parent level indentation when stack size is less than 2")
		}
		frameSiblingChild := builder.stack.pop()
		builder.handlePoppedFrame(frameSiblingChild)
		frameSibling := builder.stack.pop()
		builder.handlePoppedFrame(frameSibling)
		builder.stack.push(frame)

	case indentationRelationSibling:
		if builder.stack.isEmpty() {
			builder.stack.push(frame)
		} else {
			frameSibling := builder.stack.pop()
			builder.handlePoppedFrame(frameSibling)
			builder.stack.push(frame)
		}

	default:
		return errors.New("can not handle new indentation relationship")
	}

	return nil
}
