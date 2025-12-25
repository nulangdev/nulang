package evaluator

import (
	"math"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}
	switch node.Operator {
	case "!":
		return nativeBoolToBooleanObject(!isTruthy(right))
	case "-":
		if right.Type() != object.NUMBER_OBJ {
			return newError("unknown operator: -%s", right.Type())
		}
		return &object.Number{Value: -right.(*object.Number).Value}
	case "++":
		return evalPreIncrement(node.Right, env)
	case "--":
		return evalPreDecrement(node.Right, env)
	default:
		return newError("unknown operator: %s%s", node.Operator, right.Type())
	}
}

func evalPreIncrement(node ast.Expression, env *object.Environment) object.Object {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return newError("invalid increment operand")
	}
	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: " + ident.Value)
	}
	num, ok := val.(*object.Number)
	if !ok {
		return newError("increment requires number")
	}
	newVal := &object.Number{Value: num.Value + 1}
	env.Update(ident.Value, newVal)
	return newVal
}

func evalPreDecrement(node ast.Expression, env *object.Environment) object.Object {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return newError("invalid decrement operand")
	}
	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: " + ident.Value)
	}
	num, ok := val.(*object.Number)
	if !ok {
		return newError("decrement requires number")
	}
	newVal := &object.Number{Value: num.Value - 1}
	env.Update(ident.Value, newVal)
	return newVal
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("invalid postfix operand")
	}
	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: " + ident.Value)
	}
	num, ok := val.(*object.Number)
	if !ok {
		return newError("postfix operator requires number")
	}
	oldVal := &object.Number{Value: num.Value}
	var newVal *object.Number
	switch node.Operator {
	case "++":
		newVal = &object.Number{Value: num.Value + 1}
	case "--":
		newVal = &object.Number{Value: num.Value - 1}
	default:
		return newError("unknown postfix operator: %s", node.Operator)
	}
	env.Update(ident.Value, newVal)
	return oldVal
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	// Short-circuit operators
	if node.Operator == "&&" {
		left := Eval(node.Left, env)
		if isError(left) || !isTruthy(left) {
			return left
		}
		return Eval(node.Right, env)
	}
	if node.Operator == "||" {
		left := Eval(node.Left, env)
		if isError(left) || isTruthy(left) {
			return left
		}
		return Eval(node.Right, env)
	}
	if node.Operator == "??" {
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		if left.Type() != object.NULL_OBJ && left.Type() != object.UNDEFINED_OBJ {
			return left
		}
		return Eval(node.Right, env)
	}

	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfix(node.Operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfix(node.Operator, left, right)
	case (left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ) && node.Operator == "+":
		return &object.String{Value: objectToString(left) + objectToString(right)}
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case node.Operator == "===":
		return evalStrictEquality(left, right)
	case node.Operator == "!==":
		if evalStrictEquality(left, right) == TRUE {
			return FALSE
		}
		return TRUE
	}
	return newError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

func evalNumberInfix(operator string, left, right object.Object) object.Object {
	l := left.(*object.Number).Value
	r := right.(*object.Number).Value
	switch operator {
	case "+":
		return &object.Number{Value: l + r}
	case "-":
		return &object.Number{Value: l - r}
	case "*":
		return &object.Number{Value: l * r}
	case "/":
		if r == 0 {
			if l > 0 {
				return &object.Number{Value: math.Inf(1)}
			} else if l < 0 {
				return &object.Number{Value: math.Inf(-1)}
			}
			return &object.Number{Value: math.NaN()}
		}
		return &object.Number{Value: l / r}
	case "%":
		return &object.Number{Value: math.Mod(l, r)}
	case "**":
		return &object.Number{Value: math.Pow(l, r)}
	case "<":
		return nativeBoolToBooleanObject(l < r)
	case ">":
		return nativeBoolToBooleanObject(l > r)
	case "<=":
		return nativeBoolToBooleanObject(l <= r)
	case ">=":
		return nativeBoolToBooleanObject(l >= r)
	case "==", "===":
		return nativeBoolToBooleanObject(l == r)
	case "!=", "!==":
		return nativeBoolToBooleanObject(l != r)
	}
	return newError("unknown operator: NUMBER %s NUMBER", operator)
}

func evalStringInfix(operator string, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: l + r}
	case "==", "===":
		return nativeBoolToBooleanObject(l == r)
	case "!=", "!==":
		return nativeBoolToBooleanObject(l != r)
	case "<":
		return nativeBoolToBooleanObject(l < r)
	case ">":
		return nativeBoolToBooleanObject(l > r)
	}
	return newError("unknown operator: STRING %s STRING", operator)
}

func evalStrictEquality(left, right object.Object) *object.Boolean {
	if left.Type() != right.Type() {
		return FALSE
	}
	switch l := left.(type) {
	case *object.Number:
		return nativeBoolToBooleanObject(l.Value == right.(*object.Number).Value)
	case *object.String:
		return nativeBoolToBooleanObject(l.Value == right.(*object.String).Value)
	case *object.Boolean:
		return nativeBoolToBooleanObject(l.Value == right.(*object.Boolean).Value)
	}
	return nativeBoolToBooleanObject(left == right)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}
	return UNDEFINED
}

func evalConditionalExpression(ce *ast.ConditionalExpression, env *object.Environment) object.Object {
	condition := Eval(ce.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ce.Consequence, env)
	}
	return Eval(ce.Alternative, env)
}
