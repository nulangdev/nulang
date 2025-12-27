# Módulo Net

O módulo `net` fornece uma API para criação de servidores e clientes TCP, permitindo comunicação de rede baseada em streams.

## Importação

```javascript
const net = require("net");
```

---

## Criando um Servidor TCP

### `net.createServer(options?, connectionListener?)`

Cria um novo servidor TCP.

```javascript
const server = net.createServer((socket) => {
  console.log("Cliente conectado");

  socket.on("data", (data) => {
    console.log("Dados recebidos:", data.toString());
    socket.write("Mensagem recebida!");
  });

  socket.on("end", () => {
    console.log("Cliente desconectado");
  });
});

server.listen(3000, "127.0.0.1", () => {
  console.log("Servidor iniciado na porta 3000");
});
```

#### Opções

| Opção           | Tipo    | Descrição                                                    |
| --------------- | ------- | ------------------------------------------------------------ |
| `allowHalfOpen` | Boolean | Se true, mantém conexão aberta quando o outro lado envia FIN |

---

## Eventos do Servidor

### `'listening'`

Emitido quando o servidor começa a ouvir conexões.

```javascript
server.on("listening", () => {
  console.log("Servidor pronto para conexões");
});
```

### `'connection'`

Emitido quando uma nova conexão é estabelecida.

```javascript
server.on("connection", (socket) => {
  console.log("Nova conexão de:", socket.remoteAddress);
});
```

### `'close'`

Emitido quando o servidor é fechado.

```javascript
server.on("close", () => {
  console.log("Servidor encerrado");
});
```

### `'error'`

Emitido quando ocorre um erro.

```javascript
server.on("error", (err) => {
  console.error("Erro no servidor:", err.message);
});
```

---

## Métodos do Servidor

### `server.listen(port, host?, callback?)`

Inicia o servidor na porta e host especificados.

```javascript
server.listen(8080, "0.0.0.0", () => {
  console.log("Servidor ouvindo em 0.0.0.0:8080");
});
```

### `server.close(callback?)`

Fecha o servidor, parando de aceitar novas conexões.

```javascript
server.close(() => {
  console.log("Servidor fechado");
});
```

### `server.address()`

Retorna o endereço vinculado do servidor.

```javascript
const addr = server.address();
console.log(`Servidor em ${addr.address}:${addr.port}`);
```

---

## Criando um Cliente TCP

### `net.createConnection(options, connectListener?)`

Cria uma nova conexão TCP.

```javascript
const client = net.createConnection({ port: 3000, host: "127.0.0.1" }, () => {
  console.log("Conectado ao servidor");
  client.write("Olá, servidor!");
});

client.on("data", (data) => {
  console.log("Resposta do servidor:", data.toString());
  client.end();
});

client.on("end", () => {
  console.log("Desconectado do servidor");
});
```

### `net.connect(options, connectListener?)`

Alias para `createConnection`.

```javascript
const client = net.connect({ port: 3000 }, () => {
  console.log("Conectado!");
});
```

---

## Objeto Socket

O objeto Socket representa uma conexão TCP (tanto no cliente quanto no servidor).

### Propriedades do Socket

| Propriedade     | Tipo   | Descrição                        |
| --------------- | ------ | -------------------------------- |
| `remoteAddress` | String | Endereço IP do cliente conectado |
| `remotePort`    | Number | Porta do cliente conectado       |
| `localAddress`  | String | Endereço IP local                |
| `localPort`     | Number | Porta local                      |

```javascript
socket.on("connect", () => {
  console.log("Conectado a:", socket.remoteAddress, socket.remotePort);
});
```

### Métodos do Socket

#### `socket.write(data, encoding?, callback?)`

Envia dados através do socket.

```javascript
socket.write("Hello World");
socket.write(Buffer.from([0x48, 0x65, 0x6c, 0x6c, 0x6f]));
```

#### `socket.end(data?, encoding?, callback?)`

Encerra o socket opcionalmente enviando dados finais.

```javascript
socket.end("Goodbye");
```

#### `socket.destroy(error?)`

Destrói o socket imediatamente.

```javascript
socket.destroy();
```

#### `socket.setEncoding(encoding)`

Define a codificação para dados recebidos.

```javascript
socket.setEncoding("utf8");
```

#### `socket.setTimeout(timeout, callback?)`

Define um timeout para inatividade.

```javascript
socket.setTimeout(30000, () => {
  console.log("Socket timeout");
  socket.end();
});
```

#### `socket.setNoDelay(noDelay?)`

Desabilita o algoritmo Nagle.

```javascript
socket.setNoDelay(true);
```

#### `socket.setKeepAlive(enable?, initialDelay?)`

Habilita keep-alive.

```javascript
socket.setKeepAlive(true, 60000);
```

### Eventos do Socket

#### `'connect'`

Emitido quando a conexão é estabelecida.

```javascript
socket.on("connect", () => {
  console.log("Conectado!");
});
```

#### `'data'`

Emitido quando dados são recebidos.

```javascript
socket.on("data", (data) => {
  console.log("Recebido:", data.toString());
});
```

#### `'end'`

Emitido quando o outro lado envia FIN.

```javascript
socket.on("end", () => {
  console.log("Conexão encerrada pelo servidor");
});
```

#### `'close'`

Emitido quando o socket é completamente fechado.

```javascript
socket.on("close", (hadError) => {
  console.log("Socket fechado", hadError ? "com erro" : "normalmente");
});
```

#### `'error'`

Emitido quando ocorre um erro.

```javascript
socket.on("error", (err) => {
  console.error("Erro no socket:", err.message);
});
```

#### `'timeout'`

Emitido quando o socket fica inativo após o timeout definido.

```javascript
socket.on("timeout", () => {
  console.log("Socket inativo");
});
```

---

## Funções Utilitárias

### `net.isIP(input)`

Verifica se uma string é um endereço IP válido.

```javascript
net.isIP("127.0.0.1"); // 4 (IPv4)
net.isIP("::1"); // 6 (IPv6)
net.isIP("invalid"); // 0 (não é IP)
```

### `net.isIPv4(input)`

Verifica se é um endereço IPv4.

```javascript
net.isIPv4("192.168.1.1"); // true
net.isIPv4("::1"); // false
```

### `net.isIPv6(input)`

Verifica se é um endereço IPv6.

```javascript
net.isIPv6("::1"); // true
net.isIPv6("192.168.1.1"); // false
```

---

## Exemplo Completo: Chat Server

```javascript
const net = require("net");

const clients = [];

const server = net.createServer((socket) => {
  console.log("Novo cliente conectado");
  clients.push(socket);

  socket.on("data", (data) => {
    const message = data.toString().trim();
    console.log("Mensagem recebida:", message);

    // Broadcast para todos os outros clientes
    clients.forEach((client) => {
      if (client !== socket) {
        client.write(`Cliente disse: ${message}\n`);
      }
    });
  });

  socket.on("end", () => {
    const index = clients.indexOf(socket);
    if (index > -1) {
      clients.splice(index, 1);
    }
    console.log("Cliente desconectado");
  });

  socket.on("error", (err) => {
    console.error("Erro:", err.message);
  });
});

server.listen(4000, () => {
  console.log("Chat server iniciado na porta 4000");
});
```

---

## Exemplo: Cliente Echo

```javascript
const net = require("net");

const client = net.createConnection({ port: 4000 }, () => {
  console.log("Conectado ao servidor");
  client.write("Olá!");
});

client.on("data", (data) => {
  console.log("Servidor respondeu:", data.toString());
});

client.on("end", () => {
  console.log("Desconectado");
});

client.on("error", (err) => {
  console.error("Erro:", err.message);
});
```

---

## Notas de Compatibilidade

| Funcionalidade     | Nulang        | Node.js |
| ------------------ | ------------- | ------- |
| `createServer`     | ✅            | ✅      |
| `createConnection` | ✅            | ✅      |
| `connect`          | ✅            | ✅      |
| Socket.write/end   | ✅            | ✅      |
| Socket.destroy     | ✅            | ✅      |
| isIP/isIPv4/isIPv6 | ✅            | ✅      |
| Socket.pipe        | ⚠️ Via Stream | ✅      |
| UnixSocket         | ❌            | ✅      |
| Socket.ref/unref   | ❌            | ✅      |
