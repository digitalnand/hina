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
	Value Object
}

type BinaryTerm struct {
	Lhs Object
	Op  string
	Rhs Object
}

type LetTerm struct {
	Identifier string
	Value      Object
	Next       Object
}

type VarTerm struct {
	Text string
}

type PrintTerm struct {
	Value Object
}

type IfTerm struct {
	Condition Object
	Then      Object
	Else      Object
}

type FunctionTerm struct {
	Parameters []interface{}
	Value      Object
	Env        Environment
}

func (node FunctionTerm) String() string {
	return "<#closure>"
}

type CallTerm struct {
	Arguments []interface{}
	Callee    Object
}
