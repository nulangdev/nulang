package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nulang/nulang/evaluator"
	"github.com/nulang/nulang/lexer"
	"github.com/nulang/nulang/object"
	"github.com/nulang/nulang/parser"
)

// WatchContext holds the state for watch mode
type WatchContext struct {
	MainFile       string
	WatchedFiles   map[string]time.Time // filepath -> last mod time
	Mutex          sync.RWMutex
	StopChan       chan struct{}
	RestartChan    chan struct{}
	ActiveServers  []*object.ObjectMap // Track HTTP servers for cleanup
	ServersMutex   sync.Mutex
}

// NewWatchContext creates a new watch context
func NewWatchContext(mainFile string) *WatchContext {
	absPath, _ := filepath.Abs(mainFile)
	return &WatchContext{
		MainFile:      absPath,
		WatchedFiles:  make(map[string]time.Time),
		StopChan:      make(chan struct{}),
		RestartChan:   make(chan struct{}),
		ActiveServers: make([]*object.ObjectMap, 0),
	}
}

// RunWithWatchV2 implements an improved watch mode
func RunWithWatchV2(filename string) {
	ctx := NewWatchContext(filename)

	// Print header
	printWatchHeader(ctx.MainFile)

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial run and dependency discovery
	ctx.DiscoverAndRun()

	// Watch loop
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\033[1;33m⏹  Watch mode stopped\033[0m")
			ctx.Cleanup()
			return
		case <-ticker.C:
			if ctx.CheckForChanges() {
				ctx.Restart()
			}
		}
	}
}

// printWatchHeader prints the initial watch mode header
func printWatchHeader(filename string) {
	fmt.Println("\033[38;5;39m┌──────────────────────────────────────────────────┐\033[0m")
	fmt.Println("\033[38;5;39m│\033[0m  \033[1;36m👁  Nulang Watch Mode\033[0m                            \033[38;5;39m│\033[0m")
	fmt.Println("\033[38;5;39m├──────────────────────────────────────────────────┤\033[0m")
	fmt.Printf("\033[38;5;39m│\033[0m  📄 Main:  \033[1;33m%-36s\033[0m \033[38;5;39m│\033[0m\n", filepath.Base(filename))
	fmt.Println("\033[38;5;39m│\033[0m  \033[90mPress Ctrl+C to stop\033[0m                            \033[38;5;39m│\033[0m")
	fmt.Println("\033[38;5;39m└──────────────────────────────────────────────────┘\033[0m")
	fmt.Println()
}

// DiscoverAndRun discovers all dependencies and runs the script
func (ctx *WatchContext) DiscoverAndRun() {
	ctx.Mutex.Lock()
	// Clear watched files
	ctx.WatchedFiles = make(map[string]time.Time)
	ctx.Mutex.Unlock()

	// Add main file
	ctx.AddWatchedFile(ctx.MainFile)

	// Run the main file and collect dependencies
	ctx.RunWithDependencyDiscovery()
}

// AddWatchedFile adds a file to the watch list
func (ctx *WatchContext) AddWatchedFile(filepath string) {
	info, err := os.Stat(filepath)
	if err != nil {
		return
	}

	ctx.Mutex.Lock()
	defer ctx.Mutex.Unlock()

	if _, exists := ctx.WatchedFiles[filepath]; !exists {
		ctx.WatchedFiles[filepath] = info.ModTime()
	}
}

// CheckForChanges checks if any watched file has changed
func (ctx *WatchContext) CheckForChanges() bool {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	for filePath, lastMod := range ctx.WatchedFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		if info.ModTime().After(lastMod) {
			return true
		}
	}

	return false
}

// Restart restarts the application
func (ctx *WatchContext) Restart() {
	clearScreen()
	fmt.Printf("\033[1;32m🔄 File changed at %s\033[0m\n", time.Now().Format("15:04:05"))
	printWatchedFiles(ctx)
	fmt.Println("\033[90m" + strings.Repeat("─", 50) + "\033[0m")
	fmt.Println()

	// Cleanup previous run
	ctx.Cleanup()

	// Clear module cache
	evaluator.ClearModuleCache()

	// Run again
	ctx.DiscoverAndRun()
}

// Cleanup stops all active resources (servers, timers, etc)
func (ctx *WatchContext) Cleanup() {
	// Stop all active servers
	ctx.ServersMutex.Lock()
	for _, server := range ctx.ActiveServers {
		if closeFn, ok := server.Get("_closeGoServer"); ok {
			if fn, ok := closeFn.(*object.Builtin); ok {
				fn.Fn()
			}
		}
	}
	ctx.ActiveServers = make([]*object.ObjectMap, 0)
	ctx.ServersMutex.Unlock()

	// Clear all timers
	evaluator.ClearAllTimers()

	// Reset async task counter
	evaluator.ResetAsyncTasks()

	// Small delay to allow goroutines to finish
	time.Sleep(50 * time.Millisecond)
}

// RunWithDependencyDiscovery runs the script and discovers dependencies
func (ctx *WatchContext) RunWithDependencyDiscovery() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\033[1;31m❌ Panic: %v\033[0m\n", r)
		}
	}()

	content, err := os.ReadFile(ctx.MainFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;31m❌ Error reading file: %s\033[0m\n", err)
		printWaiting()
		return
	}

	// Set current module path for relative imports
	evaluator.CurrentModulePath = ctx.MainFile

	// Parse and discover imports BEFORE evaluation
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintf(os.Stderr, "  \033[1;31m❌ Parse Error: %s\033[0m\n", msg)
		}
		printWaiting()
		return
	}

	// Discover dependencies from the program
	dependencies := ctx.DiscoverDependencies(string(content), ctx.MainFile)
	for _, dep := range dependencies {
		ctx.AddWatchedFile(dep)
	}

	// Create environment
	env := object.NewEnvironment()
	setupGlobalEnv(env)

	// Add __filename and __dirname
	env.Set("__filename", &object.String{Value: ctx.MainFile})
	env.Set("__dirname", &object.String{Value: filepath.Dir(ctx.MainFile)})

	// Setup server tracker
	ctx.setupServerTracking(env)

	// Evaluate
	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintf(os.Stderr, "\033[1;31m❌ %s\033[0m\n", result.Inspect())
	}

	printWaiting()
}

// DiscoverDependencies finds all import/require dependencies recursively
func (ctx *WatchContext) DiscoverDependencies(content string, basePath string) []string {
	var dependencies []string
	visited := make(map[string]bool)

	ctx.discoverDepsRecursive(content, basePath, &dependencies, visited)

	return dependencies
}

// discoverDepsRecursive recursively discovers dependencies
func (ctx *WatchContext) discoverDepsRecursive(content string, basePath string, deps *[]string, visited map[string]bool) {
	baseDir := filepath.Dir(basePath)

	// Regular expressions to match import/require patterns
	importPatterns := []*regexp.Regexp{
		// ES6 imports: import ... from "..."
		regexp.MustCompile(`import\s+(?:[^"']*\s+from\s+)?["']([^"']+)["']`),
		// require: require("...")
		regexp.MustCompile(`require\s*\(\s*["']([^"']+)["']\s*\)`),
	}

	for _, pattern := range importPatterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			modulePath := match[1]

			// Skip built-in modules
			if isBuiltinModule(modulePath) {
				continue
			}

			// Resolve relative path
			var resolvedPath string
			if strings.HasPrefix(modulePath, ".") || strings.HasPrefix(modulePath, "/") {
				resolvedPath = filepath.Join(baseDir, modulePath)
				
				// Try with .nu extension if not present
				if filepath.Ext(resolvedPath) == "" {
					if _, err := os.Stat(resolvedPath + ".nu"); err == nil {
						resolvedPath = resolvedPath + ".nu"
					} else if _, err := os.Stat(filepath.Join(resolvedPath, "index.nu")); err == nil {
						resolvedPath = filepath.Join(resolvedPath, "index.nu")
					}
				}
			} else {
				// Check .nu_modules
				nuModulesPath := evaluator.FindNuModulesPath(basePath)
				if nuModulesPath != "" {
					pkgPath := filepath.Join(nuModulesPath, modulePath, "index.nu")
					if _, err := os.Stat(pkgPath); err == nil {
						resolvedPath = pkgPath
					}
				}
			}

			// Skip if already visited or doesn't exist
			if resolvedPath == "" {
				continue
			}

			absPath, err := filepath.Abs(resolvedPath)
			if err != nil {
				continue
			}

			if visited[absPath] {
				continue
			}

			if _, err := os.Stat(absPath); err != nil {
				continue
			}

			visited[absPath] = true
			*deps = append(*deps, absPath)

			// Recursively discover dependencies in this file
			depContent, err := os.ReadFile(absPath)
			if err == nil {
				ctx.discoverDepsRecursive(string(depContent), absPath, deps, visited)
			}
		}
	}
}

// isBuiltinModule checks if a module name is a built-in module
func isBuiltinModule(name string) bool {
	builtins := []string{"fs", "path", "crypto", "os", "http", "stream", "url", "querystring", "events"}
	for _, b := range builtins {
		if name == b {
			return true
		}
	}
	return false
}

// setupServerTracking injects server tracking into the environment
func (ctx *WatchContext) setupServerTracking(env *object.Environment) {
	// We'll track servers that are created
	// This is done by wrapping the http module's createServer
	originalCreateServer := evaluator.WrapHttpCreateServer(func(server *object.ObjectMap) {
		ctx.ServersMutex.Lock()
		ctx.ActiveServers = append(ctx.ActiveServers, server)
		ctx.ServersMutex.Unlock()
	})
	_ = originalCreateServer // Keep reference if needed
}

// printWatchedFiles prints the list of watched files
func printWatchedFiles(ctx *WatchContext) {
	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	count := len(ctx.WatchedFiles)
	fmt.Printf("\033[90m📂 Watching %d file(s):\033[0m\n", count)

	i := 0
	for filePath := range ctx.WatchedFiles {
		if i >= 5 { // Only show first 5 files
			fmt.Printf("\033[90m   ... and %d more\033[0m\n", count-5)
			break
		}
		rel, err := filepath.Rel(filepath.Dir(ctx.MainFile), filePath)
		if err != nil {
			rel = filepath.Base(filePath)
		}
		fmt.Printf("\033[90m   • %s\033[0m\n", rel)
		i++
	}
}

// clearScreen clears the terminal
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// printWaiting prints the waiting message
func printWaiting() {
	fmt.Println()
	fmt.Println("\033[90m" + strings.Repeat("─", 50) + "\033[0m")
	fmt.Println("\033[90m👀 Waiting for changes...\033[0m")
}
