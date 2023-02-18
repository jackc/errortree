package errortree

import (
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	Errs       []error
	Attributes map[string]*Node
	Elements   map[int]*Node
}

func (n *Node) Error() string {
	sb := &strings.Builder{}

	errs := n.AllErrors()
	for i, err := range errs {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(err.Error())
	}

	return sb.String()
}

// AllErrors returns all errors in the node and its descendents.
func (n *Node) AllErrors() []*ErrorWithPath {
	return n.errorsWithPath(nil)
}

func (n *Node) errorsWithPath(path []any) []*ErrorWithPath {
	var errs []*ErrorWithPath

	// Ensure that the path an ErrorWithPath stores is not mutated by later calls to errorsWithPath.
	{
		pathCopy := make([]any, len(path))
		copy(pathCopy, path)
		path = pathCopy
	}

	for _, err := range n.Errs {
		errs = append(errs, &ErrorWithPath{Path: path, Err: err})
	}

	path = append(path, nil)

	// Sort keys and iterate for stable order.
	fieldNames := make([]string, 0, len(n.Attributes))
	for fieldName := range n.Attributes {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		path[len(path)-1] = fieldName
		errs = append(errs, n.Attributes[fieldName].errorsWithPath(path)...)
	}

	// Sort keys and iterate for stable order.
	indexes := make([]int, 0, len(n.Elements))
	for index := range n.Elements {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	for _, index := range indexes {
		path[len(path)-1] = index
		errs = append(errs, n.Elements[index].errorsWithPath(path)...)
	}

	return errs
}

func (n *Node) Add(path []any, err error) {
	if len(path) == 0 {
		// If err is a *Node then merge the errors. Use type assertion instead of errors.As because we don't want to reach
		// into the error chain.
		if nodeErr, ok := err.(*Node); ok {
			allErrors := nodeErr.AllErrors()
			for _, errWithPath := range allErrors {
				n.Add(errWithPath.Path, errWithPath.Err)
			}
		} else {
			n.Errs = append(n.Errs, err)
		}
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

type ErrorWithPath struct {
	Path []any
	Err  error
}

func (e *ErrorWithPath) Error() string {
	sb := &strings.Builder{}
	for _, step := range e.Path {
		switch step := step.(type) {
		case string:
			sb.WriteByte('.')
			sb.WriteString(step)
		case int:
			sb.WriteByte('[')
			sb.WriteString(strconv.FormatInt(int64(step), 10))
			sb.WriteByte(']')
		default:
			panic("path elements must be string or int")
		}
	}

	sb.WriteString(": ")
	sb.WriteString(e.Err.Error())

	return sb.String()
}
