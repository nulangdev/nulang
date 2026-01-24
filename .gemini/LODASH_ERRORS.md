# Nulang Lodash Compatibility - Error Documentation

**Data**: 2026-01-24  
**Status**: 72/86 testes passando (84%)

---

## Resumo Executivo

Este documento descreve os erros e bloqueios identificados durante o trabalho de compatibilidade entre Nulang e Lodash v4.17.23. Os problemas estão organizados por categoria e prioridade.

---

## ✅ Problema Crítico Resolvido: Regex Lookahead

### Descrição

O Go's `regexp` padrão (baseado em RE2) não suporta lookahead `(?=...)` e lookbehind `(?<=...)`. Isso quebrava o `_.words` do Lodash, que usa lookahead para separar palavras em strings com hífens, underscores ou camelCase.

### Solução

Substituímos `regexp` por `regexp2` (`github.com/dlclark/regexp2`), que suporta PCRE-style regex com lookahead/lookbehind.

### Funções Lodash que Agora Funcionam

- ✅ `_.words` - agora divide strings corretamente
- ✅ `_.camelCase`
- ✅ `_.kebabCase`
- ✅ `_.snakeCase`
- ✅ `_.lowerCase`
- ✅ `_.upperCase` (antes falhava com hífens)
- ✅ `_.trim`

---

## ✅ Problema Crítico Resolvido: Regex Literal Parsing

### Descrição Original

O parsing de literais regex contendo sequências de escape (como `\s`, `\d`, `\w`) estava incorreto.

### Solução

Implementamos:

- Campo `Position` no `Token` para guardar posição inicial
- Método `SetPosition()` no lexer para restauração
- Lógica em `parseRegexLiteral()` para restaurar posição antes de `ScanRegex()`

---

## 🟠 Erros de Lógica Pendentes

### 1. `_.find` - Retorno Incorreto

**Sintoma**:

```javascript
_.find([1, 2, 3, 4], (n) => n > 2); // Retorna 4, esperado 3
```

**Status**: 🔍 Investigação pendente

---

### 2. `_.uniq` - Não Remove Duplicatas

**Sintoma**:

```javascript
_.uniq([1, 2, 1, 3, 2]); // Retorna [1, 2, 1, 3, 2], esperado [1, 2, 3]
```

**Status**: 🔍 Investigação pendente

---

### 3. `_.assign` - Não Mescla Objetos

**Sintoma**:

```javascript
_.assign({ a: 1 }, { b: 2 }, { c: 3 }); // Retorna {a: 1}, esperado {a: 1, b: 2, c: 3}
```

**Status**: 🔍 Investigação pendente

---

### 4. `_.groupBy` / `_.countBy`

**Sintoma**: Retorna `false`

**Status**: 🔍 Investigação pendente

---

### 5. `_.has`

**Sintoma**: Retorna `false`

**Status**: 🔍 Investigação pendente

---

### 6. `_.invert`

**Sintoma**: Retorna `false`

**Status**: 🔍 Investigação pendente

---

## 🔴 Erros "not a function: UNDEFINED"

Estes erros indicam que uma função helper interna do Lodash não está sendo resolvida corretamente.

### Funções Afetadas

- `_.merge`
- `_.get`
- `_.isPlainObject`
- `_.defaultsDeep`

### Causa Provável

Acesso a propriedades do `Object.prototype` ou verificação de prototype chain que não está implementada em Nulang.

**Status**: 🔍 Investigação pendente

---

## 🟡 Erros de Implementação Faltante

### 1. `_.isEqual` - Deep Equality

**Sintoma**: Retorna `false` para objetos que deveriam ser iguais.

**Causa Provável**: Comparação recursiva de objetos/arrays não implementada corretamente.

---

### 2. `_.memoize`

**Sintoma**: Retorna `false` em vez de valor cacheado.

**Causa Provável**: Problema no mecanismo de cache ou wrapper de função.

---

## ✅ Problemas Resolvidos

| Problema                           | Solução                         | Arquivo                      |
| ---------------------------------- | ------------------------------- | ---------------------------- |
| `Object(arr)` retornava `{}`       | Retornar referência original    | `evaluator/builtins.go`      |
| `function.length` não existia      | Implementar propriedade         | `evaluator/functions.go`     |
| `obj.prop++` falhava               | Suportar MemberExpression       | `evaluator/operators.go`     |
| `arr[i]++` falhava                 | Suportar IndexExpression        | `evaluator/operators.go`     |
| `Array.sort()` não existia         | Implementar método              | `evaluator/array_methods.go` |
| `NaN`/`Infinity` não existiam      | Adicionar constantes            | `evaluator/builtins.go`      |
| `str['method']()` falhava          | Delegar para evalStringProperty | `evaluator/functions.go`     |
| **Regex escape stripping**         | SetPosition + recompilação      | `parser/parser.go`           |
| **Regex lookahead não funcionava** | Migração para regexp2           | `evaluator/regexp.go`        |

---

## Próximos Passos Recomendados

1. **PRIORIDADE ALTA**: Resolver erros "not a function: UNDEFINED"
   - Investigar acesso a `Object.prototype`
   - Verificar resolução de propriedades em prototype chain

2. **PRIORIDADE MÉDIA**: Debugar `_.find` (off-by-one) e `_.uniq` (Set usage)
   - Adicionar logs nas funções internas do Lodash
   - Verificar comportamento de iteradores

3. **PRIORIDADE BAIXA**: Implementar métodos faltantes
   - Deep equality para `_.isEqual`
   - Cache mechanism para `_.memoize`

---

## Debug Commands

```bash
# Testar regex lookahead
printf 'var re = /[a-z]+(?=[A-Z])/g;\nconsole.log("helloWorld".match(re));\n' > /tmp/t.js
./nulang /tmp/t.js

# Testar _.words
printf 'import _ from "lodash";\nconsole.log(_.words("hello-world"));\n' > examples/t.js
./nulang examples/t.js

# Rodar todos os testes Lodash
./nulang examples/test_lodash.js 2>&1 | grep -E "(✓|✗|Passed|Failed)"

# Recompilar após mudanças
go build -o nulang .
```
