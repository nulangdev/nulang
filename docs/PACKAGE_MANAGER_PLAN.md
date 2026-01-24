# Plano: Sistema de Gerenciamento de Pacotes Compatível com Node.js

## Objetivo

Migrar o sistema de gerenciamento de pacotes do Nulang de `nulang.yml`/`nulang.lock` para `package.json`/`package-lock.json`, mantendo 100% de compatibilidade com o ecossistema Node.js.

## Estado Atual

| Componente   | Nulang (atual)  | Node.js (alvo)      |
| ------------ | --------------- | ------------------- |
| Manifesto    | `nulang.yml`    | `package.json`      |
| Lock file    | `nulang.lock`   | `package-lock.json` |
| Dependências | `node_modules/` | `node_modules/` ✅  |
| Formato      | YAML            | JSON                |

---

## Fases de Implementação

### 📦 Fase 1: Estruturas de Dados (package.go)

**Arquivo:** `evaluator/package.go`

#### 1.1 Criar estruturas compatíveis com package.json

```go
// PackageJSON represents the package.json structure
type PackageJSON struct {
    Name            string                     `json:"name"`
    Version         string                     `json:"version"`
    Description     string                     `json:"description,omitempty"`
    Main            string                     `json:"main,omitempty"`
    Module          string                     `json:"module,omitempty"`
    Types           string                     `json:"types,omitempty"`
    Scripts         map[string]string          `json:"scripts,omitempty"`
    Dependencies    map[string]string          `json:"dependencies,omitempty"`
    DevDependencies map[string]string          `json:"devDependencies,omitempty"`
    PeerDependencies map[string]string         `json:"peerDependencies,omitempty"`
    Keywords        []string                   `json:"keywords,omitempty"`
    Author          interface{}                `json:"author,omitempty"`      // string ou objeto
    License         string                     `json:"license,omitempty"`
    Repository      interface{}                `json:"repository,omitempty"`  // string ou objeto
    Bugs            interface{}                `json:"bugs,omitempty"`
    Homepage        string                     `json:"homepage,omitempty"`
    Engines         map[string]string          `json:"engines,omitempty"`
    Private         bool                       `json:"private,omitempty"`
    Type            string                     `json:"type,omitempty"`       // "module" ou "commonjs"
    Exports         interface{}                `json:"exports,omitempty"`
    Bin             interface{}                `json:"bin,omitempty"`
}

// PackageLockJSON represents the package-lock.json structure (v3)
type PackageLockJSON struct {
    Name            string                     `json:"name"`
    Version         string                     `json:"version"`
    LockfileVersion int                        `json:"lockfileVersion"`      // 3
    Requires        bool                       `json:"requires,omitempty"`
    Packages        map[string]PackageLockEntry `json:"packages"`
}

type PackageLockEntry struct {
    Version         string            `json:"version"`
    Resolved        string            `json:"resolved"`
    Integrity       string            `json:"integrity"`           // sha512 hash
    Dependencies    map[string]string `json:"dependencies,omitempty"`
    DevDependencies map[string]string `json:"devDependencies,omitempty"`
    Engines         map[string]string `json:"engines,omitempty"`
    Funding         interface{}       `json:"funding,omitempty"`
    License         string            `json:"license,omitempty"`
    Optional        bool              `json:"optional,omitempty"`
    Peer            bool              `json:"peer,omitempty"`
}
```

#### 1.2 Funções de Load/Save

```go
func LoadPackageJSON(path string) (*PackageJSON, error)
func SavePackageJSON(path string, pkg *PackageJSON) error
func LoadPackageLockJSON(path string) (*PackageLockJSON, error)
func SavePackageLockJSON(path string, lock *PackageLockJSON) error
```

---

### 📦 Fase 2: Resolução de Dependências

**Arquivo:** `evaluator/resolver.go` (novo)

#### 2.1 Resolução de versões semver

```go
// SemverRange parses and matches semver ranges
type SemverRange struct {
    Raw        string
    Comparators []SemverComparator
}

// ParseSemverRange parses a semver range string
// Suporte necessário:
// - Exatas: "1.0.0"
// - Ranges: "^1.0.0", "~1.0.0", ">=1.0.0", "<2.0.0"
// - Combinadas: ">=1.0.0 <2.0.0"
// - OU: "1.0.0 || 2.0.0"
// - Wildcard: "*", "1.x", "1.0.x"
func ParseSemverRange(range string) (*SemverRange, error)
func (r *SemverRange) Match(version string) bool
```

#### 2.2 Resolução de dependências

```go
// DependencyResolver handles dependency resolution
type DependencyResolver struct {
    Registry    string  // Default: "https://registry.npmjs.org"
    Cache       string  // ~/.nu/cache
}

// Resolve resolves all dependencies for a package
func (r *DependencyResolver) Resolve(pkg *PackageJSON) (*DependencyTree, error)

// DependencyTree represents the resolved dependency graph
type DependencyTree struct {
    Root     *DependencyNode
    Packages map[string]*DependencyNode
}

type DependencyNode struct {
    Name         string
    Version      string
    Resolved     string   // URL where it was resolved from
    Integrity    string   // SHA-512 hash
    Dependencies map[string]*DependencyNode
}
```

---

### 📦 Fase 3: Registry NPM

**Arquivo:** `evaluator/registry.go` (novo)

#### 3.1 Cliente NPM Registry

```go
// NPMRegistry interacts with npm registry
type NPMRegistry struct {
    BaseURL string   // https://registry.npmjs.org
    Token   string   // Optional auth token
}

// PackageMetadata from npm registry
type PackageMetadata struct {
    Name        string                        `json:"name"`
    Description string                        `json:"description"`
    DistTags    map[string]string             `json:"dist-tags"`  // "latest", "next", etc
    Versions    map[string]PackageVersionInfo `json:"versions"`
    Time        map[string]string             `json:"time"`
}

type PackageVersionInfo struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Main         string            `json:"main"`
    Dependencies map[string]string `json:"dependencies"`
    Dist         PackageDist       `json:"dist"`
}

type PackageDist struct {
    Shasum   string `json:"shasum"`
    Tarball  string `json:"tarball"`
    Integrity string `json:"integrity"`
}

// GetPackageMetadata fetches package info from registry
func (r *NPMRegistry) GetPackageMetadata(name string) (*PackageMetadata, error)

// GetPackageVersion fetches a specific version
func (r *NPMRegistry) GetPackageVersion(name, version string) (*PackageVersionInfo, error)

// DownloadPackage downloads and extracts a package tarball
func (r *NPMRegistry) DownloadPackage(name, version, destPath string) error
```

---

### 📦 Fase 4: Comandos CLI

**Arquivo:** `main.go` + `evaluator/commands.go` (novo)

#### 4.1 Comandos a implementar

| Comando                       | Descrição                       | Prioridade |
| ----------------------------- | ------------------------------- | ---------- |
| `nu init`                     | Criar `package.json` interativo | 🔴 Alta    |
| `nu install`                  | Instalar todas as deps          | 🔴 Alta    |
| `nu install <pkg>`            | Instalar pacote específico      | 🔴 Alta    |
| `nu install <pkg> --save-dev` | Instalar como devDep            | 🟡 Média   |
| `nu uninstall <pkg>`          | Remover pacote                  | 🟡 Média   |
| `nu update`                   | Atualizar dependências          | 🟡 Média   |
| `nu list`                     | Listar dependências instaladas  | 🟢 Baixa   |
| `nu outdated`                 | Mostrar deps desatualizadas     | 🟢 Baixa   |
| `nu run <script>`             | Executar script do package.json | 🔴 Alta    |
| `nu link`                     | Link local de pacote            | 🟢 Baixa   |
| `nu pack`                     | Criar tarball do pacote         | 🟢 Baixa   |
| `nu publish`                  | Publicar no registry            | 🟢 Baixa   |

#### 4.2 Implementação dos comandos

```go
// nu init
func HandleInit() error {
    // Criar package.json com campos padrão
    // Perguntar interativamente: name, version, description, main, license
}

// nu install
func HandleInstall(args []string) error {
    if len(args) == 0 {
        // Instalar todas as deps do package.json
        return InstallAllDependencies()
    }
    // Instalar pacote específico
    return InstallPackage(args[0], flags)
}

// nu run <script>
func HandleRun(script string) error {
    pkg, _ := LoadPackageJSON("package.json")
    if cmd, ok := pkg.Scripts[script]; ok {
        return execCommand(cmd)
    }
    return fmt.Errorf("script '%s' not found", script)
}
```

---

### 📦 Fase 5: Instalação de Pacotes

**Arquivo:** `evaluator/installer.go` (novo)

#### 5.1 Processo de instalação

```go
// Installer handles package installation
type Installer struct {
    Registry   *NPMRegistry
    Resolver   *DependencyResolver
    CacheDir   string   // ~/.nu/cache
    ModulesDir string   // ./node_modules
}

// Install installs all dependencies
func (i *Installer) Install() error {
    // 1. Ler package.json
    // 2. Ler package-lock.json (se existir)
    // 3. Resolver dependências
    // 4. Comparar com lock file
    // 5. Baixar pacotes faltantes
    // 6. Extrair para node_modules
    // 7. Atualizar package-lock.json
}

// InstallPackage installs a specific package
func (i *Installer) InstallPackage(name, version string, isDev bool) error {
    // 1. Buscar versão no registry
    // 2. Resolver dependências transitivas
    // 3. Atualizar package.json
    // 4. Instalar
    // 5. Atualizar package-lock.json
}
```

#### 5.2 Estrutura node_modules

```
node_modules/
├── .package-lock.json    # Cópia do lock para verificação
├── express/
│   ├── package.json
│   ├── index.js
│   └── lib/
│       └── ...
├── lodash/
│   ├── package.json
│   └── ...
└── .bin/                 # Binários dos pacotes
    ├── express
    └── ...
```

---

### 📦 Fase 6: Resolução de Módulos (modules.go)

**Arquivo:** `evaluator/modules.go`

#### 6.1 Atualizar resolução de imports

```go
func resolveModulePath(modulePath string, basePath string) string {
    // 1. Verificar se é builtin
    // 2. Se começa com ./ ou ../ -> resolver relativo
    // 3. Se não, buscar em node_modules:
    //    a. Ler package.json do pacote
    //    b. Ler campo "main" ou "module" ou "exports"
    //    c. Resolver ponto de entrada correto
}

// Adicionar suporte a exports map do package.json
func resolvePackageExports(pkgPath string, subpath string) string {
    // Suporte a:
    // "exports": "./index.js"
    // "exports": { ".": "./index.js", "./utils": "./lib/utils.js" }
    // "exports": { "import": "./esm.js", "require": "./cjs.js" }
}
```

---

### 📦 Fase 7: Cache e Performance

**Arquivo:** `evaluator/cache.go` (novo)

#### 7.1 Cache de pacotes

```go
// CacheManager handles package caching
type CacheManager struct {
    Dir string  // ~/.nu/cache
}

// Estrutura do cache:
// ~/.nu/
// ├── cache/
// │   ├── registry.npmjs.org/
// │   │   ├── express/
// │   │   │   └── 4.18.2.tgz
// │   │   └── lodash/
// │   │       └── 4.17.21.tgz
// │   └── _cacache/        # Content-addressable storage
// │       ├── content-v2/
// │       └── index-v5/
// └── logs/

func (c *CacheManager) Get(name, version string) (string, bool)
func (c *CacheManager) Put(name, version string, tarball []byte) error
func (c *CacheManager) Clean() error
func (c *CacheManager) Verify() error
```

---

### 📦 Fase 8: Migração

**Arquivo:** `evaluator/migrate.go` (novo)

#### 8.1 Migrar projetos existentes

```go
// MigrateFromYAML migrates from nulang.yml to package.json
func MigrateFromYAML(projectPath string) error {
    // 1. Ler nulang.yml
    // 2. Converter para PackageJSON
    // 3. Salvar package.json
    // 4. Converter nulang.lock para package-lock.json
    // 5. (Opcional) Remover arquivos antigos
}
```

---

## Cronograma de Implementação

| Fase | Descrição                 | Esforço | Dependências |
| ---- | ------------------------- | ------- | ------------ |
| 1    | Estruturas de dados       | 2h      | -            |
| 2    | Resolução de dependências | 4h      | Fase 1       |
| 3    | Cliente NPM Registry      | 4h      | Fase 1       |
| 4    | Comandos CLI              | 3h      | Fases 2, 3   |
| 5    | Instalador                | 6h      | Fases 2, 3   |
| 6    | Resolução de módulos      | 3h      | Fase 5       |
| 7    | Cache                     | 2h      | Fase 5       |
| 8    | Migração                  | 1h      | Fase 1       |

**Total estimado:** ~25 horas de desenvolvimento

---

## Arquivos a Criar/Modificar

### Novos Arquivos

```
evaluator/
├── package_json.go      # Estruturas e I/O de package.json
├── package_lock.go      # Estruturas e I/O de package-lock.json
├── resolver.go          # Resolução de dependências + semver
├── registry.go          # Cliente NPM Registry
├── installer.go         # Lógica de instalação
├── cache.go            # Gerenciamento de cache
├── migrate.go          # Migração de nulang.yml
└── commands.go         # Handlers dos comandos CLI
```

### Arquivos a Modificar

```
main.go                  # Novos comandos CLI
evaluator/modules.go     # Resolução com package.json exports
evaluator/package.go     # Manter retrocompatibilidade (deprecated)
```

---

## Compatibilidade

### ✅ Funcionalidades Compatíveis com npm

- `package.json` com todos os campos padrão
- `package-lock.json` v3
- `node_modules` flat structure
- Semver ranges (^, ~, >=, etc)
- npm registry como fonte padrão
- Scripts (npm run)
- Bin executables
- Exports field (ESM)

### ⚠️ Limitações Conhecidas

- **Workspaces:** Não suportado inicialmente (npm 7+)
- **npx:** Não suportado inicialmente
- **Peer dependencies:** Suporte básico
- **Optional dependencies:** Suporte básico
- **Publish:** Requer auth token

---

## Próximos Passos

1. **Aprovar este plano** ✅
2. **Implementar Fase 1** - Estruturas de dados
3. **Testes unitários** para cada fase
4. **Documentação** atualizada

---

## Referências

- [package.json spec](https://docs.npmjs.com/cli/v10/configuring-npm/package-json)
- [package-lock.json spec](https://docs.npmjs.com/cli/v10/configuring-npm/package-lock-json)
- [npm registry API](https://github.com/npm/registry/blob/main/docs/REGISTRY-API.md)
- [Node.js module resolution](https://nodejs.org/api/modules.html#all-together)
- [semver spec](https://semver.org/)
