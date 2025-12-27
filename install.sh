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
║         NULANG INSTALLER              ║
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
    echo -e "${YELLOW}⚠ Nulang suporta apenas macOS e Linux${NC}"
    exit 1
fi

# Diretório de instalação
INSTALL_DIR="/usr/local/bin"
NULANG_BIN="$INSTALL_DIR/nulang"

# URL do binário (você precisará hospedar o binário em algum lugar)
GITHUB_REPO="nulangdev/nulang"
VERSION="latest"
BINARY_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/nulang-${OS}-${ARCH}"

# Verificar se já está instalado
if [ -f "$NULANG_BIN" ]; then
    echo -e "${YELLOW}⚠ Nulang já está instalado em $NULANG_BIN${NC}"
    read -p "Deseja reinstalar? (s/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[SsYy]$ ]]; then
        echo -e "${BLUE}ℹ Instalação cancelada${NC}"
        exit 0
    fi
fi

echo -e "${BLUE}➜ Baixando Nulang...${NC}"

# Criar diretório temporário
TMP_DIR=$(mktemp -d)
TMP_BIN="$TMP_DIR/nulang"

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
echo -e "${BLUE}➜ Instalando Nulang em $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_BIN" "$NULANG_BIN"
else
    echo -e "${YELLOW}⚠ Permissões de administrador necessárias${NC}"
    sudo mv "$TMP_BIN" "$NULANG_BIN"
fi

# Limpar arquivos temporários
rm -rf "$TMP_DIR"

# Verificar instalação
if [ ! -f "$NULANG_BIN" ]; then
    echo -e "${RED}✗ Falha na instalação${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Nulang instalado com sucesso!${NC}"

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
    if grep -q "# Nulang configuration" "$SHELL_RC" 2>/dev/null; then
        echo -e "${YELLOW}⚠ Configuração do Nulang já existe em $SHELL_RC${NC}"
        return
    fi
    
    # Adicionar configuração ao shell
    cat >> "$SHELL_RC" << 'EOF'

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

# Autocompletar para arquivos .nu
if [ -n "$ZSH_VERSION" ]; then
    # Zsh completion
    _nulang_completion() {
        local -a nu_files
        nu_files=(*.nu)
        _describe 'nulang files' nu_files
    }
    compdef _nulang_completion nulang nu run_nu
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

# Exibir informações da versão
echo ""
echo -e "${BLUE}╔═══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}  ${GREEN}Instalação concluída!${NC}              ${BLUE}║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}➜ Versão instalada:${NC}"
"$NULANG_BIN" --version 2>/dev/null || echo "  Nulang (development build)"
echo ""
echo -e "${GREEN}➜ Comandos disponíveis:${NC}"
echo -e "  ${BLUE}nulang${NC} <arquivo.nu>  - Executar arquivo Nulang"
echo -e "  ${BLUE}nu${NC} <arquivo.nu>      - Alias para nulang"
echo -e "  ${BLUE}run_nu${NC} <arquivo.nu>  - Função helper"
echo ""
echo -e "${GREEN}➜ Exemplo de uso:${NC}"
echo -e "  ${BLUE}nulang index.nu${NC}"
echo -e "  ${BLUE}nu index.nu${NC}"
echo ""
echo -e "${YELLOW}➜ Para desinstalar:${NC}"
echo -e "  ${BLUE}curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/uninstall.sh | bash${NC}"
echo ""
echo -e "${GREEN}Obrigado por usar Nulang! 🚀${NC}"
echo ""
