package evaluator

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/nulang/nulang/object"
)

func initBufferConstructor() *object.ObjectMap {
	bufferObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Buffer.from(data, encoding?)
	bufferObj.Set("from", &object.Builtin{Name: "from", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Buffer{Data: []byte{}}
		}

		encoding := "utf8"
		if len(args) > 1 {
			if enc, ok := args[1].(*object.String); ok {
				encoding = enc.Value
			}
		}

		switch data := args[0].(type) {
		case *object.String:
			switch encoding {
			case "base64":
				decoded, err := base64.StdEncoding.DecodeString(data.Value)
				if err != nil {
					return newError("Invalid base64 string")
				}
				return &object.Buffer{Data: decoded}
			case "hex":
				decoded, err := hex.DecodeString(data.Value)
				if err != nil {
					return newError("Invalid hex string")
				}
				return &object.Buffer{Data: decoded}
			default:
				return &object.Buffer{Data: []byte(data.Value)}
			}
		case *object.Array:
			bytes := make([]byte, len(data.Elements))
			for i, elem := range data.Elements {
				if num, ok := elem.(*object.Number); ok {
					bytes[i] = byte(int(num.Value) & 0xFF)
				}
			}
			return &object.Buffer{Data: bytes}
		case *object.Buffer:
			newData := make([]byte, len(data.Data))
			copy(newData, data.Data)
			return &object.Buffer{Data: newData}
		}

		return &object.Buffer{Data: []byte{}}
	}})

	// Buffer.alloc(size, fill?, encoding?)
	bufferObj.Set("alloc", &object.Builtin{Name: "alloc", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Buffer.alloc requires size argument")
		}
		size := int(args[0].(*object.Number).Value)
		data := make([]byte, size)

		if len(args) > 1 {
			var fillByte byte = 0
			switch fill := args[1].(type) {
			case *object.Number:
				fillByte = byte(int(fill.Value) & 0xFF)
			case *object.String:
				if len(fill.Value) > 0 {
					fillByte = fill.Value[0]
				}
			}
			for i := range data {
				data[i] = fillByte
			}
		}

		return &object.Buffer{Data: data}
	}})

	// Buffer.concat(list, totalLength?)
	bufferObj.Set("concat", &object.Builtin{Name: "concat", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Buffer{Data: []byte{}}
		}

		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("Buffer.concat requires an array of Buffers")
		}

		var result []byte
		for _, elem := range arr.Elements {
			if buf, ok := elem.(*object.Buffer); ok {
				result = append(result, buf.Data...)
			}
		}

		return &object.Buffer{Data: result}
	}})

	// Buffer.isBuffer(obj)
	bufferObj.Set("isBuffer", &object.Builtin{Name: "isBuffer", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		_, ok := args[0].(*object.Buffer)
		return nativeBoolToBooleanObject(ok)
	}})

	// Buffer.byteLength(string, encoding?)
	bufferObj.Set("byteLength", &object.Builtin{Name: "byteLength", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: 0}
		}
		switch data := args[0].(type) {
		case *object.String:
			return &object.Number{Value: float64(len(data.Value))}
		case *object.Buffer:
			return &object.Number{Value: float64(len(data.Data))}
		}
		return &object.Number{Value: 0}
	}})

	return bufferObj
}

// evalBufferMethods handles Buffer instance methods
func evalBufferProperty(buf *object.Buffer, prop string) object.Object {
	switch prop {
	case "length":
		return &object.Number{Value: float64(len(buf.Data))}
	case "toString":
		return &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
			encoding := "utf8"
			if len(args) > 0 {
				if enc, ok := args[0].(*object.String); ok {
					encoding = enc.Value
				}
			}
			switch encoding {
			case "hex":
				return &object.String{Value: hex.EncodeToString(buf.Data)}
			case "base64":
				return &object.String{Value: base64.StdEncoding.EncodeToString(buf.Data)}
			default:
				return &object.String{Value: string(buf.Data)}
			}
		}}
	case "toJSON":
		return &object.Builtin{Name: "toJSON", Fn: func(args ...object.Object) object.Object {
			elements := make([]object.Object, len(buf.Data))
			for i, b := range buf.Data {
				elements[i] = &object.Number{Value: float64(b)}
			}
			result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			result.Set("type", &object.String{Value: "Buffer"})
			result.Set("data", &object.Array{Elements: elements})
			return result
		}}
	case "slice":
		return &object.Builtin{Name: "slice", Fn: func(args ...object.Object) object.Object {
			start, end := 0, len(buf.Data)
			if len(args) > 0 {
				start = int(args[0].(*object.Number).Value)
				if start < 0 {
					start = len(buf.Data) + start
				}
			}
			if len(args) > 1 {
				end = int(args[1].(*object.Number).Value)
				if end < 0 {
					end = len(buf.Data) + end
				}
			}
			if start < 0 {
				start = 0
			}
			if end > len(buf.Data) {
				end = len(buf.Data)
			}
			if start >= end {
				return &object.Buffer{Data: []byte{}}
			}
			newData := make([]byte, end-start)
			copy(newData, buf.Data[start:end])
			return &object.Buffer{Data: newData}
		}}
	case "copy":
		return &object.Builtin{Name: "copy", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: 0}
			}
			target, ok := args[0].(*object.Buffer)
			if !ok {
				return newError("First argument must be a Buffer")
			}
			targetStart := 0
			sourceStart := 0
			sourceEnd := len(buf.Data)
			if len(args) > 1 {
				targetStart = int(args[1].(*object.Number).Value)
			}
			if len(args) > 2 {
				sourceStart = int(args[2].(*object.Number).Value)
			}
			if len(args) > 3 {
				sourceEnd = int(args[3].(*object.Number).Value)
			}
			copied := copy(target.Data[targetStart:], buf.Data[sourceStart:sourceEnd])
			return &object.Number{Value: float64(copied)}
		}}
	case "equals":
		return &object.Builtin{Name: "equals", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			other, ok := args[0].(*object.Buffer)
			if !ok {
				return FALSE
			}
			if len(buf.Data) != len(other.Data) {
				return FALSE
			}
			for i, b := range buf.Data {
				if b != other.Data[i] {
					return FALSE
				}
			}
			return TRUE
		}}
	case "fill":
		return &object.Builtin{Name: "fill", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return buf
			}
			var fillByte byte = 0
			switch fill := args[0].(type) {
			case *object.Number:
				fillByte = byte(int(fill.Value) & 0xFF)
			case *object.String:
				if len(fill.Value) > 0 {
					fillByte = fill.Value[0]
				}
			}
			for i := range buf.Data {
				buf.Data[i] = fillByte
			}
			return buf
		}}
	case "indexOf":
		return &object.Builtin{Name: "indexOf", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: -1}
			}
			var search []byte
			switch s := args[0].(type) {
			case *object.Number:
				search = []byte{byte(int(s.Value) & 0xFF)}
			case *object.String:
				search = []byte(s.Value)
			case *object.Buffer:
				search = s.Data
			}
			if len(search) == 0 || len(search) > len(buf.Data) {
				return &object.Number{Value: -1}
			}
			for i := 0; i <= len(buf.Data)-len(search); i++ {
				found := true
				for j := 0; j < len(search); j++ {
					if buf.Data[i+j] != search[j] {
						found = false
						break
					}
				}
				if found {
					return &object.Number{Value: float64(i)}
				}
			}
			return &object.Number{Value: -1}
		}}
	}
	return UNDEFINED
}
