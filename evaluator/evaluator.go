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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = UNDEFINED
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
