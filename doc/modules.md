# Modules

O sistema de módulos em Nulang suporta tanto CommonJS quanto ES6 imports.

## ES6 Imports

### Import Default

```javascript
import modulo from "./modulo";
import fs from "fs";
```

### Import Named

```javascript
import { funcao1, funcao2 } from "./modulo";
import { readFileSync, writeFileSync } from "fs";
```

### Import Namespace

```javascript
import * as modulo from "./modulo";
console.log(modulo.funcao1());
```

### Import Combinado

```javascript
import modulo, { funcao1, funcao2 } from "./modulo";
```

### Import Side Effect

```javascript
import "./polyfills"; // Executa o módulo sem importar nada
```

## ES6 Exports

### Export Named

```javascript
// utils.nu
export function add(a, b) {
  return a + b;
}

export const PI = 3.14159;

export class Calculator {
  // ...
}
```

### Export Default

```javascript
// Calculator.nu
export default class Calculator {
  add(a, b) {
    return a + b
  }
}

// Ou com expressão
export default function(x) {
  return x * 2
}
```

### Re-export

```javascript
// index.nu
export { add, subtract } from "./math";
export { default as Calculator } from "./Calculator";
```

## CommonJS

### require()

```javascript
const modulo = require("./modulo");
const fs = require("fs");

// Desestruturação
const { add, subtract } = require("./math");
```

### module.exports

```javascript
// Exportar objeto
module.exports = {
  add: function (a, b) {
    return a + b;
  },
  subtract: function (a, b) {
    return a - b;
  },
};

// Exportar classe
module.exports = class Calculator {
  // ...
};

// Exportar função
module.exports = function (x) {
  return x * 2;
};
```

### exports

```javascript
exports.add = function (a, b) {
  return a + b;
};
exports.subtract = function (a, b) {
  return a - b;
};
exports.PI = 3.14159;
```

## Resolução de Módulos

### Ordem de Resolução

1. **Módulos built-in**: `fs`, `path`, `http`, etc.
2. **Caminho relativo/absoluto**: `./`, `../`, `/`
3. **Pacotes em .nu_modules**: `nome-pacote`

### Extensões Tentadas

```javascript
require("./utils");
// Tenta:
// 1. ./utils.nu
// 2. ./utils/index.nu
```

### Módulos Built-in

| Módulo        | Descrição               |
| ------------- | ----------------------- |
| `fs`          | Sistema de arquivos     |
| `path`        | Manipulação de caminhos |
| `crypto`      | Criptografia            |
| `os`          | Sistema operacional     |
| `http`        | HTTP cliente/servidor   |
| `stream`      | Streams de dados        |
| `url`         | Parsing de URLs         |
| `querystring` | Query strings           |
| `events`      | EventEmitter            |

## Variáveis de Módulo

### \_\_filename

Caminho absoluto do arquivo atual.

```javascript
console.log(__filename);
// /Users/admin/projeto/src/main.nu
```

### \_\_dirname

Diretório do arquivo atual.

```javascript
console.log(__dirname);
// /Users/admin/projeto/src
```

### module

Objeto representando o módulo atual.

```javascript
console.log(module.exports); // Objeto de exportação
```

### exports

Atalho para `module.exports`.

```javascript
exports.myFunction = function () {};
// Equivalente a:
module.exports.myFunction = function () {};
```

## Cache de Módulos

Módulos são cacheados após a primeira importação:

```javascript
// Primeira vez: executa o módulo
const a = require("./modulo");

// Segunda vez: retorna do cache
const b = require("./modulo");

console.log(a === b); // true
```

## Padrões de Organização

### Barrel Pattern

Agrupa exports em um único arquivo:

```javascript
// lib/index.nu
export { default as Calculator } from "./Calculator";
export { add, subtract } from "./math";
export { format } from "./utils";

// Uso
import { Calculator, add, format } from "./lib";
```

### Module Factory

```javascript
// logger.nu
export default function createLogger(prefix) {
  return {
    log: function (msg) {
      console.log(`[${prefix}] ${msg}`);
    },
    error: function (msg) {
      console.error(`[${prefix}] ERROR: ${msg}`);
    },
  };
}

// Uso
import createLogger from "./logger";
const log = createLogger("APP");
log.log("Iniciando...");
```

### Singleton Pattern

```javascript
// database.nu
let instance = null;

class Database {
  constructor() {
    if (instance) {
      return instance;
    }
    this.connection = null;
    instance = this;
  }

  connect(url) {
    this.connection = { url: url };
  }
}

export default new Database();

// Uso - sempre a mesma instância
import db from "./database";
db.connect("localhost:5432");
```

## Exemplos Práticos

### Organização de Projeto

```
projeto/
├── index.nu              # Entrada principal
├── src/
│   ├── controllers/
│   │   ├── index.nu      # Barrel
│   │   ├── userController.nu
│   │   └── postController.nu
│   ├── services/
│   │   ├── index.nu
│   │   ├── userService.nu
│   │   └── emailService.nu
│   └── utils/
│       ├── index.nu
│       ├── validation.nu
│       └── formatting.nu
└── .nu_modules/
    └── ...
```

### Arquivo Barrel

```javascript
// src/utils/index.nu
export { validateEmail, validatePhone } from "./validation";
export { formatDate, formatCurrency } from "./formatting";
```

### Importação Limpa

```javascript
// index.nu
import { UserController, PostController } from "./src/controllers";
import { formatDate, validateEmail } from "./src/utils";

const user = UserController.create({
  email: "test@test.com",
  createdAt: formatDate(new Date()),
});
```

## Ciclos de Dependência

Evite ciclos de dependência:

```javascript
// ❌ Ruim - ciclo
// a.nu
import { b } from "./b";
export const a = "a";

// b.nu
import { a } from "./a"; // a ainda não foi definido!
export const b = "b";
```

Solução: reorganize ou use injeção de dependência.

## Boas Práticas

### 1. Prefira Named Exports

```javascript
// ✅ Bom - claro e refatorável
export function processData(data) {}
export function validateData(data) {}

// Uso
import { processData, validateData } from "./data";
```

### 2. Um Export Default por Arquivo

```javascript
// ✅ Bom
// Calculator.nu
export default class Calculator {}
```

### 3. Imports no Topo

```javascript
// ✅ Bom
import fs from "fs";
import { helper } from "./utils";

const data = fs.readFileSync("file.txt");
```

### 4. Agrupe Imports

```javascript
// ✅ Bom
// Built-ins primeiro
import fs from "fs";
import path from "path";

// Depois pacotes externos
import lodash from "lodash";

// Por último, módulos locais
import { utils } from "./utils";
```

## Veja Também

- [Packages](./packages.md) - Sistema de pacotes
- [File System](./filesystem.md) - Operações com arquivos
- [Classes](./classes.md) - Definição de classes
