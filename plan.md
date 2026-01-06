Nulang é nossa linguagem de programação escrita em Golang com a sintaxe do Javascript/NodeJs:
Devemos criar um interpretador para que possamos executar o código.

Execução:
`$ nu index.js`

#### Tipos primitivos:

- number
- string
- boolean
- null
- undefined
- symbol
- bigint

#### Tipos estruturais:

- Object
- Array
- Function
- Date
- RegExp
- Map
- Set
- WeakMap
- WeakSet

#### Controle de fluxo

- if / else
- switch
- for
- while
- do...while
- break
- continue
- return
- throw
- try / catch / finally

#### Funções:

- function fn() {}
- () => {}
- async function
- async () => {}

#### Escopo e variáveis:

- var
- let
- const

#### Array:

- push
- pop
- shift
- unshift
- map
- filter
- reduce
- forEach
- find
- includes
- length
- from
- isArray
- of
- is
- join
- keys
- values
- entries
- flat
- flatMap
- sort
- reverse
- slice
- splice
- concat
- join
- pop
- push
- shift
- unshift
- findIndex
- where

#### Imports/Exports:

- `import fs from "fs"`
- `import { readFile } from "fs"`
- `import * as fs from "fs"`
- `import "fs"`
- `export default fs`
- `export { readFile } from "fs"`
- `export * as fs from "fs"`
- `export * from "fs"`

#### Object:

- keys
- values
- entries
- hasOwnProperty
- isPrototypeOf
- propertyIsEnumerable
- defineProperty
- defineProperties
- getOwnPropertyDescriptor
- getOwnPropertyDescriptors
- getPrototypeOf
- setPrototypeOf
- assign
- create
- freeze
- getOwnPropertyNames
- getOwnPropertySymbols
- is
- isExtensible
- isFrozen
- isSealed
- preventExtensions
- seal
- toLocaleString
- toString
- valueOf

#### String:

- charAt
- charCodeAt
- concat
- indexOf
- lastIndexOf
- match
- replace
- search
- slice
- split
- substr
- substring
- toLowerCase
- toUpperCase
- trim
- trimLeft
- trimRight
- valueOf

#### Number:

- toFixed
- toExponential
- toPrecision
- toLocaleString
- toString
- valueOf

#### Operadores:

- `+ - * / % **`
- `== === != !==`
- `< > <= >=`
- `&& || ??`
- `?:`
- `typeof`
- `instanceof`
- `in`
- `delete`

#### Globals:

- globalThis
- Infinity
- NaN
- undefined

#### Math:

Math.abs
Math.ceil
Math.floor
Math.round
Math.max
Math.min
Math.random
Math.pow
Math.sqrt

#### JSON:

JSON.stringify
JSON.parse

#### Reflect:

Reflect.get
Reflect.set
Reflect.apply
Reflect.construct

#### Proxy:

new Proxy(target, handler)

#### Sistema:

- process
- process.env
- process.argv
- process.exit

#### Timers:

- setTimeout
- setInterval
- setImmediate
- clearTimeout
- clearInterval

#### Eventos:

- EventEmitter

#### Buffer:

- Buffer.from
- Buffer.alloc
- Buffer.concat

#### FS:

- fs.readFile
- fs.writeFile
- fs.appendFile
- fs.unlink
- fs.mkdir
- fs.rmdir
- fs.readdir
- fs.stat
- fs.rename
- fs.copyFile
- fs.createReadStream
- fs.createWriteStream

#### HTTP / HTTPS:

- http.createServer
- https.createServer
- http.request
- https.request
- http.get
- https.get
- http.put
- https.put
- http.post
- https.post
- http.delete
- https.delete
- http.patch
- https.patch
- http.head
- https.head
- http.connect
- https.connect
- http.options
- https.options
- http.trace
- https.trace
- http.copy
- https.copy
- http.link
- https.link
- http.unlink
- https.unlink
- http.purge
- https.purge
- http.propfind
- https.propfind
- http.proppatch
- https.proppatch
- http.mkcol
- https.mkcol
- http.move
- https.move
- http.copy
- https.copy
- http.lock
- https.lock
- http.unlock
- https.unlock
- http.version

#### Path:

- path.join
- path.resolve
- path.dirname
- path.basename
- path.extname
- path.parse
- path.format

#### OS:

- os.platform
- os.release
- os.type
- os.uptime
- os.loadavg
- os.freemem
- os.totalmem
- os.cpus
- os.hostname
- os.endianness
- os.tmpdir
- os.homedir
- os.tmpdir
- os.networkInterfaces
- os.EOL

#### Crypto:

- crypto.createHash
- crypto.createHmac
- crypto.createCipher
- crypto.createCipheriv
- crypto.createDecipher
- crypto.createDecipheriv
- crypto.randomBytes
- crypto.pbkdf2
- crypto.pbkdf2Sync
- crypto.scrypt
- crypto.scryptSync
- crypto.timingSafeEqual
- crypto.timingSafeEqual

#### Child Process:

- child_process.spawn
- child_process.exec
- child_process.execFile
- child_process.fork
- child_process.spawnSync
- child_process.execSync
- child_process.execFileSync

#### Stream:

- stream.Readable
- stream.Writable
- stream.Duplex
- stream.Pipeline
- stream.Transform
- stream.PassThrough

Conceitos obrigatórios
• Call Stack
• Microtasks
• Macrotasks
• Promise queue
• Async/Await lowering
