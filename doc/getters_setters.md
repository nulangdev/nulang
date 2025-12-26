# Getters e Setters

Getters e setters permitem definir métodos que são acessados como propriedades.

## Sintaxe em Classes

```javascript
class Pessoa {
  constructor(nome) {
    this._nome = nome;
  }

  get nome() {
    console.log("Getter chamado");
    return this._nome;
  }

  set nome(valor) {
    console.log("Setter chamado");
    this._nome = valor;
  }
}

const p = new Pessoa("João");
console.log(p.nome); // "Getter chamado", "João"
p.nome = "Maria"; // "Setter chamado"
console.log(p.nome); // "Getter chamado", "Maria"
```

## Casos de Uso

### Validação

```javascript
class Produto {
  constructor() {
    this._preco = 0;
  }

  get preco() {
    return this._preco;
  }

  set preco(valor) {
    if (valor < 0) {
      throw new Error("Preço não pode ser negativo");
    }
    this._preco = valor;
  }
}

const prod = new Produto();
prod.preco = 100; // OK
prod.preco = -50; // Error: Preço não pode ser negativo
```

### Computação Lazy

```javascript
class Circulo {
  constructor(raio) {
    this._raio = raio;
  }

  get raio() {
    return this._raio;
  }

  set raio(valor) {
    this._raio = valor;
  }

  get area() {
    return Math.PI * this._raio * this._raio;
  }

  get circunferencia() {
    return 2 * Math.PI * this._raio;
  }
}

const c = new Circulo(5);
console.log(c.area); // 78.5398...
console.log(c.circunferencia); // 31.4159...
```

### Formatação Automática

```javascript
class Usuario {
  constructor(nome, email) {
    this._nome = nome;
    this._email = email;
  }

  get nome() {
    return this._nome;
  }

  set nome(valor) {
    // Capitaliza automaticamente
    this._nome = valor.charAt(0).toUpperCase() + valor.slice(1).toLowerCase();
  }

  get email() {
    return this._email;
  }

  set email(valor) {
    // Normaliza para minúsculas
    this._email = valor.toLowerCase();
  }
}

const u = new Usuario("", "");
u.nome = "JOÃO";
u.email = "JOAO@EMAIL.COM";

console.log(u.nome); // "João"
console.log(u.email); // "joao@email.com"
```

### Propriedades Read-Only

```javascript
class Configuracao {
  constructor() {
    this._versao = "1.0.0";
    this._criado = Date.now();
  }

  get versao() {
    return this._versao;
  }

  get criadoEm() {
    return new Date(this._criado);
  }

  // Sem setter = read-only
}

const config = new Configuracao();
console.log(config.versao); // "1.0.0"
console.log(config.criadoEm); // Date object
// config.versao = "2.0.0"    // Não faz nada
```

### Cache de Valores Computados

```javascript
class Dados {
  constructor(valores) {
    this._valores = valores;
    this._mediaCache = null;
  }

  get valores() {
    return this._valores;
  }

  set valores(novos) {
    this._valores = novos;
    this._mediaCache = null; // Invalida cache
  }

  get media() {
    if (this._mediaCache === null) {
      console.log("Calculando média...");
      const sum = this._valores.reduce((a, b) => a + b, 0);
      this._mediaCache = sum / this._valores.length;
    }
    return this._mediaCache;
  }
}

const d = new Dados([1, 2, 3, 4, 5]);
console.log(d.media); // "Calculando média...", 3
console.log(d.media); // 3 (do cache)
d.valores = [10, 20, 30];
console.log(d.media); // "Calculando média...", 20
```

### Propriedades Derivadas

```javascript
class Retangulo {
  constructor(largura, altura) {
    this._largura = largura;
    this._altura = altura;
  }

  get largura() {
    return this._largura;
  }
  set largura(v) {
    this._largura = v;
  }

  get altura() {
    return this._altura;
  }
  set altura(v) {
    this._altura = v;
  }

  get area() {
    return this._largura * this._altura;
  }

  set area(novaArea) {
    // Recalcula mantendo proporção
    const razao = this._largura / this._altura;
    this._altura = Math.sqrt(novaArea / razao);
    this._largura = razao * this._altura;
  }

  get perimetro() {
    return 2 * (this._largura + this._altura);
  }
}

const r = new Retangulo(4, 3);
console.log(r.area); // 12
console.log(r.perimetro); // 14

r.area = 48;
console.log(r.largura); // ~8
console.log(r.altura); // ~6
```

### Encapsulamento

```javascript
class ContaBancaria {
  constructor(titular, saldoInicial) {
    this._titular = titular;
    this._saldo = saldoInicial;
    this._transacoes = [];
  }

  get titular() {
    return this._titular;
  }

  get saldo() {
    return this._saldo;
  }

  get transacoes() {
    // Retorna cópia para proteger o original
    return this._transacoes.slice();
  }

  depositar(valor) {
    if (valor <= 0) {
      throw new Error("Valor inválido");
    }
    this._saldo += valor;
    this._transacoes.push({ tipo: "deposito", valor: valor });
  }

  sacar(valor) {
    if (valor > this._saldo) {
      throw new Error("Saldo insuficiente");
    }
    this._saldo -= valor;
    this._transacoes.push({ tipo: "saque", valor: valor });
  }
}

const conta = new ContaBancaria("João", 1000);
conta.depositar(500);
conta.sacar(200);
console.log(conta.saldo); // 1300
console.log(conta.transacoes); // [{...}, {...}]
// conta.saldo = 9999999     // Não funciona (read-only)
```

## Convenções

### Underscore para Privado

Use `_` para indicar propriedades internas:

```javascript
class Exemplo {
  constructor() {
    this._valor = 0; // Convenção: privado
  }

  get valor() {
    return this._valor;
  }
}
```

### Nome do Getter/Setter = Nome da "Propriedade"

```javascript
class Pessoa {
  get nome() {
    /* ... */
  } // Acesso: pessoa.nome
  set nome(v) {
    /* ... */
  } // Atribuição: pessoa.nome = "x"

  get nomeCompleto() {
    /* ... */
  } // pessoa.nomeCompleto
}
```

## Diferença de Métodos

```javascript
class Exemplo {
  // Método - requer parênteses
  calcular() {
    return this._valor * 2;
  }

  // Getter - sem parênteses
  get calculado() {
    return this._valor * 2;
  }
}

const e = new Exemplo();
e.calcular(); // Método
e.calculado; // Getter (parece propriedade)
```

## Boas Práticas

### ✅ Use Getters para Computações Simples

```javascript
get total() {
  return this.preco * this.quantidade
}
```

### ✅ Use Setters para Validação

```javascript
set idade(valor) {
  if (valor < 0 || valor > 150) {
    throw new Error("Idade inválida")
  }
  this._idade = valor
}
```

### ❌ Evite Efeitos Colaterais em Getters

```javascript
// ❌ Ruim
get valor() {
  this._acessos++  // Efeito colateral
  return this._valor
}

// ✅ Bom
get valor() {
  return this._valor
}
```

### ❌ Evite Operações Custosas sem Cache

```javascript
// ❌ Ruim - recalcula sempre
get soma() {
  return this.items.reduce((a, b) => a + b, 0)
}

// ✅ Bom - com cache
get soma() {
  if (this._somaCache === null) {
    this._somaCache = this.items.reduce((a, b) => a + b, 0)
  }
  return this._somaCache
}
```

## Veja Também

- [Classes](./classes.md) - Definição de classes
- [Modules](./modules.md) - Sistema de módulos
