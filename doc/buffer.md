# Buffer

`Buffer` é usado para manipular dados binários diretamente.

## Criação

### Buffer.from(data, encoding)

Cria um Buffer a partir de dados.

```javascript
// De string
const buf1 = Buffer.from("Hello");
console.log(buf1); // <Buffer 48 65 6c 6c 6f>

// De string com encoding
const buf2 = Buffer.from("SGVsbG8=", "base64");
console.log(buf2.toString()); // "Hello"

// De hex
const buf3 = Buffer.from("48656c6c6f", "hex");
console.log(buf3.toString()); // "Hello"

// De array de bytes
const buf4 = Buffer.from([72, 101, 108, 108, 111]);
console.log(buf4.toString()); // "Hello"

// De outro Buffer (cópia)
const buf5 = Buffer.from(buf1);
```

### Buffer.alloc(size, fill, encoding)

Cria um Buffer de tamanho fixo.

```javascript
// Buffer zerado
const buf1 = Buffer.alloc(10);
console.log(buf1); // <Buffer 00 00 00 00 00 00 00 00 00 00>

// Buffer preenchido com valor
const buf2 = Buffer.alloc(10, 255);
console.log(buf2); // <Buffer ff ff ff ff ff ff ff ff ff ff>

// Preenchido com caractere
const buf3 = Buffer.alloc(5, "a");
console.log(buf3.toString()); // "aaaaa"
```

### Buffer.concat(list, totalLength)

Concatena múltiplos Buffers.

```javascript
const buf1 = Buffer.from("Hello ");
const buf2 = Buffer.from("World");
const combined = Buffer.concat([buf1, buf2]);
console.log(combined.toString()); // "Hello World"
```

## Propriedades

| Propriedade | Tipo   | Descrição        |
| ----------- | ------ | ---------------- |
| `length`    | Number | Tamanho em bytes |

```javascript
const buf = Buffer.from("Hello");
console.log(buf.length); // 5
```

## Métodos Estáticos

### Buffer.isBuffer(obj)

Verifica se é um Buffer.

```javascript
console.log(Buffer.isBuffer(Buffer.from("test"))); // true
console.log(Buffer.isBuffer("test")); // false
console.log(Buffer.isBuffer([1, 2, 3])); // false
```

### Buffer.byteLength(string, encoding)

Retorna o tamanho em bytes de uma string.

```javascript
console.log(Buffer.byteLength("Hello")); // 5
console.log(Buffer.byteLength("こんにちは")); // 15 (UTF-8)
```

## Métodos de Instância

### toString(encoding)

Converte Buffer para string.

```javascript
const buf = Buffer.from("Hello");

console.log(buf.toString()); // "Hello" (utf8)
console.log(buf.toString("hex")); // "48656c6c6f"
console.log(buf.toString("base64")); // "SGVsbG8="
```

**Encodings**:
| Encoding | Descrição |
|----------|-----------|
| `utf8` | UTF-8 (padrão) |
| `hex` | Hexadecimal |
| `base64` | Base64 |

### toJSON()

Retorna representação JSON do Buffer.

```javascript
const buf = Buffer.from("Hi");
console.log(buf.toJSON());
// { type: "Buffer", data: [72, 105] }
```

### slice(start, end)

Retorna uma cópia de uma parte do Buffer.

```javascript
const buf = Buffer.from("Hello World");
const slice = buf.slice(0, 5);
console.log(slice.toString()); // "Hello"
```

### copy(target, targetStart, sourceStart, sourceEnd)

Copia dados para outro Buffer.

```javascript
const source = Buffer.from("Hello");
const target = Buffer.alloc(10);

const copied = source.copy(target, 0, 0, 5);
console.log(target.toString()); // "Hello     "
console.log(copied); // 5 (bytes copiados)
```

### equals(otherBuffer)

Compara dois Buffers.

```javascript
const buf1 = Buffer.from("Hello");
const buf2 = Buffer.from("Hello");
const buf3 = Buffer.from("World");

console.log(buf1.equals(buf2)); // true
console.log(buf1.equals(buf3)); // false
```

### fill(value)

Preenche o Buffer com um valor.

```javascript
const buf = Buffer.alloc(5);
buf.fill(65); // 'A' em ASCII
console.log(buf.toString()); // "AAAAA"

buf.fill("x");
console.log(buf.toString()); // "xxxxx"
```

### indexOf(value)

Encontra a posição de um valor.

```javascript
const buf = Buffer.from("Hello World");

console.log(buf.indexOf("o")); // 4
console.log(buf.indexOf("World")); // 6
console.log(buf.indexOf(111)); // 4 (111 = 'o' em ASCII)
console.log(buf.indexOf("xyz")); // -1
```

## Acesso por Índice

```javascript
const buf = Buffer.from("Hello");

// Leitura
console.log(buf[0]); // 72 ('H')
console.log(buf[4]); // 111 ('o')

// Escrita
buf[0] = 74; // 'J'
console.log(buf.toString()); // "Jello"
```

## Exemplos Práticos

### Ler Arquivo Binário

```javascript
const fs = require("fs");

const data = fs.readFileSync("imagem.png");
console.log(`Tamanho: ${data.length} bytes`);
console.log(`Primeiros bytes: ${data.slice(0, 8).toString("hex")}`);
```

### Converter Base64

```javascript
function encodeBase64(text) {
  return Buffer.from(text).toString("base64");
}

function decodeBase64(encoded) {
  return Buffer.from(encoded, "base64").toString();
}

const encoded = encodeBase64("Hello, World!");
console.log(encoded); // "SGVsbG8sIFdvcmxkIQ=="

const decoded = decodeBase64(encoded);
console.log(decoded); // "Hello, World!"
```

### Converter Hex

```javascript
function toHex(text) {
  return Buffer.from(text).toString("hex");
}

function fromHex(hex) {
  return Buffer.from(hex, "hex").toString();
}

console.log(toHex("Hello")); // "48656c6c6f"
console.log(fromHex("48656c6c6f")); // "Hello"
```

### Criar Cabeçalho de Protocolo

```javascript
function createHeader(version, type, length) {
  const header = Buffer.alloc(8);

  header[0] = version;
  header[1] = type;
  // Length como 2 bytes big-endian
  header[2] = (length >> 8) & 0xff;
  header[3] = length & 0xff;
  // 4 bytes reservados

  return header;
}

const header = createHeader(1, 5, 1024);
console.log(header.toString("hex"));
// "01050400000000000"
```

### XOR de Buffers

```javascript
function xor(buf1, buf2) {
  const result = Buffer.alloc(buf1.length);

  for (let i = 0; i < buf1.length; i++) {
    result[i] = buf1[i] ^ buf2[i % buf2.length];
  }

  return result;
}

const data = Buffer.from("Hello");
const key = Buffer.from("key");
const encrypted = xor(data, key);
const decrypted = xor(encrypted, key);

console.log(decrypted.toString()); // "Hello"
```

### Comparar Buffers Byte a Byte

```javascript
function compareBuffers(a, b) {
  if (a.length !== b.length) {
    return false;
  }

  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) {
      return false;
    }
  }

  return true;
}
```

### Concatenar com Delimitador

```javascript
function joinBuffers(buffers, delimiter) {
  if (buffers.length === 0) return Buffer.alloc(0);

  const delim = Buffer.from(delimiter);
  const parts = [];

  for (let i = 0; i < buffers.length; i++) {
    parts.push(buffers[i]);
    if (i < buffers.length - 1) {
      parts.push(delim);
    }
  }

  return Buffer.concat(parts);
}

const result = joinBuffers(
  [Buffer.from("a"), Buffer.from("b"), Buffer.from("c")],
  ","
);

console.log(result.toString()); // "a,b,c"
```

## Buffer vs String

| Característica | Buffer         | String             |
| -------------- | -------------- | ------------------ |
| Mutável        | Sim            | Não                |
| Tipo de dado   | Bytes (0-255)  | Caracteres Unicode |
| Uso            | Dados binários | Texto              |
| Encoding       | Múltiplos      | UTF-16 interno     |

## Boas Práticas

### ✅ Use Buffer.alloc para buffers zerados

```javascript
// ✅ Seguro - zerado
const buf = Buffer.alloc(100);
```

### ✅ Especifique encoding ao converter

```javascript
// ✅ Explícito
const str = buf.toString("utf8");
const hex = buf.toString("hex");
```

### ✅ Valide tamanhos antes de alocar

```javascript
function safeAlloc(size) {
  if (size > 10 * 1024 * 1024) {
    // 10MB
    throw new Error("Buffer muito grande");
  }
  return Buffer.alloc(size);
}
```

## Veja Também

- [File System](./filesystem.md) - Operações com arquivos
- [Crypto](./crypto.md) - Funções criptográficas
- [Stream](./stream.md) - Streams de dados
