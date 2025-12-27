# Event Loop

O Event Loop é o mecanismo que permite execução assíncrona em Nulang.

## Conceito

O Event Loop gerencia a execução de código, coleta e processa eventos, e executa sub-tarefas enfileiradas.

```
┌───────────────────────────┐
│       Call Stack          │
└───────────────────────────┘
              ↓
┌───────────────────────────┐
│       Event Loop          │ ← Verifica se há callbacks pendentes
└───────────────────────────┘
              ↓
┌───────────────────────────┐
│    Callback Queue         │ ← setTimeout, setInterval, etc.
└───────────────────────────┘
```

## Fases de Execução

1. **Código Síncrono**: Executa primeiro
2. **nextTick**: Executa após código síncrono atual
3. **Immediate**: setImmediate callbacks
4. **Timers**: setTimeout e setInterval callbacks
5. **I/O**: Callbacks de operações de I/O
6. **Poll**: Aguarda novos eventos

## Ordem de Execução

```javascript
console.log("1 - Síncrono início");

setTimeout(() => {
  console.log("5 - setTimeout 0ms");
}, 0);

setImmediate(() => {
  console.log("4 - setImmediate");
});

process.nextTick(() => {
  console.log("3 - nextTick");
});

console.log("2 - Síncrono fim");

// Output:
// 1 - Síncrono início
// 2 - Síncrono fim
// 3 - nextTick
// 4 - setImmediate
// 5 - setTimeout 0ms
```

## Tarefas Assíncronas

### Macro Tasks (Timers)

Executam em iterações separadas do event loop.

```javascript
setTimeout(() => console.log("Timeout"), 0);
setInterval(() => console.log("Interval"), 1000);
setImmediate(() => console.log("Immediate"));
```

### Micro Tasks (nextTick)

Executam imediatamente após o código síncrono atual.

```javascript
process.nextTick(() => console.log("Next tick"));
Promise.resolve().then(() => console.log("Promise"));
```

### Prioridade

```
Código Síncrono → nextTick → Promises → setImmediate → setTimeout
```

## Exemplos

### Event Loop Básico

```javascript
console.log("Início");

setTimeout(() => {
  console.log("Timeout 1");
}, 0);

setTimeout(() => {
  console.log("Timeout 2");
}, 0);

process.nextTick(() => {
  console.log("NextTick 1");
});

process.nextTick(() => {
  console.log("NextTick 2");
});

console.log("Fim");

// Output:
// Início
// Fim
// NextTick 1
// NextTick 2
// Timeout 1
// Timeout 2
```

### Promises e Event Loop

```javascript
console.log("1");

Promise.resolve().then(() => {
  console.log("3 - Promise");
});

console.log("2");

// Output:
// 1
// 2
// 3 - Promise
```

### Aninhamento

```javascript
setTimeout(() => {
  console.log("A");

  process.nextTick(() => {
    console.log("B - nextTick dentro de timeout");
  });

  console.log("C");
}, 0);

// Output:
// A
// C
// B - nextTick dentro de timeout
```

## Bloqueio do Event Loop

### ❌ Código Bloqueante

```javascript
// RUIM: Bloqueia o event loop
function fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

console.log("Calculando...");
const result = fibonacci(40); // Bloqueia por vários segundos
console.log("Resultado:", result);

// Timeouts não executam durante o bloqueio
setTimeout(() => console.log("Timer"), 100); // Atrasado!
```

### ✅ Evitar Bloqueio

```javascript
// BOM: Quebrar trabalho pesado
function processChunk(items, index, callback) {
  const chunkSize = 100;
  const end = Math.min(index + chunkSize, items.length);

  for (let i = index; i < end; i++) {
    // Processar item
  }

  if (end < items.length) {
    setImmediate(() => processChunk(items, end, callback));
  } else {
    callback();
  }
}
```

## Patterns Comuns

### Defer Execution

```javascript
function deferredInit() {
  setImmediate(() => {
    console.log("Inicialização diferida");
    // Código de inicialização
  });
}

deferredInit();
console.log("Continua imediatamente");
```

### Batch Processing

```javascript
function processBatch(items, batchSize, processFn, done) {
  let index = 0;

  function processNext() {
    const batch = items.slice(index, index + batchSize);

    if (batch.length === 0) {
      done();
      return;
    }

    batch.forEach(processFn);
    index += batchSize;

    // Cede controle para outras tarefas
    setImmediate(processNext);
  }

  processNext();
}

processBatch(
  [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
  3,
  (item) => console.log("Processando:", item),
  () => console.log("Concluído!")
);
```

### Cooperative Scheduling

```javascript
function cooperativeTask(iterations, callback) {
  let i = 0;

  function work() {
    const start = Date.now();

    while (i < iterations && Date.now() - start < 10) {
      // Trabalho
      i++;
    }

    if (i < iterations) {
      setImmediate(work); // Cede controle
    } else {
      callback();
    }
  }

  work();
}
```

## Async Tasks em Nulang

Nulang rastreia tarefas assíncronas internamente:

```javascript
// Tarefas assíncronas são registradas
setTimeout(() => {
  console.log("Isso executa");
}, 1000);

// O script não termina até todas as tarefas completarem
console.log("Script principal finalizado");

// Output:
// Script principal finalizado
// (pausa de 1 segundo)
// Isso executa
```

## Diferenças do Node.js

| Aspecto         | Node.js | Nulang              |
| --------------- | ------- | ------------------- |
| libuv           | Sim     | Não (goroutines Go) |
| Fases completas | 6 fases | Simplificado        |
| Worker Threads  | Sim     | Não                 |
| Child Processes | Sim     | Não                 |

## Boas Práticas

### ✅ Não bloqueie o Event Loop

```javascript
// Evite loops longos síncronos
// Use setImmediate para ceder controle
```

### ✅ Use nextTick para callbacks imediatos

```javascript
function asyncOperation(callback) {
  // Garante que callback seja sempre assíncrono
  process.nextTick(() => {
    callback(null, result);
  });
}
```

### ✅ Quebre trabalho pesado

```javascript
// Divida em chunks menores
// Use setImmediate entre chunks
```

### ❌ Evite aninhamento profundo

```javascript
// Ruim
setTimeout(() => {
  setTimeout(() => {
    setTimeout(() => {
      // ...
    }, 0);
  }, 0);
}, 0);

// Melhor: use sequência ou Promises
```

## Debugging

### Verificar Ordem de Execução

```javascript
function trace(label) {
  console.log(`[${Date.now()}] ${label}`);
}

trace("1 - sync");
setTimeout(() => trace("3 - timeout"), 0);
process.nextTick(() => trace("2 - nextTick"));
trace("1b - sync");
```

### Identificar Bloqueios

```javascript
const start = Date.now();

// Código potencialmente bloqueante
heavyOperation();

const duration = Date.now() - start;
if (duration > 100) {
  console.warn(`Operação bloqueou por ${duration}ms`);
}
```

## Veja Também

- [Timers](./timers.md) - setTimeout, setInterval
- [Promise](./promise.md) - Operações assíncronas
- [Process](./process.md) - process.nextTick
