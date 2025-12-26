# Numbers

O tipo `Number` representa valores numéricos (inteiros e decimais), compatível com JavaScript.

## Criação

```javascript
// Inteiros
const int = 42;

// Decimais
const float = 3.14;

// Notação científica
const sci = 1.5e10; // 15000000000

// Negativos
const neg = -273.15;
```

## Objeto Global Math

O objeto `Math` fornece funções e constantes matemáticas.

### Constantes

| Constante | Valor      | Descrição           |
| --------- | ---------- | ------------------- |
| `Math.PI` | 3.14159... | Pi (π)              |
| `Math.E`  | 2.71828... | Número de Euler (e) |

```javascript
console.log(Math.PI); // 3.141592653589793
console.log(Math.E); // 2.718281828459045
```

### Métodos de Arredondamento

#### Math.round(x)

Arredonda para o inteiro mais próximo.

```javascript
console.log(Math.round(4.5)); // 5
console.log(Math.round(4.4)); // 4
console.log(Math.round(-4.5)); // -4
```

#### Math.floor(x)

Arredonda para baixo.

```javascript
console.log(Math.floor(4.9)); // 4
console.log(Math.floor(-4.1)); // -5
```

#### Math.ceil(x)

Arredonda para cima.

```javascript
console.log(Math.ceil(4.1)); // 5
console.log(Math.ceil(-4.9)); // -4
```

### Métodos Matemáticos

#### Math.abs(x)

Retorna o valor absoluto.

```javascript
console.log(Math.abs(-5)); // 5
console.log(Math.abs(5)); // 5
console.log(Math.abs(-3.14)); // 3.14
```

#### Math.pow(base, exponent)

Calcula a potência.

```javascript
console.log(Math.pow(2, 3)); // 8
console.log(Math.pow(3, 2)); // 9
console.log(Math.pow(4, 0.5)); // 2 (raiz quadrada)
```

#### Math.sqrt(x)

Calcula a raiz quadrada.

```javascript
console.log(Math.sqrt(16)); // 4
console.log(Math.sqrt(2)); // 1.4142135623730951
```

#### Math.max(...values)

Retorna o maior valor.

```javascript
console.log(Math.max(1, 5, 3)); // 5
console.log(Math.max(-1, -5, -3)); // -1
console.log(Math.max()); // -Infinity
```

#### Math.min(...values)

Retorna o menor valor.

```javascript
console.log(Math.min(1, 5, 3)); // 1
console.log(Math.min(-1, -5, -3)); // -5
console.log(Math.min()); // Infinity
```

#### Math.random()

Retorna um número aleatório entre 0 (inclusive) e 1 (exclusivo).

```javascript
console.log(Math.random()); // 0.123456789... (varia)

// Número entre 0 e 99
console.log(Math.floor(Math.random() * 100));

// Número entre min e max
function randomBetween(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}
console.log(randomBetween(1, 10)); // 1-10
```

## Funções de Conversão

### parseInt(string)

Converte string para inteiro.

```javascript
console.log(parseInt("42")); // 42
console.log(parseInt("42.9")); // 42
console.log(parseInt("   42")); // 42
console.log(parseInt("42abc")); // 42
console.log(parseInt("abc")); // NaN
```

### parseFloat(string)

Converte string para número decimal.

```javascript
console.log(parseFloat("3.14")); // 3.14
console.log(parseFloat("3.14abc")); // 3.14
console.log(parseFloat("   3.14")); // 3.14
```

### Number(value)

Converte qualquer valor para número.

```javascript
console.log(Number("42")); // 42
console.log(Number("3.14")); // 3.14
console.log(Number(true)); // 1
console.log(Number(false)); // 0
console.log(Number(null)); // 0
console.log(Number(undefined)); // NaN
console.log(Number("abc")); // NaN
```

## Valores Especiais

### NaN (Not a Number)

Resultado de operações matemáticas inválidas.

```javascript
console.log(0 / 0); // NaN
console.log(Number("abc")); // NaN
console.log(Math.sqrt(-1)); // NaN
```

#### isNaN(value)

Verifica se é NaN.

```javascript
console.log(isNaN(NaN)); // true
console.log(isNaN("abc")); // true
console.log(isNaN(123)); // false
console.log(isNaN("123")); // false
```

### Infinity

Representa infinito.

```javascript
console.log(1 / 0); // Infinity
console.log(-1 / 0); // -Infinity

console.log(Infinity + 1); // Infinity
console.log(Infinity - Infinity); // NaN
```

#### isFinite(value)

Verifica se é um número finito.

```javascript
console.log(isFinite(42)); // true
console.log(isFinite(Infinity)); // false
console.log(isFinite(NaN)); // false
```

## Operadores Aritméticos

| Operador | Descrição      | Exemplo               |
| -------- | -------------- | --------------------- |
| `+`      | Adição         | `5 + 3` → `8`         |
| `-`      | Subtração      | `5 - 3` → `2`         |
| `*`      | Multiplicação  | `5 * 3` → `15`        |
| `/`      | Divisão        | `10 / 3` → `3.333...` |
| `%`      | Módulo (resto) | `10 % 3` → `1`        |
| `**`     | Exponenciação  | `2 ** 3` → `8`        |

```javascript
console.log(10 + 5); // 15
console.log(10 - 5); // 5
console.log(10 * 5); // 50
console.log(10 / 5); // 2
console.log(10 % 3); // 1
console.log(2 ** 10); // 1024
```

## Operadores de Incremento/Decremento

```javascript
let n = 5;

n++; // Pós-incremento
console.log(n); // 6

n--; // Pós-decremento
console.log(n); // 5

++n; // Pré-incremento
console.log(n); // 6

--n; // Pré-decremento
console.log(n); // 5
```

## Operadores de Atribuição Composta

```javascript
let n = 10;

n += 5; // n = n + 5
console.log(n); // 15

n -= 3; // n = n - 3
console.log(n); // 12

n *= 2; // n = n * 2
console.log(n); // 24

n /= 4; // n = n / 4
console.log(n); // 6

n %= 4; // n = n % 4
console.log(n); // 2
```

## Exemplos Práticos

### Calcular Média

```javascript
function average(numbers) {
  const sum = numbers.reduce((a, b) => a + b, 0);
  return sum / numbers.length;
}

console.log(average([1, 2, 3, 4, 5])); // 3
```

### Clamp (Limitar Valor)

```javascript
function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

console.log(clamp(5, 0, 10)); // 5
console.log(clamp(-5, 0, 10)); // 0
console.log(clamp(15, 0, 10)); // 10
```

### Arredondar para N Casas Decimais

```javascript
function roundTo(value, decimals) {
  const factor = Math.pow(10, decimals);
  return Math.round(value * factor) / factor;
}

console.log(roundTo(3.14159, 2)); // 3.14
console.log(roundTo(3.14159, 4)); // 3.1416
```

### Verificar se é Par/Ímpar

```javascript
function isEven(n) {
  return n % 2 === 0;
}

function isOdd(n) {
  return n % 2 !== 0;
}

console.log(isEven(4)); // true
console.log(isOdd(5)); // true
```

### Distância Entre Dois Pontos

```javascript
function distance(x1, y1, x2, y2) {
  const dx = x2 - x1;
  const dy = y2 - y1;
  return Math.sqrt(dx * dx + dy * dy);
}

console.log(distance(0, 0, 3, 4)); // 5
```

### Fatorial

```javascript
function factorial(n) {
  if (n <= 1) return 1;
  let result = 1;
  for (let i = 2; i <= n; i++) {
    result *= i;
  }
  return result;
}

console.log(factorial(5)); // 120
```

## Veja Também

- [String](./string.md) - Manipulação de texto
- [Array](./array.md) - Listas ordenadas
