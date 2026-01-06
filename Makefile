# Makefile para Nulang
# Facilita o build, instalação e distribuição do projeto

.PHONY: build install uninstall clean release help test

# Variáveis
BINARY_NAME=nu
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Diretórios
INSTALL_DIR=/usr/local/bin
DIST_DIR=releases/download

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m

# Ajuda
help:
	@echo "$(BLUE)Nu - Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Comandos disponíveis:$(NC)"
	@echo "  $(YELLOW)make build$(NC)         - Compilar o binário"
	@echo "  $(YELLOW)make install$(NC)       - Instalar localmente"
	@echo "  $(YELLOW)make uninstall$(NC)     - Desinstalar"
	@echo "  $(YELLOW)make clean$(NC)         - Limpar arquivos de build"
	@echo "  $(YELLOW)make release$(NC)       - Criar release multi-plataforma"
	@echo "  $(YELLOW)make test$(NC)          - Executar testes"
	@echo "  $(YELLOW)make test-install$(NC)  - Testar instalação local"
	@echo ""

# Build local
build:
	@echo "$(BLUE)➜ Compilando Nu...$(NC)"
	@go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Build concluído: ./$(BINARY_NAME)$(NC)"

# Instalar localmente
install: build
	@echo "$(BLUE)➜ Instalando em $(INSTALL_DIR)...$(NC)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Nu instalado com sucesso!$(NC)"
	@echo "$(YELLOW)ℹ Execute './install.sh' para configurar o shell$(NC)"

# Desinstalar
uninstall:
	@echo "$(BLUE)➜ Desinstalando Nu...$(NC)"
	@if [ -w "$(INSTALL_DIR)" ]; then \
		rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Nu desinstalado!$(NC)"

# Limpar arquivos de build
clean:
	@echo "$(BLUE)➜ Limpando arquivos...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -rf $(DIST_DIR)
	@echo "$(GREEN)✓ Limpeza concluída!$(NC)"

# Criar release para múltiplas plataformas
release: clean
	@echo "$(BLUE)➜ Criando releases...$(NC)"
	@mkdir -p $(DIST_DIR)
	
	@echo "$(YELLOW)  • macOS ARM64...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "$(YELLOW)  • macOS AMD64...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "$(YELLOW)  • Linux ARM64...$(NC)"
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	@echo "$(YELLOW)  • Linux AMD64...$(NC)"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	@echo "$(GREEN)✓ Releases criados em ./$(DIST_DIR)/$(NC)"
	@echo ""
	@echo "$(BLUE)📦 Arquivos gerados:$(NC)"
	@ls -lh $(DIST_DIR)/
	@echo ""
	@echo "$(YELLOW)➜ Gerando checksums...$(NC)"
	@cd $(DIST_DIR) && shasum -a 256 * > checksums.txt
	@echo "$(GREEN)✓ Checksums salvos em $(DIST_DIR)/checksums.txt$(NC)"

# Executar testes
test:
	@echo "$(BLUE)➜ Executando testes...$(NC)"
	@go test -v ./...

# Testar instalação local
test-install: install
	@echo "$(BLUE)➜ Testando instalação...$(NC)"
	@which $(BINARY_NAME)
	@$(BINARY_NAME) --version || echo "$(YELLOW)⚠ Flag --version não implementada$(NC)"
	@echo "$(GREEN)✓ Instalação OK!$(NC)"

# Build com informações de debug
debug:
	@echo "$(BLUE)➜ Compilando com debug...$(NC)"
	@go build -gcflags="all=-N -l" -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Build debug concluído$(NC)"

# Verificar dependências
deps:
	@echo "$(BLUE)➜ Verificando dependências...$(NC)"
	@go mod download
	@go mod verify
	@echo "$(GREEN)✓ Dependências OK!$(NC)"

# Atualizar dependências
update-deps:
	@echo "$(BLUE)➜ Atualizando dependências...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependências atualizadas!$(NC)"

# Executar um arquivo de exemplo
run-example:
	@echo "$(BLUE)➜ Executando exemplo...$(NC)"
	@./$(BINARY_NAME) examples/index.js
