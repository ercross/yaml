package yaml

import "fmt"

type (
	frameType          int8
	scalarDataType     int8
	sequenceSyntaxType int8
)

const (
	scalarDataTypeUnknown scalarDataType = iota
	scalarDataTypeString
	scalarDataTypeFloat
	scalarDataTypeInteger
	scalarDataTypeBoolean
	scalarDataTypeTypeNull
)

const (
	frameTypeUnknown frameType = iota
	frameTypeSequenceParser
	frameTypeMappingParser
	frameTypeScalarParser
	frameTypeMultilineStringParser
	frameTypeFoldedStringParser
	frameTypeComplexKeyParser
	frameTypeAnchorParser
)

const (
	sequenceSyntaxTypeUnknown sequenceSyntaxType = iota
	sequenceSyntaxTypeInline
	sequenceSyntaxTypeList
)

type parserStackFrame interface {
	Type() frameType

	// Parse the given token
	//
	// If remaining is not empty, then current frame can be taken off the stack
	// and the node resulting from parsing tokens can be returned through Node
	Parse([]token) (remaining []token, childFrame parserStackFrame, err error)

	// Builder building the underlying Node
	Builder() NodeBuilder

	// IndentationLevel is usually set by the indentation that preceded the first token
	// parsed into this frame
	IndentationLevel() int
}

type (
	sequenceFrame struct {
		previousTokens        []token
		syntaxType            sequenceSyntaxType
		allowedNextTokenTypes map[tokenType]bool
		builder               NodeBuilder
		indentationManager    *indentationManager
	}

	unknownFrame struct {
		previousTokens     []token
		key                string
		indentationManager *indentationManager
	}

	mappingFrame struct {
		previousTokens        []token
		allowedNextTokenTypes map[tokenType]bool
		builder               NodeBuilder
		indentationManager    *indentationManager
	}

	scalarFrame struct {
		builder            NodeBuilder
		indentationManager *indentationManager
	}

	multilineStringFrame struct {
		builder            NodeBuilder
		indentationManager *indentationManager
	}

	foldedStringFrame struct {
		builder            NodeBuilder
		indentationManager *indentationManager
	}

	complexKeyFrame struct {
		builder            NodeBuilder
		indentationManager *indentationManager
	}

	anchorFrame struct {
		builder            NodeBuilder
		indentationManager *indentationManager
	}
)

func newSequenceFrame(key string, syntaxType sequenceSyntaxType, idm *indentationManager) *sequenceFrame {
	return &sequenceFrame{
		allowedNextTokenTypes: map[tokenType]bool{
			// TODO initialize with allowed token types after key has been set
		},
		syntaxType: syntaxType,
		builder:    &SequenceNode{key: key},
	}
}

func (f *sequenceFrame) Type() frameType {
	return frameTypeSequenceParser
}

func (f *sequenceFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {
	// pick only the child nodes from tokens and check if some part of tokens can indicate a change in parser state
	for i := 0; i < len(tokens); i++ {
		if !f.isAnExpectedTokenType(tokens[i].Type) {
			return tokens[i:], nil, fmt.Errorf("%v: %s", ErrUnexpectedTokenType, tokens[i].position)
		}

		if tokens[i].Type != TokenTypeData {

		}
	}
}

func (f *sequenceFrame) isKeyData(targetTokenIndex int, tokens []token) bool {
	// if has appropriate number of indent
}

func (f *sequenceFrame) isScalarValueData() bool {}

func (f *sequenceFrame) Builder() NodeBuilder {
	return f.builder
}

func newMappingFrame(key string, idm *indentationManager) *mappingFrame {}

func (f *mappingFrame) Type() frameType {
	return frameTypeMappingParser
}

func (f *mappingFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *mappingFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *scalarFrame) Type() frameType {
	return frameTypeScalarParser
}

func (f *scalarFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *scalarFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *multilineStringFrame) Type() frameType {
	return frameTypeMultilineStringParser
}

func (f *multilineStringFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *multilineStringFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *foldedStringFrame) Type() frameType {
	return frameTypeFoldedStringParser
}

func (f *foldedStringFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *foldedStringFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *complexKeyFrame) Type() frameType {
	return frameTypeComplexKeyParser
}

func (f *complexKeyFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *complexKeyFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *anchorFrame) Type() frameType {
	return frameTypeAnchorParser
}

func (f *anchorFrame) Parse(tokens []token) (remaining []token, childFrame parserStackFrame, err error) {

}

func (f *anchorFrame) Builder() NodeBuilder {
	return f.builder
}

func (f *sequenceFrame) isAnExpectedTokenType(t tokenType) bool {
	_, ok := f.allowedNextTokenTypes[t]
	return ok
}

func newUnknownFrame(tokens []token) *unknownFrame {
	return &unknownFrame{
		previousTokens: tokens,
	}
}

// attemptToDetermineFrameType attempts to determine the frame type formed by parsing tokens
// and initialized the necessary data for the newly created frame
//
// If frame type is successfully determined and remainingTokens is not empty,
// then the frame can be considered complete and can be popped
//
// if createNewDocument is true, this indicates that tokens contain the YAML document-start tokens (i.e., ---)
func (f *unknownFrame) attemptToDetermineFrameType(tokens []token) (frame parserStackFrame, createNewDocument bool, remainingTokens []token) {
	f.previousTokens = append(f.previousTokens, tokens...)

	// check if it's a document start
	documentStartExpectedDashCount := 3
	if len(f.previousTokens) >= documentStartExpectedDashCount {
		isDocumentStart := true
		for i := 0; i < documentStartExpectedDashCount; i++ {
			if f.previousTokens[i].Type != TokenTypeDash {
				isDocumentStart = false
				break
			}
		}

		if isDocumentStart {
			return nil, true, f.previousTokens[documentStartExpectedDashCount:]
		}
	}

	if sf, remaining := makeSequenceFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeMappingFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeScalarFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeAnchorFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeComplexKeyFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeFoldedStringFrame(tokens); sf != nil {
		return sf, false, remaining
	}
	if sf, remaining := makeMultilineStringFrame(tokens); sf != nil {
		return sf, false, remaining
	}

	return nil, false, nil
}

func (f *unknownFrame) clearAccumulatedData() {
	f.previousTokens = []token{}
	f.key = ""
}

func (f *unknownFrame) containsData() bool {
	return len(f.previousTokens) > 0
}

func makeSequenceFrame(tokens []token) (frame *sequenceFrame, remainingTokens []token) {

}

func makeMappingFrame(tokens []token) (frame *mappingFrame, remainingTokens []token) {}

func makeScalarFrame(tokens []token) (frame *scalarFrame, remainingTokens []token) {}

func makeMultilineStringFrame(tokens []token) (frame *multilineStringFrame, remainingTokens []token) {
}

func makeFoldedStringFrame(tokens []token) (frame *foldedStringFrame, remainingTokens []token) {}

func makeComplexKeyFrame(tokens []token) (frame *complexKeyFrame, remainingTokens []token) {}

func makeAnchorFrame(tokens []token) (frame *anchorFrame, remainingTokens []token) {}
