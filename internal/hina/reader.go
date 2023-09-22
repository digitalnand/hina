package hina

import (
	"fmt"
	"strconv"
)

func inspectCall(node map[string]interface{}) (CallNode, error) {
	arguments, hasArguments := node["arguments"].([]interface{})
	callee, hasCallee := node["callee"].(map[string]interface{})
	if !hasArguments || !hasCallee {
		return CallNode{}, fmt.Errorf("'Function' node is badly structured")
	}
	return CallNode{Arguments: arguments, Callee: callee}, nil
}

func inspectFunction(node map[string]interface{}) (FunctionNode, error) {
	parameters, hasParameters := node["parameters"].([]interface{})
	value, hasValue := node["value"].(map[string]interface{})
	if !hasParameters || !hasValue {
		return FunctionNode{}, fmt.Errorf("'Function' node is badly structured")
	}
	return FunctionNode{Parameters: parameters, Value: value, Env: NewEnvironment()}, nil
}

func inspectIf(node map[string]interface{}) (IfNode, error) {
	condition, hasCondition := node["condition"].(map[string]interface{})
	then, hasThen := node["then"].(map[string]interface{})
	elseNode, hasElse := node["otherwise"].(map[string]interface{})
	if !hasCondition || !hasThen || !hasElse {
		return IfNode{}, fmt.Errorf("'If' node is badly structured")
	}
	return IfNode{Condition: condition, Then: then, Else: elseNode}, nil
}

func inspectTupleFunction(node map[string]interface{}) (TupleFunction, error) {
	kind, hasKind := node["kind"].(string)
	value, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return TupleFunction{}, fmt.Errorf("'%s' node is badly structured", kind)
	}
	return TupleFunction{Kind: kind, Value: value}, nil
}

func inspectTuple(node map[string]interface{}) (TupleNode, error) {
	first, hasFirst := node["first"].(map[string]interface{})
	second, hasSecond := node["second"].(map[string]interface{})
	if !hasFirst || !hasSecond {
		return TupleNode{}, fmt.Errorf("'Tuple' node is badly structured")
	}
	return TupleNode{First: first, Second: second}, nil
}

func inspectVar(node map[string]interface{}) (VarNode, error) {
	text, hasText := node["text"].(string)
	if !hasText {
		return VarNode{}, fmt.Errorf("'Var' node is badly structured")
	}
	return VarNode{Text: text}, nil
}

func inspectLet(node map[string]interface{}) (LetNode, error) {
	name, hasName := node["name"].(map[string]interface{})
	identifier, hasIdentifier := name["text"].(string)
	value, hasValue := node["value"].(map[string]interface{})
	next, hasNext := node["next"].(map[string]interface{})
	if !hasName || !hasValue || !hasIdentifier || !hasNext {
		return LetNode{}, fmt.Errorf("'Let' node is badly structured")
	}

	return LetNode{Identifier: identifier, Value: value, Next: next}, nil
}

func inspectBinary(node map[string]interface{}) (BinaryNode, error) {
	op, hasOp := node["op"].(string)
	lhs, hasLhs := node["lhs"].(map[string]interface{})
	rhs, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryNode{}, fmt.Errorf("'Binary' node is badly structured")
	}
	return BinaryNode{Lhs: lhs, Op: op, Rhs: rhs}, nil
}

func inspectPrint(node map[string]interface{}) (PrintNode, error) {
	value, hasValue := node["value"].(map[string]interface{})
	if !hasValue {
		return PrintNode{}, fmt.Errorf("'Print' node is badly structured")
	}
	return PrintNode{Value: value}, nil
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
