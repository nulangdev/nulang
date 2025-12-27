# Como Publicar uma Release do Nulang

## 📋 Pré-requisitos

1. Certifique-se de que todas as alterações estão commitadas
2. Verifique se o código compila corretamente: `go build .`
3. Teste localmente: `./nulang examples/index.nu`

## 🚀 Passos para Criar uma Release

### 1. Criar uma Tag de Versão

```bash
# Defina a versão (seguindo semver: vMAJOR.MINOR.PATCH)
VERSION="v1.0.0"

# Crie a tag localmente
git tag -a $VERSION -m "Release $VERSION"

# Envie a tag para o GitHub
git push origin $VERSION
```

### 2. O que Acontece Automaticamente

Quando você envia uma tag que começa com `v`, o GitHub Actions:

1. ✅ Detecta a tag e inicia o workflow
2. ✅ Compila binários para 4 plataformas:
   - macOS ARM64 (Apple Silicon)
   - macOS AMD64 (Intel)
   - Linux ARM64
   - Linux AMD64
3. ✅ Gera checksums SHA256 de todos os binários
4. ✅ Cria automaticamente uma release no GitHub
5. ✅ Anexa todos os binários à release
6. ✅ Gera notas de release automaticamente
7. ✅ Atualiza a fórmula do Homebrew

### 3. Verificar a Release

Após alguns minutos, verifique:

```bash
# URL da release
https://github.com/nulangdev/nulang/releases/tag/$VERSION
```

### 4. Testar a Instalação

```bash
# Testar o script de instalação
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

## 🔧 Build Local (Sem Publicar)

Para testar o build de todas as plataformas localmente:

```bash
# Usar o Makefile
make release

# Os binários serão criados em ./releases/download/
ls -lh releases/download/
```

## 🏷️ Convenções de Versionamento

Use **Semantic Versioning** (semver):

- **v1.0.0** - Release inicial estável
- **v1.1.0** - Novas funcionalidades (minor)
- **v1.0.1** - Correções de bugs (patch)
- **v2.0.0** - Mudanças incompatíveis (major)

## ❌ Deletar uma Tag (se necessário)

```bash
# Deletar localmente
git tag -d v1.0.0

# Deletar no GitHub
git push origin :refs/tags/v1.0.0
```

## 🔄 Atualizar uma Release

Se você precisar atualizar uma release existente:

1. Delete a tag antiga (veja acima)
2. Delete a release no GitHub (via web interface)
3. Crie uma nova tag com o mesmo nome
4. O workflow rodará novamente

## 📊 Monitorar o Workflow

Acompanhe o progresso em:

```
https://github.com/nulangdev/nulang/actions
```

## 🐛 Troubleshooting

### Erro 404 ao baixar binário

- Verifique se a release foi criada: `https://github.com/nulangdev/nulang/releases`
- Verifique se os arquivos foram anexados à release
- Aguarde alguns minutos - o workflow pode ainda estar em execução

### Workflow falhou

- Acesse: `https://github.com/nulangdev/nulang/actions`
- Clique no workflow com falha para ver os logs
- Corrija o problema e crie uma nova tag

### Binários não foram gerados

Verifique se:

- A tag começa com `v` (ex: `v1.0.0`, não `1.0.0`)
- O repositório tem permissões de escrita configuradas
- O arquivo `.github/workflows/release.yml` está na branch principal

## 📝 Exemplo Completo

```bash
# 1. Commit todas as alterações
git add .
git commit -m "Preparando release v1.0.0"
git push

# 2. Criar e enviar tag
git tag -a v1.0.0 -m "Release v1.0.0 - Primeira versão estável"
git push origin v1.0.0

# 3. Aguardar workflow (3-5 minutos)
# Acompanhar em: https://github.com/nulangdev/nulang/actions

# 4. Verificar release
open https://github.com/nulangdev/nulang/releases/tag/v1.0.0

# 5. Testar instalação
curl -fsSL https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh | bash
```

## ✅ Checklist Antes de Publicar

- [ ] Código compila sem erros
- [ ] Testes passam (se houver)
- [ ] README.md está atualizado
- [ ] Documentação está atualizada
- [ ] Versão segue semver (vX.Y.Z)
- [ ] Commit message é descritivo
- [ ] Tag message é descritivo
