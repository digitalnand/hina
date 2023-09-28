package hina

import (
	"fmt"
	"math"
	"reflect"
)

func EvalTree(tree Object, env Environment) error {
	exp, exists := tree["expression"].(map[string]interface{})
	if !exists || len(exp) == 0 {
		return fmt.Errorf("tree has no expressions")
	}
	term, inspectErr := InspectNode(exp)
	if inspectErr != nil {
		return inspectErr
	}
	_, evalErr := evalTerm(term, env)
	if evalErr != nil {
		return evalErr
	}
	return nil
}

func evalTerm(term Term, env Environment) (Term, error) {
	switch termType := term.(type) {
	case StrTerm, IntTerm, BoolTerm:
		return term, nil
	case PrintTerm:
		print := term.(PrintTerm)
		returnTerm, err := print.Eval(env)
		if err != nil {
			return nil, err
		}
		return returnTerm, nil
	case BinaryTerm:
		binary := term.(BinaryTerm)
		result, err := binary.Eval(env)
		if err != nil {
			return nil, err
		}
		return result, nil
	case LetTerm:
		let := term.(LetTerm)
		nextResult, err := let.Eval(env)
		if err != nil {
			return nil, err
		}
		return nextResult, nil
	case VarTerm:
		varTerm := term.(VarTerm)
		value, err := varTerm.Eval(env)
		if err != nil {
			return nil, err
		}
		return value, nil
	case TupleTerm:
		tuple := term.(TupleTerm)
		first, firstEvalErr := evalTerm(tuple.First, env)
		if firstEvalErr != nil {
			return nil, firstEvalErr
		}
		second, secondEvalErr := evalTerm(tuple.Second, env)
		if secondEvalErr != nil {
			return nil, secondEvalErr
		}
		tuple.First = first
		tuple.Second = second
		return tuple, nil
	case TupleFunction:
		tupleFunc := term.(TupleFunction)
		value, err := tupleFunc.Eval(env)
		if err != nil {
			return nil, err
		}
		return value, nil
	case IfTerm:
		ifTerm := term.(IfTerm)
		result, err := ifTerm.Eval(env)
		if err != nil {
			return nil, err
		}
		return result, nil
	case FunctionTerm:
		function := term.(FunctionTerm)
		function.Env.Copy(env)
		return function, nil
	case CallTerm:
		call := term.(CallTerm)
		result, err := call.Eval(env)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown term: %s", termType)
	}
}

func (print PrintTerm) Eval(env Environment) (Term, error) {
	value, err := evalTerm(print.Value, env)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	return value, nil
}

func (binary BinaryTerm) Eval(env Environment) (Term, error) {
	lhs, lhsEvalErr := evalTerm(binary.Lhs, env)
	if lhsEvalErr != nil {
		return nil, lhsEvalErr
	}
	rhs, rhsEvalErr := evalTerm(binary.Rhs, env)
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
		if !isLhsInt || !isRhsInt {
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
			return IntTerm{Value: math.Mod(intLhs.Value, intRhs.Value)}, nil
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
		if !isLhsInt || !isRhsInt {
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
		if !isLhsBool || !isRhsBool {
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
	env.Set(variable.Identifier, variable.Value)
	nextResult, nextEvalErr := evalTerm(variable.Next, env)
	if nextEvalErr != nil {
		return nil, nextEvalErr
	}
	return nextResult, nil
}

func (varCall VarTerm) Eval(env Environment) (Term, error) {
	valueTerm, exists := env.Get(varCall.Text)
	if !exists {
		return nil, fmt.Errorf("calling an undeclared variable: %s", varCall.Text)
	}
	value, err := evalTerm(valueTerm, env)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (tupleFunc TupleFunction) Eval(env Environment) (Term, error) {
	value, err := evalTerm(tupleFunc.Value, env)
	if err != nil {
		return nil, err
	}
	tuple, isTuple := value.(TupleTerm)
	if !isTuple {
		return nil, fmt.Errorf("'%s' only accepts Tuples", tupleFunc.Kind)
	}
	if tupleFunc.Kind == "First" {
		return tuple.First, nil
	} else {
		return tuple.Second, nil
	}
}

func (ifTerm IfTerm) Eval(env Environment) (Term, error) {
	conditionTerm, conditionEvalErr := evalTerm(ifTerm.Condition, env)
	if conditionEvalErr != nil {
		return nil, conditionEvalErr
	}
	condition, isBool := conditionTerm.(BoolTerm)
	if !isBool {
		return nil, fmt.Errorf("'If' only accepts Bools as condition")
	}

	var body Term
	if condition.Value {
		body = ifTerm.Then
	} else {
		body = ifTerm.Else
	}
	result, bodyEvalErr := evalTerm(body, env)
	if bodyEvalErr != nil {
		return nil, bodyEvalErr
	}
	return result, nil
}

func (call CallTerm) insertArgs(function FunctionTerm, env Environment) error {
	if len(function.Parameters) != len(call.Arguments) {
		return fmt.Errorf("expected %d arguments, received %d", len(function.Parameters), len(call.Arguments))
	}
	for index := 0; index < len(call.Arguments); index++ {
		parameter := function.Parameters[index]
		argument, err := evalTerm(call.Arguments[index], env)
		if err != nil {
			return err
		}
		function.Env.Set(parameter, argument)
	}
	return nil
}

func (call CallTerm) Eval(env Environment) (Term, error) {
	calleeTerm := call.Callee
	callee, calleeEvalErr := evalTerm(calleeTerm, env)
	if calleeEvalErr != nil {
		return nil, calleeEvalErr
	}
	function, isFunction := callee.(FunctionTerm)
	if !isFunction {
		return nil, fmt.Errorf("'Call' can only call Functions")
	}
	paramsErr := call.insertArgs(function, env)
	if paramsErr != nil {
		return nil, paramsErr
	}

	newEnv := NewEnv()
	newEnv.Copy(function.Env)
	result, resultEvalErr := evalTerm(function.Value, newEnv)
	if resultEvalErr != nil {
		return nil, resultEvalErr
	}
	return result, nil
}
