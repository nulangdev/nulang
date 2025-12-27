# Classes

Classes em Nulang seguem a sintaxe ES6, permitindo programação orientada a objetos.

## Declaração Básica

```javascript
class Pessoa {
  constructor(nome, idade) {
    this.nome = nome;
    this.idade = idade;
  }

  saudar() {
    return `Olá, meu nome é ${this.nome}`;
  }
}

const pessoa = new Pessoa("João", 30);
console.log(pessoa.saudar()); // "Olá, meu nome é João"
```

## Construtor

O método `constructor` é chamado ao criar uma nova instância.

```javascript
class Produto {
  constructor(nome, preco) {
    this.nome = nome;
    this.preco = preco;
    this.criado = Date.now();
  }
}

const prod = new Produto("Laptop", 2500);
console.log(prod.nome); // "Laptop"
console.log(prod.preco); // 2500
```

## Métodos

### Métodos de Instância

```javascript
class Calculadora {
  constructor(valor) {
    this.valor = valor;
  }

  somar(n) {
    this.valor += n;
    return this;
  }

  subtrair(n) {
    this.valor -= n;
    return this;
  }

  resultado() {
    return this.valor;
  }
}

const calc = new Calculadora(10);
const result = calc.somar(5).subtrair(3).resultado();
console.log(result); // 12
```

### Métodos Estáticos

Métodos que pertencem à classe, não às instâncias.

```javascript
class Utils {
  static somar(a, b) {
    return a + b;
  }

  static multiplicar(a, b) {
    return a * b;
  }
}

console.log(Utils.somar(5, 3)); // 8
console.log(Utils.multiplicar(4, 2)); // 8

// Não disponível em instâncias
const u = new Utils();
// u.somar(1, 2) // Erro
```

## Propriedades

### Propriedades de Instância

```javascript
class Carro {
  // Propriedade com valor padrão
  rodas = 4;

  constructor(modelo) {
    this.modelo = modelo;
  }
}

const carro = new Carro("Sedan");
console.log(carro.rodas); // 4
console.log(carro.modelo); // "Sedan"
```

### Propriedades Estáticas

```javascript
class Config {
  static versao = "1.0.0";
  static ambiente = "development";
}

console.log(Config.versao); // "1.0.0"
console.log(Config.ambiente); // "development"
```

## Getters e Setters

```javascript
class Retangulo {
  constructor(largura, altura) {
    this._largura = largura;
    this._altura = altura;
  }

  get largura() {
    return this._largura;
  }

  set largura(valor) {
    if (valor <= 0) {
      throw new Error("Largura deve ser positiva");
    }
    this._largura = valor;
  }

  get altura() {
    return this._altura;
  }

  set altura(valor) {
    if (valor <= 0) {
      throw new Error("Altura deve ser positiva");
    }
    this._altura = valor;
  }

  get area() {
    return this._largura * this._altura;
  }

  get perimetro() {
    return 2 * (this._largura + this._altura);
  }
}

const ret = new Retangulo(10, 5);
console.log(ret.area); // 50
console.log(ret.perimetro); // 30

ret.largura = 20;
console.log(ret.area); // 100
```

## Herança

### extends

```javascript
class Animal {
  constructor(nome) {
    this.nome = nome;
  }

  falar() {
    return `${this.nome} faz um som`;
  }
}

class Cachorro extends Animal {
  constructor(nome, raca) {
    super(nome); // Chama construtor pai
    this.raca = raca;
  }

  falar() {
    return `${this.nome} late!`;
  }

  buscar() {
    return `${this.nome} busca a bola`;
  }
}

const dog = new Cachorro("Rex", "Labrador");
console.log(dog.nome); // "Rex"
console.log(dog.raca); // "Labrador"
console.log(dog.falar()); // "Rex late!"
console.log(dog.buscar()); // "Rex busca a bola"
```

### super

Referência à classe pai.

```javascript
class Forma {
  constructor(cor) {
    this.cor = cor;
  }

  descrever() {
    return `Uma forma ${this.cor}`;
  }
}

class Circulo extends Forma {
  constructor(cor, raio) {
    super(cor); // Chama Forma.constructor
    this.raio = raio;
  }

  descrever() {
    return `${super.descrever()} com raio ${this.raio}`;
  }

  get area() {
    return Math.PI * this.raio * this.raio;
  }
}

const c = new Circulo("vermelho", 5);
console.log(c.descrever()); // "Uma forma vermelho com raio 5"
console.log(c.area); // 78.54...
```

## Encapsulamento

### Convenção de Privado

Use `_` para indicar propriedades privadas (convenção).

```javascript
class ContaBancaria {
  constructor(titular, saldoInicial) {
    this._titular = titular;
    this._saldo = saldoInicial;
    this._historico = [];
  }

  get saldo() {
    return this._saldo;
  }

  depositar(valor) {
    if (valor <= 0) {
      throw new Error("Valor inválido");
    }
    this._saldo += valor;
    this._historico.push({ tipo: "deposito", valor });
  }

  sacar(valor) {
    if (valor > this._saldo) {
      throw new Error("Saldo insuficiente");
    }
    this._saldo -= valor;
    this._historico.push({ tipo: "saque", valor });
  }

  get extrato() {
    return this._historico.slice(); // Retorna cópia
  }
}

const conta = new ContaBancaria("João", 1000);
conta.depositar(500);
conta.sacar(200);
console.log(conta.saldo); // 1300
console.log(conta.extrato); // [{...}, {...}]
```

## Padrões de Design

### Singleton

```javascript
class Database {
  static instance = null;

  constructor() {
    if (Database.instance) {
      return Database.instance;
    }
    this.conexao = null;
    Database.instance = this;
  }

  conectar(url) {
    this.conexao = { url: url, ativa: true };
    console.log("Conectado a:", url);
  }
}

const db1 = new Database();
const db2 = new Database();
console.log(db1 === db2); // true (mesma instância)
```

### Factory

```javascript
class Forma {
  static criar(tipo, ...args) {
    switch (tipo) {
      case "circulo":
        return new Circulo(...args);
      case "retangulo":
        return new Retangulo(...args);
      default:
        throw new Error(`Tipo desconhecido: ${tipo}`);
    }
  }
}

const circulo = Forma.criar("circulo", 5);
const retangulo = Forma.criar("retangulo", 10, 20);
```

### Observer

```javascript
class EventoObservavel {
  constructor() {
    this._observers = [];
  }

  subscribe(fn) {
    this._observers.push(fn);
    return () => {
      this._observers = this._observers.filter((obs) => obs !== fn);
    };
  }

  notify(data) {
    this._observers.forEach((fn) => fn(data));
  }
}

const eventos = new EventoObservavel();

const unsubscribe = eventos.subscribe((data) => {
  console.log("Recebido:", data);
});

eventos.notify("Hello"); // "Recebido: Hello"
unsubscribe();
eventos.notify("World"); // Não imprime
```

### Builder

```javascript
class QueryBuilder {
  constructor() {
    this._table = "";
    this._columns = ["*"];
    this._conditions = [];
  }

  from(table) {
    this._table = table;
    return this;
  }

  select(...columns) {
    this._columns = columns;
    return this;
  }

  where(condition) {
    this._conditions.push(condition);
    return this;
  }

  build() {
    let sql = `SELECT ${this._columns.join(", ")} FROM ${this._table}`;
    if (this._conditions.length > 0) {
      sql += ` WHERE ${this._conditions.join(" AND ")}`;
    }
    return sql;
  }
}

const query = new QueryBuilder()
  .from("users")
  .select("id", "name", "email")
  .where("active = true")
  .where("age > 18")
  .build();

console.log(query);
// "SELECT id, name, email FROM users WHERE active = true AND age > 18"
```

## Herança Múltipla com Mixins

```javascript
// Mixin
const Swimmer = {
  swim() {
    return `${this.nome} está nadando`;
  },
};

const Flyer = {
  fly() {
    return `${this.nome} está voando`;
  },
};

class Animal {
  constructor(nome) {
    this.nome = nome;
  }
}

class Pato extends Animal {
  constructor(nome) {
    super(nome);
    // Aplicar mixins
    Object.assign(this, Swimmer, Flyer);
  }

  falar() {
    return `${this.nome} diz quack!`;
  }
}

const pato = new Pato("Donald");
console.log(pato.swim()); // "Donald está nadando"
console.log(pato.fly()); // "Donald está voando"
console.log(pato.falar()); // "Donald diz quack!"
```

## Estendendo Classes Built-in

```javascript
class MinhaStream extends stream.Readable {
  constructor(dados) {
    super();
    this.dados = dados;
    this.indice = 0;
  }

  _read() {
    if (this.indice < this.dados.length) {
      this.push(this.dados[this.indice]);
      this.indice++;
    } else {
      this.push(null);
    }
  }
}

const s = new MinhaStream(["a", "b", "c"]);
s.on("data", (chunk) => console.log(chunk));
```

## Boas Práticas

### ✅ Use nomes descritivos

```javascript
class UserRepository {} // Bom
class UR {} // Ruim
```

### ✅ Uma responsabilidade por classe

```javascript
// Bom: classes focadas
class Logger {
  /* logging */
}
class Database {
  /* persistência */
}
class UserService {
  /* lógica de usuário */
}
```

### ✅ Prefira composição a herança profunda

```javascript
// Evite herança > 2-3 níveis
class A {}
class B extends A {}
class C extends B {}
class D extends C {} // Profundo demais
```

### ✅ Use getters para valores computados

```javascript
class Retangulo {
  get area() {
    return this.largura * this.altura;
  }
}
```

## Veja Também

- [Getters/Setters](./getters_setters.md) - Propriedades computadas
- [Modules](./modules.md) - Organização de código
- [Events](./events.md) - EventEmitter como classe base
