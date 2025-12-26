package evaluator

import (
	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/lexer"
	"github.com/nulang/nulang/object"
	"github.com/nulang/nulang/parser"
)

// initVMModule initializes the vm module
func initVMModule() *object.ObjectMap {
	vm := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// vm.runInThisContext(code, options?)
	vm.Set("runInThisContext", &object.Builtin{Name: "runInThisContext", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("runInThisContext requires code argument")
		}

		code := objectToString(args[0])
		
		// Parse the code
		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			return newError("Parse error: %s", p.Errors()[0])
		}

		// Execute in current global context
		env := object.NewEnvironment()
		
		// Copy globals from builtins
		initBuiltins()
		for key, val := range builtins {
			env.Set(key, val)
		}

		result := Eval(program, env)
		return unwrapReturnValue(result)
	}})

	// vm.runInNewContext(code, sandbox?, options?)
	vm.Set("runInNewContext", &object.Builtin{Name: "runInNewContext", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("runInNewContext requires code argument")
		}

		code := objectToString(args[0])
		
		// Create new sandbox environment
		env := object.NewEnvironment()
		
		// Add sandbox properties if provided
		if len(args) > 1 {
			if sandbox, ok := args[1].(*object.ObjectMap); ok {
				for key, pair := range sandbox.Pairs {
					env.Set(key, pair.Value)
				}
			}
		}

		// Add minimal globals
		env.Set("console", builtins["console"])
		env.Set("Math", builtins["Math"])
		env.Set("JSON", builtins["JSON"])

		// Parse the code
		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			return newError("Parse error: %s", p.Errors()[0])
		}

		result := Eval(program, env)
		return unwrapReturnValue(result)
	}})

	// vm.runInContext(code, contextifiedObject, options?)
	vm.Set("runInContext", &object.Builtin{Name: "runInContext", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("runInContext requires code and context arguments")
		}

		code := objectToString(args[0])
		
		// Use the provided context as environment
		env := object.NewEnvironment()
		
		if context, ok := args[1].(*object.ObjectMap); ok {
			for key, pair := range context.Pairs {
				env.Set(key, pair.Value)
			}
		}

		// Parse the code
		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			return newError("Parse error: %s", p.Errors()[0])
		}

		result := Eval(program, env)
		
		// Update context with new values - skip this for now as we can't iterate env
		// In a full implementation, we'd need to track set variables differently

		return unwrapReturnValue(result)
	}})

	// vm.createContext(contextObject?)
	vm.Set("createContext", &object.Builtin{Name: "createContext", Fn: func(args ...object.Object) object.Object {
		context := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		
		// Copy properties from provided object
		if len(args) > 0 {
			if obj, ok := args[0].(*object.ObjectMap); ok {
				for key, pair := range obj.Pairs {
					context.Set(key, pair.Value)
				}
			}
		}

		// Add minimal globals
		context.Set("console", builtins["console"])
		context.Set("Math", builtins["Math"])
		context.Set("JSON", builtins["JSON"])
		context.Set("Array", builtins["Array"])
		context.Set("Object", builtins["Object"])
		context.Set("String", builtins["String"])
		context.Set("Number", builtins["Number"])
		context.Set("Boolean", builtins["Boolean"])

		return context
	}})

	// vm.isContext(object)
	vm.Set("isContext", &object.Builtin{Name: "isContext", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		
		// In our implementation, any ObjectMap can be a context
		_, ok := args[0].(*object.ObjectMap)
		return nativeBoolToBooleanObject(ok)
	}})

	// vm.compileFunction(code, params?, options?)
	vm.Set("compileFunction", &object.Builtin{Name: "compileFunction", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("compileFunction requires code argument")
		}

		code := objectToString(args[0])
		var params []string

		// Parse parameters if provided
		if len(args) > 1 {
			if paramsArr, ok := args[1].(*object.Array); ok {
				for _, param := range paramsArr.Elements {
					params = append(params, objectToString(param))
				}
			}
		}

		// Build function code
		functionCode := "function("
		for i, param := range params {
			if i > 0 {
				functionCode += ", "
			}
			functionCode += param
		}
		functionCode += ") { " + code + " }"

		// Parse the function
		l := lexer.New(functionCode)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			return newError("Parse error: %s", p.Errors()[0])
		}

		// Extract the function
		if len(program.Statements) > 0 {
			if exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
				if funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
					env := object.NewEnvironment()
					return &object.Function{
						Parameters: funcLit.Parameters,
						Body:       funcLit.Body,
						Env:        env,
					}
				}
			}
		}

		return newError("Failed to compile function")
	}})

	// VM Script class
	scriptConstructor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	scriptConstructor.Set("__call__", &object.Builtin{Name: "Script", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Script constructor requires code argument")
		}

		code := objectToString(args[0])
		
		// Parse the code
		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			return newError("Parse error: %s", p.Errors()[0])
		}

		// Create script object
		script := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		script.Set("__program", &object.String{Value: code}) // Store as string for now

		// script.runInThisContext(options?)
		script.Set("runInThisContext", &object.Builtin{Name: "runInThisContext", Fn: func(args ...object.Object) object.Object {
			env := object.NewEnvironment()
			
			for key, val := range builtins {
				env.Set(key, val)
			}

			result := Eval(program, env)
			return unwrapReturnValue(result)
		}})

		// script.runInNewContext(sandbox?, options?)
		script.Set("runInNewContext", &object.Builtin{Name: "runInNewContext", Fn: func(args ...object.Object) object.Object {
			env := object.NewEnvironment()
			
			if len(args) > 0 {
				if sandbox, ok := args[0].(*object.ObjectMap); ok {
					for key, pair := range sandbox.Pairs {
						env.Set(key, pair.Value)
					}
				}
			}

			env.Set("console", builtins["console"])
			env.Set("Math", builtins["Math"])
			env.Set("JSON", builtins["JSON"])

			result := Eval(program, env)
			return unwrapReturnValue(result)
		}})

		// script.runInContext(contextifiedObject, options?)
		script.Set("runInContext", &object.Builtin{Name: "runInContext", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("runInContext requires a context argument")
			}

			env := object.NewEnvironment()
			
			if context, ok := args[0].(*object.ObjectMap); ok {
				for key, pair := range context.Pairs {
					env.Set(key, pair.Value)
				}
			}

			result := Eval(program, env)
			return unwrapReturnValue(result)
		}})

		return script
	}})

	vm.Set("Script", scriptConstructor)

	return vm
}
