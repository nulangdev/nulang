# Módulo Child Process

O módulo `child_process` permite criar e gerenciar processos filhos, executando comandos do sistema operacional e outros programas.

## Importação

```javascript
const child_process = require("child_process");
// ou
const { exec, execSync, spawn, fork } = require("child_process");
```

---

## Funções Principais

### `exec(command, options?, callback?)`

Executa um comando em um shell e armazena a saída em buffer.

```javascript
const { exec } = require("child_process");

exec("ls -la", (error, stdout, stderr) => {
  if (error) {
    console.error(`Erro: ${error.message}`);
    return;
  }
  if (stderr) {
    console.error(`stderr: ${stderr}`);
    return;
  }
  console.log(`stdout: ${stdout}`);
});

// Com opções
exec('echo "Hello"', { encoding: "utf8" }, (error, stdout, stderr) => {
  console.log(stdout);
});
```

#### Opções

| Opção       | Tipo           | Descrição                               |
| ----------- | -------------- | --------------------------------------- |
| `cwd`       | String         | Diretório de trabalho do processo filho |
| `env`       | Object         | Variáveis de ambiente                   |
| `encoding`  | String         | Codificação (padrão: 'utf8')            |
| `timeout`   | Number         | Timeout em milissegundos                |
| `maxBuffer` | Number         | Tamanho máximo do buffer de saída       |
| `shell`     | String/Boolean | Shell a usar                            |

```javascript
exec(
  "pwd",
  {
    cwd: "/tmp",
    env: { ...process.env, MY_VAR: "value" },
    timeout: 5000,
  },
  (error, stdout) => {
    console.log(stdout);
  }
);
```

---

### `execSync(command, options?)`

Versão síncrona de `exec`. Bloqueia até o comando terminar.

```javascript
const { execSync } = require("child_process");

try {
  const output = execSync("ls -la");
  console.log(output.toString());
} catch (error) {
  console.error("Comando falhou:", error.message);
}

// Com opções
const result = execSync('echo "Hello World"', {
  encoding: "utf8",
  cwd: "/tmp",
});
console.log(result);

// Capturando stderr
try {
  execSync("comando-invalido");
} catch (error) {
  console.log("stderr:", error.stderr.toString());
  console.log("Código de saída:", error.status);
}
```

---

### `spawn(command, args?, options?)`

Inicia um novo processo com um comando. Mais flexível que `exec`, usa streams.

```javascript
const { spawn } = require("child_process");

// Executar comando simples
const ls = spawn("ls", ["-la", "/tmp"]);

ls.stdout.on("data", (data) => {
  console.log(`stdout: ${data}`);
});

ls.stderr.on("data", (data) => {
  console.error(`stderr: ${data}`);
});

ls.on("close", (code) => {
  console.log(`Processo terminou com código ${code}`);
});

ls.on("error", (error) => {
  console.error(`Erro ao iniciar processo: ${error.message}`);
});
```

#### Opções do spawn

| Opção      | Tipo           | Descrição                      |
| ---------- | -------------- | ------------------------------ |
| `cwd`      | String         | Diretório de trabalho          |
| `env`      | Object         | Variáveis de ambiente          |
| `stdio`    | Array/String   | Configuração de stdio          |
| `shell`    | Boolean/String | Executar em shell              |
| `detached` | Boolean        | Executar processo independente |

```javascript
// Executar em shell
const child = spawn('echo "Hello $USER"', {
  shell: true,
});

// Herdar stdio do pai
const child = spawn("npm", ["install"], {
  stdio: "inherit",
  cwd: "./myproject",
});

// Processo detached
const child = spawn("node", ["server.js"], {
  detached: true,
  stdio: "ignore",
});
child.unref(); // Permite que o pai termine
```

---

### `spawnSync(command, args?, options?)`

Versão síncrona de `spawn`.

```javascript
const { spawnSync } = require("child_process");

const result = spawnSync("ls", ["-la"]);

console.log("stdout:", result.stdout.toString());
console.log("stderr:", result.stderr.toString());
console.log("status:", result.status);
console.log("signal:", result.signal);

// Com timeout
const result2 = spawnSync("sleep", ["10"], {
  timeout: 1000,
});
if (result2.error) {
  console.log("Timeout ou erro:", result2.error.message);
}
```

---

### `execFile(file, args?, options?, callback?)`

Similar a `exec`, mas executa um arquivo diretamente sem shell.

```javascript
const { execFile } = require("child_process");

execFile("/usr/bin/node", ["--version"], (error, stdout, stderr) => {
  console.log("Node version:", stdout);
});

// Executar script
execFile("./myscript.sh", ["arg1", "arg2"], (error, stdout, stderr) => {
  console.log(stdout);
});
```

### `execFileSync(file, args?, options?)`

Versão síncrona de `execFile`.

```javascript
const { execFileSync } = require("child_process");

const version = execFileSync("/usr/bin/node", ["--version"], {
  encoding: "utf8",
});
console.log("Node version:", version.trim());
```

---

### `fork(modulePath, args?, options?)`

Cria um novo processo Node.js e estabelece comunicação via IPC.

```javascript
const { fork } = require("child_process");

// child.js
// process.on('message', (msg) => {
//   console.log('Filho recebeu:', msg);
//   process.send({ response: 'Olá pai!' });
// });

const child = fork("./child.js");

child.send({ hello: "world" });

child.on("message", (msg) => {
  console.log("Pai recebeu:", msg);
});

child.on("close", (code) => {
  console.log("Filho terminou com código:", code);
});
```

---

## Objeto ChildProcess

O objeto retornado por `spawn`, `exec`, `execFile` e `fork`.

### Propriedades

| Propriedade | Tipo    | Descrição               |
| ----------- | ------- | ----------------------- |
| `pid`       | Number  | ID do processo          |
| `stdin`     | Stream  | Stream de entrada       |
| `stdout`    | Stream  | Stream de saída         |
| `stderr`    | Stream  | Stream de erro          |
| `killed`    | Boolean | Se o processo foi morto |
| `exitCode`  | Number  | Código de saída         |

```javascript
const child = spawn("ls");

console.log("PID:", child.pid);

child.on("exit", () => {
  console.log("Exit code:", child.exitCode);
});
```

### Métodos

#### `child.kill(signal?)`

Envia um sinal para o processo filho.

```javascript
const child = spawn("sleep", ["100"]);

setTimeout(() => {
  child.kill("SIGTERM"); // ou 'SIGKILL'
  console.log("Processo morto");
}, 1000);
```

#### `child.send(message, callback?)`

Envia uma mensagem para processos criados com `fork`.

```javascript
const child = fork("./worker.js");

child.send({ task: "processData", data: [1, 2, 3] });
```

#### `child.disconnect()`

Fecha o canal IPC.

```javascript
child.disconnect();
```

### Eventos

#### `'exit'`

Emitido quando o processo termina.

```javascript
child.on("exit", (code, signal) => {
  console.log(`Processo saiu com código ${code} e sinal ${signal}`);
});
```

#### `'close'`

Emitido quando as streams stdio são fechadas.

```javascript
child.on("close", (code, signal) => {
  console.log("Streams fechadas");
});
```

#### `'error'`

Emitido quando ocorre um erro.

```javascript
child.on("error", (error) => {
  console.error("Erro:", error.message);
});
```

#### `'message'`

Emitido quando uma mensagem é recebida via IPC (fork).

```javascript
child.on("message", (msg) => {
  console.log("Mensagem:", msg);
});
```

---

## Exemplos Práticos

### Executar Script Python

```javascript
const { spawn } = require("child_process");

const python = spawn("python3", ["script.py", "arg1", "arg2"]);

python.stdout.on("data", (data) => {
  console.log(`Python output: ${data}`);
});

python.on("close", (code) => {
  console.log(`Python script exited with code ${code}`);
});
```

### Pipeline de Comandos

```javascript
const { spawn } = require("child_process");

// cat file.txt | grep "pattern" | wc -l

const cat = spawn("cat", ["file.txt"]);
const grep = spawn("grep", ["pattern"]);
const wc = spawn("wc", ["-l"]);

cat.stdout.pipe(grep.stdin);
grep.stdout.pipe(wc.stdin);

wc.stdout.on("data", (data) => {
  console.log(`Linhas encontradas: ${data.toString().trim()}`);
});
```

### Worker Pool com Fork

```javascript
// main.js
const { fork } = require("child_process");

const workers = [];
const numWorkers = 4;

for (let i = 0; i < numWorkers; i++) {
  const worker = fork("./worker.js");
  workers.push(worker);

  worker.on("message", (result) => {
    console.log(`Worker ${i} resultado:`, result);
  });
}

// Distribuir tarefas
let taskIndex = 0;
const tasks = [1, 2, 3, 4, 5, 6, 7, 8];

tasks.forEach((task) => {
  const worker = workers[taskIndex % numWorkers];
  worker.send({ task });
  taskIndex++;
});

// worker.js
process.on("message", (msg) => {
  const result = msg.task * 2;
  process.send({ result });
});
```

### Executar Git

```javascript
const { execSync } = require("child_process");

function git(command) {
  try {
    return execSync(`git ${command}`, { encoding: "utf8" }).trim();
  } catch (error) {
    console.error("Git error:", error.message);
    return null;
  }
}

console.log("Branch atual:", git("branch --show-current"));
console.log("Status:", git("status --short"));
console.log("Último commit:", git("log -1 --oneline"));
```

---

## Notas de Compatibilidade

| Funcionalidade    | Nulang    | Node.js |
| ----------------- | --------- | ------- |
| `exec`            | ✅        | ✅      |
| `execSync`        | ✅        | ✅      |
| `execFile`        | ✅        | ✅      |
| `execFileSync`    | ✅        | ✅      |
| `spawn`           | ✅        | ✅      |
| `spawnSync`       | ✅        | ✅      |
| `fork`            | ✅        | ✅      |
| `child.kill`      | ✅        | ✅      |
| `child.send`      | ✅        | ✅      |
| `child.pid`       | ✅        | ✅      |
| stdio: 'pipe'     | ✅        | ✅      |
| stdio: 'inherit'  | ✅        | ✅      |
| IPC channel       | ⚠️ Básico | ✅      |
| `child.ref/unref` | ❌        | ✅      |
