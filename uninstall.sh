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
║       NULANG UNINSTALLER              ║
║                                       ║
╚═══════════════════════════════════════╝
EOF
echo -e "${NC}"

INSTALL_DIR="/usr/local/bin"
NULANG_BIN="$INSTALL_DIR/nulang"

# Verificar se Nulang está instalado
if [ ! -f "$NULANG_BIN" ]; then
    echo -e "${RED}✗ Nulang não está instalado em $NULANG_BIN${NC}"
    exit 1
fi

echo -e "${YELLOW}⚠ Isso irá remover o Nulang do seu sistema${NC}"
read -p "Tem certeza que deseja continuar? (s/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[SsYy]$ ]]; then
    echo -e "${BLUE}ℹ Desinstalação cancelada${NC}"
    exit 0
fi

# Remover binário
echo -e "${BLUE}➜ Removendo binário...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    rm -f "$NULANG_BIN"
else
    sudo rm -f "$NULANG_BIN"
fi

# Verificar remoção
if [ -f "$NULANG_BIN" ]; then
    echo -e "${RED}✗ Falha ao remover o binário${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Binário removido com sucesso!${NC}"

# Remover configurações do shell
remove_shell_config() {
    local SHELL_RC="$HOME/.zshrc"
    
    # Detectar qual shell está sendo usado
    if [[ "$SHELL" == *"bash"* ]]; then
        SHELL_RC="$HOME/.bashrc"
    elif [[ "$SHELL" == *"zsh"* ]]; then
        SHELL_RC="$HOME/.zshrc"
    fi
    
    if [ ! -f "$SHELL_RC" ]; then
        return
    fi
    
    echo -e "${BLUE}➜ Removendo configurações do shell...${NC}"
    
    # Criar backup
    cp "$SHELL_RC" "${SHELL_RC}.backup"
    
    # Remover configuração do Nulang
    sed -i.tmp '/# Nulang configuration/,/^$/d' "$SHELL_RC" 2>/dev/null || \
    sed -i '.tmp' '/# Nulang configuration/,/^$/d' "$SHELL_RC"
    
    rm -f "${SHELL_RC}.tmp"
    
    echo -e "${GREEN}✓ Configurações removidas (backup salvo em ${SHELL_RC}.backup)${NC}"
}

# Perguntar se deseja remover configurações do shell
echo ""
read -p "Deseja remover as configurações do shell? (S/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Nn]$ ]]; then
    remove_shell_config
fi

echo ""
echo -e "${GREEN}✓ Nulang foi desinstalado com sucesso!${NC}"
echo ""
echo -e "${YELLOW}ℹ Para reinstalar, execute:${NC}"
echo -e "  ${BLUE}curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash${NC}"
echo ""
