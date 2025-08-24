# ✨ Vex CLI - Gerenciador de Tarefas no Terminal

O **Vex CLI** é um gerenciador de tarefas interativo em modo texto (TUI), feito em Go com [Bubble Tea](https://github.com/charmbracelet/bubbletea) e [Lipgloss](https://github.com/charmbracelet/lipgloss). Ele permite gerenciar tarefas localmente e, opcionalmente, sincronizá-las com o Google Tasks.

## 🚀 Funcionalidades

- ✅ **Adicionar, listar, editar, concluir e remover tarefas**
- 🎨 Interface interativa e colorida com Bubble Tea + Lipgloss
- ☁️ Integração opcional com Google Tasks
- 💾 Persistência simples em memória (local) + suporte a salvar/editar
- 🖥️ Totalmente em linha de comando

## 📦 Instalação

### Pré-requisitos
- **Go 1.22+**
- Conta Google (opcional, apenas se quiser sincronizar com Google Tasks)

### Instalação via Go
```bash
go install github.com/itslucassouza/vex-cli@latest

## 📝 Como usar
Ao iniciar o vex-cli, você verá uma tela de boas-vindas. Pressione qualquer tecla (ou aguarde 3 segundos) para começar.

### 🔑 Comandos disponíveis
| Comando | Descrição | Exemplo |
|---------|-----------|---------|
| `help` | Mostra todos os comandos disponíveis | `help` |
| `list` | Lista todas as tarefas | `list` |
| `add <título>` | Adiciona uma nova tarefa | `add Estudar Go` |
| `done <id>` | Marca/desmarca uma tarefa como concluída | `done 2` |
| `remove <id>` | Remove uma tarefa pelo ID | `remove 3` |
| `edit <id> <novo título>` | Edita o título de uma tarefa | `edit 1 Revisar documentação` |
| `quit` / `q` / `esc` | Sai do programa | `quit` |

### 📊 Exemplo de uso
```bash
✨ Vex CLI - Todo List Interativo ✨

vex-cli> add Comprar pão
Tarefa adicionada: [ ] 1 Comprar pão

vex-cli> add Estudar Go
Tarefa adicionada: [ ] 2 Estudar Go

vex-cli> list
[ ] 1 Comprar pão
[ ] 2 Estudar Go

vex-cli> done 1
Tarefa 1 marcada como feita/não feita.

vex-cli> list
[x] 1 Comprar pão
[ ] 2 Estudar Go