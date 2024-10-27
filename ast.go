package yaml

type (
	NodeType int8
)

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeScalar
	NodeTypeDocument
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

func newScalarNodeBuilder() *ScalarNode {
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
