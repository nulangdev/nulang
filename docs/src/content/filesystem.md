# File System Module

O módulo `fs` fornece funções para operações com arquivos e diretórios, compatível com Node.js.

## Importação

```javascript
const fs = require("fs");
// ou
import fs from "fs";
```

## Operações com Arquivos

### fs.readFileSync(path, [encoding])

Lê o conteúdo de um arquivo de forma síncrona.

```javascript
// Ler como Buffer
const buffer = fs.readFileSync("arquivo.txt");

// Ler como string
const content = fs.readFileSync("arquivo.txt", "utf8");
const content2 = fs.readFileSync("arquivo.txt", { encoding: "utf8" });

console.log(content); // Conteúdo do arquivo
```

**Parâmetros**:
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `path` | String | Caminho do arquivo |
| `encoding` | String/Object | Encoding ("utf8") ou objeto com opções |

**Retorno**: Buffer ou String

### fs.writeFileSync(path, data, [encoding])

Escreve dados em um arquivo (cria ou sobrescreve).

```javascript
// Escrever string
fs.writeFileSync("arquivo.txt", "Hello, World!");

// Escrever buffer
const buffer = Buffer.from("Dados binários");
fs.writeFileSync("binario.dat", buffer);
```

**Parâmetros**:
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `path` | String | Caminho do arquivo |
| `data` | String/Buffer | Dados a escrever |
| `encoding` | String | Encoding (opcional) |

### fs.appendFileSync(path, data)

Adiciona dados ao final de um arquivo.

```javascript
fs.appendFileSync("log.txt", "Nova linha\n");
fs.appendFileSync("log.txt", "Outra linha\n");
```

### fs.existsSync(path)

Verifica se um arquivo ou diretório existe.

```javascript
if (fs.existsSync("config.json")) {
  console.log("Arquivo existe!");
} else {
  console.log("Arquivo não encontrado");
}
```

**Retorno**: Boolean

### fs.unlinkSync(path)

Remove (deleta) um arquivo.

```javascript
fs.unlinkSync("arquivo-para-deletar.txt");
```

### fs.renameSync(oldPath, newPath)

Renomeia ou move um arquivo.

```javascript
// Renomear
fs.renameSync("antigo.txt", "novo.txt");

// Mover
fs.renameSync("arquivo.txt", "pasta/arquivo.txt");
```

### fs.copyFileSync(src, dest)

Copia um arquivo.

```javascript
fs.copyFileSync("original.txt", "copia.txt");
```

## Operações com Diretórios

### fs.mkdirSync(path, [options])

Cria um diretório.

```javascript
// Criar um diretório
fs.mkdirSync("nova-pasta");

// Criar diretórios recursivamente
fs.mkdirSync("pasta/subpasta/outra", { recursive: true });
```

**Opções**:
| Opção | Tipo | Descrição |
|-------|------|-----------|
| `recursive` | Boolean | Criar diretórios pais se necessário |

### fs.rmdirSync(path, [options])

Remove um diretório.

```javascript
// Remover diretório vazio
fs.rmdirSync("pasta-vazia");

// Remover diretório com conteúdo
fs.rmdirSync("pasta-com-arquivos", { recursive: true });
```

**Opções**:
| Opção | Tipo | Descrição |
|-------|------|-----------|
| `recursive` | Boolean | Remover recursivamente |

### fs.readdirSync(path)

Lista o conteúdo de um diretório.

```javascript
const files = fs.readdirSync("./");
console.log(files); // ["arquivo1.txt", "pasta", "script.nu"]

files.forEach((file) => {
  console.log("- " + file);
});
```

**Retorno**: Array de strings (nomes de arquivos/pastas)

### fs.statSync(path)

Obtém informações sobre um arquivo ou diretório.

```javascript
const stats = fs.statSync("arquivo.txt");

console.log("Tamanho:", stats.size, "bytes");
console.log("É arquivo?", stats.isFile());
console.log("É diretório?", stats.isDirectory());
console.log("Modificado em:", stats.mtime);
```

**Retorno**: Objeto Stats

#### Objeto Stats

| Propriedade/Método | Tipo     | Descrição                       |
| ------------------ | -------- | ------------------------------- |
| `size`             | Number   | Tamanho em bytes                |
| `mtime`            | Number   | Timestamp da última modificação |
| `mode`             | Number   | Permissões do arquivo           |
| `isFile()`         | Function | Retorna true se for arquivo     |
| `isDirectory()`    | Function | Retorna true se for diretório   |

## Módulo Path

Funções para manipulação de caminhos de arquivo.

### path.join(...paths)

Une partes de um caminho.

```javascript
const path = require("path");

const filePath = path.join("pasta", "subpasta", "arquivo.txt");
console.log(filePath); // "pasta/subpasta/arquivo.txt"
```

### path.resolve(...paths)

Resolve um caminho absoluto.

```javascript
const absolute = path.resolve("pasta", "arquivo.txt");
console.log(absolute); // "/Users/admin/projeto/pasta/arquivo.txt"
```

### path.dirname(path)

Retorna o diretório de um caminho.

```javascript
console.log(path.dirname("/pasta/arquivo.txt")); // "/pasta"
```

### path.basename(path, [ext])

Retorna o nome do arquivo.

```javascript
console.log(path.basename("/pasta/arquivo.txt")); // "arquivo.txt"
console.log(path.basename("/pasta/arquivo.txt", ".txt")); // "arquivo"
```

### path.extname(path)

Retorna a extensão do arquivo.

```javascript
console.log(path.extname("arquivo.txt")); // ".txt"
console.log(path.extname("arquivo.tar.gz")); // ".gz"
```

### path.parse(path)

Parseia um caminho em componentes.

```javascript
const parsed = path.parse("/home/user/file.txt");
console.log(parsed);
// {
//   root: "/",
//   dir: "/home/user",
//   base: "file.txt",
//   ext: ".txt",
//   name: "file"
// }
```

### path.isAbsolute(path)

Verifica se é um caminho absoluto.

```javascript
console.log(path.isAbsolute("/pasta/arquivo")); // true
console.log(path.isAbsolute("./arquivo")); // false
```

### path.sep

Separador de diretórios do sistema.

```javascript
console.log(path.sep); // "/" no Unix, "\\" no Windows
```

## Exemplos Práticos

### Ler JSON

```javascript
function readJSON(filePath) {
  const content = fs.readFileSync(filePath, "utf8");
  return JSON.parse(content);
}

const config = readJSON("config.json");
console.log(config.name);
```

### Escrever JSON

```javascript
function writeJSON(filePath, data) {
  const json = JSON.stringify(data);
  fs.writeFileSync(filePath, json);
}

writeJSON("data.json", { name: "Nulang", version: "1.0" });
```

### Logger Simples

```javascript
function log(message) {
  const timestamp = new Date().toISOString();
  const line = `[${timestamp}] ${message}\n`;
  fs.appendFileSync("app.log", line);
}

log("Aplicação iniciada");
log("Processando dados...");
log("Aplicação finalizada");
```

### Listar Arquivos Recursivamente

```javascript
function listFiles(dir, prefix = "") {
  const items = fs.readdirSync(dir);

  items.forEach((item) => {
    const fullPath = path.join(dir, item);
    const stats = fs.statSync(fullPath);

    console.log(prefix + item);

    if (stats.isDirectory()) {
      listFiles(fullPath, prefix + "  ");
    }
  });
}

listFiles("./projeto");
```

### Copiar Diretório

```javascript
function copyDir(src, dest) {
  fs.mkdirSync(dest, { recursive: true });

  const items = fs.readdirSync(src);

  items.forEach((item) => {
    const srcPath = path.join(src, item);
    const destPath = path.join(dest, item);
    const stats = fs.statSync(srcPath);

    if (stats.isDirectory()) {
      copyDir(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  });
}

copyDir("origem", "destino");
```

### Limpar Diretório

```javascript
function cleanDir(dir) {
  if (!fs.existsSync(dir)) return;

  const items = fs.readdirSync(dir);

  items.forEach((item) => {
    const fullPath = path.join(dir, item);
    const stats = fs.statSync(fullPath);

    if (stats.isDirectory()) {
      cleanDir(fullPath);
      fs.rmdirSync(fullPath);
    } else {
      fs.unlinkSync(fullPath);
    }
  });
}

cleanDir("temp");
```

### Encontrar Arquivos por Extensão

```javascript
function findByExtension(dir, ext) {
  const results = [];
  const items = fs.readdirSync(dir);

  items.forEach((item) => {
    const fullPath = path.join(dir, item);
    const stats = fs.statSync(fullPath);

    if (stats.isDirectory()) {
      const subResults = findByExtension(fullPath, ext);
      subResults.forEach((r) => results.push(r));
    } else if (path.extname(item) === ext) {
      results.push(fullPath);
    }
  });

  return results;
}

const nuFiles = findByExtension("./", ".nu");
console.log(nuFiles);
```

## Variáveis Especiais de Módulo

### \_\_filename

Caminho completo do arquivo atual.

```javascript
console.log(__filename); // "/Users/admin/projeto/script.nu"
```

### \_\_dirname

Diretório do arquivo atual.

```javascript
console.log(__dirname); // "/Users/admin/projeto"

// Usado com path.join para caminhos relativos
const dataPath = path.join(__dirname, "data", "file.json");
```

## Veja Também

- [Buffer](./buffer.md) - Manipulação de dados binários
- [Modules](./modules.md) - Sistema de módulos
- [Process](./process.md) - Objeto process
