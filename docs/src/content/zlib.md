# Módulo Zlib

O módulo `zlib` fornece funcionalidades de compressão e descompressão usando Gzip e Deflate/Inflate.

## Importação

```javascript
const zlib = require("zlib");
```

## Compressão Gzip

### `zlib.gzip(buffer, callback)`

Comprime dados usando gzip.

```javascript
const data = Buffer.from("Hello World!");

zlib.gzip(data, (err, compressed) => {
  console.log("Original:", data.length, "bytes");
  console.log("Comprimido:", compressed.length, "bytes");
});
```

### `zlib.gzipSync(buffer)`

Versão síncrona de gzip.

```javascript
const data = Buffer.from("Hello World!");
const compressed = zlib.gzipSync(data);
console.log("Compressed:", compressed);
```

## Descompressão Gzip

### `zlib.gunzip(buffer, callback)`

Descomprime dados gzip.

```javascript
zlib.gzip(Buffer.from("Hello"), (err, compressed) => {
  zlib.gunzip(compressed, (err, decompressed) => {
    console.log(decompressed.toString()); // "Hello"
  });
});
```

### `zlib.gunzipSync(buffer)`

Versão síncrona.

```javascript
const compressed = zlib.gzipSync(Buffer.from("Hello"));
const original = zlib.gunzipSync(compressed);
console.log(original.toString()); // "Hello"
```

## Compressão Deflate

### `zlib.deflate(buffer, callback)`

Comprime usando deflate.

```javascript
const data = Buffer.from("Dados para comprimir");

zlib.deflate(data, (err, compressed) => {
  console.log("Comprimido:", compressed.length, "bytes");
});
```

### `zlib.deflateSync(buffer)`

Versão síncrona.

```javascript
const compressed = zlib.deflateSync(Buffer.from("Hello"));
```

## Descompressão Inflate

### `zlib.inflate(buffer, callback)`

Descomprime dados deflate.

```javascript
zlib.deflate(Buffer.from("Test"), (err, compressed) => {
  zlib.inflate(compressed, (err, decompressed) => {
    console.log(decompressed.toString());
  });
});
```

### `zlib.inflateSync(buffer)`

Versão síncrona.

```javascript
const compressed = zlib.deflateSync(Buffer.from("Hello"));
const original = zlib.inflateSync(compressed);
console.log(original.toString());
```

## Exemplos Práticos

### Comprimir Arquivo

```javascript
const fs = require("fs");
const zlib = require("zlib");

const input = fs.readFileSync("arquivo.txt");
const compressed = zlib.gzipSync(input);
fs.writeFileSync("arquivo.txt.gz", compressed);

console.log("Arquivo comprimido!");
```

### Descomprimir Arquivo

```javascript
const compressed = fs.readFileSync("arquivo.txt.gz");
const decompressed = zlib.gunzipSync(compressed);
fs.writeFileSync("arquivo_restaurado.txt", decompressed);
```

### Compressão de String JSON

```javascript
const data = { users: [], config: {} };
const json = JSON.stringify(data);

const compressed = zlib.gzipSync(Buffer.from(json));
console.log(`JSON: ${json.length} bytes`);
console.log(`Comprimido: ${compressed.length} bytes`);
console.log(
  `Redução: ${Math.round((1 - compressed.length / json.length) * 100)}%`
);
```

## Comparação Gzip vs Deflate

| Aspecto  | Gzip           | Deflate  |
| -------- | -------------- | -------- |
| Header   | Sim            | Não      |
| Checksum | CRC32          | Adler32  |
| Uso      | Arquivos, HTTP | ZIP, PNG |
| Overhead | Maior          | Menor    |

## Compatibilidade

| Funcionalidade   | Nulang | Node.js |
| ---------------- | ------ | ------- |
| `gzip`           | ✅     | ✅      |
| `gzipSync`       | ✅     | ✅      |
| `gunzip`         | ✅     | ✅      |
| `gunzipSync`     | ✅     | ✅      |
| `deflate`        | ✅     | ✅      |
| `deflateSync`    | ✅     | ✅      |
| `inflate`        | ✅     | ✅      |
| `inflateSync`    | ✅     | ✅      |
| `brotliCompress` | ❌     | ✅      |
| Streams          | ❌     | ✅      |
