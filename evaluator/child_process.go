package evaluator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/nulang/nulang/object"
)

// initChildProcessModule initializes the child_process module
func initChildProcessModule() *object.ObjectMap {
	cp := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// child_process.exec(command, options?, callback?)
	cp.Set("exec", &object.Builtin{Name: "exec", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("exec requires a command argument")
		}

		command := objectToString(args[0])
		var callback object.Object
		var options *object.ObjectMap

		// Parse arguments (command, options?, callback?)
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				options = opts
				if len(args) > 2 {
					callback = args[2]
				}
			} else {
				callback = args[1]
			}
		}

		// Parse options
		cwd := ""
		env := os.Environ()
		shell := true

		if options != nil {
			if cwdVal, ok := options.Get("cwd"); ok {
				cwd = objectToString(cwdVal)
			}
			if envVal, ok := options.Get("env"); ok {
				if envMap, ok := envVal.(*object.ObjectMap); ok {
					env = []string{}
					for key, pair := range envMap.Pairs {
						env = append(env, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
					}
				}
			}
		}


		// Execute command
		var cmd *exec.Cmd
		if shell {
			cmd = exec.Command("sh", "-c", command)
		} else {
			parts := strings.Fields(command)
			if len(parts) > 0 {
				cmd = exec.Command(parts[0], parts[1:]...)
			}
		}

		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		// Prepare callback arguments
		var errorObj object.Object = NULL
		if err != nil {
			errorObj = createErrorObject(err.Error())
		}

		stdoutObj := &object.String{Value: stdout.String()}
		stderrObj := &object.String{Value: stderr.String()}

		// Call callback if provided
		if callback != nil {
			if fn, ok := callback.(*object.Function); ok {
				fnEnv := extendFunctionEnv(fn, []object.Object{errorObj, stdoutObj, stderrObj})
				Eval(fn.Body, fnEnv)
			} else if builtin, ok := callback.(*object.Builtin); ok {
				builtin.Fn(errorObj, stdoutObj, stderrObj)
			}
		}

		return UNDEFINED
	}})

	// child_process.execSync(command, options?)
	cp.Set("execSync", &object.Builtin{Name: "execSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("execSync requires a command argument")
		}

		command := objectToString(args[0])
		var options *object.ObjectMap

		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				options = opts
			}
		}

		// Parse options
		cwd := ""
		env := os.Environ()
		encoding := "utf8"
		shell := true

		if options != nil {
			if cwdVal, ok := options.Get("cwd"); ok {
				cwd = objectToString(cwdVal)
			}
			if envVal, ok := options.Get("env"); ok {
				if envMap, ok := envVal.(*object.ObjectMap); ok {
					env = []string{}
					for key, pair := range envMap.Pairs {
						env = append(env, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
					}
				}
			}
			if encVal, ok := options.Get("encoding"); ok {
				encoding = objectToString(encVal)
			}
		}

		// Execute command
		var cmd *exec.Cmd
		if shell {
			cmd = exec.Command("sh", "-c", command)
		} else {
			parts := strings.Fields(command)
			if len(parts) > 0 {
				cmd = exec.Command(parts[0], parts[1:]...)
			}
		}

		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env

		output, err := cmd.CombinedOutput()
		if err != nil {
			return newError("execSync failed: %s", err.Error())
		}

		if encoding == "buffer" {
			return &object.Buffer{Data: output}
		}

		return &object.String{Value: string(output)}
	}})

	// child_process.spawn(command, args?, options?)
	cp.Set("spawn", &object.Builtin{Name: "spawn", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("spawn requires a command argument")
		}

		command := objectToString(args[0])
		var cmdArgs []string
		var options *object.ObjectMap

		// Parse arguments
		if len(args) > 1 {
			if arr, ok := args[1].(*object.Array); ok {
				for _, elem := range arr.Elements {
					cmdArgs = append(cmdArgs, objectToString(elem))
				}
				if len(args) > 2 {
					if opts, ok := args[2].(*object.ObjectMap); ok {
						options = opts
					}
				}
			} else if opts, ok := args[1].(*object.ObjectMap); ok {
				options = opts
			}
		}

		// Parse options
		cwd := ""
		env := os.Environ()

		if options != nil {
			if cwdVal, ok := options.Get("cwd"); ok {
				cwd = objectToString(cwdVal)
			}
			if envVal, ok := options.Get("env"); ok {
				if envMap, ok := envVal.(*object.ObjectMap); ok {
					env = []string{}
					for key, pair := range envMap.Pairs {
						env = append(env, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
					}
				}
			}
		}

		// Create command
		cmd := exec.Command(command, cmdArgs...)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env

		// Create child process object (EventEmitter)
		childProcess := createEventEmitter()

		// Setup stdio pipes
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		// Create stream objects
		stdinStream := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		stdinStream.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				data := objectToString(args[0])
				io.WriteString(stdin, data)
			}
			return TRUE
		}})
		stdinStream.Set("end", &object.Builtin{Name: "end", Fn: func(args ...object.Object) object.Object {
			stdin.Close()
			return UNDEFINED
		}})

		childProcess.Set("stdin", stdinStream)
		childProcess.Set("pid", &object.Number{Value: float64(cmd.Process.Pid)})

		// Start the process
		err := cmd.Start()
		if err != nil {
			return newError("spawn failed: %s", err.Error())
		}

		// Handle stdout
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if n > 0 {
					data := &object.String{Value: string(buf[:n])}
					emitEvent(childProcess, "data", data)
				}
				if err != nil {
					break
				}
			}
		}()

		// Handle stderr
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderr.Read(buf)
				if n > 0 {
					data := &object.String{Value: string(buf[:n])}
					emitEvent(childProcess, "error", data)
				}
				if err != nil {
					break
				}
			}
		}()

		// Handle process exit
		go func() {
			err := cmd.Wait()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}
			emitEvent(childProcess, "close", &object.Number{Value: float64(exitCode)})
		}()

		return childProcess
	}})

	// child_process.spawnSync(command, args?, options?)
	cp.Set("spawnSync", &object.Builtin{Name: "spawnSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("spawnSync requires a command argument")
		}

		command := objectToString(args[0])
		var cmdArgs []string
		var options *object.ObjectMap

		// Parse arguments
		if len(args) > 1 {
			if arr, ok := args[1].(*object.Array); ok {
				for _, elem := range arr.Elements {
					cmdArgs = append(cmdArgs, objectToString(elem))
				}
				if len(args) > 2 {
					if opts, ok := args[2].(*object.ObjectMap); ok {
						options = opts
					}
				}
			} else if opts, ok := args[1].(*object.ObjectMap); ok {
				options = opts
			}
		}

		// Parse options
		cwd := ""
		env := os.Environ()
		encoding := "utf8"

		if options != nil {
			if cwdVal, ok := options.Get("cwd"); ok {
				cwd = objectToString(cwdVal)
			}
			if envVal, ok := options.Get("env"); ok {
				if envMap, ok := envVal.(*object.ObjectMap); ok {
					env = []string{}
					for key, pair := range envMap.Pairs {
						env = append(env, fmt.Sprintf("%s=%s", key, objectToString(pair.Value)))
					}
				}
			}
			if encVal, ok := options.Get("encoding"); ok {
				encoding = objectToString(encVal)
			}
		}

		// Execute command
		cmd := exec.Command(command, cmdArgs...)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		// Create result object
		result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
		}

		result.Set("status", &object.Number{Value: float64(exitCode)})

		if encoding == "buffer" {
			result.Set("stdout", &object.Buffer{Data: stdout.Bytes()})
			result.Set("stderr", &object.Buffer{Data: stderr.Bytes()})
		} else {
			result.Set("stdout", &object.String{Value: stdout.String()})
			result.Set("stderr", &object.String{Value: stderr.String()})
		}

		if err != nil {
			result.Set("error", createErrorObject(err.Error()))
		} else {
			result.Set("error", NULL)
		}

		return result
	}})

	// child_process.fork(modulePath, args?, options?)
	cp.Set("fork", &object.Builtin{Name: "fork", Fn: func(args ...object.Object) object.Object {
		// Fork is more complex - simplified implementation
		return newError("fork is not yet fully implemented")
	}})

	return cp
}
