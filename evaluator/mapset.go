package evaluator

import (
	"fmt"
	"strings"

	"github.com/nulang/nulang/object"
)

// NuMap represents a JavaScript-like Map object
type NuMap struct {
	Store map[string]mapEntry
	Order []string // Maintain insertion order
}

type mapEntry struct {
	Key   object.Object
	Value object.Object
}

func (m *NuMap) Type() object.ObjectType { return object.MAP_OBJ }
func (m *NuMap) Inspect() string {
	pairs := []string{}
	for _, key := range m.Order {
		entry := m.Store[key]
		pairs = append(pairs, fmt.Sprintf("%s => %s", entry.Key.Inspect(), entry.Value.Inspect()))
	}
	return "Map { " + strings.Join(pairs, ", ") + " }"
}

// NuSet represents a JavaScript-like Set object
type NuSet struct {
	Store map[string]object.Object
	Order []string // Maintain insertion order
}

func (s *NuSet) Type() object.ObjectType { return object.SET_OBJ }
func (s *NuSet) Inspect() string {
	values := []string{}
	for _, key := range s.Order {
		values = append(values, s.Store[key].Inspect())
	}
	return "Set { " + strings.Join(values, ", ") + " }"
}

// getObjectKey generates a unique key for an object (for map/set storage)
func getObjectKey(obj object.Object) string {
	switch o := obj.(type) {
	case *object.String:
		return "s:" + o.Value
	case *object.Number:
		return fmt.Sprintf("n:%v", o.Value)
	case *object.Boolean:
		return fmt.Sprintf("b:%v", o.Value)
	case *object.Null:
		return "null"
	case *object.Undefined:
		return "undefined"
	default:
		// Use pointer address for reference types
		return fmt.Sprintf("r:%p", obj)
	}
}

// initMapConstructor creates the Map constructor
func initMapConstructor() *object.Builtin {
	return &object.Builtin{Name: "Map", Fn: func(args ...object.Object) object.Object {
		m := &NuMap{
			Store: make(map[string]mapEntry),
			Order: []string{},
		}

		// Initialize from iterable if provided
		if len(args) > 0 {
			if arr, ok := args[0].(*object.Array); ok {
				for _, elem := range arr.Elements {
					if pair, ok := elem.(*object.Array); ok && len(pair.Elements) >= 2 {
						key := getObjectKey(pair.Elements[0])
						m.Store[key] = mapEntry{Key: pair.Elements[0], Value: pair.Elements[1]}
						m.Order = append(m.Order, key)
					}
				}
			}
		}

		return createMapObject(m)
	}}
}

// createMapObject wraps NuMap in an ObjectMap with methods
func createMapObject(m *NuMap) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Store the map
	obj.Set("_map", m)

	// size property
	obj.Set("size", &object.Number{Value: float64(len(m.Store))})

	// set(key, value)
	obj.Set("set", &object.Builtin{Name: "set", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return obj
		}
		key := getObjectKey(args[0])
		if _, exists := m.Store[key]; !exists {
			m.Order = append(m.Order, key)
		}
		m.Store[key] = mapEntry{Key: args[0], Value: args[1]}
		obj.Set("size", &object.Number{Value: float64(len(m.Store))})
		return obj
	}})

	// get(key)
	obj.Set("get", &object.Builtin{Name: "get", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		key := getObjectKey(args[0])
		if entry, ok := m.Store[key]; ok {
			return entry.Value
		}
		return UNDEFINED
	}})

	// has(key)
	obj.Set("has", &object.Builtin{Name: "has", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		_, exists := m.Store[key]
		return nativeBoolToBooleanObject(exists)
	}})

	// delete(key)
	obj.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		if _, exists := m.Store[key]; exists {
			delete(m.Store, key)
			// Remove from order
			for i, k := range m.Order {
				if k == key {
					m.Order = append(m.Order[:i], m.Order[i+1:]...)
					break
				}
			}
			obj.Set("size", &object.Number{Value: float64(len(m.Store))})
			return TRUE
		}
		return FALSE
	}})

	// clear()
	obj.Set("clear", &object.Builtin{Name: "clear", Fn: func(args ...object.Object) object.Object {
		m.Store = make(map[string]mapEntry)
		m.Order = []string{}
		obj.Set("size", &object.Number{Value: float64(0)})
		return UNDEFINED
	}})

	// keys()
	obj.Set("keys", &object.Builtin{Name: "keys", Fn: func(args ...object.Object) object.Object {
		keys := make([]object.Object, len(m.Order))
		for i, key := range m.Order {
			keys[i] = m.Store[key].Key
		}
		return &object.Array{Elements: keys}
	}})

	// values()
	obj.Set("values", &object.Builtin{Name: "values", Fn: func(args ...object.Object) object.Object {
		values := make([]object.Object, len(m.Order))
		for i, key := range m.Order {
			values[i] = m.Store[key].Value
		}
		return &object.Array{Elements: values}
	}})

	// entries()
	obj.Set("entries", &object.Builtin{Name: "entries", Fn: func(args ...object.Object) object.Object {
		entries := make([]object.Object, len(m.Order))
		for i, key := range m.Order {
			entry := m.Store[key]
			entries[i] = &object.Array{Elements: []object.Object{entry.Key, entry.Value}}
		}
		return &object.Array{Elements: entries}
	}})

	// forEach(callback)
	obj.Set("forEach", &object.Builtin{Name: "forEach", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		fn, ok := args[0].(*object.Function)
		if !ok {
			return UNDEFINED
		}
		for _, key := range m.Order {
			entry := m.Store[key]
			applyFunction(fn, []object.Object{entry.Value, entry.Key, obj})
		}
		return UNDEFINED
	}})

	// toString()
	obj.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: m.Inspect()}
	}})

	return obj
}

// initSetConstructor creates the Set constructor
func initSetConstructor() *object.Builtin {
	return &object.Builtin{Name: "Set", Fn: func(args ...object.Object) object.Object {
		s := &NuSet{
			Store: make(map[string]object.Object),
			Order: []string{},
		}

		// Initialize from iterable if provided
		if len(args) > 0 {
			if arr, ok := args[0].(*object.Array); ok {
				for _, elem := range arr.Elements {
					key := getObjectKey(elem)
					if _, exists := s.Store[key]; !exists {
						s.Store[key] = elem
						s.Order = append(s.Order, key)
					}
				}
			}
		}

		return createSetObject(s)
	}}
}

// createSetObject wraps NuSet in an ObjectMap with methods
func createSetObject(s *NuSet) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Store the set
	obj.Set("_set", s)

	// size property
	obj.Set("size", &object.Number{Value: float64(len(s.Store))})

	// add(value)
	obj.Set("add", &object.Builtin{Name: "add", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return obj
		}
		key := getObjectKey(args[0])
		if _, exists := s.Store[key]; !exists {
			s.Store[key] = args[0]
			s.Order = append(s.Order, key)
			obj.Set("size", &object.Number{Value: float64(len(s.Store))})
		}
		return obj
	}})

	// has(value)
	obj.Set("has", &object.Builtin{Name: "has", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		_, exists := s.Store[key]
		return nativeBoolToBooleanObject(exists)
	}})

	// delete(value)
	obj.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		if _, exists := s.Store[key]; exists {
			delete(s.Store, key)
			// Remove from order
			for i, k := range s.Order {
				if k == key {
					s.Order = append(s.Order[:i], s.Order[i+1:]...)
					break
				}
			}
			obj.Set("size", &object.Number{Value: float64(len(s.Store))})
			return TRUE
		}
		return FALSE
	}})

	// clear()
	obj.Set("clear", &object.Builtin{Name: "clear", Fn: func(args ...object.Object) object.Object {
		s.Store = make(map[string]object.Object)
		s.Order = []string{}
		obj.Set("size", &object.Number{Value: float64(0)})
		return UNDEFINED
	}})

	// values() / keys() - they're the same for Set
	obj.Set("values", &object.Builtin{Name: "values", Fn: func(args ...object.Object) object.Object {
		values := make([]object.Object, len(s.Order))
		for i, key := range s.Order {
			values[i] = s.Store[key]
		}
		return &object.Array{Elements: values}
	}})

	obj.Set("keys", &object.Builtin{Name: "keys", Fn: func(args ...object.Object) object.Object {
		values := make([]object.Object, len(s.Order))
		for i, key := range s.Order {
			values[i] = s.Store[key]
		}
		return &object.Array{Elements: values}
	}})

	// entries() - returns [value, value] pairs for compatibility
	obj.Set("entries", &object.Builtin{Name: "entries", Fn: func(args ...object.Object) object.Object {
		entries := make([]object.Object, len(s.Order))
		for i, key := range s.Order {
			val := s.Store[key]
			entries[i] = &object.Array{Elements: []object.Object{val, val}}
		}
		return &object.Array{Elements: entries}
	}})

	// forEach(callback)
	obj.Set("forEach", &object.Builtin{Name: "forEach", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		fn, ok := args[0].(*object.Function)
		if !ok {
			return UNDEFINED
		}
		for _, key := range s.Order {
			val := s.Store[key]
			applyFunction(fn, []object.Object{val, val, obj})
		}
		return UNDEFINED
	}})

	// toString()
	obj.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: s.Inspect()}
	}})

	// Extra Set operations

	// union(otherSet)
	obj.Set("union", &object.Builtin{Name: "union", Fn: func(args ...object.Object) object.Object {
		newSet := &NuSet{
			Store: make(map[string]object.Object),
			Order: []string{},
		}
		// Copy current set
		for _, key := range s.Order {
			newSet.Store[key] = s.Store[key]
			newSet.Order = append(newSet.Order, key)
		}
		// Add from other set
		if len(args) > 0 {
			if other, ok := args[0].(*object.ObjectMap); ok {
				if otherSet, ok := other.Get("_set"); ok {
					if os, ok := otherSet.(*NuSet); ok {
						for _, key := range os.Order {
							if _, exists := newSet.Store[key]; !exists {
								newSet.Store[key] = os.Store[key]
								newSet.Order = append(newSet.Order, key)
							}
						}
					}
				}
			}
		}
		return createSetObject(newSet)
	}})

	// intersection(otherSet)
	obj.Set("intersection", &object.Builtin{Name: "intersection", Fn: func(args ...object.Object) object.Object {
		newSet := &NuSet{
			Store: make(map[string]object.Object),
			Order: []string{},
		}
		if len(args) > 0 {
			if other, ok := args[0].(*object.ObjectMap); ok {
				if otherSet, ok := other.Get("_set"); ok {
					if os, ok := otherSet.(*NuSet); ok {
						for _, key := range s.Order {
							if _, exists := os.Store[key]; exists {
								newSet.Store[key] = s.Store[key]
								newSet.Order = append(newSet.Order, key)
							}
						}
					}
				}
			}
		}
		return createSetObject(newSet)
	}})

	// difference(otherSet)
	obj.Set("difference", &object.Builtin{Name: "difference", Fn: func(args ...object.Object) object.Object {
		newSet := &NuSet{
			Store: make(map[string]object.Object),
			Order: []string{},
		}
		var otherStore map[string]object.Object
		if len(args) > 0 {
			if other, ok := args[0].(*object.ObjectMap); ok {
				if otherSet, ok := other.Get("_set"); ok {
					if os, ok := otherSet.(*NuSet); ok {
						otherStore = os.Store
					}
				}
			}
		}
		for _, key := range s.Order {
			if otherStore == nil || otherStore[key] == nil {
				newSet.Store[key] = s.Store[key]
				newSet.Order = append(newSet.Order, key)
			}
		}
		return createSetObject(newSet)
	}})

	return obj
}

// WeakMap implementation (simplified - doesn't actually have weak references in Go)
type WeakMap struct {
	Store map[string]object.Object
}

func (w *WeakMap) Type() object.ObjectType { return object.MAP_OBJ }
func (w *WeakMap) Inspect() string         { return "WeakMap {}" }

// initWeakMapConstructor creates WeakMap constructor
func initWeakMapConstructor() *object.Builtin {
	return &object.Builtin{Name: "WeakMap", Fn: func(args ...object.Object) object.Object {
		w := &WeakMap{Store: make(map[string]object.Object)}
		return createWeakMapObject(w)
	}}
}

func createWeakMapObject(w *WeakMap) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	obj.Set("_weakmap", w)

	obj.Set("set", &object.Builtin{Name: "set", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return obj
		}
		key := getObjectKey(args[0])
		w.Store[key] = args[1]
		return obj
	}})

	obj.Set("get", &object.Builtin{Name: "get", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}
		key := getObjectKey(args[0])
		if val, ok := w.Store[key]; ok {
			return val
		}
		return UNDEFINED
	}})

	obj.Set("has", &object.Builtin{Name: "has", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		_, exists := w.Store[key]
		return nativeBoolToBooleanObject(exists)
	}})

	obj.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		if _, exists := w.Store[key]; exists {
			delete(w.Store, key)
			return TRUE
		}
		return FALSE
	}})

	return obj
}

// WeakSet implementation
type WeakSet struct {
	Store map[string]bool
}

func (w *WeakSet) Type() object.ObjectType { return object.SET_OBJ }
func (w *WeakSet) Inspect() string         { return "WeakSet {}" }

// initWeakSetConstructor creates WeakSet constructor
func initWeakSetConstructor() *object.Builtin {
	return &object.Builtin{Name: "WeakSet", Fn: func(args ...object.Object) object.Object {
		w := &WeakSet{Store: make(map[string]bool)}
		return createWeakSetObject(w)
	}}
}

func createWeakSetObject(w *WeakSet) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	obj.Set("_weakset", w)

	obj.Set("add", &object.Builtin{Name: "add", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return obj
		}
		key := getObjectKey(args[0])
		w.Store[key] = true
		return obj
	}})

	obj.Set("has", &object.Builtin{Name: "has", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		return nativeBoolToBooleanObject(w.Store[key])
	}})

	obj.Set("delete", &object.Builtin{Name: "delete", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		key := getObjectKey(args[0])
		if w.Store[key] {
			delete(w.Store, key)
			return TRUE
		}
		return FALSE
	}})

	return obj
}
