package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nulang/nulang/evaluator"
	"github.com/nulang/nulang/lexer"
	"github.com/nulang/nulang/object"
	"github.com/nulang/nulang/parser"
)

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		startREPL()
		return
	}

	// Check for --watch flag
	watchMode := false
	filename := ""
	
	for i, arg := range os.Args[1:] {
		if arg == "--watch" || arg == "-w" {
			watchMode = true
		} else if arg != "" && filename == "" && !isCommand(arg) {
			filename = os.Args[i+1]
		}
	}

	// Handle commands
	if len(os.Args) > 1 && isCommand(os.Args[1]) {
		switch os.Args[1] {
		case "install":
			handleInstall()
		case "init":
			handleInit()
		case "version", "-v", "--version":
			fmt.Printf("Nulang v%s\n", VERSION)
		case "help", "-h", "--help":
			printCLIHelp()
		}
		return
	}

	// Run file (with or without watch)
	if filename != "" {
		if watchMode {
			runWithWatch(filename)
		} else {
			runFile(filename)
		}
	} else if len(os.Args) > 1 && !isCommand(os.Args[1]) {
		if watchMode {
			runWithWatch(os.Args[1])
		} else {
			runFile(os.Args[1])
		}
	}
}

func isCommand(arg string) bool {
	commands := []string{"install", "init", "version", "-v", "--version", "help", "-h", "--help", "--watch", "-w"}
	for _, cmd := range commands {
		if arg == cmd {
			return true
		}
	}
	return false
}

func runWithWatch(filename string) {
	fmt.Printf("\033[1;36mWatch mode enabled for: %s\033[0m\n", filename)
	fmt.Println("\033[1;33mPress Ctrl+C to stop\033[0m")
	fmt.Println()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Get initial file info
	lastModTime, err := getModTime(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;31mError: %s\033[0m\n", err)
		os.Exit(1)
	}

	// Run initially
	runFileWithRecover(filename)

	// Watch loop
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\033[1;33mWatch mode stopped\033[0m")
			return
		case <-ticker.C:
			currentModTime, err := getModTime(filename)
			if err != nil {
				continue
			}

			if currentModTime.After(lastModTime) {
				lastModTime = currentModTime
				clearScreen()
				fmt.Printf("\033[1;32mFile changed at %s\033[0m\n", time.Now().Format("15:04:05"))
				fmt.Println("\033[1;90m" + "─────────────────────────────────────" + "\033[0m")
				runFileWithRecover(filename)
			}
		}
	}
}

func getModTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func runFileWithRecover(filename string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\033[1;31mPanic recovered: %v\033[0m\n", r)
		}
	}()

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;31mError reading file: %s\033[0m\n", err)
		return
	}

	// Set current module path for relative imports
	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}
	evaluator.CurrentModulePath = absPath

	env := object.NewEnvironment()
	setupGlobalEnv(env)
	
	// Add __filename and __dirname
	env.Set("__filename", &object.String{Value: absPath})
	env.Set("__dirname", &object.String{Value: filepath.Dir(absPath)})

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintf(os.Stderr, "\033[1;31m%s\033[0m\n", result.Inspect())
	}

	fmt.Println("\033[1;90m" + "─────────────────────────────────────" + "\033[0m")
	fmt.Println("\033[1;90mWaiting for changes...\033[0m")
}

func handleInstall() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if err := evaluator.InstallDependencies(cwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func handleInit() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if err := evaluator.InitProject(cwd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func printCLIHelp() {
	fmt.Printf("Nulang v%s - JavaScript-like language written in Go\n\n", VERSION)
	fmt.Println("Usage:")
	fmt.Println("  nulang <file.nu>           Run a Nulang script")
	fmt.Println("  nulang <file.nu> --watch   Run with hot reload on changes")
	fmt.Println("  nulang                     Start the REPL")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install              Install dependencies from nulang.yml")
	fmt.Println("  init                 Create a new nulang.yml file")
	fmt.Println("  version, -v          Show version")
	fmt.Println("  help, -h             Show this help")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --watch, -w          Watch for file changes and auto-reload")
	fmt.Println()
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Set current module path for relative imports
	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}
	evaluator.CurrentModulePath = absPath

	env := object.NewEnvironment()
	setupGlobalEnv(env)
	
	// Add __filename and __dirname
	env.Set("__filename", &object.String{Value: absPath})
	env.Set("__dirname", &object.String{Value: filepath.Dir(absPath)})

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	result := evaluator.Eval(program, env)

	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintf(os.Stderr, "%s\n", result.Inspect())
		os.Exit(1)
	}

	// Wait for any async tasks (servers, timers)
	evaluator.AwaitAsyncTasks()
}

func startREPL() {
	fmt.Printf("Nulang v%s - JavaScript-like language written in Go\n", VERSION)
	fmt.Println("Type 'exit' to quit, 'help' for more info")
	fmt.Println()

	env := object.NewEnvironment()
	setupGlobalEnv(env)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("nulang> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}
		if line == "help" {
			printHelp()
			continue
		}
		if line == "" {
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(os.Stdout, p.Errors())
			continue
		}

		result := evaluator.Eval(program, env)

		if result != nil {
			if result.Type() != object.UNDEFINED_OBJ {
				fmt.Println(result.Inspect())
			}
		}
	}
}

func setupGlobalEnv(env *object.Environment) {
	// globalThis
	env.Set("globalThis", &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)})
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "  Parse Error: %s\n", msg)
	}
}

func printHelp() {
	fmt.Println("Nulang Help:")
	fmt.Println("  - Syntax is similar to JavaScript/Node.js")
	fmt.Println("  - Supports: let, const, var, function, arrow functions")
	fmt.Println("  - Control flow: if/else, for, while, try/catch")
	fmt.Println("  - Built-ins: console.log, Math, JSON, Array, Object")
	fmt.Println("  - Run file: nulang filename.nu")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  let x = 5")
	fmt.Println("  const arr = [1, 2, 3]")
	fmt.Println("  arr.map(x => x * 2)")
	fmt.Println("  console.log('Hello, World!')")
	fmt.Println()
}
