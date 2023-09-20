package hina

import (
	"fmt"
	"strconv"
)

// TODO: idk if a global symbolTable is a good thing, so improve this
var (
	symbolTable = make(map[string]interface{})
)

func inspectNode(node map[string]interface{}) (any, error) {
	kind := node["kind"]

	switch kind {
	case "Str", "Int", "Bool":
		valueNode, err := inspectValue(node)
		if err != nil {
			return nil, err
		}
		return valueNode, nil
	case "Print":
		printNode, err := inspectPrint(node)
		if err != nil {
			return nil, err
		}
		return printNode.Evaluate(), nil
	case "Binary":
		binaryNode, inspectErr := inspectBinary(node)
		if inspectErr != nil {
			return nil, inspectErr
		}

		resultNode, resultErr := binaryNode.Evaluate()
		if resultErr != nil {
			return nil, resultErr
		}
		return resultNode, nil
	case "Let":
		letNode, letErr := inspectLet(node)
		if letErr != nil {
			return nil, letErr
		}
		symbolTable[letNode.Identifier] = letNode

		_, nextErr := inspectNode(letNode.Next.(map[string]interface{}))
		if nextErr != nil {
			return nil, nextErr
		}
		return letNode, nil
	case "Var":
		varNode, varErr := inspectVar(node)
		if varErr != nil {
			return nil, varErr
		}

		letNode, hasLet := symbolTable[varNode.Text].(LetNode)
		if !hasLet {
			return nil, fmt.Errorf("calling an undeclared variable: %s", varNode.Text)
		}

		value, valueErr := inspectNode(letNode.Value.(map[string]interface{}))
		if valueErr != nil {
			return nil, valueErr
		}
		return value, nil
	case "Tuple":
		tupleNode, err := inspectTuple(node)
		if err != nil {
			return nil, err
		}
		return tupleNode, nil
	case "First", "Second":
		tupleFunction, err := inspectTupleFunction(node)
		if err != nil {
			return nil, err
		}
		return tupleFunction, nil
	case "If":
		ifNode, ifErr := inspectIf(node)
		if ifErr != nil {
			return nil, ifErr
		}

		resultNode := ifNode.Evaluate().(map[string]interface{})
		result, resultErr := inspectNode(resultNode)
		if resultErr != nil {
			return nil, resultErr
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown term: %s", kind)
	}
}

func inspectIf(node map[string]interface{}) (IfNode, error) {
	conditionNode, hasCondition := node["condition"].(map[string]interface{})
	thenNode, hasThen := node["then"].(map[string]interface{})
	elseNode, hasElse := node["otherwise"].(map[string]interface{})
	if !hasCondition || !hasThen || !hasElse {
		return IfNode{}, fmt.Errorf("'If' node is badly structured")
	}

	condition, conditionErr := inspectNode(conditionNode)
	if conditionErr != nil {
		return IfNode{}, conditionErr
	}
	boolCondition, isBool := condition.(BoolNode)
	if !isBool {
		return IfNode{}, fmt.Errorf("'If' only accepts Bools as condition")
	}

	return IfNode{Condition: boolCondition, Then: thenNode, Else: elseNode}, nil
}

func inspectTupleFunction(node map[string]interface{}) (any, error) {
	kind, hasKind := node["kind"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return nil, fmt.Errorf("'%s' node is badly structured", kind)
	}

	value, err := inspectNode(valueNode)
	if err != nil {
		return nil, err
	}
	tuple, isTuple := value.(TupleNode)
	if !isTuple {
		return nil, fmt.Errorf("'%s' only accepts Tuples", kind)
	}

	var result any
	switch kind {
	case "First":
		result = tuple.First
	case "Second":
		result = tuple.Second
	}
	return result, nil
}

func inspectTuple(node map[string]interface{}) (TupleNode, error) {
	firstNode, hasFirst := node["first"].(map[string]interface{})
	secondNode, hasSecond := node["second"].(map[string]interface{})
	if !hasFirst || !hasSecond {
		return TupleNode{}, fmt.Errorf("'Tuple' node is badly structured")
	}

	first, firstErr := inspectNode(firstNode)
	if firstErr != nil {
		return TupleNode{}, nil
	}
	second, secondErr := inspectNode(secondNode)
	if secondErr != nil {
		return TupleNode{}, nil
	}

	return TupleNode{First: first, Second: second}, nil
}

func inspectLet(node map[string]interface{}) (LetNode, error) {
	nameNode, hasName := node["name"].(map[string]interface{})
	identifier, hasIdentifier := nameNode["text"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	nextNode, hasNext := node["next"].(map[string]interface{})
	if !hasName || !hasValue || !hasIdentifier || !hasNext {
		return LetNode{}, fmt.Errorf("'Let' node is badly structured")
	}

	return LetNode{Identifier: identifier, Value: valueNode, Next: nextNode}, nil
}

func inspectVar(node map[string]interface{}) (VarNode, error) {
	text, hasText := node["text"].(string)
	if !hasText {
		return VarNode{}, fmt.Errorf("'Var' node is badly structured")
	}
	return VarNode{Text: text}, nil
}

func inspectBinary(node map[string]interface{}) (BinaryNode, error) {
	op, hasOp := node["op"]
	lhsNode, hasLhs := node["lhs"].(map[string]interface{})
	rhsNode, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryNode{}, fmt.Errorf("'Binary' node is badly structured")
	}

	lhs, lhsErr := inspectNode(lhsNode)
	if lhsErr != nil {
		return BinaryNode{}, lhsErr
	}
	rhs, rhsErr := inspectNode(rhsNode)
	if rhsErr != nil {
		return BinaryNode{}, rhsErr
	}

	return BinaryNode{Lhs: lhs, Op: fmt.Sprint(op), Rhs: rhs}, nil
}

func inspectPrint(node map[string]interface{}) (PrintNode, error) {
	value, hasValue := node["value"].(map[string]interface{})
	if !hasValue {
		return PrintNode{}, fmt.Errorf("'Print' node is badly structured")
	}

	termNode, err := inspectNode(value)
	if err != nil {
		return PrintNode{}, err
	}
	return PrintNode{Value: termNode}, nil
}

func inspectValue(node map[string]interface{}) (any, error) {
	kind, hasKind := node["kind"].(string)
	value, hasValue := node["value"]
	if !hasKind || !hasValue {
		return nil, fmt.Errorf("'%s' node is badly structured", kind)
	}

	valueStr := fmt.Sprint(value)
	var result any
	switch kind {
	case "Str":
		result = StrNode{Value: valueStr}
	case "Int":
		num, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, err
		}
		result = IntNode{Value: int32(num)}
	case "Bool":
		boolValue, err := strconv.ParseBool(valueStr)
		if err != nil {
			return nil, err
		}
		result = BoolNode{Value: boolValue}
	}
	return result, nil
}

func getExpression(tree map[string]interface{}) (map[string]interface{}, error) {
	expression, hasExpression := tree["expression"].(map[string]interface{})
	if !hasExpression || len(expression) < 1 {
		return nil, fmt.Errorf("tree has no expressions")
	}
	return expression, nil
}

func WalkTree(tree map[string]interface{}) error {
	expression, expressionErr := getExpression(tree)
	if expressionErr != nil {
		return expressionErr
	}

	_, termErr := inspectNode(expression)
	if termErr != nil {
		return termErr
	}
	return nil
}
