package evaluator

import (
	"math"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

func evalFunctionLiteral(fl *ast.FunctionLiteral, env *object.Environment) object.Object {
	fn := &object.Function{
		Parameters: fl.Parameters,
		Body:       fl.Body,
		Env:        env,
		IsAsync:    fl.IsAsync,
	}
	if fl.Name != nil {
		fn.Name = fl.Name.Value
	}
	return fn
}

func evalCallExpression(ce *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(ce.Function, env)
	if isError(function) {
		return function
	}
	args := evalExpressions(ce.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunction(function, args)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		if spread, ok := e.(*ast.SpreadExpression); ok {
			val := Eval(spread.Value, env)
			if isError(val) {
				return []object.Object{val}
			}
			if arr, ok := val.(*object.Array); ok {
				result = append(result, arr.Elements...)
			} else {
				// If not an array, treat as single value or iterator (not supported yet)
				result = append(result, val)
			}
			continue
		}

		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.ObjectMap:
		if call, ok := fn.Get("__call__"); ok {
			return applyFunction(call, args)
		}
	case *ProxyObject:
		return ProxyApply(fn, UNDEFINED, args, nil)
	}
	return newError("not a function: %s", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		if param.IsRest {
			// Rest parameter: collect all remaining arguments into an array
			restArgs := []object.Object{}
			if i < len(args) {
				restArgs = args[i:]
			}
			env.Set(param.Value, &object.Array{Elements: restArgs})
		} else if i < len(args) {
			env.Set(param.Value, args[i])
		} else {
			env.Set(param.Value, UNDEFINED)
		}
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value
	}
	return obj
}

func evalArrayLiteral(al *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(al.Elements, env)
	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}
	return &object.Array{Elements: elements}
}

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(ie.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(ie.Index, env)
	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		arr := left.(*object.Array)
		idx := int(index.(*object.Number).Value)
		if idx < 0 || idx >= len(arr.Elements) {
			return UNDEFINED
		}
		return arr.Elements[idx]
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		str := left.(*object.String)
		idx := int(index.(*object.Number).Value)
		if idx < 0 || idx >= len(str.Value) {
			return UNDEFINED
		}
		return &object.String{Value: string(str.Value[idx])}
	case left.Type() == object.OBJECT_OBJ:
		objMap := left.(*object.ObjectMap)
		key := objectToString(index)
		if val, ok := objMap.Get(key); ok {
			return val
		}
		return UNDEFINED
	}
	return newError("index operator not supported: %s", left.Type())
}

func evalObjectLiteral(ol *ast.ObjectLiteral, env *object.Environment) object.Object {
	pairs := make(map[string]object.ObjectPair)
	for keyNode, valueNode := range ol.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		keyStr := objectToString(key)
		pairs[keyStr] = object.ObjectPair{Key: key, Value: value}
	}
	return &object.ObjectMap{Pairs: pairs}
}

func evalMemberExpression(me *ast.MemberExpression, env *object.Environment) object.Object {
	obj := Eval(me.Object, env)
	if isError(obj) {
		return obj
	}
	if me.Optional && (obj.Type() == object.NULL_OBJ || obj.Type() == object.UNDEFINED_OBJ) {
		return UNDEFINED
	}
	propName := me.Property.(*ast.Identifier).Value

	// Handle Proxy objects
	if proxy, ok := obj.(*ProxyObject); ok {
		return ProxyGet(proxy, propName, env)
	}

	switch o := obj.(type) {
	case *object.Array:
		return evalArrayProperty(o, propName)
	case *object.String:
		return evalStringProperty(o, propName)
	case *object.Function:
		return evalFunctionProperty(o, propName)
	case *object.Buffer:
		return evalBufferProperty(o, propName)
	case *object.Promise:
		return evalPromiseProperty(o, propName, env)
	case *object.ObjectMap:
		if val, ok := o.Get(propName); ok {
			return val
		}
	case *Class:
		// Access static members of a class
		if val, ok := o.Static[propName]; ok {
			return val
		}
	}
	return UNDEFINED
}

func evalFunctionProperty(fn *object.Function, propName string) object.Object {
	switch propName {
	case "apply":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 2 {
					return newError("apply requires at least 2 arguments: thisArg and argsArray")
				}
				thisArg := args[0]
				argsArray := args[1] 
				
				var fnArgs []object.Object
				if arr, ok := argsArray.(*object.Array); ok {
					fnArgs = arr.Elements
				} else {
					return newError("apply expects second argument to be an array")
				}
				
				// Create a bound version of the function with the given thisArg
				// But simpler: just execute the body with a new env properly set up
				
				// We need to use extendFunctionEnv but also inject 'this'
				// The clean way is to reuse applyFunction logic but modifying the env first?
				// Actually createClassInstance logic was complex.
				
				// Let's create a temporary environment for execution
				env := extendFunctionEnv(fn, fnArgs)
				env.Set("this", thisArg)
				
				evaluated := Eval(fn.Body, env)
				return unwrapReturnValue(evaluated)
			},
		}
	case "call":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("call requires at least 1 argument: thisArg")
				}
				thisArg := args[0]
				fnArgs := args[1:]
				
				env := extendFunctionEnv(fn, fnArgs)
				env.Set("this", thisArg)
				
				evaluated := Eval(fn.Body, env)
				return unwrapReturnValue(evaluated)
			},
		}
	case "bind":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("bind requires at least 1 argument: thisArg")
				}
				thisArg := args[0]
				boundArgs := args[1:]
				
				// Return a new function that calls the original one with fixed 'this' and prepended args
				return &object.Builtin{
					Fn: func(callArgs ...object.Object) object.Object {
						// Merge bound args and call args
						allArgs := append(boundArgs, callArgs...)
						
						env := extendFunctionEnv(fn, allArgs)
						env.Set("this", thisArg)
						
						evaluated := Eval(fn.Body, env)
						return unwrapReturnValue(evaluated)
					},
				}
			},
		}
	}
	return UNDEFINED
}

func evalAssignmentExpression(ae *ast.AssignmentExpression, env *object.Environment) object.Object {
	right := Eval(ae.Right, env)
	if isError(right) {
		return right
	}

	switch left := ae.Left.(type) {
	case *ast.Identifier:
		if env.IsConst(left.Value) {
			return newError("Assignment to constant variable '%s'", left.Value)
		}
		newVal := computeAssignment(ae.Operator, left.Value, right, env)
		if _, ok := env.Update(left.Value, newVal); !ok {
			env.Set(left.Value, newVal)
		}
		return newVal

	case *ast.MemberExpression:
		obj := Eval(left.Object, env)
		if isError(obj) {
			return obj
		}
		propName := left.Property.(*ast.Identifier).Value
		
		// Handle Proxy objects
		if proxy, ok := obj.(*ProxyObject); ok {
			result := ProxySet(proxy, propName, right, env)
			if isError(result) {
				return result
			}
			return right
		}
		
		if objMap, ok := obj.(*object.ObjectMap); ok {
			objMap.Set(propName, right)
			return right
		}

	case *ast.IndexExpression:
		obj := Eval(left.Left, env)
		if isError(obj) {
			return obj
		}
		index := Eval(left.Index, env)
		if isError(index) {
			return index
		}
		if arr, ok := obj.(*object.Array); ok {
			idx := int(index.(*object.Number).Value)
			if idx >= 0 {
				// Expand array if needed
				if idx >= len(arr.Elements) {
					// Create new slice with enough capacity
					newElements := make([]object.Object, idx+1)
					// Copy existing elements
					copy(newElements, arr.Elements)
					// Fill gaps with undefined
					for i := len(arr.Elements); i < idx; i++ {
						newElements[i] = UNDEFINED
					}
					arr.Elements = newElements
				}
				arr.Elements[idx] = right
			}
			return right
		}
		if objMap, ok := obj.(*object.ObjectMap); ok {
			key := objectToString(index)
			objMap.Set(key, right)
			return right
		}
	}
	return newError("invalid assignment target")
}

func computeAssignment(op, name string, right object.Object, env *object.Environment) object.Object {
	if op == "=" {
		return right
	}
	oldVal, _ := env.Get(name)
	oldNum, isNum := oldVal.(*object.Number)
	rightNum, rightIsNum := right.(*object.Number)

	if isNum && rightIsNum {
		switch op {
		case "+=":
			return &object.Number{Value: oldNum.Value + rightNum.Value}
		case "-=":
			return &object.Number{Value: oldNum.Value - rightNum.Value}
		case "*=":
			return &object.Number{Value: oldNum.Value * rightNum.Value}
		case "/=":
			return &object.Number{Value: oldNum.Value / rightNum.Value}
		case "%=":
			return &object.Number{Value: math.Mod(oldNum.Value, rightNum.Value)}
		}
	}
	if op == "+=" {
		if oldStr, ok := oldVal.(*object.String); ok {
			return &object.String{Value: oldStr.Value + objectToString(right)}
		}
	}
	return right
}

func evalTypeofExpression(te *ast.TypeofExpression, env *object.Environment) object.Object {
	val := Eval(te.Value, env)
	if isError(val) {
		errObj := val.(*object.Error)
		if len(errObj.Message) > 21 && errObj.Message[:21] == "identifier not found:" {
			return &object.String{Value: "undefined"}
		}
		return val
	}
	switch val.Type() {
	case object.NUMBER_OBJ:
		return &object.String{Value: "number"}
	case object.STRING_OBJ:
		return &object.String{Value: "string"}
	case object.BOOLEAN_OBJ:
		return &object.String{Value: "boolean"}
	case object.NULL_OBJ:
		return &object.String{Value: "object"}
	case object.UNDEFINED_OBJ:
		return &object.String{Value: "undefined"}
	case object.FUNCTION_OBJ, object.BUILTIN_OBJ:
		return &object.String{Value: "function"}
	default:
		return &object.String{Value: "object"}
	}
}

func evalThisExpression(env *object.Environment) object.Object {
	if val, ok := env.Get("this"); ok {
		return val
	}
	return UNDEFINED
}
