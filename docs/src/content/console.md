# Console

O objeto global `console` fornece uma interface para debugging e logging, compatível com Node.js.

## Métodos de Logging

### `console.log([data][, ...args])`

Imprime mensagens no stdout.

```javascript
console.log("Hello", "World"); // Hello World
console.log({ name: "John", age: 30 }); // {name: John, age: 30}
```

### `console.info([data][, ...args])`

Alias para `console.log()`.

```javascript
console.info("Informational message");
```

### `console.debug([data][, ...args])`

Alias para `console.log()`.

```javascript
console.debug("Debug message");
```

### `console.error([data][, ...args])`

Imprime mensagens no stderr.

```javascript
console.error("An error occurred!");
```

### `console.warn([data][, ...args])`

Imprime avisos no stderr.

```javascript
console.warn("This is a warning");
```

## Asserções

### `console.assert(value[, ...message])`

Verifica se `value` é verdadeiro. Se falso, imprime uma mensagem de erro.

```javascript
console.assert(true, "This won't print");
console.assert(false, "Assertion failed!"); // Assertion failed: Assertion failed!
console.assert(1 === 2, "Math check"); // Assertion failed: Math check
```

## Contadores

### `console.count([label])`

Conta quantas vezes foi chamado com o mesmo label.

```javascript
console.count("myCounter"); // myCounter: 1
console.count("myCounter"); // myCounter: 2
console.count("myCounter"); // myCounter: 3
console.count(); // default: 1
```

### `console.countReset([label])`

Reseta o contador para zero.

```javascript
console.count("myCounter"); // myCounter: 1
console.count("myCounter"); // myCounter: 2
console.countReset("myCounter");
console.count("myCounter"); // myCounter: 1
```

## Grupos

### `console.group([...label])`

Cria um novo grupo de indentação no console.

```javascript
console.log("Level 0");
console.group("Group 1");
console.log("Level 1");
console.group("Nested");
console.log("Level 2");
console.groupEnd();
console.log("Back to Level 1");
console.groupEnd();
console.log("Back to Level 0");
```

Saída:

```
Level 0
Group 1
  Level 1
  Nested
    Level 2
  Back to Level 1
Back to Level 0
```

### `console.groupCollapsed([...label])`

Igual a `console.group()` (em terminal, funciona da mesma forma).

### `console.groupEnd()`

Encerra o grupo atual de indentação.

## Tabelas

### `console.table(tabularData[, properties])`

Exibe dados tabulares em formato de tabela.

```javascript
// Array
console.table(["apple", "banana", "cherry"]);
```

Saída:

```
┌───────┬────────┐
│ (idx) │ Values │
├───────┼────────┤
│ 0     │ apple  │
│ 1     │ banana │
│ 2     │ cherry │
└───────┴────────┘
```

```javascript
// Object
console.table({ name: "John", age: 30, city: "NYC" });
```

Saída:

```
┌──────┬───────┐
│ Key  │ Value │
├──────┼───────┤
│ name │ John  │
│ age  │ 30    │
│ city │ NYC   │
└──────┴───────┘
```

## Timers

### `console.time([label])`

Inicia um timer com o label especificado.

```javascript
console.time("myTimer");
```

### `console.timeLog([label][, ...data])`

Registra o tempo atual do timer sem encerrá-lo.

```javascript
console.time("operation");
// ... algum trabalho
console.timeLog("operation", "checkpoint 1"); // operation: 1.234ms checkpoint 1
// ... mais trabalho
console.timeLog("operation", "checkpoint 2"); // operation: 2.567ms checkpoint 2
```

### `console.timeEnd([label])`

Encerra o timer e imprime o tempo decorrido.

```javascript
console.time("myTimer");
// ... algum trabalho
console.timeEnd("myTimer"); // myTimer: 123.456ms
```

## Inspeção de Objetos

### `console.dir(obj[, options])`

Exibe uma representação do objeto de forma inspecionável.

```javascript
console.dir({ a: 1, b: { c: 2 } });
// {a: 1, b: {c: 2}}
```

### `console.dirxml(...data)`

Igual a `console.dir()` em ambientes não-DOM.

```javascript
console.dirxml({ name: "test" });
```

## Stack Trace

### `console.trace([message][, ...args])`

Imprime um stack trace.

```javascript
function myFunction() {
  console.trace("Trace from myFunction");
}
myFunction();
// Trace: Trace from myFunction
//     at <anonymous>
```

## Limpeza

### `console.clear()`

Limpa a tela do console (usando códigos ANSI).

```javascript
console.clear();
```

## Exemplo Completo

```javascript
// Logging básico
console.log("Starting application...");
console.info("Version: 1.0.0");

// Timer para medir performance
console.time("init");

// Grupos para organizar output
console.group("Configuration");
console.log("Port:", 3000);
console.log("Debug:", true);
console.groupEnd();

// Contadores para tracking
for (let i = 0; i < 3; i++) {
  console.count("loop");
}

// Tabela para dados estruturados
const users = ["Alice", "Bob", "Charlie"];
console.table(users);

// Fim do timer
console.timeEnd("init");

// Asserção para validação
const config = { valid: true };
console.assert(config.valid, "Config should be valid");
```

## Compatibilidade

| Método           | Nulang | Node.js |
| ---------------- | ------ | ------- |
| `log`            | ✅     | ✅      |
| `info`           | ✅     | ✅      |
| `debug`          | ✅     | ✅      |
| `error`          | ✅     | ✅      |
| `warn`           | ✅     | ✅      |
| `assert`         | ✅     | ✅      |
| `clear`          | ✅     | ✅      |
| `count`          | ✅     | ✅      |
| `countReset`     | ✅     | ✅      |
| `dir`            | ✅     | ✅      |
| `dirxml`         | ✅     | ✅      |
| `group`          | ✅     | ✅      |
| `groupCollapsed` | ✅     | ✅      |
| `groupEnd`       | ✅     | ✅      |
| `table`          | ✅     | ✅      |
| `time`           | ✅     | ✅      |
| `timeEnd`        | ✅     | ✅      |
| `timeLog`        | ✅     | ✅      |
| `trace`          | ✅     | ✅      |
| `profile`        | ❌     | ✅      |
| `profileEnd`     | ❌     | ✅      |
