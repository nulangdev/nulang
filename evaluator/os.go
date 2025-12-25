package evaluator

import (
	"os"
	goruntime "runtime"

	"github.com/nulang/nulang/object"
)

func initOsModule() *object.ObjectMap {
	osModule := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// os.platform()
	osModule.Set("platform", &object.Builtin{Name: "platform", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: goruntime.GOOS}
	}})

	// os.arch()
	osModule.Set("arch", &object.Builtin{Name: "arch", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: goruntime.GOARCH}
	}})

	// os.homedir()
	osModule.Set("homedir", &object.Builtin{Name: "homedir", Fn: func(args ...object.Object) object.Object {
		home, err := os.UserHomeDir()
		if err != nil {
			return newError("failed to get home directory: %s", err)
		}
		return &object.String{Value: home}
	}})

	// os.tmpdir()
	osModule.Set("tmpdir", &object.Builtin{Name: "tmpdir", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: os.TempDir()}
	}})

	// os.hostname()
	osModule.Set("hostname", &object.Builtin{Name: "hostname", Fn: func(args ...object.Object) object.Object {
		hostname, err := os.Hostname()
		if err != nil {
			return newError("failed to get hostname: %s", err)
		}
		return &object.String{Value: hostname}
	}})

	// os.type() - returns OS name
	osModule.Set("type", &object.Builtin{Name: "type", Fn: func(args ...object.Object) object.Object {
		switch goruntime.GOOS {
		case "darwin":
			return &object.String{Value: "Darwin"}
		case "linux":
			return &object.String{Value: "Linux"}
		case "windows":
			return &object.String{Value: "Windows_NT"}
		default:
			return &object.String{Value: goruntime.GOOS}
		}
	}})

	// os.cpus()
	osModule.Set("cpus", &object.Builtin{Name: "cpus", Fn: func(args ...object.Object) object.Object {
		numCPU := goruntime.NumCPU()
		cpus := make([]object.Object, numCPU)
		for i := 0; i < numCPU; i++ {
			cpu := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
			cpu.Set("model", &object.String{Value: "CPU"})
			cpu.Set("speed", &object.Number{Value: 2400})
			cpus[i] = cpu
		}
		return &object.Array{Elements: cpus}
	}})

	// os.uptime()
	osModule.Set("uptime", &object.Builtin{Name: "uptime", Fn: func(args ...object.Object) object.Object {
		// Not easily available in Go, return 0
		return &object.Number{Value: 0}
	}})

	// os.freemem()
	osModule.Set("freemem", &object.Builtin{Name: "freemem", Fn: func(args ...object.Object) object.Object {
		// Not easily available in Go, return 0
		return &object.Number{Value: 0}
	}})

	// os.totalmem()
	osModule.Set("totalmem", &object.Builtin{Name: "totalmem", Fn: func(args ...object.Object) object.Object {
		// Not easily available in Go, return 0
		return &object.Number{Value: 0}
	}})

	// os.EOL - End of line character
	eol := "\n"
	if goruntime.GOOS == "windows" {
		eol = "\r\n"
	}
	osModule.Set("EOL", &object.String{Value: eol})

	// os.devNull
	devNull := "/dev/null"
	if goruntime.GOOS == "windows" {
		devNull = "nul"
	}
	osModule.Set("devNull", &object.String{Value: devNull})

	return osModule
}
