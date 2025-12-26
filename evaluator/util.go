package evaluator

import (
	"fmt"
	"strings"

	"github.com/nulang/nulang/object"
)

// initUtilModule initializes the util module
func initUtilModule() *object.ObjectMap {
	util := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// util.promisify(fn)
	util.Set("promisify", &object.Builtin{Name: "promisify", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("promisify requires a function argument")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("promisify argument must be a function")
		}

		// Return a new function that returns a Promise
		return &object.Builtin{Name: "promisified", Fn: func(promiseArgs ...object.Object) object.Object {
			promise := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			promise.Set("_promiseState", &object.String{Value: "pending"})

			// Execute the original function with a callback
			callbackArgs := make([]object.Object, len(promiseArgs)+1)
			copy(callbackArgs, promiseArgs)

			// Add callback as last argument
			callbackArgs[len(callbackArgs)-1] = &object.Builtin{
				Name: "callback",
				Fn: func(cbArgs ...object.Object) object.Object {
					if len(cbArgs) > 0 {
						// First argument is error
						if !isNullOrUndefined(cbArgs[0]) {
							promise.Set("_promiseState", &object.String{Value: "rejected"})
							promise.Set("_promiseValue", cbArgs[0])
						} else if len(cbArgs) > 1 {
							// Second argument is result
							promise.Set("_promiseState", &object.String{Value: "fulfilled"})
							promise.Set("_promiseValue", cbArgs[1])
						}
					}
					return UNDEFINED
				},
			}

			// Call original function
			fnEnv := extendFunctionEnv(fn, callbackArgs)
			Eval(fn.Body, fnEnv)

			return promise
		}}
	}})

	// util.inspect(object, options?)
	util.Set("inspect", &object.Builtin{Name: "inspect", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: ""}
		}

		depth := 2
		colors := false
		showHidden := false

		// Parse options
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				if d, ok := opts.Get("depth"); ok {
					if num, ok := d.(*object.Number); ok {
						depth = int(num.Value)
					}
				}
				if c, ok := opts.Get("colors"); ok {
					if b, ok := c.(*object.Boolean); ok {
						colors = b.Value
					}
				}
				if sh, ok := opts.Get("showHidden"); ok {
					if b, ok := sh.(*object.Boolean); ok {
						showHidden = b.Value
					}
				}
			}
		}

		result := inspectObject(args[0], depth, 0, colors, showHidden)
		return &object.String{Value: result}
	}})

	// util.format(format, ...args)
	util.Set("format", &object.Builtin{Name: "format", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return &object.String{Value: ""}
		}

		format := objectToString(args[0])
		formatArgs := args[1:]

		result := formatString(format, formatArgs)
		return &object.String{Value: result}
	}})

	// util.types
	types := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	types.Set("isArray", &object.Builtin{Name: "isArray", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		_, ok := args[0].(*object.Array)
		return nativeBoolToBooleanObject(ok)
	}})

	types.Set("isDate", &object.Builtin{Name: "isDate", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		// Check if it's a Date object (has _date property)
		if obj, ok := args[0].(*object.ObjectMap); ok {
			_, hasDate := obj.Get("_date")
			return nativeBoolToBooleanObject(hasDate)
		}
		return FALSE
	}})

	types.Set("isPromise", &object.Builtin{Name: "isPromise", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		// Check if it's a Promise object (has _promiseState property)
		if obj, ok := args[0].(*object.ObjectMap); ok {
			_, hasState := obj.Get("_promiseState")
			return nativeBoolToBooleanObject(hasState)
		}
		return FALSE
	}})

	types.Set("isRegExp", &object.Builtin{Name: "isRegExp", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		// Check if it's a RegExp object (has _regex property)
		if obj, ok := args[0].(*object.ObjectMap); ok {
			_, hasRegex := obj.Get("_regex")
			return nativeBoolToBooleanObject(hasRegex)
		}
		return FALSE
	}})

	types.Set("isMap", &object.Builtin{Name: "isMap", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		// Check if it's a Map object (has _map property)
		if obj, ok := args[0].(*object.ObjectMap); ok {
			if mapVal, ok := obj.Get("_map"); ok {
				_, isMap := mapVal.(*NuMap)
				return nativeBoolToBooleanObject(isMap)
			}
		}
		return FALSE
	}})

	types.Set("isSet", &object.Builtin{Name: "isSet", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		// Check if it's a Set object (has _set property)
		if obj, ok := args[0].(*object.ObjectMap); ok {
			if setVal, ok := obj.Get("_set"); ok {
				_, isSet := setVal.(*NuSet)
				return nativeBoolToBooleanObject(isSet)
			}
		}
		return FALSE
	}})

	types.Set("isError", &object.Builtin{Name: "isError", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		return nativeBoolToBooleanObject(isErrorInstance(args[0]))
	}})

	types.Set("isBuffer", &object.Builtin{Name: "isBuffer", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		_, ok := args[0].(*object.Buffer)
		return nativeBoolToBooleanObject(ok)
	}})

	util.Set("types", types)

	// util.inherits(constructor, superConstructor)
	util.Set("inherits", &object.Builtin{Name: "inherits", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("inherits requires two constructor arguments")
		}
		// This is mainly for compatibility - our class system handles inheritance differently
		return UNDEFINED
	}})

	// util.callbackify(async) - opposite of promisify
	util.Set("callbackify", &object.Builtin{Name: "callbackify", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("callbackify requires a function argument")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("callbackify argument must be a function")
		}

		// Return a new function that accepts a callback
		return &object.Builtin{Name: "callbackified", Fn: func(cbArgs ...object.Object) object.Object {
			if len(cbArgs) == 0 {
				return newError("callbackified function requires a callback")
			}

			callback := cbArgs[len(cbArgs)-1]
			actualArgs := cbArgs[:len(cbArgs)-1]

			// Execute the async function
			fnEnv := extendFunctionEnv(fn, actualArgs)
			result := unwrapReturnValue(Eval(fn.Body, fnEnv))

			// If result is a promise, handle it
			if promiseObj, ok := result.(*object.ObjectMap); ok {
				if stateVal, ok := promiseObj.Get("_promiseState"); ok {
					if stateStr, ok := stateVal.(*object.String); ok {
						if stateStr.Value == "fulfilled" {
							if val, ok := promiseObj.Get("_promiseValue"); ok {
								if cb, ok := callback.(*object.Function); ok {
									cbEnv := extendFunctionEnv(cb, []object.Object{NULL, val})
									return Eval(cb.Body, cbEnv)
								} else if cb, ok := callback.(*object.Builtin); ok {
									return cb.Fn(NULL, val)
								}
							}
						} else if stateStr.Value == "rejected" {
							if val, ok := promiseObj.Get("_promiseValue"); ok {
								if cb, ok := callback.(*object.Function); ok {
									cbEnv := extendFunctionEnv(cb, []object.Object{val})
									return Eval(cb.Body, cbEnv)
								} else if cb, ok := callback.(*object.Builtin); ok {
									return cb.Fn(val)
								}
							}
						}
					}
				}
			} else {
				// Call callback with result
				if cb, ok := callback.(*object.Function); ok {
					cbEnv := extendFunctionEnv(cb, []object.Object{NULL, result})
					Eval(cb.Body, cbEnv)
				} else if cb, ok := callback.(*object.Builtin); ok {
					cb.Fn(NULL, result)
				}
			}

			return UNDEFINED
		}}
	}})

	return util
}

// inspectObject recursively inspects an object
func inspectObject(obj object.Object, maxDepth, currentDepth int, colors, showHidden bool) string {
	if currentDepth > maxDepth {
		return "[Object]"
	}

	switch o := obj.(type) {
	case *object.String:
		if colors {
			return fmt.Sprintf("\x1b[32m'%s'\x1b[0m", o.Value)
		}
		return fmt.Sprintf("'%s'", o.Value)
	case *object.Number:
		if colors {
			return fmt.Sprintf("\x1b[33m%g\x1b[0m", o.Value)
		}
		return fmt.Sprintf("%g", o.Value)
	case *object.Boolean:
		if colors {
			return fmt.Sprintf("\x1b[33m%t\x1b[0m", o.Value)
		}
		return fmt.Sprintf("%t", o.Value)
	case *object.Null:
		if colors {
			return "\x1b[1mnull\x1b[0m"
		}
		return "null"
	case *object.Undefined:
		if colors {
			return "\x1b[90mundefined\x1b[0m"
		}
		return "undefined"
	case *object.Array:
		parts := []string{}
		for _, elem := range o.Elements {
			parts = append(parts, inspectObject(elem, maxDepth, currentDepth+1, colors, showHidden))
		}
		return "[ " + strings.Join(parts, ", ") + " ]"
	case *object.ObjectMap:
		if len(o.Pairs) == 0 {
			return "{}"
		}
		parts := []string{}
		for key, pair := range o.Pairs {
			valStr := inspectObject(pair.Value, maxDepth, currentDepth+1, colors, showHidden)
			if colors {
				parts = append(parts, fmt.Sprintf("\x1b[36m%s\x1b[0m: %s", key, valStr))
			} else {
				parts = append(parts, fmt.Sprintf("%s: %s", key, valStr))
			}
		}
		return "{ " + strings.Join(parts, ", ") + " }"
	case *object.Function:
		if colors {
			return "\x1b[36m[Function]\x1b[0m"
		}
		return "[Function]"
	case *object.Builtin:
		if colors {
			return fmt.Sprintf("\x1b[36m[Function: %s]\x1b[0m", o.Name)
		}
		return fmt.Sprintf("[Function: %s]", o.Name)
	default:
		return obj.Inspect()
	}
}

// formatString formats a string with placeholders
func formatString(format string, args []object.Object) string {
	result := format
	argIndex := 0

	for i := 0; i < len(result); i++ {
		if result[i] == '%' && i+1 < len(result) && argIndex < len(args) {
			switch result[i+1] {
			case 's': // string
				replacement := objectToString(args[argIndex])
				result = result[:i] + replacement + result[i+2:]
				i += len(replacement) - 1
				argIndex++
			case 'd', 'i': // integer
				if num, ok := args[argIndex].(*object.Number); ok {
					replacement := fmt.Sprintf("%d", int(num.Value))
					result = result[:i] + replacement + result[i+2:]
					i += len(replacement) - 1
				}
				argIndex++
			case 'f': // float
				if num, ok := args[argIndex].(*object.Number); ok {
					replacement := fmt.Sprintf("%f", num.Value)
					result = result[:i] + replacement + result[i+2:]
					i += len(replacement) - 1
				}
				argIndex++
			case 'j': // JSON
				jsonStr := stringify(args[argIndex])
				result = result[:i] + jsonStr + result[i+2:]
				i += len(jsonStr) - 1
				argIndex++
			case 'o', 'O': // object
				objStr := inspectObject(args[argIndex], 2, 0, false, false)
				result = result[:i] + objStr + result[i+2:]
				i += len(objStr) - 1
				argIndex++
			case '%': // literal %
				result = result[:i] + "%" + result[i+2:]
			}
		}
	}

	// Append remaining arguments
	for argIndex < len(args) {
		result += " " + objectToString(args[argIndex])
		argIndex++
	}

	return result
}

func isNullOrUndefined(obj object.Object) bool {
	switch obj.(type) {
	case *object.Null, *object.Undefined:
		return true
	default:
		return false
	}
}
