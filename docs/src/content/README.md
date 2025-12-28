# Documentação Nulang

Bem-vindo à documentação completa do Nulang - um interpretador JavaScript/Node.js-like escrito em Go.

## Índice

### Tipos e Estruturas de Dados

| Módulo                  | Descrição                                |
| ----------------------- | ---------------------------------------- |
| [Array](./array.md)     | Listas ordenadas e métodos de array      |
| [String](./string.md)   | Manipulação de texto e template literals |
| [Numbers](./numbers.md) | Números, Math e operações matemáticas    |
| [Map](./map.md)         | Map, Set, WeakMap e WeakSet              |
| [Date](./date.md)       | Manipulação de datas e horários          |
| [RegExp](./regex.md)    | Expressões regulares                     |
| [Buffer](./buffer.md)   | Dados binários                           |
| [Blob](./blob.md)       | Dados brutos imutáveis                   |
| [File](./file.md)       | Arquivos (estende Blob)                  |

### Programação Orientada a Objetos

| Módulo                                  | Descrição                       |
| --------------------------------------- | ------------------------------- |
| [Classes](./classes.md)                 | Declaração de classes e herança |
| [Getters/Setters](./getters_setters.md) | Propriedades computadas         |
| [Decorators](./decorators.md)           | Decorators de classe e método   |

### Programação Assíncrona

| Módulo                        | Descrição                             |
| ----------------------------- | ------------------------------------- |
| [Promise](./promise.md)       | Operações assíncronas                 |
| [Events](./events.md)         | EventEmitter e padrão de eventos      |
| [Timers](./timers.md)         | setTimeout, setInterval, setImmediate |
| [Event Loop](./event_loop.md) | Modelo de execução assíncrona         |

### Tratamento de Erros

| Módulo              | Descrição                             |
| ------------------- | ------------------------------------- |
| [Error](./error.md) | Classe Error e tratamento de exceções |

### Streams

| Módulo                | Descrição                                  |
| --------------------- | ------------------------------------------ |
| [Stream](./stream.md) | Readable, Writable, Transform, PassThrough |

### Sistema de Arquivos

| Módulo                         | Descrição                                     |
| ------------------------------ | --------------------------------------------- |
| [File System](./filesystem.md) | Operações com arquivos (fs) e caminhos (path) |

### Rede

| Módulo              | Descrição                                        |
| ------------------- | ------------------------------------------------ |
| [HTTP](./http.md)   | Cliente e servidor HTTP, fetch, url, querystring |
| [Net](./net.md)     | Servidores e clientes TCP                        |
| [Dgram](./dgram.md) | Sockets UDP                                      |
| [DNS](./dns.md)     | Resolução de nomes de domínio                    |

### Sistema

| Módulo                              | Descrição                          |
| ----------------------------------- | ---------------------------------- |
| [Process](./process.md)             | Objeto global process              |
| [OS](./os.md)                       | Informações do sistema operacional |
| [Crypto](./crypto.md)               | Funções criptográficas             |
| [Child Process](./child_process.md) | Execução de processos filhos       |
| [Zlib](./zlib.md)                   | Compressão e descompressão         |
| [VM](./vm.md)                       | Execução de código em sandbox      |

### Entrada/Saída Interativa

| Módulo                    | Descrição                          |
| ------------------------- | ---------------------------------- |
| [Readline](./readline.md) | Leitura de entrada linha por linha |

### Utilitários e Testes

| Módulo                | Descrição                        |
| --------------------- | -------------------------------- |
| [Util](./util.md)     | Funções utilitárias e formatação |
| [Assert](./assert.md) | Funções de asserção para testes  |

### Metaprogramação

| Módulo                  | Descrição                             |
| ----------------------- | ------------------------------------- |
| [Proxy](./proxy.md)     | Interceptação de operações em objetos |
| [Reflect](./reflect.md) | API Reflect para operações de objetos |

### Módulos e Dependências

| Módulo                    | Descrição                               |
| ------------------------- | --------------------------------------- |
| [Modules](./modules.md)   | Sistema de módulos (ES6 e CommonJS)     |
| [Packages](./packages.md) | Gerenciamento de pacotes e dependências |

---

## Início Rápido

### Instalação

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

### Iniciar novo projeto

```bash
nulang init
```

### Executar Script

```bash
nulang script.nu
```

### Instalar dependências

```bash
nulang install
```

### Modo Watch

```bash
nulang script.nu --watch
```

### REPL

```bash
nulang
```

---

## Exemplos Básicos

### Hello World

```javascript
console.log("Hello, World!");
```

### Variáveis

```javascript
let nome = "João";
const idade = 30;
var ativo = true;
```

### Funções

```javascript
function somar(a, b) {
  return a + b;
}

const multiplicar = (a, b) => a * b;
```

### Classes

```javascript
class Pessoa {
  constructor(nome) {
    this.nome = nome;
  }

  saudar() {
    return `Olá, ${this.nome}`;
  }
}

const p = new Pessoa("Maria");
console.log(p.saudar());
```

### Módulos

```javascript
// math.nu
export function add(a, b) {
  return a + b;
}

// main.nu
import { add } from "./math";
console.log(add(1, 2));
```

### Async

```javascript
setTimeout(() => {
  console.log("Após 1 segundo");
}, 1000);

fetch("https://api.example.com/data")
  .then((res) => res.json())
  .then((data) => console.log(data));
```

---

## Objetos Globais

| Objeto         | Descrição                      |
| -------------- | ------------------------------ |
| `console`      | Logging (log, error, warn)     |
| `Math`         | Funções matemáticas            |
| `JSON`         | Serialização JSON              |
| `Array`        | Construtor e métodos estáticos |
| `Object`       | Construtor e métodos estáticos |
| `Date`         | Construtor de datas            |
| `RegExp`       | Construtor de regex            |
| `Map`, `Set`   | Coleções                       |
| `Promise`      | Operações assíncronas          |
| `Buffer`       | Dados binários                 |
| `Blob`, `File` | Dados brutos                   |
| `EventEmitter` | Eventos                        |
| `process`      | Informações do processo        |

## Funções Globais

| Função           | Descrição                |
| ---------------- | ------------------------ |
| `print()`        | Imprime no console       |
| `parseInt()`     | Converte para inteiro    |
| `parseFloat()`   | Converte para float      |
| `isNaN()`        | Verifica se é NaN        |
| `isFinite()`     | Verifica se é finito     |
| `String()`       | Converte para string     |
| `Number()`       | Converte para número     |
| `Boolean()`      | Converte para boolean    |
| `setTimeout()`   | Timeout                  |
| `setInterval()`  | Intervalo                |
| `setImmediate()` | Execução imediata        |
| `fetch()`        | Requisição HTTP          |
| `require()`      | Importar módulo CommonJS |
| `sleep()`        | Pausa síncrona           |

---

## Módulos Built-in

```javascript
const fs = require("fs");
const path = require("path");
const crypto = require("crypto");
const os = require("os");
const http = require("http");
const stream = require("stream");
const url = require("url");
const querystring = require("querystring");
const events = require("events");
const net = require("net");
const dgram = require("dgram");
const dns = require("dns");
const child_process = require("child_process");
const readline = require("readline");
const vm = require("vm");
const zlib = require("zlib");
const assert = require("assert");
const util = require("util");
```

---

## Suporte

- **GitHub**: [github.com/nulang/nulang](https://github.com/nulang/nulang)
- **Extensão**: `.nu`
- **Licença**: MIT
