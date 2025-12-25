package evaluator

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/nulang/nulang/object"
)

func initCryptoModule() *object.ObjectMap {
	crypto := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// crypto.createHash(algorithm)
	crypto.Set("createHash", &object.Builtin{Name: "createHash", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("createHash requires algorithm argument")
		}
		algorithm := objectToString(args[0])
		return createHashObject(algorithm)
	}})

	// crypto.createHmac(algorithm, key)
	crypto.Set("createHmac", &object.Builtin{Name: "createHmac", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("createHmac requires algorithm and key arguments")
		}
		algorithm := objectToString(args[0])
		key := objectToString(args[1])
		return createHmacObject(algorithm, key)
	}})

	// crypto.randomBytes(size)
	crypto.Set("randomBytes", &object.Builtin{Name: "randomBytes", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("randomBytes requires size argument")
		}
		size := int(args[0].(*object.Number).Value)
		bytes := make([]byte, size)
		_, err := rand.Read(bytes)
		if err != nil {
			return newError("failed to generate random bytes: %s", err)
		}
		return &object.Buffer{Data: bytes}
	}})

	// crypto.randomUUID()
	crypto.Set("randomUUID", &object.Builtin{Name: "randomUUID", Fn: func(args ...object.Object) object.Object {
		uuid := make([]byte, 16)
		_, err := rand.Read(uuid)
		if err != nil {
			return newError("failed to generate UUID: %s", err)
		}
		// Set version (4) and variant bits
		uuid[6] = (uuid[6] & 0x0f) | 0x40
		uuid[8] = (uuid[8] & 0x3f) | 0x80
		return &object.String{Value: fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])}
	}})

	return crypto
}

func createHashObject(algorithm string) *object.ObjectMap {
	hashObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	var data []byte

	hashObj.Set("update", &object.Builtin{Name: "update", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return hashObj
		}
		input := objectToString(args[0])
		data = append(data, []byte(input)...)
		return hashObj
	}})

	hashObj.Set("digest", &object.Builtin{Name: "digest", Fn: func(args ...object.Object) object.Object {
		encoding := "hex"
		if len(args) > 0 {
			encoding = objectToString(args[0])
		}

		var hash []byte
		switch algorithm {
		case "md5":
			h := md5.Sum(data)
			hash = h[:]
		case "sha1":
			h := sha1.Sum(data)
			hash = h[:]
		case "sha256":
			h := sha256.Sum256(data)
			hash = h[:]
		case "sha512":
			h := sha512.Sum512(data)
			hash = h[:]
		default:
			return newError("unsupported hash algorithm: %s", algorithm)
		}

		switch encoding {
		case "hex":
			return &object.String{Value: hex.EncodeToString(hash)}
		case "base64":
			return &object.String{Value: base64.StdEncoding.EncodeToString(hash)}
		default:
			return &object.Buffer{Data: hash}
		}
	}})

	return hashObj
}

func createHmacObject(algorithm, key string) *object.ObjectMap {
	hmacObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	var data []byte
	keyBytes := []byte(key)

	hmacObj.Set("update", &object.Builtin{Name: "update", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return hmacObj
		}
		input := objectToString(args[0])
		data = append(data, []byte(input)...)
		return hmacObj
	}})

	hmacObj.Set("digest", &object.Builtin{Name: "digest", Fn: func(args ...object.Object) object.Object {
		encoding := "hex"
		if len(args) > 0 {
			encoding = objectToString(args[0])
		}

		var hash []byte
		switch algorithm {
		case "sha256":
			h := hmac.New(sha256.New, keyBytes)
			h.Write(data)
			hash = h.Sum(nil)
		case "sha512":
			h := hmac.New(sha512.New, keyBytes)
			h.Write(data)
			hash = h.Sum(nil)
		case "sha1":
			h := hmac.New(sha1.New, keyBytes)
			h.Write(data)
			hash = h.Sum(nil)
		default:
			return newError("unsupported hmac algorithm: %s", algorithm)
		}

		switch encoding {
		case "hex":
			return &object.String{Value: hex.EncodeToString(hash)}
		case "base64":
			return &object.String{Value: base64.StdEncoding.EncodeToString(hash)}
		default:
			return &object.Buffer{Data: hash}
		}
	}})

	return hmacObj
}
