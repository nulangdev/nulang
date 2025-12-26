package evaluator

import (
	"fmt"
	"net"
	"time"

	"github.com/nulang/nulang/object"
)

// initNetModule initializes the net module
func initNetModule() *object.ObjectMap {
	netMod := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// net.createServer(options?, connectionListener?)
	netMod.Set("createServer", &object.Builtin{Name: "createServer", Fn: func(args ...object.Object) object.Object {
		server := createEventEmitter()
		
		var connectionListener object.Object
		
		// Parse arguments
		if len(args) > 0 {
			if fn, ok := args[0].(*object.Function); ok {
				connectionListener = fn
			} else if builtin, ok := args[0].(*object.Builtin); ok {
				connectionListener = builtin
			} else if _, ok := args[0].(*object.ObjectMap); ok {
				// Options object
				if len(args) > 1 {
					connectionListener = args[1]
				}
			}
		}

		var listener net.Listener

		// server.listen(port, host?, callback?)
		server.Set("listen", &object.Builtin{Name: "listen", Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("listen requires a port argument")
			}

			port := int(args[0].(*object.Number).Value)
			host := "0.0.0.0"
			var callback object.Object

			if len(args) > 1 {
				if str, ok := args[1].(*object.String); ok {
					host = str.Value
					if len(args) > 2 {
						callback = args[2]
					}
				} else {
					callback = args[1]
				}
			}

			address := fmt.Sprintf("%s:%d", host, port)
			ln, err := net.Listen("tcp", address)
			if err != nil {
				return newError("Failed to listen: %s", err.Error())
			}

			listener = ln

			// Call callback if provided
			if callback != nil {
				if fn, ok := callback.(*object.Function); ok {
					fnEnv := extendFunctionEnv(fn, []object.Object{})
					Eval(fn.Body, fnEnv)
				} else if builtin, ok := callback.(*object.Builtin); ok {
					builtin.Fn()
				}
			}

			emitEvent(server, "listening")

			// Accept connections in background
			RegisterAsyncTask()
			go func() {
				defer UnregisterAsyncTask()
				for {
					conn, err := listener.Accept()
					if err != nil {
						// Only emit error if not closed intentionally
						if err.Error() != "use of closed network connection" {
							emitEvent(server, "error", createErrorObject(err.Error()))
						}
						break
					}

					// Create socket object WITHOUT starting read
					socket := newSocketObject(conn)

					// Register connection listener before emitting event
					// But emitEvent executes listeners synchronously, so we are fine.
					if connectionListener != nil {
						addEventListener(server, "connection", connectionListener)
					}
					
					// Emit 'connection' - this will run connectionListener(socket) synchronously
					// Inside connectionListener, user does socket.on('data', ...)
					emitEvent(server, "connection", socket)

					// NOW it is safe to start reading
					startSocketRead(socket, conn)
				}
			}()

			return server
		}})

		// server.close(callback?)
		server.Set("close", &object.Builtin{Name: "close", Fn: func(args ...object.Object) object.Object {
			if listener != nil {
				listener.Close()
				emitEvent(server, "close")
			}

			if len(args) > 0 {
				if callback, ok := args[0].(*object.Function); ok {
					fnEnv := extendFunctionEnv(callback, []object.Object{})
					Eval(callback.Body, fnEnv)
				}
			}

			return UNDEFINED
		}})

		// server.address()
		server.Set("address", &object.Builtin{Name: "address", Fn: func(args ...object.Object) object.Object {
			if listener != nil {
				addr := listener.Addr()
				result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
				
				if tcpAddr, ok := addr.(*net.TCPAddr); ok {
					result.Set("address", &object.String{Value: tcpAddr.IP.String()})
					result.Set("port", &object.Number{Value: float64(tcpAddr.Port)})
					result.Set("family", &object.String{Value: "IPv4"})
				}
				
				return result
			}
			return NULL
		}})

		return server
	}})

	// net.connect(options, connectionListener?) or net.connect(port, host?, connectionListener?)
	netMod.Set("connect", &object.Builtin{Name: "connect", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("connect requires at least one argument")
		}

		var port int
		var host string = "localhost"
		var connectionListener object.Object

		// Parse arguments
		if opts, ok := args[0].(*object.ObjectMap); ok {
			// Options object
			if p, ok := opts.Get("port"); ok {
				if num, ok := p.(*object.Number); ok {
					port = int(num.Value)
				}
			}
			if h, ok := opts.Get("host"); ok {
				host = objectToString(h)
			}
			if len(args) > 1 {
				connectionListener = args[1]
			}
		} else if num, ok := args[0].(*object.Number); ok {
			// port, host?, callback?
			port = int(num.Value)
			if len(args) > 1 {
				if str, ok := args[1].(*object.String); ok {
					host = str.Value
					if len(args) > 2 {
						connectionListener = args[2]
					}
				} else {
					connectionListener = args[1]
				}
			}
		}

		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return newError("Failed to connect: %s", err.Error())
		}

		// Create socket WITHOUT starting read
		socket := newSocketObject(conn)

		// Register listener if provided (synchronously)
		if connectionListener != nil {
			addEventListener(socket, "connect", connectionListener)
		}

		// Emit 'connect' and start reading asynchronously
		// This ensures 'net.connect' returns the socket BEFORE the callback runs,
		// allowing usage like: var client = net.connect(..., function() { client.write(...) });
		RegisterAsyncTask()
		go func() {
			defer UnregisterAsyncTask()
			// Small yield to ensure main thread has assigned the variable
			time.Sleep(10 * time.Millisecond)
			
			emitEvent(socket, "connect")

			// NOW start reading
			startSocketRead(socket, conn)
		}()

		return socket
	}})

	// Alias
	netMod.Set("createConnection", netMod.Pairs["connect"].Value)

	return netMod
}

// newSocketObject creates a socket object (extends EventEmitter) WITHOUT starting read loop
func newSocketObject(conn net.Conn) *object.ObjectMap {
	socket := createEventEmitter()

	// socket.write(data, encoding?, callback?)
	socket.Set("write", &object.Builtin{Name: "write", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}

		var data []byte
		switch d := args[0].(type) {
		case *object.String:
			data = []byte(d.Value)
		case *object.Buffer:
			data = d.Data
		default:
			data = []byte(objectToString(args[0]))
		}

		_, err := conn.Write(data)
		if err != nil {
			emitEvent(socket, "error", createErrorObject(err.Error()))
			return FALSE
		}

		if len(args) > 1 {
			if callback, ok := args[len(args)-1].(*object.Function); ok {
				fnEnv := extendFunctionEnv(callback, []object.Object{})
				Eval(callback.Body, fnEnv)
			}
		}

		return TRUE
	}})

	// socket.end(data?, encoding?, callback?)
	socket.Set("end", &object.Builtin{Name: "end", Fn: func(args ...object.Object) object.Object {
		if len(args) > 0 {
			if str, ok := args[0].(*object.String); ok {
				conn.Write([]byte(str.Value))
			}
		}
		conn.Close()
		emitEvent(socket, "end")
		return UNDEFINED
	}})

	// socket.destroy()
	socket.Set("destroy", &object.Builtin{Name: "destroy", Fn: func(args ...object.Object) object.Object {
		conn.Close()
		emitEvent(socket, "close")
		return UNDEFINED
	}})

	// socket.pause()
	socket.Set("pause", &object.Builtin{Name: "pause", Fn: func(args ...object.Object) object.Object {
		// In a real implementation, this would pause reading
		return socket
	}})

	// socket.resume()
	socket.Set("resume", &object.Builtin{Name: "resume", Fn: func(args ...object.Object) object.Object {
		// In a real implementation, this would resume reading
		return socket
	}})

	// socket.setTimeout(timeout, callback?)
	socket.Set("setTimeout", &object.Builtin{Name: "setTimeout", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return socket
		}

		timeout := time.Duration(args[0].(*object.Number).Value) * time.Millisecond
		conn.SetDeadline(time.Now().Add(timeout))

		if len(args) > 1 {
			if callback, ok := args[1].(*object.Function); ok {
				addEventListener(socket, "timeout", callback)
			}
		}

		return socket
	}})

	// socket.setKeepAlive(enable?, initialDelay?)
	socket.Set("setKeepAlive", &object.Builtin{Name: "setKeepAlive", Fn: func(args ...object.Object) object.Object {
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			enable := true
			if len(args) > 0 {
				if b, ok := args[0].(*object.Boolean); ok {
					enable = b.Value
				}
			}
			tcpConn.SetKeepAlive(enable)
		}
		return socket
	}})

	// socket.setNoDelay(noDelay?)
	socket.Set("setNoDelay", &object.Builtin{Name: "setNoDelay", Fn: func(args ...object.Object) object.Object {
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			noDelay := true
			if len(args) > 0 {
				if b, ok := args[0].(*object.Boolean); ok {
					noDelay = b.Value
				}
			}
			tcpConn.SetNoDelay(noDelay)
		}
		return socket
	}})

	// Properties
	if tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		socket.Set("remoteAddress", &object.String{Value: tcpAddr.IP.String()})
		socket.Set("remotePort", &object.Number{Value: float64(tcpAddr.Port)})
	}

	if tcpAddr, ok := conn.LocalAddr().(*net.TCPAddr); ok {
		socket.Set("localAddress", &object.String{Value: tcpAddr.IP.String()})
		socket.Set("localPort", &object.Number{Value: float64(tcpAddr.Port)})
	}

	return socket
}

// startSocketRead starts reading data from connection in background
func startSocketRead(socket *object.ObjectMap, conn net.Conn) {
	RegisterAsyncTask()
	go func() {
		defer UnregisterAsyncTask()
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err.Error() != "EOF" && err.Error() != "use of closed network connection" {
					emitEvent(socket, "error", createErrorObject(err.Error()))
				}
				emitEvent(socket, "end")
				emitEvent(socket, "close")
				break
			}

			if n > 0 {
				data := &object.Buffer{Data: make([]byte, n)}
				copy(data.Data, buf[:n])
				emitEvent(socket, "data", data)
			}
		}
	}()
}
