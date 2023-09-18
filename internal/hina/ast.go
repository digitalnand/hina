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

type TupleNode struct {
	First  any
	Second any
}

func (node TupleNode) String() string {
	return fmt.Sprintf("(%s, %s)", node.First, node.Second)
}

type BinaryNode struct {
	Lhs any
	Op  string
	Rhs any
}

type LetNode struct {
	Identifier string
	Value      any
	Next       any
}

type VarNode struct {
	Text string
}

type PrintNode struct {
	Value any
}
