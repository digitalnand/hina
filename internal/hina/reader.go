package hina

import (
	"fmt"
	"strconv"
)

// TODO: idk if a global symbolTable is a good thing, so improve this
var (
	symbolTable = make(map[string]interface{})
)

func inspectTerm(node map[string]interface{}) (any, error) {
	kind := node["kind"]
	value := fmt.Sprint(node["value"])

	switch kind {
	case "Str":
		return StrNode{Value: value}, nil
	case "Int":
		num, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return IntNode{Value: int32(num)}, nil
	case "Bool":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return BoolNode{Value: boolValue}, nil
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

		_, nextErr := inspectTerm(letNode.Next.(map[string]interface{}))
		if nextErr != nil {
			return nil, nextErr
		}
		return letNode, nil
	case "Var":
		varNode, err := inspectVar(node)
		if err != nil {
			return nil, err
		}

		letNode, hasLet := symbolTable[varNode.Text].(LetNode)
		if !hasLet {
			return nil, fmt.Errorf("calling an undeclared variable: %s", varNode.Text)
		}
		return letNode.Value, nil
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
		ifNode, err := inspectIf(node)
		if err != nil {
			return nil, err
		}
		return ifNode.Evaluate(), nil
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

	condition, conditionErr := inspectTerm(conditionNode)
	if conditionErr != nil {
		return IfNode{}, conditionErr
	}
	conditionBool, isBool := condition.(BoolNode)
	if !isBool {
		return IfNode{}, fmt.Errorf("'If' only accepts Bools as condition")
	}
	then, thenErr := inspectTerm(thenNode)
	if thenErr != nil {
		return IfNode{}, thenErr
	}
	elseTerm, elseErr := inspectTerm(elseNode)
	if elseErr != nil {
		return IfNode{}, elseErr
	}

	return IfNode{Condition: conditionBool, Then: then, Else: elseTerm}, nil
}

func inspectTupleFunction(node map[string]interface{}) (any, error) {
	kind, hasKind := node["kind"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return nil, fmt.Errorf("'%s' node is badly structured", kind)
	}

	value, err := inspectTerm(valueNode)
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

	first, firstErr := inspectTerm(firstNode)
	if firstErr != nil {
		return TupleNode{}, nil
	}
	second, secondErr := inspectTerm(secondNode)
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

	value, err := inspectTerm(valueNode)
	if err != nil {
		return LetNode{}, err
	}

	return LetNode{Identifier: identifier, Value: value, Next: nextNode}, nil
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

	lhs, lhsErr := inspectTerm(lhsNode)
	if lhsErr != nil {
		return BinaryNode{}, lhsErr
	}
	rhs, rhsErr := inspectTerm(rhsNode)
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

	termNode, err := inspectTerm(value)
	if err != nil {
		return PrintNode{}, err
	}
	return PrintNode{Value: termNode}, nil
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

	_, termErr := inspectTerm(expression)
	if termErr != nil {
		return termErr
	}
	return nil
}
