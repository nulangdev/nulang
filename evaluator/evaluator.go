// Package evaluator implements the evaluator for Nulang.
package evaluator

import (
	"fmt"
	"math"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

// Singleton objects
var (
	NULL      = &object.Null{}
	UNDEFINED = &object.Undefined{}
	TRUE      = &object.Boolean{Value: true}
	FALSE     = &object.Boolean{Value: false}
	BREAK     = &object.Break{}
	CONTINUE  = &object.Continue{}
)

// Eval evaluates an AST node
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ConstStatement:
		return evalConstStatement(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.ForInStatement:
		return evalForInStatement(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.DoWhileStatement:
		return evalDoWhileStatement(node, env)
	case *ast.SwitchStatement:
		return evalSwitchStatement(node, env)
	case *ast.BreakStatement:
		return BREAK
	case *ast.ContinueStatement:
		return CONTINUE
	case *ast.LabeledStatement:
		// For now, just evaluate the body - proper label handling would require tracking labels
		return Eval(node.Body, env)
	case *ast.ThrowStatement:
		return evalThrowStatement(node, env)
	case *ast.TryStatement:
		return evalTryStatement(node, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.TemplateLiteral:
		return evalTemplateLiteral(node, env)
	case *ast.RegexLiteral:
		return createRegExp(node.Pattern, node.Flags)
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.NullLiteral:
		return NULL
	case *ast.UndefinedLiteral:
		return UNDEFINED
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ConditionalExpression:
		return evalConditionalExpression(node, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.ObjectLiteral:
		return evalObjectLiteral(node, env)
	case *ast.MemberExpression:
		return evalMemberExpression(node, env)
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
	case *ast.TypeofExpression:
		return evalTypeofExpression(node, env)
	case *ast.ThisExpression:
		return evalThisExpression(env)
	case *ast.NewExpression:
		return evalNewExpressionWithClass(node, env)
	case *ast.SpreadExpression:
		return Eval(node.Value, env)
	case *ast.ClassDeclaration:
		return evalDecoratedClass(node, env)
	case *ast.SuperExpression:
		return evalSuperExpression(node, env)
	case *ast.InterfaceDeclaration:
		// Interfaces are type-only, just return undefined
		return UNDEFINED
	case *ast.TypeAliasDeclaration:
		// Type aliases are type-only, just return undefined
		return UNDEFINED
	case *ast.DeclareStatement:
		// Declarations are type-only, just return undefined
		return UNDEFINED
	case *ast.AwaitExpression:
		return evalAwaitExpression(node, env)
	}
	return UNDEFINED
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = UNDEFINED
	
	// Initialize exports for module support
	if _, ok := env.Get("exports"); !ok {
		exports := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		module := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		module.Set("exports", exports)
		env.Set("exports", exports)
		env.Set("module", module)
	}
	
	// Variable and function hoisting: first hoist all var declarations to undefined,
	// then hoist function declarations (which may overwrite some vars)
	hoistVars(program.Statements, env)
	hoistFunctions(program.Statements, env)
	
	for _, statement := range program.Statements {
		// Handle import statements
		if importStmt, ok := statement.(*ast.ImportStatement); ok {
			result = evalImportStatement(importStmt, env)
			if isError(result) {
				return result
			}
			continue
		}

		// Handle export statements
		if exportStmt, ok := statement.(*ast.ExportStatement); ok {
			result = evalExportStatement(exportStmt, env)
			if isError(result) {
				return result
			}
			continue
		}

		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

// hoistVars collects all var declarations and initializes them to undefined
// before executing any code in the block, mimicking JavaScript's var hoisting behavior.
// This is essential for libraries like lodash that use circular bootstrap patterns
// where variables are referenced before their lexical definition.
// IMPORTANT: var declarations are hoisted to function scope, not block scope,
// so we must recursively scan ALL nested blocks (if, while, for, etc.)
func hoistVars(statements []ast.Statement, env *object.Environment) {
	for _, stmt := range statements {
		hoistVarsFromStatement(stmt, env)
	}
}

// hoistVarsFromStatement recursively extracts var declarations from a statement
func hoistVarsFromStatement(stmt ast.Statement, env *object.Environment) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		// Only hoist if the variable is not already defined
		// (to avoid overwriting function declarations that were hoisted first in the same scope)
		if _, ok := env.Get(s.Name.Value); !ok {
			env.Set(s.Name.Value, UNDEFINED)
		}
		// Also hoist any additional declarations in comma-separated vars
		for _, decl := range s.Declarations {
			if _, ok := env.Get(decl.Name.Value); !ok {
				env.Set(decl.Name.Value, UNDEFINED)
			}
		}
	case *ast.ForStatement:
		// Hoist var declarations inside for loop init
		if varStmt, ok := s.Init.(*ast.VarStatement); ok {
			hoistVarsFromStatement(varStmt, env)
		}
		// Also scan the body
		if s.Body != nil {
			hoistVars(s.Body.Statements, env)
		}
	case *ast.ForInStatement:
		// Hoist var declarations inside for-in loop
		if s.IsVar {
			if _, ok := env.Get(s.Key.Value); !ok {
				env.Set(s.Key.Value, UNDEFINED)
			}
		}
		// Also scan the body
		if s.Body != nil {
			hoistVars(s.Body.Statements, env)
		}
	case *ast.WhileStatement:
		// Scan while body
		if s.Body != nil {
			hoistVars(s.Body.Statements, env)
		}
	case *ast.DoWhileStatement:
		// Scan do-while body
		if s.Body != nil {
			hoistVars(s.Body.Statements, env)
		}
	case *ast.SwitchStatement:
		// Scan all case bodies
		for _, c := range s.Cases {
			if c.Body != nil {
				hoistVars(c.Body.Statements, env)
			}
		}
		// Scan default body
		if s.Default != nil {
			hoistVars(s.Default.Statements, env)
		}
	case *ast.TryStatement:
		// Scan try block
		if s.Block != nil {
			hoistVars(s.Block.Statements, env)
		}
		// Scan catch block
		if s.CatchBlock != nil {
			hoistVars(s.CatchBlock.Statements, env)
		}
		// Scan finally block
		if s.FinallyBlock != nil {
			hoistVars(s.FinallyBlock.Statements, env)
		}
	case *ast.BlockStatement:
		// Scan block contents
		hoistVars(s.Statements, env)
	case *ast.LabeledStatement:
		// Scan the labeled statement's body
		if s.Body != nil {
			hoistVarsFromStatement(s.Body, env)
		}
	case *ast.ExpressionStatement:
		// Check for if expressions inside expression statements
		if ie, ok := s.Expression.(*ast.IfExpression); ok {
			if ie.Consequence != nil {
				hoistVars(ie.Consequence.Statements, env)
			}
			if ie.Alternative != nil {
				hoistVars(ie.Alternative.Statements, env)
			}
		}
	}
}

// hoistFunctions collects function declarations and adds them to the environment
// before executing any code in the block, mimicking JavaScript's hoisting behavior
func hoistFunctions(statements []ast.Statement, env *object.Environment) {
	for _, stmt := range statements {
		switch s := stmt.(type) {
		case *ast.VarStatement:
			// Check if this is a function declaration: function name() {}
			// which is parsed as VarStatement with FunctionLiteral value
			if fl, ok := s.Value.(*ast.FunctionLiteral); ok && fl.Name != nil {
				fn := &object.Function{
					Parameters: fl.Parameters,
					Body:       fl.Body,
					Env:        env,
					IsAsync:    fl.IsAsync,
					Name:       fl.Name.Value,
				}
				env.Set(fl.Name.Value, fn)
			}
		}
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = UNDEFINED
	
	// Variable and function hoisting within block
	hoistVars(block.Statements, env)
	hoistFunctions(block.Statements, env)
	
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ ||
				rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}
	return result
}

func evalLetStatement(ls *ast.LetStatement, env *object.Environment) object.Object {
	var val object.Object = UNDEFINED
	if ls.Value != nil {
		val = Eval(ls.Value, env)
		if isError(val) {
			return val
		}
	}
	env.Set(ls.Name.Value, val)
	return val
}

func evalConstStatement(cs *ast.ConstStatement, env *object.Environment) object.Object {
	val := Eval(cs.Value, env)
	if isError(val) {
		return val
	}
	env.SetConst(cs.Name.Value, val)
	return val
}

func evalVarStatement(vs *ast.VarStatement, env *object.Environment) object.Object {
	var val object.Object = UNDEFINED
	if vs.Value != nil {
		val = Eval(vs.Value, env)
		if isError(val) {
			return val
		}
	}
	env.Set(vs.Name.Value, val)
	
	// Handle additional declarations (for comma-separated var declarations)
	for _, decl := range vs.Declarations {
		var declVal object.Object = UNDEFINED
		if decl.Value != nil {
			declVal = Eval(decl.Value, env)
			if isError(declVal) {
				return declVal
			}
		}
		env.Set(decl.Name.Value, declVal)
	}
	
	return val
}

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	if rs.ReturnValue == nil {
		return &object.ReturnValue{Value: UNDEFINED}
	}
	val := Eval(rs.ReturnValue, env)
	if isError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	if fs.Init != nil {
		Eval(fs.Init, loopEnv)
	}
	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, loopEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}
		result := Eval(fs.Body, loopEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
		}
		if fs.Update != nil {
			Eval(fs.Update, loopEnv)
		}
	}
	return UNDEFINED
}

func evalForInStatement(fs *ast.ForInStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	
	obj := Eval(fs.Object, env)
	if isError(obj) {
		return obj
	}

	var keys []string

	switch o := obj.(type) {
	case *object.ObjectMap:
		for key := range o.Pairs {
			keys = append(keys, key)
		}
	case *object.Array:
		for i := range o.Elements {
			keys = append(keys, fmt.Sprintf("%d", i))
		}
	case *object.String:
		for i := range o.Value {
			keys = append(keys, fmt.Sprintf("%d", i))
		}
	default:
		return UNDEFINED
	}

	for _, key := range keys {
		loopEnv.Set(fs.Key.Value, &object.String{Value: key})
		
		result := Eval(fs.Body, loopEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			// CONTINUE is handled by continuing the loop
		}
	}
	
	return UNDEFINED
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		result := Eval(ws.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
		}
	}
	return UNDEFINED
}

func evalDoWhileStatement(dw *ast.DoWhileStatement, env *object.Environment) object.Object {
	for {
		// Execute body first (at least once)
		result := Eval(dw.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
		}
		// Then check condition
		condition := Eval(dw.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
	}
	return UNDEFINED
}

func evalSwitchStatement(ss *ast.SwitchStatement, env *object.Environment) object.Object {
	switchVal := Eval(ss.Value, env)
	if isError(switchVal) {
		return switchVal
	}

	matched := false
	shouldFallthrough := false

	for _, caseClause := range ss.Cases {
		if !matched && !shouldFallthrough {
			testVal := Eval(caseClause.Test, env)
			if isError(testVal) {
				return testVal
			}
			// Strict equality check
			if objectsEqual(switchVal, testVal) {
				matched = true
			}
		}

		if matched || shouldFallthrough {
			result := Eval(caseClause.Body, env)
			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					return UNDEFINED
				}
			}
			// If no break, fall through to next case
			shouldFallthrough = true
		}
	}

	// Execute default if no match or fallthrough
	if ss.Default != nil && (matched || !matched && !shouldFallthrough || shouldFallthrough) {
		if !matched && !shouldFallthrough {
			// No case matched, execute default
			result := Eval(ss.Default, env)
			if result != nil && result.Type() != object.BREAK_OBJ {
				return result
			}
		} else if shouldFallthrough {
			// Falling through to default
			result := Eval(ss.Default, env)
			if result != nil && result.Type() != object.BREAK_OBJ {
				return result
			}
		}
	}

	return UNDEFINED
}

// objectsEqual checks if two objects are equal (strict equality)
func objectsEqual(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch av := a.(type) {
	case *object.Number:
		if bv, ok := b.(*object.Number); ok {
			return av.Value == bv.Value
		}
	case *object.String:
		if bv, ok := b.(*object.String); ok {
			return av.Value == bv.Value
		}
	case *object.Boolean:
		if bv, ok := b.(*object.Boolean); ok {
			return av.Value == bv.Value
		}
	case *object.Null:
		_, ok := b.(*object.Null)
		return ok
	case *object.Undefined:
		_, ok := b.(*object.Undefined)
		return ok
	}
	// Reference equality for objects
	return a == b
}

func evalThrowStatement(ts *ast.ThrowStatement, env *object.Environment) object.Object {
	val := Eval(ts.Value, env)
	if isError(val) {
		return val
	}
	
	// If it's an Error instance (ObjectMap with message property), extract the message
	message := getErrorMessage(val)
	
	return &object.Error{Message: message}
}

func evalTryStatement(ts *ast.TryStatement, env *object.Environment) object.Object {
	result := Eval(ts.Block, env)
	if isError(result) && ts.CatchBlock != nil {
		catchEnv := object.NewEnclosedEnvironment(env)
		if ts.CatchParam != nil {
			// Create an Error object to pass to the catch block
			errorObj := createErrorObject(result.(*object.Error).Message)
			catchEnv.Set(ts.CatchParam.Value, errorObj)
		}
		result = Eval(ts.CatchBlock, catchEnv)
	}
	if ts.FinallyBlock != nil {
		finallyResult := Eval(ts.FinallyBlock, env)
		if isError(finallyResult) {
			return finallyResult
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	initBuiltins()
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("%s", "identifier not found: " + node.Value)
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Null:
		return false
	case *object.Undefined:
		return false
	case *object.Boolean:
		return obj.Value
	case *object.Number:
		return obj.Value != 0 && !math.IsNaN(obj.Value)
	case *object.String:
		return obj.Value != ""
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func objectToString(obj object.Object) string {
	switch obj := obj.(type) {
	case *object.String:
		return obj.Value
	case *object.Number:
		return fmt.Sprintf("%g", obj.Value)
	case *object.Boolean:
		if obj.Value {
			return "true"
		}
		return "false"
	case *object.Null:
		return "null"
	case *object.Undefined:
		return "undefined"
	default:
		return obj.Inspect()
	}
}

// evalTemplateLiteral evaluates a template literal with ${} interpolation
func evalTemplateLiteral(tl *ast.TemplateLiteral, env *object.Environment) object.Object {
	var result string
	
	for i, part := range tl.Parts {
		result += part
		if i < len(tl.Expressions) {
			val := Eval(tl.Expressions[i], env)
			if isError(val) {
				return val
			}
			result += objectToString(val)
		}
	}
	
	return &object.String{Value: result}
}

// evalAwaitExpression evaluates an await expression
func evalAwaitExpression(ae *ast.AwaitExpression, env *object.Environment) object.Object {
	val := Eval(ae.Value, env)
	if isError(val) {
		return val
	}
	
	// If it's a promise, unwrap it
	if promise, ok := val.(*object.Promise); ok {
		switch promise.State {
		case "fulfilled":
			return promise.Value
		case "rejected":
			if promise.Reason != nil {
				if err, ok := promise.Reason.(*object.Error); ok {
					return err
				}
				return newError("%s", promise.Reason.Inspect())
			}
			return newError("Promise rejected")
		default:
			// Promise is pending - in a real implementation we'd wait
			// For now, just return undefined
			return UNDEFINED
		}
	}
	
	// If it's not a promise, just return the value directly
	return val
}
