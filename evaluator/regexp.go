package evaluator

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/nulang/nulang/object"
)

// RegExp represents a regular expression object
type RegExp struct {
	Pattern    string
	Flags      string
	Regexp     *regexp2.Regexp
	Global     bool
	IgnoreCase bool
	Multiline  bool
}

func (r *RegExp) Type() object.ObjectType { return object.REGEXP_OBJ }
func (r *RegExp) Inspect() string         { return "/" + r.Pattern + "/" + r.Flags }

// initRegExpConstructor creates the RegExp constructor
func initRegExpConstructor() *object.Builtin {
	return &object.Builtin{Name: "RegExp", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("RegExp requires a pattern")
		}

		pattern := objectToString(args[0])
		flags := ""
		if len(args) > 1 {
			flags = objectToString(args[1])
		}

		return createRegExp(pattern, flags)
	}}
}

// createRegExp creates a new RegExp object
func createRegExp(pattern, flags string) object.Object {
	re := &RegExp{
		Pattern: pattern,
		Flags:   flags,
	}

	// Parse flags and build regexp2 options
	options := regexp2.None
	for _, f := range flags {
		switch f {
		case 'i':
			re.IgnoreCase = true
			options |= regexp2.IgnoreCase
		case 'm':
			re.Multiline = true
			options |= regexp2.Multiline
		case 'g':
			re.Global = true
		case 's':
			options |= regexp2.Singleline
		}
	}

	// Convert JavaScript regex to regexp2-compatible format
	goPattern := convertJSRegexToGo(pattern)

	// Compile regex using regexp2
	compiled, err := regexp2.Compile(goPattern, options)
	if err != nil {
		// If compilation fails, create a fallback that never matches
		// This allows the code to continue running even with unsupported regex syntax
		compiled, _ = regexp2.Compile("^$impossible-match$", regexp2.None)
		re.Regexp = compiled
		return createRegExpObject(re)
	}
	re.Regexp = compiled

	// Return ObjectMap with methods
	return createRegExpObject(re)
}

// convertJSRegexToGo converts JavaScript regex escapes to Go-compatible format
func convertJSRegexToGo(pattern string) string {
	result := ""
	i := 0
	for i < len(pattern) {
		if pattern[i] == '\\' && i+1 < len(pattern) {
			next := pattern[i+1]
			switch next {
			case 'u':
				// \uXXXX -> convert to literal Unicode character
				if i+5 < len(pattern) {
					hex := pattern[i+2 : i+6]
					var codePoint int
					_, err := fmt.Sscanf(hex, "%x", &codePoint)
					if err == nil {
						result += string(rune(codePoint))
						i += 6
						continue
					}
				}
				// Invalid \u escape, skip it
				result += pattern[i : i+2]
				i += 2
			case 'x':
				// \xXX -> convert to literal character
				if i+3 < len(pattern) {
					hex := pattern[i+2 : i+4]
					var codePoint int
					_, err := fmt.Sscanf(hex, "%x", &codePoint)
					if err == nil {
						result += string(rune(codePoint))
						i += 4
						continue
					}
				}
				result += pattern[i : i+2]
				i += 2
			default:
				// Keep other escapes as-is
				result += pattern[i : i+2]
				i += 2
			}
		} else {
			result += string(pattern[i])
			i++
		}
	}
	return result
}

// createRegExpObject wraps RegExp in an ObjectMap with methods
func createRegExpObject(re *RegExp) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Store the regexp
	obj.Set("_regexp", re)
	obj.Set("source", &object.String{Value: re.Pattern})
	obj.Set("flags", &object.String{Value: re.Flags})
	obj.Set("global", nativeBoolToBooleanObject(re.Global))
	obj.Set("ignoreCase", nativeBoolToBooleanObject(re.IgnoreCase))
	obj.Set("multiline", nativeBoolToBooleanObject(re.Multiline))

	// test(string) - returns boolean
	obj.Set("test", &object.Builtin{Name: "test", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		str := objectToString(args[0])
		match, _ := re.Regexp.MatchString(str)
		return nativeBoolToBooleanObject(match)
	}})

	// exec(string) - returns array or null
	obj.Set("exec", &object.Builtin{Name: "exec", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return NULL
		}
		str := objectToString(args[0])

		match, _ := re.Regexp.FindStringMatch(str)
		if match == nil {
			return NULL
		}

		// Create result array from groups
		groups := match.Groups()
		elements := make([]object.Object, len(groups))
		for i, g := range groups {
			elements[i] = &object.String{Value: g.String()}
		}

		result := &object.Array{Elements: elements}
		return result
	}})

	// match(string) - returns array of all matches or null
	obj.Set("match", &object.Builtin{Name: "match", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return NULL
		}
		str := objectToString(args[0])

		var matches []string
		if re.Global {
			// Find all matches
			m, _ := re.Regexp.FindStringMatch(str)
			for m != nil {
				matches = append(matches, m.String())
				m, _ = re.Regexp.FindNextMatch(m)
			}
		} else {
			m, _ := re.Regexp.FindStringMatch(str)
			if m != nil {
				matches = []string{m.String()}
			}
		}

		if len(matches) == 0 {
			return NULL
		}

		elements := make([]object.Object, len(matches))
		for i, m := range matches {
			elements[i] = &object.String{Value: m}
		}
		return &object.Array{Elements: elements}
	}})

	// replace(string, replacement) - returns new string
	obj.Set("replace", &object.Builtin{Name: "replace", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("replace requires string and replacement")
		}
		str := objectToString(args[0])
		replacement := objectToString(args[1])

		var result string
		var err error
		if re.Global {
			result, err = re.Regexp.Replace(str, replacement, -1, -1)
		} else {
			result, err = re.Regexp.Replace(str, replacement, 0, 1)
		}
		if err != nil {
			return &object.String{Value: str}
		}
		return &object.String{Value: result}
	}})

	// split(string, limit?) - returns array
	obj.Set("split", &object.Builtin{Name: "split", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Array{Elements: []object.Object{}}
		}
		str := objectToString(args[0])

		limit := -1
		if len(args) > 1 {
			if num, ok := args[1].(*object.Number); ok {
				limit = int(num.Value)
			}
		}

		// regexp2 doesn't have Split, so we implement manually
		var parts []string
		lastIndex := 0
		count := 0

		m, _ := re.Regexp.FindStringMatch(str)
		for m != nil && (limit < 0 || count < limit-1) {
			parts = append(parts, str[lastIndex:m.Index])
			lastIndex = m.Index + m.Length
			count++
			m, _ = re.Regexp.FindNextMatch(m)
		}
		// Add remaining part
		if lastIndex <= len(str) {
			parts = append(parts, str[lastIndex:])
		}

		elements := make([]object.Object, len(parts))
		for i, p := range parts {
			elements[i] = &object.String{Value: p}
		}
		return &object.Array{Elements: elements}
	}})

	// toString()
	obj.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: "/" + re.Pattern + "/" + re.Flags}
	}})

	return obj
}

// Helper function to add regex methods to String
func addStringRegexMethods(str *object.String, propName string) object.Object {
	switch propName {
	case "match":
		return &object.Builtin{Name: "match", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return NULL
			}

			var re *RegExp
			switch arg := args[0].(type) {
			case *RegExp:
				re = arg
			case *object.ObjectMap:
				if reObj, ok := arg.Get("_regexp"); ok {
					if r, ok := reObj.(*RegExp); ok {
						re = r
					}
				}
			case *object.String:
				// Create regex from string
				compiled, err := regexp2.Compile(arg.Value, regexp2.None)
				if err != nil {
					return NULL
				}
				re = &RegExp{Pattern: arg.Value, Regexp: compiled}
			default:
				return NULL
			}

			var matches []string
			if re.Global {
				m, _ := re.Regexp.FindStringMatch(str.Value)
				for m != nil {
					matches = append(matches, m.String())
					m, _ = re.Regexp.FindNextMatch(m)
				}
			} else {
				m, _ := re.Regexp.FindStringMatch(str.Value)
				if m != nil {
					matches = []string{m.String()}
				}
			}

			if len(matches) == 0 {
				return NULL
			}

			elements := make([]object.Object, len(matches))
			for i, m := range matches {
				elements[i] = &object.String{Value: m}
			}
			return &object.Array{Elements: elements}
		}}

	case "replace":
		return &object.Builtin{Name: "replace", Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				// Simple string replace
				if len(args) < 1 {
					return str
				}
				return str
			}

			replacement := objectToString(args[1])

			switch arg := args[0].(type) {
			case *RegExp:
				var result string
				var err error
				if arg.Global {
					result, err = arg.Regexp.Replace(str.Value, replacement, -1, -1)
				} else {
					result, err = arg.Regexp.Replace(str.Value, replacement, 0, 1)
				}
				if err != nil {
					return str
				}
				return &object.String{Value: result}
			case *object.ObjectMap:
				if reObj, ok := arg.Get("_regexp"); ok {
					if re, ok := reObj.(*RegExp); ok {
						var result string
						var err error
						if re.Global {
							result, err = re.Regexp.Replace(str.Value, replacement, -1, -1)
						} else {
							result, err = re.Regexp.Replace(str.Value, replacement, 0, 1)
						}
						if err != nil {
							return str
						}
						return &object.String{Value: result}
					}
				}
				// Treat as regular object, convert to string for search
				searchStr := objectToString(arg)
				result := strings.Replace(str.Value, searchStr, replacement, 1)
				return &object.String{Value: result}
			case *object.String:
				// Simple string replace
				result := strings.Replace(str.Value, arg.Value, replacement, 1)
				return &object.String{Value: result}
			default:
				searchStr := objectToString(args[0])
				result := strings.Replace(str.Value, searchStr, replacement, 1)
				return &object.String{Value: result}
			}
		}}

	case "replaceAll":
		return &object.Builtin{Name: "replaceAll", Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return str
			}
			searchStr := objectToString(args[0])
			replacement := objectToString(args[1])
			result := strings.ReplaceAll(str.Value, searchStr, replacement)
			return &object.String{Value: result}
		}}

	case "search":
		return &object.Builtin{Name: "search", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return &object.Number{Value: -1}
			}

			var re *regexp2.Regexp
			switch arg := args[0].(type) {
			case *RegExp:
				re = arg.Regexp
			case *object.ObjectMap:
				if reObj, ok := arg.Get("_regexp"); ok {
					if r, ok := reObj.(*RegExp); ok {
						re = r.Regexp
					}
				}
			case *object.String:
				compiled, err := regexp2.Compile(arg.Value, regexp2.None)
				if err != nil {
					return &object.Number{Value: -1}
				}
				re = compiled
			default:
				return &object.Number{Value: -1}
			}

			match, _ := re.FindStringMatch(str.Value)
			if match == nil {
				return &object.Number{Value: -1}
			}
			return &object.Number{Value: float64(match.Index)}
		}}
	}
	return nil
}
