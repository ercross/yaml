package parser

import (
	"errors"
	"fmt"
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/token"
)

var (
	// errIncompatibleIndentationEncountered is returned if a TokenTypeIndentation with indentation length
	// longer than Frame.AllowedIndentationLevel is found.
	// This error might indicate to pause the current Frame,
	// especially if the encountered indentationLevel is lower than the current Frame.AllowedIndentationLevel
	errIncompatibleIndentationEncountered = errors.New("incompatible indentation encountered")

	// errNewKeyEncountered is returned if newly encountered key is preceded by the same level of indentation
	// as the current Frame.AllowedIndentationLevel
	errNewKeyEncountered = errors.New("new key encountered")

	errUnexpectedTokenType = errors.New("unexpected token type")

	errFrameTypeUnknown = errors.New("frame type unknown")
)

type Frame interface {
	NodeType() yaml.NodeType

	Parse([]token.Token) error

	// Builder building the underlying Node
	Builder() yaml.NodeBuilder

	// AllowedIndentationLevel is usually set by the indentation preceding the first token
	// parsed into this Frame
	AllowedIndentationLevel() int
}

type (
	scalarFrame struct {
		sequenceIterator *nodeSyntaxTraverser
		builder          *yaml.ScalarNode
		indentationLevel int
	}
)

func newScalarFrame(indentationLevel int, iterator *nodeSyntaxTraverser) *scalarFrame {
	return &scalarFrame{
		builder:          yaml.NewScalarNodeBuilder(),
		indentationLevel: indentationLevel,
		sequenceIterator: iterator,
	}
}

func (f *scalarFrame) NodeType() yaml.NodeType {
	return yaml.NodeTypeScalar
}

func (f *scalarFrame) Parse(tokens []token.Token) error {
	var hasVisitedAllowedIndentation bool
	i := 0
	for f.sequenceIterator.hasNext() && i < len(tokens) {
		i++
		if tokens[i].Type == token.TypeIndentation {
			if !hasVisitedAllowedIndentation && f.Builder().ToNode().Value() == nil {
				hasVisitedAllowedIndentation = true
				continue
			} else {
				return fmt.Errorf("indentation at %s unsupported: %w", tokens[i].Position, errUnexpectedTokenType)
			}
		}

		expected := f.sequenceIterator.next()

		if expected.optional && tokens[i].Type != expected.tokenType {
			continue
		}

		if expected.tokenType != tokens[i].Type {
			return errUnexpectedTokenType
		}

		// newline token is the last token in a scalar frame syntax
		if expected.tokenType == token.TypeNewline && f.Builder().ToNode().Value() != nil && f.Builder().ToNode().Key() != "" {
			return nil
		}

		if tokens[i].Type == token.TypeData {
			if f.builder.Key() == "" {
				f.Builder().SetKey(tokens[i].Value)
			} else {
				f.Builder().SetValue(tokens[i].Value)
			}
		}
	}

	return nil
}

func (f *scalarFrame) Builder() yaml.NodeBuilder {
	return f.builder
}

func (f *scalarFrame) AllowedIndentationLevel() int {
	return f.indentationLevel
}
