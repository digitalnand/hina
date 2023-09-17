package hina

import (
	"fmt"
	"strconv"
)

func inspectTerm(node map[string]interface{}, nodeValue *any) {
	kind := node["kind"]
	value := fmt.Sprint(node["value"])
	switch kind {
	case "Str":
		var strNode StrNode
		strNode.Value = value
		*nodeValue = strNode
	case "Int":
		var intNode IntNode
		num, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		intNode.Value = int32(num)
		*nodeValue = intNode
	case "Bool":
		var boolNode BoolNode
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		boolNode.Value = boolValue
		*nodeValue = boolNode
	case "Print":
		printNode := inspectPrint(node)
		*nodeValue = printNode
	}
}

func inspectPrint(node map[string]interface{}) PrintNode {
	var printNode PrintNode
	value := node["value"]
	inspectTerm(value.(map[string]interface{}), &printNode.Value)
	return printNode
}

func WalkTree(tree map[string]interface{}) {
	expressionNode, hasExpressionNode := tree["expression"]
	expression, hasExpression := expressionNode.(map[string]interface{})
	if !hasExpressionNode || !hasExpression || len(expression) < 1 {
		panic("tree has no expressions")
	}

	if expression["kind"] == "Print" {
		inspectPrint(expression).Evaluate()
	}
}
