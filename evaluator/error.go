package evaluator

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/nulang/nulang/object"
)

// ErrorInstance represents a JavaScript Error object instance
type ErrorInstance struct {
	Message string
	Name    string
	Stack   string
}

func (e *ErrorInstance) Type() object.ObjectType { return object.ERROR_OBJ }
func (e *ErrorInstance) Inspect() string {
	if e.Name != "" && e.Name != "Error" {
		return fmt.Sprintf("%s: %s", e.Name, e.Message)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// initErrorConstructor creates the global Error constructor
func initErrorConstructor() object.Object {
	errorClass := &Class{
		Name:          "Error",
		Properties:    make(map[string]object.Object),
		Methods:       make(map[string]*object.Function),
		Getters:       make(map[string]*object.Function),
		Setters:       make(map[string]*object.Function),
		Static:        make(map[string]object.Object),
		NativeMethods: make(map[string]NativeMethod),
	}

	// Constructor
	errorClass.NativeMethods["constructor"] = func(this object.Object, args ...object.Object) object.Object {
		objMap, ok := this.(*object.ObjectMap)
		if !ok {
			return newError("Error: invalid this context")
		}
		
		message := ""
		if len(args) > 0 {
			message = objectToString(args[0])
		}
		
		// Set message property
		objMap.Set("message", &object.String{Value: message})
		
		// Set name property
		objMap.Set("name", &object.String{Value: "Error"})
		
		// Generate stack trace
		stack := generateStackTrace(message)
		objMap.Set("stack", &object.String{Value: stack})
		
		return UNDEFINED
	}

	// toString method
	errorClass.NativeMethods["toString"] = func(this object.Object, args ...object.Object) object.Object {
		objMap, ok := this.(*object.ObjectMap)
		if !ok {
			return &object.String{Value: "Error"}
		}
		
		name := "Error"
		if nameVal, ok := objMap.Get("name"); ok {
			if nameStr, ok := nameVal.(*object.String); ok {
				name = nameStr.Value
			}
		}
		
		message := ""
		if messageVal, ok := objMap.Get("message"); ok {
			if messageStr, ok := messageVal.(*object.String); ok {
				message = messageStr.Value
			}
		}
		
		if message == "" {
			return &object.String{Value: name}
		}
		return &object.String{Value: fmt.Sprintf("%s: %s", name, message)}
	}

	return errorClass
}

// generateStackTrace creates a stack trace string
func generateStackTrace(message string) string {
	var stackLines []string
	stackLines = append(stackLines, fmt.Sprintf("Error: %s", message))
	
	// Get Go stack trace for context (limited depth)
	pcs := make([]uintptr, 10)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	
	for {
		frame, more := frames.Next()
		// Filter out internal Go runtime frames
		if !strings.Contains(frame.Function, "runtime.") &&
		   !strings.Contains(frame.Function, "evaluator.") {
			stackLines = append(stackLines, fmt.Sprintf("    at %s (%s:%d)", 
				frame.Function, frame.File, frame.Line))
		}
		if !more || len(stackLines) > 10 {
			break
		}
	}
	
	// Add a generic Nulang stack entry
	stackLines = append(stackLines, "    at <anonymous>")
	
	return strings.Join(stackLines, "\n")
}

// isErrorInstance checks if an object is an Error instance
func isErrorInstance(obj object.Object) bool {
	if objMap, ok := obj.(*object.ObjectMap); ok {
		if name, ok := objMap.Get("name"); ok {
			if nameStr, ok := name.(*object.String); ok {
				return nameStr.Value == "Error" || 
				       nameStr.Value == "TypeError" || 
				       nameStr.Value == "ReferenceError" ||
				       nameStr.Value == "SyntaxError" ||
				       nameStr.Value == "RangeError"
			}
		}
	}
	return false
}

// getErrorMessage extracts the message from an Error instance or returns the object's string representation
func getErrorMessage(obj object.Object) string {
	// Check if it's an ObjectMap (Error instance)
	if objMap, ok := obj.(*object.ObjectMap); ok {
		if message, ok := objMap.Get("message"); ok {
			if messageStr, ok := message.(*object.String); ok {
				return messageStr.Value
			}
		}
	}
	
	// Check if it's a native Error object
	if err, ok := obj.(*object.Error); ok {
		return err.Message
	}
	
	// Fallback to string representation
	return obj.Inspect()
}

// createErrorObject creates an Error ObjectMap from a message string
func createErrorObject(message string) *object.ObjectMap {
	errorObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	errorObj.Set("message", &object.String{Value: message})
	errorObj.Set("name", &object.String{Value: "Error"})
	errorObj.Set("stack", &object.String{Value: generateStackTrace(message)})
	
	// Add toString method
	errorObj.Set("toString", &object.Builtin{
		Name: "toString",
		Fn: func(args ...object.Object) object.Object {
			if message == "" {
				return &object.String{Value: "Error"}
			}
			return &object.String{Value: fmt.Sprintf("Error: %s", message)}
		},
	})
	
	return errorObj
}
