# Módulo Util

O módulo `util` fornece funções utilitárias para debugging, formatação e outras tarefas.

## Importação

```javascript
const util = require("util");
```

## Funções de Inspeção

### `util.inspect(object, options?)`

Retorna uma representação string de um objeto para debugging.

```javascript
const obj = { name: "João", age: 30, nested: { a: 1, b: 2 } };

console.log(util.inspect(obj));
// { name: 'João', age: 30, nested: { a: 1, b: 2 } }

// Com opções
console.log(
  util.inspect(obj, {
    depth: 1, // Profundidade máxima
    colors: true, // Cores no terminal
    showHidden: true, // Mostrar props não-enumeráveis
  })
);
```

#### Opções

| Opção        | Tipo    | Descrição                            |
| ------------ | ------- | ------------------------------------ |
| `depth`      | Number  | Profundidade de recursão (padrão: 2) |
| `colors`     | Boolean | Usar cores ANSI                      |
| `showHidden` | Boolean | Mostrar propriedades ocultas         |

## Funções de Formatação

### `util.format(format, ...args)`

Formata uma string com placeholders.

```javascript
util.format("Olá, %s!", "Mundo"); // "Olá, Mundo!"
util.format("Valor: %d", 42); // "Valor: 42"
util.format("JSON: %j", { a: 1 }); // "JSON: {\"a\":1}"
util.format("Objeto: %o", { x: 1 }); // "Objeto: { x: 1 }"
util.format("100%%"); // "100%"
```

#### Placeholders

| Placeholder | Descrição |
| ----------- | --------- |
| `%s`        | String    |
| `%d`        | Número    |
| `%i`        | Inteiro   |
| `%f`        | Float     |
| `%j`        | JSON      |
| `%o`        | Objeto    |
| `%%`        | Literal % |

### `util.formatWithOptions(options, format, ...args)`

Format com opções de inspect.

```javascript
util.formatWithOptions({ colors: true }, "Obj: %o", { a: 1 });
```

## Funções de Tipo

### `util.types.isDate(value)`

```javascript
util.types.isDate(new Date()); // true
util.types.isDate("2024-01-01"); // false
```

### `util.types.isRegExp(value)`

```javascript
util.types.isRegExp(/abc/); // true
util.types.isRegExp("abc"); // false
```

### `util.types.isArray(value)`

```javascript
util.types.isArray([1, 2, 3]); // true
util.types.isArray("abc"); // false
```

### `util.types.isFunction(value)`

```javascript
util.types.isFunction(() => {}); // true
util.types.isFunction({}); // false
```

### `util.types.isObject(value)`

```javascript
util.types.isObject({}); // true
util.types.isObject(null); // false
```

### `util.types.isString(value)`

```javascript
util.types.isString("hello"); // true
util.types.isString(123); // false
```

### `util.types.isNumber(value)`

```javascript
util.types.isNumber(42); // true
util.types.isNumber("42"); // false
```

### `util.types.isBoolean(value)`

```javascript
util.types.isBoolean(true); // true
util.types.isBoolean(1); // false
```

### `util.types.isNull(value)`

```javascript
util.types.isNull(null); // true
util.types.isNull(undefined); // false
```

### `util.types.isUndefined(value)`

```javascript
util.types.isUndefined(undefined); // true
util.types.isUndefined(null); // false
```

### `util.types.isNullOrUndefined(value)`

```javascript
util.types.isNullOrUndefined(null); // true
util.types.isNullOrUndefined(undefined); // true
util.types.isNullOrUndefined(0); // false
```

## Funções de Debugging

### `util.debuglog(section)`

Cria uma função de log condicional.

```javascript
const debug = util.debuglog("myapp");

debug("Mensagem de debug");
// Só aparece se NODE_DEBUG=myapp
```

### `util.deprecate(fn, message)`

Marca uma função como deprecated.

```javascript
const oldFunc = util.deprecate(() => {
  return "resultado";
}, "oldFunc está obsoleta, use newFunc");

oldFunc(); // Exibe warning
```

## Exemplos Práticos

### Logger Formatado

```javascript
const util = require("util");

function log(level, message, ...args) {
  const timestamp = new Date().toISOString();
  const formatted = util.format(message, ...args);
  console.log(`[${timestamp}] [${level}] ${formatted}`);
}

log("INFO", "Usuário %s logou", "João");
log("ERROR", "Falha: %j", { code: 500 });
```

### Debug de Objetos

```javascript
function debugObject(obj) {
  console.log("=== DEBUG ===");
  console.log(
    util.inspect(obj, {
      depth: null,
      colors: true,
    })
  );
  console.log("=============");
}

debugObject({
  users: [
    { name: "A", data: { x: 1 } },
    { name: "B", data: { x: 2 } },
  ],
});
```

## Compatibilidade

| Funcionalidade | Nulang    | Node.js |
| -------------- | --------- | ------- |
| `inspect`      | ✅        | ✅      |
| `format`       | ✅        | ✅      |
| `types.*`      | ✅        | ✅      |
| `debuglog`     | ⚠️ Básico | ✅      |
| `deprecate`    | ⚠️ Básico | ✅      |
| `promisify`    | ❌        | ✅      |
| `callbackify`  | ❌        | ✅      |
| `TextEncoder`  | ❌        | ✅      |
