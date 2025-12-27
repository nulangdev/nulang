# Date

O objeto `Date` representa datas e horários, compatível com JavaScript.

## Criação

```javascript
// Data/hora atual
const now = new Date();

// De timestamp (milissegundos desde 1970-01-01)
const fromTimestamp = new Date(1609459200000);

// De string
const fromString = new Date("2021-01-01T00:00:00Z");
const fromString2 = new Date("Jan 1, 2021");

// Com componentes
const custom = new Date(2021, 0, 15, 10, 30, 0); // Ano, Mês (0-11), Dia, Hora, Min, Seg
```

**Nota**: Meses são indexados em 0 (Janeiro = 0, Dezembro = 11).

## Métodos Estáticos

### Date.now()

Retorna o timestamp atual em milissegundos.

```javascript
const timestamp = Date.now();
console.log(timestamp); // 1703548800000
```

### Date.parse(dateString)

Parseia uma string de data e retorna o timestamp.

```javascript
const ts = Date.parse("2021-01-01");
console.log(ts); // 1609459200000
```

### Date.UTC(year, month, ...)

Retorna timestamp para data em UTC.

```javascript
const ts = Date.UTC(2021, 0, 1, 0, 0, 0);
console.log(ts); // 1609459200000
```

## Métodos Getters

### Componentes de Data

| Método                | Retorno | Descrição                  |
| --------------------- | ------- | -------------------------- |
| `getFullYear()`       | Number  | Ano (4 dígitos)            |
| `getMonth()`          | Number  | Mês (0-11)                 |
| `getDate()`           | Number  | Dia do mês (1-31)          |
| `getDay()`            | Number  | Dia da semana (0-6, Dom=0) |
| `getHours()`          | Number  | Hora (0-23)                |
| `getMinutes()`        | Number  | Minutos (0-59)             |
| `getSeconds()`        | Number  | Segundos (0-59)            |
| `getMilliseconds()`   | Number  | Milissegundos (0-999)      |
| `getTime()`           | Number  | Timestamp em ms            |
| `getTimezoneOffset()` | Number  | Offset do fuso em minutos  |

```javascript
const date = new Date(2021, 6, 15, 14, 30, 45, 123);

console.log(date.getFullYear()); // 2021
console.log(date.getMonth()); // 6 (Julho)
console.log(date.getDate()); // 15
console.log(date.getDay()); // 4 (Quinta-feira)
console.log(date.getHours()); // 14
console.log(date.getMinutes()); // 30
console.log(date.getSeconds()); // 45
console.log(date.getMilliseconds()); // 123
console.log(date.getTime()); // timestamp
```

## Métodos Setters

### Modificando Componentes

| Método                               | Descrição         |
| ------------------------------------ | ----------------- |
| `setTime(ms)`                        | Define timestamp  |
| `setFullYear(year, [month], [day])`  | Define ano        |
| `setMonth(month, [day])`             | Define mês        |
| `setDate(day)`                       | Define dia do mês |
| `setHours(hour, [min], [sec], [ms])` | Define hora       |

```javascript
const date = new Date();

date.setFullYear(2025);
date.setMonth(11); // Dezembro
date.setDate(25);
date.setHours(10, 30, 0);

console.log(date.toString());
```

**Retorno**: Timestamp após a modificação

## Métodos de Formatação

### toString()

String completa com data e hora.

```javascript
const date = new Date();
console.log(date.toString());
// "Thu Dec 26 2024 00:01:27 GMT-0300 (-03)"
```

### toISOString()

String no formato ISO 8601 (UTC).

```javascript
const date = new Date();
console.log(date.toISOString());
// "2024-12-26T03:01:27Z"
```

### toDateString()

Apenas a parte da data.

```javascript
const date = new Date();
console.log(date.toDateString());
// "Thu Dec 26 2024"
```

### toTimeString()

Apenas a parte do horário.

```javascript
const date = new Date();
console.log(date.toTimeString());
// "00:01:27 -03"
```

### toLocaleDateString()

Data formatada para o locale.

```javascript
const date = new Date();
console.log(date.toLocaleDateString());
// "12/26/2024"
```

### toLocaleTimeString()

Hora formatada para o locale.

```javascript
const date = new Date();
console.log(date.toLocaleTimeString());
// "12:01:27 AM"
```

### toLocaleString()

Data e hora formatadas para o locale.

```javascript
const date = new Date();
console.log(date.toLocaleString());
// "12/26/2024, 12:01:27 AM"
```

### toUTCString()

String em formato UTC.

```javascript
const date = new Date();
console.log(date.toUTCString());
// "Thu, 26 Dec 2024 03:01:27 GMT"
```

### toJSON()

String para serialização JSON (igual a toISOString).

```javascript
const date = new Date();
console.log(date.toJSON());
// "2024-12-26T03:01:27Z"
```

### valueOf()

Retorna o timestamp (primitivo).

```javascript
const date = new Date();
console.log(date.valueOf());
// 1703556087000
```

## Exemplos Práticos

### Calcular Diferença Entre Datas

```javascript
const date1 = new Date("2024-01-01");
const date2 = new Date("2024-12-31");

const diffMs = date2.getTime() - date1.getTime();
const diffDays = diffMs / (1000 * 60 * 60 * 24);

console.log(`Diferença: ${diffDays} dias`); // "Diferença: 365 dias"
```

### Adicionar Dias

```javascript
function addDays(date, days) {
  const result = new Date(date.getTime());
  result.setDate(result.getDate() + days);
  return result;
}

const today = new Date();
const nextWeek = addDays(today, 7);
console.log(nextWeek.toDateString());
```

### Verificar Se É Hoje

```javascript
function isToday(date) {
  const today = new Date();
  return (
    date.getFullYear() === today.getFullYear() &&
    date.getMonth() === today.getMonth() &&
    date.getDate() === today.getDate()
  );
}

console.log(isToday(new Date())); // true
```

### Formatar Data Personalizado

```javascript
function formatDate(date) {
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}/${month}/${year}`;
}

// Pad helper
String.prototype.padStart = function (length, char) {
  let result = this;
  while (result.length < length) {
    result = char + result;
  }
  return result;
};

console.log(formatDate(new Date(2024, 0, 5))); // "05/01/2024"
```

### Obter Primeiro/Último Dia do Mês

```javascript
function getFirstDayOfMonth(date) {
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function getLastDayOfMonth(date) {
  return new Date(date.getFullYear(), date.getMonth() + 1, 0);
}

const date = new Date(2024, 1, 15); // Fevereiro
console.log(getFirstDayOfMonth(date).getDate()); // 1
console.log(getLastDayOfMonth(date).getDate()); // 29 (2024 é bissexto)
```

### Verificar Ano Bissexto

```javascript
function isLeapYear(year) {
  return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0;
}

console.log(isLeapYear(2024)); // true
console.log(isLeapYear(2023)); // false
console.log(isLeapYear(2000)); // true
console.log(isLeapYear(1900)); // false
```

### Calcular Idade

```javascript
function calculateAge(birthDate) {
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthDate.getDate())
  ) {
    age--;
  }

  return age;
}

const birth = new Date(1990, 5, 15); // 15 de Junho de 1990
console.log(`Idade: ${calculateAge(birth)} anos`);
```

### Nome do Dia da Semana

```javascript
function getDayName(date) {
  const days = [
    "Domingo",
    "Segunda",
    "Terça",
    "Quarta",
    "Quinta",
    "Sexta",
    "Sábado",
  ];
  return days[date.getDay()];
}

console.log(getDayName(new Date())); // Ex: "Quinta"
```

### Nome do Mês

```javascript
function getMonthName(date) {
  const months = [
    "Janeiro",
    "Fevereiro",
    "Março",
    "Abril",
    "Maio",
    "Junho",
    "Julho",
    "Agosto",
    "Setembro",
    "Outubro",
    "Novembro",
    "Dezembro",
  ];
  return months[date.getMonth()];
}

console.log(getMonthName(new Date(2024, 0))); // "Janeiro"
```

## Formatos de Parse Suportados

Os seguintes formatos podem ser parseados:

- ISO 8601: `"2021-01-15T10:30:00Z"`
- RFC 2822: `"Thu, 15 Jan 2021 10:30:00 GMT"`
- Data simples: `"2021-01-15"`
- Data US: `"01/15/2021"`
- Data textual: `"Jan 15, 2021"`

## Veja Também

- [Numbers](./numbers.md) - Operações matemáticas
- [Timers](./timers.md) - setTimeout, setInterval
