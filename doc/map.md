# Map e Set

Estruturas de dados para coleções, compatíveis com JavaScript ES6.

## Map

`Map` é uma coleção de pares chave-valor onde chaves podem ser de qualquer tipo.

### Criação

```javascript
// Map vazio
const map = new Map();

// Map de array de pares
const map2 = new Map([
  ["chave1", "valor1"],
  ["chave2", "valor2"],
]);
```

### Propriedades

| Propriedade | Tipo   | Descrição          |
| ----------- | ------ | ------------------ |
| `size`      | Number | Número de entradas |

### Métodos

#### set(key, value)

Adiciona ou atualiza uma entrada.

```javascript
const map = new Map();
map.set("nome", "João");
map.set("idade", 30);
map.set("nome", "Maria"); // Atualiza

console.log(map.size); // 2
```

**Retorno**: O próprio Map (permite encadeamento)

```javascript
map.set("a", 1).set("b", 2).set("c", 3);
```

#### get(key)

Obtém o valor associado à chave.

```javascript
const map = new Map();
map.set("nome", "João");

console.log(map.get("nome")); // "João"
console.log(map.get("outro")); // undefined
```

#### has(key)

Verifica se a chave existe.

```javascript
const map = new Map();
map.set("nome", "João");

console.log(map.has("nome")); // true
console.log(map.has("idade")); // false
```

#### delete(key)

Remove uma entrada.

```javascript
const map = new Map();
map.set("nome", "João");
map.delete("nome");

console.log(map.has("nome")); // false
```

**Retorno**: Boolean (true se existia)

#### clear()

Remove todas as entradas.

```javascript
const map = new Map();
map.set("a", 1);
map.set("b", 2);
map.clear();

console.log(map.size); // 0
```

#### keys()

Retorna um array com as chaves.

```javascript
const map = new Map();
map.set("a", 1);
map.set("b", 2);

console.log(map.keys()); // ["a", "b"]
```

#### values()

Retorna um array com os valores.

```javascript
const map = new Map();
map.set("a", 1);
map.set("b", 2);

console.log(map.values()); // [1, 2]
```

#### entries()

Retorna um array de pares [chave, valor].

```javascript
const map = new Map();
map.set("a", 1);
map.set("b", 2);

console.log(map.entries()); // [["a", 1], ["b", 2]]
```

#### forEach(callback)

Itera sobre cada entrada.

```javascript
const map = new Map();
map.set("a", 1);
map.set("b", 2);

map.forEach((value, key, map) => {
  console.log(`${key}: ${value}`);
});
// a: 1
// b: 2
```

### Vantagens sobre Object

```javascript
// Map aceita qualquer tipo como chave
const objKey = { id: 1 };
const map = new Map();
map.set(objKey, "valor");
console.log(map.get(objKey)); // "valor"

// Map mantém ordem de inserção
const ordered = new Map();
ordered.set("z", 1);
ordered.set("a", 2);
console.log(ordered.keys()); // ["z", "a"]
```

---

## Set

`Set` é uma coleção de valores únicos.

### Criação

```javascript
// Set vazio
const set = new Set();

// Set de array
const set2 = new Set([1, 2, 3, 2, 1]);
console.log(set2.size); // 3 (duplicatas removidas)
```

### Propriedades

| Propriedade | Tipo   | Descrição           |
| ----------- | ------ | ------------------- |
| `size`      | Number | Número de elementos |

### Métodos

#### add(value)

Adiciona um valor (ignorado se já existe).

```javascript
const set = new Set();
set.add(1);
set.add(2);
set.add(1); // Ignorado

console.log(set.size); // 2
```

**Retorno**: O próprio Set (permite encadeamento)

#### has(value)

Verifica se o valor existe.

```javascript
const set = new Set([1, 2, 3]);

console.log(set.has(2)); // true
console.log(set.has(5)); // false
```

#### delete(value)

Remove um valor.

```javascript
const set = new Set([1, 2, 3]);
set.delete(2);

console.log(set.has(2)); // false
```

**Retorno**: Boolean (true se existia)

#### clear()

Remove todos os valores.

```javascript
const set = new Set([1, 2, 3]);
set.clear();

console.log(set.size); // 0
```

#### values() / keys()

Retorna um array com os valores.

```javascript
const set = new Set([1, 2, 3]);
console.log(set.values()); // [1, 2, 3]
console.log(set.keys()); // [1, 2, 3] (mesmo que values)
```

#### entries()

Retorna pares [valor, valor] (para compatibilidade com Map).

```javascript
const set = new Set([1, 2]);
console.log(set.entries()); // [[1, 1], [2, 2]]
```

#### forEach(callback)

Itera sobre cada valor.

```javascript
const set = new Set([1, 2, 3]);

set.forEach((value) => {
  console.log(value);
});
// 1
// 2
// 3
```

### Operações de Conjunto

#### union(otherSet)

União de dois conjuntos.

```javascript
const a = new Set([1, 2, 3]);
const b = new Set([3, 4, 5]);

const uniao = a.union(b);
console.log(uniao.values()); // [1, 2, 3, 4, 5]
```

#### intersection(otherSet)

Interseção de dois conjuntos.

```javascript
const a = new Set([1, 2, 3]);
const b = new Set([2, 3, 4]);

const inter = a.intersection(b);
console.log(inter.values()); // [2, 3]
```

#### difference(otherSet)

Diferença entre conjuntos.

```javascript
const a = new Set([1, 2, 3]);
const b = new Set([2, 3, 4]);

const diff = a.difference(b);
console.log(diff.values()); // [1]
```

---

## WeakMap

`WeakMap` é similar ao Map, mas permite garbage collection das chaves.

### Características

- Chaves devem ser objetos
- Não há `size` ou iteração
- Permite garbage collection

### Métodos

```javascript
const wm = new WeakMap();
const obj = { id: 1 };

wm.set(obj, "dados privados");
console.log(wm.get(obj)); // "dados privados"
console.log(wm.has(obj)); // true

wm.delete(obj);
console.log(wm.has(obj)); // false
```

---

## WeakSet

`WeakSet` é similar ao Set, mas permite garbage collection.

### Características

- Valores devem ser objetos
- Não há `size` ou iteração
- Permite garbage collection

### Métodos

```javascript
const ws = new WeakSet();
const obj = { id: 1 };

ws.add(obj);
console.log(ws.has(obj)); // true

ws.delete(obj);
console.log(ws.has(obj)); // false
```

---

## Exemplos Práticos

### Contar Frequência de Palavras

```javascript
function wordFrequency(text) {
  const words = text.toLowerCase().split(" ");
  const freq = new Map();

  words.forEach((word) => {
    const count = freq.get(word) || 0;
    freq.set(word, count + 1);
  });

  return freq;
}

const text = "o rato roeu a roupa do rei de roma";
const freq = wordFrequency(text);
console.log(freq.get("o")); // 1
console.log(freq.get("de")); // 1
console.log(freq.get("roupa")); // 1
```

### Remover Duplicatas

```javascript
function removeDuplicates(arr) {
  return Array.from(new Set(arr));
}

console.log(removeDuplicates([1, 2, 2, 3, 3, 3]));
// [1, 2, 3]
```

### Cache com Map

```javascript
const cache = new Map();

function fetchData(id) {
  if (cache.has(id)) {
    console.log("Do cache");
    return cache.get(id);
  }

  console.log("Buscando...");
  const data = { id: id, data: "..." };
  cache.set(id, data);
  return data;
}

fetchData(1); // Buscando...
fetchData(1); // Do cache
```

### Tracking de Objetos Visitados

```javascript
function detectCycle(obj, visited = new Set()) {
  if (typeof obj !== "object" || obj === null) {
    return false;
  }

  if (visited.has(obj)) {
    return true; // Ciclo detectado
  }

  visited.add(obj);

  const values = Object.values(obj);
  for (let i = 0; i < values.length; i++) {
    if (detectCycle(values[i], visited)) {
      return true;
    }
  }

  return false;
}
```

### Dados Privados com WeakMap

```javascript
const privateData = new WeakMap();

class User {
  constructor(name, password) {
    this.name = name;
    privateData.set(this, { password: password });
  }

  checkPassword(pwd) {
    return privateData.get(this).password === pwd;
  }
}

const user = new User("João", "secret123");
console.log(user.name); // "João"
console.log(user.checkPassword("secret123")); // true
// privateData não é acessível diretamente
```

## Map vs Object

| Feature     | Map                          | Object                   |
| ----------- | ---------------------------- | ------------------------ |
| Chaves      | Qualquer tipo                | String/Symbol            |
| Ordem       | Mantida                      | Não garantida            |
| Tamanho     | `.size`                      | `Object.keys().length`   |
| Iteração    | Direto                       | Via keys/values/entries  |
| Performance | Melhor para muitas operações | OK para poucos elementos |

## Set vs Array

| Feature    | Set         | Array   |
| ---------- | ----------- | ------- |
| Duplicatas | Não permite | Permite |
| Busca      | O(1)        | O(n)    |
| Ordem      | Mantida     | Mantida |
| Índices    | Não tem     | Tem     |

## Veja Também

- [Array](./array.md) - Listas ordenadas
- [Classes](./classes.md) - Definição de classes
