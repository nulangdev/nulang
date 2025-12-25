package evaluator

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/nulang/nulang/object"
)

// TimerID represents a unique timer identifier
type TimerID int64

var (
	timerCounter int64
	timers       = make(map[TimerID]*Timer)
	timersMutex  sync.RWMutex
)

// Timer represents a timeout or interval
type Timer struct {
	ID         TimerID
	Callback   *object.Function
	Env        *object.Environment
	Delay      time.Duration
	IsInterval bool
	StopChan   chan struct{}
	Running    bool
}

// nextTimerID generates the next timer ID
func nextTimerID() TimerID {
	return TimerID(atomic.AddInt64(&timerCounter, 1))
}

// initTimerBuiltins adds timer functions to builtins
func initTimerBuiltins() {
	// setTimeout(callback, delay, ...args)
	builtins["setTimeout"] = &object.Builtin{Name: "setTimeout", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("setTimeout requires a callback and delay")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			// Check if it's a builtin
			if _, ok := args[0].(*object.Builtin); ok {
				return newError("setTimeout callback must be a user-defined function")
			}
			return newError("setTimeout first argument must be a function")
		}

		var delay float64
		if num, ok := args[1].(*object.Number); ok {
			delay = num.Value
		} else {
			return newError("setTimeout delay must be a number")
		}

		// Create timer
		timer := &Timer{
			ID:         nextTimerID(),
			Callback:   fn,
			Env:        fn.Env,
			Delay:      time.Duration(delay) * time.Millisecond,
			IsInterval: false,
			StopChan:   make(chan struct{}),
			Running:    true,
		}

		timersMutex.Lock()
		timers[timer.ID] = timer
		timersMutex.Unlock()

		// Get additional arguments for callback
		callbackArgs := args[2:]

		// Start timer in goroutine
		go func() {
			select {
			case <-time.After(timer.Delay):
				timersMutex.RLock()
				t, exists := timers[timer.ID]
				timersMutex.RUnlock()
				
				if exists && t.Running {
					executeTimerCallback(fn, fn.Env, callbackArgs)
					
					timersMutex.Lock()
					delete(timers, timer.ID)
					timersMutex.Unlock()
				}
			case <-timer.StopChan:
				return
			}
		}()

		return &object.Number{Value: float64(timer.ID)}
	}}

	// clearTimeout(id)
	builtins["clearTimeout"] = &object.Builtin{Name: "clearTimeout", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}

		var id TimerID
		if num, ok := args[0].(*object.Number); ok {
			id = TimerID(int64(num.Value))
		} else {
			return UNDEFINED
		}

		timersMutex.Lock()
		if timer, ok := timers[id]; ok {
			timer.Running = false
			close(timer.StopChan)
			delete(timers, id)
		}
		timersMutex.Unlock()

		return UNDEFINED
	}}

	// setInterval(callback, delay, ...args)
	builtins["setInterval"] = &object.Builtin{Name: "setInterval", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("setInterval requires a callback and delay")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("setInterval first argument must be a function")
		}

		var delay float64
		if num, ok := args[1].(*object.Number); ok {
			delay = num.Value
		} else {
			return newError("setInterval delay must be a number")
		}

		// Minimum delay to prevent CPU hogging
		if delay < 10 {
			delay = 10
		}

		// Create timer
		timer := &Timer{
			ID:         nextTimerID(),
			Callback:   fn,
			Env:        fn.Env,
			Delay:      time.Duration(delay) * time.Millisecond,
			IsInterval: true,
			StopChan:   make(chan struct{}),
			Running:    true,
		}

		timersMutex.Lock()
		timers[timer.ID] = timer
		timersMutex.Unlock()

		// Get additional arguments for callback
		callbackArgs := args[2:]

		// Start interval in goroutine
		go func() {
			ticker := time.NewTicker(timer.Delay)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					timersMutex.RLock()
					t, exists := timers[timer.ID]
					timersMutex.RUnlock()
					
					if !exists || !t.Running {
						return
					}
					
					executeTimerCallback(fn, fn.Env, callbackArgs)
				case <-timer.StopChan:
					return
				}
			}
		}()

		return &object.Number{Value: float64(timer.ID)}
	}}

	// clearInterval(id)
	builtins["clearInterval"] = &object.Builtin{Name: "clearInterval", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}

		var id TimerID
		if num, ok := args[0].(*object.Number); ok {
			id = TimerID(int64(num.Value))
		} else {
			return UNDEFINED
		}

		timersMutex.Lock()
		if timer, ok := timers[id]; ok {
			timer.Running = false
			close(timer.StopChan)
			delete(timers, id)
		}
		timersMutex.Unlock()

		return UNDEFINED
	}}

	// setImmediate(callback, ...args) - runs on next tick
	builtins["setImmediate"] = &object.Builtin{Name: "setImmediate", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("setImmediate requires a callback")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("setImmediate first argument must be a function")
		}

		id := nextTimerID()
		callbackArgs := args[1:]

		// Run immediately in goroutine
		go func() {
			executeTimerCallback(fn, fn.Env, callbackArgs)
		}()

		return &object.Number{Value: float64(id)}
	}}

	// clearImmediate(id) - for API compatibility
	builtins["clearImmediate"] = &object.Builtin{Name: "clearImmediate", Fn: func(args ...object.Object) object.Object {
		// Immediate callbacks run, well, immediately, so this is essentially a no-op
		// but we include it for API compatibility
		return UNDEFINED
	}}

	// process.nextTick-like functionality
	builtins["nextTick"] = &object.Builtin{Name: "nextTick", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("nextTick requires a callback")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("nextTick first argument must be a function")
		}

		callbackArgs := args[1:]

		// Run in goroutine
		go func() {
			executeTimerCallback(fn, fn.Env, callbackArgs)
		}()

		return UNDEFINED
	}}
}

// executeTimerCallback executes a timer callback function
func executeTimerCallback(fn *object.Function, env *object.Environment, args []object.Object) object.Object {
	extendedEnv := extendFunctionEnv(fn, args)
	return Eval(fn.Body, extendedEnv)
}

// WaitForTimers blocks until all pending timers complete
// This is useful for scripts that use timers
func WaitForTimers(timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	
	for {
		timersMutex.RLock()
		count := len(timers)
		timersMutex.RUnlock()
		
		if count == 0 || time.Now().After(deadline) {
			return
		}
		
		time.Sleep(10 * time.Millisecond)
	}
}

// ClearAllTimers stops and removes all pending timers
func ClearAllTimers() {
	timersMutex.Lock()
	defer timersMutex.Unlock()
	
	for id, timer := range timers {
		timer.Running = false
		close(timer.StopChan)
		delete(timers, id)
	}
}

// HasPendingTimers returns true if there are pending timers
func HasPendingTimers() bool {
	timersMutex.RLock()
	defer timersMutex.RUnlock()
	return len(timers) > 0
}

// Sleep function for blocking delay
func initSleepFunction() *object.Builtin {
	return &object.Builtin{Name: "sleep", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return UNDEFINED
		}

		var delay float64
		if num, ok := args[0].(*object.Number); ok {
			delay = num.Value
		} else {
			return newError("sleep requires a number argument")
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)
		return UNDEFINED
	}}
}
