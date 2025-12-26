// Package evaluator provides initialization for Reflect and Proxy globals.
package evaluator

import (
	"github.com/nulang/nulang/object"
)

// initReflectProxyGlobals initializes Reflect and Proxy in the global environment
func initReflectProxyGlobals(env *object.Environment) {
	// Set Reflect object
	env.Set("Reflect", initReflect())

	// Set Proxy object with constructor support
	proxyObj := initProxy()
	env.Set("Proxy", proxyObj)
}
