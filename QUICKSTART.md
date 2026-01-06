# 🎯 Quick Start - Instalação e Uso do Nu

## ⚡ Instalação Rápida (1 comando)

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

Isso vai:

- ✅ Baixar o binário correto para seu sistema
- ✅ Instalar em `/usr/local/bin`
- ✅ Configurar automaticamente seu zsh/bash
- ✅ Adicionar aliases úteis

## 🎮 Uso Imediato

Após a instalação, reinicie o terminal ou execute:

```bash
source ~/.zshrc  # ou ~/.bashrc
```

Então use:

```bash
# Criar arquivo de teste
echo 'console.log("Hello, Nu! 🚀");' > test.js

# Executar
nu test.js
```

## 📦 Exemplos Práticos

### 1. Hello World

```javascript
// hello.js
console.log("Hello, World!");
```

```bash
nu hello.js
```

### 2. Servidor HTTP Simples

```javascript
// server.js
const http = require("http");

const server = http.createServer((req, res) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello from Nu! 🚀\n");
});

server.listen(3000, () => {
  console.log("Servidor rodando em http://localhost:3000");
});
```

```bash
nu server.js
```

### 3. Manipulação de Arquivos

```javascript
// files.js
const fs = require("fs");

// Escrever
fs.writeFileSync("data.txt", "Nu é incrível!");

// Ler
const content = fs.readFileSync("data.txt", "utf8");
console.log(content);

// Listar diretório
const files = fs.readdirSync(".");
console.log("Arquivos:", files);
```

```bash
nu files.js
```

### 4. Fetch API

```javascript
// fetch.js
const response = fetch("https://api.github.com/users/nulangdev");
const data = response.json();
console.log("GitHub User:", data.name);
console.log("Repos:", data.public_repos);
```

```bash
nu fetch.js
```

## 🛠️ Comandos Úteis

```bash
# Verificar instalação
which nu

# Executar arquivo
nu script.js

# REPL interativo
nu

# Desinstalar
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/uninstall.sh | bash
```

## 📚 Próximos Passos

1. 📖 Leia a [documentação completa](./doc/README.md)
2. 🎯 Veja mais [exemplos](./examples/)
3. 🚀 Comece a criar!

## ❓ Problemas?

### "Command not found: nu"

```bash
# Verifique o PATH
echo $PATH | grep /usr/local/bin

# Se não aparecer, adicione:
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Permission denied

```bash
chmod +x /usr/local/bin/nu
```

### Reinstalar

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
# Digite 's' quando perguntar se quer reinstalar
```

---

**Pronto para começar! 🚀**
