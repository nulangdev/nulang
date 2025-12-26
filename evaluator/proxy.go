// Package evaluator implements the Proxy object for Nulang.
package evaluator

import (
	"sort"

	"github.com/nulang/nulang/object"
)

// ProxyHandler represents the handler for a Proxy
type ProxyHandler struct {
	Get                      object.Object // function(target, property, receiver)
	Set                      object.Object // function(target, property, value, receiver)
	Has                      object.Object // function(target, property)
	DeleteProperty           object.Object // function(target, property)
	Apply                    object.Object // function(target, thisArg, argumentsList)
	Construct                object.Object // function(target, argumentsList, newTarget)
	GetPrototypeOf           object.Object // function(target)
	SetPrototypeOf           object.Object // function(target, prototype)
	IsExtensible             object.Object // function(target)
	PreventExtensions        object.Object // function(target)
	GetOwnPropertyDescriptor object.Object // function(target, property)
	DefineProperty           object.Object // function(target, property, descriptor)
	OwnKeys                  object.Object // function(target)
}

// ProxyObject represents a Proxy object
type ProxyObject struct {
	Target  object.Object
	Handler *ProxyHandler
	Revoked bool
}

func (p *ProxyObject) Type() object.ObjectType { return "PROXY" }
func (p *ProxyObject) Inspect() string         { return "Proxy" }

// NewProxyHandler creates a ProxyHandler from an ObjectMap
func NewProxyHandler(handlerObj *object.ObjectMap) *ProxyHandler {
	handler := &ProxyHandler{}

	if get, ok := handlerObj.Get("get"); ok {
		handler.Get = get
	}
	if set, ok := handlerObj.Get("set"); ok {
		handler.Set = set
	}
	if has, ok := handlerObj.Get("has"); ok {
		handler.Has = has
	}
	if del, ok := handlerObj.Get("deleteProperty"); ok {
		handler.DeleteProperty = del
	}
	if apply, ok := handlerObj.Get("apply"); ok {
		handler.Apply = apply
	}
	if construct, ok := handlerObj.Get("construct"); ok {
		handler.Construct = construct
	}
	if getProto, ok := handlerObj.Get("getPrototypeOf"); ok {
		handler.GetPrototypeOf = getProto
	}
	if setProto, ok := handlerObj.Get("setPrototypeOf"); ok {
		handler.SetPrototypeOf = setProto
	}
	if isExt, ok := handlerObj.Get("isExtensible"); ok {
		handler.IsExtensible = isExt
	}
	if prevExt, ok := handlerObj.Get("preventExtensions"); ok {
		handler.PreventExtensions = prevExt
	}
	if getOwnProp, ok := handlerObj.Get("getOwnPropertyDescriptor"); ok {
		handler.GetOwnPropertyDescriptor = getOwnProp
	}
	if defProp, ok := handlerObj.Get("defineProperty"); ok {
		handler.DefineProperty = defProp
	}
	if ownKeys, ok := handlerObj.Get("ownKeys"); ok {
		handler.OwnKeys = ownKeys
	}

	return handler
}

// ProxyGet handles property access on a proxy
func ProxyGet(proxy *ProxyObject, property string, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'get' on a proxy that has been revoked")
	}

	if proxy.Handler.Get != nil {
		targetObj := proxy.Target
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.Get, []object.Object{targetObj, propObj, proxy}, env)
	}

	// Default behavior: get from target
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if val, found := objMap.Get(property); found {
			return val
		}
	}
	return UNDEFINED
}

// ProxySet handles property assignment on a proxy
func ProxySet(proxy *ProxyObject, property string, value object.Object, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'set' on a proxy that has been revoked")
	}

	if proxy.Handler.Set != nil {
		targetObj := proxy.Target
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.Set, []object.Object{targetObj, propObj, value, proxy}, env)
	}

	// Default behavior: set on target
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		objMap.Set(property, value)
		return TRUE
	}
	return FALSE
}

// ProxyHas handles 'in' operator on a proxy
func ProxyHas(proxy *ProxyObject, property string, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'has' on a proxy that has been revoked")
	}

	if proxy.Handler.Has != nil {
		targetObj := proxy.Target
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.Has, []object.Object{targetObj, propObj}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if _, found := objMap.Get(property); found {
			return TRUE
		}
	}
	return FALSE
}

// ProxyDeleteProperty handles delete operator on a proxy
func ProxyDeleteProperty(proxy *ProxyObject, property string, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'deleteProperty' on a proxy that has been revoked")
	}

	if proxy.Handler.DeleteProperty != nil {
		targetObj := proxy.Target
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.DeleteProperty, []object.Object{targetObj, propObj}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		delete(objMap.Pairs, property)
		return TRUE
	}
	return FALSE
}

// ProxyApply handles function call on a proxy
func ProxyApply(proxy *ProxyObject, thisArg object.Object, args []object.Object, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'apply' on a proxy that has been revoked")
	}

	if proxy.Handler.Apply != nil {
		argsArray := &object.Array{Elements: args}
		return applyProxyTrap(proxy.Handler.Apply, []object.Object{proxy.Target, thisArg, argsArray}, env)
	}

	// Default behavior: call target function
	if fn, ok := proxy.Target.(*object.Function); ok {
		extendedEnv := extendFunctionEnv(fn, args)
		if thisObj, ok := thisArg.(*object.ObjectMap); ok {
			extendedEnv.Set("this", thisObj)
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	}

	return newError("proxy target is not callable")
}

// ProxyConstruct handles new operator on a proxy
func ProxyConstruct(proxy *ProxyObject, args []object.Object, newTarget object.Object, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'construct' on a proxy that has been revoked")
	}

	targetNew := newTarget
	if targetNew == nil {
		targetNew = proxy
	}

	if proxy.Handler.Construct != nil {
		argsArray := &object.Array{Elements: args}
		return applyProxyTrap(proxy.Handler.Construct, []object.Object{proxy.Target, argsArray, targetNew}, env)
	}

	// Default behavior: construct using target
	if class, ok := proxy.Target.(*Class); ok {
		return createClassInstance(class, args, class.Env)
	}

	return newError("proxy target is not a constructor")
}

// ProxyGetPrototypeOf handles Object.getPrototypeOf on a proxy
func ProxyGetPrototypeOf(proxy *ProxyObject, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'getPrototypeOf' on a proxy that has been revoked")
	}

	if proxy.Handler.GetPrototypeOf != nil {
		return applyProxyTrap(proxy.Handler.GetPrototypeOf, []object.Object{proxy.Target}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if objMap.Prototype != nil {
			return objMap.Prototype
		}
	}
	return NULL
}

// ProxySetPrototypeOf handles Object.setPrototypeOf on a proxy
func ProxySetPrototypeOf(proxy *ProxyObject, prototype object.Object, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'setPrototypeOf' on a proxy that has been revoked")
	}

	if proxy.Handler.SetPrototypeOf != nil {
		return applyProxyTrap(proxy.Handler.SetPrototypeOf, []object.Object{proxy.Target, prototype}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if protoMap, ok := prototype.(*object.ObjectMap); ok {
			objMap.Prototype = protoMap
			return TRUE
		}
		if prototype.Type() == object.NULL_OBJ {
			objMap.Prototype = nil
			return TRUE
		}
	}
	return FALSE
}

// ProxyIsExtensible handles Object.isExtensible on a proxy
func ProxyIsExtensible(proxy *ProxyObject, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'isExtensible' on a proxy that has been revoked")
	}

	if proxy.Handler.IsExtensible != nil {
		return applyProxyTrap(proxy.Handler.IsExtensible, []object.Object{proxy.Target}, env)
	}

	// Default behavior
	return TRUE
}

// ProxyPreventExtensions handles Object.preventExtensions on a proxy
func ProxyPreventExtensions(proxy *ProxyObject, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'preventExtensions' on a proxy that has been revoked")
	}

	if proxy.Handler.PreventExtensions != nil {
		return applyProxyTrap(proxy.Handler.PreventExtensions, []object.Object{proxy.Target}, env)
	}

	// Default behavior
	return TRUE
}

// ProxyGetOwnPropertyDescriptor handles Object.getOwnPropertyDescriptor on a proxy
func ProxyGetOwnPropertyDescriptor(proxy *ProxyObject, property string, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'getOwnPropertyDescriptor' on a proxy that has been revoked")
	}

	if proxy.Handler.GetOwnPropertyDescriptor != nil {
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.GetOwnPropertyDescriptor, []object.Object{proxy.Target, propObj}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if val, found := objMap.Get(property); found {
			descriptor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			descriptor.Set("value", val)
			descriptor.Set("writable", TRUE)
			descriptor.Set("enumerable", TRUE)
			descriptor.Set("configurable", TRUE)
			return descriptor
		}
	}
	return UNDEFINED
}

// ProxyDefineProperty handles Object.defineProperty on a proxy
func ProxyDefineProperty(proxy *ProxyObject, property string, descriptor object.Object, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'defineProperty' on a proxy that has been revoked")
	}

	if proxy.Handler.DefineProperty != nil {
		propObj := &object.String{Value: property}
		return applyProxyTrap(proxy.Handler.DefineProperty, []object.Object{proxy.Target, propObj, descriptor}, env)
	}

	// Default behavior
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		if descMap, ok := descriptor.(*object.ObjectMap); ok {
			if val, found := descMap.Get("value"); found {
				objMap.Set(property, val)
				return TRUE
			}
		}
	}
	return FALSE
}

// ProxyOwnKeys handles Object.keys / Object.getOwnPropertyNames on a proxy
func ProxyOwnKeys(proxy *ProxyObject, env *object.Environment) object.Object {
	if proxy.Revoked {
		return newError("Cannot perform 'ownKeys' on a proxy that has been revoked")
	}

	if proxy.Handler.OwnKeys != nil {
		return applyProxyTrap(proxy.Handler.OwnKeys, []object.Object{proxy.Target}, env)
	}

	// Default behavior
	keys := []object.Object{}
	if objMap, ok := proxy.Target.(*object.ObjectMap); ok {
		keyNames := make([]string, 0, len(objMap.Pairs))
		for key := range objMap.Pairs {
			keyNames = append(keyNames, key)
		}
		sort.Strings(keyNames)
		for _, key := range keyNames {
			keys = append(keys, &object.String{Value: key})
		}
	}
	return &object.Array{Elements: keys}
}

// applyProxyTrap applies a proxy trap function
func applyProxyTrap(trap object.Object, args []object.Object, _ *object.Environment) object.Object {
	switch fn := trap.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("proxy trap is not a function")
	}
}

// initProxy initializes the Proxy constructor
func initProxy() *object.ObjectMap {
	proxyObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Proxy.revocable(target, handler)
	proxyObj.Set("revocable", &object.Builtin{
		Name: "revocable",
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("Proxy.revocable requires 2 arguments")
			}
			target := args[0]
			handlerArg := args[1]

			handlerObj, ok := handlerArg.(*object.ObjectMap)
			if !ok {
				return newError("Proxy handler must be an object")
			}

			handler := NewProxyHandler(handlerObj)
			proxy := &ProxyObject{
				Target:  target,
				Handler: handler,
				Revoked: false,
			}

			result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			result.Set("proxy", proxy)
			result.Set("revoke", &object.Builtin{
				Name: "revoke",
				Fn: func(args ...object.Object) object.Object {
					proxy.Revoked = true
					return UNDEFINED
				},
			})

			return result
		},
	})

	return proxyObj
}

// createProxy creates a new Proxy object
func createProxy(args []object.Object) object.Object {
	if len(args) < 2 {
		return newError("Proxy requires 2 arguments: target and handler")
	}

	target := args[0]
	handlerArg := args[1]

	handlerObj, ok := handlerArg.(*object.ObjectMap)
	if !ok {
		return newError("Proxy handler must be an object")
	}

	handler := NewProxyHandler(handlerObj)
	return &ProxyObject{
		Target:  target,
		Handler: handler,
		Revoked: false,
	}
}

// isProxyConstructor checks if an ObjectMap is the Proxy constructor
func isProxyConstructor(obj *object.ObjectMap) bool {
	// Check if it has the revocable method which is unique to Proxy
	if _, ok := obj.Get("revocable"); ok {
		// Make sure it's the Proxy by checking it doesn't have other common methods
		if _, hasGet := obj.Get("get"); !hasGet {
			if _, hasSet := obj.Get("set"); !hasSet {
				return true
			}
		}
	}
	return false
}
