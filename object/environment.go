package object

// Environment holds variable bindings
type Environment struct {
	store  map[string]Object
	consts map[string]bool
	outer  *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, consts: c, outer: nil}
}

// NewEnclosedEnvironment creates an enclosed environment
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a variable in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// SetConst sets a constant in the environment
func (e *Environment) SetConst(name string, val Object) Object {
	e.store[name] = val
	e.consts[name] = true
	return val
}

// IsConst checks if a variable is a constant
func (e *Environment) IsConst(name string) bool {
	if isConst, ok := e.consts[name]; ok {
		return isConst
	}
	if e.outer != nil {
		return e.outer.IsConst(name)
	}
	return false
}

// Update updates an existing variable (for let/var reassignment)
func (e *Environment) Update(name string, val Object) (Object, bool) {
	if _, ok := e.store[name]; ok {
		if e.consts[name] {
			return nil, false // Cannot reassign const
		}
		e.store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}

// Has checks if a variable exists in scope
func (e *Environment) Has(name string) bool {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Has(name)
	}
	return ok
}

// Delete deletes a variable from the environment
func (e *Environment) Delete(name string) bool {
	if _, ok := e.store[name]; ok {
		if e.consts[name] {
			return false
		}
		delete(e.store, name)
		return true
	}
	return false
}

// GetStore returns the store for iteration
func (e *Environment) GetStore() map[string]Object {
	return e.store
}
