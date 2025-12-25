package evaluator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nulang/nulang/object"
)

func initFsModule() *object.ObjectMap {
	fs := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// fs.readFileSync(path, encoding?)
	fs.Set("readFileSync", &object.Builtin{Name: "readFileSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("readFileSync requires path argument")
		}
		path := objectToString(args[0])
		
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return newError("failed to read file: %s", err)
		}

		// Check for encoding option
		if len(args) > 1 {
			encoding := ""
			if str, ok := args[1].(*object.String); ok {
				encoding = str.Value
			} else if obj, ok := args[1].(*object.ObjectMap); ok {
				if enc, ok := obj.Get("encoding"); ok {
					encoding = objectToString(enc)
				}
			}
			if encoding == "utf8" || encoding == "utf-8" {
				return &object.String{Value: string(data)}
			}
		}
		return &object.Buffer{Data: data}
	}})

	// fs.writeFileSync(path, data, encoding?)
	fs.Set("writeFileSync", &object.Builtin{Name: "writeFileSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("writeFileSync requires path and data arguments")
		}
		path := objectToString(args[0])
		
		var data []byte
		switch d := args[1].(type) {
		case *object.String:
			data = []byte(d.Value)
		case *object.Buffer:
			data = d.Data
		default:
			data = []byte(objectToString(args[1]))
		}

		err := ioutil.WriteFile(path, data, 0644)
		if err != nil {
			return newError("failed to write file: %s", err)
		}
		return UNDEFINED
	}})

	// fs.appendFileSync(path, data)
	fs.Set("appendFileSync", &object.Builtin{Name: "appendFileSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("appendFileSync requires path and data arguments")
		}
		path := objectToString(args[0])
		
		var data []byte
		switch d := args[1].(type) {
		case *object.String:
			data = []byte(d.Value)
		case *object.Buffer:
			data = d.Data
		default:
			data = []byte(objectToString(args[1]))
		}

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return newError("failed to open file: %s", err)
		}
		defer f.Close()
		
		_, err = f.Write(data)
		if err != nil {
			return newError("failed to append to file: %s", err)
		}
		return UNDEFINED
	}})

	// fs.existsSync(path)
	fs.Set("existsSync", &object.Builtin{Name: "existsSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		path := objectToString(args[0])
		_, err := os.Stat(path)
		return nativeBoolToBooleanObject(!os.IsNotExist(err))
	}})

	// fs.unlinkSync(path) - delete file
	fs.Set("unlinkSync", &object.Builtin{Name: "unlinkSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("unlinkSync requires path argument")
		}
		path := objectToString(args[0])
		err := os.Remove(path)
		if err != nil {
			return newError("failed to delete file: %s", err)
		}
		return UNDEFINED
	}})

	// fs.mkdirSync(path, options?)
	fs.Set("mkdirSync", &object.Builtin{Name: "mkdirSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("mkdirSync requires path argument")
		}
		path := objectToString(args[0])
		
		recursive := false
		if len(args) > 1 {
			if obj, ok := args[1].(*object.ObjectMap); ok {
				if rec, ok := obj.Get("recursive"); ok {
					if b, ok := rec.(*object.Boolean); ok {
						recursive = b.Value
					}
				}
			}
		}

		var err error
		if recursive {
			err = os.MkdirAll(path, 0755)
		} else {
			err = os.Mkdir(path, 0755)
		}
		if err != nil {
			return newError("failed to create directory: %s", err)
		}
		return UNDEFINED
	}})

	// fs.rmdirSync(path, options?)
	fs.Set("rmdirSync", &object.Builtin{Name: "rmdirSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("rmdirSync requires path argument")
		}
		path := objectToString(args[0])
		
		recursive := false
		if len(args) > 1 {
			if obj, ok := args[1].(*object.ObjectMap); ok {
				if rec, ok := obj.Get("recursive"); ok {
					if b, ok := rec.(*object.Boolean); ok {
						recursive = b.Value
					}
				}
			}
		}

		var err error
		if recursive {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}
		if err != nil {
			return newError("failed to remove directory: %s", err)
		}
		return UNDEFINED
	}})

	// fs.readdirSync(path)
	fs.Set("readdirSync", &object.Builtin{Name: "readdirSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("readdirSync requires path argument")
		}
		path := objectToString(args[0])
		
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return newError("failed to read directory: %s", err)
		}

		elements := make([]object.Object, len(files))
		for i, file := range files {
			elements[i] = &object.String{Value: file.Name()}
		}
		return &object.Array{Elements: elements}
	}})

	// fs.statSync(path)
	fs.Set("statSync", &object.Builtin{Name: "statSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return newError("statSync requires path argument")
		}
		path := objectToString(args[0])
		
		info, err := os.Stat(path)
		if err != nil {
			return newError("failed to stat file: %s", err)
		}

		stat := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		stat.Set("size", &object.Number{Value: float64(info.Size())})
		stat.Set("isFile", &object.Builtin{Name: "isFile", Fn: func(args ...object.Object) object.Object {
			return nativeBoolToBooleanObject(!info.IsDir())
		}})
		stat.Set("isDirectory", &object.Builtin{Name: "isDirectory", Fn: func(args ...object.Object) object.Object {
			return nativeBoolToBooleanObject(info.IsDir())
		}})
		stat.Set("mtime", &object.Number{Value: float64(info.ModTime().UnixMilli())})
		stat.Set("mode", &object.Number{Value: float64(info.Mode())})

		return stat
	}})

	// fs.renameSync(oldPath, newPath)
	fs.Set("renameSync", &object.Builtin{Name: "renameSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("renameSync requires oldPath and newPath arguments")
		}
		oldPath := objectToString(args[0])
		newPath := objectToString(args[1])
		
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return newError("failed to rename: %s", err)
		}
		return UNDEFINED
	}})

	// fs.copyFileSync(src, dest)
	fs.Set("copyFileSync", &object.Builtin{Name: "copyFileSync", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return newError("copyFileSync requires src and dest arguments")
		}
		src := objectToString(args[0])
		dest := objectToString(args[1])
		
		data, err := ioutil.ReadFile(src)
		if err != nil {
			return newError("failed to read source file: %s", err)
		}
		
		err = ioutil.WriteFile(dest, data, 0644)
		if err != nil {
			return newError("failed to write destination file: %s", err)
		}
		return UNDEFINED
	}})

	return fs
}

func initPathModule() *object.ObjectMap {
	pathMod := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// path.join(...paths)
	pathMod.Set("join", &object.Builtin{Name: "join", Fn: func(args ...object.Object) object.Object {
		parts := make([]string, len(args))
		for i, arg := range args {
			parts[i] = objectToString(arg)
		}
		return &object.String{Value: filepath.Join(parts...)}
	}})

	// path.resolve(...paths)
	pathMod.Set("resolve", &object.Builtin{Name: "resolve", Fn: func(args ...object.Object) object.Object {
		if len(args) == 0 {
			cwd, _ := os.Getwd()
			return &object.String{Value: cwd}
		}
		parts := make([]string, len(args))
		for i, arg := range args {
			parts[i] = objectToString(arg)
		}
		result := filepath.Join(parts...)
		abs, _ := filepath.Abs(result)
		return &object.String{Value: abs}
	}})

	// path.dirname(path)
	pathMod.Set("dirname", &object.Builtin{Name: "dirname", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: "."}
		}
		return &object.String{Value: filepath.Dir(objectToString(args[0]))}
	}})

	// path.basename(path, ext?)
	pathMod.Set("basename", &object.Builtin{Name: "basename", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: ""}
		}
		base := filepath.Base(objectToString(args[0]))
		if len(args) > 1 {
			ext := objectToString(args[1])
			if len(base) > len(ext) && base[len(base)-len(ext):] == ext {
				base = base[:len(base)-len(ext)]
			}
		}
		return &object.String{Value: base}
	}})

	// path.extname(path)
	pathMod.Set("extname", &object.Builtin{Name: "extname", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.String{Value: ""}
		}
		return &object.String{Value: filepath.Ext(objectToString(args[0]))}
	}})

	// path.parse(path)
	pathMod.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		}
		p := objectToString(args[0])
		parsed := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}
		parsed.Set("root", &object.String{Value: filepath.VolumeName(p) + string(filepath.Separator)})
		parsed.Set("dir", &object.String{Value: filepath.Dir(p)})
		parsed.Set("base", &object.String{Value: filepath.Base(p)})
		parsed.Set("ext", &object.String{Value: filepath.Ext(p)})
		base := filepath.Base(p)
		ext := filepath.Ext(p)
		name := base
		if len(ext) > 0 {
			name = base[:len(base)-len(ext)]
		}
		parsed.Set("name", &object.String{Value: name})
		return parsed
	}})

	// path.isAbsolute(path)
	pathMod.Set("isAbsolute", &object.Builtin{Name: "isAbsolute", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return FALSE
		}
		return nativeBoolToBooleanObject(filepath.IsAbs(objectToString(args[0])))
	}})

	// path.sep
	pathMod.Set("sep", &object.String{Value: string(filepath.Separator)})

	return pathMod
}
