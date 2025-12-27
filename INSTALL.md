# 🚀 Instalação do Nulang

Guia completo de instalação e configuração do Nulang no seu sistema.

## 📋 Pré-requisitos

- **macOS** (Intel ou Apple Silicon) ou **Linux** (x64 ou ARM64)
- **Nenhuma dependência adicional** - O binário já vem compilado!

## 🎯 Métodos de Instalação

### 1️⃣ Instalação via Script (Recomendado)

O método mais simples e rápido. Baixa o binário e configura automaticamente seu shell:

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

**Ou usando wget:**

```bash
wget -qO- https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

#### O que o script faz?

- ✅ Detecta automaticamente seu sistema operacional e arquitetura
- ✅ Baixa o binário correto do Nulang
- ✅ Instala em `/usr/local/bin`
- ✅ Configura automaticamente seu shell (zsh ou bash)
- ✅ Adiciona aliases úteis (`nu`, `run_nu`)
- ✅ Configura autocomplete para arquivos `.nu`

### 2️⃣ Instalação via Homebrew (macOS/Linux)

Para usuários do Homebrew:

```bash
# Adicionar o tap (apenas na primeira vez)
brew tap nulangdev/nulang

# Instalar
brew install nulang
```

**Atualizar:**

```bash
brew upgrade nulang
```

### 3️⃣ Instalação Manual

Se preferir instalar manualmente:

#### 1. Baixar o binário

Escolha o binário para sua plataforma:

**macOS (Apple Silicon):**

```bash
curl -L https://github.com/nulangdev/nulang/releases/latest/download/nulang-darwin-arm64 -o nulang
```

**macOS (Intel):**

```bash
curl -L https://github.com/nulangdev/nulang/releases/latest/download/nulang-darwin-amd64 -o nulang
```

**Linux (x64):**

```bash
curl -L https://github.com/nulangdev/nulang/releases/latest/download/nulang-linux-amd64 -o nulang
```

**Linux (ARM64):**

```bash
curl -L https://github.com/nulangdev/nulang/releases/latest/download/nulang-linux-arm64 -o nulang
```

#### 2. Tornar executável e instalar

```bash
chmod +x nulang
sudo mv nulang /usr/local/bin/
```

#### 3. Configurar o shell manualmente

Adicione ao seu `~/.zshrc` ou `~/.bashrc`:

```bash
# Nulang configuration
export PATH="$PATH:/usr/local/bin"

# Alias para executar arquivos .nu diretamente
alias nu='nulang'

# Função para executar scripts Nulang
run_nu() {
    if [ -f "$1" ]; then
        nulang "$1"
    else
        echo "Arquivo não encontrado: $1"
    fi
}

# Autocompletar para arquivos .nu (apenas zsh)
if [ -n "$ZSH_VERSION" ]; then
    _nulang_completion() {
        local -a nu_files
        nu_files=(*.nu)
        _describe 'nulang files' nu_files
    }
    compdef _nulang_completion nulang nu run_nu
fi
```

#### 4. Recarregar o shell

```bash
source ~/.zshrc  # ou ~/.bashrc
```

## 🎮 Uso

### Executar um arquivo Nulang

```bash
# Forma padrão
nulang index.nu

# Usando o alias
nu index.nu

# Usando a função helper
run_nu index.nu
```

### Exemplos

```bash
# Criar um arquivo de exemplo
cat > hello.nu << 'EOF'
console.log("Olá, Nulang! 🚀");

const soma = (a, b) => a + b;
console.log(`2 + 3 = ${soma(2, 3)}`);
EOF

# Executar
nulang hello.nu
```

## 🔧 Verificar Instalação

```bash
# Verificar se está instalado
which nulang

# Visualizar localização
type nulang

# Testar execução
nulang --version  # (se implementado)
```

## 🗑️ Desinstalação

### Via Script

```bash
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/uninstall.sh | bash
```

### Via Homebrew

```bash
brew uninstall nulang
```

### Manual

```bash
# Remover binário
sudo rm /usr/local/bin/nulang

# Remover configurações do shell
# Edite ~/.zshrc ou ~/.bashrc e remova a seção "# Nulang configuration"
```

## 🐛 Solução de Problemas

### "Permission denied"

Se você receber um erro de permissão ao executar o binário:

```bash
chmod +x /usr/local/bin/nulang
```

### "Command not found"

Certifique-se de que `/usr/local/bin` está no seu PATH:

```bash
echo $PATH | grep -q "/usr/local/bin" && echo "OK" || echo "Adicione ao PATH"
```

Para adicionar:

```bash
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Autocomplete não funciona

Certifique-se de que está usando zsh e recarregue o shell:

```bash
exec zsh
```

## 📚 Próximos Passos

Após a instalação:

1. 📖 Leia a [documentação completa](./doc/README.md)
2. 🎯 Explore os [exemplos](./examples/)
3. 🚀 Comece a criar seus próprios scripts!

## 🆘 Suporte

- 📝 [Issues no GitHub](https://github.com/nulangdev/nulang/issues)
- 📧 Contato: [seu-email@exemplo.com]

---

**Desenvolvido com ❤️ pela comunidade Nulang**
