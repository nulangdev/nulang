# Reflect

O objeto `Reflect` fornece métodos estáticos para operações interceptáveis em JavaScript. Os métodos são os mesmos que os das traps de Proxy.

## Métodos

### `Reflect.get(target, property [, receiver])`

Retorna o valor de uma propriedade.

```javascript
const obj = { x: 10, y: 20 };

console.log(Reflect.get(obj, "x")); // 10
console.log(Reflect.get(obj, "y")); // 20
console.log(Reflect.get(obj, "z")); // undefined
```

#### Com Arrays

```javascript
const arr = [1, 2, 3];

console.log(Reflect.get(arr, 0)); // 1
console.log(Reflect.get(arr, "length")); // 3
```

#### Com Proxy

```javascript
const proxy = new Proxy(
  { name: "test" },
  {
    get: (t, p) => `[${t[p]}]`,
  }
);

console.log(Reflect.get(proxy, "name")); // "[test]"
```

### `Reflect.set(target, property, value [, receiver])`

Define o valor de uma propriedade. Retorna `true` se a operação foi bem-sucedida.

```javascript
const obj = {};

Reflect.set(obj, "x", 10);
console.log(obj.x); // 10

const result = Reflect.set(obj, "y", 20);
console.log(result); // true
```

#### Validação

```javascript
const obj = Object.freeze({ x: 1 });
const result = Reflect.set(obj, "x", 2);
console.log(result); // false (objeto congelado)
```

### `Reflect.has(target, property)`

Verifica se uma propriedade existe. Equivalente ao operador `in`.

```javascript
const obj = { x: 10 };

console.log(Reflect.has(obj, "x")); // true
console.log(Reflect.has(obj, "y")); // false
console.log(Reflect.has(obj, "toString")); // true (herdado)
```

### `Reflect.deleteProperty(target, property)`

Remove uma propriedade. Retorna `true` se a operação foi bem-sucedida.

```javascript
const obj = { x: 10, y: 20 };

console.log(Reflect.deleteProperty(obj, "x")); // true
console.log(obj); // { y: 20 }
```

### `Reflect.ownKeys(target)`

Retorna um array com todas as chaves (próprias) do objeto.

```javascript
const obj = { a: 1, b: 2, c: 3 };

console.log(Reflect.ownKeys(obj)); // ["a", "b", "c"]
```

### `Reflect.apply(target, thisArgument, argumentsList)`

Chama uma função com um `this` e argumentos específicos.

```javascript
function greet(greeting, name) {
  return `${greeting}, ${name}!`;
}

const result = Reflect.apply(greet, null, ["Hello", "World"]);
console.log(result); // "Hello, World!"
```

#### Com contexto `this`

```javascript
const obj = {
  name: "João",
  greet: function (greeting) {
    return `${greeting}, ${this.name}!`;
  },
};

const result = Reflect.apply(obj.greet, obj, ["Olá"]);
console.log(result); // "Olá, João!"
```

### `Reflect.construct(target, argumentsList [, newTarget])`

Equivalente ao operador `new`. Cria uma nova instância.

```javascript
class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }
}

const person = Reflect.construct(Person, ["Maria", 25]);
console.log(person.name); // "Maria"
console.log(person.age); // 25
```

#### Com função construtora

```javascript
function Point(x, y) {
  this.x = x;
  this.y = y;
}

const point = Reflect.construct(Point, [10, 20]);
console.log(point.x); // 10
console.log(point.y); // 20
```

### `Reflect.getPrototypeOf(target)`

Retorna o prototype do objeto.

```javascript
const obj = {};
console.log(Reflect.getPrototypeOf(obj)); // Object.prototype

class MyClass {}
const instance = new MyClass();
console.log(Reflect.getPrototypeOf(instance)); // MyClass.prototype
```

### `Reflect.setPrototypeOf(target, prototype)`

Define o prototype do objeto. Retorna `true` se bem-sucedido.

```javascript
const obj = {};
const proto = { greet: () => "Hello" };

Reflect.setPrototypeOf(obj, proto);
console.log(obj.greet()); // "Hello"
```

#### Removendo o prototype

```javascript
const obj = { x: 1 };
Reflect.setPrototypeOf(obj, null);
console.log(Object.getPrototypeOf(obj)); // null
```

### `Reflect.isExtensible(target)`

Verifica se novas propriedades podem ser adicionadas ao objeto.

```javascript
const obj = {};
console.log(Reflect.isExtensible(obj)); // true

Object.preventExtensions(obj);
console.log(Reflect.isExtensible(obj)); // false
```

### `Reflect.preventExtensions(target)`

Impede que novas propriedades sejam adicionadas ao objeto.

```javascript
const obj = { x: 1 };
Reflect.preventExtensions(obj);

obj.y = 2; // Não funciona (silenciosamente falha)
console.log(obj.y); // undefined
```

### `Reflect.getOwnPropertyDescriptor(target, property)`

Retorna o descriptor de uma propriedade.

```javascript
const obj = { x: 10 };

const descriptor = Reflect.getOwnPropertyDescriptor(obj, "x");
console.log(descriptor);
// {
//   value: 10,
//   writable: true,
//   enumerable: true,
//   configurable: true
// }
```

### `Reflect.defineProperty(target, property, descriptor)`

Define uma propriedade com um descriptor específico.

```javascript
const obj = {};

Reflect.defineProperty(obj, "x", {
  value: 42,
  writable: false,
  enumerable: true,
  configurable: false,
});

console.log(obj.x); // 42
obj.x = 100; // Não faz nada (não é writable)
console.log(obj.x); // 42
```

## Uso com Proxy

O Reflect é frequentemente usado dentro de handlers de Proxy para executar a operação padrão:

```javascript
const handler = {
  get: function (target, property, receiver) {
    console.log(`Acessando: ${property}`);
    // Usa Reflect para a operação padrão
    return Reflect.get(target, property, receiver);
  },

  set: function (target, property, value, receiver) {
    console.log(`Definindo ${property} = ${value}`);
    return Reflect.set(target, property, value, receiver);
  },
};

const obj = new Proxy({ x: 1 }, handler);
obj.x; // "Acessando: x"
obj.y = 2; // "Definindo y = 2"
```

## Comparação com Métodos de Object

| Reflect                              | Object                              |
| ------------------------------------ | ----------------------------------- |
| `Reflect.get(obj, prop)`             | `obj[prop]`                         |
| `Reflect.set(obj, prop, val)`        | `obj[prop] = val`                   |
| `Reflect.has(obj, prop)`             | `prop in obj`                       |
| `Reflect.deleteProperty(obj, prop)`  | `delete obj[prop]`                  |
| `Reflect.ownKeys(obj)`               | `Object.keys(obj)` + símbolos       |
| `Reflect.getPrototypeOf(obj)`        | `Object.getPrototypeOf(obj)`        |
| `Reflect.setPrototypeOf(obj, proto)` | `Object.setPrototypeOf(obj, proto)` |

## Vantagens do Reflect

1. **Retornos consistentes**: Reflect retorna `true`/`false` em vez de lançar erros
2. **Funções em vez de operadores**: Permite usar operadores como funções
3. **Compatibilidade com Proxy**: Mesmos métodos que as traps de Proxy
4. **Mais limpo que Object**: API mais previsível

### Exemplo de Tratamento de Erro

```javascript
// Com Object (pode lançar erro)
try {
  Object.defineProperty(frozenObj, "x", { value: 1 });
} catch (e) {
  console.log("Falhou");
}

// Com Reflect (retorna boolean)
const success = Reflect.defineProperty(frozenObj, "x", { value: 1 });
if (!success) {
  console.log("Falhou");
}
```

## Casos de Uso

### Função genérica de acesso a propriedades

```javascript
function getProperty(obj, prop, defaultValue) {
  if (Reflect.has(obj, prop)) {
    return Reflect.get(obj, prop);
  }
  return defaultValue;
}

const config = { debug: true };
console.log(getProperty(config, "debug", false)); // true
console.log(getProperty(config, "verbose", false)); // false
```

### Clonagem profunda

```javascript
function deepClone(obj) {
  if (obj === null || typeof obj !== "object") {
    return obj;
  }

  const clone = Reflect.construct(obj.constructor, []);

  for (const key of Reflect.ownKeys(obj)) {
    const descriptor = Reflect.getOwnPropertyDescriptor(obj, key);
    if (descriptor.value !== undefined) {
      descriptor.value = deepClone(descriptor.value);
    }
    Reflect.defineProperty(clone, key, descriptor);
  }

  return clone;
}
```

### Verificação de propriedade segura

```javascript
function safeGet(obj, path, defaultValue) {
  const keys = path.split(".");
  let current = obj;

  for (const key of keys) {
    if (current === null || current === undefined) {
      return defaultValue;
    }
    if (!Reflect.has(current, key)) {
      return defaultValue;
    }
    current = Reflect.get(current, key);
  }

  return current;
}

const data = { user: { name: "João" } };
console.log(safeGet(data, "user.name", "N/A")); // "João"
console.log(safeGet(data, "user.email", "N/A")); // "N/A"
console.log(safeGet(data, "company.name", "N/A")); // "N/A"
```

## Ver Também

- [Proxy](./proxy.md) - Objeto Proxy
- [Classes](./classes.md) - Declaração de classes
- [Map](./map.md) - Coleções Map e Set
