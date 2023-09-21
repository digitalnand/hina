package hina

import (
	"fmt"
	"reflect"
)

func EvalTree(tree map[string]interface{}, env Environment) error {
	expression, expressionErr := getExpression(tree)
	if expressionErr != nil {
		return expressionErr
	}

	_, termErr := evalNode(expression, env)
	if termErr != nil {
		return termErr
	}
	return nil
}

func getExpression(tree map[string]interface{}) (map[string]interface{}, error) {
	expression, hasExpression := tree["expression"].(map[string]interface{})
	if !hasExpression || len(expression) < 1 {
		return nil, fmt.Errorf("tree has no expressions")
	}
	return expression, nil
}

func evalNode(node map[string]interface{}, env Environment) (any, error) {
	kind := node["kind"]

	switch kind {
	case "Str", "Int", "Bool":
		literalNode, err := inspectLiteral(node)
		if err != nil {
			return nil, err
		}
		return literalNode, nil
	case "Print":
		printNode, printErr := inspectPrint(node)
		if printErr != nil {
			return nil, printErr
		}
		resultNode, resultErr := printNode.Evaluate(env)
		if resultErr != nil {
			return nil, resultErr
		}
		return resultNode, nil
	case "Binary":
		binaryNode, inspectErr := inspectBinary(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		resultNode, resultErr := binaryNode.Evaluate(env)
		if resultErr != nil {
			return nil, resultErr
		}
		return resultNode, nil
	case "Let":
		letNode, letErr := inspectLet(node)
		if letErr != nil {
			return nil, letErr
		}
		resultErr := letNode.Evaluate(env)
		if resultErr != nil {
			return nil, resultErr
		}
		return letNode, nil
	case "Var":
		varNode, varErr := inspectVar(node)
		if varErr != nil {
			return nil, varErr
		}
		value, valueErr := varNode.Evaluate(env)
		if valueErr != nil {
			return nil, valueErr
		}
		return value, nil
	case "Tuple":
		tupleNode, err := inspectTuple(node)
		if err != nil {
			return nil, err
		}
		tupleNode, err = tupleNode.Evaluate(env)
		if err != nil {
			return nil, err
		}
		return tupleNode, nil
	case "First", "Second":
		node, nodeErr := inspectTupleFunction(node)
		if nodeErr != nil {
			return nil, nodeErr
		}
		value, valueErr := node.Evaluate(env)
		if valueErr != nil {
			return nil, valueErr
		}
		return value, nil
	case "If":
		ifNode, ifErr := inspectIf(node)
		if ifErr != nil {
			return nil, ifErr
		}
		resultNode, nodeErr := ifNode.Evaluate(env)
		if nodeErr != nil {
			return nil, nodeErr
		}
		result, resultErr := evalNode(resultNode.(map[string]interface{}), env)
		if resultErr != nil {
			return nil, resultErr
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown term: %s", kind)
	}
}

func (print PrintNode) Evaluate(env Environment) (any, error) {
	value, valueErr := evalNode(print.Value.(map[string]interface{}), env)
	if valueErr != nil {
		return nil, valueErr
	}

	fmt.Println(value)
	return value, nil
}

func (binary BinaryNode) Evaluate(env Environment) (any, error) {
	lhs, lhsErr := evalNode(binary.Lhs.(map[string]interface{}), env)
	if lhsErr != nil {
		return nil, lhsErr
	}
	rhs, rhsErr := evalNode(binary.Rhs.(map[string]interface{}), env)
	if rhsErr != nil {
		return nil, rhsErr
	}

	switch binary.Op {
	case "Add":
		// TODO: improve this
		_, isLhsString := lhs.(StrNode)
		_, isRhsString := rhs.(StrNode)
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		if isLhsInt && isRhsInt {
			return IntNode{Value: intLhs.Value + intRhs.Value}, nil
		}
		if (isLhsString || isLhsInt) && (isRhsInt || isRhsString) {
			return StrNode{Value: fmt.Sprintf("%s%s", lhs, rhs)}, nil
		}
		return nil, fmt.Errorf("'Add' operator can only be used with Ints and/or Strs")
	case "Sub", "Mul", "Div", "Rem":
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		var result int32
		if isLhsInt && isRhsInt {
			switch binary.Op {
			case "Sub":
				result = intLhs.Value - intRhs.Value
			case "Mul":
				result = intLhs.Value * intRhs.Value
			case "Div":
				result = intLhs.Value / intRhs.Value
			case "Rem":
				result = intLhs.Value % intRhs.Value
			}
			return IntNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' operator can only be used with Ints", binary.Op)
	case "Eq", "Neq":
		hasSameValue := lhs == rhs
		hasSameType := reflect.TypeOf(lhs) == reflect.TypeOf(rhs)
		var result bool
		switch binary.Op {
		case "Eq":
			result = hasSameValue && hasSameType
		case "Neq":
			result = !hasSameValue || !hasSameType
		}
		return BoolNode{Value: result}, nil
	case "Lt", "Gt", "Lte", "Gte":
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		var result bool
		if isLhsInt && isRhsInt {
			switch binary.Op {
			case "Lt":
				result = intLhs.Value < intRhs.Value
			case "Gt":
				result = intLhs.Value > intRhs.Value
			case "Lte":
				result = intLhs.Value <= intRhs.Value
			case "Gte":
				result = intLhs.Value >= intRhs.Value
			}
			return BoolNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' comparison can only be done with Ints", binary.Op)
	case "And", "Or":
		boolLhs, isLhsBool := lhs.(BoolNode)
		boolRhs, isRhsBool := rhs.(BoolNode)
		var result bool
		if isLhsBool && isRhsBool {
			switch binary.Op {
			case "And":
				result = boolLhs.Value && boolRhs.Value
			case "Or":
				result = boolLhs.Value || boolRhs.Value
			}
			return BoolNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' operator can only be used with Bool", binary.Op)
	default:
		return nil, fmt.Errorf("unknown binary operator: '%s'", binary.Op)
	}
}

func (variable LetNode) Evaluate(env Environment) error {
	env.Set(variable.Identifier, variable)
	_, nextErr := evalNode(variable.Next.(map[string]interface{}), env)
	if nextErr != nil {
		return nextErr
	}
	return nil
}

func (varCall VarNode) Evaluate(env Environment) (any, error) {
	node, hasNode := env.Get(varCall.Text)
	if !hasNode {
		return nil, fmt.Errorf("calling an undeclared variable: %s", varCall.Text)
	}
	variable, isLet := node.(LetNode)
	if !isLet {
		return nil, fmt.Errorf("'Var' can only call Let")
	}

	value, valueErr := evalNode(variable.Value.(map[string]interface{}), env)
	if valueErr != nil {
		return nil, valueErr
	}
	return value, nil
}

func (tuple TupleNode) Evaluate(env Environment) (TupleNode, error) {
	first, firstErr := evalNode(tuple.First.(map[string]interface{}), env)
	if firstErr != nil {
		return TupleNode{}, nil
	}
	second, secondErr := evalNode(tuple.Second.(map[string]interface{}), env)
	if secondErr != nil {
		return TupleNode{}, nil
	}
	return TupleNode{First: first, Second: second}, nil
}

func (tupleFunc TupleFunction) Evaluate(env Environment) (any, error) {
	value, valueErr := evalNode(tupleFunc.Value.(map[string]interface{}), env)
	if valueErr != nil {
		return nil, valueErr
	}
	tuple, isTuple := value.(TupleNode)
	if !isTuple {
		return nil, fmt.Errorf("'%s' only accepts Tuples", tupleFunc.Kind)
	}

	switch tupleFunc.Kind {
	case "First":
		return tuple.First, nil
	case "Second":
		return tuple.Second, nil
	default:
		return nil, fmt.Errorf("'%s' isn't a Tuple function", tupleFunc.Kind)
	}
}

func (ifTerm IfNode) Evaluate(env Environment) (any, error) {
	conditionNode, nodeErr := evalNode(ifTerm.Condition.(map[string]interface{}), env)
	if nodeErr != nil {
		return nil, nodeErr
	}
	condition, isBool := conditionNode.(BoolNode)
	if !isBool {
		return nil, fmt.Errorf("'If' only accepts Bools as condition")
	}

	if condition.Value {
		return ifTerm.Then, nil
	}
	return ifTerm.Else, nil
}
