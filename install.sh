#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Banner
echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════╗
║                                       ║
║            NU INSTALLER               ║
║                                       ║
╚═══════════════════════════════════════╝
EOF
echo -e "${NC}"

# Detectar sistema operacional e arquitetura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}✗ Arquitetura não suportada: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}➜ Sistema detectado: ${GREEN}${OS}-${ARCH}${NC}"

# Verificar se está no macOS
if [[ "$OS" != "darwin" && "$OS" != "linux" ]]; then
    echo -e "${RED}✗ Sistema operacional não suportado: $OS${NC}"
    echo -e "${YELLOW}⚠ Nu suporta apenas macOS e Linux${NC}"
    exit 1
fi

# Diretório de instalação
INSTALL_DIR="/usr/local/bin"
NU_BIN="$INSTALL_DIR/nu"

# URL do binário
GITHUB_REPO="nulangdev/nulang"

# Buscar a última versão da API do GitHub
echo -e "${BLUE}➜ Buscando última versão...${NC}"
VERSION=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}✗ Não foi possível obter a versão mais recente${NC}"
    echo -e "${YELLOW}ℹ Tentando com a versão padrão v1.0.0...${NC}"
    VERSION="v1.0.0"
fi

echo -e "${GREEN}➜ Versão encontrada: ${VERSION}${NC}"
BINARY_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/nu-${OS}-${ARCH}"


# Verificar se já está instalado
if [ -f "$NU_BIN" ]; then
    echo -e "${YELLOW}⚠ Nu já está instalado em $NU_BIN${NC}"
    read -p "Deseja reinstalar? (s/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[SsYy]$ ]]; then
        echo -e "${BLUE}ℹ Instalação cancelada${NC}"
        exit 0
    fi
fi

echo -e "${BLUE}➜ Baixando Nu...${NC}"

# Criar diretório temporário
TMP_DIR=$(mktemp -d)
TMP_BIN="$TMP_DIR/nu"

# Baixar binário
if command -v curl &> /dev/null; then
    curl -fsSL "$BINARY_URL" -o "$TMP_BIN"
elif command -v wget &> /dev/null; then
    wget -q "$BINARY_URL" -O "$TMP_BIN"
else
    echo -e "${RED}✗ curl ou wget não encontrado. Por favor, instale um deles.${NC}"
    exit 1
fi

# Verificar se o download foi bem-sucedido
if [ ! -f "$TMP_BIN" ]; then
    echo -e "${RED}✗ Falha ao baixar o binário${NC}"
    echo -e "${YELLOW}ℹ URL tentada: $BINARY_URL${NC}"
    exit 1
fi

# Tornar executável
chmod +x "$TMP_BIN"

# Instalar (pode precisar de sudo)
echo -e "${BLUE}➜ Instalando Nu em $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_BIN" "$NU_BIN"
else
    echo -e "${YELLOW}⚠ Permissões de administrador necessárias${NC}"
    sudo mv "$TMP_BIN" "$NU_BIN"
fi

# Limpar arquivos temporários
rm -rf "$TMP_DIR"

# Verificar instalação
if [ ! -f "$NU_BIN" ]; then
    echo -e "${RED}✗ Falha na instalação${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Nu instalado com sucesso!${NC}"

# Configurar shell (zsh)
configure_shell() {
    local SHELL_RC="$HOME/.zshrc"
    
    # Detectar qual shell está sendo usado
    if [[ "$SHELL" == *"bash"* ]]; then
        SHELL_RC="$HOME/.bashrc"
    elif [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    echo -e "${BLUE}➜ Configurando shell ($SHELL_RC)...${NC}"
    
    # Verificar se já existe configuração
    if grep -q "# Nu configuration" "$SHELL_RC" 2>/dev/null; then
        echo -e "${YELLOW}⚠ Configuração do Nu já existe em $SHELL_RC${NC}"
        return
    fi
    
    # Adicionar configuração ao shell
    cat >> "$SHELL_RC" << 'EOF'

# Nu configuration
export PATH="$PATH:/usr/local/bin"

# Autocompletar para arquivos .js e .ts
if [ -n "$ZSH_VERSION" ]; then
    # Zsh completion
    _nu_completion() {
        local -a js_files ts_files
        js_files=(*.js)
        ts_files=(*.ts)
        _describe 'nu files' js_files ts_files
    }
    compdef _nu_completion nu
fi
EOF
    
    echo -e "${GREEN}✓ Shell configurado com sucesso!${NC}"
    echo -e "${YELLOW}ℹ Execute 'source $SHELL_RC' ou reinicie o terminal para aplicar as mudanças${NC}"
}

# Perguntar se deseja configurar o shell
echo ""
read -p "Deseja configurar o shell automaticamente? (S/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Nn]$ ]]; then
    configure_shell
fi

echo ""
echo -e "${BLUE}╔═══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  ${GREEN}Instalação concluída!${NC}              ${BLUE}║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}➜ Versão instalada:${NC}"
"$NU_BIN" --version 2>/dev/null || echo "  Nu (development build)"
echo ""
echo -e "${GREEN}➜ Comandos disponíveis:${NC}"
echo -e "  ${BLUE}nu${NC} <arquivo.js>  - Executar arquivo JavaScript"
echo -e "  ${BLUE}nu${NC} <arquivo.ts>  - Executar arquivo TypeScript"
echo -e "  ${BLUE}nu${NC}              - REPL interativo"
echo ""
echo -e "${GREEN}➜ Exemplo de uso:${NC}"
echo -e "  ${BLUE}nu index.js${NC}"
echo -e "  ${BLUE}nu app.ts${NC}"
echo ""
echo -e "${YELLOW}➜ Para desinstalar:${NC}"
echo -e "  ${BLUE}curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/uninstall.sh | bash${NC}"
echo ""
echo -e "${GREEN}Obrigado por usar Nu! 🚀${NC}"
echo ""
