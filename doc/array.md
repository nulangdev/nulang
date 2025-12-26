# Array

O tipo `Array` representa uma lista ordenada de elementos, compatível com JavaScript.

## Criação

```javascript
// Array literal
const arr = [1, 2, 3, 4, 5];

// Array vazio
const empty = [];

// Array misto
const mixed = [1, "texto", true, { key: "value" }];

// Com Array.from
const chars = Array.from("hello");
console.log(chars); // ["h", "e", "l", "l", "o"]
```

## Propriedades

| Propriedade | Tipo   | Descrição                    |
| ----------- | ------ | ---------------------------- |
| `length`    | Number | Número de elementos no array |

```javascript
const arr = [1, 2, 3];
console.log(arr.length); // 3
```

## Métodos de Mutação

### push(...elements)

Adiciona elementos ao final do array.

```javascript
const arr = [1, 2];
arr.push(3, 4);
console.log(arr); // [1, 2, 3, 4]
```

**Retorno**: Number (novo length)

### pop()

Remove e retorna o último elemento.

```javascript
const arr = [1, 2, 3];
const last = arr.pop();
console.log(last); // 3
console.log(arr); // [1, 2]
```

**Retorno**: O elemento removido ou `undefined`

### shift()

Remove e retorna o primeiro elemento.

```javascript
const arr = [1, 2, 3];
const first = arr.shift();
console.log(first); // 1
console.log(arr); // [2, 3]
```

**Retorno**: O elemento removido ou `undefined`

### unshift(...elements)

Adiciona elementos no início do array.

```javascript
const arr = [2, 3];
arr.unshift(0, 1);
console.log(arr); // [0, 1, 2, 3]
```

**Retorno**: Number (novo length)

### splice(start, deleteCount, ...items)

Remove/adiciona elementos em uma posição específica.

```javascript
const arr = [1, 2, 3, 4, 5];

// Remover 2 elementos a partir do índice 1
const removed = arr.splice(1, 2);
console.log(removed); // [2, 3]
console.log(arr); // [1, 4, 5]

// Inserir elementos
arr.splice(1, 0, "a", "b");
console.log(arr); // [1, "a", "b", 4, 5]

// Substituir elementos
arr.splice(1, 2, "x");
console.log(arr); // [1, "x", 4, 5]
```

**Retorno**: Array com elementos removidos

### reverse()

Inverte a ordem dos elementos (modifica o array original).

```javascript
const arr = [1, 2, 3];
arr.reverse();
console.log(arr); // [3, 2, 1]
```

**Retorno**: O array invertido

## Métodos de Iteração

### map(callback)

Cria um novo array transformando cada elemento.

```javascript
const nums = [1, 2, 3, 4];
const doubled = nums.map((x) => x * 2);
console.log(doubled); // [2, 4, 6, 8]

// Com índice
const indexed = nums.map((x, i) => `${i}: ${x}`);
console.log(indexed); // ["0: 1", "1: 2", "2: 3", "3: 4"]
```

**Callback**: `(element, index, array) => newElement`
**Retorno**: Novo Array

### filter(callback)

Cria um novo array com elementos que passam no teste.

```javascript
const nums = [1, 2, 3, 4, 5, 6];
const evens = nums.filter((x) => x % 2 === 0);
console.log(evens); // [2, 4, 6]
```

**Callback**: `(element, index, array) => boolean`
**Retorno**: Novo Array

### forEach(callback)

Executa uma função para cada elemento.

```javascript
const arr = ["a", "b", "c"];
arr.forEach((item, index) => {
  console.log(`${index}: ${item}`);
});
// Output:
// 0: a
// 1: b
// 2: c
```

**Callback**: `(element, index, array) => void`
**Retorno**: `undefined`

### find(callback)

Retorna o primeiro elemento que passa no teste.

```javascript
const users = [
  { name: "João", age: 25 },
  { name: "Maria", age: 30 },
  { name: "Pedro", age: 35 },
];

const found = users.find((u) => u.age > 28);
console.log(found); // { name: "Maria", age: 30 }
```

**Callback**: `(element, index, array) => boolean`
**Retorno**: Elemento encontrado ou `undefined`

### findIndex(callback)

Retorna o índice do primeiro elemento que passa no teste.

```javascript
const nums = [5, 12, 8, 130, 44];
const index = nums.findIndex((x) => x > 10);
console.log(index); // 1
```

**Callback**: `(element, index, array) => boolean`
**Retorno**: Number (índice) ou -1

### reduce(callback, initialValue)

Reduz o array a um único valor.

```javascript
const nums = [1, 2, 3, 4];

// Soma
const sum = nums.reduce((acc, x) => acc + x, 0);
console.log(sum); // 10

// Produto
const product = nums.reduce((acc, x) => acc * x, 1);
console.log(product); // 24

// Objeto
const grouped = nums.reduce(
  (acc, x) => {
    acc[x % 2 === 0 ? "even" : "odd"].push(x);
    return acc;
  },
  { even: [], odd: [] }
);
console.log(grouped); // { even: [2, 4], odd: [1, 3] }
```

**Callback**: `(accumulator, element, index, array) => newAccumulator`
**Retorno**: Valor acumulado

## Métodos de Busca

### includes(element)

Verifica se o array contém um elemento.

```javascript
const arr = [1, 2, 3];
console.log(arr.includes(2)); // true
console.log(arr.includes(5)); // false
```

**Retorno**: Boolean

## Métodos de Transformação

### join(separator)

Junta todos os elementos em uma string.

```javascript
const arr = ["a", "b", "c"];
console.log(arr.join()); // "a,b,c"
console.log(arr.join("-")); // "a-b-c"
console.log(arr.join("")); // "abc"
```

**Retorno**: String

### slice(start, end)

Retorna uma cópia superficial de uma parte do array.

```javascript
const arr = [1, 2, 3, 4, 5];

console.log(arr.slice(1, 3)); // [2, 3]
console.log(arr.slice(2)); // [3, 4, 5]
console.log(arr.slice(-2)); // [4, 5]
console.log(arr.slice(0, -1)); // [1, 2, 3, 4]
```

**Retorno**: Novo Array

### concat(...arrays)

Combina arrays em um novo array.

```javascript
const a = [1, 2];
const b = [3, 4];
const c = a.concat(b, [5, 6]);
console.log(c); // [1, 2, 3, 4, 5, 6]
```

**Retorno**: Novo Array

## Objeto Global Array

### Array.isArray(value)

Verifica se um valor é um array.

```javascript
console.log(Array.isArray([1, 2, 3])); // true
console.log(Array.isArray("hello")); // false
console.log(Array.isArray({ 0: 1 })); // false
```

### Array.from(iterable)

Cria um array a partir de um iterável.

```javascript
// De string
console.log(Array.from("abc")); // ["a", "b", "c"]

// De outro array (cópia)
const original = [1, 2, 3];
const copy = Array.from(original);
console.log(copy); // [1, 2, 3]
```

## Exemplos Práticos

### Filtrar e Mapear

```javascript
const nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

const result = nums.filter((x) => x % 2 === 0).map((x) => x * 2);

console.log(result); // [4, 8, 12, 16, 20]
```

### Calcular Estatísticas

```javascript
const grades = [85, 90, 78, 92, 88];

const sum = grades.reduce((a, b) => a + b, 0);
const avg = sum / grades.length;
const max = grades.reduce((a, b) => (a > b ? a : b));
const min = grades.reduce((a, b) => (a < b ? a : b));

console.log(`Sum: ${sum}, Avg: ${avg}, Max: ${max}, Min: ${min}`);
// Sum: 433, Avg: 86.6, Max: 92, Min: 78
```

### Achatar Array Aninhado

```javascript
const nested = [[1, 2], [3, 4], [5]];
const flat = nested.reduce((acc, x) => acc.concat(x), []);
console.log(flat); // [1, 2, 3, 4, 5]
```

### Agrupar por Propriedade

```javascript
const people = [
  { name: "João", age: 25 },
  { name: "Maria", age: 30 },
  { name: "Pedro", age: 25 },
];

const grouped = people.reduce((acc, person) => {
  const age = person.age;
  if (!acc[age]) acc[age] = [];
  acc[age].push(person);
  return acc;
}, {});

console.log(grouped);
// {
//   25: [{ name: "João", age: 25 }, { name: "Pedro", age: 25 }],
//   30: [{ name: "Maria", age: 30 }]
// }
```

## Veja Também

- [String](./string.md) - Manipulação de texto
- [Map](./map.md) - Coleção de pares chave-valor
- [Set](./map.md#set) - Coleção de valores únicos
