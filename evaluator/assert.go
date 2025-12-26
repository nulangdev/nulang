package evaluator

import (
	"fmt"

	"github.com/nulang/nulang/object"
)

// initAssertModule initializes the assert module
func initAssertModule() *object.ObjectMap {
	assert := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// assert(value, message?)
	assertFn := &object.Builtin{Name: "assert", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("AssertionError: No value provided")
		}

		if !isTruthy(args[0]) {
			message := "Assertion failed"
			if len(args) > 1 {
				message = objectToString(args[1])
			}
			return newError("AssertionError: %s", message)
		}

		return UNDEFINED
	}}

	assert.Set("__call__", assertFn)

	// assert.ok(value, message?)
	assert.Set("ok", assertFn)

	// assert.equal(actual, expected, message?)
	assert.Set("equal", &object.Builtin{Name: "equal", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.equal requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if !isLooselyEqual(actual, expected) {
			message := "Expected values to be equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.notEqual(actual, expected, message?)
	assert.Set("notEqual", &object.Builtin{Name: "notEqual", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.notEqual requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if isLooselyEqual(actual, expected) {
			message := "Expected values to be not equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected not: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.strictEqual(actual, expected, message?)
	assert.Set("strictEqual", &object.Builtin{Name: "strictEqual", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.strictEqual requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if !isStrictlyEqual(actual, expected) {
			message := "Expected values to be strictly equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.notStrictEqual(actual, expected, message?)
	assert.Set("notStrictEqual", &object.Builtin{Name: "notStrictEqual", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.notStrictEqual requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if isStrictlyEqual(actual, expected) {
			message := "Expected values to be not strictly equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected not: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.deepEqual(actual, expected, message?)
	assert.Set("deepEqual", &object.Builtin{Name: "deepEqual", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.deepEqual requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if !isDeepEqual(actual, expected) {
			message := "Expected values to be deeply equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.notDeepEqual(actual, expected, message?)
	assert.Set("notDeepEqual", &object.Builtin{Name: "notDeepEqual", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("AssertionError: assert.notDeepEqual requires two arguments")
		}

		actual := args[0]
		expected := args[1]

		if isDeepEqual(actual, expected) {
			message := "Expected values to be not deeply equal"
			if len(args) > 2 {
				message = objectToString(args[2])
			}
			return newError("AssertionError: %s\nActual: %s\nExpected not: %s",
				message, actual.Inspect(), expected.Inspect())
		}

		return UNDEFINED
	}})

	// assert.throws(fn, error?, message?)
	assert.Set("throws", &object.Builtin{Name: "throws", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("AssertionError: assert.throws requires a function")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			if builtin, ok := args[0].(*object.Builtin); ok {
				result := builtin.Fn()
				if !isError(result) {
					return newError("AssertionError: Function did not throw")
				}
				return UNDEFINED
			}
			return newError("AssertionError: First argument must be a function")
		}

		fnEnv := object.NewEnvironment()
		result := Eval(fn.Body, fnEnv)

		if !isError(result) {
			message := "Function did not throw"
			if len(args) > 1 {
				message = objectToString(args[len(args)-1])
			}
			return newError("AssertionError: %s", message)
		}

		return UNDEFINED
	}})

	// assert.doesNotThrow(fn, message?)
	assert.Set("doesNotThrow", &object.Builtin{Name: "doesNotThrow", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("AssertionError: assert.doesNotThrow requires a function")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			if builtin, ok := args[0].(*object.Builtin); ok {
				result := builtin.Fn()
				if isError(result) {
					return newError("AssertionError: Function threw: %s", result.Inspect())
				}
				return UNDEFINED
			}
			return newError("AssertionError: First argument must be a function")
		}

		fnEnv := object.NewEnvironment()
		result := Eval(fn.Body, fnEnv)

		if isError(result) {
			message := "Function threw"
			if len(args) > 1 {
				message = objectToString(args[1])
			}
			return newError("AssertionError: %s: %s", message, result.Inspect())
		}

		return UNDEFINED
	}})

	// assert.ifError(value)
	assert.Set("ifError", &object.Builtin{Name: "ifError", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}

		if !isNullOrUndefined(args[0]) {
			return newError("AssertionError: ifError got unwanted error: %s", args[0].Inspect())
		}

		return UNDEFINED
	}})

	// assert.fail(message?)
	assert.Set("fail", &object.Builtin{Name: "fail", Fn: func(args ...object.Object) object.Object {
		message := "Failed"
		if len(args) > 0 {
			message = objectToString(args[0])
		}
		return newError("AssertionError: %s", message)
	}})

	return assert
}

// isLooselyEqual checks for == equality
func isLooselyEqual(a, b object.Object) bool {
	// Same type check
	if a.Type() == b.Type() {
		return isStrictlyEqual(a, b)
	}

	// null == undefined
	if (a.Type() == object.NULL_OBJ && b.Type() == object.UNDEFINED_OBJ) ||
		(a.Type() == object.UNDEFINED_OBJ && b.Type() == object.NULL_OBJ) {
		return true
	}

	// Number comparisons
	aNum, aIsNum := a.(*object.Number)
	bNum, bIsNum := b.(*object.Number)

	if aIsNum && bIsNum {
		return aNum.Value == bNum.Value
	}

	// String to number
	if aIsNum {
		if bStr, ok := b.(*object.String); ok {
			var val float64
			_, err := fmt.Sscanf(bStr.Value, "%f", &val)
			return err == nil && aNum.Value == val
		}
	}
	if bIsNum {
		if aStr, ok := a.(*object.String); ok {
			var val float64
			_, err := fmt.Sscanf(aStr.Value, "%f", &val)
			return err == nil && bNum.Value == val
		}
	}

	return false
}

// isStrictlyEqual checks for === equality
func isStrictlyEqual(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *object.Number:
		return a.Value == b.(*object.Number).Value
	case *object.String:
		return a.Value == b.(*object.String).Value
	case *object.Boolean:
		return a.Value == b.(*object.Boolean).Value
	case *object.Null:
		return true
	case *object.Undefined:
		return true
	default:
		// Reference equality for objects
		return a == b
	}
}

// isDeepEqual checks for deep equality
func isDeepEqual(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *object.Number:
		return a.Value == b.(*object.Number).Value
	case *object.String:
		return a.Value == b.(*object.String).Value
	case *object.Boolean:
		return a.Value == b.(*object.Boolean).Value
	case *object.Null, *object.Undefined:
		return true
	case *object.Array:
		bArr := b.(*object.Array)
		if len(a.Elements) != len(bArr.Elements) {
			return false
		}
		for i := range a.Elements {
			if !isDeepEqual(a.Elements[i], bArr.Elements[i]) {
				return false
			}
		}
		return true
	case *object.ObjectMap:
		bMap := b.(*object.ObjectMap)
		if len(a.Pairs) != len(bMap.Pairs) {
			return false
		}
		for key, aPair := range a.Pairs {
			bPair, ok := bMap.Pairs[key]
			if !ok || !isDeepEqual(aPair.Value, bPair.Value) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
