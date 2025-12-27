# OS Module

O módulo `os` fornece informações sobre o sistema operacional.

## Importação

```javascript
const os = require("os");
// ou
import os from "os";
```

## Métodos

### os.platform()

Retorna a plataforma do sistema operacional.

```javascript
console.log(os.platform());
// "darwin" (macOS), "linux", "windows"
```

### os.arch()

Retorna a arquitetura da CPU.

```javascript
console.log(os.arch());
// "x64", "arm64", "arm", etc.
```

### os.homedir()

Retorna o diretório home do usuário.

```javascript
console.log(os.homedir());
// "/Users/admin" (macOS)
// "C:\\Users\\admin" (Windows)
// "/home/admin" (Linux)
```

### os.tmpdir()

Retorna o diretório de arquivos temporários.

```javascript
console.log(os.tmpdir());
// "/tmp" (Unix)
// "C:\\Users\\admin\\AppData\\Local\\Temp" (Windows)
```

### os.hostname()

Retorna o hostname da máquina.

```javascript
console.log(os.hostname());
// "my-macbook.local"
```

### os.type()

Retorna o nome do sistema operacional.

```javascript
console.log(os.type());
// "Darwin" (macOS)
// "Linux"
// "Windows_NT"
```

### os.cpus()

Retorna informações sobre as CPUs.

```javascript
const cpus = os.cpus();
console.log(`Número de CPUs: ${cpus.length}`);

cpus.forEach((cpu, i) => {
  console.log(`CPU ${i}: ${cpu.model} @ ${cpu.speed}MHz`);
});
```

**Retorno**: Array de objetos com:
| Propriedade | Tipo | Descrição |
|-------------|------|-----------|
| `model` | String | Modelo da CPU |
| `speed` | Number | Velocidade em MHz |

### os.uptime()

Retorna o tempo de atividade do sistema em segundos.

```javascript
console.log(`Uptime: ${os.uptime()} segundos`);
```

### os.freemem()

Retorna a memória livre em bytes.

```javascript
const free = os.freemem();
console.log(`Memória livre: ${free / 1024 / 1024} MB`);
```

### os.totalmem()

Retorna a memória total em bytes.

```javascript
const total = os.totalmem();
console.log(`Memória total: ${total / 1024 / 1024 / 1024} GB`);
```

## Constantes

### os.EOL

Caractere de fim de linha do sistema.

```javascript
console.log(os.EOL === "\n"); // true no Unix
console.log(os.EOL === "\r\n"); // true no Windows

// Usar ao escrever arquivos
const lines = ["linha1", "linha2", "linha3"];
const content = lines.join(os.EOL);
```

### os.devNull

Caminho para o dispositivo null.

```javascript
console.log(os.devNull);
// "/dev/null" (Unix)
// "nul" (Windows)
```

## Exemplos Práticos

### Informações do Sistema

```javascript
function systemInfo() {
  console.log("=== Informações do Sistema ===");
  console.log(`SO: ${os.type()} ${os.arch()}`);
  console.log(`Plataforma: ${os.platform()}`);
  console.log(`Hostname: ${os.hostname()}`);
  console.log(`CPUs: ${os.cpus().length}`);
  console.log(`Home: ${os.homedir()}`);
  console.log(`Temp: ${os.tmpdir()}`);
}

systemInfo();
```

### Verificar Sistema Operacional

```javascript
function isWindows() {
  return os.platform() === "windows";
}

function isMac() {
  return os.platform() === "darwin";
}

function isLinux() {
  return os.platform() === "linux";
}

// Uso
if (isMac()) {
  console.log("Executando no macOS");
} else if (isWindows()) {
  console.log("Executando no Windows");
} else if (isLinux()) {
  console.log("Executando no Linux");
}
```

### Caminho Multiplataforma

```javascript
const path = require("path");

function getConfigPath() {
  const home = os.homedir();

  if (os.platform() === "windows") {
    return path.join(home, "AppData", "Local", "MeuApp");
  } else if (os.platform() === "darwin") {
    return path.join(home, "Library", "Application Support", "MeuApp");
  } else {
    return path.join(home, ".config", "meuapp");
  }
}

console.log(getConfigPath());
```

### Monitor de Recursos

```javascript
function resourceMonitor() {
  const total = os.totalmem();
  const free = os.freemem();
  const used = total - free;
  const usagePercent = ((used / total) * 100).toFixed(1);

  console.log(`Memória: ${usagePercent}% em uso`);
  console.log(`  Total: ${(total / 1024 / 1024 / 1024).toFixed(2)} GB`);
  console.log(`  Usado: ${(used / 1024 / 1024 / 1024).toFixed(2)} GB`);
  console.log(`  Livre: ${(free / 1024 / 1024 / 1024).toFixed(2)} GB`);
}

resourceMonitor();
```

### Criar Arquivo Temporário

```javascript
const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

function createTempFile(content) {
  const tmpDir = os.tmpdir();
  const filename = `temp_${crypto.randomUUID()}.tmp`;
  const filepath = path.join(tmpDir, filename);

  fs.writeFileSync(filepath, content);
  return filepath;
}

const tempPath = createTempFile("Conteúdo temporário");
console.log(`Arquivo criado: ${tempPath}`);
```

### Info para Logs

```javascript
function getSystemContext() {
  return {
    platform: os.platform(),
    arch: os.arch(),
    hostname: os.hostname(),
    cpus: os.cpus().length,
    nodeVersion: "nulang",
    timestamp: new Date().toISOString(),
  };
}

function logWithContext(message) {
  const context = getSystemContext();
  console.log(`[${context.timestamp}] [${context.hostname}] ${message}`);
}

logWithContext("Aplicação iniciada");
```

## Tabela de Plataformas

| os.platform() | os.type()    | Sistema |
| ------------- | ------------ | ------- |
| `darwin`      | `Darwin`     | macOS   |
| `linux`       | `Linux`      | Linux   |
| `windows`     | `Windows_NT` | Windows |

## Tabela de Arquiteturas

| os.arch() | Descrição  |
| --------- | ---------- |
| `x64`     | 64-bit x86 |
| `x86`     | 32-bit x86 |
| `arm64`   | 64-bit ARM |
| `arm`     | 32-bit ARM |

## Veja Também

- [Process](./process.md) - Objeto process
- [File System](./filesystem.md) - Operações com arquivos
