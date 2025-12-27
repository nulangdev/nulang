# Timers

Funções para executar código com atraso ou em intervalos regulares.

## setTimeout

Executa uma função após um atraso.

```javascript
setTimeout(callback, delay, ...args);
```

### Parâmetros

| Parâmetro  | Tipo     | Descrição                             |
| ---------- | -------- | ------------------------------------- |
| `callback` | Function | Função a executar                     |
| `delay`    | Number   | Atraso em milissegundos               |
| `...args`  | any      | Argumentos para o callback (opcional) |

### Retorno

Number - ID do timer (para cancelamento)

### Exemplo

```javascript
// Executa após 2 segundos
setTimeout(() => {
  console.log("2 segundos passaram!");
}, 2000);

// Com argumentos
setTimeout(
  (nome) => {
    console.log(`Olá, ${nome}!`);
  },
  1000,
  "João"
);
// Após 1 segundo: "Olá, João!"
```

## clearTimeout

Cancela um setTimeout pendente.

```javascript
clearTimeout(timeoutId);
```

### Exemplo

```javascript
const id = setTimeout(() => {
  console.log("Isso não será executado");
}, 5000);

// Cancela antes de executar
clearTimeout(id);
```

## setInterval

Executa uma função repetidamente em intervalos regulares.

```javascript
setInterval(callback, delay, ...args);
```

### Parâmetros

| Parâmetro  | Tipo     | Descrição                              |
| ---------- | -------- | -------------------------------------- |
| `callback` | Function | Função a executar                      |
| `delay`    | Number   | Intervalo em milissegundos (mín: 10ms) |
| `...args`  | any      | Argumentos para o callback (opcional)  |

### Retorno

Number - ID do intervalo (para cancelamento)

### Exemplo

```javascript
// Executa a cada 1 segundo
let contador = 0;
const id = setInterval(() => {
  contador++;
  console.log(`Tick ${contador}`);

  if (contador >= 5) {
    clearInterval(id);
    console.log("Parado!");
  }
}, 1000);
```

## clearInterval

Cancela um setInterval.

```javascript
clearInterval(intervalId);
```

### Exemplo

```javascript
const id = setInterval(() => {
  console.log("Repetindo...");
}, 1000);

// Para após 5 segundos
setTimeout(() => {
  clearInterval(id);
  console.log("Intervalo cancelado");
}, 5000);
```

## setImmediate

Executa uma função imediatamente após o código síncrono atual.

```javascript
setImmediate(callback, ...args);
```

### Exemplo

```javascript
console.log("1");

setImmediate(() => {
  console.log("3 - immediate");
});

console.log("2");

// Output:
// 1
// 2
// 3 - immediate
```

## clearImmediate

Cancela um setImmediate (para compatibilidade).

```javascript
clearImmediate(immediateId);
```

## nextTick

Similar a setImmediate, executa na próxima iteração.

```javascript
nextTick(callback, ...args);
```

### Exemplo

```javascript
console.log("Início");

nextTick(() => {
  console.log("Next tick");
});

console.log("Fim");

// Output:
// Início
// Fim
// Next tick
```

## sleep

Bloqueia a execução por um tempo (síncrono).

```javascript
sleep(ms);
```

### Exemplo

```javascript
console.log("Antes");
sleep(2000); // Pausa por 2 segundos
console.log("Depois de 2 segundos");
```

**Nota**: Use com cuidado, pois bloqueia a thread.

## Exemplos Práticos

### Debounce

```javascript
function debounce(fn, delay) {
  let timeoutId = null;

  return function (...args) {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      fn(...args);
      timeoutId = null;
    }, delay);
  };
}

// Uso
const buscar = debounce((termo) => {
  console.log("Buscando:", termo);
}, 300);

buscar("a");
buscar("ab");
buscar("abc"); // Só este executa
```

### Throttle

```javascript
function throttle(fn, delay) {
  let lastCall = 0;

  return function (...args) {
    const now = Date.now();
    if (now - lastCall >= delay) {
      fn(...args);
      lastCall = now;
    }
  };
}

// Uso
const scroll = throttle(() => {
  console.log("Scroll event");
}, 100);
```

### Retry com Backoff

```javascript
function retryWithBackoff(fn, maxRetries, initialDelay) {
  let retries = 0;

  function attempt() {
    try {
      return fn();
    } catch (error) {
      if (retries >= maxRetries) {
        throw error;
      }

      const delay = initialDelay * Math.pow(2, retries);
      retries++;

      console.log(`Retry ${retries} em ${delay}ms`);
      sleep(delay);
      return attempt();
    }
  }

  return attempt();
}
```

### Polling

```javascript
function poll(fn, interval, maxAttempts) {
  let attempts = 0;

  const intervalId = setInterval(() => {
    attempts++;
    const result = fn();

    if (result || attempts >= maxAttempts) {
      clearInterval(intervalId);
      if (result) {
        console.log("Sucesso após", attempts, "tentativas");
      } else {
        console.log("Timeout após", attempts, "tentativas");
      }
    }
  }, interval);
}

// Uso
poll(
  () => {
    // Retorna true quando pronto
    return Math.random() > 0.8;
  },
  1000,
  10
);
```

### Animação Simples

```javascript
function animate(duration, callback) {
  const start = Date.now();

  const frame = () => {
    const elapsed = Date.now() - start;
    const progress = Math.min(elapsed / duration, 1);

    callback(progress);

    if (progress < 1) {
      setImmediate(frame);
    }
  };

  frame();
}

// Uso
animate(1000, (progress) => {
  const value = Math.floor(progress * 100);
  console.log(`Progresso: ${value}%`);
});
```

### Sequência com Delays

```javascript
function sequence(tasks, delay) {
  let index = 0;

  function next() {
    if (index < tasks.length) {
      tasks[index]();
      index++;
      setTimeout(next, delay);
    }
  }

  next();
}

// Uso
sequence(
  [
    () => console.log("Tarefa 1"),
    () => console.log("Tarefa 2"),
    () => console.log("Tarefa 3"),
  ],
  1000
);
// Executa cada tarefa com 1 segundo de intervalo
```

### Countdown

```javascript
function countdown(seconds, callback) {
  let remaining = seconds;

  const id = setInterval(() => {
    callback(remaining);
    remaining--;

    if (remaining < 0) {
      clearInterval(id);
      console.log("Tempo esgotado!");
    }
  }, 1000);

  // Mostra imediatamente
  callback(remaining);
  remaining--;
}

// Uso
countdown(5, (s) => console.log(`${s} segundos restantes`));
```

### Rate Limiter

```javascript
function createRateLimiter(maxCalls, period) {
  const calls = [];

  return function (fn) {
    const now = Date.now();

    // Remove chamadas antigas
    while (calls.length > 0 && calls[0] < now - period) {
      calls.shift();
    }

    if (calls.length < maxCalls) {
      calls.push(now);
      fn();
    } else {
      console.log("Rate limit excedido");
    }
  };
}

// Uso: máximo 3 chamadas por 10 segundos
const limiter = createRateLimiter(3, 10000);

limiter(() => console.log("Chamada 1")); // OK
limiter(() => console.log("Chamada 2")); // OK
limiter(() => console.log("Chamada 3")); // OK
limiter(() => console.log("Chamada 4")); // Rate limit excedido
```

## Ordem de Execução

```javascript
console.log("1 - Síncrono");

setImmediate(() => {
  console.log("3 - Immediate");
});

setTimeout(() => {
  console.log("4 - Timeout 0ms");
}, 0);

nextTick(() => {
  console.log("2 - NextTick");
});

console.log("1b - Síncrono");

// Output:
// 1 - Síncrono
// 1b - Síncrono
// 2 - NextTick
// 3 - Immediate
// 4 - Timeout 0ms
```

## Veja Também

- [Promise](./promise.md) - Operações assíncronas
- [Events](./events.md) - Padrão de eventos
- [Event Loop](./event_loop.md) - Modelo de execução
