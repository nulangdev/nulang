package evaluator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nulang/nulang/object"
)

// initHttpModule creates the http module
func initHttpModule() *object.ObjectMap {
	httpModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// http.get(url, options?)
	httpModule.Set("get", &object.Builtin{Name: "get", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.get requires a URL argument")
		}
		url := objectToString(args[0])
		return httpRequest("GET", url, nil, args)
	}})

	// http.post(url, body, options?)
	httpModule.Set("post", &object.Builtin{Name: "post", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.post requires a URL argument")
		}
		url := objectToString(args[0])
		var body io.Reader
		if len(args) > 1 && args[1] != nil {
			bodyStr := objectToString(args[1])
			body = strings.NewReader(bodyStr)
		}
		return httpRequest("POST", url, body, args)
	}})

	// http.put(url, body, options?)
	httpModule.Set("put", &object.Builtin{Name: "put", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.put requires a URL argument")
		}
		url := objectToString(args[0])
		var body io.Reader
		if len(args) > 1 && args[1] != nil {
			bodyStr := objectToString(args[1])
			body = strings.NewReader(bodyStr)
		}
		return httpRequest("PUT", url, body, args)
	}})

	// http.delete(url, options?)
	httpModule.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.delete requires a URL argument")
		}
		url := objectToString(args[0])
		return httpRequest("DELETE", url, nil, args)
	}})

	// http.patch(url, body, options?)
	httpModule.Set("patch", &object.Builtin{Name: "patch", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.patch requires a URL argument")
		}
		url := objectToString(args[0])
		var body io.Reader
		if len(args) > 1 && args[1] != nil {
			bodyStr := objectToString(args[1])
			body = strings.NewReader(bodyStr)
		}
		return httpRequest("PATCH", url, body, args)
	}})

	// http.request(options) - full control
	httpModule.Set("request", &object.Builtin{Name: "request", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("http.request requires an options object")
		}

		opts, ok := args[0].(*object.ObjectMap)
		if !ok {
			return newError("http.request requires an options object")
		}

		method := "GET"
		if m, ok := opts.Get("method"); ok {
			method = strings.ToUpper(objectToString(m))
		}

		urlObj, ok := opts.Get("url")
		if !ok {
			return newError("http.request requires a url in options")
		}
		url := objectToString(urlObj)

		var body io.Reader
		if b, ok := opts.Get("body"); ok {
			bodyStr := objectToString(b)
			body = strings.NewReader(bodyStr)
		}

		return httpRequestWithOptions(method, url, body, opts)
	}})

	return httpModule
}

// initFetchFunction creates the global fetch function
func initFetchFunction() *object.Builtin {
	return &object.Builtin{Name: "fetch", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("fetch requires a URL argument")
		}

		url := objectToString(args[0])
		method := "GET"
		var body io.Reader
		var headers map[string]string

		// Parse options if provided
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				if m, ok := opts.Get("method"); ok {
					method = strings.ToUpper(objectToString(m))
				}
				if b, ok := opts.Get("body"); ok {
					bodyStr := objectToString(b)
					body = strings.NewReader(bodyStr)
				}
				if h, ok := opts.Get("headers"); ok {
					if headersMap, ok := h.(*object.ObjectMap); ok {
						headers = make(map[string]string)
						for key, pair := range headersMap.Pairs {
							headers[key] = objectToString(pair.Value)
						}
					}
				}
			}
		}

		// Create HTTP request
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return createRejectedPromise(newError("Failed to create request: %s", err.Error()))
		}

		// Set headers
		if headers != nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}

		// Default content type for POST/PUT/PATCH
		if (method == "POST" || method == "PUT" || method == "PATCH") && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		// Execute request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return createRejectedPromise(newError("Request failed: %s", err.Error()))
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return createRejectedPromise(newError("Failed to read response: %s", err.Error()))
		}

		// Create Response object
		response := createResponseObject(resp, respBody)
		return createResolvedPromise(response)
	}}
}

func httpRequest(method string, url string, body io.Reader, args []object.Object) object.Object {
	var headers map[string]string

	// Check for options argument
	optIndex := 1
	if method == "POST" || method == "PUT" || method == "PATCH" {
		optIndex = 2
	}
	if len(args) > optIndex {
		if opts, ok := args[optIndex].(*object.ObjectMap); ok {
			if h, ok := opts.Get("headers"); ok {
				if headersMap, ok := h.(*object.ObjectMap); ok {
					headers = make(map[string]string)
					for key, pair := range headersMap.Pairs {
						headers[key] = objectToString(pair.Value)
					}
				}
			}
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return newError("Failed to create request: %s", err.Error())
	}

	// Set headers
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Default content type
	if (method == "POST" || method == "PUT" || method == "PATCH") && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return newError("Request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return newError("Failed to read response: %s", err.Error())
	}

	return createResponseObject(resp, respBody)
}

func httpRequestWithOptions(method string, url string, body io.Reader, opts *object.ObjectMap) object.Object {
	var headers map[string]string
	timeout := 30 * time.Second

	if h, ok := opts.Get("headers"); ok {
		if headersMap, ok := h.(*object.ObjectMap); ok {
			headers = make(map[string]string)
			for key, pair := range headersMap.Pairs {
				headers[key] = objectToString(pair.Value)
			}
		}
	}

	if t, ok := opts.Get("timeout"); ok {
		if num, ok := t.(*object.Number); ok {
			timeout = time.Duration(num.Value) * time.Millisecond
		}
	}

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return newError("Failed to create request: %s", err.Error())
	}

	// Set headers
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Execute
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return newError("Request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return newError("Failed to read response: %s", err.Error())
	}

	return createResponseObject(resp, respBody)
}

func createResponseObject(resp *http.Response, body []byte) *object.ObjectMap {
	response := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Status
	response.Set("status", &object.Number{Value: float64(resp.StatusCode)})
	response.Set("statusText", &object.String{Value: resp.Status})
	response.Set("ok", nativeBoolToBooleanObject(resp.StatusCode >= 200 && resp.StatusCode < 300))

	// Headers
	headers := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	for key, values := range resp.Header {
		headers.Set(strings.ToLower(key), &object.String{Value: strings.Join(values, ", ")})
	}
	response.Set("headers", headers)

	// Body as string
	bodyStr := string(body)
	response.Set("body", &object.String{Value: bodyStr})

	// text() method
	response.Set("text", &object.Builtin{Name: "text", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: bodyStr}
	}})

	// json() method
	response.Set("json", &object.Builtin{Name: "json", Fn: func(args ...object.Object) object.Object {
		var result interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return newError("Failed to parse JSON: %s", err.Error())
		}
		return goValueToObject(result)
	}})

	// arrayBuffer() method - returns Buffer
	response.Set("arrayBuffer", &object.Builtin{Name: "arrayBuffer", Fn: func(args ...object.Object) object.Object {
		return &object.Buffer{Data: body}
	}})

	return response
}

func createResolvedPromise(value object.Object) *object.Promise {
	return &object.Promise{
		State: "fulfilled",
		Value: value,
	}
}

func createRejectedPromise(err object.Object) *object.Promise {
	return &object.Promise{
		State:  "rejected",
		Reason: err,
	}
}

// goValueToObject converts Go values to Nulang objects
func goValueToObject(v interface{}) object.Object {
	switch val := v.(type) {
	case nil:
		return NULL
	case bool:
		return nativeBoolToBooleanObject(val)
	case float64:
		return &object.Number{Value: val}
	case string:
		return &object.String{Value: val}
	case []interface{}:
		elements := make([]object.Object, len(val))
		for i, elem := range val {
			elements[i] = goValueToObject(elem)
		}
		return &object.Array{Elements: elements}
	case map[string]interface{}:
		obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		for key, value := range val {
			obj.Set(key, goValueToObject(value))
		}
		return obj
	default:
		return &object.String{Value: fmt.Sprintf("%v", val)}
	}
}

// URL module
func initURLModule() *object.ObjectMap {
	urlModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// URL constructor-like function
	urlModule.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("URL.parse requires a URL string")
		}

		urlStr := objectToString(args[0])
		return parseURL(urlStr)
	}})

	// URLSearchParams
	urlModule.Set("searchParams", &object.Builtin{Name: "searchParams", Fn: func(args ...object.Object) object.Object {
		params := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

		if len(args) > 0 {
			if str, ok := args[0].(*object.String); ok {
				// Parse query string
				queryStr := str.Value
				if strings.HasPrefix(queryStr, "?") {
					queryStr = queryStr[1:]
				}
				for _, pair := range strings.Split(queryStr, "&") {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						params.Set(parts[0], &object.String{Value: parts[1]})
					} else if len(parts) == 1 {
						params.Set(parts[0], &object.String{Value: ""})
					}
				}
			}
		}

		// get method
		params.Set("get", &object.Builtin{Name: "get", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return NULL
			}
			key := objectToString(args[0])
			if val, ok := params.Get(key); ok {
				return val
			}
			return NULL
		}})

		// toString method
		params.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
			var parts []string
			for key, pair := range params.Pairs {
				if key == "get" || key == "toString" || key == "set" || key == "has" {
					continue
				}
				parts = append(parts, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
			}
			return &object.String{Value: strings.Join(parts, "&")}
		}})

		return params
	}})

	return urlModule
}

func parseURL(urlStr string) object.Object {
	urlObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Simple URL parsing
	urlObj.Set("href", &object.String{Value: urlStr})

	// Protocol
	protocolEnd := strings.Index(urlStr, "://")
	if protocolEnd > 0 {
		urlObj.Set("protocol", &object.String{Value: urlStr[:protocolEnd+1]})
		urlStr = urlStr[protocolEnd+3:]
	}

	// Query and hash
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

	// Host and pathname
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

	// Port
	portIdx := strings.Index(host, ":")
	port := ""
	if portIdx >= 0 {
		port = host[portIdx+1:]
		urlObj.Set("hostname", &object.String{Value: host[:portIdx]})
	}
	urlObj.Set("port", &object.String{Value: port})

	return urlObj
}

// QueryString module
func initQueryStringModule() *object.ObjectMap {
	qs := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// stringify
	qs.Set("stringify", &object.Builtin{Name: "stringify", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: ""}
		}

		obj, ok := args[0].(*object.ObjectMap)
		if !ok {
			return &object.String{Value: ""}
		}

		var parts []string
		for key, pair := range obj.Pairs {
			parts = append(parts, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
		}
		return &object.String{Value: strings.Join(parts, "&")}
	}})

	// parse
	qs.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		}

		str := objectToString(args[0])
		if strings.HasPrefix(str, "?") {
			str = str[1:]
		}

		result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		for _, pair := range strings.Split(str, "&") {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				result.Set(parts[0], &object.String{Value: parts[1]})
			} else if len(parts) == 1 && parts[0] != "" {
				result.Set(parts[0], &object.String{Value: ""})
			}
		}
		return result
	}})

	return qs
}

// FormData like object for multipart requests
func createFormData() *object.ObjectMap {
	formData := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	data := make(map[string]object.Object)

	formData.Set("append", &object.Builtin{Name: "append", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return UNDEFINED
		}
		key := objectToString(args[0])
		data[key] = args[1]
		return UNDEFINED
	}})

	formData.Set("get", &object.Builtin{Name: "get", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return NULL
		}
		key := objectToString(args[0])
		if val, ok := data[key]; ok {
			return val
		}
		return NULL
	}})

	formData.Set("has", &object.Builtin{Name: "has", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := objectToString(args[0])
		_, ok := data[key]
		return nativeBoolToBooleanObject(ok)
	}})

	formData.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		key := objectToString(args[0])
		delete(data, key)
		return UNDEFINED
	}})

	formData.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		var buf bytes.Buffer
		for key, val := range data {
			if buf.Len() > 0 {
				buf.WriteString("&")
			}
			buf.WriteString(fmt.Sprintf("%s=%s", key, objectToString(val)))
		}
		return &object.String{Value: buf.String()}
	}})

	return formData
}
