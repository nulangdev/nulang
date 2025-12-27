# 🎯 Quick Start - Instalação e Uso do Nulang

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
echo 'console.log("Hello, Nulang! 🚀");' > test.nu

# Executar - 3 formas diferentes
nulang test.nu   # Forma padrão
nu test.nu       # Usando alias
run_nu test.nu   # Usando função helper
```

## 📦 Exemplos Práticos

### 1. Hello World

```javascript
// hello.nu
console.log("Hello, World!");
```

```bash
nu hello.nu
```

### 2. Servidor HTTP Simples

```javascript
// server.nu
const http = require("http");

const server = http.createServer((req, res) => {
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("Hello from Nulang! 🚀\n");
});

server.listen(3000, () => {
  console.log("Servidor rodando em http://localhost:3000");
});
```

```bash
nu server.nu
```

### 3. Manipulação de Arquivos

```javascript
// files.nu
const fs = require("fs");

// Escrever
fs.writeFileSync("data.txt", "Nulang é incrível!");

// Ler
const content = fs.readFileSync("data.txt", "utf8");
console.log(content);

// Listar diretório
const files = fs.readdirSync(".");
console.log("Arquivos:", files);
```

```bash
nu files.nu
```

### 4. Fetch API

```javascript
// fetch.nu
const response = fetch("https://api.github.com/users/nulangdev");
const data = response.json();
console.log("GitHub User:", data.name);
console.log("Repos:", data.public_repos);
```

```bash
nu fetch.nu
```

## 🛠️ Comandos Úteis

```bash
# Verificar instalação
which nulang

# Executar arquivo
nulang script.nu

# REPL interativo
nulang

# Desinstalar
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/uninstall.sh | bash
```

## 📚 Próximos Passos

1. 📖 Leia a [documentação completa](./doc/README.md)
2. 🎯 Veja mais [exemplos](./examples/)
3. 🚀 Comece a criar!

## ❓ Problemas?

### "Command not found: nulang"

```bash
# Verifique o PATH
echo $PATH | grep /usr/local/bin

# Se não aparecer, adicione:
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Permission denied

```bash
chmod +x /usr/local/bin/nulang
```

### Reinstalar

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
# Digite 's' quando perguntar se quer reinstalar
```

---

**Pronto para começar! 🚀**
