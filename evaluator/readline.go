package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nulang/nulang/object"
)

// initReadlineModule initializes the readline module
func initReadlineModule() *object.ObjectMap {
	readline := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// readline.createInterface(options)
	readline.Set("createInterface", &object.Builtin{Name: "createInterface", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("createInterface requires an options object")
		}

		options, ok := args[0].(*object.ObjectMap)
		if !ok {
			return newError("createInterface requires an options object")
		}

		// Create readline interface (EventEmitter)
		rl := createEventEmitter()

		reader := bufio.NewReader(os.Stdin)
		promptStr := "> "

		if ps, ok := options.Get("prompt"); ok {
			promptStr = objectToString(ps)
		}

		// Control channel and state for the readline interface
		closeChan := make(chan struct{})
		var closed bool
		var closeMutex sync.Mutex
		var paused bool
		var pauseMutex sync.Mutex

		rl.Set("prompt", &object.Builtin{Name: "prompt", Fn: func(args ...object.Object) object.Object {
			closeMutex.Lock()
			isClosed := closed
			closeMutex.Unlock()
			if isClosed {
				return UNDEFINED
			}
			fmt.Print(promptStr)
			return UNDEFINED
		}})

		rl.Set("question", &object.Builtin{Name: "question", Fn: func(args ...object.Object) object.Object {
			closeMutex.Lock()
			isClosed := closed
			closeMutex.Unlock()
			if isClosed {
				return &object.String{Value: ""}
			}

			if len(args) < 1 {
				return newError("question requires a query string")
			}

			query := objectToString(args[0])
			var callback object.Object
			if len(args) > 1 {
				callback = args[1]
			}

			fmt.Print(query)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)

			answerObj := &object.String{Value: answer}

			if callback != nil {
				if fn, ok := callback.(*object.Function); ok {
					fnEnv := extendFunctionEnv(fn, []object.Object{answerObj})
					Eval(fn.Body, fnEnv)
				} else if builtin, ok := callback.(*object.Builtin); ok {
					builtin.Fn(answerObj)
				}
			}

			return answerObj
		}})

		rl.Set("close", &object.Builtin{Name: "close", Fn: func(args ...object.Object) object.Object {
			closeMutex.Lock()
			if closed {
				closeMutex.Unlock()
				return UNDEFINED
			}
			closed = true
			closeMutex.Unlock()

			close(closeChan)
			emitEvent(rl, "close")
			UnregisterAsyncTask() // Allow program to exit
			return UNDEFINED
		}})

		rl.Set("pause", &object.Builtin{Name: "pause", Fn: func(args ...object.Object) object.Object {
			pauseMutex.Lock()
			paused = true
			pauseMutex.Unlock()
			emitEvent(rl, "pause")
			return UNDEFINED
		}})

		rl.Set("resume", &object.Builtin{Name: "resume", Fn: func(args ...object.Object) object.Object {
			pauseMutex.Lock()
			paused = false
			pauseMutex.Unlock()
			emitEvent(rl, "resume")
			return UNDEFINED
		}})

		rl.Set("setPrompt", &object.Builtin{Name: "setPrompt", Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				promptStr = objectToString(args[0])
			}
			return UNDEFINED
		}})

		rl.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
			if len(args) > 0 {
				data := objectToString(args[0])
				fmt.Print(data)
			}
			return UNDEFINED
		}})

		// Register async task to keep program running
		RegisterAsyncTask()

		// Start reading lines in background
		go func() {
			for {
				// Check if closed
				closeMutex.Lock()
				isClosed := closed
				closeMutex.Unlock()
				if isClosed {
					return
				}

				// Check if paused
				pauseMutex.Lock()
				isPaused := paused
				pauseMutex.Unlock()
				if isPaused {
					// Wait a bit before checking again
					select {
					case <-closeChan:
						return
					default:
						continue
					}
				}

				line, err := reader.ReadString('\n')
				if err != nil {
					// Check if interface was closed intentionally
					closeMutex.Lock()
					isClosed := closed
					closeMutex.Unlock()
					if !isClosed {
						closeMutex.Lock()
						closed = true
						closeMutex.Unlock()
						close(closeChan)
						emitEvent(rl, "close")
						UnregisterAsyncTask() // Allow program to exit
					}
					return
				}

				line = strings.TrimSpace(line)
				emitEvent(rl, "line", &object.String{Value: line})
			}
		}()

		return rl
	}})

	// readline.clearLine(stream, dir, callback?)
	readline.Set("clearLine", &object.Builtin{Name: "clearLine", Fn: func(args ...object.Object) object.Object {
		// ANSI escape to clear line
		fmt.Print("\033[2K")
		if len(args) > 2 {
			if callback, ok := args[2].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{})
				Eval(callback.Body, fnEnv)
			}
		}
		return UNDEFINED
	}})

	// readline.clearScreenDown(stream, callback?)
	readline.Set("clearScreenDown", &object.Builtin{Name: "clearScreenDown", Fn: func(args ...object.Object) object.Object {
		// ANSI escape to clear screen from cursor down
		fmt.Print("\033[J")
		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{})
				Eval(callback.Body, fnEnv)
			}
		}
		return UNDEFINED
	}})

	// readline.cursorTo(stream, x, y?, callback?)
	readline.Set("cursorTo", &object.Builtin{Name: "cursorTo", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("cursorTo requires stream and x coordinate")
		}

		x := 0
		y := 0

		if num, ok := args[1].(*object.Number); ok {
			x = int(num.Value)
		}

		if len(args) > 2 {
			if num, ok := args[2].(*object.Number); ok {
				y = int(num.Value)
			}
		}

		// ANSI escape to move cursor
		if y > 0 {
			fmt.Printf("\033[%d;%dH", y+1, x+1)
		} else {
			fmt.Printf("\033[%dG", x+1)
		}

		return UNDEFINED
	}})

	// readline.moveCursor(stream, dx, dy, callback?)
	readline.Set("moveCursor", &object.Builtin{Name: "moveCursor", Fn: func(args ...object.Object) object.Object {
		if len(args) < 3 {
			return newError("moveCursor requires stream, dx, and dy")
		}

		dx := 0
		dy := 0

		if num, ok := args[1].(*object.Number); ok {
			dx = int(num.Value)
		}
		if num, ok := args[2].(*object.Number); ok {
			dy = int(num.Value)
		}

		// ANSI escape to move cursor
		if dx > 0 {
			fmt.Printf("\033[%dC", dx)
		} else if dx < 0 {
			fmt.Printf("\033[%dD", -dx)
		}

		if dy > 0 {
			fmt.Printf("\033[%dB", dy)
		} else if dy < 0 {
			fmt.Printf("\033[%dA", -dy)
		}

		return UNDEFINED
	}})

	return readline
}
