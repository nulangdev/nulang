package evaluator

import (
	"github.com/nulang/nulang/object"
)

// Promise states
const (
	PromisePending   = "pending"
	PromiseFulfilled = "fulfilled"
	PromiseRejected  = "rejected"
)

func initPromiseConstructor() *object.ObjectMap {
	promiseObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Promise.resolve(value)
	promiseObj.Set("resolve", &object.Builtin{Name: "resolve", Fn: func(args ...object.Object) object.Object {
		var value object.Object = UNDEFINED
		if len(args) > 0 {
			value = args[0]
		}
		return &object.Promise{
			State: PromiseFulfilled,
			Value: value,
		}
	}})

	// Promise.reject(reason)
	promiseObj.Set("reject", &object.Builtin{Name: "reject", Fn: func(args ...object.Object) object.Object {
		var reason object.Object = UNDEFINED
		if len(args) > 0 {
			reason = args[0]
		}
		return &object.Promise{
			State:  PromiseRejected,
			Reason: reason,
		}
	}})

	// Promise.all(iterable)
	promiseObj.Set("all", &object.Builtin{Name: "all", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Promise.all requires an iterable argument")
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("Promise.all requires an array argument")
		}

		results := make([]object.Object, len(arr.Elements))
		for i, elem := range arr.Elements {
			if promise, ok := elem.(*object.Promise); ok {
				if promise.State == PromiseRejected {
					return &object.Promise{
						State:  PromiseRejected,
						Reason: promise.Reason,
					}
				}
				results[i] = promise.Value
			} else {
				results[i] = elem
			}
		}

		return &object.Promise{
			State: PromiseFulfilled,
			Value: &object.Array{Elements: results},
		}
	}})

	// Promise.race(iterable)
	promiseObj.Set("race", &object.Builtin{Name: "race", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Promise.race requires an iterable argument")
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("Promise.race requires an array argument")
		}

		if len(arr.Elements) == 0 {
			return &object.Promise{State: PromisePending}
		}

		// Return the first settled promise or value
		for _, elem := range arr.Elements {
			if promise, ok := elem.(*object.Promise); ok {
				if promise.State != PromisePending {
					return promise
				}
			} else {
				return &object.Promise{
					State: PromiseFulfilled,
					Value: elem,
				}
			}
		}

		return &object.Promise{State: PromisePending}
	}})

	// Promise.allSettled(iterable)
	promiseObj.Set("allSettled", &object.Builtin{Name: "allSettled", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("Promise.allSettled requires an iterable argument")
		}
		arr, ok := args[0].(*object.Array)
		if !ok {
			return newError("Promise.allSettled requires an array argument")
		}

		results := make([]object.Object, len(arr.Elements))
		for i, elem := range arr.Elements {
			result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			if promise, ok := elem.(*object.Promise); ok {
				result.Set("status", &object.String{Value: promise.State})
				if promise.State == PromiseFulfilled {
					result.Set("value", promise.Value)
				} else if promise.State == PromiseRejected {
					result.Set("reason", promise.Reason)
				}
			} else {
				result.Set("status", &object.String{Value: PromiseFulfilled})
				result.Set("value", elem)
			}
			results[i] = result
		}

		return &object.Promise{
			State: PromiseFulfilled,
			Value: &object.Array{Elements: results},
		}
	}})

	return promiseObj
}

// CreatePromise creates a new Promise by executing the executor function
func CreatePromise(executor *object.Function, env *object.Environment) *object.Promise {
	promise := &object.Promise{
		State: PromisePending,
	}

	// Create resolve function
	resolveFn := &object.Builtin{Name: "resolve", Fn: func(args ...object.Object) object.Object {
		if promise.State != PromisePending {
			return UNDEFINED
		}
		promise.State = PromiseFulfilled
		if len(args) > 0 {
			promise.Value = args[0]
		} else {
			promise.Value = UNDEFINED
		}
		return UNDEFINED
	}}

	// Create reject function
	rejectFn := &object.Builtin{Name: "reject", Fn: func(args ...object.Object) object.Object {
		if promise.State != PromisePending {
			return UNDEFINED
		}
		promise.State = PromiseRejected
		if len(args) > 0 {
			promise.Reason = args[0]
		} else {
			promise.Reason = UNDEFINED
		}
		return UNDEFINED
	}}

	// Execute the executor
	execEnv := extendFunctionEnv(executor, []object.Object{resolveFn, rejectFn})
	result := Eval(executor.Body, execEnv)

	// If executor throws, reject promise
	if errObj, ok := result.(*object.Error); ok {
		if promise.State == PromisePending {
			promise.State = PromiseRejected
			promise.Reason = &object.String{Value: errObj.Message}
		}
	}

	return promise
}

// PromiseThen adds a then handler
func PromiseThen(promise *object.Promise, onFulfilled, onRejected *object.Function, env *object.Environment) *object.Promise {
	newPromise := &object.Promise{State: PromisePending}

	switch promise.State {
	case PromiseFulfilled:
		if onFulfilled != nil {
			fnEnv := extendFunctionEnv(onFulfilled, []object.Object{promise.Value})
			result := unwrapReturnValue(Eval(onFulfilled.Body, fnEnv))
			if errObj, ok := result.(*object.Error); ok {
				newPromise.State = PromiseRejected
				newPromise.Reason = &object.String{Value: errObj.Message}
			} else {
				newPromise.State = PromiseFulfilled
				newPromise.Value = result
			}
		} else {
			newPromise.State = PromiseFulfilled
			newPromise.Value = promise.Value
		}
	case PromiseRejected:
		if onRejected != nil {
			fnEnv := extendFunctionEnv(onRejected, []object.Object{promise.Reason})
			result := unwrapReturnValue(Eval(onRejected.Body, fnEnv))
			if errObj, ok := result.(*object.Error); ok {
				newPromise.State = PromiseRejected
				newPromise.Reason = &object.String{Value: errObj.Message}
			} else {
				newPromise.State = PromiseFulfilled
				newPromise.Value = result
			}
		} else {
			newPromise.State = PromiseRejected
			newPromise.Reason = promise.Reason
		}
	}

	return newPromise
}

// evalPromiseProperty handles promise instance methods like .then(), .catch(), .finally()
func evalPromiseProperty(promise *object.Promise, prop string, env *object.Environment) object.Object {
	switch prop {
	case "then":
		return &object.Builtin{Name: "then", Fn: func(args ...object.Object) object.Object {
			var onFulfilled, onRejected *object.Function
			if len(args) > 0 {
				if fn, ok := args[0].(*object.Function); ok {
					onFulfilled = fn
				}
			}
			if len(args) > 1 {
				if fn, ok := args[1].(*object.Function); ok {
					onRejected = fn
				}
			}
			return PromiseThen(promise, onFulfilled, onRejected, env)
		}}
	case "catch":
		return &object.Builtin{Name: "catch", Fn: func(args ...object.Object) object.Object {
			var onRejected *object.Function
			if len(args) > 0 {
				if fn, ok := args[0].(*object.Function); ok {
					onRejected = fn
				}
			}
			return PromiseThen(promise, nil, onRejected, env)
		}}
	case "finally":
		return &object.Builtin{Name: "finally", Fn: func(args ...object.Object) object.Object {
			var onFinally *object.Function
			if len(args) > 0 {
				if fn, ok := args[0].(*object.Function); ok {
					onFinally = fn
				}
			}
			
			newPromise := &object.Promise{State: promise.State}
			
			if onFinally != nil {
				fnEnv := extendFunctionEnv(onFinally, []object.Object{})
				Eval(onFinally.Body, fnEnv)
			}
			
			if promise.State == PromiseFulfilled {
				newPromise.Value = promise.Value
			} else if promise.State == PromiseRejected {
				newPromise.Reason = promise.Reason
			}
			
			return newPromise
		}}
	case "state":
		return &object.String{Value: promise.State}
	case "value":
		if promise.State == PromiseFulfilled {
			return promise.Value
		}
		return UNDEFINED
	case "reason":
		if promise.State == PromiseRejected {
			return promise.Reason
		}
		return UNDEFINED
	}
	return UNDEFINED
}
