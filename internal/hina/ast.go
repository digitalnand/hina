package hina

import "fmt"

type StrNode struct {
	Value string
}

func (node StrNode) String() string {
	return node.Value
}

type IntNode struct {
	Value int32
}

func (node IntNode) String() string {
	return fmt.Sprintf("%d", node.Value)
}

type BoolNode struct {
	Value bool
}

func (node BoolNode) String() string {
	return fmt.Sprintf("%t", node.Value)
}

type BinaryNode struct {
	Lhs any
	Op  string
	Rhs any
}

type PrintNode struct {
	Value any
}
