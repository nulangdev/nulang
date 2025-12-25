# AGENTS.md - Nulang Interpreter

Este documento descreve a arquitetura e implementação do interpretador Nulang para orientar agentes de IA e desenvolvedores que precisem entender ou modificar o projeto.

## Visão Geral

**Nulang** é uma linguagem de programação interpretada com sintaxe JavaScript/Node.js, escrita em Go. O objetivo é fornecer uma runtime JavaScript-like leve e embeddable.

## Arquitetura

O interpretador segue a arquitetura clássica de interpretadores:

```
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Output
```

### Fluxo de Execução

1. **Entrada**: Código fonte `.nu` ou input do REPL
2. **Lexer**: Converte texto em tokens
3. **Parser**: Constrói AST (Abstract Syntax Tree) usando Pratt Parser
4. **Evaluator**: Percorre a AST e executa as instruções
5. **Object System**: Representa valores em runtime

## Estrutura do Projeto

```
nulang/
├── main.go                    # CLI e REPL
├── go.mod                     # Go module
├── token/
│   └── token.go               # Definição de tipos de tokens
├── lexer/
│   └── lexer.go               # Analisador léxico
├── ast/
│   └── ast.go                 # Nós da AST
├── parser/
│   └── parser.go              # Parser (Pratt Parser)
├── object/
│   ├── object.go              # Sistema de objetos runtime
│   └── environment.go         # Gerenciamento de escopo
├── evaluator/
│   ├── evaluator.go           # Avaliador principal
│   ├── operators.go           # Operadores (prefix, infix, postfix)
│   ├── functions.go           # Funções, arrays, objetos, atribuições
│   ├── array_methods.go       # Métodos de Array (map, filter, etc.)
│   ├── string_methods.go      # Métodos de String
│   ├── builtins.go            # Objetos globais (console, Math, JSON)
│   ├── modules.go             # Sistema de módulos (require, import/export)
│   ├── fs.go                  # Módulo fs e path
│   ├── crypto.go              # Módulo crypto
│   ├── buffer.go              # Buffer
│   ├── promise.go             # Promises
│   └── os.go                  # Módulo os
└── examples/
    ├── example.nu             # Exemplo completo
    ├── simple.nu              # Exemplo básico
    ├── new_features.nu        # Testes de fs, crypto, buffer, promise
    ├── test_es6.nu            # Testes de import ES6
    ├── test_modules.nu        # Testes de módulos
    └── modules/
        └── math_utils.nu      # Módulo de exemplo
```

## Componentes Detalhados

### 1. Token (`token/token.go`)

Define todos os tipos de tokens da linguagem:

- **Literais**: `NUMBER`, `STRING`, `IDENT`, `REGEX`
- **Operadores**: `+`, `-`, `*`, `/`, `%`, `**`, `==`, `===`, `!=`, `!==`, etc.
- **Delimitadores**: `(`, `)`, `{`, `}`, `[`, `]`, `,`, `;`, `:`
- **Keywords**: `let`, `const`, `var`, `function`, `if`, `else`, `for`, `while`, `return`, `try`, `catch`, `finally`, `throw`, `import`, `export`, `from`, `async`, `await`, etc.

```go
type TokenType string

type Token struct {
    Type    TokenType
    Literal string
    Line    int
    Column  int
}
```

### 2. Lexer (`lexer/lexer.go`)

Converte código fonte em stream de tokens:

- Suporta strings com `"`, `'`, `` ` ``
- Processa escape sequences (`\n`, `\t`, `\\`, etc.)
- Ignora comentários (`//` e `/* */`)
- Rastreia linha e coluna para erros

**Pontos importantes:**

- `readString()` processa escape sequences corretamente
- Após ler string, chama `readChar()` para consumir a aspas de fechamento

### 3. AST (`ast/ast.go`)

Define todos os nós da árvore sintática:

**Interfaces base:**

```go
type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface { Node; statementNode() }
type Expression interface { Node; expressionNode() }
```

**Statements principais:**

- `LetStatement`, `ConstStatement`, `VarStatement`
- `ReturnStatement`, `BreakStatement`, `ContinueStatement`
- `ForStatement`, `WhileStatement`
- `TryStatement`, `ThrowStatement`
- `ImportStatement`, `ExportStatement`

**Expressions principais:**

- Literais: `NumberLiteral`, `StringLiteral`, `BooleanLiteral`, `NullLiteral`, `UndefinedLiteral`
- `Identifier`, `FunctionLiteral`, `ArrayLiteral`, `ObjectLiteral`
- `PrefixExpression`, `InfixExpression`, `PostfixExpression`
- `CallExpression`, `MemberExpression`, `IndexExpression`
- `IfExpression`, `ConditionalExpression` (ternário)
- `AssignmentExpression`, `NewExpression`, `ThisExpression`
- `TypeofExpression`, `AwaitExpression`, `SpreadExpression`

### 4. Parser (`parser/parser.go`)

Implementa um **Pratt Parser** (Top-Down Operator Precedence):

**Precedências (do menor para maior):**

```go
const (
    _ int = iota
    LOWEST
    COMMA       // ,
    ASSIGN      // = += -= *= /= %=
    TERNARY     // ?:
    NULLISH_C   // ??
    OR          // ||
    AND         // &&
    EQUALS      // == === != !==
    LESSGREATER // > >= < <=
    SUM         // + -
    PRODUCT     // * / %
    POWER       // **
    PREFIX      // -X !X ++X --X typeof
    POSTFIX     // X++ X--
    CALL        // fn()
    MEMBER      // obj.prop obj[prop]
)
```

**Pontos importantes:**

- `parseExpressionList()` usa precedência `ASSIGN` para evitar conflitos com vírgula
- Arrow functions são detectadas em `parseGroupedExpression()` e `parseIdentifier()`
- Suporta async functions com `parseAsyncFunction()`

### 5. Object System (`object/object.go`)

Define os tipos de runtime:

```go
type Object interface {
    Type() ObjectType
    Inspect() string
}
```

**Tipos implementados:**

- Primitivos: `Number`, `String`, `Boolean`, `Null`, `Undefined`
- Estruturais: `Array`, `ObjectMap`, `Function`, `Builtin`
- Especiais: `ReturnValue`, `Error`, `Break`, `Continue`
- Adicionais: `Buffer`, `Promise`, `Symbol`, `BigInt`, `Date`, `RegExp`, `Map`, `Set`

### 6. Environment (`object/environment.go`)

Gerencia escopo de variáveis:

```go
type Environment struct {
    store  map[string]Object
    consts map[string]bool
    outer  *Environment
}
```

**Métodos:**

- `Get(name)` - busca variável (sobe na cadeia de escopos)
- `Set(name, val)` - define variável
- `SetConst(name, val)` - define constante
- `Update(name, val)` - atualiza variável existente
- `IsConst(name)` - verifica se é constante

### 7. Evaluator (`evaluator/evaluator.go`)

Função principal `Eval(node, env)` que despacha para funções específicas:

**Singletons reutilizados:**

```go
var (
    NULL      = &object.Null{}
    UNDEFINED = &object.Undefined{}
    TRUE      = &object.Boolean{Value: true}
    FALSE     = &object.Boolean{Value: false}
    BREAK     = &object.Break{}
    CONTINUE  = &object.Continue{}
)
```

**`evalProgram` trata:**

- Import statements (`evalImportStatement`)
- Export statements (`evalExportStatement`)
- Statements regulares

**Helpers importantes:**

- `isTruthy(obj)` - converte para boolean
- `isError(obj)` - verifica se é erro
- `objectToString(obj)` - converte para string
- `nativeBoolToBooleanObject(bool)` - converte bool Go para Boolean Nulang

### 8. Módulos (`evaluator/modules.go`)

Sistema de módulos CommonJS e ES6:

**Módulos built-in:**

```go
var builtinModules = map[string]func() *object.ObjectMap{
    "fs":     initFsModule,
    "path":   initPathModule,
    "crypto": initCryptoModule,
    "os":     initOsModule,
}
```

**Import ES6 suportado:**

```javascript
import fs from "fs"; // default import
import { join } from "path"; // named imports
import * as os from "os"; // namespace import
import "module"; // side-effect import
```

**CommonJS suportado:**

```javascript
const fs = require("fs");
exports.foo = bar;
module.exports = { ... };
```

**Variáveis de módulo:**

- `__filename` - caminho absoluto do arquivo
- `__dirname` - diretório do arquivo
- `exports` - objeto de exports
- `module.exports` - exports do módulo

### 9. Módulo fs (`evaluator/fs.go`)

Operações de arquivo síncronas:

| Função                          | Descrição            |
| ------------------------------- | -------------------- |
| `readFileSync(path, encoding?)` | Ler arquivo          |
| `writeFileSync(path, data)`     | Escrever arquivo     |
| `appendFileSync(path, data)`    | Adicionar ao arquivo |
| `existsSync(path)`              | Verificar existência |
| `unlinkSync(path)`              | Deletar arquivo      |
| `mkdirSync(path, {recursive})`  | Criar diretório      |
| `rmdirSync(path, {recursive})`  | Remover diretório    |
| `readdirSync(path)`             | Listar diretório     |
| `statSync(path)`                | Info do arquivo      |
| `renameSync(old, new)`          | Renomear             |
| `copyFileSync(src, dest)`       | Copiar               |

**Módulo path:**

| Função                 | Descrição            |
| ---------------------- | -------------------- |
| `join(...paths)`       | Juntar paths         |
| `resolve(...paths)`    | Resolver absoluto    |
| `dirname(path)`        | Diretório pai        |
| `basename(path, ext?)` | Nome do arquivo      |
| `extname(path)`        | Extensão             |
| `parse(path)`          | Analisar path        |
| `isAbsolute(path)`     | É absoluto?          |
| `sep`                  | Separador do sistema |

### 10. Módulo crypto (`evaluator/crypto.go`)

| Função                  | Descrição                              |
| ----------------------- | -------------------------------------- |
| `createHash(algo)`      | Criar hash (md5, sha1, sha256, sha512) |
| `createHmac(algo, key)` | Criar HMAC                             |
| `randomBytes(size)`     | Gerar bytes aleatórios                 |
| `randomUUID()`          | Gerar UUID v4                          |

**Uso:**

```javascript
let hash = crypto.createHash("sha256");
hash.update("data");
console.log(hash.digest("hex"));
```

### 11. Buffer (`evaluator/buffer.go`)

| Método                         | Descrição                     |
| ------------------------------ | ----------------------------- |
| `Buffer.from(data, encoding?)` | Criar de string/array         |
| `Buffer.alloc(size, fill?)`    | Alocar buffer                 |
| `Buffer.concat(list)`          | Concatenar buffers            |
| `Buffer.isBuffer(obj)`         | Verificar tipo                |
| `.toString(encoding?)`         | Converter (utf8, hex, base64) |
| `.slice(start, end)`           | Fatiar                        |
| `.fill(value)`                 | Preencher                     |
| `.copy(target, ...)`           | Copiar para outro buffer      |
| `.equals(other)`               | Comparar                      |
| `.indexOf(value)`              | Buscar                        |

### 12. Promise (`evaluator/promise.go`)

Promises síncronas (não-async):

| Método                            | Descrição          |
| --------------------------------- | ------------------ |
| `Promise.resolve(value)`          | Criar fulfilled    |
| `Promise.reject(reason)`          | Criar rejected     |
| `Promise.all(array)`              | Todos fulfill      |
| `Promise.race(array)`             | Primeiro settled   |
| `Promise.allSettled(array)`       | Status de todos    |
| `.then(onFulfilled, onRejected?)` | Handler de sucesso |
| `.catch(onRejected)`              | Handler de erro    |
| `.finally(onFinally)`             | Handler final      |

### 13. Módulo os (`evaluator/os.go`)

| Função/Propriedade | Descrição                       |
| ------------------ | ------------------------------- |
| `platform()`       | "darwin", "linux", "windows"    |
| `arch()`           | "arm64", "amd64", etc.          |
| `homedir()`        | Diretório home                  |
| `tmpdir()`         | Diretório temporário            |
| `hostname()`       | Nome do host                    |
| `type()`           | "Darwin", "Linux", "Windows_NT" |
| `cpus()`           | Array de CPUs                   |
| `EOL`              | "\n" ou "\r\n"                  |
| `devNull`          | "/dev/null" ou "nul"            |

## Features da Linguagem

### Variáveis

```javascript
let x = 5; // block-scoped, mutável
const y = 10; // block-scoped, imutável
var z = 15; // function-scoped
```

### Tipos de Dados

```javascript
// Primitivos
42, 3.14         // number
"hello", 'hi'    // string
true, false      // boolean
null             // null
undefined        // undefined

// Estruturais
[1, 2, 3]                    // array
{a: 1, b: 2}                 // object
function(x) { return x; }    // function
(x) => x * 2                 // arrow function
Buffer.from("data")          // buffer
Promise.resolve(42)          // promise
```

### Operadores

```javascript
// Aritméticos
+ - * / % **

// Comparação
== != === !== < > <= >=

// Lógicos
&& || !

// Especiais
?? (nullish coalescing)
?: (ternário)
typeof
++ --
```

### Controle de Fluxo

```javascript
if (cond) { } else if (cond) { } else { }
for (let i = 0; i < n; i++) { }
while (cond) { }
break; continue;
try { } catch (e) { } finally { }
throw "error";
```

### Funções

```javascript
function add(a, b) {
  return a + b;
}
const mul = (a, b) => a * b;
const inc = (x) => x + 1;
async function fetch() {}
```

### Métodos de Array

```javascript
arr.push(x), arr.pop();
arr.shift(), arr.unshift(x);
arr.map(fn), arr.filter(fn), arr.reduce(fn, init);
arr.forEach(fn), arr.find(fn), arr.findIndex(fn);
arr.includes(x), arr.indexOf(x);
arr.slice(start, end), arr.concat(other);
arr.join(sep), arr.reverse();
```

### Métodos de String

```javascript
str.length;
str.toUpperCase(), str.toLowerCase();
str.trim(), str.split(sep);
str.charAt(i), str.indexOf(sub);
str.slice(start, end), str.substring(start, end);
str.includes(sub), str.startsWith(s), str.endsWith(s);
str.replace(s, r), str.repeat(n);
```

### Objetos Globais

```javascript
console.log(), console.error(), console.warn()
Math.PI, Math.E, Math.abs(), Math.ceil(), Math.floor()
Math.round(), Math.max(), Math.min(), Math.pow(), Math.sqrt()
Math.random(), Math.sin(), Math.cos(), Math.log()
JSON.stringify(), JSON.parse()
Array.isArray(), Array.from()
Object.keys(), Object.values(), Object.entries(), Object.assign()
parseInt(), parseFloat(), isNaN(), isFinite()
typeof, String(), Number(), Boolean()
```

## Compilação e Execução

```bash
# Compilar
go build -o nulang .

# Executar arquivo
./nulang arquivo.nu

# REPL interativo
./nulang
```

## Convenções de Código

1. **Nomes de arquivos**: snake_case para Go, kebab-case ou snake_case para exemplos
2. **Funções exportadas**: CamelCase (Go convention)
3. **Variáveis locais**: camelCase
4. **Constantes**: UPPER_SNAKE_CASE para tokens

## Padrões de Desenvolvimento

### Adicionando novo operador

1. Definir token em `token/token.go`
2. Adicionar ao lexer em `lexer/lexer.go`
3. Definir precedência em `parser/parser.go`
4. Registrar parse function (prefix ou infix)
5. Implementar avaliação em `evaluator/operators.go`

### Adicionando novo built-in

1. Criar função em `evaluator/builtins.go`
2. Registrar no mapa `builtins` dentro de `initBuiltins()`
3. Se for módulo, criar arquivo separado e adicionar em `builtinModules`

### Adicionando método de objeto

1. Para Array: `evaluator/array_methods.go` → `evalArrayProperty()`
2. Para String: `evaluator/string_methods.go` → `evalStringProperty()`
3. Para Buffer: `evaluator/buffer.go` → `evalBufferProperty()`
4. Para Promise: `evaluator/promise.go` → `evalPromiseProperty()`

### Adicionando novo statement

1. Definir struct em `ast/ast.go`
2. Implementar `statementNode()`, `TokenLiteral()`, `String()`
3. Adicionar keyword em `token/token.go`
4. Adicionar parsing em `parser/parser.go` → `parseStatement()`
5. Implementar avaliação em `evaluator/evaluator.go` → `Eval()`

## Limitações Conhecidas

1. **Promises são síncronas** - não há event loop real
2. **Async/await** - parseado mas não executa assincronamente
3. **RegExp** - tipo definido mas não implementado
4. **Date** - tipo definido mas não implementado
5. **Map/Set** - tipos definidos mas não implementados
6. **Classes** - não implementadas
7. **HTTP** - não implementado
8. **Timers** - setTimeout/setInterval não implementados

## Testes

Executar exemplos para verificar funcionalidade:

```bash
./nulang examples/example.nu       # Features básicas
./nulang examples/new_features.nu  # fs, crypto, buffer, promise
./nulang examples/test_es6.nu      # Import ES6
./nulang examples/test_modules.nu  # Sistema de módulos
```

## Histórico de Implementação

1. **Fase 1**: Token, Lexer, AST básica
2. **Fase 2**: Parser com Pratt Parser
3. **Fase 3**: Object system e Environment
4. **Fase 4**: Evaluator com operadores e controle de fluxo
5. **Fase 5**: Métodos de Array e String
6. **Fase 6**: Built-ins (console, Math, JSON)
7. **Fase 7**: Sistema de módulos (require, exports)
8. **Fase 8**: Módulos fs, path, crypto
9. **Fase 9**: Buffer e Promise
10. **Fase 10**: Import ES6 e módulo os

---

_Documento gerado em 25 de Dezembro de 2024_
