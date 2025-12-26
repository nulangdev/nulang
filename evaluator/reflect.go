// Package evaluator implements the Reflect API for Nulang.
package evaluator

import (
	"fmt"
	"sort"

	"github.com/nulang/nulang/object"
)

// initReflect initializes the Reflect global object
func initReflect() *object.ObjectMap {
	reflectObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Reflect.get(target, property [, receiver])
	reflectObj.Set("get", &object.Builtin{
		Name: "get",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.get requires at least 2 arguments")
			}
			target := args[0]
			property := args[1]

			propStr := objectToString(property)

			// Check if target is a Proxy
			if proxy, ok := target.(*ProxyObject); ok {
				return ProxyGet(proxy, propStr, nil)
			}

			// Regular object access
			if objMap, ok := target.(*object.ObjectMap); ok {
				if val, found := objMap.Get(propStr); found {
					return val
				}
				return UNDEFINED
			}

			if arr, ok := target.(*object.Array); ok {
				if propStr == "length" {
					return &object.Number{Value: float64(len(arr.Elements))}
				}
				// Try as index
				idx := 0
				if _, err := fmt.Sscanf(propStr, "%d", &idx); err == nil {
					if idx >= 0 && idx < len(arr.Elements) {
						return arr.Elements[idx]
					}
				}
			}

			return UNDEFINED
		},
	})

	// Reflect.set(target, property, value [, receiver])
	reflectObj.Set("set", &object.Builtin{
		Name: "set",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 3 {
				return newError("Reflect.set requires at least 3 arguments")
			}
			target := args[0]
			property := args[1]
			value := args[2]

			propStr := objectToString(property)

			// Check if target is a Proxy
			if proxy, ok := target.(*ProxyObject); ok {
				return ProxySet(proxy, propStr, value, nil)
			}

			// Regular object set
			if objMap, ok := target.(*object.ObjectMap); ok {
				objMap.Set(propStr, value)
				return TRUE
			}

			return FALSE
		},
	})

	// Reflect.has(target, property)
	reflectObj.Set("has", &object.Builtin{
		Name: "has",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.has requires 2 arguments")
			}
			target := args[0]
			property := args[1]

			propStr := objectToString(property)

			// Check if target is a Proxy
			if proxy, ok := target.(*ProxyObject); ok {
				return ProxyHas(proxy, propStr, nil)
			}

			// Regular object check
			if objMap, ok := target.(*object.ObjectMap); ok {
				if _, found := objMap.Get(propStr); found {
					return TRUE
				}
			}

			if arr, ok := target.(*object.Array); ok {
				if propStr == "length" {
					return TRUE
				}
				idx := 0
				if _, err := fmt.Sscanf(propStr, "%d", &idx); err == nil {
					if idx >= 0 && idx < len(arr.Elements) {
						return TRUE
					}
				}
			}

			return FALSE
		},
	})

	// Reflect.deleteProperty(target, property)
	reflectObj.Set("deleteProperty", &object.Builtin{
		Name: "deleteProperty",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.deleteProperty requires 2 arguments")
			}
			target := args[0]
			property := args[1]

			propStr := objectToString(property)

			if objMap, ok := target.(*object.ObjectMap); ok {
				delete(objMap.Pairs, propStr)
				return TRUE
			}

			return FALSE
		},
	})

	// Reflect.ownKeys(target)
	reflectObj.Set("ownKeys", &object.Builtin{
		Name: "ownKeys",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("Reflect.ownKeys requires 1 argument")
			}
			target := args[0]

			keys := []object.Object{}

			if objMap, ok := target.(*object.ObjectMap); ok {
				keyNames := make([]string, 0, len(objMap.Pairs))
				for key := range objMap.Pairs {
					keyNames = append(keyNames, key)
				}
				sort.Strings(keyNames)
				for _, key := range keyNames {
					keys = append(keys, &object.String{Value: key})
				}
			}

			return &object.Array{Elements: keys}
		},
	})

	// Reflect.apply(target, thisArgument, argumentsList)
	reflectObj.Set("apply", &object.Builtin{
		Name: "apply",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 3 {
				return newError("Reflect.apply requires 3 arguments")
			}
			fn := args[0]
			thisArg := args[1]
			argList := args[2]

			fnArgs := []object.Object{}
			if arr, ok := argList.(*object.Array); ok {
				fnArgs = arr.Elements
			}

			switch function := fn.(type) {
			case *object.Function:
				env := object.NewEnclosedEnvironment(function.Env)
				if objMap, ok := thisArg.(*object.ObjectMap); ok {
					env.Set("this", objMap)
				}
				for i, param := range function.Parameters {
					if i < len(fnArgs) {
						env.Set(param.Value, fnArgs[i])
					} else {
						env.Set(param.Value, UNDEFINED)
					}
				}
				evaluated := Eval(function.Body, env)
				return unwrapReturnValue(evaluated)
			case *object.Builtin:
				return function.Fn(fnArgs...)
			default:
				return newError("Reflect.apply: first argument must be a function")
			}
		},
	})

	// Reflect.construct(target, argumentsList [, newTarget])
	reflectObj.Set("construct", &object.Builtin{
		Name: "construct",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.construct requires at least 2 arguments")
			}
			target := args[0]
			argList := args[1]

			fnArgs := []object.Object{}
			if arr, ok := argList.(*object.Array); ok {
				fnArgs = arr.Elements
			}

			// Check if it's a Class
			if class, ok := target.(*Class); ok {
				return createClassInstance(class, fnArgs, class.Env)
			}

			// Check if it's a function (constructor function pattern)
			if fn, ok := target.(*object.Function); ok {
				instance := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
				env := object.NewEnclosedEnvironment(fn.Env)
				env.Set("this", instance)
				for i, param := range fn.Parameters {
					if i < len(fnArgs) {
						env.Set(param.Value, fnArgs[i])
					} else {
						env.Set(param.Value, UNDEFINED)
					}
				}
				Eval(fn.Body, env)
				return instance
			}

			return newError("Reflect.construct: target must be a constructor")
		},
	})

	// Reflect.getPrototypeOf(target)
	reflectObj.Set("getPrototypeOf", &object.Builtin{
		Name: "getPrototypeOf",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("Reflect.getPrototypeOf requires 1 argument")
			}
			target := args[0]

			if objMap, ok := target.(*object.ObjectMap); ok {
				if objMap.Prototype != nil {
					return objMap.Prototype
				}
			}

			return NULL
		},
	})

	// Reflect.setPrototypeOf(target, prototype)
	reflectObj.Set("setPrototypeOf", &object.Builtin{
		Name: "setPrototypeOf",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.setPrototypeOf requires 2 arguments")
			}
			target := args[0]
			prototype := args[1]

			if objMap, ok := target.(*object.ObjectMap); ok {
				if protoMap, ok := prototype.(*object.ObjectMap); ok {
					objMap.Prototype = protoMap
					return TRUE
				}
				if prototype.Type() == object.NULL_OBJ {
					objMap.Prototype = nil
					return TRUE
				}
			}

			return FALSE
		},
	})

	// Reflect.isExtensible(target)
	reflectObj.Set("isExtensible", &object.Builtin{
		Name: "isExtensible",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("Reflect.isExtensible requires 1 argument")
			}
			// In our implementation, objects are always extensible unless frozen
			return TRUE
		},
	})

	// Reflect.preventExtensions(target)
	reflectObj.Set("preventExtensions", &object.Builtin{
		Name: "preventExtensions",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("Reflect.preventExtensions requires 1 argument")
			}
			// We'll just return true; proper implementation would need object flags
			return TRUE
		},
	})

	// Reflect.getOwnPropertyDescriptor(target, property)
	reflectObj.Set("getOwnPropertyDescriptor", &object.Builtin{
		Name: "getOwnPropertyDescriptor",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Reflect.getOwnPropertyDescriptor requires 2 arguments")
			}
			target := args[0]
			property := args[1]

			propStr := objectToString(property)

			if objMap, ok := target.(*object.ObjectMap); ok {
				if val, found := objMap.Get(propStr); found {
					descriptor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
					descriptor.Set("value", val)
					descriptor.Set("writable", TRUE)
					descriptor.Set("enumerable", TRUE)
					descriptor.Set("configurable", TRUE)
					return descriptor
				}
			}

			return UNDEFINED
		},
	})

	// Reflect.defineProperty(target, property, descriptor)
	reflectObj.Set("defineProperty", &object.Builtin{
		Name: "defineProperty",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 3 {
				return newError("Reflect.defineProperty requires 3 arguments")
			}
			target := args[0]
			property := args[1]
			descriptor := args[2]

			propStr := objectToString(property)

			if objMap, ok := target.(*object.ObjectMap); ok {
				if descMap, ok := descriptor.(*object.ObjectMap); ok {
					if val, found := descMap.Get("value"); found {
						objMap.Set(propStr, val)
						return TRUE
					}
				}
			}

			return FALSE
		},
	})

	return reflectObj
}
