import React, { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeHighlight from "rehype-highlight";
import { Menu, X, Book, ChevronRight } from "lucide-react";
import { docsList } from "../data/docsList";
import "highlight.js/styles/github-dark.css";

export default function DocsPage() {
  const { slug } = useParams();
  const [markdown, setMarkdown] = useState("");
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const currentSlug = slug || "README";

  useEffect(() => {
    // Usar import dinâmico do Vite para carregar arquivos markdown
    const loadMarkdown = async () => {
      try {
        // Importar o arquivo markdown usando import dinâmico
        const module = await import(`../content/${currentSlug}.md?raw`);
        setMarkdown(module.default);
      } catch (error) {
        console.error("Erro ao carregar markdown:", error);
        setMarkdown(
          "# Documento não encontrado\\n\\nEste documento está em desenvolvimento."
        );
      }
    };

    loadMarkdown();
  }, [currentSlug]);

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      {/* Header */}
      <header className="fixed top-0 w-full z-50 bg-slate-950/80 backdrop-blur-md border-b border-indigo-500/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-tr from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Book className="w-5 h-5 text-white" />
              </div>
              <span className="font-bold text-xl tracking-tight">
                Nulang Docs
              </span>
            </Link>

            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="lg:hidden p-2 text-gray-400 hover:text-white"
            >
              {isMenuOpen ? <X /> : <Menu />}
            </button>

            <div className="hidden lg:flex items-center gap-4">
              <Link
                to="/"
                className="text-slate-400 hover:text-white transition-colors"
              >
                Home
              </Link>
              <a
                href="https://github.com/nulangdev/nulang"
                target="_blank"
                rel="noreferrer"
                className="bg-indigo-600 hover:bg-indigo-700 px-4 py-2 rounded-full text-sm font-medium transition-all"
              >
                GitHub
              </a>
            </div>
          </div>
        </div>
      </header>

      {/* Mobile Sidebar Overlay */}
      {isMenuOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setIsMenuOpen(false)}
        />
      )}

      <div className="flex pt-16">
        {/* Sidebar */}
        <aside
          className={`fixed lg:sticky top-16 h-[calc(100vh-4rem)] w-64 bg-slate-900 border-r border-slate-800 overflow-y-auto z-40 transition-transform lg:translate-x-0 ${
            isMenuOpen ? "translate-x-0" : "-translate-x-full"
          }`}
        >
          <nav className="p-4 space-y-6">
            {docsList.map((category) => (
              <div key={category.category}>
                <h3 className="text-xs font-bold text-slate-500 uppercase tracking-wider mb-2">
                  {category.category}
                </h3>
                <ul className="space-y-1">
                  {category.items.map((item) => (
                    <li key={item.slug}>
                      <Link
                        to={`/docs/${item.slug}`}
                        onClick={() => setIsMenuOpen(false)}
                        className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm transition-colors ${
                          currentSlug === item.slug
                            ? "bg-indigo-600 text-white"
                            : "text-slate-400 hover:bg-slate-800 hover:text-white"
                        }`}
                      >
                        <ChevronRight className="w-4 h-4" />
                        {item.title}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </nav>
        </aside>

        {/* Main Content */}
        <main className="flex-1 min-w-0">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <article className="prose prose-invert prose-indigo max-w-none">
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeHighlight]}
                components={{
                  h1: ({ ...props }) => (
                    <h1
                      className="text-4xl font-bold mb-6 bg-clip-text text-transparent bg-gradient-to-r from-indigo-400 to-purple-500"
                      {...props}
                    />
                  ),
                  h2: ({ ...props }) => (
                    <h2
                      className="text-3xl font-bold mt-12 mb-4 text-indigo-300"
                      {...props}
                    />
                  ),
                  h3: ({ ...props }) => (
                    <h3
                      className="text-2xl font-bold mt-8 mb-3 text-indigo-300"
                      {...props}
                    />
                  ),
                  p: ({ ...props }) => (
                    <p
                      className="text-slate-300 leading-relaxed mb-4"
                      {...props}
                    />
                  ),
                  code: ({ inline, className, children, ...props }) => {
                    if (inline) {
                      return (
                        <code
                          className="bg-slate-800 text-indigo-300 px-1.5 py-0.5 rounded text-sm"
                          {...props}
                        >
                          {children}
                        </code>
                      );
                    }
                    return (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                  pre: ({ ...props }) => (
                    <pre
                      className="bg-slate-900 border border-slate-800 rounded-xl p-4 overflow-x-auto my-6"
                      {...props}
                    />
                  ),
                  a: ({ ...props }) => (
                    <a
                      className="text-indigo-400 hover:text-indigo-300 underline"
                      {...props}
                    />
                  ),
                  ul: ({ ...props }) => (
                    <ul
                      className="list-disc list-inside space-y-2 text-slate-300 mb-4"
                      {...props}
                    />
                  ),
                  ol: ({ ...props }) => (
                    <ol
                      className="list-decimal list-inside space-y-2 text-slate-300 mb-4"
                      {...props}
                    />
                  ),
                  blockquote: ({ ...props }) => (
                    <blockquote
                      className="border-l-4 border-indigo-500 pl-4 italic text-slate-400 my-4"
                      {...props}
                    />
                  ),
                  table: ({ ...props }) => (
                    <div className="overflow-x-auto my-6">
                      <table
                        className="min-w-full border border-slate-700 rounded-lg"
                        {...props}
                      />
                    </div>
                  ),
                  th: ({ ...props }) => (
                    <th
                      className="bg-slate-800 border border-slate-700 px-4 py-2 text-left font-semibold"
                      {...props}
                    />
                  ),
                  td: ({ ...props }) => (
                    <td
                      className="border border-slate-700 px-4 py-2"
                      {...props}
                    />
                  ),
                }}
              >
                {markdown}
              </ReactMarkdown>
            </article>
          </div>
        </main>
      </div>
    </div>
  );
}
