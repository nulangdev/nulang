# Módulo VM

O módulo `vm` permite compilar e executar código JavaScript em contextos isolados.

## Importação

```javascript
const vm = require("vm");
```

## Funções

### `vm.runInThisContext(code, options?)`

Executa código no contexto global atual.

```javascript
const result = vm.runInThisContext("2 + 2");
console.log(result); // 4

vm.runInThisContext(`
  const x = 10;
  console.log(x * 2);
`);
```

### `vm.runInNewContext(code, sandbox?, options?)`

Executa código em um novo contexto com sandbox.

```javascript
const sandbox = { x: 10, y: 20 };

const result = vm.runInNewContext("x + y", sandbox);
console.log(result); // 30

// Modificando o sandbox
vm.runInNewContext("z = x * y", sandbox);
console.log(sandbox.z); // 200
```

### `vm.runInContext(code, context, options?)`

Executa código em um contexto existente.

```javascript
const context = vm.createContext({ value: 100 });

vm.runInContext("result = value * 2", context);
console.log(context.result); // 200
```

### `vm.createContext(sandbox?)`

Cria um contexto para execução de código.

```javascript
const context = vm.createContext({
  console: console,
  Math: Math,
  myVar: 42,
});

vm.runInContext("console.log(myVar)", context);
```

### `vm.isContext(object)`

Verifica se um objeto é um contexto.

```javascript
const context = vm.createContext({});
console.log(vm.isContext(context)); // true
console.log(vm.isContext({})); // false
```

### `vm.compileFunction(code, params?, options?)`

Compila código como uma função.

```javascript
const fn = vm.compileFunction("return a + b", ["a", "b"]);
console.log(fn(2, 3)); // 5

const greet = vm.compileFunction('return "Hello, " + name', ["name"]);
console.log(greet("World")); // "Hello, World"
```

## Classe Script

### `new vm.Script(code, options?)`

Compila código para execução posterior.

```javascript
const script = new vm.Script("x + y");

// Executar múltiplas vezes
const context1 = vm.createContext({ x: 1, y: 2 });
console.log(script.runInContext(context1)); // 3

const context2 = vm.createContext({ x: 10, y: 20 });
console.log(script.runInContext(context2)); // 30
```

### Métodos do Script

```javascript
const script = new vm.Script("result = value * 2");

// runInThisContext
script.runInThisContext();

// runInNewContext
script.runInNewContext({ value: 10 });

// runInContext
const ctx = vm.createContext({ value: 5 });
script.runInContext(ctx);
```

## Exemplos

### Avaliador de Expressões

```javascript
function evaluate(expression, variables) {
  return vm.runInNewContext(expression, variables);
}

console.log(evaluate("x * 2 + y", { x: 5, y: 3 })); // 13
console.log(evaluate("Math.sqrt(a)", { a: 16, Math })); // 4
```

### Sandbox Seguro

```javascript
const sandbox = vm.createContext({
  console: { log: console.log },
  result: null,
});

const code = `
  result = 'Executado no sandbox';
  console.log('Hello from sandbox');
`;

vm.runInContext(code, sandbox);
console.log(sandbox.result);
```

### Template Engine Simples

```javascript
function render(template, data) {
  const code = `\`${template}\``;
  return vm.runInNewContext(code, data);
}

const html = render("Olá, ${name}! Você tem ${age} anos.", {
  name: "João",
  age: 30,
});
console.log(html);
```

## Compatibilidade

| Funcionalidade     | Nulang | Node.js |
| ------------------ | ------ | ------- |
| `runInThisContext` | ✅     | ✅      |
| `runInNewContext`  | ✅     | ✅      |
| `runInContext`     | ✅     | ✅      |
| `createContext`    | ✅     | ✅      |
| `isContext`        | ✅     | ✅      |
| `compileFunction`  | ✅     | ✅      |
| `Script` class     | ✅     | ✅      |
| `Module` class     | ❌     | ✅      |
| `timeout` option   | ❌     | ✅      |
