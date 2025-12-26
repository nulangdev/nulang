package evaluator

import (
	"os"
	"strings"

	"github.com/nulang/nulang/object"
)

// initProcessObject creates the global 'process' object
func initProcessObject() *object.ObjectMap {
	proc := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// process.argv
	argv := &object.Array{Elements: []object.Object{}}
	for _, arg := range os.Args {
		argv.Elements = append(argv.Elements, &object.String{Value: arg})
	}
	proc.Set("argv", argv)

	// process.env
	envMap := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			envMap.Set(pair[0], &object.String{Value: pair[1]})
		}
	}
	proc.Set("env", envMap)

	// process.platform
	proc.Set("platform", &object.String{Value: "darwin"}) // Hardcoded for now based on user context, ideally runtime check

	// process.cwd()
	proc.Set("cwd", &object.Builtin{Name: "cwd", Fn: func(args ...object.Object) object.Object {
		dir, err := os.Getwd()
		if err != nil {
			return newError("failed to get cwd: %s", err.Error())
		}
		return &object.String{Value: dir}
	}})

	// process.chdir(directory)
	proc.Set("chdir", &object.Builtin{Name: "chdir", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("chdir requires a directory path")
		}
		dir := objectToString(args[0])
		err := os.Chdir(dir)
		if err != nil {
			return newError("failed to chdir: %s", err.Error())
		}
		return UNDEFINED
	}})

	// process.exit(code)
	proc.Set("exit", &object.Builtin{Name: "exit", Fn: func(args ...object.Object) object.Object {
		code := 0
		if len(args) > 0 {
			if num, ok := args[0].(*object.Number); ok {
				code = int(num.Value)
			}
		}
		os.Exit(code)
		return UNDEFINED
	}})

	// process.nextTick(callback, ...args)
	proc.Set("nextTick", &object.Builtin{Name: "nextTick", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("nextTick requires a callback")
		}

		fn, ok := args[0].(*object.Function)
		if !ok {
			return newError("nextTick first argument must be a function")
		}

		callbackArgs := args[1:]

		// Use the same mechanism as setImmediate/timers
		go func() {
			executeTimerCallback(fn, fn.Env, callbackArgs)
		}()

		return UNDEFINED
	}})
	
	// process.stdout (basic wrapper)
	stdout := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	stdout.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
		for _, arg := range args {
			os.Stdout.WriteString(objectToString(arg))
		}
		return TRUE
	}})
	proc.Set("stdout", stdout)

	// process.stderr (basic wrapper)
	stderr := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
	stderr.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
		for _, arg := range args {
			os.Stderr.WriteString(objectToString(arg))
		}
		return TRUE
	}})
	proc.Set("stderr", stderr)

	return proc
}
