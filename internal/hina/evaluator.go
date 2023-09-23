package hina

import (
	"fmt"
	"reflect"
)

func EvalTree(tree Object, env Environment) error {
	expression, exists := tree["expression"].(map[string]interface{})
	if !exists || len(expression) == 0 {
		return fmt.Errorf("tree has no expressions")
	}
	_, err := evalNode(expression, env)
	if err != nil {
		return err
	}
	return nil
}

func evalNode(node Object, env Environment) (Term, error) {
	kind := node["kind"]

	switch kind {
	case "Str", "Int", "Bool":
		literal, err := inspectLiteral(node)
		if err != nil {
			return nil, err
		}
		return literal, nil
	case "Print":
		print, inspectErr := inspectPrint(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		returnTerm, evalErr := print.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return returnTerm, nil
	case "Binary":
		binary, inspectErr := inspectBinary(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		result, evalErr := binary.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return result, nil
	case "Let":
		let, inspectErr := inspectLet(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		nextResult, evalErr := let.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return nextResult, nil
	case "Var":
		varTerm, inspectErr := inspectVar(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		value, evalErr := varTerm.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return value, nil
	case "Tuple":
		tuple, err := inspectTuple(node)
		if err != nil {
			return nil, err
		}
		tuple, err = tuple.Eval(env)
		if err != nil {
			return nil, err
		}
		return tuple, nil
	case "First", "Second":
		tupleFunc, inspectErr := inspectTupleFunction(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		value, evalErr := tupleFunc.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return value, nil
	case "If":
		ifTerm, inspectErr := inspectIf(node)
		if inspectErr != nil {
			return nil, inspectErr
		}
		result, evalErr := ifTerm.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return result, nil
	case "Function":
		function, err := inspectFunction(node)
		if err != nil {
			return nil, err
		}
		return function, nil
	case "Call":
		call, inspectTerm := inspectCall(node)
		if inspectTerm != nil {
			return nil, inspectTerm
		}
		result, evalErr := call.Eval(env)
		if evalErr != nil {
			return nil, evalErr
		}
		return result, nil
	}

	return nil, fmt.Errorf("unknown term: %s", kind)
}

func (print PrintTerm) Eval(env Environment) (Term, error) {
	value, err := evalNode(print.Value, env)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	return value, nil
}

func (binary BinaryTerm) Eval(env Environment) (Term, error) {
	lhs, lhsEvalErr := evalNode(binary.Lhs, env)
	if lhsEvalErr != nil {
		return nil, lhsEvalErr
	}
	rhs, rhsEvalErr := evalNode(binary.Rhs, env)
	if rhsEvalErr != nil {
		return nil, rhsEvalErr
	}

	switch binary.Op {
	case "Add":
		// TODO: improve this
		_, isLhsString := lhs.(StrTerm)
		_, isRhsString := rhs.(StrTerm)
		intLhs, isLhsInt := lhs.(IntTerm)
		intRhs, isRhsInt := rhs.(IntTerm)
		if isLhsString || isRhsString {
			return StrTerm{Value: fmt.Sprintf("%s%s", lhs, rhs)}, nil
		} else if isLhsInt && isRhsInt {
			return IntTerm{Value: intLhs.Value + intRhs.Value}, nil
		}
		return nil, fmt.Errorf("'Add' operator can only be used with Ints and/or Strs")
	case "Sub", "Mul", "Div", "Rem":
		intLhs, isLhsInt := lhs.(IntTerm)
		intRhs, isRhsInt := rhs.(IntTerm)
		if !isLhsInt && !isRhsInt {
			return nil, fmt.Errorf("'%s' operator can only be used with Ints", binary.Op)
		}
		switch binary.Op {
		case "Sub":
			return IntTerm{Value: intLhs.Value - intRhs.Value}, nil
		case "Mul":
			return IntTerm{Value: intLhs.Value * intRhs.Value}, nil
		case "Div":
			return IntTerm{Value: intLhs.Value / intRhs.Value}, nil
		case "Rem":
			return IntTerm{Value: intLhs.Value % intRhs.Value}, nil
		}
	case "Eq", "Neq":
		hasSameValue := lhs == rhs
		hasSameType := reflect.TypeOf(lhs) == reflect.TypeOf(rhs)
		switch binary.Op {
		case "Eq":
			return BoolTerm{Value: hasSameValue && hasSameType}, nil
		case "Neq":
			return BoolTerm{Value: !hasSameValue || !hasSameType}, nil
		}
	case "Lt", "Gt", "Lte", "Gte":
		intLhs, isLhsInt := lhs.(IntTerm)
		intRhs, isRhsInt := rhs.(IntTerm)
		if !isLhsInt && !isRhsInt {
			return nil, fmt.Errorf("'%s' comparison can only be done with Ints", binary.Op)
		}
		switch binary.Op {
		case "Lt":
			return BoolTerm{Value: intLhs.Value < intRhs.Value}, nil
		case "Gt":
			return BoolTerm{Value: intLhs.Value > intRhs.Value}, nil
		case "Lte":
			return BoolTerm{Value: intLhs.Value <= intRhs.Value}, nil
		case "Gte":
			return BoolTerm{Value: intLhs.Value >= intRhs.Value}, nil
		}
	case "And", "Or":
		boolLhs, isLhsBool := lhs.(BoolTerm)
		boolRhs, isRhsBool := rhs.(BoolTerm)
		if !isLhsBool && !isRhsBool {
			return nil, fmt.Errorf("'%s' operator can only be used with Bool", binary.Op)
		}
		switch binary.Op {
		case "And":
			return BoolTerm{Value: boolLhs.Value && boolRhs.Value}, nil
		case "Or":
			return BoolTerm{Value: boolLhs.Value || boolRhs.Value}, nil
		}
	}

	return nil, fmt.Errorf("unknown binary operator: '%s'", binary.Op)
}

func (variable LetTerm) Eval(env Environment) (Term, error) {
	value, valueEvalError := evalNode(variable.Value, env)
	if valueEvalError != nil {
		return nil, valueEvalError
	}
	env.Set(variable.Identifier, value)
	nextResult, nextEvalErr := evalNode(variable.Next, env)
	if nextEvalErr != nil {
		return nil, nextEvalErr
	}
	return nextResult, nil
}

func (varCall VarTerm) Eval(env Environment) (Term, error) {
	value, exists := env.Get(varCall.Text)
	if !exists {
		return nil, fmt.Errorf("calling an undeclared variable: %s", varCall.Text)
	}
	return value, nil
}

func (tuple TupleTerm) Eval(env Environment) (TupleTerm, error) {
	first, firstEvalErr := evalNode(tuple.First.(map[string]interface{}), env)
	if firstEvalErr != nil {
		return TupleTerm{}, nil
	}
	second, secondEvalErr := evalNode(tuple.Second.(map[string]interface{}), env)
	if secondEvalErr != nil {
		return TupleTerm{}, nil
	}
	return TupleTerm{First: first, Second: second}, nil
}

func (tupleFunc TupleFunction) Eval(env Environment) (Term, error) {
	value, err := evalNode(tupleFunc.Value, env)
	if err != nil {
		return nil, err
	}
	tuple, isTuple := value.(TupleTerm)
	if !isTuple {
		return nil, fmt.Errorf("'%s' only accepts Tuples", tupleFunc.Kind)
	}
	if tupleFunc.Kind == "Second" {
		return tuple.Second, nil
	}
	return tuple.First, nil
}

func (ifTerm IfTerm) Eval(env Environment) (Term, error) {
	conditionTerm, conditionEvalErr := evalNode(ifTerm.Condition, env)
	if conditionEvalErr != nil {
		return nil, conditionEvalErr
	}
	condition, isBool := conditionTerm.(BoolTerm)
	if !isBool {
		return nil, fmt.Errorf("'If' only accepts Bools as condition")
	}

	var body Object
	if condition.Value {
		body = ifTerm.Then
	} else {
		body = ifTerm.Else
	}

	result, bodyEvalErr := evalNode(body, env)
	if bodyEvalErr != nil {
		return nil, bodyEvalErr
	}
	return result, nil
}

func (function FunctionTerm) captureEnv(env Environment) {
	for key, value := range env.SymbolTable {
		if _, exists := function.Env.Get(key); exists {
			continue
		}
		function.Env.Set(key, value)
	}
}

func (function FunctionTerm) setParameters(arguments []interface{}, env Environment) error {
	if len(function.Parameters) != len(arguments) {
		return fmt.Errorf("expected %d arguments, received %d", len(function.Parameters), len(arguments))
	}

	for argIndex := 0; argIndex < len(arguments); argIndex++ {
		parameter, hasParameter := function.Parameters[argIndex].(map[string]interface{})
		parameterName, parameterHasName := parameter["text"].(string)
		if !hasParameter || !parameterHasName {
			return fmt.Errorf("malformed parameter in index %d", argIndex)
		}

		argumentTerm, hasArgument := arguments[argIndex].(map[string]interface{})
		if !hasArgument {
			return fmt.Errorf("malformed argument in index %d", argIndex)
		}
		argument, evalErr := evalNode(argumentTerm, env)
		if evalErr != nil {
			return evalErr
		}

		if _, exists := function.Env.Get(parameterName); exists {
			return fmt.Errorf("mixed parameter: %s", parameterName)
		}
		function.Env.Set(parameterName, argument)
	}
	return nil
}

func (call CallTerm) Eval(env Environment) (Term, error) {
	calleeTerm := call.Callee
	callee, calleeEvalErr := evalNode(calleeTerm, env)
	if calleeEvalErr != nil {
		return nil, calleeEvalErr
	}
	function, isFunction := callee.(FunctionTerm)
	if !isFunction {
		return nil, fmt.Errorf("'Call' can only call Functions")
	}

	function.Env = NewEnvironment()
	parametersErr := function.setParameters(call.Arguments, env)
	if parametersErr != nil {
		return nil, parametersErr
	}
	function.captureEnv(env)

	result, resultEvalErr := evalNode(function.Value, function.Env)
	if resultEvalErr != nil {
		return nil, resultEvalErr
	}
	return result, nil
}
