# Blob Class

A classe `Blob` representa dados brutos imutáveis, similar a um arquivo. Compatível com a Web API.

## Importação

A classe `Blob` está disponível globalmente e não precisa ser importada.

```javascript
const blob = new Blob(["Hello, World!"], { type: "text/plain" });
```

## Construtor

```javascript
new Blob(blobParts, [options]);
```

### Parâmetros

| Parâmetro   | Tipo   | Descrição                                                           |
| ----------- | ------ | ------------------------------------------------------------------- |
| `blobParts` | Array  | Array de partes (strings, Buffers, outros Blobs, Arrays de números) |
| `options`   | Object | Opções opcionais                                                    |

### Opções

| Propriedade | Tipo   | Descrição         |
| ----------- | ------ | ----------------- |
| `type`      | String | MIME type do blob |

## Propriedades

| Propriedade | Tipo   | Descrição        |
| ----------- | ------ | ---------------- |
| `size`      | Number | Tamanho em bytes |
| `type`      | String | MIME type        |

## Métodos

### text()

Lê o conteúdo do blob como texto.

```javascript
const blob = new Blob(["Hello, World!"]);
blob.text().then((text) => {
  console.log(text); // "Hello, World!"
});
```

**Retorno**: `Promise<String>`

### arrayBuffer()

Lê o conteúdo do blob como ArrayBuffer.

```javascript
const blob = new Blob(["ABC"]);
blob.arrayBuffer().then((buffer) => {
  console.log(buffer.length); // 3
});
```

**Retorno**: `Promise<Buffer>`

### slice(start, end, contentType)

Cria um novo Blob contendo uma parte dos dados.

```javascript
const blob = new Blob(["Hello, World!"]);
const slice = blob.slice(0, 5);
slice.text().then((text) => {
  console.log(text); // "Hello"
});
```

**Parâmetros**:
| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `start` | Number | Índice inicial (opcional, padrão: 0) |
| `end` | Number | Índice final (opcional, padrão: size) |
| `contentType` | String | MIME type do novo Blob (opcional) |

**Retorno**: `Blob`

### stream()

Retorna uma stream ReadableStream do conteúdo.

```javascript
const blob = new Blob(["Streaming data"]);
const stream = blob.stream();
```

**Retorno**: `Object` (stream simplificada)

## Exemplos

### Criar Blob de String

```javascript
const blob = new Blob(["Conteúdo de texto"], {
  type: "text/plain",
});

console.log(blob.size); // 18
console.log(blob.type); // "text/plain"
```

### Criar Blob de Múltiplas Partes

```javascript
const parte1 = "Primeira ";
const parte2 = "Segunda ";
const parte3 = new Blob(["Terceira"]);

const blob = new Blob([parte1, parte2, parte3]);
blob.text().then((text) => {
  console.log(text); // "Primeira Segunda Terceira"
});
```

### Criar Blob de Buffer

```javascript
const buffer = Buffer.from([72, 101, 108, 108, 111]);
const blob = new Blob([buffer]);

blob.text().then((text) => {
  console.log(text); // "Hello"
});
```

### Criar Blob de Array de Números

```javascript
// Array representando bytes
const bytes = [65, 66, 67, 68, 69];
const blob = new Blob([bytes]);

blob.text().then((text) => {
  console.log(text); // "ABCDE"
});
```

### Slicing com Índices Negativos

```javascript
const blob = new Blob(["ABCDEFGHIJ"]);

// Últimos 3 bytes
const slice = blob.slice(-3);
slice.text().then((text) => {
  console.log(text); // "HIJ"
});
```

### Combinar Blobs

```javascript
const blob1 = new Blob(["Hello, "]);
const blob2 = new Blob(["World!"]);

const combined = new Blob([blob1, blob2]);
combined.text().then((text) => {
  console.log(text); // "Hello, World!"
});
```

### Converter para Base64

```javascript
const blob = new Blob(["Hello"]);
blob.arrayBuffer().then((buffer) => {
  const base64 = buffer.toString("base64");
  console.log(base64); // "SGVsbG8="
});
```

## Padrões de Uso

### JSON Blob

```javascript
const dados = { nome: "João", idade: 30 };
const blob = new Blob([JSON.stringify(dados)], {
  type: "application/json",
});

blob.text().then((text) => {
  const obj = JSON.parse(text);
  console.log(obj.nome); // "João"
});
```

### Binary Data

```javascript
// Criar blob com dados binários
const header = new Uint8Array([0x89, 0x50, 0x4e, 0x47]);
const blob = new Blob([header], {
  type: "application/octet-stream",
});

console.log(blob.size); // 4
```

### Processar Chunk por Chunk

```javascript
async function processBlockByBlock(blob, chunkSize) {
  const chunks = [];
  let offset = 0;

  while (offset < blob.size) {
    const chunk = blob.slice(offset, offset + chunkSize);
    const text = await chunk.text();
    chunks.push(text);
    offset += chunkSize;
  }

  return chunks;
}

const blob = new Blob(["AABBCCDD"]);
processBlockByBlock(blob, 2).then((chunks) => {
  console.log(chunks); // ["AA", "BB", "CC", "DD"]
});
```

## Comparação com Buffer

| Feature       | Blob              | Buffer              |
| ------------- | ----------------- | ------------------- |
| Imutável      | Sim               | Não                 |
| Async read    | Sim (Promise)     | Não (sync)          |
| MIME type     | Sim               | Não                 |
| Slicing       | Retorna novo Blob | Retorna novo Buffer |
| Uso principal | Web APIs          | Node.js APIs        |

## Veja Também

- [File](./file.md) - Estende Blob com nome e metadados
- [Buffer](./buffer.md) - Manipulação de dados binários
- [File System](./filesystem.md) - Operações com arquivos
