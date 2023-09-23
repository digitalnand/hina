package hina

import (
	"fmt"
	"strconv"
)

func inspectCall(node Object) (CallTerm, error) {
	arguments, hasArguments := node["arguments"].([]interface{})
	callee, hasCallee := node["callee"].(map[string]interface{})
	if !hasArguments || !hasCallee {
		return CallTerm{}, fmt.Errorf("'Function' node is badly structured")
	}
	return CallTerm{Arguments: arguments, Callee: callee}, nil
}

func inspectFunction(node Object) (FunctionTerm, error) {
	parameters, hasParameters := node["parameters"].([]interface{})
	value, hasValue := node["value"].(map[string]interface{})
	if !hasParameters || !hasValue {
		return FunctionTerm{}, fmt.Errorf("'Function' node is badly structured")
	}
	return FunctionTerm{Parameters: parameters, Value: value, Env: Environment{}}, nil
}

func inspectIf(node Object) (IfTerm, error) {
	condition, hasCondition := node["condition"].(map[string]interface{})
	then, hasThen := node["then"].(map[string]interface{})
	elseTerm, hasElse := node["otherwise"].(map[string]interface{})
	if !hasCondition || !hasThen || !hasElse {
		return IfTerm{}, fmt.Errorf("'If' node is badly structured")
	}
	return IfTerm{Condition: condition, Then: then, Else: elseTerm}, nil
}

func inspectTupleFunction(node Object) (TupleFunction, error) {
	kind, hasKind := node["kind"].(string)
	value, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return TupleFunction{}, fmt.Errorf("'%s' node is badly structured", kind)
	}
	return TupleFunction{Kind: kind, Value: value}, nil
}

func inspectTuple(node Object) (TupleTerm, error) {
	first, hasFirst := node["first"].(map[string]interface{})
	second, hasSecond := node["second"].(map[string]interface{})
	if !hasFirst || !hasSecond {
		return TupleTerm{}, fmt.Errorf("'Tuple' node is badly structured")
	}
	return TupleTerm{First: first, Second: second}, nil
}

func inspectVar(node Object) (VarTerm, error) {
	text, hasText := node["text"].(string)
	if !hasText {
		return VarTerm{}, fmt.Errorf("'Var' node is badly structured")
	}
	return VarTerm{Text: text}, nil
}

func inspectLet(node Object) (LetTerm, error) {
	name, hasName := node["name"].(map[string]interface{})
	identifier, hasIdentifier := name["text"].(string)
	value, hasValue := node["value"].(map[string]interface{})
	next, hasNext := node["next"].(map[string]interface{})
	if !hasName || !hasValue || !hasIdentifier || !hasNext {
		return LetTerm{}, fmt.Errorf("'Let' node is badly structured")
	}

	return LetTerm{Identifier: identifier, Value: value, Next: next}, nil
}

func inspectBinary(node Object) (BinaryTerm, error) {
	op, hasOp := node["op"].(string)
	lhs, hasLhs := node["lhs"].(map[string]interface{})
	rhs, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryTerm{}, fmt.Errorf("'Binary' node is badly structured")
	}
	return BinaryTerm{Lhs: lhs, Op: op, Rhs: rhs}, nil
}

func inspectPrint(node Object) (PrintTerm, error) {
	value, hasValue := node["value"].(map[string]interface{})
	if !hasValue {
		return PrintTerm{}, fmt.Errorf("'Print' node is badly structured")
	}
	return PrintTerm{Value: value}, nil
}

func inspectLiteral(node Object) (Term, error) {
	kind, hasKind := node["kind"].(string)
	value, hasValue := node["value"]
	if !hasKind || !hasValue {
		return nil, fmt.Errorf("'%s' node is badly structured", kind)
	}

	valueStr := fmt.Sprint(value)
	var result Term
	switch kind {
	case "Str":
		result = StrTerm{Value: valueStr}
	case "Int":
		num, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, err
		}
		result = IntTerm{Value: int32(num)}
	case "Bool":
		boolValue, err := strconv.ParseBool(valueStr)
		if err != nil {
			return nil, err
		}
		result = BoolTerm{Value: boolValue}
	}
	return result, nil
}
