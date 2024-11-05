package yaml

type (
	NodeType int8
)

const (
	NodeTypeUnknown NodeType = iota // Node type is unknown and used as a placeholder for uninitialized or invalid nodes.

	// NodeTypeScalar represents a single value node, such as a string, integer, or boolean.
	// Example:
	//   key: "Hello, World"
	NodeTypeScalar

	// NodeTypeDocument is the root of a YAML document. It represents the whole document structure and may contain sequences or mappings.
	// Example:
	//   ---
	//   title: "YAML Example"
	//   ---
	NodeTypeDocument

	// NodeTypeMultilineString represents a multi-line string, often using the `|` symbol to preserve line breaks.
	// Example:
	//   description: |
	//     This is a multi-line
	//     string in YAML.
	NodeTypeMultilineString

	// NodeTypeFoldedString is a multi-line text node that uses the `>` symbol to fold lines. Line breaks within folded text are converted to spaces.
	// Example:
	//   note: >
	//     This text will be folded
	//     into a single line when parsed.
	NodeTypeFoldedString

	// NodeTypeAnchor represents an anchor node, allowing values to be reused or referenced elsewhere in the document.
	// Example:
	//   base: &baseAnchor "Base Value"
	NodeTypeAnchor

	// NodeTypeSequenceFlowStyle represents a sequence in flow style, denoted by square brackets `[]`. This style is non-nestable.
	// Example:
	//   items: [1, 2, 3]
	NodeTypeSequenceFlowStyle

	// NodeTypeSequenceBlockStyle represents a sequence in block style, using `-` to denote each item. This style can contain nested elements.
	// Example:
	//   items:
	//     - name: "Item 1"
	//     - name: "Item 2"
	NodeTypeSequenceBlockStyle

	// NodeTypeMappingFlowStyle represents a mapping in flow style, using curly braces `{}`. This style is non-nestable.
	// Example:
	//   info: { key1: "value1", key2: "value2" }
	NodeTypeMappingFlowStyle

	// NodeTypeMappingBlockStyle represents a mapping in block style, where each key-value pair is on a new line.
	// This style can contain nested mappings, sequences, or scalars.
	// Example:
	//   user:
	//     name: "Alice"
	//     age: 30
	NodeTypeMappingBlockStyle

	// NodeTypeAlias represents a reference to an anchor node, allowing the reuse of values defined by an anchor.
	// Example:
	//   base: &baseAnchor "Base Value"
	//   alias: *baseAnchor
	NodeTypeAlias
)

type (
	Node interface {
		Key() string
		Type() NodeType
		Value() interface{}
		Children() any
	}

	NodeBuilder interface {
		AddChild(Node)
		SetValue(any)
		ToNode() Node
		SetKey(key string)
	}
)

type (
	// DocumentNode is usually the root of an AbstractSyntaxTree
	DocumentNode struct {
		children []Node
	}

	ScalarNode struct {
		value any
		key   string
	}
)

type AbstractSyntaxTree struct {
	documents []*DocumentNode
}

// IsNestable checks if NodeType can serve as a parent node to other child nodes.
// In the context of YAML, only certain node types can nest other nodes.
// Specifically, NodeTypeSequenceBlockStyle and NodeTypeMappingBlockStyle are nestable,
// while other types such as scalars, flow styles, and aliases are not.
func (nt NodeType) IsNestable() bool {
	switch nt {
	case NodeTypeSequenceBlockStyle, NodeTypeMappingBlockStyle:
		return true
	default:
		return false
	}
}

func NewAbstractSyntaxTree() *AbstractSyntaxTree {
	return &AbstractSyntaxTree{
		documents: []*DocumentNode{
			newDocumentNode(),
		},
	}
}

// addChild adds n to the yaml document currently being built
func (ast *AbstractSyntaxTree) addChild(n Node) {
	currentlyParsedDocumentIndex := len(ast.documents) - 1
	ast.documents[currentlyParsedDocumentIndex].AddChild(n)
}

func newDocumentNode() *DocumentNode {
	return &DocumentNode{}
}

func (n *DocumentNode) Key() string {
	return ""
}

func (n *DocumentNode) Type() NodeType {
	return NodeTypeDocument
}

func (n *DocumentNode) Value() any {
	return nil
}

func (n *DocumentNode) Children() []Node {
	return n.children
}

func (n *DocumentNode) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *DocumentNode) SetValue(interface{}) {
	return
}

func NewScalarNodeBuilder() *ScalarNode {
	return &ScalarNode{}
}

func (n *ScalarNode) Key() string {
	return n.key
}

func (n *ScalarNode) Type() NodeType {
	return NodeTypeScalar
}

func (n *ScalarNode) Value() any {
	return n.value
}

func (n *ScalarNode) Children() any {
	panic("can not get children of scalar node")
}

func (n *ScalarNode) AddChild(_ Node) {
	panic("can not add child to scalar node")
}
func (n *ScalarNode) SetValue(v any) {
	n.value = v
}

func (n *ScalarNode) SetKey(key string) {
	n.key = key
}

func (n *ScalarNode) ToNode() Node {
	return n
}
