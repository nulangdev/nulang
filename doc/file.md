# File Class

A classe `File` representa um arquivo com dados binários, compatível com a Web API. Estende a classe `Blob`.

## Importação

A classe `File` está disponível globalmente e não precisa ser importada.

```javascript
const file = new File(["conteúdo"], "arquivo.txt", { type: "text/plain" });
```

## Construtor

```javascript
new File(fileBits, fileName, [options]);
```

### Parâmetros

| Parâmetro  | Tipo   | Descrição                                                    |
| ---------- | ------ | ------------------------------------------------------------ |
| `fileBits` | Array  | Array de partes do arquivo (strings, Blobs, Buffers, Arrays) |
| `fileName` | String | Nome do arquivo                                              |
| `options`  | Object | Opções opcionais                                             |

### Opções

| Propriedade    | Tipo   | Descrição                               |
| -------------- | ------ | --------------------------------------- |
| `type`         | String | MIME type do arquivo                    |
| `lastModified` | Number | Timestamp da última modificação (em ms) |

## Propriedades

| Propriedade          | Tipo   | Descrição                                  |
| -------------------- | ------ | ------------------------------------------ |
| `name`               | String | Nome do arquivo                            |
| `size`               | Number | Tamanho em bytes                           |
| `type`               | String | MIME type                                  |
| `lastModified`       | Number | Timestamp da última modificação            |
| `webkitRelativePath` | String | Caminho relativo (sempre vazio em Node.js) |

## Métodos

### Herdados de Blob

| Método                           | Retorno                | Descrição                      |
| -------------------------------- | ---------------------- | ------------------------------ |
| `text()`                         | Promise\<String\>      | Lê o conteúdo como texto       |
| `arrayBuffer()`                  | Promise\<ArrayBuffer\> | Lê o conteúdo como ArrayBuffer |
| `slice(start, end, contentType)` | Blob                   | Retorna uma parte do arquivo   |
| `stream()`                       | ReadableStream         | Retorna uma stream do conteúdo |

## Exemplos

### Criando um Arquivo de Texto

```javascript
const conteudo = "Hello, World!";
const arquivo = new File([conteudo], "hello.txt", {
  type: "text/plain",
  lastModified: Date.now(),
});

console.log(arquivo.name); // "hello.txt"
console.log(arquivo.size); // 13
console.log(arquivo.type); // "text/plain"
console.log(arquivo.lastModified); // 1703548800000
```

### Lendo o Conteúdo

```javascript
const file = new File(["Conteúdo do arquivo"], "doc.txt");

// Ler como texto
file.text().then((text) => {
  console.log(text); // "Conteúdo do arquivo"
});

// Ler como ArrayBuffer
file.arrayBuffer().then((buffer) => {
  console.log(buffer.length); // 21
});
```

### Criando de Múltiplas Partes

```javascript
const parte1 = "Primeira parte. ";
const parte2 = "Segunda parte.";
const blob = new Blob(["Terceira parte."]);

const file = new File([parte1, parte2, blob], "combined.txt", {
  type: "text/plain",
});

file.text().then((text) => {
  console.log(text); // "Primeira parte. Segunda parte.Terceira parte."
});
```

### Slicing

```javascript
const file = new File(["ABCDEFGHIJ"], "letters.txt");

// Criar um novo Blob com parte do conteúdo
const slice = file.slice(2, 5);

slice.text().then((text) => {
  console.log(text); // "CDE"
});
```

### Criando de Buffer

```javascript
const buffer = Buffer.from([72, 101, 108, 108, 111]); // "Hello"

const file = new File([buffer], "from-buffer.txt", {
  type: "application/octet-stream",
});

file.text().then((text) => {
  console.log(text); // "Hello"
});
```

### Criando de Outro Blob/File

```javascript
const blob = new Blob(["Conteúdo original"]);
const file = new File([blob], "from-blob.txt");

console.log(file.size); // 18
```

## Diferenças do Browser

| Feature              | Browser        | Nulang       |
| -------------------- | -------------- | ------------ |
| `webkitRelativePath` | Pode ter valor | Sempre vazio |
| Input de arquivo     | Sim            | Não          |
| Drag & Drop          | Sim            | Não          |

## Casos de Uso

### Upload Simulado

```javascript
const file = new File(["dados"], "upload.txt", {
  type: "text/plain",
});

// Simular upload
fetch("http://api.example.com/upload", {
  method: "POST",
  body: file,
});
```

### Processar Conteúdo

```javascript
async function processFile(file) {
  const content = await file.text();
  const lines = content.split("\n");

  console.log(`Arquivo: ${file.name}`);
  console.log(`Linhas: ${lines.length}`);
  console.log(`Tamanho: ${file.size} bytes`);
}

const f = new File(["linha1\nlinha2\nlinha3"], "dados.txt");
processFile(f);
```

### Criar de Dados JSON

```javascript
const dados = { nome: "João", idade: 30 };
const json = JSON.stringify(dados);

const file = new File([json], "dados.json", {
  type: "application/json",
  lastModified: Date.now(),
});

console.log(file.name); // "dados.json"
console.log(file.type); // "application/json"
```

## Veja Também

- [Blob](./blob.md) - Classe base para dados binários
- [Buffer](./buffer.md) - Manipulação de dados binários
- [File System](./filesystem.md) - Operações com arquivos no sistema
