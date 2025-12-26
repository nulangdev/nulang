# HTTP Module

O módulo `http` fornece funcionalidades para criar servidores HTTP e fazer requisições HTTP, compatível com a API Node.js.

## Importação

```javascript
import http from "http";
// ou
const http = require("http");
```

## Classes

### http.Server

Classe que representa um servidor HTTP. Estende `EventEmitter`.

#### Criando um Servidor

```javascript
const server = http.createServer((req, res) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello World!");
});
```

#### Métodos

| Método                             | Descrição                               |
| ---------------------------------- | --------------------------------------- |
| `listen(port, [host], [callback])` | Inicia o servidor na porta especificada |
| `close([callback])`                | Encerra o servidor                      |

#### Eventos

| Evento    | Descrição                                |
| --------- | ---------------------------------------- |
| `request` | Emitido quando uma requisição é recebida |
| `close`   | Emitido quando o servidor é fechado      |

### http.IncomingMessage

Representa uma requisição HTTP de entrada. Estende `stream.Readable`.

#### Propriedades

| Propriedade   | Tipo   | Descrição                                  |
| ------------- | ------ | ------------------------------------------ |
| `method`      | String | Método HTTP (GET, POST, etc.)              |
| `url`         | String | URL da requisição                          |
| `httpVersion` | String | Versão do protocolo HTTP                   |
| `headers`     | Object | Headers da requisição                      |
| `rawHeaders`  | Array  | Headers em formato array [key, value, ...] |

#### Eventos

| Evento | Descrição                                     |
| ------ | --------------------------------------------- |
| `data` | Emitido quando dados são recebidos            |
| `end`  | Emitido quando todos os dados foram recebidos |

### http.ServerResponse

Representa a resposta HTTP. Estende `stream.Writable`.

#### Propriedades

| Propriedade     | Tipo    | Descrição                           |
| --------------- | ------- | ----------------------------------- |
| `statusCode`    | Number  | Código de status HTTP (padrão: 200) |
| `statusMessage` | String  | Mensagem de status                  |
| `headersSent`   | Boolean | Se os headers já foram enviados     |

#### Métodos

| Método                             | Descrição            |
| ---------------------------------- | -------------------- |
| `setHeader(name, value)`           | Define um header     |
| `getHeader(name)`                  | Obtém um header      |
| `writeHead(statusCode, [headers])` | Define o status code |
| `write(data)`                      | Escreve dados        |
| `end([data])`                      | Finaliza a resposta  |

### http.ClientRequest

Representa uma requisição HTTP de saída. Estende `stream.Writable`.

#### Métodos

| Método        | Descrição                            |
| ------------- | ------------------------------------ |
| `write(data)` | Escreve dados no corpo da requisição |
| `end([data])` | Finaliza a requisição                |

#### Eventos

| Evento     | Descrição                            |
| ---------- | ------------------------------------ |
| `response` | Emitido quando a resposta é recebida |
| `error`    | Emitido em caso de erro              |

### http.Agent

Gerencia conexões HTTP reutilizáveis.

## Funções

### http.createServer([requestListener])

Cria um novo servidor HTTP.

```javascript
const server = http.createServer((req, res) => {
  res.end("OK");
});

server.listen(3000, () => {
  console.log("Server running on port 3000");
});
```

### http.request(options, [callback])

Faz uma requisição HTTP.

```javascript
const req = http.request(
  {
    method: "POST",
    url: "http://example.com/api",
  },
  (res) => {
    res.on("data", (chunk) => {
      console.log(chunk.toString());
    });
  }
);

req.write("Hello");
req.end();
```

### http.get(url, [callback])

Faz uma requisição GET.

```javascript
http.get("http://example.com", (res) => {
  res.on("data", (chunk) => {
    console.log(chunk.toString());
  });
});
```

## Propriedades

### http.METHODS

Array com todos os métodos HTTP suportados.

```javascript
console.log(http.METHODS);
// ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"]
```

### http.STATUS_CODES

Objeto com todos os códigos de status HTTP e suas descrições.

```javascript
console.log(http.STATUS_CODES["200"]);
// "OK"
```

### http.globalAgent

Agent global usado para requisições.

## Exemplo Completo

```javascript
import http from "http";

// Criar servidor
const server = http.createServer((req, res) => {
  if (req.method === "GET" && req.url === "/") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ message: "Hello World!" }));
  } else {
    res.writeHead(404);
    res.end("Not Found");
  }
});

// Iniciar servidor
server.listen(3000, () => {
  console.log("Server running at http://localhost:3000/");
});

// Fazer requisição
http.get("http://localhost:3000", (res) => {
  let data = "";
  res.on("data", (chunk) => {
    data += chunk.toString();
  });
  res.on("end", () => {
    console.log("Response:", data);
  });
});
```

## Funções Auxiliares

### fetch(url, [options])

Função global para fazer requisições HTTP de forma simplificada (retorna Promise).

```javascript
const response = await fetch("https://api.example.com/data");
const data = await response.json();
console.log(data);
```

#### Opções

| Opção    | Tipo   | Descrição           |
| -------- | ------ | ------------------- |
| `method` | String | Método HTTP         |
| `body`   | String | Corpo da requisição |

## Módulos Relacionados

### url

Módulo para parsing de URLs.

```javascript
const url = require("url");
const parsed = url.parse("http://example.com:8080/path?query=value#hash");
console.log(parsed.hostname); // "example.com"
console.log(parsed.port); // "8080"
console.log(parsed.pathname); // "/path"
```

### querystring

Módulo para parsing de query strings.

```javascript
const qs = require("querystring");

// Parse
const obj = qs.parse("foo=bar&baz=qux");
console.log(obj); // { foo: "bar", baz: "qux" }

// Stringify
const str = qs.stringify({ name: "John", age: 30 });
console.log(str); // "name=John&age=30"
```
