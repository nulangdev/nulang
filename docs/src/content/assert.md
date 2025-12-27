# Módulo Assert

O módulo `assert` fornece funções de asserção para testes e validações.

## Importação

```javascript
const assert = require("assert");
```

## Funções de Asserção

### `assert(value, message?)`

Verifica se um valor é truthy.

```javascript
assert(true); // OK
assert(1); // OK
assert("string"); // OK
assert(false, "Valor é falso"); // Erro: Valor é falso
```

### `assert.ok(value, message?)`

Alias para `assert()`.

```javascript
assert.ok(true);
assert.ok(1 === 1);
assert.ok([], "Array vazio é truthy");
```

### `assert.equal(actual, expected, message?)`

Verifica igualdade com coerção de tipo (==).

```javascript
assert.equal(1, 1); // OK
assert.equal(1, "1"); // OK (coerção)
assert.equal(true, 1); // OK (coerção)
```

### `assert.notEqual(actual, expected, message?)`

Verifica desigualdade com coerção.

```javascript
assert.notEqual(1, 2); // OK
assert.notEqual("a", "b"); // OK
```

### `assert.strictEqual(actual, expected, message?)`

Verifica igualdade estrita (===).

```javascript
assert.strictEqual(1, 1); // OK
assert.strictEqual("a", "a"); // OK
assert.strictEqual(1, "1"); // Erro: 1 !== '1'
```

### `assert.notStrictEqual(actual, expected, message?)`

Verifica desigualdade estrita.

```javascript
assert.notStrictEqual(1, "1"); // OK (tipos diferentes)
assert.notStrictEqual(1, 2); // OK
```

### `assert.deepEqual(actual, expected, message?)`

Verifica igualdade profunda de objetos.

```javascript
assert.deepEqual({ a: 1 }, { a: 1 }); // OK
assert.deepEqual([1, 2, 3], [1, 2, 3]); // OK
assert.deepEqual({ a: { b: 1 } }, { a: { b: 1 } }); // OK
```

### `assert.notDeepEqual(actual, expected, message?)`

Verifica desigualdade profunda.

```javascript
assert.notDeepEqual({ a: 1 }, { a: 2 }); // OK
assert.notDeepEqual([1, 2], [1, 2, 3]); // OK
```

### `assert.deepStrictEqual(actual, expected, message?)`

Igualdade profunda estrita.

```javascript
assert.deepStrictEqual({ a: 1 }, { a: 1 }); // OK
```

### `assert.throws(fn, error?, message?)`

Verifica se uma função lança erro.

```javascript
assert.throws(() => {
  throw new Error("test");
});

assert.throws(
  () => {
    throw new Error("specific");
  },
  Error,
  "Deveria lançar Error"
);
```

### `assert.doesNotThrow(fn, message?)`

Verifica se uma função não lança erro.

```javascript
assert.doesNotThrow(() => {
  return 1 + 1;
});
```

### `assert.fail(message?)`

Força uma falha.

```javascript
if (condition) {
  assert.fail("Condição não deveria ser verdadeira");
}
```

### `assert.ifError(value)`

Lança se o valor não for null/undefined.

```javascript
assert.ifError(null); // OK
assert.ifError(undefined); // OK
assert.ifError(new Error()); // Erro lançado
```

## Exemplo: Testes Unitários

```javascript
const assert = require("assert");

function soma(a, b) {
  return a + b;
}

function divide(a, b) {
  if (b === 0) throw new Error("Divisão por zero");
  return a / b;
}

// Testes
console.log("Testando soma...");
assert.strictEqual(soma(2, 3), 5);
assert.strictEqual(soma(-1, 1), 0);
assert.strictEqual(soma(0, 0), 0);

console.log("Testando divisão...");
assert.strictEqual(divide(10, 2), 5);
assert.throws(() => divide(1, 0));

console.log("Todos os testes passaram!");
```

## Strict Mode

```javascript
const assert = require("assert").strict;

// Todas as asserções usam comparação estrita
assert.equal(1, 1); // Usa ===
assert.deepEqual({}, {}); // Usa comparação estrita
```

## Compatibilidade

| Funcionalidade | Nulang | Node.js |
| -------------- | ------ | ------- |
| `assert()`     | ✅     | ✅      |
| `ok`           | ✅     | ✅      |
| `equal`        | ✅     | ✅      |
| `notEqual`     | ✅     | ✅      |
| `strictEqual`  | ✅     | ✅      |
| `deepEqual`    | ✅     | ✅      |
| `throws`       | ✅     | ✅      |
| `doesNotThrow` | ✅     | ✅      |
| `fail`         | ✅     | ✅      |
| `ifError`      | ✅     | ✅      |
| `rejects`      | ❌     | ✅      |
| `match`        | ❌     | ✅      |
