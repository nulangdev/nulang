class Nulang < Formula
  desc "Runtime JavaScript moderno compatível com Node.js, escrito em Go"
  homepage "https://github.com/nulangdev/nulang"
  version "0.1"
  
  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/nulangdev/nulang/releases/download/v#{version}/nulang-darwin-amd64"
      sha256 "a8bd05aec45d9d7187418b35d7a2ab6663358962bc2514fe3c58db29602de975"
    elsif Hardware::CPU.arm?
      url "https://github.com/nulangdev/nulang/releases/download/v#{version}/nulang-darwin-arm64"
      sha256 "163e3305b40392d96ebe48d1b9aac13bae47c339279ddb75f3078e03925eb84f"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/nulangdev/nulang/releases/download/v#{version}/nulang-linux-amd64"
      sha256 "792e2b556ac42de9aeb427cd5cd783827c96b9304849bfdb75b438acb680c20d"
    elsif Hardware::CPU.arm?
      url "https://github.com/nulangdev/nulang/releases/download/v#{version}/nulang-linux-arm64"
      sha256 "c7940575e2bfa68a3356bb36994df76366c86acbca505db5cf2dab1da3f264c6"
    end
  end

  def install
    bin.install "nulang-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}" => "nulang"
  end

  def caveats
    <<~EOS
      Nulang foi instalado com sucesso! 🚀
      
      Para começar a usar, execute:
        nulang index.nu
      
      Aliases e funções helper foram adicionados ao seu shell.
      Execute 'source ~/.zshrc' ou reinicie o terminal.
      
      Comandos disponíveis:
        nulang <arquivo.nu>  - Executar arquivo Nulang
        nu <arquivo.nu>      - Alias para nulang
      
      Para mais informações, visite:
        https://github.com/nulangdev/nulang
    EOS
  end

  def post_install
    # Configurar shell automaticamente
    zshrc = "#{Dir.home}/.zshrc"
    bashrc = "#{Dir.home}/.bashrc"
    
    shell_config = <<~'SHELL'
      
      # Nulang configuration
      export PATH="$PATH:#{HOMEBREW_PREFIX}/bin"
      
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
    SHELL
    
    # Adicionar ao .zshrc se existir e não estiver configurado
    if File.exist?(zshrc) && !File.read(zshrc).include?("# Nulang configuration")
      File.open(zshrc, "a") { |f| f.write(shell_config) }
      ohai "Configuração adicionada ao ~/.zshrc"
    end
    
    # Adicionar ao .bashrc se existir e não estiver configurado
    if File.exist?(bashrc) && !File.read(bashrc).include?("# Nulang configuration")
      File.open(bashrc, "a") { |f| f.write(shell_config) }
      ohai "Configuração adicionada ao ~/.bashrc"
    end
  end

  test do
    # Testar se o binário executa
    system "#{bin}/nulang", "--version"
  end
end
