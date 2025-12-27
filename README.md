# Nulang 🚀

Uma linguagem de programação com sintaxe JavaScript/Node.js escrita em Go.

## Características

- ✅ Sintaxe JavaScript/Node.js
- ✅ Tipos primitivos: `number`, `string`, `boolean`, `null`, `undefined`
- ✅ Tipos estruturais: `Array`, `Object`, `Function`, `Buffer`, `Promise`
- ✅ Variáveis: `let`, `const`, `var`
- ✅ Funções: `function`, arrow functions (`=>`)
- ✅ Controle de fluxo: `if/else`, `for`, `while`, `break`, `continue`
- ✅ Tratamento de erros: `try/catch/finally`, `throw`
- ✅ Operadores: aritméticos, comparação, lógicos, ternário, nullish coalescing
- ✅ Métodos de Array: `push`, `pop`, `map`, `filter`, `reduce`, `forEach`, `find`, etc.
- ✅ Métodos de String: `toUpperCase`, `toLowerCase`, `split`, `trim`, `indexOf`, etc.
- ✅ Objetos globais: `console`, `Math`, `JSON`, `Array`, `Object`, `Buffer`, `Promise`
- ✅ Sistema de módulos: `require()`, `exports`, `module.exports`
- ✅ Módulo `fs` - Sistema de arquivos (readFileSync, writeFileSync, etc.)
- ✅ Módulo `path` - Manipulação de paths (join, dirname, basename, etc.)
- ✅ Módulo `crypto` - Criptografia (createHash, createHmac, randomBytes, randomUUID)
- ✅ `Buffer` - Manipulação de dados binários
- ✅ `Promise` - Promises síncronas (resolve, reject, all, race, then, catch, finally)
- ✅ REPL interativo
- ✅ Variáveis `__filename` e `__dirname`
- ✅ Objeto `process` (argv, env, cwd, exit)
- ✅ **Classes ES6** - `class`, `constructor`, `extends`, `static`, getters/setters
- ✅ **HTTP/HTTPS** - `http.get`, `http.post`, `fetch`
- ✅ **Timers** - `setTimeout`, `setInterval`, `clearTimeout`, `clearInterval`, `sleep`
- ✅ **Streams** - `stream.Readable`, `stream.Writable`, `stream.Transform`
- ✅ **URL/QueryString** - `url.parse`, `querystring.stringify/parse`
- ✅ **Interfaces/Types** - Parsing de tipos TypeScript-like (não enforced)
- ✅ **RegExp** - Regular expressions com `test`, `match`, `replace`, `split`
- ✅ **Date** - Date object com todos os métodos padrão
- ✅ **Map/Set** - Map, Set, WeakMap, WeakSet com operações completas

## Instalação

```bash
# Clone o repositório
git clone https://github.com/nulangdev/nulang.git
cd nulang

# Compile
go build -o nulang .

# (Opcional) Instale globalmente
go install
```

## Uso

### Executando arquivos

```bash
nulang arquivo.nu
```

### REPL Interativo

```bash
nulang
```

```
Nulang v0.1.0 - JavaScript-like language written in Go
Type 'exit' to quit, 'help' for more info

nulang> let x = 5
nulang> console.log(x * 2)
10
nulang> exit
```

## Exemplos

### Hello World

```javascript
console.log("Hello, World!");
```

### Variáveis

```javascript
let nome = "Nulang";
const versao = "0.1.0";
var contador = 0;

console.log(nome, versao);
```

### Funções

```javascript
// Função tradicional
function soma(a, b) {
  return a + b;
}

// Arrow function
const multiplica = (a, b) => a * b;

console.log(soma(5, 3)); // 8
console.log(multiplica(4, 7)); // 28
```

### Arrays

```javascript
let numeros = [1, 2, 3, 4, 5];

// Map
let dobrados = numeros.map((n) => n * 2);
console.log(dobrados); // [2, 4, 6, 8, 10]

// Filter
let pares = numeros.filter((n) => n % 2 === 0);
console.log(pares); // [2, 4]

// Reduce
let soma = numeros.reduce((acc, n) => acc + n, 0);
console.log(soma); // 15
```

### Objetos

```javascript
let pessoa = {
  nome: "João",
  idade: 30,
  cidade: "São Paulo",
};

console.log(pessoa.nome); // João
console.log(pessoa["idade"]); // 30
```

### Controle de Fluxo

```javascript
// If/Else
let idade = 20;
if (idade >= 18) {
  console.log("Maior de idade");
} else {
  console.log("Menor de idade");
}

// For loop
for (let i = 0; i < 5; i++) {
  console.log(i);
}

// While loop
let contador = 5;
while (contador > 0) {
  console.log(contador);
  contador--;
}
```

### Sistema de Arquivos (fs)

```javascript
// Escrever arquivo
fs.writeFileSync("arquivo.txt", "Conteúdo do arquivo");

// Ler arquivo
let conteudo = fs.readFileSync("arquivo.txt", "utf8");
console.log(conteudo);

// Verificar se existe
console.log(fs.existsSync("arquivo.txt")); // true

// Informações do arquivo
let stats = fs.statSync("arquivo.txt");
console.log(stats.size); // tamanho em bytes
console.log(stats.isFile()); // true
console.log(stats.isDirectory()); // false

// Criar diretório
fs.mkdirSync("pasta", { recursive: true });

// Listar diretório
console.log(fs.readdirSync(".")); // ['arquivo.txt', 'pasta', ...]

// Deletar arquivo
fs.unlinkSync("arquivo.txt");
```

### Criptografia (crypto)

```javascript
// Hash SHA256
let hash = crypto.createHash("sha256");
hash.update("Hello, World!");
console.log(hash.digest("hex"));

// HMAC
let hmac = crypto.createHmac("sha256", "secret-key");
hmac.update("message");
console.log(hmac.digest("hex"));

// Bytes aleatórios
let bytes = crypto.randomBytes(16);
console.log(bytes.toString("hex"));

// UUID aleatório
console.log(crypto.randomUUID());
```

### Buffer

```javascript
// Criar buffer de string
let buf = Buffer.from("Hello");
console.log(buf.toString()); // Hello
console.log(buf.toString("hex")); // 48656c6c6f
console.log(buf.toString("base64")); // SGVsbG8=

// Criar buffer de array
let buf2 = Buffer.from([72, 101, 108, 108, 111]);
console.log(buf2.toString()); // Hello

// Alocar buffer
let buf3 = Buffer.alloc(10, 0);
buf3.fill(65);
console.log(buf3.toString()); // AAAAAAAAAA
```

### Promises

```javascript
// Promise.resolve
let p1 = Promise.resolve(42);
console.log(p1.value); // 42

// then/catch
let result = p1.then((x) => x * 2);
console.log(result.value); // 84

// Promise.all
let all = Promise.all([
  Promise.resolve(1),
  Promise.resolve(2),
  Promise.resolve(3),
]);
console.log(all.value); // [1, 2, 3]
```

### Módulos

```javascript
// math_utils.nu
exports.add = (a, b) => a + b;
exports.PI = 3.14159;

// main.nu
let math = require("./math_utils.nu");
console.log(math.add(5, 3)); // 8
console.log(math.PI); // 3.14159
```

## Estrutura do Projeto

```
nulang/
├── main.go           # Ponto de entrada, CLI e REPL
├── token/
│   └── token.go      # Definição de tokens
├── lexer/
│   └── lexer.go      # Analisador léxico
├── ast/
│   └── ast.go        # Árvore Sintática Abstrata
├── parser/
│   └── parser.go     # Analisador sintático (Pratt Parser)
├── object/
│   ├── object.go     # Sistema de objetos do runtime
│   └── environment.go # Ambientes (scopes)
├── evaluator/
│   ├── evaluator.go  # Avaliador principal
│   ├── operators.go  # Operadores
│   ├── functions.go  # Funções e expressões
│   ├── array_methods.go # Métodos de Array
│   ├── string_methods.go # Métodos de String
│   ├── builtins.go   # Funções e objetos built-in
│   ├── modules.go    # Sistema de módulos
│   ├── fs.go         # Módulo fs e path
│   ├── crypto.go     # Módulo crypto
│   ├── buffer.go     # Buffer
│   └── promise.go    # Promises
└── examples/
    ├── example.nu    # Exemplo completo
    ├── simple.nu     # Exemplo simples
    ├── new_features.nu # Novos recursos
    └── modules/      # Exemplos de módulos
```

## Roadmap

- [x] Classes e herança ✅
- [ ] Async/Await (assíncrono real)
- [x] Import/Export ES6 syntax ✅
- [x] HTTP/HTTPS client ✅
- [x] Regular Expressions ✅
- [x] Date ✅
- [x] Map e Set ✅
- [x] Timers (setTimeout, setInterval) ✅
- [x] Streams ✅

## Licença

MIT License

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

---

Feito com ❤️ em Go
