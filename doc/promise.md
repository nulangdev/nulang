# Promise

Promises representam operações assíncronas, compatível com JavaScript ES6+.

## Criação

### Usando o Construtor

```javascript
const promise = new Promise((resolve, reject) => {
  // Operação assíncrona
  if (success) {
    resolve(result);
  } else {
    reject(error);
  }
});
```

### Métodos Estáticos

#### Promise.resolve(value)

Cria uma Promise resolvida.

```javascript
const p = Promise.resolve(42);
p.then((value) => console.log(value)); // 42
```

#### Promise.reject(reason)

Cria uma Promise rejeitada.

```javascript
const p = Promise.reject(new Error("Falhou"));
p.catch((err) => console.log(err.message)); // "Falhou"
```

## Estados

Uma Promise pode estar em um de três estados:

| Estado      | Descrição                       |
| ----------- | ------------------------------- |
| `pending`   | Estado inicial, aguardando      |
| `fulfilled` | Operação completada com sucesso |
| `rejected`  | Operação falhou                 |

```javascript
const p = Promise.resolve(42);
console.log(p.state); // "fulfilled"
console.log(p.value); // 42
```

## Métodos de Instância

### then(onFulfilled, onRejected)

Adiciona callbacks para quando a Promise é resolvida ou rejeitada.

```javascript
promise.then(
  (value) => console.log("Sucesso:", value),
  (error) => console.log("Erro:", error)
);
```

**Retorno**: Nova Promise

### catch(onRejected)

Adiciona callback para quando a Promise é rejeitada.

```javascript
promise
  .then((value) => console.log(value))
  .catch((error) => console.log("Erro:", error));
```

**Retorno**: Nova Promise

### finally(onFinally)

Adiciona callback que sempre executa (sucesso ou falha).

```javascript
promise
  .then((value) => console.log(value))
  .catch((error) => console.log(error))
  .finally(() => console.log("Finalizado"));
```

**Retorno**: Nova Promise

## Métodos Estáticos de Composição

### Promise.all(iterable)

Aguarda todas as Promises serem resolvidas.

```javascript
const p1 = Promise.resolve(1);
const p2 = Promise.resolve(2);
const p3 = Promise.resolve(3);

Promise.all([p1, p2, p3]).then((values) => {
  console.log(values); // [1, 2, 3]
});
```

**Comportamento**:

- Resolve apenas quando todas as Promises resolvem
- Rejeita imediatamente se qualquer Promise rejeitar

```javascript
const p1 = Promise.resolve(1);
const p2 = Promise.reject(new Error("Falhou"));
const p3 = Promise.resolve(3);

Promise.all([p1, p2, p3]).catch((error) => {
  console.log(error.message); // "Falhou"
});
```

### Promise.race(iterable)

Retorna a primeira Promise a ser resolvida ou rejeitada.

```javascript
const p1 = new Promise((resolve) => {
  setTimeout(() => resolve("lento"), 1000);
});

const p2 = Promise.resolve("rápido");

Promise.race([p1, p2]).then((value) => {
  console.log(value); // "rápido"
});
```

### Promise.allSettled(iterable)

Aguarda todas as Promises serem finalizadas (resolvidas ou rejeitadas).

```javascript
const p1 = Promise.resolve(1);
const p2 = Promise.reject(new Error("Erro"));
const p3 = Promise.resolve(3);

Promise.allSettled([p1, p2, p3]).then((results) => {
  console.log(results);
  // [
  //   { status: "fulfilled", value: 1 },
  //   { status: "rejected", reason: Error("Erro") },
  //   { status: "fulfilled", value: 3 }
  // ]
});
```

## Encadeamento

Promises podem ser encadeadas, com cada `.then()` retornando uma nova Promise.

```javascript
Promise.resolve(1)
  .then((x) => x + 1)
  .then((x) => x * 2)
  .then((x) => console.log(x)); // 4
```

### Propagação de Erros

Erros propagam pela cadeia até encontrar um `.catch()`.

```javascript
Promise.resolve(1)
  .then((x) => {
    throw new Error("Erro no meio");
  })
  .then((x) => {
    // Não executa
    console.log("Não chega aqui");
  })
  .catch((error) => {
    console.log("Capturado:", error.message); // "Erro no meio"
  });
```

### Recuperação de Erros

```javascript
Promise.reject(new Error("Falhou"))
  .catch((error) => {
    console.log("Recuperando de:", error.message);
    return "valor de recuperação";
  })
  .then((value) => {
    console.log("Continuando com:", value);
  });
// Output:
// Recuperando de: Falhou
// Continuando com: valor de recuperação
```

## Async/Await Pattern

Embora Nulang não tenha async/await nativo, você pode simular com `.then()`:

```javascript
function fetchData() {
  return fetch("https://api.example.com/data")
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      return data;
    })
    .catch((error) => {
      console.error("Erro:", error);
    });
}
```

## Exemplos Práticos

### Timeout

```javascript
function timeout(ms) {
  return new Promise((resolve) => {
    setTimeout(() => resolve(), ms);
  });
}

timeout(2000).then(() => {
  console.log("2 segundos passaram");
});
```

### Retry com Delays

```javascript
function retryWithDelay(fn, retries, delay) {
  return fn().catch((error) => {
    if (retries > 0) {
      return timeout(delay).then(() => {
        return retryWithDelay(fn, retries - 1, delay);
      });
    }
    throw error;
  });
}
```

### Fetch com Fallback

```javascript
const primaryUrl = "https://primary.api.com/data";
const fallbackUrl = "https://fallback.api.com/data";

fetch(primaryUrl)
  .then((res) => {
    if (res.status !== 200) {
      throw new Error("Primary failed");
    }
    return res.json();
  })
  .catch(() => {
    console.log("Tentando fallback...");
    return fetch(fallbackUrl).then((res) => res.json());
  })
  .then((data) => {
    console.log("Dados:", data);
  });
```

### Sequência de Operações

```javascript
function processItems(items) {
  let result = Promise.resolve([]);

  items.forEach((item) => {
    result = result.then((acc) => {
      return processItem(item).then((processed) => {
        acc.push(processed);
        return acc;
      });
    });
  });

  return result;
}

function processItem(item) {
  return Promise.resolve(item * 2);
}

processItems([1, 2, 3]).then((results) => {
  console.log(results); // [2, 4, 6]
});
```

### Carregamento Paralelo com Limite

```javascript
function parallelLimit(tasks, limit) {
  const results = [];
  let running = 0;
  let index = 0;

  return new Promise((resolve) => {
    function next() {
      while (running < limit && index < tasks.length) {
        const currentIndex = index;
        const task = tasks[index];
        index++;
        running++;

        task().then((result) => {
          results[currentIndex] = result;
          running--;
          if (index >= tasks.length && running === 0) {
            resolve(results);
          } else {
            next();
          }
        });
      }
    }
    next();
  });
}
```

## Integração com Outros Módulos

### Com fetch

```javascript
fetch("https://api.example.com/users")
  .then((response) => response.json())
  .then((users) => {
    users.forEach((user) => console.log(user.name));
  })
  .catch((error) => {
    console.error("Erro ao buscar usuários:", error);
  });
```

### Com Blob/File

```javascript
const blob = new Blob(["Hello"]);
blob.text().then((text) => {
  console.log(text); // "Hello"
});
```

## Padrões Anti-Pattern

### ❌ Promise dentro de Promise (evitar)

```javascript
// Ruim
promise.then((value) => {
  return new Promise((resolve) => {
    resolve(value + 1);
  });
});

// Bom
promise.then((value) => value + 1);
```

### ❌ Não retornar em then (evitar)

```javascript
// Ruim - valor perdido
promise.then((value) => {
  doSomething(value); // sem return
});

// Bom
promise.then((value) => {
  return doSomething(value);
});
```

## Veja Também

- [Events](./events.md) - Padrão de eventos
- [Timers](./timers.md) - setTimeout, setInterval
- [HTTP](./http.md) - Requisições HTTP
