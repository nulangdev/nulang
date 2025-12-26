# Packages e Dependencies

O sistema de módulos Nulang suporta tanto módulos internos quanto módulos de terceiros.

## Módulos Built-in

Módulos que vêm integrados com Nulang:

| Módulo        | Descrição                          |
| ------------- | ---------------------------------- |
| `fs`          | Operações com sistema de arquivos  |
| `path`        | Manipulação de caminhos            |
| `crypto`      | Funções criptográficas             |
| `os`          | Informações do sistema operacional |
| `http`        | Cliente e servidor HTTP            |
| `stream`      | Streams de dados                   |
| `url`         | Parsing de URLs                    |
| `querystring` | Parsing de query strings           |
| `events`      | EventEmitter                       |

```javascript
// Importar módulo built-in
const fs = require("fs");
const path = require("path");

// ou com ES6 imports
import fs from "fs";
import { join } from "path";
```

## Estrutura de Pacotes

### Diretório .nu_modules

Pacotes de terceiros são armazenados em `.nu_modules`:

```
projeto/
├── .nu_modules/
│   ├── meu-pacote/
│   │   ├── index.nu
│   │   └── lib/
│   │       └── utils.nu
│   └── outro-pacote/
│       └── index.nu
├── src/
│   └── main.nu
└── index.nu
```

### Estrutura de um Pacote

Um pacote deve ter ao menos um `index.nu`:

```
meu-pacote/
├── index.nu      # Entrada principal
├── lib/          # Código fonte
│   ├── utils.nu
│   └── helpers.nu
└── README.md     # Documentação
```

## Importando Pacotes

### De .nu_modules

```javascript
// Importa .nu_modules/meu-pacote/index.nu
const pkg = require("meu-pacote");
import pkg from "meu-pacote";
```

### Caminhos Relativos

```javascript
// Arquivo no mesmo diretório
const utils = require("./utils");
import utils from "./utils";

// Arquivo em outro diretório
const lib = require("../lib/helper");
import { helper } from "../lib/helper";
```

## Criando um Módulo

### Exportando com CommonJS

```javascript
// utils.nu

function add(a, b) {
  return a + b;
}

function multiply(a, b) {
  return a * b;
}

// Exportar individualmente
exports.add = add;
exports.multiply = multiply;

// Ou exportar objeto
module.exports = {
  add: add,
  multiply: multiply,
};
```

### Exportando com ES6

```javascript
// utils.nu

export function add(a, b) {
  return a + b;
}

export function multiply(a, b) {
  return a * b;
}

// Export default
export default {
  add: add,
  multiply: multiply,
};
```

## Importando um Módulo

### Importação CommonJS

```javascript
// Importar módulo completo
const utils = require("./utils");
console.log(utils.add(1, 2)); // 3

// Desestruturação
const { add, multiply } = require("./utils");
console.log(add(1, 2)); // 3
```

### Importação ES6

```javascript
// Importar default export
import utils from "./utils";

// Importar named exports
import { add, multiply } from "./utils";

// Importar tudo
import * as utils from "./utils";

// Importar default e named
import utils, { add } from "./utils";
```

## Resolução de Módulos

O sistema de módulos resolve caminhos na seguinte ordem:

1. **Módulos built-in**: `fs`, `path`, `http`, etc.
2. **Caminho relativo**: `./`, `../`
3. **Diretório .nu_modules**: `nome-pacote`

### Extensões

Se não especificada, Nulang tenta:

1. `nome.nu`
2. `nome/index.nu`

```javascript
// Estes são equivalentes
require("./utils");
require("./utils.nu");
```

## Variáveis de Módulo

### module.exports

Objeto que será exportado.

```javascript
module.exports = {
  name: "Meu Módulo",
  version: "1.0.0",
  doSomething: function () {
    return "done";
  },
};
```

### exports

Atalho para `module.exports`.

```javascript
exports.name = "Meu Módulo";
exports.doSomething = function () {
  return "done";
};
```

**Nota**: Não reatribua `exports` diretamente!

```javascript
// ❌ Errado
exports = { name: "test" };

// ✅ Correto
module.exports = { name: "test" };
// ou
exports.name = "test";
```

### \_\_filename

Caminho absoluto do arquivo atual.

```javascript
console.log(__filename);
// /Users/admin/projeto/src/utils.nu
```

### \_\_dirname

Diretório do arquivo atual.

```javascript
console.log(__dirname);
// /Users/admin/projeto/src
```

## Exemplos de Pacotes

### Pacote de Utilitários

```javascript
// .nu_modules/my-utils/index.nu

export function capitalize(str) {
  if (str.length === 0) return str;
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
}

export function slugify(str) {
  return str.toLowerCase().replaceAll(" ", "-");
}

export function truncate(str, len) {
  if (str.length <= len) return str;
  return str.slice(0, len - 3) + "...";
}
```

### Pacote de Validação

```javascript
// .nu_modules/validator/index.nu

export function isEmail(str) {
  return str.includes("@") && str.includes(".");
}

export function isNumeric(str) {
  return !isNaN(parseFloat(str));
}

export function minLength(str, min) {
  return str.length >= min;
}

export function maxLength(str, max) {
  return str.length <= max;
}
```

### Pacote de Logger

```javascript
// .nu_modules/logger/index.nu
const fs = require("fs");

class Logger {
  constructor(filepath) {
    this.filepath = filepath;
  }

  log(level, message) {
    const timestamp = new Date().toISOString();
    const line = `[${timestamp}] [${level}] ${message}\n`;
    fs.appendFileSync(this.filepath, line);
    console.log(line.trim());
  }

  info(message) {
    this.log("INFO", message);
  }

  warn(message) {
    this.log("WARN", message);
  }

  error(message) {
    this.log("ERROR", message);
  }
}

export default Logger;
```

## Uso no Projeto

```javascript
// main.nu
import { capitalize, slugify } from "my-utils";
import { isEmail, minLength } from "validator";
import Logger from "logger";

const logger = new Logger("app.log");

function validateUser(user) {
  if (!isEmail(user.email)) {
    logger.error("Email inválido");
    return false;
  }
  if (!minLength(user.password, 8)) {
    logger.error("Senha muito curta");
    return false;
  }
  logger.info("Usuário validado");
  return true;
}

const username = capitalize("joão");
const slug = slugify("Meu Artigo Legal");

console.log(username); // "João"
console.log(slug); // "meu-artigo-legal"
```

## Boas Práticas

### 1. Um Conceito por Módulo

```javascript
// ✅ Bom - módulo focado
// math.nu
export function add(a, b) {
  return a + b;
}
export function subtract(a, b) {
  return a - b;
}

// ❌ Ruim - módulo com responsabilidades misturadas
// utils.nu
export function add(a, b) {
  return a + b;
}
export function sendEmail(to, msg) {
  /* ... */
}
```

### 2. Exportar no Final

```javascript
// ✅ Bom
function privateHelper() {
  /* ... */
}

function publicFunction() {
  privateHelper();
}

export { publicFunction };
```

### 3. Nomear Exports Claramente

```javascript
// ✅ Bom
export function validateEmail(email) {
  /* ... */
}

// ❌ Ruim
export function validate(x) {
  /* ... */
}
```

## Veja Também

- [Modules](./modules.md) - Sistema de módulos detalhado
- [File System](./filesystem.md) - Operações com arquivos
