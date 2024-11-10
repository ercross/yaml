package parser

import (
	"fmt"
	"github.com/ercross/yaml"
	"github.com/ercross/yaml/token"
)

type (

	// tokenTrie is a trie that holds the syntax definition for all yaml.Node.
	// Essentially tokenTrie is a nodeSyntax holder
	//
	// Each path from the root to a leaf represents the allowed sequence of tokens
	// that defines the syntax a specific node type in the YAML specification.
	// A leaf in the tokenTrie contains the yaml.NodeType
	tokenTrie struct {
		root *tokenTrieNode
	}

	tokenTrieNode struct {
		tokenType token.Type
		children  []*tokenTrieNode

		// nodeType not empty indicates this node is a leaf
		nodeType yaml.NodeType
	}

	// nodeTypeFinder is a tokenTrie traverser capable of finding which yaml.Node
	// a sequence of token.Token attempts to represent
	nodeTypeFinder struct {
		position *tokenTrieNode

		// done indicates that nodeTypeFinder has concluded its finding.
		// The result of its search can be found in position.nodeType
		done bool
	}
)

type (
	nodeSyntaxToken struct {
		optional  bool
		tokenType token.Type
		next      *nodeSyntaxToken
	}

	// nodeSyntax (implemented as a singly-LinkedList) is a sequence of tokens
	// that defines the syntax a specific yaml.NodeType in the YAML specification.
	//
	// Note that nodeSyntax should not include token.TypeIndentation
	nodeSyntax struct {
		head *nodeSyntaxToken
		size int
	}

	nodeSyntaxTraverser struct {
		current *nodeSyntaxToken
	}
)

var st *tokenTrie

func initTokenTrie() {
	st = &tokenTrie{
		root: &tokenTrieNode{},
	}

	st.insertNodeSyntax(scalarNodeSyntax(), yaml.NodeTypeScalar)
}

func (t *tokenTrie) insertNodeSyntax(ts *nodeSyntax, f yaml.NodeType) {
	i := newNodeSyntaxTraverser(ts.head)
	currentNode := t.root
	var next *nodeSyntaxToken
	for i.hasNext() {
		next = i.next()
		if currentNode.tokenType == next.tokenType {
			continue
		}

		if index := containsTokenType(currentNode.children, next.tokenType); index != -1 {
			currentNode = currentNode.children[index]
		} else {
			currentNode.children = append(currentNode.children, &tokenTrieNode{
				tokenType: next.tokenType,
				children:  nil,
			})
			currentNode = currentNode.children[len(currentNode.children)-1]
		}
	}

	currentNode.children = append(currentNode.children, &tokenTrieNode{
		nodeType: f,
	})
}

func containsTokenType(nodes []*tokenTrieNode, t token.Type) int {
	for i, node := range nodes {
		if node.tokenType == t {
			return i
		}
	}
	return -1
}

func newNodeTypeFinder() *nodeTypeFinder {
	return &nodeTypeFinder{
		position: st.root,
	}
}

// search for Frame using Depth-First search algorithm
func (f *nodeTypeFinder) findMatch(tokens []token.Token) {

	for _, next := range tokens {

		// indentation is not part of syntax tree
		if next.Type == token.TypeIndentation {
			continue
		}
		if f.position.nodeType != yaml.NodeTypeUnknown {
			f.done = true
			return
		}
		if f.position.tokenType == next.Type {
			continue
		}

		for _, child := range f.position.children {
			if child.tokenType == next.Type {
				f.position = child
				if f.position.nodeType != yaml.NodeTypeUnknown {
					f.done = true
					return
				}
				if len(f.position.children) == 1 && f.position.children[0].nodeType != yaml.NodeTypeUnknown {
					f.position.nodeType = f.position.children[0].nodeType
				}
				break
			}
		}
	}
}

func (f *nodeTypeFinder) nodeType() yaml.NodeType {
	if !f.done {
		panic("finder not done")
	}
	return f.position.nodeType
}

func newNodeSyntaxTraverser(start *nodeSyntaxToken) *nodeSyntaxTraverser {
	return &nodeSyntaxTraverser{current: start}
}

func (i *nodeSyntaxTraverser) hasNext() bool {
	return i.current != nil
}

// next moves the iterator to the next nodeSyntaxToken and returns it
func (i *nodeSyntaxTraverser) next() *nodeSyntaxToken {
	next := i.current
	i.current = i.current.next
	return next
}

func newNodeSyntax(head *nodeSyntaxToken) *nodeSyntax {
	return &nodeSyntax{
		head: head,
		size: 1,
	}
}

func (ts *nodeSyntax) insert(n *nodeSyntaxToken) *nodeSyntax {
	current := ts.head
	for current.next != nil {
		current = current.next
	}
	current.next = n
	ts.size++
	return ts
}

// createCycle inside nodeSyntax
// Note that nodeSyntax.head is at index 0
func (ts *nodeSyntax) createCycle(from, to int) *nodeSyntax {
	if ts.size < from {
		panic(fmt.Sprintf("can not create cycle: from is out of range: from(%d) > size(%d)", from, ts.size))
	}
	if ts.size < to {
		panic(fmt.Sprintf("can not create cycle: to is out of range: to(%d) > size(%d)", to, ts.size))
	}
	var (
		fromNode *nodeSyntaxToken
		toNode   *nodeSyntaxToken
	)

	current := ts.head
	for i := 0; i <= from; i++ {
		current = current.next
	}
	fromNode = current

	current = ts.head
	for i := 0; i <= to; i++ {
		current = current.next
	}
	toNode = current

	fromNode.next = toNode

	return ts
}

func scalarNodeSyntax() *nodeSyntax {
	return newNodeSyntax(&nodeSyntaxToken{optional: false, tokenType: token.TypeData}).
		insert(&nodeSyntaxToken{optional: false, tokenType: token.TypeColon}).
		insert(&nodeSyntaxToken{optional: false, tokenType: token.TypeData}).
		insert(&nodeSyntaxToken{optional: false, tokenType: token.TypeNewline})
}
