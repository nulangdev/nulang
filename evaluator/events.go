package evaluator

import (
	"github.com/nulang/nulang/object"
)

var eventEmitterClass *Class

// initEventsModule initializes the 'events' module
func initEventsModule() *object.ObjectMap {
	if eventEmitterClass == nil {
		createEventEmitterClass()
	}

	eventsModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	eventsModule.Set("EventEmitter", eventEmitterClass)
	return eventsModule
}

func createEventEmitterClass() {
	eventEmitterClass = &Class{
		Name:          "EventEmitter",
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(), // Empty environment
	}

	// Static property: defaultMaxListeners
	eventEmitterClass.Static["defaultMaxListeners"] = &object.Number{Value: 10}

	// constructor
	eventEmitterClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("EventEmitter: this is not an object")
		}
		// Initialize _events map
		instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
		instance.Set("_maxListeners", UNDEFINED)
		return UNDEFINED
	}

	// on(event, listener)
	eventEmitterClass.NativeMethods["on"] = func(this object.Object, args ...object.Object) object.Object {
		return addListener(this, args, false)
	}

	// addListener alias
	eventEmitterClass.NativeMethods["addListener"] = eventEmitterClass.NativeMethods["on"]

	// prependListener(event, listener)
	eventEmitterClass.NativeMethods["prependListener"] = func(this object.Object, args ...object.Object) object.Object {
		return addListener(this, args, true)
	}

	// once(event, listener)
	eventEmitterClass.NativeMethods["once"] = func(this object.Object, args ...object.Object) object.Object {
		return addOnceListener(this, args, false)
	}

	// prependOnceListener(event, listener)
	eventEmitterClass.NativeMethods["prependOnceListener"] = func(this object.Object, args ...object.Object) object.Object {
		return addOnceListener(this, args, true)
	}

	// off alias for removeListener
	eventEmitterClass.NativeMethods["off"] = func(this object.Object, args ...object.Object) object.Object {
		return removeListenerMethod(this, args)
	}

	// removeListener(event, listener)
	eventEmitterClass.NativeMethods["removeListener"] = func(this object.Object, args ...object.Object) object.Object {
		return removeListenerMethod(this, args)
	}

	// removeAllListeners([event])
	eventEmitterClass.NativeMethods["removeAllListeners"] = func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		eventsObj, ok := instance.Get("_events")
		if !ok {
			return this
		}
		eventsMap := eventsObj.(*object.ObjectMap)

		if len(args) == 0 {
			// Remove all
			instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
		} else {
			eventName := objectToString(args[0])
			if _, ok := eventsMap.Get(eventName); ok {
				eventsMap.Pairs[eventName] = object.ObjectPair{Key: &object.String{Value: eventName}, Value: &object.Array{Elements: []object.Object{}}}
				delete(eventsMap.Pairs, eventName)
			}
		}

		return this
	}

	// setMaxListeners(n)
	eventEmitterClass.NativeMethods["setMaxListeners"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return this
		}
		instance := this.(*object.ObjectMap)
		instance.Set("_maxListeners", args[0])
		return this
	}

	// getMaxListeners()
	eventEmitterClass.NativeMethods["getMaxListeners"] = func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if max, ok := instance.Get("_maxListeners"); ok && max != UNDEFINED {
			return max
		}
		// Return default
		if def, ok := eventEmitterClass.Static["defaultMaxListeners"]; ok {
			return def
		}
		return &object.Number{Value: 10}
	}

	// listeners(event)
	eventEmitterClass.NativeMethods["listeners"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		return getListeners(this, args[0], true) // true = unwrap
	}

	// rawListeners(event)
	eventEmitterClass.NativeMethods["rawListeners"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		return getListeners(this, args[0], false) // false = keep wrappers
	}

	// listenerCount(event)
	eventEmitterClass.NativeMethods["listenerCount"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: 0}
		}
		list := getListeners(this, args[0], false)
		if arr, ok := list.(*object.Array); ok {
			return &object.Number{Value: float64(len(arr.Elements))}
		}
		return &object.Number{Value: 0}
	}

	// eventNames()
	eventEmitterClass.NativeMethods["eventNames"] = func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		ensureEventsMap(instance)
		events, _ := instance.Get("_events")
		eventsMap := events.(*object.ObjectMap)

		names := &object.Array{Elements: []object.Object{}}
		for _, pair := range eventsMap.Pairs {
			if arr, ok := pair.Value.(*object.Array); ok && len(arr.Elements) > 0 {
				names.Elements = append(names.Elements, pair.Key)
			}
		}
		return names
	}

	// emit(event, ...args)
	eventEmitterClass.NativeMethods["emit"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		instance := this.(*object.ObjectMap)
		eventName := objectToString(args[0])
		emitArgs := args[1:]

		// Handle 'error' event specifically
		if eventName == "error" {
			hasListeners := checkHasListeners(instance, "error")
			if !hasListeners {
				// Throw error
				var err object.Object
				if len(emitArgs) > 0 {
					err = emitArgs[0]
				} else {
					err = newError("Unhandled 'error' event")
				}
				if errObj, ok := err.(*object.Error); ok {
					return errObj
				}
				return newError("Unhandled 'error' event: %s", err.Inspect())
			}
		}

		events, ok := instance.Get("_events")
		if !ok {
			return FALSE
		}
		eventsMap := events.(*object.ObjectMap)

		listenersObj, ok := eventsMap.Get(eventName)
		if !ok {
			return FALSE
		}

		listenersArr, ok := listenersObj.(*object.Array)
		if !ok || len(listenersArr.Elements) == 0 {
			return FALSE
		}

		// Copy listeners to avoid issues during mutation
		listenersCopy := make([]object.Object, len(listenersArr.Elements))
		copy(listenersCopy, listenersArr.Elements)

		for _, handler := range listenersCopy {
			applyHandler(handler, emitArgs)
		}

		return TRUE
	}
}

// Helpers

func addListener(this object.Object, args []object.Object, prepend bool) object.Object {
	if len(args) < 2 {
		return this
	}
	instance := this.(*object.ObjectMap)
	eventName := objectToString(args[0])
	listener := args[1]

	// Emit 'newListener' before adding
	if checkHasListeners(instance, "newListener") {
		emitFn := eventEmitterClass.NativeMethods["emit"]
		emitFn(this, &object.String{Value: "newListener"}, &object.String{Value: eventName}, listener)
	}

	ensureEventsMap(instance)
	events, _ := instance.Get("_events")
	eventsMap := events.(*object.ObjectMap)

	var listeners *object.Array
	if existing, ok := eventsMap.Get(eventName); ok {
		if arr, ok := existing.(*object.Array); ok {
			listeners = arr
		}
	}

	if listeners == nil {
		listeners = &object.Array{Elements: []object.Object{}}
		eventsMap.Set(eventName, listeners)
	}

	if prepend {
		listeners.Elements = append([]object.Object{listener}, listeners.Elements...)
	} else {
		listeners.Elements = append(listeners.Elements, listener)
	}

	return this
}

func addOnceListener(this object.Object, args []object.Object, prepend bool) object.Object {
	if len(args) < 2 {
		return this
	}
	// eventName removed
	originalListener := args[1]

	wrapper := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	wrapper.Set("listener", originalListener) // Store original for removeListener

	// Wrapper function
	wrapperFn := func(wArgs ...object.Object) object.Object {
		// Remove self
		removeListenerMethod(this, []object.Object{args[0], wrapper})
		// Call original
		return applyHandler(originalListener, wArgs)
	}

	// Make wrapper callable
	wrapper.Set("__call__", &object.Builtin{Fn: wrapperFn})

	// Add the wrapper
	return addListener(this, []object.Object{args[0], wrapper}, prepend)
}

func removeListenerMethod(this object.Object, args []object.Object) object.Object {
	if len(args) < 2 {
		return this
	}
	instance := this.(*object.ObjectMap)
	eventName := objectToString(args[0])
	target := args[1]

	events, ok := instance.Get("_events")
	if !ok {
		return this
	}
	eventsMap := events.(*object.ObjectMap)

	listenersObj, ok := eventsMap.Get(eventName)
	if !ok {
		return this
	}

	listenersArr, ok := listenersObj.(*object.Array)
	if !ok {
		return this
	}

	newElements := []object.Object{}
	found := false
	var removedListener object.Object

	for _, l := range listenersArr.Elements {
		if found {
			newElements = append(newElements, l)
			continue
		}

		match := false
		if l == target {
			match = true
		} else if lObj, ok := l.(*object.ObjectMap); ok {
			// Check if it's a wrapper with .listener == target
			if orig, ok := lObj.Get("listener"); ok && orig == target {
				match = true
			}
		}

		if match {
			found = true
			removedListener = l
			continue
		}
		newElements = append(newElements, l)
	}

	listenersArr.Elements = newElements
	
	if len(newElements) == 0 {
		delete(eventsMap.Pairs, eventName)
	}

	if found && checkHasListeners(instance, "removeListener") {
		emitFn := eventEmitterClass.NativeMethods["emit"]
		emitFn(this, &object.String{Value: "removeListener"}, &object.String{Value: eventName}, removedListener)
	}

	return this
}

func getListeners(this object.Object, eventArg object.Object, unwrap bool) object.Object {
	instance := this.(*object.ObjectMap)
	eventName := objectToString(eventArg)
	
	ensureEventsMap(instance)
	events, _ := instance.Get("_events")
	eventsMap := events.(*object.ObjectMap)
	
	existing, ok := eventsMap.Get(eventName)
	if !ok {
		return &object.Array{Elements: []object.Object{}}
	}
	arr, ok := existing.(*object.Array)
	if !ok {
		return &object.Array{Elements: []object.Object{}}
	}
	
	result := &object.Array{Elements: make([]object.Object, len(arr.Elements))}
	for i, l := range arr.Elements {
		val := l
		if unwrap {
			if lObj, ok := l.(*object.ObjectMap); ok {
				if orig, ok := lObj.Get("listener"); ok {
					val = orig
				}
			}
		}
		result.Elements[i] = val
	}
	return result
}

func checkHasListeners(instance *object.ObjectMap, eventName string) bool {
	events, ok := instance.Get("_events")
	if !ok { return false }
	eventsMap := events.(*object.ObjectMap)
	l, ok := eventsMap.Get(eventName)
	if !ok { return false }
	arr, ok := l.(*object.Array)
	return ok && len(arr.Elements) > 0
}

func ensureEventsMap(instance *object.ObjectMap) {
	if _, ok := instance.Get("_events"); !ok {
		instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
	}
}

func applyHandler(handler object.Object, args []object.Object) object.Object {
	switch fn := handler.(type) {
	case *object.Function:
		env := extendFunctionEnv(fn, args)
		return Eval(fn.Body, env)
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.ObjectMap:
		if callFn, ok := fn.Get("__call__"); ok {
			return applyHandler(callFn, args)
		}
	}
	return UNDEFINED
}

// createEventEmitter creates a new EventEmitter instance
func createEventEmitter() *object.ObjectMap {
	if eventEmitterClass == nil {
		createEventEmitterClass()
	}
	
	instance := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	// Call constructor
	if constructor, ok := eventEmitterClass.NativeMethods["constructor"]; ok {
		constructor(instance)
	}
	
	// Add all methods
	for name, method := range eventEmitterClass.NativeMethods {
		methodName := name
		nativeMethod := method
		instance.Set(methodName, &object.Builtin{
			Name: methodName,
			Fn: func(args ...object.Object) object.Object {
				return nativeMethod(instance, args...)
			},
		})
	}
	
	return instance
}

// emitEvent emits an event on an EventEmitter instance
func emitEvent(instance *object.ObjectMap, eventName string, args ...object.Object) object.Object {
	if emit, ok := instance.Get("emit"); ok {
		if emitBuiltin, ok := emit.(*object.Builtin); ok {
			allArgs := make([]object.Object, len(args)+1)
			allArgs[0] = &object.String{Value: eventName}
			copy(allArgs[1:], args)
			return emitBuiltin.Fn(allArgs...)
		}
	}
	return FALSE
}

// addEventListener adds an event listener to an EventEmitter instance
func addEventListener(instance *object.ObjectMap, eventName string, listener object.Object) object.Object {
	if on, ok := instance.Get("on"); ok {
		if onBuiltin, ok := on.(*object.Builtin); ok {
			return onBuiltin.Fn(&object.String{Value: eventName}, listener)
		}
	}
	return instance
}

