package hina

import (
	"fmt"
	"strconv"
)

func inspectIf(node map[string]interface{}) (IfNode, error) {
	conditionNode, hasCondition := node["condition"].(map[string]interface{})
	thenNode, hasThen := node["then"].(map[string]interface{})
	elseNode, hasElse := node["otherwise"].(map[string]interface{})
	if !hasCondition || !hasThen || !hasElse {
		return IfNode{}, fmt.Errorf("'If' node is badly structured")
	}
	return IfNode{Condition: conditionNode, Then: thenNode, Else: elseNode}, nil
}

func inspectTupleFunction(node map[string]interface{}) (TupleFunction, error) {
	kind, hasKind := node["kind"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return TupleFunction{}, fmt.Errorf("'%s' node is badly structured", kind)
	}
	return TupleFunction{Kind: kind, Value: valueNode}, nil
}

func inspectTuple(node map[string]interface{}) (TupleNode, error) {
	firstNode, hasFirst := node["first"].(map[string]interface{})
	secondNode, hasSecond := node["second"].(map[string]interface{})
	if !hasFirst || !hasSecond {
		return TupleNode{}, fmt.Errorf("'Tuple' node is badly structured")
	}
	return TupleNode{First: firstNode, Second: secondNode}, nil
}

func inspectVar(node map[string]interface{}) (VarNode, error) {
	text, hasText := node["text"].(string)
	if !hasText {
		return VarNode{}, fmt.Errorf("'Var' node is badly structured")
	}
	return VarNode{Text: text}, nil
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

func inspectBinary(node map[string]interface{}) (BinaryNode, error) {
	op, hasOp := node["op"].(string)
	lhsNode, hasLhs := node["lhs"].(map[string]interface{})
	rhsNode, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryNode{}, fmt.Errorf("'Binary' node is badly structured")
	}
	return BinaryNode{Lhs: lhsNode, Op: op, Rhs: rhsNode}, nil
}

func inspectPrint(node map[string]interface{}) (PrintNode, error) {
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasValue {
		return PrintNode{}, fmt.Errorf("'Print' node is badly structured")
	}
	return PrintNode{Value: valueNode}, nil
}

func inspectLiteral(node map[string]interface{}) (any, error) {
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
