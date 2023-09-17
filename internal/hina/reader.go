package hina

import (
	"fmt"
	"strconv"
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
	default:
		return nil, fmt.Errorf("unknown term: %s", kind)
	}
}

func inspectBinary(node map[string]interface{}) (BinaryNode, error) {
	op, hasOp := node["op"]
	lhsNode, hasLhs := node["lhs"].(map[string]interface{})
	rhsNode, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryNode{}, fmt.Errorf("binary node is badly structured")
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
	value := node["value"].(map[string]interface{})
	termNode, err := inspectTerm(value)
	if err != nil {
		return PrintNode{}, err
	}
	return PrintNode{Value: termNode}, nil
}

func getExpression(tree map[string]interface{}) (map[string]interface{}, bool) {
	expression, ok := tree["expression"].(map[string]interface{})
	if !ok || len(expression) < 1 {
		return nil, false
	}
	return expression, true
}

func WalkTree(tree map[string]interface{}) error {
	expression, ok := getExpression(tree)
	if !ok {
		return fmt.Errorf("tree has no expressions")
	}
	_, err := inspectTerm(expression)
	if err != nil {
		return err
	}
	return nil
}
