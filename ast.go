package yaml

type NodeType string

const (
	NodeTypeScalar   NodeType = "scalar"
	NodeTypeDocument NodeType = "document"
	NodeTypeMap               = "map"
	NodeTypeSequence          = "sequence"
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
	}
)

type (
	// DocumentNode is usually the root of an abstractSyntaxTree
	DocumentNode struct {
		children []Node
	}

	ScalarNode struct {
		value any
		key   string
	}

	MappingNode struct {
		key      string
		children map[string]Node
	}

	SequenceNode struct {
		key      string
		children []Node
	}
)

type abstractSyntaxTree struct {
	documents []*DocumentNode
}

func newAbstractSyntaxTree() *abstractSyntaxTree {
	return &abstractSyntaxTree{
		documents: []*DocumentNode{
			newDocumentNode(),
		},
	}
}

func (ast *abstractSyntaxTree) addChild(n Node) {
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

func newScalarNodeBuilder(key string) NodeBuilder {
	return &ScalarNode{
		key: key,
	}
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
	return nil
}

func (n *ScalarNode) AddChild(_ Node) {
	// no-op because scalars do not have children nodes
	return
}
func (n *ScalarNode) SetValue(v any) {
	n.value = v
}

func (n *ScalarNode) ToNode() Node {
	return n
}

func newMappingNodeBuilder(key string) NodeBuilder {
	return &MappingNode{
		key: key,
	}
}

func (n *MappingNode) AddChild(node Node) {

	n.children[node.Key()] = node
}

func (n *MappingNode) SetValue(v any) {
	return
}

func (n *MappingNode) ToNode() Node {
	return n
}

func (n *MappingNode) Key() string {
	return n.key
}

func (n *MappingNode) Type() NodeType {
	return NodeTypeMap
}

func (n *MappingNode) Value() any {
	return nil
}

func (n *MappingNode) Children() any {
	return n.children
}

func newSequenceNodeBuilder(key string) NodeBuilder {
	return &SequenceNode{
		key: key,
	}
}

func (n *SequenceNode) Key() string {
	return n.key
}

func (n *SequenceNode) Type() NodeType {
	return NodeTypeSequence
}

func (n *SequenceNode) Value() interface{} {
	return nil
}

func (n *SequenceNode) Children() any {
	return n.children
}

func (n *SequenceNode) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *SequenceNode) SetValue(v any) {
	return
}

func (n *SequenceNode) ToNode() Node {
	return n
}
