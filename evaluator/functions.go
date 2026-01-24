package evaluator

import (
	"fmt"
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
	// Check if this is a method call (obj.method())
	// If so, we need to bind 'this' to the object
	var thisArg object.Object
	var function object.Object
	
	if me, ok := ce.Function.(*ast.MemberExpression); ok {
		// Evaluate the object first
		thisArg = Eval(me.Object, env)
		if isError(thisArg) {
			return thisArg
		}
		
		// Handle optional chaining
		if me.Optional && (thisArg.Type() == object.NULL_OBJ || thisArg.Type() == object.UNDEFINED_OBJ) {
			return UNDEFINED
		}
		
		// Get the method from the object
		function = evalMemberExpression(me, env)
		if isError(function) {
			return function
		}
	} else {
		function = Eval(ce.Function, env)
		if isError(function) {
			return function
		}
	}
	

	args := evalExpressions(ce.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunctionWithThis(function, args, thisArg)
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
	return applyFunctionWithThis(fn, args, nil)
}

// applyFunctionWithThis applies a function with an optional 'this' context
// This is used for method calls where obj.method() should have 'this' bound to 'obj'
func applyFunctionWithThis(fn object.Object, args []object.Object, thisArg object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		// Bind 'this' if provided
		if thisArg != nil {
			extendedEnv.Set("this", thisArg)
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.ObjectMap:
		if call, ok := fn.Get("__call__"); ok {
			return applyFunctionWithThis(call, args, thisArg)
		}
	case *ProxyObject:
		return ProxyApply(fn, UNDEFINED, args, nil)
	}

	return newError("not a function: %s", fn.Type())
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	
	// Create the 'arguments' object (array-like object containing all passed arguments)
	// In JavaScript, 'arguments' is available in all non-arrow functions
	argumentsArr := &object.Array{Elements: args}
	env.Set("arguments", argumentsArr)
	
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
	case left.Type() == object.STRING_OBJ:
		// String property access (e.g., str['length'])
		key := objectToString(index)
		str := left.(*object.String)
		if key == "length" {
			return &object.Number{Value: float64(len(str.Value))}
		}
		// Support string method access via bracket notation (e.g., str['toUpperCase'])
		// This is needed for lodash patterns like chr[methodName]()
		return evalStringProperty(str, key)
	case left.Type() == object.ARRAY_OBJ && index.Type() != object.NUMBER_OBJ:
		// Array property access with non-numeric index (e.g., arr['length'], arr['push'])
		arr := left.(*object.Array)
		key := objectToString(index)
		switch key {
		case "length":
			return &object.Number{Value: float64(len(arr.Elements))}
		default:
			// Return array method if available via evalArrayProperty
			return evalArrayProperty(arr, key)
		}
	case left.Type() == object.OBJECT_OBJ:
		objMap := left.(*object.ObjectMap)
		key := objectToString(index)
		if val, ok := objMap.Get(key); ok {
			return val
		}
		return UNDEFINED
	case left.Type() == object.FUNCTION_OBJ:
		// Functions are objects in JavaScript and can have properties accessed via bracket notation
		fn := left.(*object.Function)
		key := objectToString(index)
		if val, ok := fn.Get(key); ok {
			return val
		}
		return UNDEFINED
	case left.Type() == object.BUILTIN_OBJ:
		// Builtin functions can also have properties accessed via bracket notation
		fn := left.(*object.Builtin)
		key := objectToString(index)
		return evalBuiltinProperty(fn, key)
	case left.Type() == object.UNDEFINED_OBJ || left.Type() == object.NULL_OBJ:
		// In strict JavaScript, accessing properties of undefined/null throws TypeError
		// But for library compatibility, we'll return undefined instead
		// This allows optional chaining patterns and guard expressions to work
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
	case *object.Builtin:
		return evalBuiltinProperty(o, propName)
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
	case "length":
		// In JavaScript, function.length returns the number of declared parameters
		// This is critical for lodash's baseRest/overRest which use func.length to determine
		// where rest parameters start
		return &object.Number{Value: float64(len(fn.Parameters))}
	case "name":
		// Function name property
		name := fn.Name
		if name == "" {
			name = "anonymous"
		}
		return &object.String{Value: name}
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
	
	// Check for custom properties attached to the function
	if val, ok := fn.Get(propName); ok {
		return val
	}
	
	// Special handling for prototype - auto-create if it doesn't exist
	// In JavaScript, all functions have a prototype property
	if propName == "prototype" {
		prototype := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		prototype.Set("constructor", fn)
		fn.Set("prototype", prototype)
		return prototype
	}
	
	return UNDEFINED
}

// evalBuiltinProperty handles property access on builtin functions (.call, .apply, .bind)
func evalBuiltinProperty(fn *object.Builtin, propName string) object.Object {
	switch propName {
	case "apply":
		return &object.Builtin{
			Name: "apply",
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return fn.Fn()
				}
				// thisArg is ignored for builtins, just use the arguments
				if len(args) >= 2 {
					if arr, ok := args[1].(*object.Array); ok {
						return fn.Fn(arr.Elements...)
					}
				}
				return fn.Fn()
			},
		}
	case "call":
		return &object.Builtin{
			Name: "call",
			Fn: func(args ...object.Object) object.Object {
				// For .call(thisArg, arg1, arg2, ...), we need to invoke the function
				// with thisArg as context. For builtins like Object.prototype.toString,
				// we pass thisArg as the first argument so it can detect the type.
				if len(args) < 1 {
					return fn.Fn()
				}
				// Pass thisArg as first argument, followed by any additional arguments
				// This enables Object.prototype.toString.call(value) to work correctly
				allArgs := args // Include thisArg as first arg
				return fn.Fn(allArgs...)
			},
		}
	case "bind":
		return &object.Builtin{
			Name: "bind",
			Fn: func(args ...object.Object) object.Object {
				// Return a new builtin that calls the original with bound args
				boundArgs := args[1:] // Skip thisArg
				return &object.Builtin{
					Name: fn.Name,
					Fn: func(callArgs ...object.Object) object.Object {
						allArgs := append(boundArgs, callArgs...)
						return fn.Fn(allArgs...)
					},
				}
			},
		}
	case "name":
		return &object.String{Value: fn.Name}
	case "length":
		return &object.Number{Value: 0} // Builtins don't track arity
	case "toString":
		return &object.Builtin{
			Name: "toString",
			Fn: func(args ...object.Object) object.Object {
				return &object.String{Value: fmt.Sprintf("function %s() { [native code] }", fn.Name)}
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
		
		// Handle Function objects (allow setting properties like func.VERSION = "1.0")
		if fn, ok := obj.(*object.Function); ok {
			fn.Set(propName, right)
			return right
		}
		
		// Handle Builtin objects (allow setting properties like lodash.prototype = ...)
		if _, ok := obj.(*object.Builtin); ok {
			// Builtins don't support property assignment, but we should not error
			// Just silently ignore - this matches JS behavior for native functions
			return right
		}
		
		// Fallback - for any other types, return error with type info for debugging
		return newError("cannot assign to property '%s' of %s", propName, obj.Type())

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
	return newError("invalid assignment target: %T for %s", ae.Left, ae.Left.String())
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
		case "&=":
			return &object.Number{Value: float64(int32(oldNum.Value) & int32(rightNum.Value))}
		case "|=":
			return &object.Number{Value: float64(int32(oldNum.Value) | int32(rightNum.Value))}
		case "^=":
			return &object.Number{Value: float64(int32(oldNum.Value) ^ int32(rightNum.Value))}
		case "<<=":
			return &object.Number{Value: float64(int32(oldNum.Value) << uint32(rightNum.Value))}
		case ">>=":
			return &object.Number{Value: float64(int32(oldNum.Value) >> uint32(rightNum.Value))}
		case ">>>=":
			return &object.Number{Value: float64(uint32(oldNum.Value) >> uint32(rightNum.Value))}
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
