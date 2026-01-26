# Nulang Lodash Compatibility - Error Documentation

**Última Atualização**: 2026-01-26  
**Status**: 76/86 testes passando (88%)

---

## 🎯 Resumo Executivo

Este documento descreve os erros e as soluções implementadas durante o trabalho de compatibilidade entre Nulang e Lodash v4.17.23.

**Progresso Total**: 67/86 (78%) → 77/86 (90%) → 76/86 (88%)

Os problemas estão organizados por categoria e prioridade.

---

## ✅ Correções Implementadas (2026-01-26)

### 1. Prototype Chain e Object.prototype

Implementamos a busca de propriedades no prototype chain:

- **getObjectPrototype()**: Função helper para obter Object.prototype
- **evalMemberExpression**: Fallback para Object.prototype quando propriedade não existe
- **evalIndexExpression**: Mesmo fallback para bracket notation

Métodos do Object.prototype implementados:

- `toString()` - Retorna `[object Type]` para type detection
- `hasOwnProperty(prop)` - Verifica propriedades próprias
- `valueOf()` - Retorna o valor do objeto
- `constructor` - Aponta para Object constructor

### 2. typeof para ObjectMap callable

Modificamos `evalTypeofExpression` para retornar `"function"` para ObjectMaps que têm `__call__`, como o construtor Object.

### 3. Array.from para Sets e Maps

Expandimos `Array.from` para suportar:

- Sets (via `_set` property)
- Maps (via `_map` property)
- Iterables com method `values()`

### 4. hasOwnProperty Multi-Pattern

Implementamos `hasOwnProperty` para funcionar com múltiplos padrões de chamada:

- `obj.hasOwnProperty('prop')` - wrapper prepends object
- `Object.prototype.hasOwnProperty.call(obj, 'prop')` - .call pattern

### 5. toString Type Detection

Corrigimos `Object.prototype.toString` para encontrar o objeto correto nos argumentos quando chamado via `.call()`, evitando retornar `[object Object]` para funções.

---

## 🟠 Testes Pendentes (10 testes)

### 1. `_.find` - Off-by-One Bug

**Sintoma**:

```javascript
_.find([1, 2, 3, 4], (n) => n > 2); // Retorna 4, esperado 3
```

**Análise**: O `_.findIndex` retorna o índice correto (2), mas `_.find` retorna o elemento seguinte. Parece ser um bug no timing de captura do valor dentro do lodash.

---

### 2. `_.uniq` - Não Remove Duplicatas

**Sintoma**:

```javascript
_.uniq([1, 2, 1, 3, 2]); // Retorna [1, 2, 1, 3, 2]
```

**Análise**: `Array.from(set)` agora funciona, mas o lodash usa padrões internos diferentes que ainda não estão funcionando.

---

### 3. `_.merge` - Deep Merge Falha

**Sintoma**:

```javascript
_.merge({ a: { x: 1 } }, { a: { y: 2 } }); // Retorna {a: {y: 2}}
```

**Análise**: Não está fazendo merge profundo corretamente.

---

### 4. `_.get` com String Path - "not a function: UNDEFINED"

**Sintoma**:

```javascript
_.get(obj, "a.b.c"); // Error: not a function: UNDEFINED
_.get(obj, ["a", "b", "c"]); // Funciona!
```

**Análise**: O problema está na conversão de string path para array. Array path funciona.

---

### 5. `_.entries` - "not a function: UNDEFINED"

**Sintoma**: Erro interno no lodash mesmo que `Object.entries` funciona.

---

### 6. `_.invert` - Keys Incorretas

**Sintoma**:

```javascript
_.invert({ a: "1", b: "2" }); // Retorna {[object String]: ...}
```

**Análise**: Problema na conversão de valores para strings como keys de objeto.

---

### 7. `_.isEqual` - "not a function: UNDEFINED"

**Análise**: Erro interno no lodash.

---

### 8. `_.isPlainObject` - Retorna false para {}

**Análise**: Todos os checks manuais passam (typeof, constructor, hasOwnProperty), mas algo interno no lodash ainda falha.

---

### 9. `_.defaultsDeep` - "not a function: UNDEFINED"

**Análise**: Erro interno no lodash.

---

## ✅ Problemas Resolvidos Anteriormente

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

1. **PRIORIDADE ALTA**: Investigar erros "not a function: UNDEFINED"
   - Relacionado ao acesso interno de propriedades no lodash
   - Pode ser um bug na conversão de string path

2. **PRIORIDADE MÉDIA**: Debugar `_.find` (off-by-one)
   - Investigar o loop interno do lodash
3. **PRIORIDADE MÉDIA**: Investigar `_.isPlainObject`
   - Verificar qual check exato está falhando dentro do lodash

---

## Debug Commands

```bash
# Testar _.isPlainObject
printf 'import _ from "lodash";\nconsole.log(_.isPlainObject({}));\n' > examples/t.js
./nulang examples/t.js

# Rodar todos os testes Lodash
./nulang examples/test_lodash.js 2>&1 | grep -E "(✓|✗|Passed|Failed)"

# Recompilar após mudanças
go build -o nulang .
```
