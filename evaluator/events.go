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

	// constructor
	eventEmitterClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("EventEmitter: this is not an object")
		}
		// Initialize _events map
		instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
		return UNDEFINED
	}

	// on(event, listener)
	eventEmitterClass.NativeMethods["on"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 2 {
			return this
		}
		instance := this.(*object.ObjectMap)
		eventName := objectToString(args[0])
		listener := args[1]

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
		}
		listeners.Elements = append(listeners.Elements, listener)
		eventsMap.Set(eventName, listeners)

		return this
	}

	// addListener alias
	eventEmitterClass.NativeMethods["addListener"] = eventEmitterClass.NativeMethods["on"]

	// emit(event, ...args)
	eventEmitterClass.NativeMethods["emit"] = func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		instance := this.(*object.ObjectMap)
		eventName := objectToString(args[0])
		emitArgs := args[1:]

		events, ok := instance.Get("_events")
		if !ok {
			return FALSE
		}
		eventsMap := events.(*object.ObjectMap)
		
		listenersObj, ok := eventsMap.Get(eventName)
		if !ok {
			return FALSE // No listeners
		}
		
		listenersArr, ok := listenersObj.(*object.Array)
		if !ok || len(listenersArr.Elements) == 0 {
			return FALSE
		}

		// Copy listeners to avoid issues during mutation
		listenersCopy := make([]object.Object, len(listenersArr.Elements))
		copy(listenersCopy, listenersArr.Elements)

		for _, handler := range listenersCopy {
			// Call handler
			applyHandler(handler, emitArgs)
		}
		
		return TRUE
	}
    
    // once(event, listener)
    eventEmitterClass.NativeMethods["once"] = func(this object.Object, args ...object.Object) object.Object {
        if len(args) < 2 {
            return this
        }
        instance := this.(*object.ObjectMap)
        eventName := objectToString(args[0])
        originalListener := args[1]
        
        // We need to create a wrapper function that calls 'off' then 'original'
        // Since we are in Go, we create a Builtin closure
        var wrapper *object.Builtin
        wrapper = &object.Builtin{
            Fn: func(wArgs ...object.Object) object.Object {
                // Remove self (wrapper)
                removeListener(instance, eventName, wrapper)
                
                // Call original
                return applyHandler(originalListener, wArgs)
            },
        }
        
        // Add wrapper
        // Use 'on' method logic directly or call 'on'?
        // Calling 'on' is safer
        if onMethod, ok := instance.Get("on"); ok {
             if builtin, ok := onMethod.(*object.Builtin); ok {
                 builtin.Fn(&object.String{Value: eventName}, wrapper)
             }
        }
        
        return this
    }

    // removeListener(event, listener)
    eventEmitterClass.NativeMethods["removeListener"] = func(this object.Object, args ...object.Object) object.Object {
        if len(args) < 2 {
            return this
        }
        instance := this.(*object.ObjectMap)
        eventName := objectToString(args[0])
        target := args[1]
        
        removeListener(instance, eventName, target)
        return this
    }
    
    // removeAllListeners
    eventEmitterClass.NativeMethods["removeAllListeners"] = func(this object.Object, args ...object.Object) object.Object {
        instance := this.(*object.ObjectMap)
        if len(args) > 0 {
            eventName := objectToString(args[0])
            if events, ok := instance.Get("_events"); ok {
                events.(*object.ObjectMap).Pairs[eventName] = object.ObjectPair{Key: &object.String{Value:eventName}, Value: &object.Array{Elements: []object.Object{}}} // Clear specific
            }
        } else {
             instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
        }
        return this
    }
}

func ensureEventsMap(instance *object.ObjectMap) {
	if _, ok := instance.Get("_events"); !ok {
		instance.Set("_events", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
	}
}

func removeListener(instance *object.ObjectMap, eventName string, target object.Object) {
    events, ok := instance.Get("_events")
    if !ok { return }
    eventsMap := events.(*object.ObjectMap)
    
    listenersObj, ok := eventsMap.Get(eventName)
    if !ok { return }
    
    listenersArr, ok := listenersObj.(*object.Array)
    if !ok { return }
    
    newElements := []object.Object{}
    found := false
    for _, l := range listenersArr.Elements {
        if !found && l == target {
            found = true
            continue
        }
        newElements = append(newElements, l)
    }
    listenersArr.Elements = newElements
}

func applyHandler(handler object.Object, args []object.Object) object.Object {
    switch fn := handler.(type) {
    case *object.Function:
        env := extendFunctionEnv(fn, args)
        return Eval(fn.Body, env)
    case *object.Builtin:
        return fn.Fn(args...)
    }
    return UNDEFINED
}
