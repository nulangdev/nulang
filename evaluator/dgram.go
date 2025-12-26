package evaluator

import (
	"fmt"
	"net"

	"github.com/nulang/nulang/object"
)

// initDgramModule initializes the dgram module (UDP)
func initDgramModule() *object.ObjectMap {
	dgram := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// dgram.createSocket(type, callback?)
	dgram.Set("createSocket", &object.Builtin{Name: "createSocket", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("createSocket requires a type argument")
		}

		socketType := objectToString(args[0])
		var callback object.Object

		if len(args) > 1 {
			callback = args[1]
		}

		// Create UDP socket (EventEmitter)
		socket := createEventEmitter()

		var conn *net.UDPConn
		network := "udp4"

		if socketType == "udp6" {
			network = "udp6"
		}

		// socket.bind(port?, address?, callback?)
		socket.Set("bind", &object.Builtin{Name: "bind", Fn: func(args ...object.Object) object.Object {
			port := 0
			address := "0.0.0.0"
			var bindCallback object.Object

			// Parse arguments (very flexible in Node.js)
			for _, arg := range args {
				switch v := arg.(type) {
				case *object.Number:
					port = int(v.Value)
				case *object.String:
					address = v.Value
				case *object.Function, *object.Builtin:
					bindCallback = arg
				case *object.ObjectMap:
					if p, ok := v.Get("port"); ok {
						if num, ok := p.(*object.Number); ok {
							port = int(num.Value)
						}
					}
					if a, ok := v.Get("address"); ok {
						address = objectToString(a)
					}
				}
			}

			udpAddr, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s:%d", address, port))
			if err != nil {
				return newError("Failed to resolve UDP address: %s", err.Error())
			}

			udpConn, err := net.ListenUDP(network, udpAddr)
			if err != nil {
				return newError("Failed to bind: %s", err.Error())
			}

			conn = udpConn

			if bindCallback != nil {
				callFunction(bindCallback, []object.Object{})
			}

			emitEvent(socket, "listening")

			// Start receiving messages
			RegisterAsyncTask()
			go func() {
				defer UnregisterAsyncTask()
				buf := make([]byte, 65536)
				for {
					n, remoteAddr, err := conn.ReadFromUDP(buf)
					if err != nil {
						emitEvent(socket, "error", createErrorObject(err.Error()))
						break
					}

					if n > 0 {
						data := &object.Buffer{Data: make([]byte, n)}
						copy(data.Data, buf[:n])

						rinfo := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
						rinfo.Set("address", &object.String{Value: remoteAddr.IP.String()})
						rinfo.Set("port", &object.Number{Value: float64(remoteAddr.Port)})
						rinfo.Set("family", &object.String{Value: "IPv4"})

						emitEvent(socket, "message", data, rinfo)
					}
				}
			}()

			return UNDEFINED
		}})

		// socket.send(msg, offset?, length?, port, address?, callback?)
		socket.Set("send", &object.Builtin{Name: "send", Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("send requires at least message and port")
			}

			var msg []byte
			offset := 0
			length := -1
			var port int
			address := "127.0.0.1"
			var sendCallback object.Object

			// Parse message
			switch m := args[0].(type) {
			case *object.String:
				msg = []byte(m.Value)
			case *object.Buffer:
				msg = m.Data
			case *object.Array:
				// Array of buffers or strings
				for _, elem := range m.Elements {
					if str, ok := elem.(*object.String); ok {
						msg = append(msg, []byte(str.Value)...)
					} else if buf, ok := elem.(*object.Buffer); ok {
						msg = append(msg, buf.Data...)
					}
				}
			}

			// Parse remaining arguments
			argIdx := 1
			if len(args) > argIdx {
				if num, ok := args[argIdx].(*object.Number); ok {
					// Could be offset or port
					if len(args) > argIdx+2 {
						// It's offset
						offset = int(num.Value)
						argIdx++
						if num2, ok := args[argIdx].(*object.Number); ok {
							length = int(num2.Value)
							argIdx++
						}
					} else {
						// It's port
						port = int(num.Value)
						argIdx++
					}
				}
			}

			// Get port if not set yet
			if port == 0 && len(args) > argIdx {
				if num, ok := args[argIdx].(*object.Number); ok {
					port = int(num.Value)
					argIdx++
				}
			}

			// Get address
			if len(args) > argIdx {
				if str, ok := args[argIdx].(*object.String); ok {
					address = str.Value
					argIdx++
				}
			}

			// Get callback
			if len(args) > argIdx {
				sendCallback = args[argIdx]
			}

			// Apply offset and length
			if length < 0 {
				length = len(msg) - offset
			}
			if offset > 0 || length < len(msg) {
				msg = msg[offset : offset+length]
			}

			// Send message
			udpAddr, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s:%d", address, port))
			if err != nil {
				errorObj := createErrorObject(err.Error())
				if sendCallback != nil {
					callFunction(sendCallback, []object.Object{errorObj})
				}
				return errorObj
			}

			if conn == nil {
				// Create a temporary connection for sending
				tempConn, err := net.DialUDP(network, nil, udpAddr)
				if err != nil {
					errorObj := createErrorObject(err.Error())
					if sendCallback != nil {
						callFunction(sendCallback, []object.Object{errorObj})
					}
					return errorObj
				}
				defer tempConn.Close()
				conn = tempConn
			}

			_, err = conn.WriteToUDP(msg, udpAddr)
			if err != nil {
				errorObj := createErrorObject(err.Error())
				if sendCallback != nil {
					callFunction(sendCallback, []object.Object{errorObj})
				}
				return errorObj
			}

			if sendCallback != nil {
				callFunction(sendCallback, []object.Object{NULL, &object.Number{Value: float64(len(msg))}})
			}

			return UNDEFINED
		}})

		// socket.close(callback?)
		socket.Set("close", &object.Builtin{Name: "close", Fn: func(args ...object.Object) object.Object {
			if conn != nil {
				conn.Close()
				emitEvent(socket, "close")
			}

			if len(args) > 0 {
				callFunction(args[0], []object.Object{})
			}

			return UNDEFINED
		}})

		// socket.address()
		socket.Set("address", &object.Builtin{Name: "address", Fn: func(args ...object.Object) object.Object {
			if conn != nil {
				addr := conn.LocalAddr()
				result := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

				if udpAddr, ok := addr.(*net.UDPAddr); ok {
					result.Set("address", &object.String{Value: udpAddr.IP.String()})
					result.Set("port", &object.Number{Value: float64(udpAddr.Port)})
					result.Set("family", &object.String{Value: "IPv4"})
				}

				return result
			}
			return NULL
		}})

		// socket.setBroadcast(flag)
		socket.Set("setBroadcast", &object.Builtin{Name: "setBroadcast", Fn: func(args ...object.Object) object.Object {
			// Go UDP connections don't have a direct broadcast flag method
			// This is a placeholder
			return UNDEFINED
		}})

		// socket.setTTL(ttl)
		socket.Set("setTTL", &object.Builtin{Name: "setTTL", Fn: func(args ...object.Object) object.Object {
			// Placeholder - would need to be implemented at OS level
			return UNDEFINED
		}})

		if callback != nil {
			addEventListener(socket, "message", callback)
		}

		return socket
	}})

	return dgram
}
