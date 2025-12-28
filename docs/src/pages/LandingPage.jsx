import React, { useState } from "react";
import { Link } from "react-router-dom";
import {
  Terminal,
  Zap,
  Box,
  Layers,
  ArrowRight,
  Github,
  Book,
  Menu,
  X,
  BarChart3,
  Flame,
  Code2,
  Users,
  Sparkles,
} from "lucide-react";

export default function LandingPage() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const benchmarkData = [
    {
      name: "Nulang",
      requests: 142000,
      color: "bg-indigo-500",
      percentage: 100,
    },
    { name: "Bun", requests: 128000, color: "bg-orange-500", percentage: 90 },
    { name: "Go", requests: 98000, color: "bg-cyan-500", percentage: 69 },
    { name: "Node.js", requests: 42000, color: "bg-green-500", percentage: 30 },
  ];

  return (
    <div className="min-h-screen bg-slate-950 text-white font-sans selection:bg-indigo-500 selection:text-white">
      {/* Navbar */}
      <nav className="fixed w-full z-50 bg-slate-950/80 backdrop-blur-md border-b border-indigo-500/10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-tr from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Terminal className="w-5 h-5 text-white" />
              </div>
              <span className="font-bold text-xl tracking-tight">Nulang</span>
            </div>

            <div className="hidden md:block">
              <div className="ml-10 flex items-baseline space-x-8">
                <a
                  href="#features"
                  className="hover:text-indigo-400 transition-colors px-3 py-2 rounded-md text-sm font-medium"
                >
                  Recursos
                </a>
                <a
                  href="#benchmarks"
                  className="hover:text-indigo-400 transition-colors px-3 py-2 rounded-md text-sm font-medium"
                >
                  Benchmarks
                </a>
                <Link
                  to="/docs"
                  className="hover:text-indigo-400 transition-colors px-3 py-2 rounded-md text-sm font-medium"
                >
                  Documentação
                </Link>
                <a
                  href="#community"
                  className="hover:text-indigo-400 transition-colors px-3 py-2 rounded-md text-sm font-medium"
                >
                  Comunidade
                </a>
              </div>
            </div>

            <div className="hidden md:flex items-center gap-4">
              <a
                href="https://github.com/nulangdev/nulang"
                target="_blank"
                rel="noreferrer"
                className="p-2 hover:bg-slate-800 rounded-full transition-colors"
              >
                <Github className="w-5 h-5" />
              </a>
              <Link
                to="/docs"
                className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-full text-sm font-medium transition-all shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40"
              >
                Começar Agora
              </Link>
            </div>

            <div className="md:hidden">
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className="p-2 text-gray-400 hover:text-white"
              >
                {isMenuOpen ? <X /> : <Menu />}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile menu */}
        {isMenuOpen && (
          <div className="md:hidden bg-slate-900 border-b border-gray-800">
            <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
              <a
                href="#features"
                className="block px-3 py-2 rounded-md text-base font-medium hover:bg-slate-800"
              >
                Recursos
              </a>
              <a
                href="#benchmarks"
                className="block px-3 py-2 rounded-md text-base font-medium hover:bg-slate-800"
              >
                Benchmarks
              </a>
              <Link
                to="/docs"
                className="block px-3 py-2 rounded-md text-base font-medium hover:bg-slate-800"
              >
                Documentação
              </Link>
              <Link
                to="/docs"
                className="block w-full mt-4 bg-indigo-600 px-4 py-2 rounded-lg text-center"
              >
                Começar Agora
              </Link>
            </div>
          </div>
        )}
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto flex flex-col items-center text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-900/30 border border-indigo-500/30 text-indigo-300 text-sm mb-8 animate-fade-in-up">
          <span className="flex h-2 w-2 relative">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span>
          </span>
          v1.0 já está disponível
        </div>

        <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-8 bg-clip-text text-transparent bg-gradient-to-r from-white via-indigo-200 to-indigo-400">
          Construa o futuro com <br className="hidden md:block" />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-purple-500">
            Nulang
          </span>
        </h1>

        <p className="max-w-2xl text-lg md:text-xl text-slate-400 mb-10 leading-relaxed">
          Uma linguagem moderna e de alta performance desenhada para a próxima
          geração de software. Simples, expressiva e incrivelmente rápida.
        </p>

        <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
          <Link
            to="/docs"
            className="px-8 py-4 bg-white text-slate-950 rounded-full font-bold text-lg hover:bg-indigo-50 transition-all flex items-center justify-center gap-2"
          >
            Começar a Codar <ArrowRight className="w-5 h-5" />
          </Link>
          <Link
            to="/docs"
            className="px-8 py-4 bg-slate-800 text-white rounded-full font-bold text-lg hover:bg-slate-700 transition-all border border-slate-700 flex items-center justify-center gap-2"
          >
            <Book className="w-5 h-5" /> Ler a Doc
          </Link>
        </div>

        {/* Code Preview */}
        <div className="mt-20 w-full max-w-4xl relative group">
          <div className="absolute -inset-1 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl blur opacity-25 group-hover:opacity-50 transition duration-1000"></div>
          <div className="relative bg-slate-900 rounded-xl border border-slate-800 shadow-2xl overflow-hidden text-left">
            <div className="flex items-center px-4 py-3 bg-slate-900 border-b border-slate-800 gap-2">
              <div className="w-3 h-3 rounded-full bg-red-500/50"></div>
              <div className="w-3 h-3 rounded-full bg-yellow-500/50"></div>
              <div className="w-3 h-3 rounded-full bg-green-500/50"></div>
              <div className="ml-4 text-xs text-slate-500 font-mono">
                examples/server.nu
              </div>
            </div>
            <div className="p-6 overflow-x-auto">
              <pre className="font-mono text-sm leading-relaxed text-slate-200">
                <code>
                  <span className="text-purple-400">import</span>{" "}
                  <span className="text-blue-400">http</span>{" "}
                  <span className="text-purple-400">from</span>{" "}
                  <span className="text-green-400">"http"</span>;{"\n\n"}
                  <span className="text-slate-500">
                    // Servidor HTTP simples e performático
                  </span>
                  {"\n"}
                  <span className="text-purple-400">const</span>{" "}
                  <span className="text-blue-300">server</span> = http.
                  <span className="text-yellow-300">createServer</span>((req,
                  res) ={">"} {"{"}
                  {"\n"}
                  {"  "}res.<span className="text-yellow-300">writeHead</span>(
                  <span className="text-orange-400">200</span>, {"{"}{" "}
                  <span className="text-green-400">'Content-Type'</span>:{" "}
                  <span className="text-green-400">'application/json'</span>{" "}
                  {"}"});{"\n"}
                  {"  "}res.<span className="text-yellow-300">end</span>(
                  <span className="text-blue-400">JSON</span>.
                  <span className="text-yellow-300">stringify</span>({"{"}{" "}
                  {"\n"}
                  {"    "}
                  <span className="text-blue-300">message</span>:{" "}
                  <span className="text-green-400">"Olá do Nulang!"</span>,
                  {"\n"}
                  {"    "}
                  <span className="text-blue-300">timestamp</span>:{" "}
                  <span className="text-purple-400">new</span>{" "}
                  <span className="text-blue-400">Date</span>().
                  <span className="text-yellow-300">toISOString</span>(){"\n"}
                  {"  "}
                  {"}"}
                  {"}"}));{"\n"}
                  {"}"});{"\n\n"}
                  server.<span className="text-yellow-300">listen</span>(
                  <span className="text-orange-400">3000</span>, () ={">"} {"{"}
                  {"\n"}
                  {"  "}
                  <span className="text-blue-400">console</span>.
                  <span className="text-yellow-300">log</span>(
                  <span className="text-green-400">
                    "🚀 Servidor rodando em http://localhost:3000"
                  </span>
                  );{"\n"}
                  {"}"});
                </code>
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section
        id="features"
        className="py-24 bg-slate-900/50 border-y border-slate-800/50"
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-200 to-white">
              Por que Nulang?
            </h2>
            <p className="mt-4 text-slate-400">
              Tudo que você precisa para construir aplicações escaláveis.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <FeatureCard
              icon={<Zap className="w-8 h-8 text-yellow-400" />}
              title="Velocidade Relâmpago"
              description="Performance nativa com runtime otimizado. Até 3x mais rápido que Node.js em benchmarks."
            />
            <FeatureCard
              icon={<Code2 className="w-8 h-8 text-blue-400" />}
              title="Sintaxe Moderna"
              description="Compatível com JavaScript/TypeScript. Migre seus projetos sem reescrever código."
            />
            <FeatureCard
              icon={<Layers className="w-8 h-8 text-purple-400" />}
              title="Lib Padrão Rica"
              description="APIs compatíveis com Node.js: HTTP, File System, Crypto, Streams e muito mais."
            />
            <FeatureCard
              icon={<Box className="w-8 h-8 text-green-400" />}
              title="Zero Dependências"
              description="Runtime completo em um único binário. Sem node_modules pesados."
            />
            <FeatureCard
              icon={<Sparkles className="w-8 h-8 text-pink-400" />}
              title="TypeScript Nativo"
              description="Suporte nativo a TypeScript sem configuração adicional."
            />
            <FeatureCard
              icon={<Flame className="w-8 h-8 text-red-400" />}
              title="Hot Reload"
              description="Desenvolvimento rápido com hot reload integrado e watch mode."
            />
          </div>
        </div>
      </section>

      {/* Benchmarks Section */}
      <section id="benchmarks" className="py-24">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-yellow-900/30 border border-yellow-500/30 text-yellow-300 text-sm mb-4">
              <BarChart3 className="w-4 h-4" />
              Benchmarks
            </div>
            <h2 className="text-3xl md:text-5xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-200 to-white mb-4">
              Performance que Impressiona
            </h2>
            <p className="mt-4 text-slate-400 max-w-2xl mx-auto">
              Nulang supera os runtimes mais populares em requisições HTTP por
              segundo. Teste você mesmo!
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="bg-slate-900 rounded-2xl border border-slate-800 p-8">
              <h3 className="text-xl font-bold mb-6 text-slate-200">
                Requisições HTTP/s{" "}
                <span className="text-slate-500 text-sm font-normal">
                  (maior é melhor)
                </span>
              </h3>

              <div className="space-y-6">
                {benchmarkData.map((item, index) => (
                  <div key={item.name} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <div
                          className={`w-3 h-3 rounded-full ${item.color}`}
                        ></div>
                        <span className="font-semibold text-slate-200">
                          {item.name}
                        </span>
                      </div>
                      <span className="text-slate-400 font-mono text-sm">
                        {item.requests.toLocaleString()} req/s
                      </span>
                    </div>
                    <div className="relative h-8 bg-slate-800 rounded-lg overflow-hidden">
                      <div
                        className={`h-full ${item.color} transition-all duration-1000 ease-out flex items-center justify-end px-3`}
                        style={{ width: `${item.percentage}%` }}
                      >
                        {index === 0 && (
                          <span className="text-xs font-bold text-white">
                            🚀 {item.percentage}%
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="mt-8 pt-6 border-t border-slate-800">
                <p className="text-xs text-slate-500 text-center">
                  Benchmark: Servidor HTTP simples com JSON response • Apple M4
                  Pro • 10 conexões simultâneas • wrk -t12 -c400 -d30s
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Installation Section */}
      <section className="py-24 bg-slate-900/50 border-y border-slate-800/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl md:text-5xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-200 to-white mb-6">
              Comece em Segundos
            </h2>
            <p className="text-slate-400 mb-8">
              Instale Nulang e execute seu primeiro programa em menos de um
              minuto.
            </p>

            <div className="bg-slate-900 rounded-2xl border border-slate-800 p-8 text-left">
              <div className="space-y-4">
                <div>
                  <p className="text-slate-400 text-sm mb-2">
                    Instalação via script:
                  </p>
                  <div className="bg-slate-950 rounded-lg p-4 font-mono text-sm text-slate-200 border border-slate-800">
                    curl -fsSL
                    https://raw.githubusercontent.com/nulangdev/nulang/main/install.sh
                    | bash
                  </div>
                </div>

                <div>
                  <p className="text-slate-400 text-sm mb-2">
                    Ou baixe o binário:
                  </p>
                  <div className="bg-slate-950 rounded-lg p-4 font-mono text-sm text-slate-200 border border-slate-800">
                    wget https://github.com/nulangdev/nulang/releases/latest
                  </div>
                </div>

                <div>
                  <p className="text-slate-400 text-sm mb-2">
                    Execute seu código:
                  </p>
                  <div className="bg-slate-950 rounded-lg p-4 font-mono text-sm text-slate-200 border border-slate-800">
                    nulang run app.nu
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Community Section */}
      <section id="community" className="py-24">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-5xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-200 to-white mb-4">
              Junte-se à Comunidade
            </h2>
            <p className="mt-4 text-slate-400 max-w-2xl mx-auto">
              Milhares de desenvolvedores já estão construindo o futuro com
              Nulang.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            <CommunityCard
              icon={<Github className="w-8 h-8" />}
              title="GitHub"
              description="Contribua com o projeto open source"
              link="https://github.com/nulangdev/nulang"
            />
            <CommunityCard
              icon={<Users className="w-8 h-8" />}
              title="Discord"
              description="Converse com a comunidade em tempo real"
              link="#"
            />
            <CommunityCard
              icon={<Book className="w-8 h-8" />}
              title="Documentação"
              description="Aprenda tudo sobre Nulang"
              link="/docs"
              internal
            />
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-slate-800 bg-slate-950 py-12">
        <div className="max-w-7xl mx-auto px-4">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 bg-gradient-to-tr from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                  <Terminal className="w-5 h-5 text-white" />
                </div>
                <span className="font-bold text-xl">Nulang</span>
              </div>
              <p className="text-slate-500 text-sm">
                Uma linguagem moderna para a próxima geração de software.
              </p>
            </div>

            <div>
              <h4 className="font-semibold mb-4">Recursos</h4>
              <ul className="space-y-2 text-slate-400 text-sm">
                <li>
                  <Link
                    to="/docs"
                    className="hover:text-white transition-colors"
                  >
                    Documentação
                  </Link>
                </li>
                <li>
                  <a
                    href="#benchmarks"
                    className="hover:text-white transition-colors"
                  >
                    Benchmarks
                  </a>
                </li>
                <li>
                  <a
                    href="#features"
                    className="hover:text-white transition-colors"
                  >
                    Features
                  </a>
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold mb-4">Comunidade</h4>
              <ul className="space-y-2 text-slate-400 text-sm">
                <li>
                  <a
                    href="https://github.com/nulangdev/nulang"
                    className="hover:text-white transition-colors"
                  >
                    GitHub
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Discord
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Twitter
                  </a>
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-semibold mb-4">Legal</h4>
              <ul className="space-y-2 text-slate-400 text-sm">
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Licença
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Privacidade
                  </a>
                </li>
              </ul>
            </div>
          </div>

          <div className="pt-8 border-t border-slate-800 text-center text-slate-600 text-sm">
            <p>© 2025 Projeto Nulang. Código aberto e transparente.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({ icon, title, description }) {
  return (
    <div className="p-6 rounded-2xl bg-slate-900 border border-slate-800 hover:border-indigo-500/50 transition-colors group">
      <div className="mb-4 p-3 bg-slate-800 w-fit rounded-xl group-hover:scale-110 transition-transform">
        {icon}
      </div>
      <h3 className="text-xl font-bold text-white mb-2">{title}</h3>
      <p className="text-slate-400 leading-relaxed">{description}</p>
    </div>
  );
}

function CommunityCard({ icon, title, description, link, internal }) {
  const content = (
    <div className="p-8 rounded-2xl bg-slate-900 border border-slate-800 hover:border-indigo-500/50 transition-all group hover:scale-105">
      <div className="text-indigo-400 mb-4 group-hover:scale-110 transition-transform">
        {icon}
      </div>
      <h3 className="text-xl font-bold text-white mb-2">{title}</h3>
      <p className="text-slate-400">{description}</p>
      <div className="mt-4 flex items-center text-indigo-400 text-sm font-medium">
        Acessar <ArrowRight className="w-4 h-4 ml-1" />
      </div>
    </div>
  );

  if (internal) {
    return <Link to={link}>{content}</Link>;
  }

  return (
    <a href={link} target="_blank" rel="noreferrer">
      {content}
    </a>
  );
}
