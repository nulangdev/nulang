package evaluator

import (
	"github.com/nulang/nulang/object"
)

func evalStringProperty(str *object.String, prop string) object.Object {
	switch prop {
	case "length":
		return &object.Number{Value: float64(len(str.Value))}
	case "toUpperCase":
		return &object.Builtin{Name: "toUpperCase", Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: toUpper(str.Value)}
		}}
	case "toLowerCase":
		return &object.Builtin{Name: "toLowerCase", Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: toLower(str.Value)}
		}}
	case "trim":
		return &object.Builtin{Name: "trim", Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: trim(str.Value)}
		}}
	case "split":
		return &object.Builtin{Name: "split", Fn: func(args ...object.Object) object.Object {
			sep := ""
			if len(args) > 0 {
				if s, ok := args[0].(*object.String); ok {
					sep = s.Value
				}
			}
			parts := split(str.Value, sep)
			elements := make([]object.Object, len(parts))
			for i, part := range parts {
				elements[i] = &object.String{Value: part}
			}
			return &object.Array{Elements: elements}
		}}
	case "charAt":
		return &object.Builtin{Name: "charAt", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.String{Value: ""}
			}
			idx := int(args[0].(*object.Number).Value)
			if idx < 0 || idx >= len(str.Value) {
				return &object.String{Value: ""}
			}
			return &object.String{Value: string(str.Value[idx])}
		}}
	case "indexOf":
		return &object.Builtin{Name: "indexOf", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: -1}
			}
			search := objectToString(args[0])
			return &object.Number{Value: float64(indexOf(str.Value, search))}
		}}
	case "slice":
		return createStringSlice(str)
	case "substring":
		return createStringSubstring(str)
	case "includes":
		return &object.Builtin{Name: "includes", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			search := objectToString(args[0])
			if indexOf(str.Value, search) >= 0 {
				return TRUE
			}
			return FALSE
		}}
	case "startsWith":
		return &object.Builtin{Name: "startsWith", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			search := objectToString(args[0])
			if len(search) > len(str.Value) {
				return FALSE
			}
			if str.Value[:len(search)] == search {
				return TRUE
			}
			return FALSE
		}}
	case "endsWith":
		return &object.Builtin{Name: "endsWith", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return FALSE
			}
			search := objectToString(args[0])
			if len(search) > len(str.Value) {
				return FALSE
			}
			if str.Value[len(str.Value)-len(search):] == search {
				return TRUE
			}
			return FALSE
		}}
	case "replace":
		return &object.Builtin{Name: "replace", Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return str
			}
			search := objectToString(args[0])
			replacement := objectToString(args[1])
			return &object.String{Value: replace(str.Value, search, replacement)}
		}}
	case "repeat":
		return &object.Builtin{Name: "repeat", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.String{Value: ""}
			}
			count := int(args[0].(*object.Number).Value)
			if count <= 0 {
				return &object.String{Value: ""}
			}
			result := ""
			for i := 0; i < count; i++ {
				result += str.Value
			}
			return &object.String{Value: result}
		}}
	case "concat":
		return &object.Builtin{Name: "concat", Fn: func(args ...object.Object) object.Object {
			result := str.Value
			for _, arg := range args {
				result += objectToString(arg)
			}
			return &object.String{Value: result}
		}}
	}
	return UNDEFINED
}

func createStringSlice(str *object.String) *object.Builtin {
	return &object.Builtin{Name: "slice", Fn: func(args ...object.Object) object.Object {
		start, end := 0, len(str.Value)
		if len(args) > 0 {
			start = int(args[0].(*object.Number).Value)
			if start < 0 {
				start = len(str.Value) + start
			}
		}
		if len(args) > 1 {
			end = int(args[1].(*object.Number).Value)
			if end < 0 {
				end = len(str.Value) + end
			}
		}
		if start < 0 {
			start = 0
		}
		if end > len(str.Value) {
			end = len(str.Value)
		}
		if start >= end {
			return &object.String{Value: ""}
		}
		return &object.String{Value: str.Value[start:end]}
	}}
}

func createStringSubstring(str *object.String) *object.Builtin {
	return &object.Builtin{Name: "substring", Fn: func(args ...object.Object) object.Object {
		start, end := 0, len(str.Value)
		if len(args) > 0 {
			start = int(args[0].(*object.Number).Value)
		}
		if len(args) > 1 {
			end = int(args[1].(*object.Number).Value)
		}
		if start < 0 {
			start = 0
		}
		if end > len(str.Value) {
			end = len(str.Value)
		}
		if start > end {
			start, end = end, start
		}
		return &object.String{Value: str.Value[start:end]}
	}}
}

// String helpers
func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && isWhitespace(s[start]) {
		start++
	}
	for end > start && isWhitespace(s[end-1]) {
		end--
	}
	return s[start:end]
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func split(s, sep string) []string {
	if sep == "" {
		result := make([]string, len(s))
		for i := 0; i < len(s); i++ {
			result[i] = string(s[i])
		}
		return result
	}
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func indexOf(s, search string) int {
	if len(search) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(search); i++ {
		if s[i:i+len(search)] == search {
			return i
		}
	}
	return -1
}

func replace(s, search, replacement string) string {
	idx := indexOf(s, search)
	if idx < 0 {
		return s
	}
	return s[:idx] + replacement + s[idx+len(search):]
}
