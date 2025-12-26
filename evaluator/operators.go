package evaluator

import (
	"math"
	"strconv"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	if node.Operator == "delete" {
		return evalDeleteExpression(node.Right, env)
	}

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

// evalDeleteExpression handles the delete operator
func evalDeleteExpression(node ast.Expression, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.MemberExpression:
		// obj.prop
		left := Eval(n.Object, env)
		if isError(left) {
			return left
		}

		propName := ""
		if ident, ok := n.Property.(*ast.Identifier); ok {
			propName = ident.Value
		} else {
			return newError("invalid delete operand")
		}

		// Handle Proxy
		if proxy, ok := left.(*ProxyObject); ok {
			return ProxyDeleteProperty(proxy, propName, nil)
		}

		if obj, ok := left.(*object.ObjectMap); ok {
			delete(obj.Pairs, propName)
			return TRUE
		}
		return TRUE

	case *ast.IndexExpression:
		// obj["prop"]
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(n.Index, env)
		if isError(index) {
			return index
		}
		propName := objectToString(index)

		// Handle Proxy
		if proxy, ok := left.(*ProxyObject); ok {
			return ProxyDeleteProperty(proxy, propName, nil)
		}

		if obj, ok := left.(*object.ObjectMap); ok {
			delete(obj.Pairs, propName)
			return TRUE
		}
		return TRUE

	case *ast.Identifier:
		// delete variable - Nulang defaults to returning false/undefined for var deletion or we can implement it if needed.
		// For now, let's treat it as no-op return false (strict mode-ish).
		return FALSE

	default:
		return TRUE
	}
}

func evalPreIncrement(node ast.Expression, env *object.Environment) object.Object {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return newError("invalid increment operand")
	}
	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("%s", "identifier not found: " + ident.Value)
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
		return newError("%s", "identifier not found: " + ident.Value)
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
		return newError("%s", "identifier not found: " + ident.Value)
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
	case node.Operator == "instanceof":
		return evalInstanceof(left, right)
	case node.Operator == "in":
		return evalInOperator(left, right)
	}
	return newError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

// evalInstanceof checks if left is an instance of right
func evalInstanceof(left, right object.Object) object.Object {
	// Check if right is a Class
	if class, ok := right.(*Class); ok {
		// Check if left is an ObjectMap (instance)
		if objMap, ok := left.(*object.ObjectMap); ok {
			// Check if the object has properties/methods from this class
			// by checking prototype chain or special markers
			if className, found := objMap.Get("__class__"); found {
				if classNameStr, ok := className.(*object.String); ok {
					if classNameStr.Value == class.Name {
						return TRUE
					}
				}
			}
			// Walk prototype chain to check for class methods
			for methodName := range class.Methods {
				if _, found := objMap.Get(methodName); !found {
					return FALSE
				}
			}
			if len(class.Methods) > 0 {
				return TRUE
			}
		}
		return FALSE
	}

	// Check for Error class (special case using isErrorInstance)
	if builtin, ok := right.(*object.Builtin); ok && builtin.Name == "Error" {
		return nativeBoolToBooleanObject(isErrorInstance(left))
	}

	// Check for ObjectMap with __name__ (built-in classes like Error, Date, etc.)
	if rightMap, ok := right.(*object.ObjectMap); ok {
		if nameVal, found := rightMap.Get("__name__"); found {
			if nameStr, ok := nameVal.(*object.String); ok {
				switch nameStr.Value {
				case "Error":
					return nativeBoolToBooleanObject(isErrorInstance(left))
				case "Array":
					_, isArray := left.(*object.Array)
					return nativeBoolToBooleanObject(isArray)
				case "Promise":
					_, isPromise := left.(*object.Promise)
					return nativeBoolToBooleanObject(isPromise)
				case "RegExp":
					_, isRegex := left.(*object.RegExp)
					return nativeBoolToBooleanObject(isRegex)
				}
			}
		}
		// Check if it's an Error class by name property
		if nameVal, found := rightMap.Get("name"); found {
			if nameStr, ok := nameVal.(*object.String); ok && nameStr.Value == "Error" {
				return nativeBoolToBooleanObject(isErrorInstance(left))
			}
		}
	}

	// Check Array
	if _, ok := left.(*object.Array); ok {
		if ident, ok := right.(*object.Builtin); ok && ident.Name == "Array" {
			return TRUE
		}
	}

	// Check Promise  
	if _, ok := left.(*object.Promise); ok {
		if ident, ok := right.(*object.Builtin); ok && ident.Name == "Promise" {
			return TRUE
		}
	}

	// Check RegExp
	if _, ok := left.(*object.RegExp); ok {
		if ident, ok := right.(*object.Builtin); ok && ident.Name == "RegExp" {
			return TRUE
		}
	}

	return FALSE
}

// evalInOperator checks if property exists in object
func evalInOperator(left, right object.Object) object.Object {
	propStr := objectToString(left)

	// Check Proxy
	if proxy, ok := right.(*ProxyObject); ok {
		return ProxyHas(proxy, propStr, nil)
	}

	// Check ObjectMap
	if obj, ok := right.(*object.ObjectMap); ok {
		current := obj
		for current != nil {
			if _, found := current.Get(propStr); found {
				return TRUE
			}
			current = current.Prototype
		}
		return FALSE
	}

	// Check Array
	if arr, ok := right.(*object.Array); ok {
		if propStr == "length" {
			return TRUE
		}
		// Indices
		if idx, err := strconv.Atoi(propStr); err == nil {
			if idx >= 0 && idx < len(arr.Elements) {
				return TRUE
			}
		}
	}

	return FALSE
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
