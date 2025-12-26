package evaluator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nulang/nulang/object"
)

var (
	httpServerClass          *Class
	httpIncomingMessageClass *Class
	httpServerResponseClass  *Class
	httpClientRequestClass   *Class
	httpAgentClass           *Class
)

// initHttpModule creates the http module
func initHttpModule() *object.ObjectMap {
	if eventEmitterClass == nil {
		initEventsModule()
	}
	initStreamModule()
	createHttpClasses()

	httpModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	httpModule.Set("createServer", &object.Builtin{Name: "createServer", Fn: httpCreateServer})
	httpModule.Set("request", &object.Builtin{Name: "request", Fn: httpRequestFn})
	httpModule.Set("get", &object.Builtin{Name: "get", Fn: httpGetFn})

	httpModule.Set("Server", httpServerClass)
	httpModule.Set("IncomingMessage", httpIncomingMessageClass)
	httpModule.Set("ServerResponse", httpServerResponseClass)
	httpModule.Set("ClientRequest", httpClientRequestClass)
	httpModule.Set("Agent", httpAgentClass)

	methods := &object.Array{Elements: []object.Object{}}
	for _, m := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"} {
		methods.Elements = append(methods.Elements, &object.String{Value: m})
	}
	httpModule.Set("METHODS", methods)

	statusCodes := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	for i := 100; i < 600; i++ {
		txt := http.StatusText(i)
		if txt != "" {
			statusCodes.Set(fmt.Sprintf("%d", i), &object.String{Value: txt})
		}
	}
	httpModule.Set("STATUS_CODES", statusCodes)

	httpModule.Set("globalAgent", createClassInstance(httpAgentClass, []object.Object{}, httpAgentClass.Env))

	return httpModule
}

func createHttpClasses() {
	if httpServerClass != nil {
		return
	}

	bind := func(c *Class, name string, fn NativeMethod) {
		c.NativeMethods[name] = fn
	}

	// --- Server ---
	httpServerClass = &Class{
		Name:          "Server",
		SuperClass:    eventEmitterClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	bind(httpServerClass, "constructor", func(this object.Object, args ...object.Object) object.Object {
		if superCtor, ok := eventEmitterClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		if len(args) > 0 {
			if listener, ok := args[0].(*object.Function); ok {
				nativeCall(this.(*object.ObjectMap), "on", &object.String{Value: "request"}, listener)
			}
		}
		this.(*object.ObjectMap).Set("_listening", FALSE)
		return UNDEFINED
	})
	bind(httpServerClass, "listen", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if len(args) < 1 {
			return newError("Server.listen requires port")
		}
		var port int64
		if pNum, ok := args[0].(*object.Number); ok {
			port = int64(pNum.Value)
		} else {
			return newError("Port must be a number")
		}
		host := ""
		var callback object.Object
		for i := 1; i < len(args); i++ {
			if str, ok := args[i].(*object.String); ok {
				host = str.Value
			} else if fn, ok := args[i].(*object.Function); ok {
				callback = fn
			} else if builtin, ok := args[i].(*object.Builtin); ok {
				callback = builtin
			}
		}
		addr := fmt.Sprintf("%s:%d", host, port)
		server := &http.Server{
			Addr: addr,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqObj := createIncomingMessage(r)
				resObj, doneChan := createServerResponse(w)
				nativeCall(instance, "emit", &object.String{Value: "request"}, reqObj, resObj)
				<-doneChan
			}),
		}

		
		// Capture server for closing
		closeFn := &object.Builtin{Fn: func(cArgs ...object.Object) object.Object {
			// Close the server (this will cause ListenAndServe to return)
			server.Close()
			return UNDEFINED
		}}
		instance.Set("_closeGoServer", closeFn)

		instance.Set("_listening", TRUE)
		RegisterAsyncTask()
		go func() {
			defer UnregisterAsyncTask()
			if callback != nil {
				applyHandler(callback, []object.Object{})
			}
			server.ListenAndServe()
		}()
		return instance
	})
	bind(httpServerClass, "close", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		instance.Set("_listening", FALSE)
		
		// Stop Go Server
		if closeFn, ok := instance.Get("_closeGoServer"); ok {
			if c, ok := closeFn.(*object.Builtin); ok {
				c.Fn()
			}
		}

		if len(args) > 0 {
			applyHandler(args[0], []object.Object{})
		}
		nativeCall(instance, "emit", &object.String{Value: "close"})
		return instance
	})

	// --- IncomingMessage ---
	httpIncomingMessageClass = &Class{
		Name:          "IncomingMessage",
		SuperClass:    readableClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	bind(httpIncomingMessageClass, "constructor", func(this object.Object, args ...object.Object) object.Object {
		if superCtor, ok := readableClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		return UNDEFINED
	})
	bind(httpIncomingMessageClass, "on", func(this object.Object, args ...object.Object) object.Object {
		if onFn, ok := eventEmitterClass.NativeMethods["on"]; ok {
			onFn(this, args...)
		}
		if len(args) > 0 {
			if str, ok := args[0].(*object.String); ok && str.Value == "data" {
				instance := this.(*object.ObjectMap)
				if _, ok := instance.Get("_reading"); !ok {
					instance.Set("_reading", TRUE)
					nativeCall(instance, "read")
				}
			}
		}
		return this
	})
	bind(httpIncomingMessageClass, "_read", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if bodyReaderObj, ok := instance.Get("_bodyReader"); ok {
			if readFn, ok := bodyReaderObj.(*object.Builtin); ok {
				readFn.Fn()
			}
		} else {
			nativeCall(instance, "push", NULL)
		}
		return UNDEFINED
	})

	// --- ServerResponse ---
	httpServerResponseClass = &Class{
		Name:          "ServerResponse",
		SuperClass:    writableClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	bind(httpServerResponseClass, "constructor", func(this object.Object, args ...object.Object) object.Object {
		if superCtor, ok := writableClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		instance := this.(*object.ObjectMap)
		instance.Set("statusCode", &object.Number{Value: 200})
		instance.Set("statusMessage", &object.String{Value: "OK"})
		instance.Set("headersSent", FALSE)
		instance.Set("_headers", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
		return UNDEFINED
	})
	bind(httpServerResponseClass, "setHeader", func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 2 { return UNDEFINED }
		instance := this.(*object.ObjectMap)
		name := strings.ToLower(objectToString(args[0]))
		value := objectToString(args[1])
		if h, ok := instance.Get("_headers"); ok {
			h.(*object.ObjectMap).Set(name, &object.String{Value: value})
		}
		return UNDEFINED
	})
	bind(httpServerResponseClass, "getHeader", func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 { return UNDEFINED }
		instance := this.(*object.ObjectMap)
		name := strings.ToLower(objectToString(args[0]))
		if h, ok := instance.Get("_headers"); ok {
			if val, ok := h.(*object.ObjectMap).Get(name); ok {
				return val
			}
		}
		return UNDEFINED
	})
	bind(httpServerResponseClass, "writeHead", func(this object.Object, args ...object.Object) object.Object {
		if len(args) < 1 { return this }
		instance := this.(*object.ObjectMap)
		if code, ok := args[0].(*object.Number); ok {
			instance.Set("statusCode", code)
		}
		return this
	})
	bind(httpServerResponseClass, "_write", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		chunk := args[0]
		if writeFn, ok := instance.Get("_writeToResponse"); ok {
			if wFn, ok := writeFn.(*object.Builtin); ok {
				wFn.Fn(chunk)
			}
		}
		if len(args) > 2 {
			if cb, ok := args[2].(*object.Builtin); ok { cb.Fn() }
		}
		return UNDEFINED
	})
	bind(httpServerResponseClass, "end", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if superEnd, ok := writableClass.NativeMethods["end"]; ok {
			superEnd(this, args...)
		}
		if doneFn, ok := instance.Get("_doneFn"); ok {
			if fn, ok := doneFn.(*object.Builtin); ok {
				fn.Fn()
			}
		}
		return instance
	})

	// --- ClientRequest ---
	httpClientRequestClass = &Class{
		Name:          "ClientRequest",
		SuperClass:    writableClass,
		Methods:       make(map[string]*object.Function),
		NativeMethods: make(map[string]NativeMethod),
		Properties:    make(map[string]object.Object),
		Static:        make(map[string]object.Object),
		Env:           object.NewEnvironment(),
	}
	bind(httpClientRequestClass, "constructor", func(this object.Object, args ...object.Object) object.Object {
		if superCtor, ok := writableClass.NativeMethods["constructor"]; ok {
			superCtor(this, args...)
		}
		return UNDEFINED
	})
	bind(httpClientRequestClass, "end", func(this object.Object, args ...object.Object) object.Object {
		instance := this.(*object.ObjectMap)
		if superEnd, ok := writableClass.NativeMethods["end"]; ok {
			superEnd(this, args...)
		}
		if exec, ok := instance.Get("_execute"); ok {
			if fn, ok := exec.(*object.Builtin); ok {
				fn.Fn()
			}
		}
		return instance
	})

	// --- Agent ---
	httpAgentClass = &Class{
		Name: "Agent",
	}
}

// ServerCreatedCallback is called when a new server is created (for watch mode tracking)
var ServerCreatedCallback func(*object.ObjectMap)

// WrapHttpCreateServer sets a callback function that gets called when HTTP servers are created
// This is used by watch mode to track and cleanup servers on restart
func WrapHttpCreateServer(callback func(*object.ObjectMap)) func(args ...object.Object) object.Object {
	ServerCreatedCallback = callback
	return httpCreateServer
}

func httpCreateServer(args ...object.Object) object.Object {
	instance := createClassInstance(httpServerClass, args, httpServerClass.Env)
	
	// Notify watch mode if callback is set
	if ServerCreatedCallback != nil {
		if objMap, ok := instance.(*object.ObjectMap); ok {
			ServerCreatedCallback(objMap)
		}
	}
	
	return instance
}

func createIncomingMessage(r *http.Request) *object.ObjectMap {
	msgObj := createClassInstance(httpIncomingMessageClass, []object.Object{}, httpIncomingMessageClass.Env)
	msg := msgObj.(*object.ObjectMap)
	msg.Set("method", &object.String{Value: r.Method})
	msg.Set("url", &object.String{Value: r.URL.String()})
	msg.Set("httpVersion", &object.String{Value: fmt.Sprintf("%d.%d", r.ProtoMajor, r.ProtoMinor)})
	headers := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	rawHeaders := &object.Array{Elements: []object.Object{}}
	for k, v := range r.Header {
		joined := strings.Join(v, ", ")
		headers.Set(strings.ToLower(k), &object.String{Value: joined})
		rawHeaders.Elements = append(rawHeaders.Elements, &object.String{Value: k}, &object.String{Value: joined})
	}
	msg.Set("headers", headers)
	msg.Set("rawHeaders", rawHeaders)
	readerFn := &object.Builtin{Fn: func(args ...object.Object) object.Object {
		go func() {
			buf := make([]byte, 4096)
			for {
				n, err := r.Body.Read(buf)
				if n > 0 {
					data := make([]byte, n)
					copy(data, buf[:n])
					nativeCall(msg, "push", &object.Buffer{Data: data})
				}
				if err != nil {
					nativeCall(msg, "push", NULL)
					return
				}
			}
		}()
		return UNDEFINED
	}}
	msg.Set("_bodyReader", readerFn)
	return msg
}

func createServerResponse(w http.ResponseWriter) (*object.ObjectMap, chan bool) {
	done := make(chan bool)
	resObj := createClassInstance(httpServerResponseClass, []object.Object{}, httpServerResponseClass.Env)
	res := resObj.(*object.ObjectMap)
	var headersWritten = false
	var headersMutex sync.Mutex
	writeFn := &object.Builtin{Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 { return UNDEFINED }
		headersMutex.Lock()
		if !headersWritten {
			code := 200
			if c, ok := res.Get("statusCode"); ok {
				if num, ok := c.(*object.Number); ok {
					code = int(num.Value)
				}
			}
			if h, ok := res.Get("_headers"); ok {
				if hMap, ok := h.(*object.ObjectMap); ok {
					for k, pair := range hMap.Pairs {
						w.Header().Set(k, objectToString(pair.Value))
					}
				}
			}
			w.WriteHeader(code)
			headersWritten = true
			res.Set("headersSent", TRUE)
		}
		headersMutex.Unlock()
		var data []byte
		switch v := args[0].(type) {
		case *object.String:
			data = []byte(v.Value)
		case *object.Buffer:
			data = v.Data
		}
		w.Write(data)
		return UNDEFINED
	}}
	res.Set("_writeToResponse", writeFn)
	doneFn := &object.Builtin{Fn: func(args ...object.Object) object.Object {
		headersMutex.Lock()
		if !headersWritten {
			code := 200
			if c, ok := res.Get("statusCode"); ok {
				if num, ok := c.(*object.Number); ok {
					code = int(num.Value)
				}
			}
			if h, ok := res.Get("_headers"); ok {
				if hMap, ok := h.(*object.ObjectMap); ok {
					for k, pair := range hMap.Pairs {
						w.Header().Set(k, objectToString(pair.Value))
					}
				}
			}
			w.WriteHeader(code)
			headersWritten = true
			res.Set("headersSent", TRUE)
		}
		headersMutex.Unlock()
		// Avoid panic if closed twice
		select {
		case <-done:
		default:
			close(done)
		}
		return UNDEFINED
	}}
	res.Set("_doneFn", doneFn)
	return res, done
}

func httpRequestFn(args ...object.Object) object.Object {
	if len(args) < 1 { return newError("request expects options") }
	var method = "GET"
	var urlStr = ""
	var headers = make(map[string]string)
	var options *object.ObjectMap
	if optObj, ok := args[0].(*object.ObjectMap); ok {
		options = optObj
	} else if optStr, ok := args[0].(*object.String); ok {
		urlStr = optStr.Value
		if len(args) > 1 {
			if o, ok := args[1].(*object.ObjectMap); ok {
				options = o
			}
		}
	}
	if options != nil {
		if m, ok := options.Get("method"); ok { method = objectToString(m) }
		if u, ok := options.Get("url"); ok { urlStr = objectToString(u) }
		if h, ok := options.Get("headers"); ok {
			if hMap, ok := h.(*object.ObjectMap); ok {
				for k, p := range hMap.Pairs {
					headers[k] = objectToString(p.Value)
				}
			}
		}
	}
	reqObjRes := createClassInstance(httpClientRequestClass, []object.Object{}, httpClientRequestClass.Env)
	reqObj := reqObjRes.(*object.ObjectMap)
	var bodyBuf bytes.Buffer
	_writeFn := &object.Builtin{Fn: func(wArgs ...object.Object) object.Object {
		chunk := wArgs[0]
		if str, ok := chunk.(*object.String); ok {
			bodyBuf.WriteString(str.Value)
		} else if buf, ok := chunk.(*object.Buffer); ok {
			bodyBuf.Write(buf.Data)
		}
		if len(wArgs) > 2 {
			if cb, ok := wArgs[2].(*object.Builtin); ok { cb.Fn() }
		}
		return UNDEFINED
	}}
	reqObj.Set("_write", _writeFn)
	_executeFn := &object.Builtin{Fn: func(eArgs ...object.Object) object.Object {
		go func() {
			client := &http.Client{}
			r, err := http.NewRequest(method, urlStr, &bodyBuf)
			if err != nil {
				nativeCall(reqObj, "emit", &object.String{Value: "error"}, &object.Error{Message: err.Error()})
				return
			}
			for k, v := range headers {
				r.Header.Set(k, v)
			}
			resp, err := client.Do(r)
			if err != nil {
				nativeCall(reqObj, "emit", &object.String{Value: "error"}, &object.Error{Message: err.Error()})
				return
			}
			respObj := createIncomingMessage(&http.Request{Method: method, URL: r.URL, ProtoMajor: resp.ProtoMajor, ProtoMinor: resp.ProtoMinor, Header: resp.Header, Body: resp.Body})
			respObj.Set("statusCode", &object.Number{Value: float64(resp.StatusCode)})
			respObj.Set("statusMessage", &object.String{Value: resp.Status})
			nativeCall(reqObj, "emit", &object.String{Value: "response"}, respObj)
			if (len(args) > 1) { // Call callback if provided initially
				if cb, ok := args[len(args)-1].(*object.Function); ok {
					applyHandler(cb, []object.Object{respObj})
				}
			}
		}()
		return UNDEFINED
	}}
	reqObj.Set("_execute", _executeFn)
	return reqObj
}

func httpGetFn(args ...object.Object) object.Object {
	req := httpRequestFn(args...)
	nativeCall(req.(*object.ObjectMap), "end")
	return req
}

func initFetchFunction() *object.Builtin {
	return &object.Builtin{Name: "fetch", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 { return newError("fetch requires URL") }
		url := objectToString(args[0])
		method := "GET"
		var body io.Reader
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				if m, ok := opts.Get("method"); ok { method = objectToString(m) }
				if b, ok := opts.Get("body"); ok {
					body = strings.NewReader(objectToString(b))
				}
			}
		}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return createRejectedPromise(newError("%s", err.Error()))
		}
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return createRejectedPromise(newError("%s", err.Error()))
		}
		respObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		respObj.Set("status", &object.Number{Value: float64(resp.StatusCode)})
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		respObj.Set("text", &object.Builtin{Fn: func(a ...object.Object) object.Object {
			return createResolvedPromise(&object.String{Value: string(b)})
		}})
		respObj.Set("json", &object.Builtin{Fn: func(a ...object.Object) object.Object {
			var d interface{}
			json.Unmarshal(b, &d)
			return createResolvedPromise(goValueToObject(d))
		}})
		return createResolvedPromise(respObj)
	}}
}

func goValueToObject(v interface{}) object.Object {
	switch val := v.(type) {
	case nil:
		return NULL
	case string:
		return &object.String{Value: val}
	case float64:
		return &object.Number{Value: val}
	case bool:
		if val {
			return TRUE
		}
		return FALSE
	case []interface{}:
		arr := &object.Array{Elements: make([]object.Object, len(val))}
		for i, elem := range val {
			arr.Elements[i] = goValueToObject(elem)
		}
		return arr
	case map[string]interface{}:
		o := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		for k, v := range val {
			o.Set(k, goValueToObject(v))
		}
		return o
	}
	return NULL
}

func createResolvedPromise(val object.Object) *object.Promise {
	return &object.Promise{State: "fulfilled", Value: val}
}
func createRejectedPromise(val object.Object) *object.Promise {
	return &object.Promise{State: "rejected", Reason: val}
}

func initURLModule() *object.ObjectMap {
	m := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	m.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 { return newError("URL.parse requires url") }
		urlStr := objectToString(args[0])
		return parseURL(urlStr)
	}})
	return m
}

func parseURL(urlStr string) object.Object {
	urlObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	urlObj.Set("href", &object.String{Value: urlStr})
	protocolEnd := strings.Index(urlStr, "://")
	if protocolEnd > 0 {
		urlObj.Set("protocol", &object.String{Value: urlStr[:protocolEnd+1]})
		urlStr = urlStr[protocolEnd+3:]
	}
	hashIdx := strings.Index(urlStr, "#")
	hash := ""
	if hashIdx >= 0 {
		hash = urlStr[hashIdx:]
		urlStr = urlStr[:hashIdx]
	}
	urlObj.Set("hash", &object.String{Value: hash})
	queryIdx := strings.Index(urlStr, "?")
	search := ""
	if queryIdx >= 0 {
		search = urlStr[queryIdx:]
		urlStr = urlStr[:queryIdx]
	}
	urlObj.Set("search", &object.String{Value: search})
	pathIdx := strings.Index(urlStr, "/")
	host := urlStr
	pathname := "/"
	if pathIdx >= 0 {
		host = urlStr[:pathIdx]
		pathname = urlStr[pathIdx:]
	}
	urlObj.Set("host", &object.String{Value: host})
	urlObj.Set("hostname", &object.String{Value: host})
	urlObj.Set("pathname", &object.String{Value: pathname})
	portIdx := strings.Index(host, ":")
	port := ""
	if portIdx >= 0 {
		port = host[portIdx+1:]
		urlObj.Set("hostname", &object.String{Value: host[:portIdx]})
	}
	urlObj.Set("port", &object.String{Value: port})
	return urlObj
}

func initQueryStringModule() *object.ObjectMap {
	qs := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	qs.Set("stringify", &object.Builtin{Name: "stringify", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 { return &object.String{Value: ""} }
		obj, ok := args[0].(*object.ObjectMap)
		if !ok { return &object.String{Value: ""} }
		var parts []string
		for key, pair := range obj.Pairs {
			parts = append(parts, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
		}
		return &object.String{Value: strings.Join(parts, "&")}
	}})
	qs.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 { return &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)} }
		str := objectToString(args[0])
		if strings.HasPrefix(str, "?") { str = str[1:] }
		result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		for _, pair := range strings.Split(str, "&") {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				result.Set(parts[0], &object.String{Value: parts[1]})
			}
		}
		return result
	}})
	return qs
}
