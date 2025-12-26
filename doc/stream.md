# Stream Module

O módulo `stream` fornece classes base para trabalhar com streams de dados, compatível com a API Node.js.

## Importação

```javascript
import stream from "stream";
// ou
const stream = require("stream");
```

## Classes

### stream.Readable

Stream de leitura. Estende `EventEmitter`.

#### Propriedades

| Propriedade | Tipo    | Descrição                |
| ----------- | ------- | ------------------------ |
| `readable`  | Boolean | Se a stream está legível |

#### Métodos

| Método              | Descrição                                 |
| ------------------- | ----------------------------------------- |
| `read([size])`      | Lê dados da stream                        |
| `push(chunk)`       | Adiciona dados à stream (usar em `_read`) |
| `pipe(destination)` | Conecta a uma stream de escrita           |
| `destroy([error])`  | Destrói a stream                          |

#### Eventos

| Evento  | Descrição                              |
| ------- | -------------------------------------- |
| `data`  | Emitido quando dados estão disponíveis |
| `end`   | Emitido quando não há mais dados       |
| `close` | Emitido quando a stream é fechada      |

#### Implementação Customizada

```javascript
class MinhaStream extends stream.Readable {
  constructor() {
    super();
    this.data = ["item1", "item2", "item3"];
    this.index = 0;
  }

  _read() {
    if (this.index < this.data.length) {
      this.push(this.data[this.index]);
      this.index++;
    } else {
      this.push(null); // Sinaliza fim da stream
    }
  }
}

const s = new MinhaStream();
s.on("data", (chunk) => console.log(chunk.toString()));
s.on("end", () => console.log("Fim!"));
```

### stream.Writable

Stream de escrita. Estende `EventEmitter`.

#### Propriedades

| Propriedade | Tipo    | Descrição                            |
| ----------- | ------- | ------------------------------------ |
| `writable`  | Boolean | Se a stream está pronta para escrita |

#### Métodos

| Método                                 | Descrição               |
| -------------------------------------- | ----------------------- |
| `write(chunk, [encoding], [callback])` | Escreve dados na stream |
| `end([chunk])`                         | Finaliza a stream       |
| `cork()`                               | Pausa o envio de dados  |
| `uncork()`                             | Retoma o envio de dados |
| `destroy([error])`                     | Destrói a stream        |

#### Eventos

| Evento   | Descrição                                           |
| -------- | --------------------------------------------------- |
| `drain`  | Emitido quando a stream está pronta para mais dados |
| `finish` | Emitido quando `end()` é chamado                    |
| `close`  | Emitido quando a stream é fechada                   |

#### Implementação Customizada

```javascript
class MinhaWritable extends stream.Writable {
  _write(chunk, encoding, callback) {
    console.log("Recebido:", chunk.toString());
    callback(); // Sinaliza que está pronto para mais dados
  }
}

const w = new MinhaWritable();
w.write("Hello");
w.write("World");
w.end();
```

### stream.Transform

Stream de transformação (leitura e escrita). Estende `Readable` com métodos de `Writable`.

#### Métodos

| Método                                      | Descrição           |
| ------------------------------------------- | ------------------- |
| `_transform(chunk, encoding, callback)`     | Transforma os dados |
| Todos os métodos de `Readable` e `Writable` | -                   |

#### Implementação Customizada

```javascript
class UpperCase extends stream.Transform {
  _transform(chunk, encoding, callback) {
    this.push(chunk.toString().toUpperCase());
    callback();
  }
}

const upper = new UpperCase();
upper.on("data", (chunk) => console.log(chunk.toString()));
upper.write("hello"); // HELLO
upper.write("world"); // WORLD
```

### stream.PassThrough

Stream que simplesmente passa os dados sem modificação. Útil para testes e pipelines.

```javascript
const pass = new stream.PassThrough();

pass.on("data", (chunk) => {
  console.log("Passando:", chunk.toString());
});

pass.write("Dados");
pass.end();
```

## Funções Utilitárias

### stream.pipeline(source, ...transforms, destination)

Conecta múltiplas streams em sequência.

```javascript
const source = new stream.Readable();
const transform = new stream.Transform();
const dest = new stream.Writable();

const result = stream.pipeline(source, transform, dest);
```

## Padrões de Uso

### Piping

```javascript
const readable = new stream.Readable();
const writable = new stream.Writable();

// Conecta readable -> writable
readable.pipe(writable);
```

### Consumindo uma Readable Stream

```javascript
// Modo Flowing (eventos)
readable.on("data", (chunk) => {
  console.log("Chunk:", chunk);
});

readable.on("end", () => {
  console.log("Fim da stream");
});

// Modo Paused (pull)
let chunk;
while ((chunk = readable.read()) !== null) {
  console.log("Chunk:", chunk);
}
```

### Escrevendo em uma Writable Stream

```javascript
writable.write("Primeiro chunk");
writable.write("Segundo chunk");
writable.end("Último chunk"); // Finaliza a stream
```

## Exemplo Completo

```javascript
import stream from "stream";

// Criar stream customizada que gera números
class NumberStream extends stream.Readable {
  constructor(max) {
    super();
    this.max = max;
    this.current = 0;
  }

  _read() {
    if (this.current < this.max) {
      this.push(String(this.current));
      this.current++;
    } else {
      this.push(null);
    }
  }
}

// Criar transformador que dobra números
class DoubleTransform extends stream.Transform {
  _transform(chunk, encoding, callback) {
    const num = parseInt(chunk.toString());
    this.push(String(num * 2));
    callback();
  }
}

// Criar consumidor
class LogWriter extends stream.Writable {
  _write(chunk, encoding, callback) {
    console.log("Resultado:", chunk.toString());
    callback();
  }
}

// Conectar tudo
const numbers = new NumberStream(5);
const doubler = new DoubleTransform();
const logger = new LogWriter();

numbers.pipe(doubler).pipe(logger);
// Output:
// Resultado: 0
// Resultado: 2
// Resultado: 4
// Resultado: 6
// Resultado: 8
```

## Integração com HTTP

```javascript
import http from "http";

http
  .createServer((req, res) => {
    // req é uma Readable stream
    // res é uma Writable stream

    let body = "";
    req.on("data", (chunk) => {
      body += chunk.toString();
    });

    req.on("end", () => {
      res.write("Você enviou: " + body);
      res.end();
    });
  })
  .listen(3000);
```

## Buffer Integration

Streams trabalham nativamente com `Buffer`:

```javascript
const readable = new stream.Readable();

readable.push(Buffer.from("Hello "));
readable.push(Buffer.from("World"));
readable.push(null);

readable.on("data", (chunk) => {
  console.log(chunk.toString()); // "Hello World"
});
```
