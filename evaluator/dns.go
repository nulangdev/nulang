package evaluator

import (
	"net"

	"github.com/nulang/nulang/object"
)

// initDNSModule initializes the dns module
func initDNSModule() *object.ObjectMap {
	dns := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// dns.lookup(hostname, options?, callback)
	dns.Set("lookup", &object.Builtin{Name: "lookup", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("lookup requires a hostname argument")
		}

		hostname := objectToString(args[0])
		var callback object.Object
		family := 0 // 0 = both, 4 = IPv4, 6 = IPv6

		// Parse arguments
		if len(args) > 1 {
			if opts, ok := args[1].(*object.ObjectMap); ok {
				if fam, ok := opts.Get("family"); ok {
					if num, ok := fam.(*object.Number); ok {
						family = int(num.Value)
					}
				}
				if len(args) > 2 {
					callback = args[2]
				}
			} else {
				callback = args[1]
			}
		}

		// Perform lookup
		ips, err := net.LookupIP(hostname)
		if err != nil {
			errorObj := createErrorObject(err.Error())
			if callback != nil {
				callFunction(callback, []object.Object{errorObj})
			}
			return errorObj
		}

		// Filter by family if specified
		var selectedIP net.IP
		for _, ip := range ips {
			if family == 0 {
				selectedIP = ip
				break
			} else if family == 4 && ip.To4() != nil {
				selectedIP = ip
				break
			} else if family == 6 && ip.To16() != nil && ip.To4() == nil {
				selectedIP = ip
				break
			}
		}

		if selectedIP == nil && len(ips) > 0 {
			selectedIP = ips[0]
		}

		address := selectedIP.String()
		familyNum := 4
		if selectedIP.To4() == nil {
			familyNum = 6
		}

		if callback != nil {
			callFunction(callback, []object.Object{
				NULL,
				&object.String{Value: address},
				&object.Number{Value: float64(familyNum)},
			})
		}

		return &object.String{Value: address}
	}})

	// dns.resolve(hostname, rrtype?, callback)
	dns.Set("resolve", &object.Builtin{Name: "resolve", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("resolve requires a hostname argument")
		}

		hostname := objectToString(args[0])
		rrtype := "A" // Default to A records
		var callback object.Object

		if len(args) > 1 {
			if str, ok := args[1].(*object.String); ok {
				rrtype = str.Value
				if len(args) > 2 {
					callback = args[2]
				}
			} else {
				callback = args[1]
			}
		}

		var addresses []object.Object
		var err error

		switch rrtype {
		case "A", "AAAA":
			ips, lookupErr := net.LookupIP(hostname)
			err = lookupErr
			if err == nil {
				for _, ip := range ips {
					if rrtype == "A" && ip.To4() != nil {
						addresses = append(addresses, &object.String{Value: ip.String()})
					} else if rrtype == "AAAA" && ip.To4() == nil {
						addresses = append(addresses, &object.String{Value: ip.String()})
					}
				}
			}
		case "CNAME":
			cname, lookupErr := net.LookupCNAME(hostname)
			err = lookupErr
			if err == nil {
				addresses = append(addresses, &object.String{Value: cname})
			}
		case "MX":
			mxs, lookupErr := net.LookupMX(hostname)
			err = lookupErr
			if err == nil {
				for _, mx := range mxs {
					mxObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
					mxObj.Set("exchange", &object.String{Value: mx.Host})
					mxObj.Set("priority", &object.Number{Value: float64(mx.Pref)})
					addresses = append(addresses, mxObj)
				}
			}
		case "NS":
			nss, lookupErr := net.LookupNS(hostname)
			err = lookupErr
			if err == nil {
				for _, ns := range nss {
					addresses = append(addresses, &object.String{Value: ns.Host})
				}
			}
		case "TXT":
			txts, lookupErr := net.LookupTXT(hostname)
			err = lookupErr
			if err == nil {
				for _, txt := range txts {
					addresses = append(addresses, &object.String{Value: txt})
				}
			}
		default:
			err = &net.DNSError{Err: "unsupported record type", Name: hostname}
		}

		result := &object.Array{Elements: addresses}

		if err != nil {
			errorObj := createErrorObject(err.Error())
			if callback != nil {
				callFunction(callback, []object.Object{errorObj})
			}
			return errorObj
		}

		if callback != nil {
			callFunction(callback, []object.Object{NULL, result})
		}

		return result
	}})

	// dns.resolve4(hostname, options?, callback)
	dns.Set("resolve4", &object.Builtin{Name: "resolve4", Fn: func(args ...object.Object) object.Object {
		newArgs := make([]object.Object, len(args)+1)
		newArgs[0] = args[0]
		newArgs[1] = &object.String{Value: "A"}
		if len(args) > 1 {
			copy(newArgs[2:], args[1:])
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.resolve6(hostname, options?, callback)
	dns.Set("resolve6", &object.Builtin{Name: "resolve6", Fn: func(args ...object.Object) object.Object {
		newArgs := make([]object.Object, len(args)+1)
		newArgs[0] = args[0]
		newArgs[1] = &object.String{Value: "AAAA"}
		if len(args) > 1 {
			copy(newArgs[2:], args[1:])
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.resolveMx(hostname, callback)
	dns.Set("resolveMx", &object.Builtin{Name: "resolveMx", Fn: func(args ...object.Object) object.Object {
		newArgs := []object.Object{args[0], &object.String{Value: "MX"}}
		if len(args) > 1 {
			newArgs = append(newArgs, args[1:]...)
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.resolveNs(hostname, callback)
	dns.Set("resolveNs", &object.Builtin{Name: "resolveNs", Fn: func(args ...object.Object) object.Object {
		newArgs := []object.Object{args[0], &object.String{Value: "NS"}}
		if len(args) > 1 {
			newArgs = append(newArgs, args[1:]...)
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.resolveTxt(hostname, callback)
	dns.Set("resolveTxt", &object.Builtin{Name: "resolveTxt", Fn: func(args ...object.Object) object.Object {
		newArgs := []object.Object{args[0], &object.String{Value: "TXT"}}
		if len(args) > 1 {
			newArgs = append(newArgs, args[1:]...)
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.resolveCname(hostname, callback)
	dns.Set("resolveCname", &object.Builtin{Name: "resolveCname", Fn: func(args ...object.Object) object.Object {
		newArgs := []object.Object{args[0], &object.String{Value: "CNAME"}}
		if len(args) > 1 {
			newArgs = append(newArgs, args[1:]...)
		}
		return dns.Pairs["resolve"].Value.(*object.Builtin).Fn(newArgs...)
	}})

	// dns.reverse(ip, callback)
	dns.Set("reverse", &object.Builtin{Name: "reverse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("reverse requires an IP address argument")
		}

		ip := objectToString(args[0])
		var callback object.Object
		if len(args) > 1 {
			callback = args[1]
		}

		names, err := net.LookupAddr(ip)
		if err != nil {
			errorObj := createErrorObject(err.Error())
			if callback != nil {
				callFunction(callback, []object.Object{errorObj})
			}
			return errorObj
		}

		addresses := make([]object.Object, len(names))
		for i, name := range names {
			addresses[i] = &object.String{Value: name}
		}

		result := &object.Array{Elements: addresses}

		if callback != nil {
			callFunction(callback, []object.Object{NULL, result})
		}

		return result
	}})

	return dns
}

func callFunction(fn object.Object, args []object.Object) object.Object {
	if function, ok := fn.(*object.Function); ok {
		fnEnv := extendFunctionEnv(function, args)
		return Eval(function.Body, fnEnv)
	} else if builtin, ok := fn.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}
	return UNDEFINED
}
