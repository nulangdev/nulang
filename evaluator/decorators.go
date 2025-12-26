// Package evaluator implements Decorators evaluation for Nulang.
package evaluator

import (
	"reflect"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

// evalDecorator evaluates a single decorator
func evalDecorator(decorator *ast.Decorator, env *object.Environment) object.Object {
	// Get the decorator function
	decoratorFn, ok := env.Get(decorator.Name.Value)
	if !ok {
		return newError("decorator '%s' is not defined", decorator.Name.Value)
	}

	// If decorator has arguments, call the decorator factory first
	if len(decorator.Arguments) > 0 {
		args := []object.Object{}
		for _, arg := range decorator.Arguments {
			evaluated := Eval(arg, env)
			if isError(evaluated) {
				return evaluated
			}
			args = append(args, evaluated)
		}

		// Call the decorator factory
		switch fn := decoratorFn.(type) {
		case *object.Function:
			extendedEnv := extendFunctionEnv(fn, args)
			evaluated := Eval(fn.Body, extendedEnv)
			decoratorFn = unwrapReturnValue(evaluated)
		case *object.Builtin:
			decoratorFn = fn.Fn(args...)
		default:
			return newError("decorator '%s' is not a function", decorator.Name.Value)
		}
	}

	return decoratorFn
}

// applyClassDecorators applies decorators to a class
func applyClassDecorators(class *Class, decorators []*ast.Decorator, env *object.Environment) object.Object {
	result := object.Object(class)

	// Apply decorators in reverse order (bottom-up)
	for i := len(decorators) - 1; i >= 0; i-- {
		decorator := decorators[i]
		decoratorFn := evalDecorator(decorator, env)
		if isError(decoratorFn) {
			return decoratorFn
		}

		// Call the decorator with the class
		switch fn := decoratorFn.(type) {
		case *object.Function:
			args := []object.Object{result}
			extendedEnv := extendFunctionEnv(fn, args)
			evaluated := Eval(fn.Body, extendedEnv)
			newResult := unwrapReturnValue(evaluated)
			// Decorators can return a new class or modify the existing one
			if newResult != nil && newResult.Type() != object.UNDEFINED_OBJ && newResult.Type() != object.NULL_OBJ {
				result = newResult
			}
		case *object.Builtin:
			newResult := fn.Fn(result)
			if newResult != nil && newResult.Type() != object.UNDEFINED_OBJ && newResult.Type() != object.NULL_OBJ {
				result = newResult
			}
		default:
			return newError("decorator must be a function, got %s", reflect.TypeOf(decoratorFn))
		}
	}

	return result
}

// applyMethodDecorators applies decorators to a class method
func applyMethodDecorators(method *object.Function, methodName string, decorators []*ast.Decorator, target object.Object, env *object.Environment) object.Object {
	result := object.Object(method)

	// Create descriptor object
	descriptor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	descriptor.Set("value", method)
	descriptor.Set("writable", TRUE)
	descriptor.Set("enumerable", FALSE)
	descriptor.Set("configurable", TRUE)

	// Apply decorators in reverse order (bottom-up)
	for i := len(decorators) - 1; i >= 0; i-- {
		decorator := decorators[i]
		decoratorFn := evalDecorator(decorator, env)
		if isError(decoratorFn) {
			return decoratorFn
		}

		// Call the decorator with (target, methodName, descriptor)
		propKey := &object.String{Value: methodName}

		switch fn := decoratorFn.(type) {
		case *object.Function:
			args := []object.Object{target, propKey, descriptor}
			extendedEnv := extendFunctionEnv(fn, args)
			evaluated := Eval(fn.Body, extendedEnv)
			newDescriptor := unwrapReturnValue(evaluated)
			// If decorator returns a new descriptor, use it
			if descMap, ok := newDescriptor.(*object.ObjectMap); ok {
				descriptor = descMap
			}
		case *object.Builtin:
			newDescriptor := fn.Fn(target, propKey, descriptor)
			if descMap, ok := newDescriptor.(*object.ObjectMap); ok {
				descriptor = descMap
			}
		}
	}

	// Get the final method from the descriptor
	if val, found := descriptor.Get("value"); found {
		result = val
	}

	return result
}

// applyPropertyDecorators applies decorators to a class property
func applyPropertyDecorators(propertyValue object.Object, propertyName string, decorators []*ast.Decorator, target object.Object, env *object.Environment) object.Object {
	result := propertyValue

	// Create descriptor object
	descriptor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	if result != nil {
		descriptor.Set("value", result)
	}
	descriptor.Set("writable", TRUE)
	descriptor.Set("enumerable", TRUE)
	descriptor.Set("configurable", TRUE)

	// Apply decorators in reverse order
	for i := len(decorators) - 1; i >= 0; i-- {
		decorator := decorators[i]
		decoratorFn := evalDecorator(decorator, env)
		if isError(decoratorFn) {
			return decoratorFn
		}

		propKey := &object.String{Value: propertyName}

		switch fn := decoratorFn.(type) {
		case *object.Function:
			args := []object.Object{target, propKey, descriptor}
			extendedEnv := extendFunctionEnv(fn, args)
			evaluated := Eval(fn.Body, extendedEnv)
			newDescriptor := unwrapReturnValue(evaluated)
			if descMap, ok := newDescriptor.(*object.ObjectMap); ok {
				descriptor = descMap
			}
		case *object.Builtin:
			newDescriptor := fn.Fn(target, propKey, descriptor)
			if descMap, ok := newDescriptor.(*object.ObjectMap); ok {
				descriptor = descMap
			}
		}
	}

	// Get the final value from the descriptor
	if val, found := descriptor.Get("value"); found {
		result = val
	}

	return result
}

// evalDecoratedClass evaluates a class with decorators
func evalDecoratedClass(cd *ast.ClassDeclaration, env *object.Environment) object.Object {
	// First evaluate the class normally
	class := evalClassDeclaration(cd, env)
	if isError(class) {
		return class
	}

	// If there are no decorators, return the class as is
	if len(cd.Decorators) == 0 {
		return class
	}

	// Apply class decorators
	classObj, ok := class.(*Class)
	if !ok {
		return class
	}

	return applyClassDecorators(classObj, cd.Decorators, env)
}
