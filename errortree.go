package errortree

import (
	"fmt"
)

type Node struct {
	Errs       []error
	Attributes map[string]*Node
	Elements   map[int]*Node
}

func (n *Node) Add(path []any, err error) {
	if len(path) == 0 {
		n.Errs = append(n.Errs, err)
		return
	}

	switch step := path[0].(type) {
	case string:
		if n.Attributes == nil {
			if n.Elements != nil {
				panic("node is already acting as list")
			}
			n.Attributes = make(map[string]*Node)
		}

		var nextNode *Node
		nextNode = n.Attributes[step]
		if nextNode == nil {
			nextNode = &Node{}
			n.Attributes[step] = nextNode
		}
		nextNode.Add(path[1:], err)

	case int:
		if n.Elements == nil {
			if n.Attributes != nil {
				panic("node is already acting as object")
			}
			n.Elements = make(map[int]*Node)
		}

		var nextNode *Node
		nextNode = n.Elements[step]
		if nextNode == nil {
			nextNode = &Node{}
			n.Elements[step] = nextNode
		}
		nextNode.Add(path[1:], err)

	default:
		panic("path elements must be string or int")
	}
}

func (n *Node) Get(path []any) []error {
	if len(path) == 0 {
		return n.Errs
	}

	switch step := path[0].(type) {
	case string:
		if n.Attributes == nil {
			return nil
		}

		nextNode := n.Attributes[step]
		if nextNode == nil {
			return nil
		}
		return nextNode.Get(path[1:])

	case int:
		if n.Elements == nil {
			return nil
		}

		nextNode := n.Elements[step]
		if nextNode == nil {
			return nil
		}
		return nextNode.Get(path[1:])

	default:
		panic("path elements must be string or int")
	}
}

func Example() {
	fmt.Println("placeholder")

	// Output:
	// placeholder
}
