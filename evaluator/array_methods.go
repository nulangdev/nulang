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
			
			fromIndex := 0
			if len(args) > 1 {
				if n, ok := args[1].(*object.Number); ok {
					fromIndex = int(n.Value)
				}
			}
			if fromIndex < 0 {
				fromIndex = len(arr.Elements) + fromIndex
				if fromIndex < 0 {
					fromIndex = 0
				}
			}
			
			for i := fromIndex; i < len(arr.Elements); i++ {
				elem := arr.Elements[i]
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
	case "splice":
		return &object.Builtin{Name: "splice", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Array{Elements: []object.Object{}}
			}
			startNum, ok := args[0].(*object.Number)
			if !ok {
				return newError("splice start index must be a number")
			}
			start := int(startNum.Value)
			if start < 0 {
				start = len(arr.Elements) + start
				if start < 0 {
					start = 0
				}
			}
			if start > len(arr.Elements) {
				start = len(arr.Elements)
			}

			deleteCount := len(arr.Elements) - start
			if len(args) > 1 {
				if dn, ok := args[1].(*object.Number); ok {
					deleteCount = int(dn.Value)
					if deleteCount < 0 {
						deleteCount = 0
					}
					if start+deleteCount > len(arr.Elements) {
						deleteCount = len(arr.Elements) - start
					}
				}
			}

			itemsToInsert := []object.Object{}
			if len(args) > 2 {
				itemsToInsert = args[2:]
			}

			removed := make([]object.Object, deleteCount)
			copy(removed, arr.Elements[start:start+deleteCount])

			// Update array elements
			newElements := append(arr.Elements[:start], append(itemsToInsert, arr.Elements[start+deleteCount:]...)...)
			arr.Elements = newElements

			return &object.Array{Elements: removed}
		}}
	case "sort":
		return &object.Builtin{Name: "sort", Fn: func(args ...object.Object) object.Object {
			// In-place sort using simple bubble sort
			// Optional compareFn can be passed
			n := len(arr.Elements)
			
			// If no compare function, use default string comparison
			if len(args) == 0 {
				// Default sort: convert to strings and compare
				for i := 0; i < n-1; i++ {
					for j := 0; j < n-i-1; j++ {
						a := objectToString(arr.Elements[j])
						b := objectToString(arr.Elements[j+1])
						if a > b {
							arr.Elements[j], arr.Elements[j+1] = arr.Elements[j+1], arr.Elements[j]
						}
					}
				}
			} else if fn, ok := args[0].(*object.Function); ok {
				// Sort with compare function
				for i := 0; i < n-1; i++ {
					for j := 0; j < n-i-1; j++ {
						fnEnv := extendFunctionEnv(fn, []object.Object{arr.Elements[j], arr.Elements[j+1]})
						result := unwrapReturnValue(Eval(fn.Body, fnEnv))
						if num, ok := result.(*object.Number); ok {
							if num.Value > 0 {
								arr.Elements[j], arr.Elements[j+1] = arr.Elements[j+1], arr.Elements[j]
							}
						}
					}
				}
			} else if builtin, ok := args[0].(*object.Builtin); ok {
				// Sort with builtin compare function
				for i := 0; i < n-1; i++ {
					for j := 0; j < n-i-1; j++ {
						result := builtin.Fn(arr.Elements[j], arr.Elements[j+1])
						if num, ok := result.(*object.Number); ok {
							if num.Value > 0 {
								arr.Elements[j], arr.Elements[j+1] = arr.Elements[j+1], arr.Elements[j]
							}
						}
					}
				}
			}
			return arr
		}}
	case "indexOf":
		return &object.Builtin{Name: "indexOf", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: -1}
			}
			searchElement := args[0]
			fromIndex := 0
			if len(args) > 1 {
				if n, ok := args[1].(*object.Number); ok {
					fromIndex = int(n.Value)
				}
			}
			if fromIndex < 0 {
				fromIndex = len(arr.Elements) + fromIndex
				if fromIndex < 0 {
					fromIndex = 0
				}
			}

			for i := fromIndex; i < len(arr.Elements); i++ {
				if evalStrictEquality(arr.Elements[i], searchElement) == TRUE {
					return &object.Number{Value: float64(i)}
				}
			}
			return &object.Number{Value: -1}
		}}
	case "lastIndexOf":
		return &object.Builtin{Name: "lastIndexOf", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: -1}
			}
			searchElement := args[0]
			fromIndex := len(arr.Elements) - 1
			if len(args) > 1 {
				if n, ok := args[1].(*object.Number); ok {
					fromIndex = int(n.Value)
				}
			}
			if fromIndex < 0 {
				fromIndex = len(arr.Elements) + fromIndex
			}
			if fromIndex >= len(arr.Elements) {
				fromIndex = len(arr.Elements) - 1
			}

			for i := fromIndex; i >= 0; i-- {
				if evalStrictEquality(arr.Elements[i], searchElement) == TRUE {
					return &object.Number{Value: float64(i)}
				}
			}
			return &object.Number{Value: -1}
		}}
	case "some":
		return &object.Builtin{Name: "some", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			fn, ok := args[0].(*object.Function)
			if !ok {
				return newError("some callback must be a function")
			}
			var thisArg object.Object
			if len(args) > 1 {
				thisArg = args[1]
			}
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					return TRUE
				}
			}
			return FALSE
		}}
	case "every":
		return &object.Builtin{Name: "every", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return TRUE
			}
			fn, ok := args[0].(*object.Function)
			if !ok {
				return newError("every callback must be a function")
			}
			var thisArg object.Object
			if len(args) > 1 {
				thisArg = args[1]
			}
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				if !isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					return FALSE
				}
			}
			return TRUE
		}}
	case "fill":
		return &object.Builtin{Name: "fill", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return arr
			}
			value := args[0]
			start := 0
			end := len(arr.Elements)
			if len(args) > 1 {
				if n, ok := args[1].(*object.Number); ok {
					start = int(n.Value)
					if start < 0 {
						start = len(arr.Elements) + start
					}
				}
			}
			if len(args) > 2 {
				if n, ok := args[2].(*object.Number); ok {
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
			for i := start; i < end; i++ {
				arr.Elements[i] = value
			}
			return arr
		}}
	case "flat":
		return &object.Builtin{Name: "flat", Fn: func(args ...object.Object) object.Object {
			depth := 1
			if len(args) > 0 {
				if n, ok := args[0].(*object.Number); ok {
					depth = int(n.Value)
				}
			}
			result := flattenArray(arr.Elements, depth)
			return &object.Array{Elements: result}
		}}
	case "flatMap":
		return &object.Builtin{Name: "flatMap", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("flatMap requires a callback function")
			}
			fn, ok := args[0].(*object.Function)
			if !ok {
				return newError("flatMap callback must be a function")
			}
			var thisArg object.Object
			if len(args) > 1 {
				thisArg = args[1]
			}
			var result []object.Object
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				mapped := unwrapReturnValue(Eval(fn.Body, fnEnv))
				if mappedArr, ok := mapped.(*object.Array); ok {
					result = append(result, mappedArr.Elements...)
				} else {
					result = append(result, mapped)
				}
			}
			return &object.Array{Elements: result}
		}}
	case "copyWithin":
		return &object.Builtin{Name: "copyWithin", Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return arr
			}
			target := 0
			start := 0
			end := len(arr.Elements)
			
			if n, ok := args[0].(*object.Number); ok {
				target = int(n.Value)
				if target < 0 {
					target = len(arr.Elements) + target
				}
			}
			if n, ok := args[1].(*object.Number); ok {
				start = int(n.Value)
				if start < 0 {
					start = len(arr.Elements) + start
				}
			}
			if len(args) > 2 {
				if n, ok := args[2].(*object.Number); ok {
					end = int(n.Value)
					if end < 0 {
						end = len(arr.Elements) + end
					}
				}
			}
			
			// Copy elements
			count := end - start
			if target+count > len(arr.Elements) {
				count = len(arr.Elements) - target
			}
			for i := 0; i < count; i++ {
				arr.Elements[target+i] = arr.Elements[start+i]
			}
			return arr
		}}
	}
	return UNDEFINED
}

// Helper function to flatten arrays
func flattenArray(elements []object.Object, depth int) []object.Object {
	if depth == 0 {
		return elements
	}
	var result []object.Object
	for _, elem := range elements {
		if arr, ok := elem.(*object.Array); ok {
			result = append(result, flattenArray(arr.Elements, depth-1)...)
		} else {
			result = append(result, elem)
		}
	}
	return result
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

		var thisArg object.Object
		if len(args) > 1 {
			thisArg = args[1]
		}

		switch method {
		case "map":
			result := make([]object.Object, len(arr.Elements))
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				result[i] = unwrapReturnValue(Eval(fn.Body, fnEnv))
			}
			return &object.Array{Elements: result}
		case "filter":
			var result []object.Object
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					result = append(result, elem)
				}
			}
			return &object.Array{Elements: result}
		case "forEach":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				Eval(fn.Body, fnEnv)
			}
			return UNDEFINED
		case "find":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
				if isTruthy(unwrapReturnValue(Eval(fn.Body, fnEnv))) {
					return elem
				}
			}
			return UNDEFINED
		case "findIndex":
			for i, elem := range arr.Elements {
				fnEnv := extendFunctionEnv(fn, []object.Object{elem, &object.Number{Value: float64(i)}, arr})
				if thisArg != nil {
					fnEnv.Set("this", thisArg)
				}
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
