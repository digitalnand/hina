package hina

import (
	"fmt"
	"strconv"
)

func InspectNode(node Object) (Term, error) {
	kind := node["kind"]
	switch kind {
	case "Str", "Int", "Bool":
		literal, err := inspectLiteral(node)
		if err != nil {
			return nil, err
		}
		return literal, nil
	case "Print":
		print, err := inspectPrint(node)
		if err != nil {
			return nil, err
		}
		return print, nil
	case "Binary":
		binary, err := inspectBinary(node)
		if err != nil {
			return nil, err
		}
		return binary, nil
	case "Let":
		let, err := inspectLet(node)
		if err != nil {
			return nil, err
		}
		return let, nil
	case "Var":
		varTerm, err := inspectVar(node)
		if err != nil {
			return nil, err
		}
		return varTerm, nil
	case "Tuple":
		tuple, err := inspectTuple(node)
		if err != nil {
			return nil, err
		}
		return tuple, nil
	case "First", "Second":
		tupleFunc, err := inspectTupleFunction(node)
		if err != nil {
			return nil, err
		}
		return tupleFunc, nil
	case "If":
		ifTerm, err := inspectIf(node)
		if err != nil {
			return nil, err
		}
		return ifTerm, nil
	case "Function":
		function, err := inspectFunction(node)
		if err != nil {
			return nil, err
		}
		return function, nil
	case "Call":
		call, err := inspectCall(node)
		if err != nil {
			return nil, err
		}
		return call, nil
	default:
		return nil, fmt.Errorf("unknown node: %s", kind)
	}
}

func inspectCall(node Object) (CallTerm, error) {
	argumentsNode, hasArguments := node["arguments"].([]interface{})
	calleeNode, hasCallee := node["callee"].(map[string]interface{})
	if !hasArguments || !hasCallee {
		return CallTerm{}, fmt.Errorf("'Function' node is badly structured")
	}

	arguments, argumentsInspectErr := inspectCallArguments(argumentsNode)
	if argumentsInspectErr != nil {
		return CallTerm{}, argumentsInspectErr
	}
	callee, calleeInspectErr := InspectNode(calleeNode)
	if calleeInspectErr != nil {
		return CallTerm{}, calleeInspectErr
	}
	return CallTerm{Arguments: arguments, Callee: callee}, nil
}

func inspectCallArguments(argumentsNode []interface{}) ([]Term, error) {
	var arguments []Term
	for index, argument := range argumentsNode {
		argumentNode, isNode := argument.(map[string]interface{})
		if !isNode {
			return nil, fmt.Errorf("argument in index %d isn't a Node", index)
		}
		term, err := InspectNode(argumentNode)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, term)
	}
	return arguments, nil
}

func inspectFunction(node Object) (FunctionTerm, error) {
	parametersNode, hasParameters := node["parameters"].([]interface{})
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasParameters || !hasValue {
		return FunctionTerm{}, fmt.Errorf("'Function' node is badly structured")
	}

	parameters, parametersInspectErr := inspectFunctionParameters(parametersNode)
	if parametersInspectErr != nil {
		return FunctionTerm{}, parametersInspectErr
	}
	value, valueInspectErr := InspectNode(valueNode)
	if valueInspectErr != nil {
		return FunctionTerm{}, valueInspectErr
	}
	return FunctionTerm{Parameters: parameters, Value: value, Env: Environment{}}, nil
}

func inspectFunctionParameters(parametersNode []interface{}) ([]string, error) {
	var parameters []string
	for index, parameter := range parametersNode {
		parameterNode, isNode := parameter.(map[string]interface{})
		if !isNode {
			return nil, fmt.Errorf("parameter in index %d isn't a Node", index)
		}
		parameterName, exists := parameterNode["text"].(string)
		if !exists {
			return nil, fmt.Errorf("malformed parameter in index %d", index)
		}
		parameters = append(parameters, parameterName)
	}
	return parameters, nil
}

func inspectIf(node Object) (IfTerm, error) {
	conditionNode, hasCondition := node["condition"].(map[string]interface{})
	thenNode, hasThen := node["then"].(map[string]interface{})
	elseNode, hasElse := node["otherwise"].(map[string]interface{})
	if !hasCondition || !hasThen || !hasElse {
		return IfTerm{}, fmt.Errorf("'If' node is badly structured")
	}

	condition, conditionInspectErr := InspectNode(conditionNode)
	if conditionInspectErr != nil {
		return IfTerm{}, conditionInspectErr
	}
	then, thenInspectErr := InspectNode(thenNode)
	if thenInspectErr != nil {
		return IfTerm{}, thenInspectErr
	}
	elseTerm, elseInspectErr := InspectNode(elseNode)
	if elseInspectErr != nil {
		return IfTerm{}, elseInspectErr
	}
	return IfTerm{Condition: condition, Then: then, Else: elseTerm}, nil
}

func inspectTupleFunction(node Object) (TupleFunction, error) {
	kind, hasKind := node["kind"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	if !hasValue || !hasKind {
		return TupleFunction{}, fmt.Errorf("'%s' node is badly structured", kind)
	}
	value, err := InspectNode(valueNode)
	if err != nil {
		return TupleFunction{}, err
	}
	return TupleFunction{Kind: kind, Value: value}, nil
}

func inspectTuple(node Object) (TupleTerm, error) {
	firstNode, hasFirst := node["first"].(map[string]interface{})
	secondNode, hasSecond := node["second"].(map[string]interface{})
	if !hasFirst || !hasSecond {
		return TupleTerm{}, fmt.Errorf("'Tuple' node is badly structured")
	}

	first, firstInspectErr := InspectNode(firstNode)
	if firstInspectErr != nil {
		return TupleTerm{}, firstInspectErr
	}
	second, secondInspectErr := InspectNode(secondNode)
	if secondInspectErr != nil {
		return TupleTerm{}, secondInspectErr
	}
	return TupleTerm{First: first, Second: second}, nil
}

func inspectVar(node Object) (VarTerm, error) {
	text, exists := node["text"].(string)
	if !exists {
		return VarTerm{}, fmt.Errorf("'Var' node is badly structured")
	}
	return VarTerm{Text: text}, nil
}

func inspectLet(node Object) (LetTerm, error) {
	name, hasName := node["name"].(map[string]interface{})
	identifier, hasIdentifier := name["text"].(string)
	valueNode, hasValue := node["value"].(map[string]interface{})
	nextNode, hasNext := node["next"].(map[string]interface{})
	if !hasName || !hasValue || !hasIdentifier || !hasNext {
		return LetTerm{}, fmt.Errorf("'Let' node is badly structured")
	}

	value, valueInspectErr := InspectNode(valueNode)
	if valueInspectErr != nil {
		return LetTerm{}, valueInspectErr
	}
	next, nextInspectErr := InspectNode(nextNode)
	if nextInspectErr != nil {
		return LetTerm{}, nextInspectErr
	}
	return LetTerm{Identifier: identifier, Value: value, Next: next}, nil
}

func inspectBinary(node Object) (BinaryTerm, error) {
	op, hasOp := node["op"].(string)
	lhsNode, hasLhs := node["lhs"].(map[string]interface{})
	rhsNode, hasRhs := node["rhs"].(map[string]interface{})
	if !hasOp || !hasLhs || !hasRhs {
		return BinaryTerm{}, fmt.Errorf("'Binary' node is badly structured")
	}

	lhs, lhsInspectErr := InspectNode(lhsNode)
	if lhsInspectErr != nil {
		return BinaryTerm{}, lhsInspectErr
	}
	rhs, rhsInspectErr := InspectNode(rhsNode)
	if lhsInspectErr != nil {
		return BinaryTerm{}, rhsInspectErr
	}
	return BinaryTerm{Lhs: lhs, Op: op, Rhs: rhs}, nil
}

func inspectPrint(node Object) (PrintTerm, error) {
	valueNode, exists := node["value"].(map[string]interface{})
	if !exists {
		return PrintTerm{}, fmt.Errorf("'Print' node is badly structured")
	}
	value, err := InspectNode(valueNode)
	if err != nil {
		return PrintTerm{}, err
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
