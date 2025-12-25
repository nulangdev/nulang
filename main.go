package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	filename := os.Args[1]
	runFile(filename)
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
	// Process object
	processObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	
	// process.argv
	argv := &object.Array{Elements: make([]object.Object, len(os.Args))}
	for i, arg := range os.Args {
		argv.Elements[i] = &object.String{Value: arg}
	}
	processObj.Set("argv", argv)
	
	// process.env
	envMap := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	for _, e := range os.Environ() {
		for j := 0; j < len(e); j++ {
			if e[j] == '=' {
				key := e[:j]
				value := e[j+1:]
				envMap.Set(key, &object.String{Value: value})
				break
			}
		}
	}
	processObj.Set("env", envMap)
	
	// process.exit
	processObj.Set("exit", &object.Builtin{
		Name: "exit",
		Fn: func(args ...object.Object) object.Object {
			code := 0
			if len(args) > 0 {
				if num, ok := args[0].(*object.Number); ok {
					code = int(num.Value)
				}
			}
			os.Exit(code)
			return nil
		},
	})
	
	// process.cwd
	processObj.Set("cwd", &object.Builtin{
		Name: "cwd",
		Fn: func(args ...object.Object) object.Object {
			cwd, _ := os.Getwd()
			return &object.String{Value: cwd}
		},
	})
	
	env.Set("process", processObj)
	
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
