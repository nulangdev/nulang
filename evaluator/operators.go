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
	case "+":
		// Unary plus converts to number
		switch r := right.(type) {
		case *object.Number:
			return r
		case *object.String:
			if val, err := strconv.ParseFloat(r.Value, 64); err == nil {
				return &object.Number{Value: val}
			}
			return &object.Number{Value: math.NaN()}
		case *object.Boolean:
			if r.Value {
				return &object.Number{Value: 1}
			}
			return &object.Number{Value: 0}
		case *object.Null, *object.Undefined:
			return &object.Number{Value: 0}
		default:
			return &object.Number{Value: math.NaN()}
		}
	case "-":
		if right.Type() != object.NUMBER_OBJ {
			return newError("unknown operator: -%s", right.Type())
		}
		return &object.Number{Value: -right.(*object.Number).Value}
	case "~":
		if right.Type() != object.NUMBER_OBJ {
			return newError("unknown operator: ~%s", right.Type())
		}
		// Bitwise NOT: convert to int32, NOT, convert back
		return &object.Number{Value: float64(^int32(right.(*object.Number).Value))}
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
	return evalPrefixIncDec(node, env, 1)
}

func evalPreDecrement(node ast.Expression, env *object.Environment) object.Object {
	return evalPrefixIncDec(node, env, -1)
}

// evalPrefixIncDec handles both ++x and --x for various left-hand side expressions
func evalPrefixIncDec(node ast.Expression, env *object.Environment, delta float64) object.Object {
	switch left := node.(type) {
	case *ast.Identifier:
		// Simple variable: ++x
		val, ok := env.Get(left.Value)
		if !ok {
			return newError("%s", "identifier not found: "+left.Value)
		}
		num, ok := val.(*object.Number)
		if !ok {
			return newError("increment/decrement requires number")
		}
		newVal := &object.Number{Value: num.Value + delta}
		env.Update(left.Value, newVal)
		return newVal

	case *ast.MemberExpression:
		// Object property: ++obj.prop
		obj := Eval(left.Object, env)
		if isError(obj) {
			return obj
		}
		propName := left.Property.(*ast.Identifier).Value
		
		switch o := obj.(type) {
		case *object.ObjectMap:
			if val, ok := o.Get(propName); ok {
				if num, ok := val.(*object.Number); ok {
					newVal := &object.Number{Value: num.Value + delta}
					o.Set(propName, newVal)
					return newVal
				}
				return newError("increment/decrement requires number")
			}
			return newError("property not found: %s", propName)
		default:
			return newError("cannot use increment/decrement on %s", obj.Type())
		}

	case *ast.IndexExpression:
		// Array index: ++arr[i]
		obj := Eval(left.Left, env)
		if isError(obj) {
			return obj
		}
		idx := Eval(left.Index, env)
		if isError(idx) {
			return idx
		}

		switch o := obj.(type) {
		case *object.Array:
			if idxNum, ok := idx.(*object.Number); ok {
				i := int(idxNum.Value)
				if i < 0 || i >= len(o.Elements) {
					return newError("array index out of bounds")
				}
				if num, ok := o.Elements[i].(*object.Number); ok {
					newVal := &object.Number{Value: num.Value + delta}
					o.Elements[i] = newVal
					return newVal
				}
				return newError("increment/decrement requires number")
			}
			return newError("array index must be a number")
		case *object.ObjectMap:
			key := objectToString(idx)
			if val, ok := o.Get(key); ok {
				if num, ok := val.(*object.Number); ok {
					newVal := &object.Number{Value: num.Value + delta}
					o.Set(key, newVal)
					return newVal
				}
				return newError("increment/decrement requires number")
			}
			return newError("property not found: %s", key)
		default:
			return newError("cannot use increment/decrement on %s", obj.Type())
		}
	}

	return newError("invalid increment/decrement operand")
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	// Calculate the increment/decrement based on operator
	var delta float64
	switch node.Operator {
	case "++":
		delta = 1
	case "--":
		delta = -1
	default:
		return newError("unknown postfix operator: %s", node.Operator)
	}

	// Handle different types of left-hand side expressions
	switch left := node.Left.(type) {
	case *ast.Identifier:
		// Simple variable: x++
		val, ok := env.Get(left.Value)
		if !ok {
			return newError("%s", "identifier not found: "+left.Value)
		}
		num, ok := val.(*object.Number)
		if !ok {
			return newError("postfix operator requires number")
		}
		oldVal := &object.Number{Value: num.Value}
		newVal := &object.Number{Value: num.Value + delta}
		env.Update(left.Value, newVal)
		return oldVal

	case *ast.MemberExpression:
		// Object property: obj.prop++
		obj := Eval(left.Object, env)
		if isError(obj) {
			return obj
		}
		propName := left.Property.(*ast.Identifier).Value
		
		var oldNum float64
		switch o := obj.(type) {
		case *object.ObjectMap:
			if val, ok := o.Get(propName); ok {
				if num, ok := val.(*object.Number); ok {
					oldNum = num.Value
					o.Set(propName, &object.Number{Value: oldNum + delta})
				} else {
					return newError("postfix operator requires number")
				}
			} else {
				return newError("property not found: %s", propName)
			}
		default:
			return newError("cannot use postfix operator on %s", obj.Type())
		}
		return &object.Number{Value: oldNum}

	case *ast.IndexExpression:
		// Array index: arr[i]++
		obj := Eval(left.Left, env)
		if isError(obj) {
			return obj
		}
		idx := Eval(left.Index, env)
		if isError(idx) {
			return idx
		}

		switch o := obj.(type) {
		case *object.Array:
			if idxNum, ok := idx.(*object.Number); ok {
				i := int(idxNum.Value)
				if i < 0 || i >= len(o.Elements) {
					return newError("array index out of bounds")
				}
				if num, ok := o.Elements[i].(*object.Number); ok {
					oldVal := num.Value
					o.Elements[i] = &object.Number{Value: oldVal + delta}
					return &object.Number{Value: oldVal}
				}
				return newError("postfix operator requires number")
			}
			return newError("array index must be a number")
		case *object.ObjectMap:
			key := objectToString(idx)
			if val, ok := o.Get(key); ok {
				if num, ok := val.(*object.Number); ok {
					oldVal := num.Value
					o.Set(key, &object.Number{Value: oldVal + delta})
					return &object.Number{Value: oldVal}
				}
				return newError("postfix operator requires number")
			}
			return newError("property not found: %s", key)
		default:
			return newError("cannot use postfix operator on %s", obj.Type())
		}
	}

	return newError("invalid postfix operand")
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
		return evalLooseEquality(left, right)
	case node.Operator == "!=":
		if evalLooseEquality(left, right) == TRUE {
			return FALSE
		}
		return TRUE
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
	
	// Handle arithmetic operators with type coercion to numbers
	// This matches JavaScript behavior where e.g. undefined - 10 = NaN
	if isArithmeticOperator(node.Operator) {
		leftNum := toNumber(left)
		rightNum := toNumber(right)
		return evalNumberInfix(node.Operator, &object.Number{Value: leftNum}, &object.Number{Value: rightNum})
	}
	
	// Handle comparison operators with type coercion
	if isComparisonOperator(node.Operator) {
		leftNum := toNumber(left)
		rightNum := toNumber(right)
		return evalNumberInfix(node.Operator, &object.Number{Value: leftNum}, &object.Number{Value: rightNum})
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
	case "&":
		return &object.Number{Value: float64(int32(l) & int32(r))}
	case "|":
		return &object.Number{Value: float64(int32(l) | int32(r))}
	case "^":
		return &object.Number{Value: float64(int32(l) ^ int32(r))}
	case "<<":
		return &object.Number{Value: float64(int32(l) << uint32(r))}
	case ">>":
		return &object.Number{Value: float64(int32(l) >> uint32(r))}
	case ">>>":
		return &object.Number{Value: float64(uint32(l) >> uint32(r))}
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

// evalLooseEquality implements JavaScript's Abstract Equality Comparison (==)
// Key rules:
// - null == undefined is true
// - undefined == null is true
// - Same types compare values
// - String vs Number: convert string to number
// - Boolean vs anything: convert boolean to number first
func evalLooseEquality(left, right object.Object) *object.Boolean {
	// Rule 1: same types - use strict equality for value comparison
	if left.Type() == right.Type() {
		return evalStrictEquality(left, right)
	}

	// Rule 2: null == undefined is true (both directions)
	if (left.Type() == object.NULL_OBJ && right.Type() == object.UNDEFINED_OBJ) ||
		(left.Type() == object.UNDEFINED_OBJ && right.Type() == object.NULL_OBJ) {
		return TRUE
	}

	// Rule 3: null and undefined are only equal to each other, not to anything else
	if left.Type() == object.NULL_OBJ || left.Type() == object.UNDEFINED_OBJ ||
		right.Type() == object.NULL_OBJ || right.Type() == object.UNDEFINED_OBJ {
		return FALSE
	}

	// Rule 4: Number vs String - convert string to number
	if left.Type() == object.NUMBER_OBJ && right.Type() == object.STRING_OBJ {
		rightNum := toNumber(right)
		return nativeBoolToBooleanObject(left.(*object.Number).Value == rightNum)
	}
	if left.Type() == object.STRING_OBJ && right.Type() == object.NUMBER_OBJ {
		leftNum := toNumber(left)
		return nativeBoolToBooleanObject(leftNum == right.(*object.Number).Value)
	}

	// Rule 5: Boolean vs anything - convert boolean to number and retry
	if left.Type() == object.BOOLEAN_OBJ {
		leftNum := toNumber(left)
		return evalLooseEquality(&object.Number{Value: leftNum}, right)
	}
	if right.Type() == object.BOOLEAN_OBJ {
		rightNum := toNumber(right)
		return evalLooseEquality(left, &object.Number{Value: rightNum})
	}

	// Rule 6: Object vs primitive - for now, use reference equality
	// Full JS would call ToPrimitive here
	return nativeBoolToBooleanObject(left == right)
}

// toNumber converts an object to a number (JavaScript ToNumber)
func toNumber(obj object.Object) float64 {
	switch o := obj.(type) {
	case *object.Number:
		return o.Value
	case *object.String:
		val, err := strconv.ParseFloat(o.Value, 64)
		if err != nil {
			return math.NaN()
		}
		return val
	case *object.Boolean:
		if o.Value {
			return 1
		}
		return 0
	case *object.Null:
		return 0
	case *object.Undefined:
		return math.NaN()
	default:
		return math.NaN()
	}
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

// isArithmeticOperator checks if the operator is an arithmetic operator
func isArithmeticOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%", "**", "&", "|", "^", "<<", ">>", ">>>":
		return true
	}
	return false
}

// isComparisonOperator checks if the operator is a comparison operator  
func isComparisonOperator(op string) bool {
	switch op {
	case "<", ">", "<=", ">=":
		return true
	}
	return false
}
