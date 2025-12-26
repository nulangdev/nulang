package evaluator

import (
	"github.com/nulang/nulang/object"
)

var (
	readableClass    *Class
	writableClass    *Class
	transformClass   *Class
	passThroughClass *Class
)

// initStreamModule creates the stream module
func initStreamModule() *object.ObjectMap {
	// Ensure events module classes are ready
	if eventEmitterClass == nil {
		createEventEmitterClass() 
	}
	
	createStreamClasses()

	streamModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	streamModule.Set("Readable", readableClass)
	streamModule.Set("Writable", writableClass)
	streamModule.Set("Transform", transformClass)
	streamModule.Set("PassThrough", passThroughClass)
	
	// Pipeline utility
	streamModule.Set("pipeline", &object.Builtin{Name: "pipeline", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("pipeline requires at least 2 streams")
		}
		
		// Simple pipeline: connect streams
		for i := 0; i < len(args)-1; i++ {
			source, ok := args[i].(*object.ObjectMap)
			if !ok {
				continue
			}
			dest, ok := args[i+1].(*object.ObjectMap)
			if !ok {
				continue
			}
			
			// Pipe data from source to dest
			if pipeMethod, ok := source.Get("pipe"); ok {
				if builtin, ok := pipeMethod.(*object.Builtin); ok {
					builtin.Fn(dest)
				} else if fn, ok := pipeMethod.(*object.Function); ok {
				    env := object.NewEnclosedEnvironment(fn.Env)
				    Eval(fn.Body, env)
				}
			}
		}
		
		return args[len(args)-1]
	}})

	return streamModule
}

func createStreamClasses() {
	if readableClass != nil {
		return
	}

	// Readable
	readableClass = &Class{
		Name:          "Readable",
		SuperClass:    eventEmitterClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	
	readableClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		// Call super constructor
		if superCtor, ok := eventEmitterClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		
		instance := this.(*object.ObjectMap)
		// Init buffer
		instance.Set("_buffer", &object.Buffer{Data: []byte{}})
		instance.Set("readable", TRUE)
		instance.Set("_ended", FALSE)
		instance.Set("_paused", FALSE)
		
		return UNDEFINED
	}
	
	readableClass.NativeMethods["push"] = func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if len(args) < 1 {
			return FALSE
		}
		
		if args[0] == NULL {
			instance.Set("_ended", TRUE)
			nativeCall(instance, "emit", &object.String{Value: "end"})
			return TRUE
		}
		
		var data []byte
		switch v := args[0].(type) {
		case *object.String:
			data = []byte(v.Value)
		case *object.Buffer:
			data = v.Data
		default:
			data = []byte(objectToString(args[0]))
		}
		
		// Append to buffer
		currentBufObj, _ := instance.Get("_buffer")
		currentBuf := currentBufObj.(*object.Buffer)
		currentBuf.Data = append(currentBuf.Data, data...)
		
		// Emit data if flowing
		pausedObj, _ := instance.Get("_paused")
		if pausedObj == FALSE {
		    nativeCall(instance, "emit", &object.String{Value: "data"}, &object.Buffer{Data: data})
		}
		
		return TRUE
	}
	
	// read(size)
	readableClass.NativeMethods["read"] = func(this object.Object, args ...object.Object) object.Object {
	    instance := this.(*object.ObjectMap)
	    
	    // Call _read
	    if _read, ok := instance.Get("_read"); ok {
            invokeMethod(instance, _read, args...)
	    }
	    
	    // Retrieve from buffer
	    bufObj, _ := instance.Get("_buffer")
	    buf := bufObj.(*object.Buffer)
	    
	    if len(buf.Data) == 0 {
	        ended, _ := instance.Get("_ended")
	        if ended == TRUE {
	            return NULL
	        }
	        return NULL
	    }
	    
	    res := &object.Buffer{Data: buf.Data}
	    buf.Data = []byte{}
	    return res
	}
	
	// pipe
	readableClass.NativeMethods["pipe"] = func(this object.Object, args ...object.Object) object.Object {
	    if len(args) < 1 { return newError("pipe expects dest") }
	    dest := args[0]
	    instance := this.(*object.ObjectMap)
	    
	    // Add 'data' listener that writes to dest
	    dataListener := &object.Builtin{
	        Fn: func(dArgs ...object.Object) object.Object {
	            if len(dArgs) > 0 {
	                 if destObj, ok := dest.(*object.ObjectMap); ok {
	                     nativeCall(destObj, "write", dArgs[0])
	                 }
	            }
	            return UNDEFINED
	        },
	    }
	    
	    nativeCall(instance, "on", &object.String{Value: "data"}, dataListener)
	    
	    // Handle end -> dest.end()
	    endListener := &object.Builtin{
	        Fn: func(dArgs ...object.Object) object.Object {
	             if destObj, ok := dest.(*object.ObjectMap); ok {
                     nativeCall(destObj, "end")
                 }
	             return UNDEFINED
	        },
	    }
	    nativeCall(instance, "on", &object.String{Value: "end"}, endListener)
	    
	    return dest
	}

	// destroy(error?)
	readableClass.NativeMethods["destroy"] = func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		instance.Set("_ended", TRUE)
		instance.Set("destroyed", TRUE)
		// emit close
		nativeCall(instance, "emit", &object.String{Value: "close"})
		return instance
	}
	
	// Writable
	writableClass = &Class{
		Name:          "Writable",
		SuperClass:    eventEmitterClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	
	writableClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		if superCtor, ok := eventEmitterClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		instance := this.(*object.ObjectMap)
		instance.Set("writable", TRUE)
		instance.Set("_writableState", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}) // Use object map for state
		return UNDEFINED
	}
	
	writableClass.NativeMethods["write"] = func(this object.Object, args ...object.Object) object.Object {
	    if len(args) < 1 { return FALSE }
	    instance := this.(*object.ObjectMap)
	    chunk := args[0]
	    
	    cb := &object.Builtin{Fn: func(cArgs ...object.Object) object.Object {
	        nativeCall(instance, "emit", &object.String{Value: "drain"})
	        return UNDEFINED
	    }}
	    
	    if len(args) > 2 {
	        if userCb, ok := args[len(args)-1].(*object.Function); ok {
	            cb = &object.Builtin{Fn: func(cArgs ...object.Object) object.Object {
	                 invokeMethod(instance, userCb, cArgs...)
	                 nativeCall(instance, "emit", &object.String{Value: "drain"})
	                 return UNDEFINED
	            }}
	        }
	    }
	    
	    if _write, ok := instance.Get("_write"); ok {
	        invokeMethod(instance, _write, chunk, &object.String{Value: "utf8"}, cb)
	    } else {
            cb.Fn()
	    }
	    
	    return TRUE
	}
	
	writableClass.NativeMethods["end"] = func(this object.Object, args ...object.Object) object.Object {
	    instance := this.(*object.ObjectMap)
	    if len(args) > 0 {
	        nativeCall(instance, "write", args[0])
	    }
	    
	    nativeCall(instance, "emit", &object.String{Value: "finish"})
	    return instance
	}

	writableClass.NativeMethods["cork"] = func(this object.Object, args ...object.Object) object.Object {
		return TRUE
	}
	
	writableClass.NativeMethods["uncork"] = func(this object.Object, args ...object.Object) object.Object {
		return TRUE
	}

	writableClass.NativeMethods["destroy"] = readableClass.NativeMethods["destroy"]
	
	// Transform
	transformClass = &Class{
		Name:          "Transform",
		SuperClass:    readableClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	
	// Copy Writable methods
	transformClass.NativeMethods["write"] = writableClass.NativeMethods["write"]
	transformClass.NativeMethods["end"] = writableClass.NativeMethods["end"]
    
    transformClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
        // Call super (Readable)
        if superCtor, ok := readableClass.NativeMethods["constructor"]; ok {
            superCtor(this, args...)
        }
        // Also init Writable state manually since we don't inherit it
        if _, ok := writableClass.NativeMethods["constructor"]; ok {
             instance := this.(*object.ObjectMap)
             instance.Set("writable", TRUE)
             instance.Set("_writableState", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
        }
        return UNDEFINED
    }
	
	// Default _write for Transform calls _transform
	transformClass.NativeMethods["_write"] = func(this object.Object, args ...object.Object) object.Object {
	    instance := this.(*object.ObjectMap)
	    // _transform(chunk, encoding, callback)
	    if _transform, ok := instance.Get("_transform"); ok {
	        invokeMethod(instance, _transform, args...)
	    } else {
	        if len(args) > 2 {
	            if cb, ok := args[2].(*object.Builtin); ok { cb.Fn() }
	        }
	    }
	    return UNDEFINED
	}
	
	// PassThrough
	passThroughClass = &Class{
		Name:          "PassThrough",
		SuperClass:    transformClass, 
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
    
    passThroughClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
        if superCtor, ok := transformClass.NativeMethods["constructor"]; ok {
            superCtor(this, args...)
        }
        return UNDEFINED
    }

	// PassThrough _transform just pushes
	passThroughClass.NativeMethods["_transform"] = func(this object.Object, args ...object.Object) object.Object {
	    instance := this.(*object.ObjectMap)
	    chunk := args[0]
	    cb := args[2] // callback
	    
	    nativeCall(instance, "push", chunk)
	    
	    if builtinCb, ok := cb.(*object.Builtin); ok {
	        builtinCb.Fn()
	    }
	    return UNDEFINED
	}
}

// Helpers

func invokeMethod(instance *object.ObjectMap, method object.Object, args ...object.Object) object.Object {
    switch fn := method.(type) {
    case *object.Builtin:
        return fn.Fn(args...)
    case *object.Function:
        env := object.NewEnclosedEnvironment(fn.Env)
        env.Set("this", instance)
        for i, param := range fn.Parameters {
            if i < len(args) {
                env.Set(param.Value, args[i])
            } else {
                env.Set(param.Value, UNDEFINED)
            }
        }
        return Eval(fn.Body, env)
    }
    return UNDEFINED
}

func nativeCall(o *object.ObjectMap, methodName string, args ...object.Object) object.Object {
    if m, ok := o.Get(methodName); ok {
        return invokeMethod(o, m, args...)
    }
    return newError("Method %s not found", methodName)
}
