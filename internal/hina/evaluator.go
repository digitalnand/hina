package hina

import (
	"fmt"
	"reflect"
)

func (print PrintNode) Evaluate() any {
	fmt.Println(print.Value)
	return print.Value
}

func (binary BinaryNode) Evaluate() (any, error) {
	lhs := binary.Lhs
	rhs := binary.Rhs

	switch binary.Op {
	case "Add":
		// TODO: improve this
		_, isLhsString := lhs.(StrNode)
		_, isRhsString := rhs.(StrNode)
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		if isLhsInt && isRhsInt {
			return IntNode{Value: intLhs.Value + intRhs.Value}, nil
		}
		if (isLhsString || isLhsInt) && (isRhsInt || isRhsString) {
			return StrNode{Value: fmt.Sprintf("%s%s", lhs, rhs)}, nil
		}
		return nil, fmt.Errorf("'Add' operator can only be used with Ints and/or Strs")
	case "Sub", "Mul", "Div", "Rem":
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		var result int32
		if isLhsInt && isRhsInt {
			switch binary.Op {
			case "Sub":
				result = intLhs.Value - intRhs.Value
			case "Mul":
				result = intLhs.Value * intRhs.Value
			case "Div":
				result = intLhs.Value / intRhs.Value
			case "Rem":
				result = intLhs.Value % intRhs.Value
			}
			return IntNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' operator can only be used with Ints", binary.Op)
	case "Eq", "Neq":
		hasSameValue := lhs == rhs
		hasSameType := reflect.TypeOf(lhs) == reflect.TypeOf(rhs)
		var result bool
		switch binary.Op {
		case "Eq":
			result = hasSameValue && hasSameType
		case "Neq":
			result = !hasSameValue || !hasSameType
		}
		return BoolNode{Value: result}, nil
	case "Lt", "Gt", "Lte", "Gte":
		intLhs, isLhsInt := lhs.(IntNode)
		intRhs, isRhsInt := rhs.(IntNode)
		var result bool
		if isLhsInt && isRhsInt {
			switch binary.Op {
			case "Lt":
				result = intLhs.Value < intRhs.Value
			case "Gt":
				result = intLhs.Value > intRhs.Value
			case "Lte":
				result = intLhs.Value <= intRhs.Value
			case "Gte":
				result = intLhs.Value >= intRhs.Value
			}
			return BoolNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' comparison can only be done with Ints", binary.Op)
	case "And", "Or":
		boolLhs, isLhsBool := lhs.(BoolNode)
		boolRhs, isRhsBool := rhs.(BoolNode)
		var result bool
		if isLhsBool && isRhsBool {
			switch binary.Op {
			case "And":
				result = boolLhs.Value && boolRhs.Value
			case "Or":
				result = boolLhs.Value || boolRhs.Value
			}
			return BoolNode{Value: result}, nil
		}
		return nil, fmt.Errorf("'%s' operator can only be used with Bool", binary.Op)
	default:
		return nil, fmt.Errorf("unknown binary operator: '%s'", binary.Op)
	}
}
