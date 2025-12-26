# Decorators

Decorators são uma funcionalidade que permite modificar classes, métodos e propriedades de forma declarativa, seguindo a proposta do ECMAScript Decorators.

## Sintaxe Básica

```javascript
@decoratorName
class MinhaClasse {
  @decoratorMethod
  meuMetodo() {}

  @decoratorProperty
  minhaPropriedade = value;
}
```

## Decorators de Classe

Decorators de classe recebem a classe como argumento e podem retornar uma nova classe modificada.

```javascript
function logged(target) {
  console.log(`Classe ${target.name} foi criada`);
  return target;
}

@logged
class MinhaClasse {
  constructor() {
    console.log("Instância criada");
  }
}

const obj = new MinhaClasse();
// Output:
// Classe MinhaClasse foi criada
// Instância criada
```

### Decorator Factory (com argumentos)

```javascript
function setVersion(version) {
  return function (target) {
    target.version = version;
    return target;
  };
}

@setVersion("1.0.0")
class App {
  static version;
}

console.log(App.version); // "1.0.0"
```

## Decorators de Método

Decorators de método recebem três argumentos:

- `target`: O objeto alvo (classe ou protótipo)
- `propertyKey`: O nome do método
- `descriptor`: O property descriptor

```javascript
function log(target, propertyKey, descriptor) {
  const originalMethod = descriptor.value;

  descriptor.value = function (...args) {
    console.log(`Chamando ${propertyKey} com args:`, args);
    const result = originalMethod.apply(this, args);
    console.log(`Resultado:`, result);
    return result;
  };

  return descriptor;
}

class Calculator {
  @log
  add(a, b) {
    return a + b;
  }
}

const calc = new Calculator();
calc.add(2, 3);
// Output:
// Chamando add com args: [2, 3]
// Resultado: 5
```

### Timing Decorator

```javascript
function timing(target, propertyKey, descriptor) {
  const originalMethod = descriptor.value;

  descriptor.value = function (...args) {
    const start = Date.now();
    const result = originalMethod.apply(this, args);
    const end = Date.now();
    console.log(`${propertyKey} executou em ${end - start}ms`);
    return result;
  };

  return descriptor;
}

class Worker {
  @timing
  processData(data) {
    // Processamento pesado
    let sum = 0;
    for (let i = 0; i < 1000000; i++) {
      sum += i;
    }
    return sum;
  }
}
```

## Decorators de Propriedade

Decorators de propriedade permitem modificar ou validar valores de propriedades.

```javascript
function readonly(target, propertyKey, descriptor) {
  descriptor.writable = false;
  return descriptor;
}

function validate(min, max) {
  return function (target, propertyKey, descriptor) {
    let value = descriptor.value;

    const getter = function () {
      return value;
    };

    const setter = function (newVal) {
      if (newVal < min || newVal > max) {
        throw new Error(`Valor deve estar entre ${min} e ${max}`);
      }
      value = newVal;
    };

    descriptor.get = getter;
    descriptor.set = setter;
    delete descriptor.value;
    delete descriptor.writable;

    return descriptor;
  };
}

class Config {
  @readonly
  version = "1.0.0";

  @validate(0, 100)
  volume = 50;
}
```

## Composição de Decorators

Múltiplos decorators podem ser aplicados. Eles são avaliados de cima para baixo, mas aplicados de baixo para cima.

```javascript
function first() {
  console.log("first(): avaliado");
  return function (target, propertyKey, descriptor) {
    console.log("first(): aplicado");
    return descriptor;
  };
}

function second() {
  console.log("second(): avaliado");
  return function (target, propertyKey, descriptor) {
    console.log("second(): aplicado");
    return descriptor;
  };
}

class Example {
  @first()
  @second()
  method() {}
}

// Output:
// first(): avaliado
// second(): avaliado
// second(): aplicado
// first(): aplicado
```

## Casos de Uso Comuns

### Memoização

```javascript
function memoize(target, propertyKey, descriptor) {
  const originalMethod = descriptor.value;
  const cache = new Map();

  descriptor.value = function (...args) {
    const key = JSON.stringify(args);
    if (cache.has(key)) {
      return cache.get(key);
    }
    const result = originalMethod.apply(this, args);
    cache.set(key, result);
    return result;
  };

  return descriptor;
}

class Math {
  @memoize
  fibonacci(n) {
    if (n <= 1) return n;
    return this.fibonacci(n - 1) + this.fibonacci(n - 2);
  }
}
```

### Deprecation Warning

```javascript
function deprecated(message) {
  return function (target, propertyKey, descriptor) {
    const originalMethod = descriptor.value;

    descriptor.value = function (...args) {
      console.warn(`DEPRECATED: ${propertyKey} - ${message}`);
      return originalMethod.apply(this, args);
    };

    return descriptor;
  };
}

class API {
  @deprecated("Use newMethod() instead")
  oldMethod() {
    // ...
  }
}
```

### Bind Automático

```javascript
function autobind(target, propertyKey, descriptor) {
  const originalMethod = descriptor.value;

  return {
    configurable: true,
    get() {
      const bound = originalMethod.bind(this);
      Object.defineProperty(this, propertyKey, {
        value: bound,
        configurable: true,
        writable: true,
      });
      return bound;
    },
  };
}

class Button {
  constructor() {
    this.text = "Click me";
  }

  @autobind
  handleClick() {
    console.log(this.text);
  }
}

const button = new Button();
const handler = button.handleClick;
handler(); // "Click me" (this está correto)
```

## Property Descriptor

O descriptor passado para decorators de método/propriedade contém:

| Propriedade    | Descrição                             |
| -------------- | ------------------------------------- |
| `value`        | O valor da propriedade (para métodos) |
| `writable`     | Se o valor pode ser alterado          |
| `enumerable`   | Se aparece em for...in                |
| `configurable` | Se pode ser deletado ou reconfigurado |
| `get`          | Função getter (para accessors)        |
| `set`          | Função setter (para accessors)        |

## Notas de Implementação

- Decorators são aplicados em tempo de definição da classe
- A ordem de aplicação é bottom-up (de baixo para cima)
- Decorators de classe são aplicados após os decorators de membros
- Decorators podem ser funções ou factories que retornam funções

## Ver Também

- [Classes](./classes.md) - Declaração de classes
- [Getters/Setters](./getters_setters.md) - Propriedades computadas
- [Reflect](./reflect.md) - API Reflect
