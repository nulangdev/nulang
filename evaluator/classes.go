package evaluator

import (
	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/object"
)

// Class represents a class definition
type Class struct {
	Name       string
	SuperClass *Class
	Properties map[string]object.Object
	Methods    map[string]*object.Function
	Getters    map[string]*object.Function
	Setters    map[string]*object.Function
	Static        map[string]object.Object
	Env           *object.Environment
	NativeMethods map[string]NativeMethod
}

// NativeMethod is a function that can be attached to a class and access 'this'
type NativeMethod func(this object.Object, args ...object.Object) object.Object

func (c *Class) Type() object.ObjectType { return object.CLASS_OBJ }
func (c *Class) Inspect() string         { return "[class " + c.Name + "]" }

// evalClassDeclaration evaluates a class declaration
func evalClassDeclaration(cd *ast.ClassDeclaration, env *object.Environment) object.Object {
	class := &Class{
		Properties: make(map[string]object.Object),
		Methods:    make(map[string]*object.Function),
		Getters:    make(map[string]*object.Function),
		Setters:    make(map[string]*object.Function),
		Static:        make(map[string]object.Object),
		NativeMethods: make(map[string]NativeMethod),
		Env:           env,
	}

	if cd.Name != nil {
		class.Name = cd.Name.Value
	}

	// Handle extends
	if cd.SuperClass != nil {
		superObj := Eval(cd.SuperClass, env)
		if isError(superObj) {
			return superObj
		}
		
		if superClass, ok := superObj.(*Class); ok {
			class.SuperClass = superClass
		} else {
			return newError("Superclass is not a class. Got %s", superObj.Type())
		}
	}

	// Process class body
	if cd.Body != nil {
		for _, member := range cd.Body.Members {
			if member.Name == nil {
				continue
			}

			memberName := member.Name.Value

			// Handle getter
			if member.IsGetter {
				if fn, ok := member.Value.(*ast.FunctionLiteral); ok {
					class.Getters[memberName] = &object.Function{
						Parameters: fn.Parameters,
						Body:       fn.Body,
						Env:        env,
					}
				}
				continue
			}

			// Handle setter
			if member.IsSetter {
				if fn, ok := member.Value.(*ast.FunctionLiteral); ok {
					class.Setters[memberName] = &object.Function{
						Parameters: fn.Parameters,
						Body:       fn.Body,
						Env:        env,
					}
				}
				continue
			}

			// Handle method or property
			if fn, ok := member.Value.(*ast.FunctionLiteral); ok {
				method := &object.Function{
					Name:       memberName,
					Parameters: fn.Parameters,
					Body:       fn.Body,
					Env:        env,
				}

				// Apply method decorators if present
				var finalMethod object.Object = method
				if len(member.Decorators) > 0 {
					// Create a temporary class object as target for decorator
					targetClass := &Class{Name: class.Name}
					decorated := applyMethodDecorators(method, memberName, member.Decorators, targetClass, env)
					if isError(decorated) {
						return decorated
					}
					finalMethod = decorated
				}

				if member.IsStatic {
					if fn, ok := finalMethod.(*object.Function); ok {
						class.Static[memberName] = fn
					} else {
						class.Static[memberName] = finalMethod
					}
				} else {
					if fn, ok := finalMethod.(*object.Function); ok {
						class.Methods[memberName] = fn
					} else {
						// If decorator returned something other than a function,
						// store it as a property instead
						class.Properties[memberName] = finalMethod
					}
				}
			} else if member.Value != nil {
				// Property
				value := Eval(member.Value, env)
				if isError(value) {
					return value
				}

				// Apply property decorators if present
				if len(member.Decorators) > 0 {
					// Create a temporary class object as target for decorator
					targetClass := &Class{Name: class.Name}
					decorated := applyPropertyDecorators(value, memberName, member.Decorators, targetClass, env)
					if isError(decorated) {
						return decorated
					}
					value = decorated
				}

				if member.IsStatic {
					class.Static[memberName] = value
				} else {
					class.Properties[memberName] = value
				}
			}
		}
	}

	// Register class in environment if named
	if class.Name != "" {
		env.Set(class.Name, class)
	}

	return class
}

// createClassInstance creates a new instance of a class
func createClassInstance(class *Class, args []object.Object, _ *object.Environment) object.Object {
	instance := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Create instance environment
	instanceEnv := object.NewEnclosedEnvironment(class.Env)

	// Copy properties from superclass chain
	if class.SuperClass != nil {
		copyClassProperties(instance, class.SuperClass)
	}

	// Copy default property values
	for name, value := range class.Properties {
		instance.Set(name, value)
	}

	// Set 'super' if there's a superclass
	var superObj *object.ObjectMap
	if class.SuperClass != nil {
		superObj = createSuperObject(class.SuperClass, instance, instanceEnv)
		instanceEnv.Set("super", superObj)
	}

	// Add methods as properties
	for name, method := range class.Methods {
		// Create bound method environment that preserves the original closure (method.Env)
		// but also provides 'this' and 'super'.
		methodEnv := object.NewEnclosedEnvironment(method.Env)
		methodEnv.Set("this", instance)
		if superObj != nil {
			methodEnv.Set("super", superObj)
		}

		// Create bound method
		boundMethod := &object.Function{
			Name:       method.Name,
			Parameters: method.Parameters,
			Body:       method.Body,
			Env:        methodEnv,
		}
		instance.Set(name, boundMethod)
	}

	// Add native methods
	for name, method := range class.NativeMethods {
		boundMethod := bindNativeMethod(method, instance)
		instance.Set(name, boundMethod)
	}

	// Copy superclass methods
	if class.SuperClass != nil {
		// Note: copyClassMethods might need Env to bind methods? 
		// Actually copyClassMethods creates bound methods too.
		// Let's check copyClassMethods implementation later, for now pass instanceEnv 
		// but keep in mind superclass methods are usually not closures from decorators.
		copyClassMethods(instance, class.SuperClass, instanceEnv)
	}

	// Set up getters and setters (store in prototype)
	instance.Prototype = createPrototype(class, instanceEnv)

	// Set 'this' reference in instanceEnv (still needed for constructor)
	instanceEnv.Set("this", instance)

	// Call constructor if exists
	if constructor, ok := class.Methods["constructor"]; ok {
		constructorEnv := object.NewEnclosedEnvironment(instanceEnv)
		constructorEnv.Set("this", instance)

		// Bind parameters
		for i, param := range constructor.Parameters {
			if i < len(args) {
				constructorEnv.Set(param.Value, args[i])
			} else {
				constructorEnv.Set(param.Value, UNDEFINED)
			}
		}

		// Execute constructor
		result := Eval(constructor.Body, constructorEnv)
		if isError(result) {
			return result
		}
	} else if nativeConstructor, ok := class.NativeMethods["constructor"]; ok {
		result := nativeConstructor(instance, args...)
		if isError(result) {
			return result
		}
	}

	return instance
}

// copyClassProperties copies properties from superclass to instance
func copyClassProperties(instance *object.ObjectMap, class *Class) {
	if class.SuperClass != nil {
		copyClassProperties(instance, class.SuperClass)
	}
	for name, value := range class.Properties {
		instance.Set(name, value)
	}
}

// copyClassMethods copies methods from superclass to instance
func copyClassMethods(instance *object.ObjectMap, class *Class, env *object.Environment) {
	if class.SuperClass != nil {
		copyClassMethods(instance, class.SuperClass, env)
	}
	for name, method := range class.Methods {
		if name == "constructor" {
			continue
		}
		// Don't override if already exists
		if _, ok := instance.Get(name); !ok {
			boundMethod := &object.Function{
				Name:       method.Name,
				Parameters: method.Parameters,
				Body:       method.Body,
				Env:        env,
			}
			instance.Set(name, boundMethod)
		}
	}

	// Copy native methods from superclass
	for name, method := range class.NativeMethods {
		if name == "constructor" {
			continue
		}
		if _, ok := instance.Get(name); !ok {
			boundMethod := bindNativeMethod(method, instance)
			instance.Set(name, boundMethod)
		}
	}
}

// bindNativeMethod creates a Builtin that calls the NativeMethod with 'this' context
func bindNativeMethod(method NativeMethod, instance object.Object) *object.Builtin {
	return &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return method(instance, args...)
		},
	}
}

// createPrototype creates a prototype object with getters/setters
func createPrototype(class *Class, env *object.Environment) *object.ObjectMap {
	prototype := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Add getters as special properties
	for name, getter := range class.Getters {
		prototype.Set("__get_"+name, &object.Function{
			Parameters: getter.Parameters,
			Body:       getter.Body,
			Env:        env,
		})
	}

	// Add setters
	for name, setter := range class.Setters {
		prototype.Set("__set_"+name, &object.Function{
			Parameters: setter.Parameters,
			Body:       setter.Body,
			Env:        env,
		})
	}

	// Chain prototypes
	if class.SuperClass != nil {
		prototype.Prototype = createPrototype(class.SuperClass, env)
	}

	return prototype
}

// createSuperObject creates a 'super' object for calling parent methods
// createSuperObject creates a 'super' object for calling parent methods
func createSuperObject(superClass *Class, instance *object.ObjectMap, env *object.Environment) *object.ObjectMap {
	superObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Add parent methods bound to the instance
	for name, method := range superClass.Methods {
		boundMethod := &object.Builtin{
			Name: name,
			Fn: func(m *object.Function) object.BuiltinFunction {
				return func(args ...object.Object) object.Object {
					methodEnv := object.NewEnclosedEnvironment(env)
					methodEnv.Set("this", instance)

					for i, param := range m.Parameters {
						if i < len(args) {
							methodEnv.Set(param.Value, args[i])
						} else {
							methodEnv.Set(param.Value, UNDEFINED)
						}
					}

					result := Eval(m.Body, methodEnv)
					if rv, ok := result.(*object.ReturnValue); ok {
						return rv.Value
					}
					return result
				}
			}(method),
		}
		superObj.Set(name, boundMethod)
		
		if name == "constructor" {
			superObj.Set("__call__", boundMethod)
		}
	}
	
	// Add native methods
	for name, method := range superClass.NativeMethods {
		boundMethod := bindNativeMethod(method, instance)
		superObj.Set(name, boundMethod)
		
		if name == "constructor" {
			superObj.Set("__call__", boundMethod)
		}
	}

	return superObj
}

// evalNewExpression handles 'new ClassName(args)'
func evalNewExpressionWithClass(ne *ast.NewExpression, env *object.Environment) object.Object {
	// Evaluate the class expression
	classObj := Eval(ne.Class, env)
	if isError(classObj) {
		return classObj
	}

	// Evaluate arguments
	args := evalExpressions(ne.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	// Check if it's a class
	if class, ok := classObj.(*Class); ok {
		return createClassInstance(class, args, env)
	}

	// Check if it's a builtin constructor (e.g., RegExp)
	if builtin, ok := classObj.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}

	// Check for ObjectMap with __call__ (e.g., Date, Map, Set, Proxy)
	if objMap, ok := classObj.(*object.ObjectMap); ok {
		// Special handling for Proxy
		if isProxyConstructor(objMap) {
			return createProxy(args)
		}
		if callFn, ok := objMap.Get("__call__"); ok {
			if builtin, ok := callFn.(*object.Builtin); ok {
				return builtin.Fn(args...)
			}
		}
	}

	// Check for function constructor pattern
	if fn, ok := classObj.(*object.Function); ok {
		// Create instance
		instance := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		instanceEnv := object.NewEnclosedEnvironment(fn.Env)
		instanceEnv.Set("this", instance)

		// Bind parameters
		for i, param := range fn.Parameters {
			if i < len(args) {
				instanceEnv.Set(param.Value, args[i])
			} else {
				instanceEnv.Set(param.Value, UNDEFINED)
			}
		}

		// Execute function
		result := Eval(fn.Body, instanceEnv)
		if isError(result) {
			return result
		}

		// If function returns an object, use that instead
		if obj, ok := result.(*object.ObjectMap); ok {
			return obj
		}
		if rv, ok := result.(*object.ReturnValue); ok {
			if obj, ok := rv.Value.(*object.ObjectMap); ok {
				return obj
			}
		}

		return instance
	}

	return newError("'%s' is not a constructor", classObj.Type())
}

// evalSuperExpression handles 'super' keyword
func evalSuperExpression(_ *ast.SuperExpression, env *object.Environment) object.Object {
	superObj, ok := env.Get("super")
	if !ok {
		return newError("'super' keyword is only valid inside a class method")
	}
	return superObj
}
