# Módulo Readline

O módulo `readline` fornece uma interface para leitura de dados linha por linha.

## Importação

```javascript
const readline = require("readline");
```

## Criando uma Interface

```javascript
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  prompt: "> ",
});
```

## Métodos

### `rl.question(query, callback)`

```javascript
rl.question("Nome: ", (answer) => {
  console.log(`Olá, ${answer}!`);
  rl.close();
});
```

### `rl.prompt()`

```javascript
rl.prompt();
rl.on("line", (line) => {
  console.log(`Você digitou: ${line}`);
  rl.prompt();
});
```

### `rl.setPrompt(prompt)`

```javascript
rl.setPrompt(">>> ");
```

### `rl.write(data)`

```javascript
rl.write("Texto de saída");
```

### `rl.close()`

```javascript
rl.close();
```

### `rl.pause()` / `rl.resume()`

```javascript
rl.pause();
rl.resume();
```

## Eventos

### `'line'`

```javascript
rl.on("line", (input) => {
  console.log(`Recebido: ${input}`);
});
```

### `'close'`

```javascript
rl.on("close", () => {
  console.log("Até logo!");
});
```

## Funções de Tela

### `readline.clearLine(stream, dir)`

```javascript
readline.clearLine(process.stdout, 0);
```

### `readline.clearScreenDown(stream)`

```javascript
readline.clearScreenDown(process.stdout);
```

### `readline.cursorTo(stream, x, y?)`

```javascript
readline.cursorTo(process.stdout, 0, 5);
```

### `readline.moveCursor(stream, dx, dy)`

```javascript
readline.moveCursor(process.stdout, 3, 2);
```

## Exemplo: REPL

```javascript
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  prompt: "nulang> ",
});

rl.prompt();

rl.on("line", (line) => {
  if (line.trim() === "exit") {
    rl.close();
    return;
  }
  console.log(`Comando: ${line}`);
  rl.prompt();
});

rl.on("close", () => {
  console.log("Bye!");
});
```

## Compatibilidade

| Funcionalidade    | Nulang | Node.js |
| ----------------- | ------ | ------- |
| `createInterface` | ✅     | ✅      |
| `question`        | ✅     | ✅      |
| `prompt`          | ✅     | ✅      |
| `clearLine`       | ✅     | ✅      |
| `cursorTo`        | ✅     | ✅      |
| Histórico         | ❌     | ✅      |
