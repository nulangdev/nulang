# String

O tipo `String` representa uma sequência de caracteres, compatível com JavaScript.

## Criação

```javascript
// String literal
const str1 = "Hello, World!";
const str2 = "Single quotes";

// Template literal (interpolação)
const name = "João";
const greeting = `Hello, ${name}!`;
console.log(greeting); // "Hello, João!"
```

## Propriedades

| Propriedade | Tipo   | Descrição            |
| ----------- | ------ | -------------------- |
| `length`    | Number | Número de caracteres |

```javascript
const str = "Hello";
console.log(str.length); // 5
```

## Métodos de Caso

### toUpperCase()

Converte para maiúsculas.

```javascript
const str = "hello";
console.log(str.toUpperCase()); // "HELLO"
```

### toLowerCase()

Converte para minúsculas.

```javascript
const str = "HELLO";
console.log(str.toLowerCase()); // "hello"
```

## Métodos de Busca

### indexOf(searchValue)

Retorna o índice da primeira ocorrência.

```javascript
const str = "Hello, World!";
console.log(str.indexOf("o")); // 4
console.log(str.indexOf("World")); // 7
console.log(str.indexOf("xyz")); // -1
```

### includes(searchValue)

Verifica se contém a substring.

```javascript
const str = "Hello, World!";
console.log(str.includes("World")); // true
console.log(str.includes("xyz")); // false
```

### startsWith(searchValue)

Verifica se começa com a substring.

```javascript
const str = "Hello, World!";
console.log(str.startsWith("Hello")); // true
console.log(str.startsWith("World")); // false
```

### endsWith(searchValue)

Verifica se termina com a substring.

```javascript
const str = "Hello, World!";
console.log(str.endsWith("!")); // true
console.log(str.endsWith("World")); // false
console.log(str.endsWith("d!")); // true
```

### search(regexp)

Procura por uma expressão regular.

```javascript
const str = "Hello123World";
const regex = RegExp("[0-9]+");
console.log(str.search(regex)); // 5
```

## Métodos de Extração

### charAt(index)

Retorna o caractere na posição.

```javascript
const str = "Hello";
console.log(str.charAt(0)); // "H"
console.log(str.charAt(4)); // "o"
console.log(str.charAt(10)); // ""
```

### slice(start, end)

Extrai uma parte da string.

```javascript
const str = "Hello, World!";

console.log(str.slice(0, 5)); // "Hello"
console.log(str.slice(7)); // "World!"
console.log(str.slice(-6)); // "World!"
console.log(str.slice(0, -1)); // "Hello, World"
console.log(str.slice(-6, -1)); // "World"
```

### substring(start, end)

Similar ao slice, mas trata índices negativos diferente.

```javascript
const str = "Hello, World!";

console.log(str.substring(0, 5)); // "Hello"
console.log(str.substring(7)); // "World!"

// Diferença: substring troca start e end se start > end
console.log(str.substring(5, 0)); // "Hello"
```

## Métodos de Transformação

### trim()

Remove espaços do início e fim.

```javascript
const str = "   Hello   ";
console.log(str.trim()); // "Hello"
```

### split(separator)

Divide a string em um array.

```javascript
const str = "a,b,c,d";
console.log(str.split(",")); // ["a", "b", "c", "d"]

const words = "Hello World";
console.log(words.split(" ")); // ["Hello", "World"]

// Cada caractere
console.log("abc".split("")); // ["a", "b", "c"]
```

### replace(search, replacement)

Substitui a primeira ocorrência.

```javascript
const str = "Hello, World!";
console.log(str.replace("World", "Nulang")); // "Hello, Nulang!"

// Com regex
const regex = RegExp("[0-9]+");
const text = "Item 123 Price 456";
console.log(text.replace(regex, "XXX")); // "Item XXX Price 456"
```

### replaceAll(search, replacement)

Substitui todas as ocorrências.

```javascript
const str = "a-b-c-d";
console.log(str.replaceAll("-", "_")); // "a_b_c_d"
```

### repeat(count)

Repete a string n vezes.

```javascript
const str = "abc";
console.log(str.repeat(3)); // "abcabcabc"
console.log(str.repeat(0)); // ""
```

### concat(...strings)

Concatena strings.

```javascript
const str = "Hello";
console.log(str.concat(", ", "World", "!")); // "Hello, World!"
```

## Métodos com Regex

### match(regexp)

Encontra correspondências com regex.

```javascript
const str = "The rain in SPAIN stays mainly in the plain";
const regex = RegExp("ain", "gi"); // global, case-insensitive

const matches = str.match(regex);
console.log(matches); // ["ain", "AIN", "ain", "ain"]
```

## Template Literals

Template literals permitem interpolação de expressões.

```javascript
const name = "João";
const age = 30;

// Interpolação simples
console.log(`Nome: ${name}`); // "Nome: João"

// Expressões
console.log(`Idade: ${age * 2}`); // "Idade: 60"

// Multi-linha
const html = `
  <div>
    <h1>${name}</h1>
    <p>Idade: ${age}</p>
  </div>
`;

// Expressões complexas
const items = ["a", "b", "c"];
console.log(`Items: ${items.join(", ")}`); // "Items: a, b, c"
```

## Exemplos Práticos

### Capitalizar Primeira Letra

```javascript
function capitalize(str) {
  if (str.length === 0) return str;
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
}

console.log(capitalize("hello")); // "Hello"
console.log(capitalize("WORLD")); // "World"
```

### Contar Palavras

```javascript
function countWords(str) {
  return str
    .trim()
    .split(" ")
    .filter((w) => w.length > 0).length;
}

console.log(countWords("Hello World")); // 2
console.log(countWords("  Multiple   spaces  ")); // 2
```

### Reverter String

```javascript
function reverse(str) {
  return str.split("").reverse().join("");
}

console.log(reverse("hello")); // "olleh"
```

### Truncar com Ellipsis

```javascript
function truncate(str, maxLength) {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength - 3) + "...";
}

console.log(truncate("Hello, World!", 10)); // "Hello, ..."
console.log(truncate("Short", 10)); // "Short"
```

### Slug URL

```javascript
function slugify(str) {
  return str.toLowerCase().trim().replace(" ", "-").replaceAll(" ", "-");
}

console.log(slugify("Hello World Test")); // "hello-world-test"
```

### Validar Email (Básico)

```javascript
function isValidEmail(email) {
  return email.includes("@") && email.includes(".");
}

console.log(isValidEmail("test@example.com")); // true
console.log(isValidEmail("invalid")); // false
```

## Objeto Global String

### String(value)

Converte um valor para string.

```javascript
console.log(String(123)); // "123"
console.log(String(true)); // "true"
console.log(String(null)); // "null"
console.log(String(undefined)); // "undefined"
```

## Acesso por Índice

```javascript
const str = "Hello";
console.log(str[0]); // "H"
console.log(str[4]); // "o"
```

## Iteração

```javascript
const str = "ABC";

// forEach via Array.from
Array.from(str).forEach((char, i) => {
  console.log(`${i}: ${char}`);
});
// 0: A
// 1: B
// 2: C
```

## Veja Também

- [Array](./array.md) - Listas ordenadas
- [RegExp](./regex.md) - Expressões regulares
- [Numbers](./numbers.md) - Números
