// ============================================================
// NULANG CRUD EXAMPLE - Task Manager
// Um exemplo completo de CRUD com persistência em arquivo JSON
// Demonstra: Classes, Módulos, File System, Arrays, Objects
// ============================================================

// Importações
const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

// Configurações
const DATA_FILE = path.join(__dirname, "tasks.json");

// ============================================================
// Classe Task - Representa uma tarefa
// ============================================================
class Task {
  constructor(title, description, priority) {
    this.id = crypto.randomUUID();
    this.title = title;
    this.description = description || "";
    this.priority = priority || "medium";
    this.completed = false;
    this.createdAt = Date.now();
    this.updatedAt = Date.now();
  }

  // Método para marcar como concluída
  complete() {
    this.completed = true;
    this.updatedAt = Date.now();
  }

  // Método para atualizar dados
  update(data) {
    if (data.title) {
      this.title = data.title;
    }
    if (data.description) {
      this.description = data.description;
    }
    if (data.priority) {
      this.priority = data.priority;
    }
    if (data.completed !== undefined) {
      this.completed = data.completed;
    }
    this.updatedAt = Date.now();
  }

  // Formatar para exibição
  toString() {
    let status = this.completed ? "✅" : "⏳";
    let priorityIcon = "🔵";
    if (this.priority === "high") {
      priorityIcon = "🔴";
    } else if (this.priority === "low") {
      priorityIcon = "🟢";
    }
    return (
      status +
      " " +
      priorityIcon +
      " [" +
      this.id.slice(0, 8) +
      "] " +
      this.title
    );
  }
}

// ============================================================
// Classe TaskManager - Gerencia as tarefas (CRUD)
// ============================================================
class TaskManager {
  constructor() {
    this.tasks = [];
    this.load();
  }

  // CREATE - Criar nova tarefa
  create(title: string, description: string, priority: string) {
    if (!title || title.trim() === "") {
      console.error("❌ Erro: Título é obrigatório");
      return null;
    }

    // Validar prioridade
    let validPriorities = ["low", "medium", "high"];
    if (priority && !validPriorities.includes(priority)) {
      console.error("❌ Erro: Prioridade inválida. Use: low, medium, high");
      return null;
    }

    let task = new Task(title.trim(), description, priority);
    this.tasks.push(task);
    this.save();
    console.log("✅ Tarefa criada: " + task.toString());
    return task;
  }

  // READ - Buscar tarefa por ID
  findById(id) {
    return this.tasks.find((task) => task.id === id || task.id.startsWith(id));
  }

  // READ - Listar todas as tarefas
  list(filter) {
    let result = this.tasks;

    if (filter === "completed") {
      result = this.tasks.filter((task) => task.completed === true);
    } else if (filter === "pending") {
      result = this.tasks.filter((task) => task.completed === false);
    } else if (filter === "high") {
      result = this.tasks.filter((task) => task.priority === "high");
    } else if (filter === "medium") {
      result = this.tasks.filter((task) => task.priority === "medium");
    } else if (filter === "low") {
      result = this.tasks.filter((task) => task.priority === "low");
    }

    return result;
  }

  // UPDATE - Atualizar tarefa
  update(id, data) {
    let task = this.findById(id);
    if (!task) {
      console.error("❌ Erro: Tarefa não encontrada com ID: " + id);
      return null;
    }

    task.update(data);
    this.save();
    console.log("✅ Tarefa atualizada: " + task.toString());
    return task;
  }

  // UPDATE - Marcar como concluída
  complete(id) {
    let task = this.findById(id);
    if (!task) {
      console.error("❌ Erro: Tarefa não encontrada com ID: " + id);
      return null;
    }

    task.complete();
    this.save();
    console.log("✅ Tarefa concluída: " + task.toString());
    return task;
  }

  // DELETE - Remover tarefa
  delete(id) {
    let index = this.tasks.findIndex(
      (task) => task.id === id || task.id.startsWith(id)
    );
    if (index === -1) {
      console.error("❌ Erro: Tarefa não encontrada com ID: " + id);
      return false;
    }

    let removed = this.tasks.splice(index, 1);
    this.save();
    console.log("🗑️  Tarefa removida: " + removed[0].title);
    return true;
  }

  // DELETE - Remover todas as tarefas concluídas
  clearCompleted() {
    let count = this.tasks.filter((task) => task.completed === true).length;
    this.tasks = this.tasks.filter((task) => task.completed === false);
    this.save();
    console.log("🧹 " + count + " tarefa(s) concluída(s) removida(s)");
    return count;
  }

  // Estatísticas
  stats() {
    let total = this.tasks.length;
    let completed = this.tasks.filter((t) => t.completed === true).length;
    let pending = total - completed;
    let highPriority = this.tasks.filter(
      (t) => t.priority === "high" && !t.completed
    ).length;

    return {
      total: total,
      completed: completed,
      pending: pending,
      highPriority: highPriority,
      completionRate: total > 0 ? Math.round((completed / total) * 100) : 0,
    };
  }

  // Persistência - Salvar no arquivo
  save() {
    try {
      let data = JSON.stringify(this.tasks);
      fs.writeFileSync(DATA_FILE, data);
    } catch (e) {
      console.error("❌ Erro ao salvar: " + e);
    }
  }

  // Persistência - Carregar do arquivo
  load() {
    try {
      if (fs.existsSync(DATA_FILE)) {
        let data = fs.readFileSync(DATA_FILE, "utf8");
        let parsed = JSON.parse(data);
        // Reconstrói os objetos como Tasks
        this.tasks = parsed.map((item) => {
          let task = new Task(item.title, item.description, item.priority);
          task.id = item.id;
          task.completed = item.completed;
          task.createdAt = item.createdAt;
          task.updatedAt = item.updatedAt;
          return task;
        });
        console.log(
          "📂 " + this.tasks.length + " tarefa(s) carregada(s) do arquivo"
        );
      } else {
        console.log("📂 Nenhum arquivo de dados encontrado. Iniciando vazio.");
      }
    } catch (e) {
      console.error("❌ Erro ao carregar dados: " + e);
      this.tasks = [];
    }
  }

  // Imprimir lista formatada
  print(filter) {
    let tasks = this.list(filter);

    console.log("");
    console.log("═══════════════════════════════════════════════════════");
    console.log("                    📋 LISTA DE TAREFAS                 ");
    if (filter) {
      console.log("                    Filtro: " + filter);
    }
    console.log("═══════════════════════════════════════════════════════");
    console.log("");

    if (tasks.length === 0) {
      console.log("    Nenhuma tarefa encontrada.");
    } else {
      for (let i = 0; i < tasks.length; i++) {
        let task = tasks[i];
        console.log("  " + (i + 1) + ". " + task.toString());
        if (task.description && task.description !== "") {
          console.log("      📝 " + task.description);
        }
      }
    }

    console.log("");
    console.log("═══════════════════════════════════════════════════════");

    // Mostrar estatísticas
    let s = this.stats();
    console.log(
      "  📊 Total: " +
        s.total +
        " | ✅ Concluídas: " +
        s.completed +
        " | ⏳ Pendentes: " +
        s.pending
    );
    console.log(
      "  🔴 Alta prioridade pendente: " +
        s.highPriority +
        " | 📈 Taxa de conclusão: " +
        s.completionRate +
        "%"
    );
    console.log("═══════════════════════════════════════════════════════");
    console.log("");
  }
}

// ============================================================
// Demonstração do CRUD
// ============================================================
console.log("");
console.log("🚀 NULANG CRUD EXAMPLE - Task Manager");
console.log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
console.log("");

// Criar instância do gerenciador
let manager = new TaskManager();

// ═══════════════════════════════════════════════════════
// CREATE - Criar tarefas de exemplo
// ═══════════════════════════════════════════════════════
console.log("📌 1. CREATE - Criando tarefas...");
console.log("--------------------------------");

manager.create(
  "Estudar Nulang",
  "Ler a documentação e praticar exemplos",
  "high"
);
manager.create("Fazer compras", "Ir ao supermercado comprar frutas", "medium");
manager.create("Exercícios", "Fazer 30 minutos de caminhada", "low");
manager.create("Reunião importante", "Reunião com a equipe às 14h", "high");
manager.create("Responder emails", "Verificar caixa de entrada", "medium");

console.log("");

// ═══════════════════════════════════════════════════════
// READ - Listar todas as tarefas
// ═══════════════════════════════════════════════════════
console.log("📌 2. READ - Listando todas as tarefas...");
console.log("--------------------------------");
manager.print();

// ═══════════════════════════════════════════════════════
// UPDATE - Atualizar uma tarefa
// ═══════════════════════════════════════════════════════
console.log("📌 3. UPDATE - Atualizando tarefas...");
console.log("--------------------------------");

// Buscar primeira tarefa e atualizar
let firstTask = manager.tasks[0];
if (firstTask) {
  manager.update(firstTask.id.slice(0, 8), {
    description: "Documentação atualizada com novos exemplos de CRUD",
  });
}

// Marcar algumas como concluídas
let secondTask = manager.tasks[1];
if (secondTask) {
  manager.complete(secondTask.id.slice(0, 8));
}

let thirdTask = manager.tasks[2];
if (thirdTask) {
  manager.complete(thirdTask.id.slice(0, 8));
}

console.log("");

// ═══════════════════════════════════════════════════════
// READ - Filtrar tarefas
// ═══════════════════════════════════════════════════════
console.log("📌 4. READ - Filtrando tarefas pendentes...");
console.log("--------------------------------");
manager.print("pending");

console.log("📌 5. READ - Filtrando tarefas de alta prioridade...");
console.log("--------------------------------");
manager.print("high");

// ═══════════════════════════════════════════════════════
// DELETE - Remover uma tarefa
// ═══════════════════════════════════════════════════════
console.log("📌 6. DELETE - Removendo uma tarefa...");
console.log("--------------------------------");

let lastTask = manager.tasks[manager.tasks.length - 1];
if (lastTask) {
  manager.delete(lastTask.id.slice(0, 8));
}

console.log("");

// ═══════════════════════════════════════════════════════
// Listar resultado final
// ═══════════════════════════════════════════════════════
console.log("📌 7. RESULTADO FINAL");
console.log("--------------------------------");
manager.print();

// ═══════════════════════════════════════════════════════
// Limpar tarefas concluídas
// ═══════════════════════════════════════════════════════
console.log("📌 8. DELETE - Limpando tarefas concluídas...");
console.log("--------------------------------");
manager.clearCompleted();
manager.print();

// ═══════════════════════════════════════════════════════
// Informações finais
// ═══════════════════════════════════════════════════════
console.log("✨ Demonstração completa!");
console.log("");
console.log("📁 Dados salvos em: " + DATA_FILE);
console.log("");
console.log("💡 Funcionalidades demonstradas:");
console.log("   • Classes ES6 (Task, TaskManager)");
console.log("   • Métodos de classe e instância");
console.log("   • Sistema de módulos (require)");
console.log("   • Módulo fs (readFileSync, writeFileSync)");
console.log("   • Módulo path (join)");
console.log("   • Módulo crypto (randomUUID)");
console.log("   • JSON (stringify, parse)");
console.log("   • Métodos de Array (map, filter, find, findIndex, splice)");
console.log("   • Controle de fluxo (if/else, for)");
console.log("   • Try/catch para tratamento de erros");
console.log("   • Operador ternário");
console.log("   • Template strings");
console.log("   • Variáveis __dirname, __filename");
console.log("");
