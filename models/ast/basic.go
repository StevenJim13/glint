package ast

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
)

// QueryCondition used to determine whether a node meets certain conditions
// It returns true when the conditions are met else false
type QueryCondition func(n *sitter.Node) bool

// QueryChild Query child node
func QueryChild(node *sitter.Node, condition QueryCondition) *sitter.Node {
	count := int(node.ChildCount())
	if count == 0 {
		return nil
	}
	for i := 0; i < count; i++ {
		child := node.Child(i)
		if condition(child) {
			return child
		}
	}
	return nil
}

// QueryAncestor search for a node that meets the condition in the ancestor nodes,
// and return immediately when found.
func QueryAncestor(node *sitter.Node, condition QueryCondition) *sitter.Node {
	for node != nil {
		if condition(node) {
			return node
		} else {
			node = node.Parent()
		}
	}
	return nil
}

// ApplyLevelNodes ...
func ApplyLevelNodes(node *sitter.Node, fn func(sub *sitter.Node)) {
	for node != nil {
		fn(node)
		node = node.NextSibling()
	}
}

// ApplyChildrenNodes ...
func ApplyChildrenNodes(node *sitter.Node, fn func(sub *sitter.Node)) {
	count := int(node.ChildCount())
	for i := 0; i < count; i++ {
		fn(node.Child(i))
	}
}

func InspectNode(node *sitter.Node, content []byte, indent string) {
	fmt.Printf("%stype: %s, content: %s\n", indent, node.Type(), node.Content(content))
}

func InspectChildren(node *sitter.Node, content []byte) {
	count := int(node.ChildCount())
	InspectNode(node, content, "")
	for i := 0; i < count; i++ {
		child := node.Child(i)
		InspectNode(child, content, "    ")
	}
}

func InspectTree(node *sitter.Node, content []byte, indent string) {
	fmt.Println(indent, node.Type(), node.Content(content), node.StartPoint())
	count := int(node.ChildCount())
	for i := 0; i < count; i++ {
		child := node.Child(i)
		InspectTree(child, content, indent+"    ")
	}
}

const (
	Contine = iota
	Break
	Skip
)

func DFVisit(node *sitter.Node, f func(node *sitter.Node) int) {

	if node == nil {
		return
	}
	count := int(node.ChildCount())
	for i := 0; i < count; i++ {
		child := node.Child(i)
		switch f(child) {
		case Contine:
			DFVisit(child, f)
		case Break:
			return
		case Skip:
		}
	}
}

func BFVisit(node *sitter.Node, f func(node *sitter.Node) int) {
	if node == nil {
		return
	}
	queue := []*sitter.Node{}
	for len(queue) > 0 {
		size := len(queue)
		for i := 0; i < size; i++ {
			child := queue[0]
			queue = queue[1:]
			switch f(child) {
			case Break:
				return
			case Contine:
				count := int(node.ChildCount())
				for j := 0; j < count; j++ {
					queue = append(queue, child.Child(j))
				}
			}
		}
	}
}

// NodeLines ...
func NodeLines(node *sitter.Node) int {
	if node == nil {
		return 0
	}
	return int(node.EndPoint().Row - node.StartPoint().Row)
}
