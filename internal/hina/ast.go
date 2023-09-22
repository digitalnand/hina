package hina

import "fmt"

type Term interface{}
type Object map[string]interface{}

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
	First  Term
	Second Term
}

func (node TupleNode) String() string {
	return fmt.Sprintf("(%s, %s)", node.First, node.Second)
}

type TupleFunction struct {
	Kind  string
	Value Term
}

type BinaryNode struct {
	Lhs Term
	Op  string
	Rhs Term
}

type LetNode struct {
	Identifier string
	Value      Term
	Next       Term
}

type VarNode struct {
	Text string
}

type PrintNode struct {
	Value Term
}

type IfNode struct {
	Condition Term
	Then      Term
	Else      Term
}

type FunctionNode struct {
	Parameters []interface{}
	Value      Term
	Env        Environment
}

func (node FunctionNode) String() string {
	return "<#closure>"
}

type CallNode struct {
	Arguments []interface{}
	Callee    Term
}
