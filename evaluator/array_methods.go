package evaluator

import (
	"github.com/nulang/nulang/object"
)

func evalArrayProperty(arr *object.Array, prop string) object.Object {
	switch prop {
	case "length":
		return &object.Number{Value: float64(len(arr.Elements))}
	case "push":
		return &object.Builtin{Name: "push", Fn: func(args ...object.Object) object.Object {
			arr.Elements = append(arr.Elements, args...)
			return &object.Number{Value: float64(len(arr.Elements))}
		}}
	case "pop":
		return &object.Builtin{Name: "pop", Fn: func(args ...object.Object) object.Object {
			if len(arr.Elements) == 0 {
				return UNDEFINED
			}
			elem := arr.Elements[len(arr.Elements)-1]
			arr.Elements = arr.Elements[:len(arr.Elements)-1]
			return elem
		}}
	case "shift":
		return &object.Builtin{Name: "shift", Fn: func(args ...object.Object) object.Object {
			if len(arr.Elements) == 0 {
				return UNDEFINED
			}
			elem := arr.Elements[0]
			arr.Elements = arr.Elements[1:]
			return elem
		}}
	case "unshift":
		return &object.Builtin{Name: "unshift", Fn: func(args ...object.Object) object.Object {
			arr.Elements = append(args, arr.Elements...)
			return &object.Number{Value: float64(len(arr.Elements))}
		}}
	case "map":
		return createArrayIterator(arr, "map")
	case "filter":
		return createArrayIterator(arr, "filter")
	case "forEach":
		return createArrayIterator(arr, "forEach")
	case "find":
		return createArrayIterator(arr, "find")
	case "findIndex":
		return createArrayIterator(arr, "findIndex")
	case "reduce":
		return createArrayReduce(arr)
	case "includes":
		return &object.Builtin{Name: "includes", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			for _, elem := range arr.Elements {
				if evalStrictEquality(elem, args[0]) == TRUE {
					return TRUE
				}
			}
			return FALSE
		}}
	case "join":
		return &object.Builtin{Name: "join", Fn: func(args ...object.Object) object.Object {
			sep := ","
			if len(args) > 0 {
				if s, ok := args[0].(*object.String); ok {
					sep = s.Value
				}
			}
			result := ""
			for i, elem := range arr.Elements {
				if i > 0 {
					result += sep
				}
				result += objectToString(elem)
			}
			return &object.String{Value: result}
		}}
	case "reverse":
		return &object.Builtin{Name: "reverse", Fn: func(args ...object.Object) object.Object {
			n := len(arr.Elements)
			for i := 0; i < n/2; i++ {
				arr.Elements[i], arr.Elements[n-1-i] = arr.Elements[n-1-i], arr.Elements[i]
			}
			return arr
		}}
	case "slice":
		return createArraySlice(arr)
	case "concat":
		return &object.Builtin{Name: "concat", Fn: func(args ...object.Object) object.Object {
			newElements := make([]object.Object, len(arr.Elements))
			copy(newElements, arr.Elements)
			for _, arg := range args {
				if a, ok := arg.(*object.Array); ok {
					newElements = append(newElements, a.Elements...)
				} else {
					newElements = append(newElements, arg)
				}
			}
			return &object.Array{Elements: newElements}
		}}
	}
	return UNDEFINED
}

func createArrayIterator(arr *object.Array, method string) *object.Builtin {
	return &object.Builtin{Name: method, Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("%s requires a callback function", method)
		}
		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("%s callback must be a function", method)
		}

		switch method {
		case "map":
			result := make([]object.Object, len(arr.Elements))
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				result[i] = unwrapReturnValue(Eval(fn.Body, fnEnv))
			}
			return &object.Array{Elements: result}
		case "filter":
			var result []object.Object
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					result = append(result, elem)
				}
			}
			return &object.Array{Elements: result}
		case "forEach":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				Eval(fn.Body, fnEnv)
			}
			return UNDEFINED
		case "find":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					return elem
				}
			}
			return UNDEFINED
		case "findIndex":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					return &object.Number{Value: float64(i)}
				}
			}
			return &object.Number{Value: -1}
		}
		return UNDEFINED
	}}
}

func createArrayReduce(arr *object.Array) *object.Builtin {
	return &object.Builtin{Name: "reduce", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("reduce requires a callback function")
		}
		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("reduce callback must be a function")
		}
		var acc object.Object
		startIdx := 0
		if len(args) > 1 {
			acc = args[1]
		} else if len(arr.Elements) > 0 {
			acc = arr.Elements[0]
			startIdx = 1
		} else {
			return newError("reduce of empty array with no initial value")
		}
		for i := startIdx; i < len(arr.Elements); i++ {
			fnEnv := extendFunctionEnv(fn, []object.Object{acc, arr.Elements[i], &object.Number{Value: float64(i)}, arr})
			acc = unwrapReturnValue(Eval(fn.Body, fnEnv))
		}
		return acc
	}}
}

func createArraySlice(arr *object.Array) *object.Builtin {
	return &object.Builtin{Name: "slice", Fn: func(args ...object.Object) object.Object {
		start, end := 0, len(arr.Elements)
		if len(args) > 0 {
			if n, ok := args[0].(*object.Number); ok {
				start = int(n.Value)
				if start < 0 {
					start = len(arr.Elements) + start
				}
			}
		}
		if len(args) > 1 {
			if n, ok := args[1].(*object.Number); ok {
				end = int(n.Value)
				if end < 0 {
					end = len(arr.Elements) + end
				}
			}
		}
		if start < 0 {
			start = 0
		}
		if end > len(arr.Elements) {
			end = len(arr.Elements)
		}
		if start >= end {
			return &object.Array{Elements: []object.Object{}}
		}
		newElements := make([]object.Object, end-start)
		copy(newElements, arr.Elements[start:end])
		return &object.Array{Elements: newElements}
	}}
}
