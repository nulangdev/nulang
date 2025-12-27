# Proxy

O objeto `Proxy` é usado para definir comportamento customizado para operações fundamentais em objetos (leitura de propriedades, atribuição, enumeração, invocação de funções, etc.).

## Construtor

```javascript
new Proxy(target, handler);
```

### Parâmetros

| Parâmetro | Tipo   | Descrição                             |
| --------- | ------ | ------------------------------------- |
| `target`  | Object | O objeto a ser proxeado               |
| `handler` | Object | Objeto com as traps (interceptadores) |

### Exemplo Básico

```javascript
const target = {
  message: "Hello, World!",
};

const handler = {
  get: function (target, property, receiver) {
    console.log(`Acessando propriedade: ${property}`);
    return target[property];
  },
};

const proxy = new Proxy(target, handler);
console.log(proxy.message);
// Output:
// Acessando propriedade: message
// Hello, World!
```

## Handler Traps

### `get(target, property, receiver)`

Intercepta leitura de propriedades.

```javascript
const handler = {
  get: function (target, property, receiver) {
    if (property in target) {
      return target[property];
    }
    return `Propriedade '${property}' não existe`;
  },
};

const obj = { x: 10 };
const proxy = new Proxy(obj, handler);

console.log(proxy.x); // 10
console.log(proxy.y); // "Propriedade 'y' não existe"
```

### `set(target, property, value, receiver)`

Intercepta atribuição de propriedades.

```javascript
const handler = {
  set: function (target, property, value, receiver) {
    if (typeof value !== "number") {
      throw new Error("Apenas números são permitidos");
    }
    target[property] = value;
    return true;
  },
};

const obj = {};
const proxy = new Proxy(obj, handler);

proxy.x = 10; // OK
proxy.y = "abc"; // Error: Apenas números são permitidos
```

### `has(target, property)`

Intercepta o operador `in`.

```javascript
const handler = {
  has: function (target, property) {
    if (property.startsWith("_")) {
      return false; // Esconde propriedades privadas
    }
    return property in target;
  },
};

const obj = { _private: 1, public: 2 };
const proxy = new Proxy(obj, handler);

console.log("public" in proxy); // true
console.log("_private" in proxy); // false
```

### `deleteProperty(target, property)`

Intercepta o operador `delete`.

```javascript
const handler = {
  deleteProperty: function (target, property) {
    if (property.startsWith("_")) {
      console.log("Não é possível deletar propriedades privadas");
      return false;
    }
    delete target[property];
    return true;
  },
};

const obj = { _id: 1, name: "test" };
const proxy = new Proxy(obj, handler);

delete proxy.name; // OK
delete proxy._id; // "Não é possível deletar propriedades privadas"
```

### `apply(target, thisArg, argumentsList)`

Intercepta chamadas de função.

```javascript
const handler = {
  apply: function (target, thisArg, argumentsList) {
    console.log(`Chamando com args: ${argumentsList}`);
    return target.apply(thisArg, argumentsList);
  },
};

function sum(a, b) {
  return a + b;
}

const proxy = new Proxy(sum, handler);
console.log(proxy(1, 2));
// Output:
// Chamando com args: 1,2
// 3
```

### `construct(target, argumentsList, newTarget)`

Intercepta o operador `new`.

```javascript
const handler = {
  construct: function (target, argumentsList, newTarget) {
    console.log(`Construindo com args: ${argumentsList}`);
    return new target(...argumentsList);
  },
};

class Person {
  constructor(name) {
    this.name = name;
  }
}

const ProxyPerson = new Proxy(Person, handler);
const p = new ProxyPerson("João");
// Output: Construindo com args: João
console.log(p.name); // "João"
```

### `getPrototypeOf(target)`

Intercepta `Object.getPrototypeOf`.

```javascript
const handler = {
  getPrototypeOf: function (target) {
    return { custom: true };
  },
};

const obj = {};
const proxy = new Proxy(obj, handler);
console.log(Object.getPrototypeOf(proxy)); // { custom: true }
```

### `setPrototypeOf(target, prototype)`

Intercepta `Object.setPrototypeOf`.

```javascript
const handler = {
  setPrototypeOf: function (target, prototype) {
    console.log("Tentativa de alterar prototype bloqueada");
    return false;
  },
};

const obj = {};
const proxy = new Proxy(obj, handler);
Object.setPrototypeOf(proxy, {}); // "Tentativa de alterar prototype bloqueada"
```

### `isExtensible(target)`

Intercepta `Object.isExtensible`.

```javascript
const handler = {
  isExtensible: function (target) {
    return true;
  },
};

const obj = {};
const proxy = new Proxy(obj, handler);
console.log(Object.isExtensible(proxy)); // true
```

### `preventExtensions(target)`

Intercepta `Object.preventExtensions`.

```javascript
const handler = {
  preventExtensions: function (target) {
    Object.preventExtensions(target);
    return true;
  },
};
```

### `getOwnPropertyDescriptor(target, property)`

Intercepta `Object.getOwnPropertyDescriptor`.

```javascript
const handler = {
  getOwnPropertyDescriptor: function (target, property) {
    return {
      value: target[property],
      writable: true,
      enumerable: true,
      configurable: true,
    };
  },
};
```

### `defineProperty(target, property, descriptor)`

Intercepta `Object.defineProperty`.

```javascript
const handler = {
  defineProperty: function (target, property, descriptor) {
    console.log(`Definindo ${property}`);
    Object.defineProperty(target, property, descriptor);
    return true;
  },
};
```

### `ownKeys(target)`

Intercepta `Object.keys`, `Object.getOwnPropertyNames`, etc.

```javascript
const handler = {
  ownKeys: function (target) {
    return Object.keys(target).filter((key) => !key.startsWith("_"));
  },
};

const obj = { _hidden: 1, visible: 2 };
const proxy = new Proxy(obj, handler);
console.log(Object.keys(proxy)); // ["visible"]
```

## Proxy.revocable()

Cria um Proxy revogável.

```javascript
const { proxy, revoke } = Proxy.revocable(target, handler);

// Use o proxy normalmente
console.log(proxy.foo);

// Revogue o proxy
revoke();

// Qualquer operação agora causa erro
proxy.foo; // Error: Cannot perform 'get' on a proxy that has been revoked
```

### Exemplo Completo

```javascript
const target = { name: "Test" };
const handler = {
  get: (t, p) => t[p],
};

const { proxy, revoke } = Proxy.revocable(target, handler);

console.log(proxy.name); // "Test"

revoke();

try {
  console.log(proxy.name);
} catch (e) {
  console.log(e.message); // "Cannot perform 'get' on a proxy that has been revoked"
}
```

## Casos de Uso

### Validação de Dados

```javascript
function createValidator(schema) {
  return {
    set: function (target, property, value) {
      if (property in schema) {
        const type = schema[property];
        if (typeof value !== type) {
          throw new Error(`${property} deve ser ${type}`);
        }
      }
      target[property] = value;
      return true;
    },
  };
}

const userSchema = {
  name: "string",
  age: "number",
  active: "boolean",
};

const user = new Proxy({}, createValidator(userSchema));

user.name = "João"; // OK
user.age = 30; // OK
user.age = "trinta"; // Error: age deve ser number
```

### Observador de Mudanças

```javascript
function observe(obj, callback) {
  return new Proxy(obj, {
    set: function (target, property, value) {
      const oldValue = target[property];
      target[property] = value;
      callback(property, oldValue, value);
      return true;
    },
  });
}

const data = observe({ count: 0 }, (prop, oldVal, newVal) => {
  console.log(`${prop} mudou de ${oldVal} para ${newVal}`);
});

data.count = 1; // "count mudou de 0 para 1"
data.count = 2; // "count mudou de 1 para 2"
```

### Cache/Memoização

```javascript
function memoize(fn) {
  const cache = new Map();

  return new Proxy(fn, {
    apply: function (target, thisArg, args) {
      const key = JSON.stringify(args);
      if (cache.has(key)) {
        return cache.get(key);
      }
      const result = target.apply(thisArg, args);
      cache.set(key, result);
      return result;
    },
  });
}

const expensiveFn = memoize(function (n) {
  console.log("Calculando...");
  return n * n;
});

expensiveFn(5); // "Calculando..." -> 25
expensiveFn(5); // 25 (do cache)
```

### Propriedades Negativas de Array

```javascript
const arr = [1, 2, 3, 4, 5];

const handler = {
  get: function (target, property) {
    const index = Number(property);
    if (!isNaN(index) && index < 0) {
      return target[target.length + index];
    }
    return target[property];
  },
};

const pArr = new Proxy(arr, handler);

console.log(pArr[-1]); // 5
console.log(pArr[-2]); // 4
```

### Objeto com Valores Default

```javascript
function withDefaults(target, defaults) {
  return new Proxy(target, {
    get: function (obj, property) {
      if (property in obj) {
        return obj[property];
      }
      return defaults[property];
    },
  });
}

const config = withDefaults({ debug: true }, { timeout: 1000, retries: 3 });

console.log(config.debug); // true (do objeto)
console.log(config.timeout); // 1000 (do default)
console.log(config.retries); // 3 (do default)
```

## Notas de Implementação

- Todas as traps são opcionais; quando não definidas, a operação padrão é executada
- Proxies podem ser aninhados para compor comportamentos
- `Proxy.revocable` é útil para controle de acesso temporário
- Proxies são transparentes: o código que usa o proxy não sabe que é um proxy

## Ver Também

- [Reflect](./reflect.md) - API para operações em objetos
- [Classes](./classes.md) - Declaração de classes
- [Getters/Setters](./getters_setters.md) - Propriedades computadas
