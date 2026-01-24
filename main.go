package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		} else if arg != "" && filename == "" && !isCommand(arg) && !strings.HasPrefix(arg, "-") {
			filename = os.Args[i+1]
		}
	}

	// Handle commands
	if len(os.Args) > 1 && isCommand(os.Args[1]) {
		switch os.Args[1] {
		case "install", "i", "add":
			handleInstallNew()
		case "uninstall", "remove", "rm":
			handleUninstall()
		case "init":
			handleInitNew()
		case "list", "ls":
			handleList()
		case "run":
			handleRun()
		case "update", "upgrade":
			handleUpdate()
		case "prune":
			handlePrune()
		case "outdated":
			handleOutdated()
		case "version", "-v", "--version":
			fmt.Printf("Nu v%s\n", VERSION)
		case "help", "-h", "--help":
			printCLIHelp()
		}
		return
	}

	// Run file (with or without watch)
	if filename != "" {
		if watchMode {
			RunWithWatchV2(filename)
		} else {
			runFile(filename)
		}
	} else if len(os.Args) > 1 && !isCommand(os.Args[1]) {
		if watchMode {
			RunWithWatchV2(os.Args[1])
		} else {
			runFile(os.Args[1])
		}
	}
}

func isCommand(arg string) bool {
	commands := []string{
		"install", "i", "add",
		"uninstall", "remove", "rm",
		"init",
		"list", "ls",
		"run",
		"update", "upgrade",
		"prune",
		"outdated",
		"version", "-v", "--version",
		"help", "-h", "--help",
		"--watch", "-w",
	}
	for _, cmd := range commands {
		if arg == cmd {
			return true
		}
	}
	return false
}

// handleInstallNew handles install command with npm-compatible package.json
func handleInstallNew() {
	// Parse flags and packages
	flags := make(map[string]bool)
	var packages []string

	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "--") {
			flags[strings.TrimPrefix(arg, "--")] = true
		} else if strings.HasPrefix(arg, "-") {
			flags[strings.TrimPrefix(arg, "-")] = true
		} else {
			packages = append(packages, arg)
		}
	}

	if err := evaluator.HandleInstallCommand(packages, flags); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handleUninstall handles uninstall command
func handleUninstall() {
	packages := os.Args[2:]
	if err := evaluator.HandleUninstallCommand(packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handleInitNew handles init command with npm-compatible package.json
func handleInitNew() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Check for -y flag for non-interactive
	interactive := true
	for _, arg := range os.Args[2:] {
		if arg == "-y" || arg == "--yes" {
			interactive = false
			break
		}
	}

	var initErr error
	if interactive {
		initErr = evaluator.InitProjectJSONInteractive(cwd)
	} else {
		initErr = evaluator.InitProjectJSON(cwd)
	}

	if initErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", initErr)
		os.Exit(1)
	}
}

// handleList handles list command
func handleList() {
	if err := evaluator.HandleListCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handleRun handles run command
func handleRun() {
	scriptName := ""
	var args []string

	if len(os.Args) > 2 {
		scriptName = os.Args[2]
		if len(os.Args) > 3 {
			args = os.Args[3:]
		}
	}

	if err := evaluator.HandleRunCommand(scriptName, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handleUpdate handles update command
func handleUpdate() {
	packages := os.Args[2:]
	if err := evaluator.HandleUpdateCommand(packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handlePrune handles prune command
func handlePrune() {
	if err := evaluator.HandlePruneCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// handleOutdated handles outdated command
func handleOutdated() {
	if err := evaluator.HandleOutdatedCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func printCLIHelp() {
	fmt.Printf("Nu v%s - JavaScript-like language written in Go\n\n", VERSION)
	fmt.Println("Usage:")
	fmt.Println("  nu <file.js>           Run a script (.js or .ts)")
	fmt.Println("  nu <file.ts> --watch   Run with hot reload on changes")
	fmt.Println("  nu                     Start the REPL")
	fmt.Println()
	fmt.Println("Package Management (npm-compatible):")
	fmt.Println("  init                   Create package.json")
	fmt.Println("  init -y                Create package.json with defaults")
	fmt.Println("  install, i             Install all dependencies")
	fmt.Println("  install <pkg>          Install a package (adds to package.json)")
	fmt.Println("  install <pkg> -D       Install as devDependency")
	fmt.Println("  uninstall <pkg>        Remove a package")
	fmt.Println("  update [pkg...]        Update packages")
	fmt.Println("  list, ls               List installed packages")
	fmt.Println("  outdated               Show outdated packages")
	fmt.Println("  prune                  Remove extraneous packages")
	fmt.Println("  run <script>           Run a script from package.json")
	fmt.Println()
	fmt.Println("Other Commands:")
	fmt.Println("  version, -v            Show version")
	fmt.Println("  help, -h               Show this help")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --watch, -w            Watch for file changes and auto-reload")
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
	fmt.Printf("Nu v%s\n", VERSION)
	fmt.Println("Type 'exit' to quit, 'help' for more info")
	fmt.Println()

	env := object.NewEnvironment()
	setupGlobalEnv(env)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("nu> ")
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
	fmt.Println("Nu Help:")
	fmt.Println("  - Syntax is similar to JavaScript/Node.js")
	fmt.Println("  - Supports: let, const, var, function, arrow functions")
	fmt.Println("  - Control flow: if/else, for, while, try/catch")
	fmt.Println("  - Built-ins: console.log, Math, JSON, Array, Object")
	fmt.Println("  - Run file: nu filename.js (or .ts)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  let x = 5")
	fmt.Println("  const arr = [1, 2, 3]")
	fmt.Println("  arr.map(x => x * 2)")
	fmt.Println("  console.log('Hello, World!')")
	fmt.Println()
}
