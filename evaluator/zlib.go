package evaluator

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/nulang/nulang/object"
)

// initZlibModule initializes the zlib module
func initZlibModule() *object.ObjectMap {
	zlibMod := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// zlib.gzip(buffer, callback)
	zlibMod.Set("gzip", &object.Builtin{Name: "gzip", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("gzip requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("gzip requires a buffer or string")
		}

		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return newError("gzip failed: %s", err.Error())
		}
		writer.Close()

		result := &object.Buffer{Data: buf.Bytes()}

		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{NULL, result})
				Eval(callback.Body, fnEnv)
			}
		}

		return result
	}})

	// zlib.gzipSync(buffer)
	zlibMod.Set("gzipSync", &object.Builtin{Name: "gzipSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("gzipSync requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("gzipSync requires a buffer or string")
		}

		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return newError("gzipSync failed: %s", err.Error())
		}
		writer.Close()

		return &object.Buffer{Data: buf.Bytes()}
	}})

	// zlib.gunzip(buffer, callback)
	zlibMod.Set("gunzip", &object.Builtin{Name: "gunzip", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("gunzip requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("gunzip requires a buffer or string")
		}

		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return newError("gunzip failed: %s", err.Error())
		}
		defer reader.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, reader)
		if err != nil {
			return newError("gunzip failed: %s", err.Error())
		}

		result := &object.Buffer{Data: buf.Bytes()}

		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{NULL, result})
				Eval(callback.Body, fnEnv)
			}
		}

		return result
	}})

	// zlib.gunzipSync(buffer)
	zlibMod.Set("gunzipSync", &object.Builtin{Name: "gunzipSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("gunzipSync requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("gunzipSync requires a buffer or string")
		}

		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return newError("gunzipSync failed: %s", err.Error())
		}
		defer reader.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, reader)
		if err != nil {
			return newError("gunzipSync failed: %s", err.Error())
		}

		return &object.Buffer{Data: buf.Bytes()}
	}})

	// zlib.deflate(buffer, callback)
	zlibMod.Set("deflate", &object.Builtin{Name: "deflate", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("deflate requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("deflate requires a buffer or string")
		}

		var buf bytes.Buffer
		writer := zlib.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return newError("deflate failed: %s", err.Error())
		}
		writer.Close()

		result := &object.Buffer{Data: buf.Bytes()}

		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{NULL, result})
				Eval(callback.Body, fnEnv)
			}
		}

		return result
	}})

	// zlib.deflateSync(buffer)
	zlibMod.Set("deflateSync", &object.Builtin{Name: "deflateSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("deflateSync requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("deflateSync requires a buffer or string")
		}

		var buf bytes.Buffer
		writer := zlib.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return newError("deflateSync failed: %s", err.Error())
		}
		writer.Close()

		return &object.Buffer{Data: buf.Bytes()}
	}})

	// zlib.inflate(buffer, callback)
	zlibMod.Set("inflate", &object.Builtin{Name: "inflate", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("inflate requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("inflate requires a buffer or string")
		}

		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return newError("inflate failed: %s", err.Error())
		}
		defer reader.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, reader)
		if err != nil {
			return newError("inflate failed: %s", err.Error())
		}

		result := &object.Buffer{Data: buf.Bytes()}

		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{NULL, result})
				Eval(callback.Body, fnEnv)
			}
		}

		return result
	}})

	// zlib.inflateSync(buffer)
	zlibMod.Set("inflateSync", &object.Builtin{Name: "inflateSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("inflateSync requires a buffer argument")
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.Buffer:
			data = d.Data
		case *object.String:
			data = []byte(d.Value)
		default:
			return newError("inflateSync requires a buffer or string")
		}

		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return newError("inflateSync failed: %s", err.Error())
		}
		defer reader.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, reader)
		if err != nil {
			return newError("inflateSync failed: %s", err.Error())
		}

		return &object.Buffer{Data: buf.Bytes()}
	}})

	return zlibMod
}
