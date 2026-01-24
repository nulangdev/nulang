package evaluator

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
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
	
	// Console - Full Node.js compatible implementation
	console := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	// Internal state for console
	consoleCounters := make(map[string]int)
	consoleTimers := make(map[string]time.Time)
	consoleGroupLevel := 0
	
	// Helper to get indentation based on group level
	getIndent := func() string {
		indent := ""
		for i := 0; i < consoleGroupLevel; i++ {
			indent += "  "
		}
		return indent
	}
	
	// Helper to print with indentation
	printWithIndent := func(args ...object.Object) {
		indent := getIndent()
		fmt.Print(indent)
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(arg.Inspect())
		}
		fmt.Println()
	}
	
	// console.log([data][, ...args])
	console.Set("log", &object.Builtin{Name: "log", Fn: func(args ...object.Object) object.Object {
		printWithIndent(args...)
		return UNDEFINED
	}})
	
	// console.info([data][, ...args]) - alias for log
	console.Set("info", &object.Builtin{Name: "info", Fn: func(args ...object.Object) object.Object {
		printWithIndent(args...)
		return UNDEFINED
	}})
	
	// console.debug([data][, ...args]) - alias for log
	console.Set("debug", &object.Builtin{Name: "debug", Fn: func(args ...object.Object) object.Object {
		printWithIndent(args...)
		return UNDEFINED
	}})
	
	// console.error([data][, ...args])
	console.Set("error", &object.Builtin{Name: "error", Fn: func(args ...object.Object) object.Object {
		indent := getIndent()
		fmt.Fprint(getStderr(), indent)
		for i, arg := range args {
			if i > 0 {
				fmt.Fprint(getStderr(), " ")
			}
			fmt.Fprint(getStderr(), arg.Inspect())
		}
		fmt.Fprintln(getStderr())
		return UNDEFINED
	}})
	
	// console.warn([data][, ...args])
	console.Set("warn", &object.Builtin{Name: "warn", Fn: func(args ...object.Object) object.Object {
		indent := getIndent()
		fmt.Fprint(getStderr(), indent)
		for i, arg := range args {
			if i > 0 {
				fmt.Fprint(getStderr(), " ")
			}
			fmt.Fprint(getStderr(), arg.Inspect())
		}
		fmt.Fprintln(getStderr())
		return UNDEFINED
	}})
	
	// console.assert(value[, ...message])
	console.Set("assert", &object.Builtin{Name: "assert", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return UNDEFINED
		}
		
		if !isTruthy(args[0]) {
			indent := getIndent()
			fmt.Fprint(getStderr(), indent, "Assertion failed:")
			if len(args) > 1 {
				for _, arg := range args[1:] {
					fmt.Fprint(getStderr(), " ", arg.Inspect())
				}
			}
			fmt.Fprintln(getStderr())
		}
		return UNDEFINED
	}})
	
	// console.clear()
	console.Set("clear", &object.Builtin{Name: "clear", Fn: func(args ...object.Object) object.Object {
		// ANSI escape codes to clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")
		return UNDEFINED
	}})
	
	// console.count([label])
	console.Set("count", &object.Builtin{Name: "count", Fn: func(args ...object.Object) object.Object {
		label := "default"
		if len(args) > 0 {
			label = objectToString(args[0])
		}
		consoleCounters[label]++
		indent := getIndent()
		fmt.Printf("%s%s: %d\n", indent, label, consoleCounters[label])
		return UNDEFINED
	}})
	
	// console.countReset([label])
	console.Set("countReset", &object.Builtin{Name: "countReset", Fn: func(args ...object.Object) object.Object {
		label := "default"
		if len(args) > 0 {
			label = objectToString(args[0])
		}
		if _, exists := consoleCounters[label]; exists {
			consoleCounters[label] = 0
		} else {
			fmt.Fprintf(getStderr(), "Count for '%s' does not exist\n", label)
		}
		return UNDEFINED
	}})
	
	// console.dir(obj[, options])
	console.Set("dir", &object.Builtin{Name: "dir", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return UNDEFINED
		}
		indent := getIndent()
		// For now, use Inspect() - could be enhanced with options later
		fmt.Printf("%s%s\n", indent, args[0].Inspect())
		return UNDEFINED
	}})
	
	// console.dirxml(...data) - same as dir for non-DOM environments
	console.Set("dirxml", &object.Builtin{Name: "dirxml", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return UNDEFINED
		}
		indent := getIndent()
		fmt.Printf("%s%s\n", indent, args[0].Inspect())
		return UNDEFINED
	}})
	
	// console.group([...label])
	console.Set("group", &object.Builtin{Name: "group", Fn: func(args ...object.Object) object.Object {
		if len(args) > 0 {
			printWithIndent(args...)
		}
		consoleGroupLevel++
		return UNDEFINED
	}})
	
	// console.groupCollapsed([...label]) - same as group in terminal
	console.Set("groupCollapsed", &object.Builtin{Name: "groupCollapsed", Fn: func(args ...object.Object) object.Object {
		if len(args) > 0 {
			printWithIndent(args...)
		}
		consoleGroupLevel++
		return UNDEFINED
	}})
	
	// console.groupEnd()
	console.Set("groupEnd", &object.Builtin{Name: "groupEnd", Fn: func(args ...object.Object) object.Object {
		if consoleGroupLevel > 0 {
			consoleGroupLevel--
		}
		return UNDEFINED
	}})
	
	// console.table(tabularData[, properties])
	console.Set("table", &object.Builtin{Name: "table", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			return UNDEFINED
		}
		
		indent := getIndent()
		data := args[0]
		
		switch d := data.(type) {
		case *object.Array:
			if len(d.Elements) == 0 {
				fmt.Println(indent + "┌─────────┐")
				fmt.Println(indent + "│ (empty) │")
				fmt.Println(indent + "└─────────┘")
				return UNDEFINED
			}
			
			// Calculate column widths
			indexWidth := len(fmt.Sprintf("%d", len(d.Elements)-1))
			if indexWidth < 5 {
				indexWidth = 5
			}
			valueWidth := 5
			for _, elem := range d.Elements {
				w := len(elem.Inspect())
				if w > valueWidth {
					valueWidth = w
				}
			}
			if valueWidth > 50 {
				valueWidth = 50
			}
			
			// Print table
			fmt.Printf("%s┌%s┬%s┐\n", indent, repeatChar('─', indexWidth+2), repeatChar('─', valueWidth+2))
			fmt.Printf("%s│ %-*s │ %-*s │\n", indent, indexWidth, "(idx)", valueWidth, "Values")
			fmt.Printf("%s├%s┼%s┤\n", indent, repeatChar('─', indexWidth+2), repeatChar('─', valueWidth+2))
			for i, elem := range d.Elements {
				val := truncateString(elem.Inspect(), valueWidth)
				fmt.Printf("%s│ %-*d │ %-*s │\n", indent, indexWidth, i, valueWidth, val)
			}
			fmt.Printf("%s└%s┴%s┘\n", indent, repeatChar('─', indexWidth+2), repeatChar('─', valueWidth+2))
			
		case *object.ObjectMap:
			if len(d.Pairs) == 0 {
				fmt.Println(indent + "┌─────────┐")
				fmt.Println(indent + "│ (empty) │")
				fmt.Println(indent + "└─────────┘")
				return UNDEFINED
			}
			
			// Calculate column widths
			keyWidth := 3
			valueWidth := 5
			for key, pair := range d.Pairs {
				if len(key) > keyWidth {
					keyWidth = len(key)
				}
				w := len(pair.Value.Inspect())
				if w > valueWidth {
					valueWidth = w
				}
			}
			if keyWidth > 30 {
				keyWidth = 30
			}
			if valueWidth > 50 {
				valueWidth = 50
			}
			
			// Print table
			fmt.Printf("%s┌%s┬%s┐\n", indent, repeatChar('─', keyWidth+2), repeatChar('─', valueWidth+2))
			fmt.Printf("%s│ %-*s │ %-*s │\n", indent, keyWidth, "Key", valueWidth, "Value")
			fmt.Printf("%s├%s┼%s┤\n", indent, repeatChar('─', keyWidth+2), repeatChar('─', valueWidth+2))
			for key, pair := range d.Pairs {
				k := truncateString(key, keyWidth)
				v := truncateString(pair.Value.Inspect(), valueWidth)
				fmt.Printf("%s│ %-*s │ %-*s │\n", indent, keyWidth, k, valueWidth, v)
			}
			fmt.Printf("%s└%s┴%s┘\n", indent, repeatChar('─', keyWidth+2), repeatChar('─', valueWidth+2))
			
		default:
			fmt.Printf("%s%s\n", indent, data.Inspect())
		}
		
		return UNDEFINED
	}})
	
	// console.time([label])
	console.Set("time", &object.Builtin{Name: "time", Fn: func(args ...object.Object) object.Object {
		label := "default"
		if len(args) > 0 {
			label = objectToString(args[0])
		}
		if _, exists := consoleTimers[label]; exists {
			fmt.Fprintf(getStderr(), "Timer '%s' already exists\n", label)
		} else {
			consoleTimers[label] = time.Now()
		}
		return UNDEFINED
	}})
	
	// console.timeEnd([label])
	console.Set("timeEnd", &object.Builtin{Name: "timeEnd", Fn: func(args ...object.Object) object.Object {
		label := "default"
		if len(args) > 0 {
			label = objectToString(args[0])
		}
		if startTime, exists := consoleTimers[label]; exists {
			elapsed := time.Since(startTime)
			indent := getIndent()
			fmt.Printf("%s%s: %.3fms\n", indent, label, float64(elapsed.Microseconds())/1000.0)
			delete(consoleTimers, label)
		} else {
			fmt.Fprintf(getStderr(), "Timer '%s' does not exist\n", label)
		}
		return UNDEFINED
	}})
	
	// console.timeLog([label][, ...data])
	console.Set("timeLog", &object.Builtin{Name: "timeLog", Fn: func(args ...object.Object) object.Object {
		label := "default"
		if len(args) > 0 {
			label = objectToString(args[0])
		}
		if startTime, exists := consoleTimers[label]; exists {
			elapsed := time.Since(startTime)
			indent := getIndent()
			fmt.Printf("%s%s: %.3fms", indent, label, float64(elapsed.Microseconds())/1000.0)
			// Print additional data
			if len(args) > 1 {
				for _, arg := range args[1:] {
					fmt.Print(" ", arg.Inspect())
				}
			}
			fmt.Println()
		} else {
			fmt.Fprintf(getStderr(), "Timer '%s' does not exist\n", label)
		}
		return UNDEFINED
	}})
	
	// console.trace([message][, ...args])
	console.Set("trace", &object.Builtin{Name: "trace", Fn: func(args ...object.Object) object.Object {
		indent := getIndent()
		fmt.Fprint(getStderr(), indent, "Trace")
		if len(args) > 0 {
			fmt.Fprint(getStderr(), ":")
			for _, arg := range args {
				fmt.Fprint(getStderr(), " ", arg.Inspect())
			}
		}
		fmt.Fprintln(getStderr())
		// Print a simplified stack trace
		fmt.Fprintln(getStderr(), indent+"    at <anonymous>")
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
		str, ok := args[0].(*object.String)
		if !ok {
			return newError("JSON.parse: argument must be a string")
		}
		var native interface{}
		err := json.Unmarshal([]byte(str.Value), &native)
		if err != nil {
			return newError("JSON.parse: %s", err.Error())
		}
		return nativeToObject(native)
	}})
	builtins["JSON"] = jsonObj
	
	// Array constructor
	arrayObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	// Constructor for new Array()
	arrayObj.Set("__call__", &object.Builtin{Name: "Array", Fn: func(args ...object.Object) object.Object {
		if len(args) == 1 {
			if num, ok := args[0].(*object.Number); ok {
				// Array(len)
				size := int(num.Value)
				if size < 0 {
					return newError("Invalid array length")
				}
				elements := make([]object.Object, size)
				for i := 0; i < size; i++ {
					elements[i] = UNDEFINED
				}
				return &object.Array{Elements: elements}
			}
		}
		// Array(element0, element1, ...)
		elements := make([]object.Object, len(args))
		copy(elements, args)
		return &object.Array{Elements: elements}
	}})
	
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
		
		var elements []object.Object
		
		// 1. Get elements from array-like or iterable
		switch arg := args[0].(type) {
		case *object.Array:
			elements = make([]object.Object, len(arg.Elements))
			copy(elements, arg.Elements)
		case *object.String:
			elements = make([]object.Object, len(arg.Value))
			for i, r := range arg.Value {
				elements[i] = &object.String{Value: string(r)}
			}
		case *object.ObjectMap:
			// Check for length property
			if lenProp, ok := arg.Get("length"); ok {
				if num, ok := lenProp.(*object.Number); ok {
					l := int(num.Value)
					if l < 0 { l = 0 }
					elements = make([]object.Object, l)
					for i := 0; i < l; i++ {
						key := fmt.Sprintf("%d", i)
						if val, ok := arg.Get(key); ok {
							elements[i] = val
						} else {
							elements[i] = UNDEFINED
						}
					}
				} else {
					elements = []object.Object{}
				}
			} else {
				// Not array-like (no length), return empty array
				elements = []object.Object{}
			}
		default:
			elements = []object.Object{}
		}
		
		// 2. Map function and thisArg
		if len(args) > 1 {
			if mapFn, ok := args[1].(*object.Function); ok {
				var thisArg object.Object
				if len(args) > 2 {
					thisArg = args[2]
				}
				
				mappedElements := make([]object.Object, len(elements))
				for i, elem := range elements {
					fnEnv := extendFunctionEnv(mapFn, []object.Object{elem, &object.Number{Value: float64(i)}})
					if thisArg != nil {
						fnEnv.Set("this", thisArg)
					}
					mappedElements[i] = unwrapReturnValue(Eval(mapFn.Body, fnEnv))
				}
				elements = mappedElements
			}
		}
		
		return &object.Array{Elements: elements}
	}})
	builtins["Array"] = arrayObj
	
	// Object constructor
	objectObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	objectObj.Set("keys", &object.Builtin{Name: "keys", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		if proxy, ok := args[0].(*ProxyObject); ok {
			return ProxyOwnKeys(proxy, nil)
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
		if proxy, ok := args[0].(*ProxyObject); ok {
			keysObj := ProxyOwnKeys(proxy, nil)
			keysArr, _ := keysObj.(*object.Array)
			values := []object.Object{}
			for _, key := range keysArr.Elements {
				if keyStr, ok := key.(*object.String); ok {
					values = append(values, ProxyGet(proxy, keyStr.Value, nil))
				}
			}
			return &object.Array{Elements: values}
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
		if proxy, ok := args[0].(*ProxyObject); ok {
			keysObj := ProxyOwnKeys(proxy, nil)
			keysArr, _ := keysObj.(*object.Array)
			entries := []object.Object{}
			for _, key := range keysArr.Elements {
				if keyStr, ok := key.(*object.String); ok {
					val := ProxyGet(proxy, keyStr.Value, nil)
					entry := &object.Array{Elements: []object.Object{key, val}}
					entries = append(entries, entry)
				}
			}
			return &object.Array{Elements: entries}
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
	// Function constructor - handles dynamic function creation
	// Special case: Function('return this')() is used by lodash to get global object
	functionObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	// Create a global object that mimics the global scope
	globalObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	globalObj.Set("Object", builtins["Object"])
	globalObj.Set("Array", builtins["Array"])
	globalObj.Set("String", builtins["String"])
	globalObj.Set("Number", builtins["Number"])
	globalObj.Set("Boolean", builtins["Boolean"])
	globalObj.Set("Math", builtins["Math"])
	globalObj.Set("JSON", builtins["JSON"])
	globalObj.Set("Date", builtins["Date"])
	globalObj.Set("Error", builtins["Error"])
	globalObj.Set("RegExp", builtins["RegExp"])
	globalObj.Set("Promise", builtins["Promise"])
	globalObj.Set("Map", builtins["Map"])
	globalObj.Set("Set", builtins["Set"])
	globalObj.Set("Symbol", builtins["Symbol"])
	globalObj.Set("console", builtins["console"])
	globalObj.Set("undefined", UNDEFINED)
	// Mark as global
	globalObj.Set("global", globalObj)
	globalObj.Set("globalThis", globalObj)
	globalObj.Set("self", globalObj)
	globalObj.Set("window", globalObj)
	
	functionObj.Set("__call__", &object.Builtin{Name: "Function", Fn: func(args ...object.Object) object.Object {
		// Check for the common pattern: Function('return this')
		if len(args) >= 1 {
			if str, ok := args[len(args)-1].(*object.String); ok {
				body := strings.TrimSpace(str.Value)
				if body == "return this" || body == "return this;" {
					// Return a function that returns the global object
					return &object.Builtin{Name: "anonymous", Fn: func(args ...object.Object) object.Object {
						return globalObj
					}}
				}
			}
		}
		// Return a no-op function for other cases
		return &object.Builtin{Name: "anonymous", Fn: func(args ...object.Object) object.Object {
			return UNDEFINED
		}}
	}})
	functionObj.Set("prototype", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
	builtins["Function"] = functionObj
	builtins["global"] = globalObj
	builtins["globalThis"] = globalObj
	
	// Built-in modules available globally
	builtins["fs"] = initFsModule()
	builtins["path"] = initPathModule()
	builtins["crypto"] = initCryptoModule()
	builtins["Promise"] = initPromiseConstructor()
	
	// Buffer constructor
	builtins["Buffer"] = initBufferConstructor()
	
	// HTTP module
	builtins["http"] = initHttpModule()
	
	// Stream module
	builtins["stream"] = initStreamModule()
	
	// URL module
	builtins["url"] = initURLModule()
	
	// Query string module
	builtins["querystring"] = initQueryStringModule()
	
	// Fetch function (global)
	builtins["fetch"] = initFetchFunction()
	
	// Sleep function
	builtins["sleep"] = initSleepFunction()
	
	// Initialize timer functions
	initTimerBuiltins()
	
	// RegExp constructor
	builtins["RegExp"] = initRegExpConstructor()
	
	// Date constructor and static methods
	builtins["Date"] = initDateStaticMethods()
	
	// Map constructor
	builtins["Map"] = initMapConstructor()
	
	// Set constructor
	builtins["Set"] = initSetConstructor()
	
	// WeakMap constructor
	builtins["WeakMap"] = initWeakMapConstructor()
	
	// WeakSet constructor
	builtins["WeakSet"] = initWeakSetConstructor()
	
	// Blob constructor
	builtins["Blob"] = initBlobConstructor()
	
	// File constructor
	builtins["File"] = initFileConstructor()
	
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

	// process global object
	builtins["process"] = initProcessObject()

	// Error constructor
	builtins["Error"] = initErrorConstructor()

	// Reflect API
	builtins["Reflect"] = initReflect()

	// Proxy constructor
	builtins["Proxy"] = initProxy()

	// Symbol constructor and static methods
	symbolRegistry := make(map[string]*object.Symbol)
	symbolReverseRegistry := make(map[uint64]string)

	symbolConstructor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Symbol(description) - creates a new unique symbol
	symbolConstructor.Set("__call__", &object.Builtin{Name: "Symbol", Fn: func(args ...object.Object) object.Object {
		description := ""
		if len(args) > 0 {
			description = objectToString(args[0])
		}
		return object.NewSymbol(description)
	}})

	// Symbol.for(key) - returns a symbol from the global registry
	symbolConstructor.Set("for", &object.Builtin{Name: "for", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Symbol.for requires a key argument")
		}
		key := objectToString(args[0])
		if sym, exists := symbolRegistry[key]; exists {
			return sym
		}
		sym := object.NewSymbol(key)
		symbolRegistry[key] = sym
		symbolReverseRegistry[sym.ID] = key
		return sym
	}})

	// Symbol.keyFor(symbol) - returns the key for a symbol in the global registry
	symbolConstructor.Set("keyFor", &object.Builtin{Name: "keyFor", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Symbol.keyFor requires a symbol argument")
		}
		sym, ok := args[0].(*object.Symbol)
		if !ok {
			return newError("Symbol.keyFor: argument must be a symbol")
		}
		if key, exists := symbolReverseRegistry[sym.ID]; exists {
			return &object.String{Value: key}
		}
		return UNDEFINED
	}})

	builtins["Symbol"] = symbolConstructor

	// BigInt constructor
	builtins["BigInt"] = &object.Builtin{Name: "BigInt", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("BigInt requires a value argument")
		}
		switch v := args[0].(type) {
		case *object.Number:
			// Convert to integer string
			intVal := int64(v.Value)
			return &object.BigInt{Value: fmt.Sprintf("%d", intVal)}
		case *object.String:
			// Parse as integer string
			value := v.Value
			// Remove 'n' suffix if present
			if len(value) > 0 && value[len(value)-1] == 'n' {
				value = value[:len(value)-1]
			}
			return &object.BigInt{Value: value}
		case *object.BigInt:
			return v
		default:
			return newError("Cannot convert %s to BigInt", args[0].Type())
		}
	}}
}

func stringify(obj object.Object) string {
	switch o := obj.(type) {
	case *object.String:
		// Basic escape of double quotes
		return fmt.Sprintf("\"%s\"", o.Value)
	case *object.Number:
		return fmt.Sprintf("%g", o.Value)
	case *object.Boolean:
		return fmt.Sprintf("%t", o.Value)
	case *object.Null:
		return "null"
	case *object.Array:
		parts := []string{}
		for _, elem := range o.Elements {
			s := stringify(elem)
			if s == "" { // For functions in arrays, JS uses null
				parts = append(parts, "null")
			} else {
				parts = append(parts, s)
			}
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
			s := stringify(pair.Value)
			// Skip functions/undefined in objects
			if s != "" {
				pairs = append(pairs, fmt.Sprintf("\"%s\":%s", key, s))
			}
		}
		result := "{"
		for i, pair := range pairs {
			if i > 0 {
				result += ","
			}
			result += pair
		}
		return result + "}"
	case *object.Function, *object.Builtin, *object.Undefined:
		return "" // Special marker to skip in objects
	}
	return "null"
}

func nativeToObject(native interface{}) object.Object {
	switch v := native.(type) {
	case nil:
		return NULL
	case bool:
		return nativeBoolToBooleanObject(v)
	case float64:
		return &object.Number{Value: v}
	case string:
		return &object.String{Value: v}
	case []interface{}:
		elements := make([]object.Object, len(v))
		for i, item := range v {
			elements[i] = nativeToObject(item)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[string]object.ObjectPair)
		for key, val := range v {
			objVal := nativeToObject(val)
			pairs[key] = object.ObjectPair{
				Key:   &object.String{Value: key},
				Value: objVal,
			}
		}
		return &object.ObjectMap{Pairs: pairs}
	}
	return UNDEFINED
}

// getStderr returns os.Stderr for console output
func getStderr() *os.File {
	return os.Stderr
}

// repeatChar repeats a character n times and returns the resulting string
func repeatChar(char rune, count int) string {
	result := make([]rune, count)
	for i := 0; i < count; i++ {
		result[i] = char
	}
	return string(result)
}

// truncateString truncates a string to maxLen characters, adding "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
