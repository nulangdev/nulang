# Events Module

O módulo `events` fornece a classe `EventEmitter` para implementar o padrão de eventos, compatível com Node.js.

## Importação

```javascript
import { EventEmitter } from "events";
// ou
const events = require("events");
const EventEmitter = events.EventEmitter;
```

## Classe EventEmitter

### Criando um EventEmitter

```javascript
const emitter = new EventEmitter();

emitter.on("evento", (data) => {
  console.log("Recebido:", data);
});

emitter.emit("evento", "Hello!");
// Output: Recebido: Hello!
```

## Métodos

### on(eventName, listener)

Adiciona um listener para o evento.

```javascript
emitter.on("message", (msg) => {
  console.log("Mensagem:", msg);
});
```

**Alias**: `addListener`

### once(eventName, listener)

Adiciona um listener que é removido após a primeira execução.

```javascript
emitter.once("initialize", () => {
  console.log("Inicializado (apenas uma vez)");
});

emitter.emit("initialize"); // Dispara
emitter.emit("initialize"); // Não dispara
```

### emit(eventName, ...args)

Dispara um evento com os argumentos fornecidos.

```javascript
emitter.on("sum", (a, b) => {
  console.log("Soma:", a + b);
});

emitter.emit("sum", 5, 3); // Soma: 8
```

**Retorno**: Boolean (true se havia listeners)

### off(eventName, listener)

Remove um listener específico.

```javascript
const handler = (msg) => console.log(msg);

emitter.on("test", handler);
emitter.off("test", handler);

emitter.emit("test", "Hello"); // Não dispara
```

**Alias**: `removeListener`

### removeAllListeners([eventName])

Remove todos os listeners (opcionalmente de um evento específico).

```javascript
// Remove todos os listeners de um evento
emitter.removeAllListeners("message");

// Remove todos os listeners
emitter.removeAllListeners();
```

### prependListener(eventName, listener)

Adiciona um listener no início da fila.

```javascript
emitter.on("test", () => console.log("Segundo"));
emitter.prependListener("test", () => console.log("Primeiro"));

emitter.emit("test");
// Output:
// Primeiro
// Segundo
```

### prependOnceListener(eventName, listener)

Combina `prependListener` e `once`.

```javascript
emitter.prependOnceListener("init", () => {
  console.log("Primeiro e apenas uma vez");
});
```

### listeners(eventName)

Retorna uma cópia do array de listeners.

```javascript
const handlers = emitter.listeners("test");
console.log(handlers.length);
```

### rawListeners(eventName)

Retorna os listeners incluindo wrappers (de `once`).

```javascript
emitter.once("test", () => {});
const raw = emitter.rawListeners("test");
```

### listenerCount(eventName)

Retorna o número de listeners para um evento.

```javascript
emitter.on("test", () => {});
emitter.on("test", () => {});
console.log(emitter.listenerCount("test")); // 2
```

### eventNames()

Retorna um array com os nomes dos eventos que têm listeners.

```javascript
emitter.on("foo", () => {});
emitter.on("bar", () => {});
console.log(emitter.eventNames()); // ["foo", "bar"]
```

### setMaxListeners(n)

Define o limite máximo de listeners por evento.

```javascript
emitter.setMaxListeners(20);
```

### getMaxListeners()

Retorna o limite máximo de listeners.

```javascript
console.log(emitter.getMaxListeners()); // 10 (padrão)
```

## Eventos Especiais

### newListener

Emitido quando um novo listener é adicionado.

```javascript
emitter.on("newListener", (event, listener) => {
  console.log(`Novo listener para: ${event}`);
});

emitter.on("test", () => {});
// Output: Novo listener para: test
```

### removeListener

Emitido quando um listener é removido.

```javascript
emitter.on("removeListener", (event, listener) => {
  console.log(`Listener removido de: ${event}`);
});
```

### error

Se um evento `error` é emitido sem listeners, um erro é lançado.

```javascript
// Sempre adicione um handler de erro
emitter.on("error", (err) => {
  console.error("Erro:", err.message);
});

emitter.emit("error", new Error("Algo deu errado"));
// Output: Erro: Algo deu errado
```

## Propriedade Estática

### EventEmitter.defaultMaxListeners

Define o limite padrão de listeners para todas as instâncias.

```javascript
EventEmitter.defaultMaxListeners = 15;
```

## Herança

Você pode criar classes que estendem EventEmitter.

```javascript
class MyEmitter extends EventEmitter {
  constructor() {
    super();
    this.name = "MyEmitter";
  }

  doSomething() {
    this.emit("action", { timestamp: Date.now() });
  }
}

const my = new MyEmitter();
my.on("action", (data) => {
  console.log("Ação em:", data.timestamp);
});

my.doSomething();
```

## Exemplos Práticos

### Pub/Sub Simples

```javascript
const bus = new EventEmitter();

// Subscriber 1
bus.on("user:created", (user) => {
  console.log(`Enviando email para: ${user.email}`);
});

// Subscriber 2
bus.on("user:created", (user) => {
  console.log(`Criando perfil para: ${user.name}`);
});

// Publisher
function createUser(name, email) {
  const user = { name, email };
  bus.emit("user:created", user);
  return user;
}

createUser("João", "joao@email.com");
// Output:
// Enviando email para: joao@email.com
// Criando perfil para: João
```

### Timeout com Eventos

```javascript
class Timer extends EventEmitter {
  constructor(duration) {
    super();
    this.duration = duration;
  }

  start() {
    this.emit("start");
    setTimeout(() => {
      this.emit("end");
    }, this.duration);
  }
}

const timer = new Timer(2000);
timer.on("start", () => console.log("Timer iniciado"));
timer.on("end", () => console.log("Timer finalizado"));
timer.start();
```

### Cadeia de Eventos

```javascript
const pipeline = new EventEmitter();

pipeline.on("step1", (data) => {
  console.log("Passo 1:", data);
  pipeline.emit("step2", data + " -> step1");
});

pipeline.on("step2", (data) => {
  console.log("Passo 2:", data);
  pipeline.emit("step3", data + " -> step2");
});

pipeline.on("step3", (data) => {
  console.log("Passo 3:", data);
  pipeline.emit("done", data + " -> step3");
});

pipeline.on("done", (result) => {
  console.log("Resultado final:", result);
});

pipeline.emit("step1", "início");
// Output:
// Passo 1: início
// Passo 2: início -> step1
// Passo 3: início -> step1 -> step2
// Resultado final: início -> step1 -> step2 -> step3
```

### Contador de Eventos

```javascript
class Counter extends EventEmitter {
  constructor() {
    super();
    this.count = 0;
  }

  increment() {
    this.count++;
    this.emit("change", this.count);
    if (this.count % 10 === 0) {
      this.emit("milestone", this.count);
    }
  }
}

const counter = new Counter();
counter.on("change", (n) => console.log("Count:", n));
counter.on("milestone", (n) => console.log("🎉 Milestone:", n));

for (let i = 0; i < 15; i++) {
  counter.increment();
}
```

## Integração com Streams

Todas as classes de Stream herdam de EventEmitter.

```javascript
import stream from "stream";

const readable = new stream.Readable();
readable.on("data", (chunk) => console.log("Dados:", chunk));
readable.on("end", () => console.log("Fim"));
```

## Veja Também

- [Stream](./stream.md) - Streams de dados
- [HTTP](./http.md) - Servidor HTTP (usa eventos)
- [Process](./process.md) - Objeto process (usa eventos)
