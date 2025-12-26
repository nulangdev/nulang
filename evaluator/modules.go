package evaluator

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/nulang/nulang/ast"
	"github.com/nulang/nulang/lexer"
	"github.com/nulang/nulang/object"
	"github.com/nulang/nulang/parser"
)

// ModuleCache stores loaded modules
var moduleCache = make(map[string]*object.ObjectMap)
var moduleCacheMutex sync.RWMutex

// CurrentModulePath tracks the current module being evaluated
var CurrentModulePath string = ""

// ClearModuleCache clears the module cache for watch mode restarts
func ClearModuleCache() {
	moduleCacheMutex.Lock()
	defer moduleCacheMutex.Unlock()
	moduleCache = make(map[string]*object.ObjectMap)
}

// getBuiltinModule returns a builtin module by name if it exists
func getBuiltinModule(name string) (*object.ObjectMap, bool) {
	switch name {
	case "fs":
		return initFsModule(), true
	case "path":
		return initPathModule(), true
	case "crypto":
		return initCryptoModule(), true
	case "os":
		return initOsModule(), true
	case "http":
		return initHttpModule(), true
	case "stream":
		return initStreamModule(), true
	case "url":
		return initURLModule(), true
	case "querystring":
		return initQueryStringModule(), true
	case "events":
		return initEventsModule(), true
	case "util":
		return initUtilModule(), true
	case "child_process":
		return initChildProcessModule(), true
	case "assert":
		return initAssertModule(), true
	case "readline":
		return initReadlineModule(), true
	case "zlib":
		return initZlibModule(), true
	case "net":
		return initNetModule(), true
	case "dns":
		return initDNSModule(), true
	case "dgram":
		return initDgramModule(), true
	case "vm":
		return initVMModule(), true
	default:
		return nil, false
	}
}

// LoadModule loads a module from a path
func LoadModule(modulePath string, basePath string) (object.Object, error) {
	// Check for built-in modules
	if mod, ok := getBuiltinModule(modulePath); ok {
		return mod, nil
	}

	// Resolve the module path
	resolvedPath := resolveModulePath(modulePath, basePath)

	// Check cache
	moduleCacheMutex.RLock()
	if cached, ok := moduleCache[resolvedPath]; ok {
		moduleCacheMutex.RUnlock()
		return cached, nil
	}
	moduleCacheMutex.RUnlock()

	// Read the module file
	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, err
	}

	// Save current module path
	prevModulePath := CurrentModulePath
	CurrentModulePath = resolvedPath

	// Parse and evaluate the module
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		CurrentModulePath = prevModulePath
		return nil, &ModuleError{Message: p.Errors()[0]}
	}

	// Create a new environment for the module
	moduleEnv := object.NewEnvironment()
	
	// Create exports object
	exports := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	module := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	module.Set("exports", exports)
	
	moduleEnv.Set("exports", exports)
	moduleEnv.Set("module", module)
	moduleEnv.Set("__filename", &object.String{Value: resolvedPath})
	moduleEnv.Set("__dirname", &object.String{Value: filepath.Dir(resolvedPath)})

	// Evaluate the module
	result := evalModuleProgram(program, moduleEnv)
	
	// Restore module path
	CurrentModulePath = prevModulePath

	if errObj, ok := result.(*object.Error); ok {
		return nil, &ModuleError{Message: errObj.Message}
	}

	// Get the final exports
	finalExports, _ := moduleEnv.Get("exports")
	if modObj, ok := moduleEnv.Get("module"); ok {
		if modMap, ok := modObj.(*object.ObjectMap); ok {
			if exp, ok := modMap.Get("exports"); ok {
				finalExports = exp
			}
		}
	}

	// Cache the module
	if exportMap, ok := finalExports.(*object.ObjectMap); ok {
		moduleCacheMutex.Lock()
		moduleCache[resolvedPath] = exportMap
		moduleCacheMutex.Unlock()
		return exportMap, nil
	}

	return finalExports, nil
}

func resolveModulePath(modulePath string, basePath string) string {
	// If it starts with ./ or ../, resolve relative to basePath
	if len(modulePath) > 0 && (modulePath[0] == '.' || modulePath[0] == '/') {
		if basePath != "" {
			baseDir := filepath.Dir(basePath)
			resolved := filepath.Join(baseDir, modulePath)
			
			// Try with .nu extension if not present
			if filepath.Ext(resolved) == "" {
				if _, err := os.Stat(resolved + ".nu"); err == nil {
					return resolved + ".nu"
				}
				// Try index.nu in directory
				indexPath := filepath.Join(resolved, "index.nu")
				if _, err := os.Stat(indexPath); err == nil {
					return indexPath
				}
			}
			return resolved
		}
	}
	
	// Check .nu_modules directory for package imports
	nuModulesPath := FindNuModulesPath(basePath)
	if nuModulesPath != "" {
		pkgPath := filepath.Join(nuModulesPath, modulePath, "index.nu")
		if _, err := os.Stat(pkgPath); err == nil {
			return pkgPath
		}
	}
	
	// Try current directory
	if filepath.Ext(modulePath) == "" {
		if _, err := os.Stat(modulePath + ".nu"); err == nil {
			return modulePath + ".nu"
		}
	}
	
	return modulePath
}

// evalModuleProgram evaluates a module program with export handling
func evalModuleProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = UNDEFINED

	for _, statement := range program.Statements {
		// Handle import statements
		if importStmt, ok := statement.(*ast.ImportStatement); ok {
			result = evalImportStatement(importStmt, env)
			if isError(result) {
				return result
			}
			continue
		}

		// Handle export statements
		if exportStmt, ok := statement.(*ast.ExportStatement); ok {
			result = evalExportStatement(exportStmt, env)
			if isError(result) {
				return result
			}
			continue
		}

		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalImportStatement(is *ast.ImportStatement, env *object.Environment) object.Object {
	source := is.Source.Value
	
	moduleObj, err := LoadModule(source, CurrentModulePath)
	if err != nil {
		return newError("failed to import module '%s': %s", source, err.Error())
	}

	// Handle different import types
	if is.NamespaceAs != nil {
		// import * as name from "module"
		env.Set(is.NamespaceAs.Value, moduleObj)
	} else if is.Default != nil {
		// import name from "module"
		if modMap, ok := moduleObj.(*object.ObjectMap); ok {
			if defaultExport, ok := modMap.Get("default"); ok {
				env.Set(is.Default.Value, defaultExport)
			} else {
				env.Set(is.Default.Value, moduleObj)
			}
		} else {
			env.Set(is.Default.Value, moduleObj)
		}
	} else if len(is.Named) > 0 {
		// import { a, b } or { a as b } from "module"
		if modMap, ok := moduleObj.(*object.ObjectMap); ok {
			for _, importName := range is.Named {
				if val, ok := modMap.Get(importName.Name.Value); ok {
					// Use alias if present, otherwise use original name
					localName := importName.Name.Value
					if importName.Alias != nil {
						localName = importName.Alias.Value
					}
					env.Set(localName, val)
				} else {
					return newError("module '%s' does not export '%s'", source, importName.Name.Value)
				}
			}
		}
	} else {
		// import "module" - side effect only
		// Module was already loaded and executed
	}

	return UNDEFINED
}

func evalExportStatement(es *ast.ExportStatement, env *object.Environment) object.Object {
	exports, _ := env.Get("exports")
	exportsMap, ok := exports.(*object.ObjectMap)
	if !ok {
		exportsMap = &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		env.Set("exports", exportsMap)
	}

	if es.Default != nil {
		// export default expression
		value := Eval(es.Default, env)
		if isError(value) {
			return value
		}
		exportsMap.Set("default", value)
	}

	// Handle export statements for named exports
	for _, stmt := range es.Named {
		switch s := stmt.(type) {
		case *ast.LetStatement:
			value := Eval(s, env)
			if isError(value) {
				return value
			}
			exportsMap.Set(s.Name.Value, value)
		case *ast.ConstStatement:
			value := Eval(s, env)
			if isError(value) {
				return value
			}
			exportsMap.Set(s.Name.Value, value)
		case *ast.VarStatement:
			value := Eval(s, env)
			if isError(value) {
				return value
			}
			// For function declarations
			if s.Value != nil {
				if fn, ok := s.Value.(*ast.FunctionLiteral); ok && fn.Name != nil {
					exportsMap.Set(fn.Name.Value, value)
				} else {
					exportsMap.Set(s.Name.Value, value)
				}
			} else {
				exportsMap.Set(s.Name.Value, value)
			}
		}
	}

	return UNDEFINED
}

// ModuleError represents a module loading error
type ModuleError struct {
	Message string
}

func (e *ModuleError) Error() string {
	return e.Message
}
