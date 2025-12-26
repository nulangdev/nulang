package evaluator

import (
	"encoding/base64"
	"time"

	"github.com/nulang/nulang/object"
)

// BlobClass is the Blob constructor (Node.js compatible)
var BlobClass *Class

// FileClass is the File constructor (Node.js compatible - extends Blob)
var FileClass *Class

// initBlobFileClasses initializes the Blob and File classes
func initBlobFileClasses() {
	if BlobClass != nil {
		return
	}

	// Define Blob class
	BlobClass = &Class{
		Name:          "Blob",
		Properties:    make(map[string]object.Object),
		Methods:       make(map[string]*object.Function),
		Getters:       make(map[string]*object.Function),
		Setters:       make(map[string]*object.Function),
		Static:        make(map[string]object.Object),
		NativeMethods: make(map[string]NativeMethod),
	}

	// Blob constructor
	BlobClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("Blob: invalid this context")
		}

		// Process parts argument
		var data []byte
		if len(args) > 0 {
			if arr, ok := args[0].(*object.Array); ok {
				for _, elem := range arr.Elements {
					switch e := elem.(type) {
					case *object.String:
						data = append(data, []byte(e.Value)...)
					case *object.ObjectMap:
						// Could be another Blob or ArrayBuffer
						if blobData, exists := e.Get("_data"); exists {
							if buffer, ok := blobData.(*object.Buffer); ok {
								data = append(data, buffer.Data...)
							}
						}
					case *object.Buffer:
						data = append(data, e.Data...)
					case *object.Array:
						// Uint8Array-like
						for _, b := range e.Elements {
							if num, ok := b.(*object.Number); ok {
								data = append(data, byte(int(num.Value)%256))
							}
						}
					}
				}
			}
		}

		// Process options argument
		contentType := ""
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				if typeVal, exists := opts.Get("type"); exists {
					if str, ok := typeVal.(*object.String); ok {
						contentType = str.Value
					}
				}
			}
		}

		// Store internal data
		instance.Set("_data", &object.Buffer{Data: data})
		instance.Set("_type", &object.String{Value: contentType})

		return UNDEFINED
	}

	// Blob.prototype.size (getter)
	BlobClass.NativeMethods["size"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return &object.Number{Value: 0}
		}
		if data, exists := instance.Get("_data"); exists {
			if buffer, ok := data.(*object.Buffer); ok {
				return &object.Number{Value: float64(len(buffer.Data))}
			}
		}
		return &object.Number{Value: 0}
	}

	// Blob.prototype.type (getter)
	BlobClass.NativeMethods["type"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return &object.String{Value: ""}
		}
		if typeVal, exists := instance.Get("_type"); exists {
			if str, ok := typeVal.(*object.String); ok {
				return str
			}
		}
		return &object.String{Value: ""}
	}

	// Blob.prototype.text() - returns Promise<string>
	BlobClass.NativeMethods["text"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return blobCreateRejectedPromise(newError("Blob.text: invalid this context"))
		}
		if data, exists := instance.Get("_data"); exists {
			if buffer, ok := data.(*object.Buffer); ok {
				return blobCreateResolvedPromise(&object.String{Value: string(buffer.Data)})
			}
		}
		return blobCreateResolvedPromise(&object.String{Value: ""})
	}

	// Blob.prototype.arrayBuffer() - returns Promise<ArrayBuffer>
	BlobClass.NativeMethods["arrayBuffer"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return blobCreateRejectedPromise(newError("Blob.arrayBuffer: invalid this context"))
		}
		if data, exists := instance.Get("_data"); exists {
			if buffer, ok := data.(*object.Buffer); ok {
				return blobCreateResolvedPromise(buffer)
			}
		}
		return blobCreateResolvedPromise(&object.Buffer{Data: []byte{}})
	}

	// Blob.prototype.slice(start, end, contentType)
	BlobClass.NativeMethods["slice"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("Blob.slice: invalid this context")
		}

		var data []byte
		if dataVal, exists := instance.Get("_data"); exists {
			if buffer, ok := dataVal.(*object.Buffer); ok {
				data = buffer.Data
			}
		}

		start := 0
		end := len(data)
		contentType := ""

		if len(args) > 0 {
			if num, ok := args[0].(*object.Number); ok {
				start = int(num.Value)
				if start < 0 {
					start = len(data) + start
				}
				if start < 0 {
					start = 0
				}
			}
		}

		if len(args) > 1 {
			if num, ok := args[1].(*object.Number); ok {
				end = int(num.Value)
				if end < 0 {
					end = len(data) + end
				}
			}
		}

		if len(args) > 2 {
			if str, ok := args[2].(*object.String); ok {
				contentType = str.Value
			}
		}

		// Normalize indices
		if start > len(data) {
			start = len(data)
		}
		if end > len(data) {
			end = len(data)
		}
		if start > end {
			start = end
		}

		slicedData := data[start:end]

		// Create new Blob with sliced data
		newBlob := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		newBlob.Set("_data", &object.Buffer{Data: slicedData})
		newBlob.Set("_type", &object.String{Value: contentType})

		// Add methods
		for name, method := range BlobClass.NativeMethods {
			if name != "constructor" && name != "size" && name != "type" {
				newBlob.Set(name, bindNativeMethod(method, newBlob))
			}
		}

		// Add size and type as direct properties
		newBlob.Set("size", &object.Number{Value: float64(len(slicedData))})
		newBlob.Set("type", &object.String{Value: contentType})

		return newBlob
	}

	// Blob.prototype.stream() - returns ReadableStream (simplified)
	BlobClass.NativeMethods["stream"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("Blob.stream: invalid this context")
		}

		// Create a simplified stream object
		stream := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

		if data, exists := instance.Get("_data"); exists {
			stream.Set("_data", data)
		}

		return stream
	}

	// Blob.prototype.toDataURL() - returns a data URL (non-standard but useful)
	BlobClass.NativeMethods["toDataURL"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("Blob.toDataURL: invalid this context")
		}

		// Get content type
		contentType := "application/octet-stream"
		if typeVal, exists := instance.Get("_type"); exists {
			if str, ok := typeVal.(*object.String); ok && str.Value != "" {
				contentType = str.Value
			}
		}

		// Convert to base64
		base64Data := blobToBase64(instance)
		dataURL := "data:" + contentType + ";base64," + base64Data

		return &object.String{Value: dataURL}
	}

	// Blob.prototype.bytes() - returns Promise<Uint8Array> (Node.js 20+)
	BlobClass.NativeMethods["bytes"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return blobCreateRejectedPromise(newError("Blob.bytes: invalid this context"))
		}

		if data, exists := instance.Get("_data"); exists {
			if buffer, ok := data.(*object.Buffer); ok {
				// Convert to Array of numbers (Uint8Array-like)
				elements := make([]object.Object, len(buffer.Data))
				for i, b := range buffer.Data {
					elements[i] = &object.Number{Value: float64(b)}
				}
				return blobCreateResolvedPromise(&object.Array{Elements: elements})
			}
		}
		return blobCreateResolvedPromise(&object.Array{Elements: []object.Object{}})
	}

	// Define File class (extends Blob)
	FileClass = &Class{
		Name:          "File",
		SuperClass:    BlobClass,
		Properties:    make(map[string]object.Object),
		Methods:       make(map[string]*object.Function),
		Getters:       make(map[string]*object.Function),
		Setters:       make(map[string]*object.Function),
		Static:        make(map[string]object.Object),
		NativeMethods: make(map[string]NativeMethod),
	}

	// File constructor
	FileClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("File: invalid this context")
		}

		// Process fileBits argument (array of parts)
		var data []byte
		if len(args) > 0 {
			if arr, ok := args[0].(*object.Array); ok {
				for _, elem := range arr.Elements {
					switch e := elem.(type) {
					case *object.String:
						data = append(data, []byte(e.Value)...)
					case *object.ObjectMap:
						if blobData, exists := e.Get("_data"); exists {
							if buffer, ok := blobData.(*object.Buffer); ok {
								data = append(data, buffer.Data...)
							}
						}
					case *object.Buffer:
						data = append(data, e.Data...)
					case *object.Array:
						for _, b := range e.Elements {
							if num, ok := b.(*object.Number); ok {
								data = append(data, byte(int(num.Value)%256))
							}
						}
					}
				}
			}
		}

		// Process fileName argument
		fileName := ""
		if len(args) > 1 {
			if str, ok := args[1].(*object.String); ok {
				fileName = str.Value
			}
		}

		// Process options argument
		contentType := ""
		lastModified := time.Now().UnixMilli()

		if len(args) > 2 {
			if opts, ok := args[2].(*object.ObjectMap); ok {
				if typeVal, exists := opts.Get("type"); exists {
					if str, ok := typeVal.(*object.String); ok {
						contentType = str.Value
					}
				}
				if lmVal, exists := opts.Get("lastModified"); exists {
					if num, ok := lmVal.(*object.Number); ok {
						lastModified = int64(num.Value)
					}
				}
			}
		}

		// Store data
		instance.Set("_data", &object.Buffer{Data: data})
		instance.Set("_type", &object.String{Value: contentType})
		instance.Set("_name", &object.String{Value: fileName})
		instance.Set("_lastModified", &object.Number{Value: float64(lastModified)})

		return UNDEFINED
	}

	// File.prototype.name (getter)
	FileClass.NativeMethods["name"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return &object.String{Value: ""}
		}
		if nameVal, exists := instance.Get("_name"); exists {
			if str, ok := nameVal.(*object.String); ok {
				return str
			}
		}
		return &object.String{Value: ""}
	}

	// File.prototype.lastModified (getter)
	FileClass.NativeMethods["lastModified"] = func(this object.Object, args ...object.Object) object.Object {
		instance, ok := this.(*object.ObjectMap)
		if !ok {
			return &object.Number{Value: 0}
		}
		if lmVal, exists := instance.Get("_lastModified"); exists {
			if num, ok := lmVal.(*object.Number); ok {
				return num
			}
		}
		return &object.Number{Value: float64(time.Now().UnixMilli())}
	}

	// File.prototype.webkitRelativePath (getter) - always empty in Node.js
	FileClass.NativeMethods["webkitRelativePath"] = func(this object.Object, args ...object.Object) object.Object {
		return &object.String{Value: ""}
	}
}

// createBlobInstance creates a new Blob instance with the given data and type
func createBlobInstance(data []byte, contentType string) *object.ObjectMap {
	blob := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	blob.Set("_data", &object.Buffer{Data: data})
	blob.Set("_type", &object.String{Value: contentType})

	// Add methods
	for name, method := range BlobClass.NativeMethods {
		if name != "constructor" {
			blob.Set(name, bindNativeMethod(method, blob))
		}
	}

	// Add getters as properties
	if sizeMethod := BlobClass.NativeMethods["size"]; sizeMethod != nil {
		blob.Set("size", sizeMethod(blob))
	}
	if typeMethod := BlobClass.NativeMethods["type"]; typeMethod != nil {
		blob.Set("type", typeMethod(blob))
	}

	return blob
}

// createFileInstance creates a new File instance
func createFileInstance(data []byte, fileName, contentType string, lastModified int64) *object.ObjectMap {
	file := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	file.Set("_data", &object.Buffer{Data: data})
	file.Set("_type", &object.String{Value: contentType})
	file.Set("_name", &object.String{Value: fileName})
	file.Set("_lastModified", &object.Number{Value: float64(lastModified)})

	// Add Blob methods
	for name, method := range BlobClass.NativeMethods {
		if name != "constructor" {
			file.Set(name, bindNativeMethod(method, file))
		}
	}

	// Add File methods
	for name, method := range FileClass.NativeMethods {
		if name != "constructor" {
			file.Set(name, bindNativeMethod(method, file))
		}
	}

	// Add getters as properties
	file.Set("size", &object.Number{Value: float64(len(data))})
	file.Set("type", &object.String{Value: contentType})
	file.Set("name", &object.String{Value: fileName})
	file.Set("lastModified", &object.Number{Value: float64(lastModified)})
	file.Set("webkitRelativePath", &object.String{Value: ""})

	return file
}

// blobCreateResolvedPromise creates a resolved Promise with the given value (Blob-specific to avoid redeclaration)
func blobCreateResolvedPromise(value object.Object) object.Object {
	promise := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	promise.Set("_state", &object.String{Value: "fulfilled"})
	promise.Set("_value", value)
	promise.Set("then", &object.Builtin{
		Name: "then",
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				if fn, ok := args[0].(*object.Function); ok {
					return applyFunction(fn, []object.Object{value})
				}
				if builtin, ok := args[0].(*object.Builtin); ok {
					return builtin.Fn(value)
				}
			}
			return promise
		},
	})
	promise.Set("catch", &object.Builtin{
		Name: "catch",
		Fn: func(args ...object.Object) object.Object {
			return promise
		},
	})
	return promise
}

// blobCreateRejectedPromise creates a rejected Promise with the given error (Blob-specific to avoid redeclaration)
func blobCreateRejectedPromise(err object.Object) object.Object {
	promise := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	promise.Set("_state", &object.String{Value: "rejected"})
	promise.Set("_reason", err)
	promise.Set("then", &object.Builtin{
		Name: "then",
		Fn: func(args ...object.Object) object.Object {
			return promise
		},
	})
	promise.Set("catch", &object.Builtin{
		Name: "catch",
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				if fn, ok := args[0].(*object.Function); ok {
					return applyFunction(fn, []object.Object{err})
				}
				if builtin, ok := args[0].(*object.Builtin); ok {
					return builtin.Fn(err)
				}
			}
			return promise
		},
	})
	return promise
}

// initBlobConstructor initializes the Blob constructor for builtins
func initBlobConstructor() object.Object {
	initBlobFileClasses()

	blobConstructor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Make it callable with new
	blobConstructor.Set("__call__", &object.Builtin{
		Name: "Blob",
		Fn: func(args ...object.Object) object.Object {
			// Process parts argument
			var data []byte
			if len(args) > 0 {
				if arr, ok := args[0].(*object.Array); ok {
					for _, elem := range arr.Elements {
						switch e := elem.(type) {
						case *object.String:
							data = append(data, []byte(e.Value)...)
						case *object.ObjectMap:
							if blobData, exists := e.Get("_data"); exists {
								if buffer, ok := blobData.(*object.Buffer); ok {
									data = append(data, buffer.Data...)
								}
							}
						case *object.Buffer:
							data = append(data, e.Data...)
						}
					}
				}
			}

			// Process options argument
			contentType := ""
			if len(args) > 1 {
				if opts, ok := args[1].(*object.ObjectMap); ok {
					if typeVal, exists := opts.Get("type"); exists {
						if str, ok := typeVal.(*object.String); ok {
							contentType = str.Value
						}
					}
				}
			}

			return createBlobInstance(data, contentType)
		},
	})

	return blobConstructor
}

// initFileConstructor initializes the File constructor for builtins
func initFileConstructor() object.Object {
	initBlobFileClasses()

	fileConstructor := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Make it callable with new
	fileConstructor.Set("__call__", &object.Builtin{
		Name: "File",
		Fn: func(args ...object.Object) object.Object {
			// Process fileBits argument
			var data []byte
			if len(args) > 0 {
				if arr, ok := args[0].(*object.Array); ok {
					for _, elem := range arr.Elements {
						switch e := elem.(type) {
						case *object.String:
							data = append(data, []byte(e.Value)...)
						case *object.ObjectMap:
							if blobData, exists := e.Get("_data"); exists {
								if buffer, ok := blobData.(*object.Buffer); ok {
									data = append(data, buffer.Data...)
								}
							}
						case *object.Buffer:
							data = append(data, e.Data...)
						}
					}
				}
			}

			// Process fileName argument
			fileName := ""
			if len(args) > 1 {
				if str, ok := args[1].(*object.String); ok {
					fileName = str.Value
				}
			}

			// Process options argument
			contentType := ""
			lastModified := time.Now().UnixMilli()
			if len(args) > 2 {
				if opts, ok := args[2].(*object.ObjectMap); ok {
					if typeVal, exists := opts.Get("type"); exists {
						if str, ok := typeVal.(*object.String); ok {
							contentType = str.Value
						}
					}
					if lmVal, exists := opts.Get("lastModified"); exists {
						if num, ok := lmVal.(*object.Number); ok {
							lastModified = int64(num.Value)
						}
					}
				}
			}

			return createFileInstance(data, fileName, contentType, lastModified)
		},
	})

	return fileConstructor
}

// blobToBase64 converts a Blob's data to base64
func blobToBase64(blob *object.ObjectMap) string {
	if data, exists := blob.Get("_data"); exists {
		if buffer, ok := data.(*object.Buffer); ok {
			return base64.StdEncoding.EncodeToString(buffer.Data)
		}
	}
	return ""
}
