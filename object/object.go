// Package object defines the object system for Nulang runtime.
package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/nulang/nulang/ast"
)

// ObjectType represents the type of an object
type ObjectType string

const (
	NUMBER_OBJ       ObjectType = "NUMBER"
	STRING_OBJ       ObjectType = "STRING"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	UNDEFINED_OBJ    ObjectType = "UNDEFINED"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	OBJECT_OBJ       ObjectType = "OBJECT"
	BREAK_OBJ        ObjectType = "BREAK"
	CONTINUE_OBJ     ObjectType = "CONTINUE"
	SYMBOL_OBJ       ObjectType = "SYMBOL"
	BIGINT_OBJ       ObjectType = "BIGINT"
	PROMISE_OBJ      ObjectType = "PROMISE"
	DATE_OBJ         ObjectType = "DATE"
	REGEXP_OBJ       ObjectType = "REGEXP"
	MAP_OBJ          ObjectType = "MAP"
	SET_OBJ          ObjectType = "SET"
	BUFFER_OBJ       ObjectType = "BUFFER"
	CLASS_OBJ        ObjectType = "CLASS"
	STREAM_OBJ       ObjectType = "STREAM"
)

// Object interface
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Hashable interface for objects that can be used as hash keys
type Hashable interface {
	HashKey() HashKey
}

// HashKey represents a hash key for map lookups
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Number represents a number value
type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) Inspect() string  { return fmt.Sprintf("%g", n.Value) }
func (n *Number) HashKey() HashKey {
	return HashKey{Type: n.Type(), Value: uint64(n.Value)}
}

// String represents a string value
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// Null represents null value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Undefined represents undefined value
type Undefined struct{}

func (u *Undefined) Type() ObjectType { return UNDEFINED_OBJ }
func (u *Undefined) Inspect() string  { return "undefined" }

// ReturnValue wraps a returned value
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents an error
type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("Error at line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// Break represents break statement signal
type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

// Continue represents continue statement signal
type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }

// Function represents a function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Name       string
	IsAsync    bool
	Properties map[string]Object // Properties attached to the function (e.g., VERSION)
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	if f.IsAsync {
		out.WriteString("async ")
	}
	out.WriteString("function")
	if f.Name != "" {
		out.WriteString(" " + f.Name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Get retrieves a property from the function
func (f *Function) Get(key string) (Object, bool) {
	if f.Properties == nil {
		return nil, false
	}
	val, ok := f.Properties[key]
	return val, ok
}

// Set sets a property on the function
func (f *Function) Set(key string, value Object) {
	if f.Properties == nil {
		f.Properties = make(map[string]Object)
	}
	f.Properties[key] = value
}

// BuiltinFunction represents a built-in function signature
type BuiltinFunction func(args ...Object) Object

// Builtin represents a built-in function
type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return fmt.Sprintf("function %s() { [native code] }", b.Name) }

// Array represents an array
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// ObjectPair represents a key-value pair in an object
type ObjectPair struct {
	Key   Object
	Value Object
}

// ObjectMap represents a JavaScript object
type ObjectMap struct {
	Pairs     map[string]ObjectPair
	Prototype *ObjectMap
}

func (o *ObjectMap) Type() ObjectType { return OBJECT_OBJ }
func (o *ObjectMap) Inspect() string {
	return o.InspectWithDepth(0)
}

// InspectWithDepth returns a string representation with recursion depth limit
func (o *ObjectMap) InspectWithDepth(depth int) string {
	if o.Pairs == nil {
		return "{}"
	}
	if depth > 2 {
		return "{...}"
	}
	var out bytes.Buffer
	pairs := []string{}
	for key, pair := range o.Pairs {
		keyStr := key
		valStr := "undefined"
		if pair.Value != nil {
			// Check if it's a self-reference or circular
			if innerMap, ok := pair.Value.(*ObjectMap); ok {
				if innerMap == o {
					valStr = "[Circular]"
				} else {
					valStr = innerMap.InspectWithDepth(depth + 1)
				}
			} else {
				valStr = pair.Value.Inspect()
			}
		}
		pairs = append(pairs, fmt.Sprintf("%s: %s", keyStr, valStr))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Get retrieves a property from the object or its prototype chain
func (o *ObjectMap) Get(key string) (Object, bool) {
	if o.Pairs == nil {
		if o.Prototype != nil {
			return o.Prototype.Get(key)
		}
		return nil, false
	}
	if pair, ok := o.Pairs[key]; ok {
		return pair.Value, true
	}
	if o.Prototype != nil {
		return o.Prototype.Get(key)
	}
	return nil, false
}

// Set sets a property on the object
func (o *ObjectMap) Set(key string, value Object) {
	o.Pairs[key] = ObjectPair{
		Key:   &String{Value: key},
		Value: value,
	}
}

// Symbol represents a Symbol value
type Symbol struct {
	Description string
	ID          uint64
}

var symbolCounter uint64 = 0

func NewSymbol(description string) *Symbol {
	symbolCounter++
	return &Symbol{Description: description, ID: symbolCounter}
}

func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Inspect() string  { return fmt.Sprintf("Symbol(%s)", s.Description) }

// BigInt represents a bigint value
type BigInt struct {
	Value string // stored as string for arbitrary precision
}

func (b *BigInt) Type() ObjectType { return BIGINT_OBJ }
func (b *BigInt) Inspect() string  { return b.Value + "n" }

// Promise represents a Promise
type Promise struct {
	State    string // "pending", "fulfilled", "rejected"
	Value    Object
	Reason   Object
	Handlers []func(Object)
}

func (p *Promise) Type() ObjectType { return PROMISE_OBJ }
func (p *Promise) Inspect() string  { return fmt.Sprintf("Promise { <%s> }", p.State) }

// Date represents a Date object
type Date struct {
	Value int64 // Unix timestamp in milliseconds
}

func (d *Date) Type() ObjectType { return DATE_OBJ }
func (d *Date) Inspect() string  { return fmt.Sprintf("Date(%d)", d.Value) }

// RegExp represents a RegExp object
type RegExp struct {
	Pattern string
	Flags   string
}

func (r *RegExp) Type() ObjectType { return REGEXP_OBJ }
func (r *RegExp) Inspect() string  { return fmt.Sprintf("/%s/%s", r.Pattern, r.Flags) }

// Map represents a Map object
type Map struct {
	Store map[HashKey]ObjectPair
}

func NewMap() *Map {
	return &Map{Store: make(map[HashKey]ObjectPair)}
}

func (m *Map) Type() ObjectType { return MAP_OBJ }
func (m *Map) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range m.Store {
		pairs = append(pairs, fmt.Sprintf("%s => %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("Map(")
	out.WriteString(fmt.Sprintf("%d", len(m.Store)))
	out.WriteString(") { ")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")
	return out.String()
}

// Set represents a Set object
type Set struct {
	Store map[HashKey]Object
}

func NewSet() *Set {
	return &Set{Store: make(map[HashKey]Object)}
}

func (s *Set) Type() ObjectType { return SET_OBJ }
func (s *Set) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, elem := range s.Store {
		elements = append(elements, elem.Inspect())
	}
	out.WriteString("Set(")
	out.WriteString(fmt.Sprintf("%d", len(s.Store)))
	out.WriteString(") { ")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString(" }")
	return out.String()
}

// Buffer represents a Buffer object
type Buffer struct {
	Data []byte
}

func (b *Buffer) Type() ObjectType { return BUFFER_OBJ }
func (b *Buffer) Inspect() string {
	return fmt.Sprintf("<Buffer %x>", b.Data)
}
