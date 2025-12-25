package evaluator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/nulang/nulang/object"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// builtins contains all global built-in objects and functions
var builtins map[string]object.Object

func initBuiltins() {
	if builtins != nil {
		return
	}
	builtins = make(map[string]object.Object)
	
	// Console
	console := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	console.Set("log", &object.Builtin{Name: "log", Fn: func(args ...object.Object) object.Object {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return UNDEFINED
	}})
	console.Set("error", &object.Builtin{Name: "error", Fn: func(args ...object.Object) object.Object {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return UNDEFINED
	}})
	console.Set("warn", &object.Builtin{Name: "warn", Fn: func(args ...object.Object) object.Object {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return UNDEFINED
	}})
	builtins["console"] = console
	
	// Math object
	mathObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	mathObj.Set("PI", &object.Number{Value: math.Pi})
	mathObj.Set("E", &object.Number{Value: math.E})
	mathObj.Set("abs", &object.Builtin{Name: "abs", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		if num, ok := args[0].(*object.Number); ok {
			return &object.Number{Value: math.Abs(num.Value)}
		}
		return &object.Number{Value: math.NaN()}
	}})
	mathObj.Set("ceil", &object.Builtin{Name: "ceil", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		if num, ok := args[0].(*object.Number); ok {
			return &object.Number{Value: math.Ceil(num.Value)}
		}
		return &object.Number{Value: math.NaN()}
	}})
	mathObj.Set("floor", &object.Builtin{Name: "floor", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		if num, ok := args[0].(*object.Number); ok {
			return &object.Number{Value: math.Floor(num.Value)}
		}
		return &object.Number{Value: math.NaN()}
	}})
	mathObj.Set("round", &object.Builtin{Name: "round", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		if num, ok := args[0].(*object.Number); ok {
			return &object.Number{Value: math.Round(num.Value)}
		}
		return &object.Number{Value: math.NaN()}
	}})
	mathObj.Set("max", &object.Builtin{Name: "max", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return &object.Number{Value: math.Inf(-1)}
		}
		max := math.Inf(-1)
		for _, arg := range args {
			if num, ok := arg.(*object.Number); ok && num.Value > max {
				max = num.Value
			}
		}
		return &object.Number{Value: max}
	}})
	mathObj.Set("min", &object.Builtin{Name: "min", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return &object.Number{Value: math.Inf(1)}
		}
		min := math.Inf(1)
		for _, arg := range args {
			if num, ok := arg.(*object.Number); ok && num.Value < min {
				min = num.Value
			}
		}
		return &object.Number{Value: min}
	}})
	mathObj.Set("random", &object.Builtin{Name: "random", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: rand.Float64()}
	}})
	mathObj.Set("pow", &object.Builtin{Name: "pow", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return &object.Number{Value: math.NaN()}
		}
		base, ok1 := args[0].(*object.Number)
		exp, ok2 := args[1].(*object.Number)
		if !ok1 || !ok2 {
			return &object.Number{Value: math.NaN()}
		}
		return &object.Number{Value: math.Pow(base.Value, exp.Value)}
	}})
	mathObj.Set("sqrt", &object.Builtin{Name: "sqrt", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		if num, ok := args[0].(*object.Number); ok {
			return &object.Number{Value: math.Sqrt(num.Value)}
		}
		return &object.Number{Value: math.NaN()}
	}})
	builtins["Math"] = mathObj
	
	// JSON object
	jsonObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	jsonObj.Set("stringify", &object.Builtin{Name: "stringify", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		return &object.String{Value: stringify(args[0])}
	}})
	jsonObj.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("JSON.parse requires a string argument")
		}
		return args[0]
	}})
	builtins["JSON"] = jsonObj
	
	// Array constructor
	arrayObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	arrayObj.Set("isArray", &object.Builtin{Name: "isArray", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		_, ok := args[0].(*object.Array)
		return nativeBoolToBooleanObject(ok)
	}})
	arrayObj.Set("from", &object.Builtin{Name: "from", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		if arr, ok := args[0].(*object.Array); ok {
			newElements := make([]object.Object, len(arr.Elements))
			copy(newElements, arr.Elements)
			return &object.Array{Elements: newElements}
		}
		if str, ok := args[0].(*object.String); ok {
			elements := make([]object.Object, len(str.Value))
			for i := 0; i < len(str.Value); i++ {
				elements[i] = &object.String{Value: string(str.Value[i])}
			}
			return &object.Array{Elements: elements}
		}
		return &object.Array{Elements: []object.Object{}}
	}})
	builtins["Array"] = arrayObj
	
	// Object constructor
	objectObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	objectObj.Set("keys", &object.Builtin{Name: "keys", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		if obj, ok := args[0].(*object.ObjectMap); ok {
			keys := make([]object.Object, 0, len(obj.Pairs))
			for key := range obj.Pairs {
				keys = append(keys, &object.String{Value: key})
			}
			return &object.Array{Elements: keys}
		}
		return &object.Array{Elements: []object.Object{}}
	}})
	objectObj.Set("values", &object.Builtin{Name: "values", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		if obj, ok := args[0].(*object.ObjectMap); ok {
			values := make([]object.Object, 0, len(obj.Pairs))
			for _, pair := range obj.Pairs {
				values = append(values, pair.Value)
			}
			return &object.Array{Elements: values}
		}
		return &object.Array{Elements: []object.Object{}}
	}})
	objectObj.Set("entries", &object.Builtin{Name: "entries", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		if obj, ok := args[0].(*object.ObjectMap); ok {
			entries := make([]object.Object, 0, len(obj.Pairs))
			for key, pair := range obj.Pairs {
				entry := &object.Array{Elements: []object.Object{&object.String{Value: key}, pair.Value}}
				entries = append(entries, entry)
			}
			return &object.Array{Elements: entries}
		}
		return &object.Array{Elements: []object.Object{}}
	}})
	builtins["Object"] = objectObj
	
	// Global functions
	builtins["print"] = &object.Builtin{Name: "print", Fn: func(args ...object.Object) object.Object {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
		return UNDEFINED
	}}
	builtins["parseInt"] = &object.Builtin{Name: "parseInt", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		str := objectToString(args[0])
		var result float64
		fmt.Sscanf(str, "%f", &result)
		return &object.Number{Value: float64(int(result))}
	}}
	builtins["parseFloat"] = &object.Builtin{Name: "parseFloat", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: math.NaN()}
		}
		str := objectToString(args[0])
		var result float64
		fmt.Sscanf(str, "%f", &result)
		return &object.Number{Value: result}
	}}
	builtins["isNaN"] = &object.Builtin{Name: "isNaN", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return TRUE
		}
		if num, ok := args[0].(*object.Number); ok {
			return nativeBoolToBooleanObject(math.IsNaN(num.Value))
		}
		return TRUE
	}}
	builtins["isFinite"] = &object.Builtin{Name: "isFinite", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		if num, ok := args[0].(*object.Number); ok {
			return nativeBoolToBooleanObject(!math.IsInf(num.Value, 0) && !math.IsNaN(num.Value))
		}
		return FALSE
	}}
	builtins["String"] = &object.Builtin{Name: "String", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: ""}
		}
		return &object.String{Value: objectToString(args[0])}
	}}
	builtins["Number"] = &object.Builtin{Name: "Number", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: 0}
		}
		if num, ok := args[0].(*object.Number); ok {
			return num
		}
		if str, ok := args[0].(*object.String); ok {
			var result float64
			_, err := fmt.Sscanf(str.Value, "%f", &result)
			if err != nil {
				return &object.Number{Value: math.NaN()}
			}
			return &object.Number{Value: result}
		}
		return &object.Number{Value: math.NaN()}
	}}
	builtins["Boolean"] = &object.Builtin{Name: "Boolean", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		return nativeBoolToBooleanObject(isTruthy(args[0]))
	}}
	
	// Built-in modules available globally
	builtins["fs"] = initFsModule()
	builtins["path"] = initPathModule()
	builtins["crypto"] = initCryptoModule()
	builtins["Promise"] = initPromiseConstructor()
	
	// Buffer constructor
	builtins["Buffer"] = initBufferConstructor()
	
	// require function for CommonJS-style imports
	builtins["require"] = &object.Builtin{Name: "require", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("require() requires a module path")
		}
		modulePath := objectToString(args[0])
		module, err := LoadModule(modulePath, CurrentModulePath)
		if err != nil {
			return newError("Cannot find module '%s': %s", modulePath, err.Error())
		}
		return module
	}}
}

func stringify(obj object.Object) string {
	switch o := obj.(type) {
	case *object.String:
		return fmt.Sprintf("\"%s\"", o.Value)
	case *object.Number:
		return fmt.Sprintf("%g", o.Value)
	case *object.Boolean:
		return fmt.Sprintf("%t", o.Value)
	case *object.Null:
		return "null"
	case *object.Array:
		parts := make([]string, len(o.Elements))
		for i, elem := range o.Elements {
			parts[i] = stringify(elem)
		}
		result := "["
		for i, part := range parts {
			if i > 0 {
				result += ","
			}
			result += part
		}
		return result + "]"
	case *object.ObjectMap:
		pairs := []string{}
		for key, pair := range o.Pairs {
			pairs = append(pairs, fmt.Sprintf("\"%s\":%s", key, stringify(pair.Value)))
		}
		result := "{"
		for i, pair := range pairs {
			if i > 0 {
				result += ","
			}
			result += pair
		}
		return result + "}"
	}
	return "undefined"
}
