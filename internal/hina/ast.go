package hina

import "fmt"

type Term interface{}
type Object map[string]interface{}

type StrTerm struct {
	Value string
}

func (node StrTerm) String() string {
	return node.Value
}

type IntTerm struct {
	Value int32
}

func (node IntTerm) String() string {
	return fmt.Sprintf("%d", node.Value)
}

type BoolTerm struct {
	Value bool
}

func (node BoolTerm) String() string {
	return fmt.Sprintf("%t", node.Value)
}

type TupleTerm struct {
	First  Term
	Second Term
}

func (node TupleTerm) String() string {
	return fmt.Sprintf("(%s, %s)", node.First, node.Second)
}

type TupleFunction struct {
	Kind  string
	Value Term
}

type BinaryTerm struct {
	Lhs Term
	Op  string
	Rhs Term
}

type LetTerm struct {
	Identifier string
	Value      Term
	Next       Term
}

type VarTerm struct {
	Text string
}

type PrintTerm struct {
	Value Term
}

type IfTerm struct {
	Condition Term
	Then      Term
	Else      Term
}

type FunctionTerm struct {
	Parameters []string
	Value      Term
	Env        Environment
}

func (node FunctionTerm) String() string {
	return "<#closure>"
}

type CallTerm struct {
	FunctionCalled string
	Arguments      []Term
	Callee         Term
}
