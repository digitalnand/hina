package hina

import (
	"fmt"
	"strconv"
)

func inspectTerm(node map[string]interface{}, nodeValue *any) error {
	kind := node["kind"]
	value := fmt.Sprint(node["value"])

	switch kind {
	case "Str":
		*nodeValue = StrNode{Value: value}
	case "Int":
		num, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*nodeValue = IntNode{Value: int32(num)}
	case "Bool":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*nodeValue = BoolNode{Value: boolValue}
	case "Print":
		printNode, err := inspectPrint(node)
		if err != nil {
			return err
		}
		*nodeValue = printNode
	default:
		return fmt.Errorf("unknown term: %s", kind)
	}

	return nil
}

func inspectNode(node map[string]interface{}) error {
	kind := node["kind"]
	switch kind {
	case "Print":
		printNode, err := inspectPrint(node)
		if err != nil {
			return err
		}
		printNode.Evaluate()
	default:
		return fmt.Errorf("unknown node: %s", kind)
	}
	return nil
}

func inspectPrint(node map[string]interface{}) (PrintNode, error) {
	value := node["value"].(map[string]interface{})
	var printNode PrintNode
	err := inspectTerm(value, &printNode.Value)
	if err != nil {
		return PrintNode{}, err
	}
	return printNode, nil
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
	err := inspectNode(expression)
	if err != nil {
		return err
	}
	return nil
}
