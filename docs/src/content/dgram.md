# Módulo Dgram (UDP)

O módulo `dgram` fornece uma implementação de sockets UDP (User Datagram Protocol), permitindo comunicação connectionless (sem conexão).

## Importação

```javascript
const dgram = require("dgram");
```

---

## Criando um Socket UDP

### `dgram.createSocket(type, callback?)`

Cria um novo socket UDP.

```javascript
// UDP IPv4
const socket = dgram.createSocket("udp4");

// UDP IPv6
const socket6 = dgram.createSocket("udp6");

// Com callback para mensagens
const socket = dgram.createSocket("udp4", (msg, rinfo) => {
  console.log(`Mensagem de ${rinfo.address}:${rinfo.port}`);
  console.log(msg.toString());
});
```

#### Tipos de Socket

| Tipo     | Descrição       |
| -------- | --------------- |
| `'udp4'` | Socket UDP IPv4 |
| `'udp6'` | Socket UDP IPv6 |

---

## Eventos do Socket

### `'message'`

Emitido quando uma mensagem é recebida.

```javascript
socket.on("message", (msg, rinfo) => {
  console.log(`Mensagem: ${msg.toString()}`);
  console.log(`De: ${rinfo.address}:${rinfo.port}`);
  console.log(`Família: ${rinfo.family}`);
});
```

#### Objeto `rinfo`

| Propriedade | Tipo   | Descrição                |
| ----------- | ------ | ------------------------ |
| `address`   | String | Endereço IP do remetente |
| `port`      | Number | Porta do remetente       |
| `family`    | String | "IPv4" ou "IPv6"         |

### `'listening'`

Emitido quando o socket começa a escutar.

```javascript
socket.on("listening", () => {
  const address = socket.address();
  console.log(`Socket escutando em ${address.address}:${address.port}`);
});
```

### `'close'`

Emitido quando o socket é fechado.

```javascript
socket.on("close", () => {
  console.log("Socket fechado");
});
```

### `'error'`

Emitido quando ocorre um erro.

```javascript
socket.on("error", (err) => {
  console.error("Erro no socket:", err.message);
  socket.close();
});
```

---

## Métodos do Socket

### `socket.bind(port?, address?, callback?)`

Vincula o socket a uma porta e endereço para receber dados.

```javascript
// Bind em porta específica
socket.bind(41234);

// Bind com endereço
socket.bind(41234, "0.0.0.0");

// Bind com callback
socket.bind(41234, () => {
  console.log("Socket pronto");
});

// Bind com objeto de opções
socket.bind({
  port: 41234,
  address: "0.0.0.0",
});
```

### `socket.send(msg, offset?, length?, port, address?, callback?)`

Envia uma mensagem UDP.

```javascript
// Envio simples
socket.send("Hello", 41234, "127.0.0.1");

// Com Buffer
const message = Buffer.from("Hello UDP");
socket.send(message, 41234, "127.0.0.1");

// Com offset e length
socket.send(message, 0, message.length, 41234, "127.0.0.1");

// Com callback
socket.send("Hello", 41234, "127.0.0.1", (error, bytes) => {
  if (error) {
    console.error("Erro ao enviar:", error.message);
  } else {
    console.log(`${bytes} bytes enviados`);
  }
});

// Enviando múltiplos buffers
socket.send(["Hello", " ", "World"], 41234, "127.0.0.1");
```

### `socket.close(callback?)`

Fecha o socket.

```javascript
socket.close(() => {
  console.log("Socket fechado");
});
```

### `socket.address()`

Retorna informações do endereço local.

```javascript
const address = socket.address();
console.log(`Endereço: ${address.address}`);
console.log(`Porta: ${address.port}`);
console.log(`Família: ${address.family}`);
```

### `socket.setBroadcast(flag)`

Habilita ou desabilita broadcast.

```javascript
socket.setBroadcast(true);
```

### `socket.setTTL(ttl)`

Define o Time-To-Live dos pacotes.

```javascript
socket.setTTL(128);
```

---

## Exemplo: Servidor UDP

```javascript
const dgram = require("dgram");

const server = dgram.createSocket("udp4");

server.on("message", (msg, rinfo) => {
  console.log(`Servidor recebeu: ${msg} de ${rinfo.address}:${rinfo.port}`);

  // Responde ao cliente
  const response = Buffer.from("Mensagem recebida!");
  server.send(response, rinfo.port, rinfo.address);
});

server.on("listening", () => {
  const address = server.address();
  console.log(`Servidor UDP escutando em ${address.address}:${address.port}`);
});

server.on("error", (err) => {
  console.error(`Erro no servidor: ${err.message}`);
  server.close();
});

server.bind(41234);
```

---

## Exemplo: Cliente UDP

```javascript
const dgram = require("dgram");

const client = dgram.createSocket("udp4");

const message = Buffer.from("Olá, servidor UDP!");

client.send(message, 41234, "127.0.0.1", (err) => {
  if (err) {
    console.error("Erro:", err.message);
    client.close();
    return;
  }
  console.log("Mensagem enviada");
});

// Receber resposta do servidor
client.on("message", (msg, rinfo) => {
  console.log(`Resposta do servidor: ${msg}`);
  client.close();
});

// Bind em porta aleatória para receber resposta
client.bind();
```

---

## Exemplo: Broadcast UDP

```javascript
const dgram = require("dgram");

const broadcaster = dgram.createSocket("udp4");

broadcaster.bind(() => {
  broadcaster.setBroadcast(true);

  setInterval(() => {
    const message = Buffer.from("Broadcast message");
    broadcaster.send(message, 41234, "255.255.255.255", (err) => {
      if (err) console.error(err);
      console.log("Mensagem broadcast enviada");
    });
  }, 3000);
});
```

---

## Exemplo: Multicast UDP

```javascript
const dgram = require("dgram");

// Servidor Multicast
const server = dgram.createSocket("udp4");

server.on("message", (msg, rinfo) => {
  console.log(`Servidor: ${msg} de ${rinfo.address}`);
});

server.on("listening", () => {
  // Nota: addMembership pode não estar disponível
  console.log("Servidor multicast pronto");
});

server.bind(5007);
```

---

## UDP vs TCP

| Característica | UDP (dgram)           | TCP (net)           |
| -------------- | --------------------- | ------------------- |
| Conexão        | Connectionless        | Connection-oriented |
| Confiabilidade | Não garantida         | Garantida           |
| Ordem          | Não garantida         | Garantida           |
| Velocidade     | Mais rápido           | Mais lento          |
| Uso            | Games, streaming, DNS | HTTP, FTP, email    |
| Overhead       | Baixo                 | Alto                |

---

## Notas de Compatibilidade

| Funcionalidade    | Nulang | Node.js |
| ----------------- | ------ | ------- |
| `createSocket`    | ✅     | ✅      |
| `bind`            | ✅     | ✅      |
| `send`            | ✅     | ✅      |
| `close`           | ✅     | ✅      |
| `address`         | ✅     | ✅      |
| `setBroadcast`    | ✅     | ✅      |
| `setTTL`          | ✅     | ✅      |
| `addMembership`   | ❌     | ✅      |
| `dropMembership`  | ❌     | ✅      |
| `setMulticastTTL` | ❌     | ✅      |
| `ref/unref`       | ❌     | ✅      |
