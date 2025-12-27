# Error

A classe `Error` representa erros em tempo de execução e é usada para tratamento de exceções em Nulang.

## Construtor

```javascript
new Error([message]);
```

### Parâmetros

| Parâmetro | Tipo   | Descrição                   |
| --------- | ------ | --------------------------- |
| `message` | string | Mensagem de erro (opcional) |

### Exemplo

```javascript
const error = new Error("Algo deu errado");
console.log(error.message); // "Algo deu errado"
console.log(error.name); // "Error"
```

## Propriedades de Instância

### `message`

A mensagem de erro.

```javascript
const error = new Error("Falha na operação");
console.log(error.message); // "Falha na operação"
```

### `name`

O nome do tipo de erro.

```javascript
const error = new Error("Falha");
console.log(error.name); // "Error"
```

### `stack`

Stack trace do erro.

```javascript
const error = new Error("Erro de teste");
console.log(error.stack);
// Error: Erro de teste
//     at <anonymous>
```

## Métodos de Instância

### `toString()`

Retorna uma representação em string do erro.

```javascript
const error = new Error("Mensagem de erro");
console.log(error.toString()); // "Error: Mensagem de erro"
```

## Tipos de Erro

Além do `Error` base, Nulang reconhece os seguintes tipos de erro:

| Tipo             | Descrição                          |
| ---------------- | ---------------------------------- |
| `Error`          | Erro genérico                      |
| `TypeError`      | Tipo de valor incorreto            |
| `ReferenceError` | Referência a variável não definida |
| `SyntaxError`    | Erro de sintaxe                    |
| `RangeError`     | Valor fora do intervalo permitido  |

### Criando Erros Específicos

```javascript
// Criar erro com nome customizado
const typeError = new Error("Esperava um número");
typeError.name = "TypeError";

console.log(typeError.toString()); // "TypeError: Esperava um número"
```

## Throw e Try/Catch

### throw

Lança uma exceção.

```javascript
function divide(a, b) {
  if (b === 0) {
    throw new Error("Divisão por zero");
  }
  return a / b;
}

// Também pode lançar strings
throw "Erro simples";

// Ou objetos
throw { code: 500, message: "Erro interno" };
```

### try/catch/finally

```javascript
try {
  const result = riskyOperation();
  console.log("Resultado:", result);
} catch (error) {
  console.error("Erro capturado:", error.message);
} finally {
  console.log("Sempre executa");
}
```

### Rethrowing

```javascript
try {
  doSomething();
} catch (error) {
  if (error.message.includes("crítico")) {
    throw error; // Re-lança o erro
  }
  // Trata outros erros
  console.log("Erro tratado:", error.message);
}
```

## Erros em Promises

### catch() em Promises

```javascript
fetch("https://api.example.com/data")
  .then((response) => response.json())
  .then((data) => console.log(data))
  .catch((error) => {
    console.error("Erro na requisição:", error.message);
  });
```

### Promise.reject()

```javascript
const failedPromise = Promise.reject(new Error("Promise rejeitada"));

failedPromise.catch((error) => {
  console.log(error.message); // "Promise rejeitada"
});
```

## Padrões de Error Handling

### Validação com Erros

```javascript
function validateUser(user) {
  if (!user) {
    throw new Error("Usuário é obrigatório");
  }
  if (!user.name) {
    throw new Error("Nome é obrigatório");
  }
  if (!user.email) {
    throw new Error("Email é obrigatório");
  }
  return true;
}

try {
  validateUser({ name: "João" });
} catch (e) {
  console.error(e.message); // "Email é obrigatório"
}
```

### Error Wrapper

```javascript
function wrapError(fn, context) {
  return function (...args) {
    try {
      return fn.apply(this, args);
    } catch (error) {
      const wrappedError = new Error(`${context}: ${error.message}`);
      wrappedError.originalError = error;
      throw wrappedError;
    }
  };
}

const safeParseJSON = wrapError(JSON.parse, "JSON Parse Error");

try {
  safeParseJSON("invalid json");
} catch (e) {
  console.log(e.message); // "JSON Parse Error: ..."
}
```

### Guard Clause Pattern

```javascript
function processData(data) {
  if (!data) {
    throw new Error("Data is required");
  }

  if (!Array.isArray(data)) {
    throw new Error("Data must be an array");
  }

  if (data.length === 0) {
    throw new Error("Data cannot be empty");
  }

  // Processamento normal
  return data.map((item) => item * 2);
}
```

## Funções Auxiliares

### isError (verificação interna)

O Nulang verifica internamente se um objeto é uma instância de Error:

```javascript
function handleValue(value) {
  // O sistema verifica se:
  // - É um ObjectMap com propriedade "name"
  // - O name é "Error", "TypeError", "ReferenceError", etc.
}
```

### getErrorMessage

Extrai a mensagem de erro de qualquer objeto:

```javascript
// Funciona com Error instances
const error = new Error("Mensagem");
// getErrorMessage(error) retorna "Mensagem"

// Funciona com objetos simples
// getErrorMessage({ message: "Texto" }) retorna "Texto"

// Fallback para outros tipos
// getErrorMessage("string") retorna "string"
```

## Stack Trace

O stack trace é gerado automaticamente quando um Error é criado:

```javascript
function outer() {
  function inner() {
    throw new Error("Erro interno");
  }
  inner();
}

try {
  outer();
} catch (e) {
  console.log(e.stack);
}
// Output:
// Error: Erro interno
//     at inner
//     at outer
//     at <anonymous>
```

## Boas Práticas

1. **Use mensagens descritivas**

```javascript
// ❌ Ruim
throw new Error("Erro");

// ✅ Bom
throw new Error("Falha ao conectar com o servidor: timeout após 30s");
```

2. **Capture erros específicos**

```javascript
try {
  // operação
} catch (error) {
  if (error.name === "TypeError") {
    // Trata tipo específico
  } else {
    throw error; // Re-lança outros
  }
}
```

3. **Sempre limpe recursos no finally**

```javascript
let file;
try {
  file = openFile("data.txt");
  processFile(file);
} finally {
  if (file) {
    closeFile(file);
  }
}
```

## Ver Também

- [Promise](./promise.md) - Promise.reject e .catch()
- [Classes](./classes.md) - Criando classes de erro customizadas
