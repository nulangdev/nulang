# Nulang - Documentação Oficial

Site de documentação oficial da linguagem Nulang, construído com Vite, React, Tailwind CSS e React Router.

## 🚀 Começando

### Instalação

```bash
npm install
```

### Desenvolvimento

Para iniciar o servidor de desenvolvimento:

```bash
npm run dev
```

O site estará disponível em `http://localhost:5173`

### Build de Produção

Para criar a build de produção:

```bash
npm run build
```

Para visualizar a build de produção localmente:

```bash
npm run preview
```

## 📁 Estrutura do Projeto

```
docs/
├── src/
│   ├── content/          # Arquivos markdown com a documentação
│   ├── data/             # Dados estruturados (lista de documentos, etc)
│   ├── pages/            # Páginas da aplicação
│   │   ├── LandingPage.jsx   # Página inicial
│   │   └── DocsPage.jsx      # Visualizador de documentação
│   ├── App.jsx           # Configuração de rotas
│   ├── main.jsx          # Entry point
│   └── index.css         # Estilos globais (Tailwind)
├── public/               # Arquivos estáticos
└── releases/download/                 # Build de produção (gerado)
```

## 📝 Adicionando Nova Documentação

1. **Adicione o arquivo markdown** em `src/content/`:

   ```bash
   touch src/content/meu-topico.md
   ```

2. **Atualize a lista de documentos** em `src/data/docsList.js`:

   ```javascript
   {
     category: "Categoria",
     items: [
       { title: "Meu Tópico", slug: "meu-topico" }
     ]
   }
   ```

3. A documentação estará automaticamente disponível em `/docs/meu-topico`

## 🎨 Tecnologias

- **Vite** - Build tool e dev server
- **React** - Framework UI
- **Tailwind CSS** - Framework CSS utilitário
- **React Router** - Roteamento
- **React Markdown** - Renderização de markdown
- **Remark GFM** - Suporte a GitHub Flavored Markdown
- **Rehype Highlight** - Syntax highlighting para código
- **Lucide React** - Ícones

## 🌐 Deploy

O projeto pode ser facilmente deployado em:

- **Vercel**: `vercel deploy`
- **Netlify**: Conecte o repositório e configure o build command como `npm run build` e o publish directory como `dist`
- **GitHub Pages**: Use o action do GitHub para build e deploy automático

## 📄 Licença

Este projeto é parte do projeto Nulang - Código aberto e transparente.
