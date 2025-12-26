# Process

O objeto global `process` fornece informações e controle sobre o processo atual.

## Disponibilidade

`process` está disponível globalmente, sem necessidade de importação.

```javascript
console.log(process.platform);
```

## Propriedades

### process.argv

Array com argumentos de linha de comando.

```javascript
// Executando: nulang script.nu arg1 arg2
console.log(process.argv);
// ["nulang", "script.nu", "arg1", "arg2"]

// Pegar apenas argumentos do script
const args = process.argv.slice(2);
console.log(args); // ["arg1", "arg2"]
```

### process.env

Objeto com variáveis de ambiente.

```javascript
console.log(process.env.HOME); // "/Users/admin"
console.log(process.env.PATH); // "/usr/bin:/bin:..."
console.log(process.env.NODE_ENV); // "development" ou "production"

// Verificar se existe
if (process.env.API_KEY) {
  console.log("API Key configurada");
}
```

### process.platform

String identificando a plataforma.

```javascript
console.log(process.platform);
// "darwin" (macOS), "linux", "windows"
```

### process.stdout

Stream de saída padrão.

```javascript
process.stdout.write("Hello ");
process.stdout.write("World\n");
// Output: Hello World
```

#### process.stdout.write(data)

Escreve dados no stdout sem quebra de linha.

```javascript
process.stdout.write("Carregando");
for (let i = 0; i < 3; i++) {
  sleep(500);
  process.stdout.write(".");
}
process.stdout.write(" Pronto!\n");
// Output: Carregando... Pronto!
```

### process.stderr

Stream de erro padrão.

```javascript
process.stderr.write("ERRO: Algo deu errado\n");
```

## Métodos

### process.cwd()

Retorna o diretório de trabalho atual.

```javascript
console.log(process.cwd());
// "/Users/admin/Desktop/projetos/nulang"
```

### process.chdir(directory)

Muda o diretório de trabalho atual.

```javascript
console.log(process.cwd()); // "/Users/admin"
process.chdir("/tmp");
console.log(process.cwd()); // "/tmp"
```

### process.exit(code)

Encerra o processo com um código de saída.

```javascript
if (erro) {
  console.error("Erro fatal!");
  process.exit(1); // Código != 0 indica erro
}

// Saída bem-sucedida
process.exit(0);
```

| Código | Significado     |
| ------ | --------------- |
| `0`    | Sucesso         |
| `1`    | Erro genérico   |
| `> 0`  | Erro específico |

### process.nextTick(callback, ...args)

Agenda uma função para a próxima iteração do event loop.

```javascript
console.log("1");

process.nextTick(() => {
  console.log("3 - nextTick");
});

console.log("2");

// Output:
// 1
// 2
// 3 - nextTick
```

## Exemplos Práticos

### Parser de Argumentos Simples

```javascript
function parseArgs(args) {
  const result = { flags: {}, positional: [] };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];

    if (arg.startsWith("--")) {
      const key = arg.slice(2);
      // Verificar se próximo arg é um valor
      if (i + 1 < args.length && !args[i + 1].startsWith("-")) {
        result.flags[key] = args[i + 1];
        i++;
      } else {
        result.flags[key] = true;
      }
    } else if (arg.startsWith("-")) {
      const key = arg.slice(1);
      result.flags[key] = true;
    } else {
      result.positional.push(arg);
    }
  }

  return result;
}

// Uso
const args = parseArgs(process.argv.slice(2));
console.log(args);
// nulang script.nu --verbose -d --output file.txt input.txt
// { flags: { verbose: true, d: true, output: "file.txt" }, positional: ["input.txt"] }
```

### Configuração por Ambiente

```javascript
const config = {
  development: {
    debug: true,
    apiUrl: "http://localhost:3000",
  },
  production: {
    debug: false,
    apiUrl: "https://api.example.com",
  },
};

const env = process.env.NODE_ENV || "development";
const currentConfig = config[env];

console.log(`Ambiente: ${env}`);
console.log(`API URL: ${currentConfig.apiUrl}`);
```

### Progress Bar

```javascript
function progressBar(current, total, width = 40) {
  const percent = current / total;
  const filled = Math.floor(width * percent);
  const empty = width - filled;

  const bar = "█".repeat(filled) + "░".repeat(empty);
  const percentText = (percent * 100).toFixed(1);

  process.stdout.write(`\r[${bar}] ${percentText}%`);

  if (current === total) {
    process.stdout.write("\n");
  }
}

// Uso
for (let i = 0; i <= 100; i++) {
  progressBar(i, 100);
  sleep(50);
}
```

### Verificar Variáveis Obrigatórias

```javascript
function requireEnv(name) {
  const value = process.env[name];
  if (!value) {
    console.error(`Erro: Variável de ambiente ${name} não definida`);
    process.exit(1);
  }
  return value;
}

const apiKey = requireEnv("API_KEY");
const dbUrl = requireEnv("DATABASE_URL");
```

### CLI com Help

```javascript
const args = process.argv.slice(2);

if (args.includes("--help") || args.includes("-h")) {
  console.log(`
Uso: nulang meu-script.nu [opções] [arquivo]

Opções:
  --help, -h     Mostra esta ajuda
  --version, -v  Mostra a versão
  --verbose      Modo verboso
  --output FILE  Arquivo de saída
  `);
  process.exit(0);
}

if (args.includes("--version") || args.includes("-v")) {
  console.log("v1.0.0");
  process.exit(0);
}
```

### Tratamento de Erros Global

```javascript
function main() {
  // Código principal
  const args = process.argv.slice(2);

  if (args.length === 0) {
    throw new Error("Nenhum argumento fornecido");
  }

  console.log("Processando:", args.join(", "));
}

try {
  main();
} catch (error) {
  console.error("Erro:", error.message);
  process.exit(1);
}
```

### Spinner de Loading

```javascript
function spinner(duration) {
  const frames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
  let i = 0;
  const start = Date.now();

  while (Date.now() - start < duration) {
    process.stdout.write(`\r${frames[i]} Carregando...`);
    i = (i + 1) % frames.length;
    sleep(100);
  }

  process.stdout.write("\r✓ Concluído!    \n");
}

spinner(3000);
```

### Múltiplos Comandos

```javascript
const command = process.argv[2];
const args = process.argv.slice(3);

const commands = {
  build: () => {
    console.log("Construindo projeto...");
  },

  test: () => {
    console.log("Executando testes...");
  },

  deploy: () => {
    console.log("Fazendo deploy...");
  },

  help: () => {
    console.log("Comandos disponíveis: build, test, deploy");
  },
};

if (commands[command]) {
  commands[command](...args);
} else {
  console.error(`Comando desconhecido: ${command}`);
  commands.help();
  process.exit(1);
}
```

## Diferenças do Node.js

| Feature                  | Node.js | Nulang |
| ------------------------ | ------- | ------ |
| `process.argv`           | ✅      | ✅     |
| `process.env`            | ✅      | ✅     |
| `process.cwd()`          | ✅      | ✅     |
| `process.chdir()`        | ✅      | ✅     |
| `process.exit()`         | ✅      | ✅     |
| `process.nextTick()`     | ✅      | ✅     |
| `process.stdout.write()` | ✅      | ✅     |
| `process.stderr.write()` | ✅      | ✅     |
| `process.stdin`          | ✅      | ❌     |
| `process.on('exit')`     | ✅      | ❌     |
| `process.memoryUsage()`  | ✅      | ❌     |

## Veja Também

- [OS](./os.md) - Informações do sistema
- [File System](./filesystem.md) - Operações com arquivos
- [Event Loop](./event_loop.md) - Modelo de execução
