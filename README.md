# Nulang 🚀

Uma linguagem de programação com sintaxe JavaScript/Node.js escrita em Go, pensada para ser leve, embutível e com uma biblioteca padrão rica inspirada no Node.js.

## Visão geral

- **Pipeline completo**: Lexer ➜ Parser (Pratt) ➜ AST ➜ Evaluator ➜ Event Loop completo.
- **Execução flexível**: rode arquivos `.js` ou `.ts`, use o REPL interativo ou o modo watch para desenvolvimento.
- **Compatível com JS**: sintaxe, tipos e APIs modeladas no ecossistema Node.js sempre que possível.

## Funcionalidades

### Linguagem

- Sintaxe JavaScript/Node.js com `let`, `const`, `var`, funções tradicionais e **arrow functions**.
- Controle de fluxo: `if/else`, `for`, `while`, `break`, `continue` e operador ternário.
- Tratamento de erros: `try/catch/finally` e `throw` com objeto **Error**.
- **Classes ES6**: `class`, `constructor`, `extends`, `static`, **getters/setters** e suporte a decorators.
- Operadores: aritméticos, comparação, lógicos, **nullish coalescing**, atribuição composta e incremento/decremento.
- Literais e tipos: `number`, `string`, `boolean`, `null`, `undefined`, template literals, arrays, objetos e parsing opcional de tipos estilo TypeScript.
- **RegExp**, **Date**, **Map/Set/WeakMap/WeakSet** com API completa.

### Runtime e objetos globais

- Objetos globais: `console`, `Math`, `JSON`, `Array`, `Object`, `Buffer`, `Promise`, `Reflect`, `Proxy`.
- Variáveis especiais: `__filename`, `__dirname` e **process** (argv, env, cwd, exit, stdin/stdout/stderr básicos).
- **Promises realmente assíncronas** com `resolve`, `reject`, `all`, `race`, `then`, `catch`, `finally`.
- Event loop completo com **Timers** (`setTimeout`, `setInterval`, `clearTimeout`, `clearInterval`, `sleep`) e integração com Promises.

### Biblioteca padrão

- **Módulos de núcleo** (via `require` ou `import`):
  - `fs` e `path` – leitura/escrita de arquivos, diretórios, stat, paths, workdir.
  - `crypto` – `createHash`, `createHmac`, `randomBytes`, `randomUUID` e utilitários.
  - `buffer` – criação, conversão e manipulação de dados binários.
  - `http`/`https` e `fetch` – servidor e cliente HTTP, headers, métodos e body streaming.
  - `os`, `process`, `child_process` – informações do SO, processo atual e execução de comandos.
  - `events` – `EventEmitter` compatível.
  - `stream` – `Readable`, `Writable`, `Transform` com piping.
  - `timers` – helpers para temporizadores (interno ao runtime).
  - `url` e `querystring` – parsing/formatting de URLs e query strings.
  - `assert` – asserções para testes.
  - `dns`, `net`, `dgram` – networking TCP/UDP/DNS.
  - `zlib` – compressão/descompressão.
  - `readline` – CLI interativo.
  - `vm` – execução dinâmica de código.
  - `util` – `promisify`, `inspect`, `format`, `types`.
  - `package` – utilidades de gerenciamento de pacotes.
- **Blob/File** e **Buffer** para arquivos e binários, com conversões para `hex`, `base64`, etc.
- **Streams** com suporte a backpressure e transformação.
- **HTTP/HTTPS** com suporte a headers, métodos, corpo de resposta e fetch-like API.

### Sistema de módulos

- Suporte a **CommonJS** (`require`, `exports`, `module.exports`) e **ES Modules** (`import/export`).
- Resolução de caminhos relativa e absoluta, com cache de módulos carregados.
- Variáveis de contexto por módulo: `__filename`, `__dirname`, `module`, `exports`, `require`.

### Ferramentas de desenvolvimento

- **REPL interativo** com ajuda integrada.
- **Watch mode** (`go run watch.go <arquivo.nu>`) recompila e reexecuta automaticamente.
- Exemplos em `examples/` cobrindo sintaxe básica, módulos, ES6 e novos recursos.

## Instalação

```bash
  curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

## Uso

### CLI e gerenciamento de pacotes

```bash
# Criar um novo projeto (gera nulang.yml)
nu init

# Instalar dependências listadas em nulang.yml e gerar/atualizar nulang.lock
nu install

# Rodar um arquivo normalmente
nu caminho/para/arquivo.js
```

- **nulang.yml**: manifest YAML com `name`, `version`, `main` e `dependencies`.
- **nulang.lock**: lockfile gerado automaticamente contendo URL, commit e checksum de cada dependência instalada.
- **node_modules/**: diretório onde os pacotes baixados são armazenados (semelhante ao `node_modules` do Node.js).

### Executar arquivo

```bash
nu arquivo.js
```

### REPL interativo

```bash
nu
```

```text
Nu v0.1.0 - JavaScript-like language written in Go
Type 'exit' to quit, 'help' for more info

nu> let x = 5
nu> console.log(x * 2)
10
nu> exit
```

### Watch mode

```bash
nu index.js --watch
```

### Estrutura de pacote

```yaml
# nulang.yml
name: meu-projeto
version: 1.0.0
main: index.js
dependencies:
  math-tools: https://github.com/exemplo/math-tools
  http-utils: https://github.com/exemplo/http-utils
```

```text
meu-projeto/
├── nulang.yml          # manifest (criado por nu init)
├── nulang.lock         # lockfile (criado/atualizado por nu install)
├── node_modules/       # pacotes instalados
│   └── math-tools/
│       └── index.js
├── index.js
└── lib/
    └── api.js
```

### Exemplos de uso real com pacotes

```javascript
// index.js
import math from "math-tools"; // resolvido a partir de node_modules
import { request } from "http-utils"; // usa fetch e event loop completo

async function main() {
  const doubled = math.times(21, 2);
  const res = await request({ url: "https://api.example.com/ping" });
  console.log({ doubled, status: res.status });
}

main();
```

## Exemplos rápidos

### Hello World

```javascript
console.log("Hello, World!");
```

### Funções e arrow functions

```javascript
function soma(a, b) {
  return a + b;
}

const multiplica = (a, b) => a * b;

console.log(soma(5, 3)); // 8
console.log(multiplica(4, 7)); // 28
```

### Arrays e métodos utilitários

```javascript
let numeros = [1, 2, 3, 4, 5];
console.log(numeros.map((n) => n * 2));
console.log(numeros.filter((n) => n % 2 === 0));
console.log(numeros.reduce((acc, n) => acc + n, 0));
```

### Módulos CommonJS

```javascript
// math_utils.js
exports.add = (a, b) => a + b;
exports.PI = 3.14159;

// main.js
let math = require("./math_utils.js");
console.log(math.add(5, 3));
console.log(math.PI);
```

### Filesystem e Buffer

```javascript
fs.writeFileSync("arquivo.txt", "Conteúdo");
let data = fs.readFileSync("arquivo.txt", "utf8");
console.log(data);

let buf = Buffer.from("Hello");
console.log(buf.toString("hex"));
```

### HTTP e fetch

```javascript
http.get("https://example.com", (res) => {
  res.on("data", (chunk) => console.log(chunk.toString()));
});

let response = await fetch("https://example.com");
console.log(response.status);
```

## Estrutura do projeto

```text
nulang/
├── main.go              # Entrada do CLI/REPL
├── watch.go             # Watch mode
├── token/               # Definição de tokens
├── lexer/               # Lexer
├── ast/                 # Nós da AST
├── parser/              # Parser (Pratt)
├── object/              # Sistema de objetos e environment
├── evaluator/           # Avaliador, operadores, built-ins e módulos
├── doc/                 # Documentação detalhada (Arrays, Strings, HTTP, etc.)
├── examples/            # Exemplos de uso
└── docs/                # Documentação adicional e releases
```

## Roadmap e limitações

- ✅ Classes e herança
- ✅ Import/Export ES6
- ✅ HTTP/HTTPS client & server
- ✅ Streams, timers, EventEmitter
- ✅ RegExp, Date, Map/Set
- ✅ Promises com event loop completo
- ⏳ Async/Await com scheduling mais avançado (futuro)

## Recursos adicionais

- Documentação completa em `doc/` (arrays, strings, promessas, streams, HTTP, eventos, fs, crypto, etc.).
- Exemplos em `examples/` cobrindo módulos, ES6, novas APIs e funcionalidades avançadas.
- Histórico de releases em `RELEASE.md` e notas rápidas em `QUICKSTART.md`.

## Licença

MIT License.

---

Feito com ❤️ em Go.
