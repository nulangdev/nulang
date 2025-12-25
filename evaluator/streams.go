package evaluator

import (
	"bytes"
	"sync"

	"github.com/nulang/nulang/object"
)

// Stream represents a basic stream
type Stream struct {
	buffer    bytes.Buffer
	mutex     sync.Mutex
	listeners map[string][]object.Object
	readable  bool
	writable  bool
	ended     bool
	paused    bool
}

// initStreamModule creates the stream module
func initStreamModule() *object.ObjectMap {
	streamModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Readable stream constructor
	streamModule.Set("Readable", &object.Builtin{Name: "Readable", Fn: func(args ...object.Object) object.Object {
		return createReadableStream()
	}})

	// Writable stream constructor
	streamModule.Set("Writable", &object.Builtin{Name: "Writable", Fn: func(args ...object.Object) object.Object {
		return createWritableStream()
	}})

	// Transform stream constructor
	streamModule.Set("Transform", &object.Builtin{Name: "Transform", Fn: func(args ...object.Object) object.Object {
		return createTransformStream()
	}})

	// PassThrough stream
	streamModule.Set("PassThrough", &object.Builtin{Name: "PassThrough", Fn: func(args ...object.Object) object.Object {
		return createPassThroughStream()
	}})

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
				}
			}
		}
		
		return args[len(args)-1]
	}})

	return streamModule
}

// createReadableStream creates a new readable stream object
func createReadableStream() *object.ObjectMap {
	stream := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	buffer := &bytes.Buffer{}
	listeners := make(map[string][]object.Object)
	paused := false
	ended := false

	// readable property
	stream.Set("readable", TRUE)
	stream.Set("writable", FALSE)

	// read(size?)
	stream.Set("read", &object.Builtin{Name: "read", Fn: func(args ...object.Object) object.Object {
		if ended && buffer.Len() == 0 {
			return NULL
		}

		size := buffer.Len()
		if len(args) > 0 {
			if num, ok := args[0].(*object.Number); ok {
				size = int(num.Value)
				if size > buffer.Len() {
					size = buffer.Len()
				}
			}
		}

		if size == 0 {
			return NULL
		}

		data := make([]byte, size)
		buffer.Read(data)
		return &object.Buffer{Data: data}
	}})

	// push(chunk)
	stream.Set("push", &object.Builtin{Name: "push", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}

		if args[0] == NULL {
			// Signal end of stream
			ended = true
			emitEvent(listeners, "end", []object.Object{})
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

		buffer.Write(data)
		
		if !paused {
			emitEvent(listeners, "data", []object.Object{&object.Buffer{Data: data}})
		}
		
		return TRUE
	}})

	// on(event, callback)
	stream.Set("on", &object.Builtin{Name: "on", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return stream
		}
		
		event := objectToString(args[0])
		if event != "" {
			listeners[event] = append(listeners[event], args[1])
		}
		return stream
	}})

	// once(event, callback)
	stream.Set("once", &object.Builtin{Name: "once", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return stream
		}
		event := objectToString(args[0])
		// For simplicity, once is treated same as on
		listeners[event] = append(listeners[event], args[1])
		return stream
	}})

	// removeListener(event, callback)
	stream.Set("removeListener", &object.Builtin{Name: "removeListener", Fn: func(args ...object.Object) object.Object {
		return stream
	}})

	// pause()
	stream.Set("pause", &object.Builtin{Name: "pause", Fn: func(args ...object.Object) object.Object {
		paused = true
		return stream
	}})

	// resume()
	stream.Set("resume", &object.Builtin{Name: "resume", Fn: func(args ...object.Object) object.Object {
		paused = false
		return stream
	}})

	// isPaused()
	stream.Set("isPaused", &object.Builtin{Name: "isPaused", Fn: func(args ...object.Object) object.Object {
		return nativeBoolToBooleanObject(paused)
	}})

	// pipe(destination)
	stream.Set("pipe", &object.Builtin{Name: "pipe", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("pipe requires a destination")
		}
		
		dest, ok := args[0].(*object.ObjectMap)
		if !ok {
			return newError("pipe destination must be a writable stream")
		}

		// Add data listener to pipe data to destination
		var dataListener object.Object
		dataListener = &object.Builtin{Name: "dataListener", Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				if writeMethod, ok := dest.Get("write"); ok {
					if builtin, ok := writeMethod.(*object.Builtin); ok {
						builtin.Fn(args[0])
					}
				}
			}
			return UNDEFINED
		}}
		
		listeners["data"] = append(listeners["data"], dataListener)

		// Handle end
		var endListener object.Object
		endListener = &object.Builtin{Name: "endListener", Fn: func(args ...object.Object) object.Object {
			if endMethod, ok := dest.Get("end"); ok {
				if builtin, ok := endMethod.(*object.Builtin); ok {
					builtin.Fn()
				}
			}
			return UNDEFINED
		}}
		
		listeners["end"] = append(listeners["end"], endListener)

		return dest
	}})

	// unpipe(destination?)
	stream.Set("unpipe", &object.Builtin{Name: "unpipe", Fn: func(args ...object.Object) object.Object {
		return stream
	}})

	// setEncoding(encoding)
	stream.Set("setEncoding", &object.Builtin{Name: "setEncoding", Fn: func(args ...object.Object) object.Object {
		return stream
	}})

	// destroy()
	stream.Set("destroy", &object.Builtin{Name: "destroy", Fn: func(args ...object.Object) object.Object {
		ended = true
		buffer.Reset()
		return stream
	}})

	return stream
}

// createWritableStream creates a new writable stream object
func createWritableStream() *object.ObjectMap {
	stream := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	buffer := &bytes.Buffer{}
	listeners := make(map[string][]object.Object)
	ended := false

	// Properties
	stream.Set("readable", FALSE)
	stream.Set("writable", TRUE)

	// write(chunk, encoding?, callback?)
	stream.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
		if ended {
			return FALSE
		}

		if len(args) < 1 {
			return FALSE
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

		buffer.Write(data)
		emitEvent(listeners, "drain", []object.Object{})
		
		return TRUE
	}})

	// end(chunk?, encoding?, callback?)
	stream.Set("end", &object.Builtin{Name: "end", Fn: func(args ...object.Object) object.Object {
		if len(args) > 0 && args[0] != nil {
			// Write final chunk
			var data []byte
			switch v := args[0].(type) {
			case *object.String:
				data = []byte(v.Value)
			case *object.Buffer:
				data = v.Data
			default:
				data = []byte(objectToString(args[0]))
			}
			buffer.Write(data)
		}

		ended = true
		emitEvent(listeners, "finish", []object.Object{})
		
		return stream
	}})

	// on(event, callback)
	stream.Set("on", &object.Builtin{Name: "on", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return stream
		}
		event := objectToString(args[0])
		listeners[event] = append(listeners[event], args[1])
		return stream
	}})

	// cork()
	stream.Set("cork", &object.Builtin{Name: "cork", Fn: func(args ...object.Object) object.Object {
		return stream
	}})

	// uncork()
	stream.Set("uncork", &object.Builtin{Name: "uncork", Fn: func(args ...object.Object) object.Object {
		return stream
	}})

	// destroy()
	stream.Set("destroy", &object.Builtin{Name: "destroy", Fn: func(args ...object.Object) object.Object {
		ended = true
		buffer.Reset()
		return stream
	}})

	// _getData() - internal method to get buffer contents
	stream.Set("_getData", &object.Builtin{Name: "_getData", Fn: func(args ...object.Object) object.Object {
		return &object.Buffer{Data: buffer.Bytes()}
	}})

	// toString()
	stream.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: buffer.String()}
	}})

	return stream
}

// createTransformStream creates a transform stream
func createTransformStream() *object.ObjectMap {
	stream := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	outputBuffer := &bytes.Buffer{}
	listeners := make(map[string][]object.Object)
	var transformFn object.Object

	// Properties
	stream.Set("readable", TRUE)
	stream.Set("writable", TRUE)

	// _transform(chunk, encoding, callback)
	stream.Set("_transform", &object.Builtin{Name: "_transform", Fn: func(args ...object.Object) object.Object {
		if len(args) > 0 {
			transformFn = args[0]
		}
		return stream
	}})

	// write(chunk)
	stream.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
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

		// Apply transform if set
		if transformFn != nil {
			if fn, ok := transformFn.(*object.Function); ok {
				result := applyFunction(fn, []object.Object{&object.Buffer{Data: data}})
				if buf, ok := result.(*object.Buffer); ok {
					data = buf.Data
				} else if str, ok := result.(*object.String); ok {
					data = []byte(str.Value)
				}
			}
		}

		outputBuffer.Write(data)
		emitEvent(listeners, "data", []object.Object{&object.Buffer{Data: data}})
		
		return TRUE
	}})

	// read(size?)
	stream.Set("read", &object.Builtin{Name: "read", Fn: func(args ...object.Object) object.Object {
		if outputBuffer.Len() == 0 {
			return NULL
		}
		data := outputBuffer.Bytes()
		outputBuffer.Reset()
		return &object.Buffer{Data: data}
	}})

	// on(event, callback)
	stream.Set("on", &object.Builtin{Name: "on", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return stream
		}
		event := objectToString(args[0])
		listeners[event] = append(listeners[event], args[1])
		return stream
	}})

	// push(chunk)
	stream.Set("push", &object.Builtin{Name: "push", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}

		if args[0] == NULL {
			emitEvent(listeners, "end", []object.Object{})
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

		outputBuffer.Write(data)
		emitEvent(listeners, "data", []object.Object{&object.Buffer{Data: data}})
		
		return TRUE
	}})

	// pipe(destination)
	stream.Set("pipe", &object.Builtin{Name: "pipe", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("pipe requires a destination")
		}
		
		dest, ok := args[0].(*object.ObjectMap)
		if !ok {
			return newError("pipe destination must be a writable stream")
		}

		var dataListener object.Object
		dataListener = &object.Builtin{Name: "dataListener", Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				if writeMethod, ok := dest.Get("write"); ok {
					if builtin, ok := writeMethod.(*object.Builtin); ok {
						builtin.Fn(args[0])
					}
				}
			}
			return UNDEFINED
		}}
		
		listeners["data"] = append(listeners["data"], dataListener)

		return dest
	}})

	// end()
	stream.Set("end", &object.Builtin{Name: "end", Fn: func(args ...object.Object) object.Object {
		emitEvent(listeners, "finish", []object.Object{})
		return stream
	}})

	return stream
}

// createPassThroughStream creates a pass-through stream
func createPassThroughStream() *object.ObjectMap {
	stream := createTransformStream()
	// PassThrough just passes data without transformation
	return stream
}

// emitEvent emits an event to registered listeners
func emitEvent(listeners map[string][]object.Object, event string, args []object.Object) {
	if handlers, ok := listeners[event]; ok {
		for _, handler := range handlers {
			switch h := handler.(type) {
			case *object.Function:
				// Execute callback directly without using executeTimerCallback to avoid cycle
				callbackEnv := extendFunctionEnvForCallback(h, args)
				Eval(h.Body, callbackEnv)
			case *object.Builtin:
				h.Fn(args...)
			}
		}
	}
}

// extendFunctionEnvForCallback creates env for callback execution
func extendFunctionEnvForCallback(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Set(param.Value, args[i])
		} else {
			env.Set(param.Value, UNDEFINED)
		}
	}
	return env
}
