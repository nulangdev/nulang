# Regular Expressions (RegExp)

Expressões regulares para correspondência de padrões em strings.

## Criação

### Usando o Construtor

```javascript
const regex = RegExp("padrão");
const regexWithFlags = RegExp("padrão", "gi");
```

### Com new

```javascript
const regex = new RegExp("padrão");
const regexWithFlags = new RegExp("padrão", "gi");
```

## Flags

| Flag | Nome       | Descrição                          |
| ---- | ---------- | ---------------------------------- |
| `g`  | global     | Encontra todas as correspondências |
| `i`  | ignoreCase | Ignora maiúsculas/minúsculas       |
| `m`  | multiline  | ^ e $ correspondem a linhas        |
| `s`  | dotAll     | . corresponde a \n                 |

```javascript
const regex = RegExp("hello", "gi");
console.log(regex.global); // true
console.log(regex.ignoreCase); // true
```

## Propriedades

| Propriedade  | Tipo    | Descrição         |
| ------------ | ------- | ----------------- |
| `source`     | String  | Padrão da regex   |
| `flags`      | String  | Flags usadas      |
| `global`     | Boolean | Flag g está ativa |
| `ignoreCase` | Boolean | Flag i está ativa |
| `multiline`  | Boolean | Flag m está ativa |

```javascript
const regex = RegExp("test", "gi");
console.log(regex.source); // "test"
console.log(regex.flags); // "gi"
```

## Métodos da RegExp

### test(string)

Verifica se a regex corresponde à string.

```javascript
const regex = RegExp("hello");
console.log(regex.test("hello world")); // true
console.log(regex.test("hi there")); // false
```

**Retorno**: Boolean

### exec(string)

Executa a busca e retorna detalhes da correspondência.

```javascript
const regex = RegExp("(\\d+)");
const result = regex.exec("Idade: 25 anos");

if (result) {
  console.log(result[0]); // "25" (match completo)
  console.log(result[1]); // "25" (grupo 1)
}
```

**Retorno**: Array ou null

### match(string)

Encontra todas as correspondências.

```javascript
const regex = RegExp("\\d+", "g");
const matches = regex.match("Itens: 10, 20, 30");
console.log(matches); // ["10", "20", "30"]
```

### replace(string, replacement)

Substitui correspondências.

```javascript
const regex = RegExp("\\s+", "g");
const result = regex.replace("hello   world", " ");
console.log(result); // "hello world"
```

### split(string, limit)

Divide a string pelo padrão.

```javascript
const regex = RegExp("[,;]");
const parts = regex.split("a,b;c,d");
console.log(parts); // ["a", "b", "c", "d"]
```

## Métodos de String com RegExp

### string.match(regex)

```javascript
const str = "The rain in SPAIN";
const regex = RegExp("ain", "gi");
const matches = str.match(regex);
console.log(matches); // ["ain", "AIN", "ain"]
```

### string.replace(regex, replacement)

```javascript
const str = "Hello World";
const regex = RegExp("world", "i");
const result = str.replace(regex, "Nulang");
console.log(result); // "Hello Nulang"
```

### string.search(regex)

Retorna o índice da primeira correspondência.

```javascript
const str = "Hello 123 World";
const regex = RegExp("\\d+");
console.log(str.search(regex)); // 6
```

### string.replaceAll(search, replacement)

Substitui todas as ocorrências.

```javascript
const str = "a-b-c-d";
console.log(str.replaceAll("-", "_")); // "a_b_c_d"
```

## Sintaxe de Padrões

### Caracteres Especiais

| Padrão | Descrição                                |
| ------ | ---------------------------------------- |
| `.`    | Qualquer caractere (exceto \n)           |
| `\d`   | Dígito (0-9)                             |
| `\w`   | Caractere de palavra (a-z, A-Z, 0-9, \_) |
| `\s`   | Espaço em branco                         |
| `\D`   | Não-dígito                               |
| `\W`   | Não-palavra                              |
| `\S`   | Não-espaço                               |

```javascript
const digits = RegExp("\\d+");
console.log(digits.test("abc123")); // true

const words = RegExp("\\w+", "g");
console.log("Hello World".match(words)); // ["Hello", "World"]
```

### Âncoras

| Padrão | Descrição              |
| ------ | ---------------------- |
| `^`    | Início da string/linha |
| `$`    | Fim da string/linha    |
| `\b`   | Limite de palavra      |

```javascript
const startsWith = RegExp("^Hello");
console.log(startsWith.test("Hello World")); // true
console.log(startsWith.test("Say Hello")); // false

const endsWith = RegExp("World$");
console.log(endsWith.test("Hello World")); // true
```

### Quantificadores

| Padrão  | Descrição    |
| ------- | ------------ |
| `*`     | 0 ou mais    |
| `+`     | 1 ou mais    |
| `?`     | 0 ou 1       |
| `{n}`   | Exatamente n |
| `{n,}`  | n ou mais    |
| `{n,m}` | Entre n e m  |

```javascript
const oneOrMore = RegExp("a+");
console.log(oneOrMore.test("baaab")); // true

const exactly3 = RegExp("\\d{3}");
console.log(exactly3.test("12")); // false
console.log(exactly3.test("123")); // true
```

### Grupos e Classes

| Padrão    | Descrição                  |
| --------- | -------------------------- |
| `[abc]`   | Qualquer um de a, b, c     |
| `[^abc]`  | Qualquer um EXCETO a, b, c |
| `[a-z]`   | Range de a até z           |
| `(abc)`   | Grupo de captura           |
| `(?:abc)` | Grupo sem captura          |
| `a\|b`    | a OU b                     |

```javascript
const vowels = RegExp("[aeiou]", "gi");
console.log("Hello".match(vowels)); // ["e", "o"]

const range = RegExp("[A-Z]", "g");
console.log("Hello World".match(range)); // ["H", "W"]

const group = RegExp("(\\d{2})-(\\d{2})-(\\d{4})");
const result = group.exec("15-06-2024");
console.log(result[1]); // "15"
console.log(result[2]); // "06"
console.log(result[3]); // "2024"
```

## Exemplos Práticos

### Validar Email

```javascript
function isValidEmail(email) {
  const regex = RegExp("^[\\w.-]+@[\\w.-]+\\.[a-z]{2,}$", "i");
  return regex.test(email);
}

console.log(isValidEmail("test@example.com")); // true
console.log(isValidEmail("invalid-email")); // false
```

### Validar Telefone

```javascript
function isValidPhone(phone) {
  // Formato: (99) 99999-9999 ou 99999999999
  const regex = RegExp("^\\(?\\d{2}\\)?\\s?\\d{4,5}-?\\d{4}$");
  return regex.test(phone);
}

console.log(isValidPhone("(11) 99999-9999")); // true
console.log(isValidPhone("11999999999")); // true
```

### Extrair Números

```javascript
function extractNumbers(str) {
  const regex = RegExp("\\d+", "g");
  return str.match(regex) || [];
}

console.log(extractNumbers("Item 1: R$ 50,00 | Item 2: R$ 100,00"));
// ["1", "50", "00", "2", "100", "00"]
```

### Remover Tags HTML

```javascript
function stripHtml(html) {
  const regex = RegExp("<[^>]*>", "g");
  return html.replace(regex, "");
}

console.log(stripHtml("<p>Hello <b>World</b></p>"));
// "Hello World"
```

### Validar CPF (formato)

```javascript
function isValidCpfFormat(cpf) {
  const regex = RegExp("^\\d{3}\\.\\d{3}\\.\\d{3}-\\d{2}$");
  return regex.test(cpf);
}

console.log(isValidCpfFormat("123.456.789-00")); // true
console.log(isValidCpfFormat("12345678900")); // false
```

### Extrair URLs

```javascript
function extractUrls(text) {
  const regex = RegExp("https?://[\\w.-]+(?:/[\\w.-]*)*", "g");
  return text.match(regex) || [];
}

console.log(extractUrls("Visite https://example.com e http://test.org"));
// ["https://example.com", "http://test.org"]
```

### Substituição com Grupos

```javascript
// Inverter data de DD/MM/YYYY para YYYY-MM-DD
function formatDate(dateStr) {
  const regex = RegExp("(\\d{2})/(\\d{2})/(\\d{4})");
  const match = regex.exec(dateStr);
  if (match) {
    return `${match[3]}-${match[2]}-${match[1]}`;
  }
  return dateStr;
}

console.log(formatDate("25/12/2024")); // "2024-12-25"
```

### Tokenização

```javascript
function tokenize(expression) {
  const regex = RegExp("[\\d.]+|[+\\-*/()]", "g");
  return expression.match(regex) || [];
}

console.log(tokenize("(10 + 20) * 3"));
// ["(", "10", "+", "20", ")", "*", "3"]
```

## Escapando Caracteres Especiais

Caracteres que precisam ser escapados: `. \ + * ? [ ^ ] $ ( ) { } = ! < > | : -`

```javascript
// Errado - . é especial
const wrong = RegExp("file.txt");
console.log(wrong.test("filextxt")); // true (indesejado)

// Correto
const correct = RegExp("file\\.txt");
console.log(correct.test("filextxt")); // false
console.log(correct.test("file.txt")); // true
```

## Veja Também

- [String](./string.md) - Métodos de string
- [Array](./array.md) - Trabalhar com resultados de match
