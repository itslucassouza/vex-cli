package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/seunome/bcli/internal/tasks"
)

type model struct {
	tasks       *tasks.TaskService
	input       string
	output      []string
	errMsg      string
	showWelcome bool
	width       int
	height      int
}

// Cores alinhadas com o design system do React
var (
	// Cores principais
	accentNeon   = lipgloss.Color("#00f2ff")
	accentPurple = lipgloss.Color("#8a2be2")
	accentPink   = lipgloss.Color("#ff2a6d")
	textLight    = lipgloss.Color("#f0f3ff")
	textMuted    = lipgloss.Color("#a0a8d0")

	// Estilos com largura consistente
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentNeon).
			Align(lipgloss.Center).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPurple).
			Width(60). // Largura fixa para consistência
			MarginBottom(1)

	welcomeStyle = lipgloss.NewStyle().
			Foreground(accentNeon).
			Align(lipgloss.Center).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPurple).
			Width(60). // Mesma largura do título
			Margin(1, 0)

	promptStyle = lipgloss.NewStyle().
			Foreground(accentPurple).
			Bold(true).
			MarginTop(1)

	outputStyle = lipgloss.NewStyle().
			Foreground(textLight).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentPurple).
			Width(60). // Mesma largura para consistência
			Margin(1, 0)

	errorStyle = lipgloss.NewStyle().
			Foreground(accentPink).
			Bold(true).
			Padding(0, 1).
			Italic(true)

	// Estilo para lista de tarefas
	taskStyle = lipgloss.NewStyle().
			Foreground(textLight).
			MarginLeft(2)

	doneStyle = lipgloss.NewStyle().
			Foreground(accentPurple).
			Strikethrough(true).
			MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginLeft(2)
)

func initialModel() model {
	svc := tasks.NewTaskService()

	return model{
		tasks:       svc,
		input:       "",
		output:      []string{},
		errMsg:      "",
		showWelcome: true,
		width:       80,
		height:      24,
	}
}

type hideWelcomeMsg struct{}
type resizeMsg struct{ width, height int }

func (m model) Init() bubbletea.Cmd {
	return bubbletea.Tick(time.Second*2, func(t time.Time) bubbletea.Msg {
		return hideWelcomeMsg{}
	})
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {

	case hideWelcomeMsg:
		m.showWelcome = false
		m.output = []string{"Digite 'help' para ver os comandos disponíveis."}
		return m, nil

	case bubbletea.KeyMsg:
		if m.showWelcome {
			m.showWelcome = false
			m.output = []string{"Digite 'help' para ver os comandos disponíveis."}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, bubbletea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "enter":
			m = m.processCommand(m.input)
			m.input = ""
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}

	case bubbletea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Ajustar a largura dos estilos baseado na janela
		contentWidth := min(msg.Width-4, 76) // Margem de 2 de cada lado
		titleStyle = titleStyle.Width(contentWidth)
		outputStyle = outputStyle.Width(contentWidth)
		welcomeStyle = welcomeStyle.Width(contentWidth)
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.showWelcome {
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			welcomeScreenView(),
		)
	}

	// Construir a interface principal
	var sb strings.Builder

	// Título centralizado
	sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		titleStyle.Render("✨ VEX - Gerenciador de Tarefas ✨"),
	))
	sb.WriteString("\n\n")

	// Área de output
	if len(m.output) > 0 {
		outputContent := strings.Join(m.output, "\n")
		sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
			outputStyle.Render(outputContent),
		))
		sb.WriteString("\n")
	}

	// Área de input
	inputLine := promptStyle.Render("vex> ") + m.input
	sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		lipgloss.NewStyle().Width(60).Render(inputLine),
	))

	// Mensagem de erro
	if m.errMsg != "" {
		sb.WriteString("\n")
		sb.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
			errorStyle.Render("⚠️  "+m.errMsg),
		))
	}

	return sb.String()
}

func welcomeScreenView() string {
	welcomeText := `
╔══════════════════════════════════════════════════╗
║                  ██╗   ██╗███████╗██╗  ██╗       ║
║                  ██║   ██║██╔════╝╚██╗██╔╝       ║
║                  ██║   ██║█████╗   ╚███╔╝        ║
║                  ╚██╗ ██╔╝██╔══╝   ██╔██╗        ║
║                   ╚████╔╝ ███████╗██╔╝ ██╗       ║
║                    ╚═══╝  ╚══════╝╚═╝  ╚═╝       ║
╠══════════════════════════════════════════════════╣
║                                                  ║
║         Bem-vindo ao VEX Terminal!               ║
║         Seu gerenciador de tarefas moderno       ║
║                                                  ║
║         Pressione qualquer tecla para começar    ║
║                                                  ║
╚══════════════════════════════════════════════════╝
`
	return welcomeStyle.Render(welcomeText)
}

func (m model) processCommand(cmd string) model {
	m.errMsg = ""
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return m
	}

	args := strings.SplitN(cmd, " ", 2)
	command := strings.ToLower(args[0])
	param := ""
	if len(args) > 1 {
		param = args[1]
	}

	switch command {
	case "help", "?":
		m.output = []string{
			"📋 Comandos Disponíveis:",
			helpStyle.Render("list              - Listar todas as tarefas"),
			helpStyle.Render("add <tarefa>      - Adicionar nova tarefa"),
			helpStyle.Render("done <id>         - Marcar/desmarcar tarefa"),
			helpStyle.Render("remove <id>       - Remover tarefa"),
			helpStyle.Render("edit <id> <novo>  - Editar tarefa"),
			helpStyle.Render("quit / q / esc    - Sair do VEX"),
		}

	case "list", "ls", "l":
		tasks := m.tasks.ListTasks()
		if len(tasks) == 0 {
			m.output = []string{"🎉 Nenhuma tarefa encontrada! Use 'add' para criar uma."}
		} else {
			lines := []string{"📝 Suas Tarefas:"}
			for _, t := range tasks {
				check := "◻️ "
				style := taskStyle
				if t.Done {
					check = "✅"
					style = doneStyle
				}
				idStr := lipgloss.NewStyle().Foreground(accentNeon).Render(fmt.Sprintf("%2d", t.ID))
				lines = append(lines, fmt.Sprintf("%s %s %s", check, idStr, style.Render(t.Title)))
			}
			m.output = lines
		}

	case "add", "a", "new":
		if param == "" {
			m.errMsg = "Use: add <título da tarefa>"
		} else {
			task := m.tasks.AddTask(param)
			idStr := lipgloss.NewStyle().Foreground(accentNeon).Render(fmt.Sprintf("%d", task.ID))
			m.output = []string{
				"✨ Tarefa adicionada:",
				taskStyle.Render(fmt.Sprintf("   %s %s", idStr, task.Title)),
			}
		}

	case "done", "d", "check":
		id, err := strconv.Atoi(param)
		if err != nil {
			m.errMsg = "Use: done <id>"
		} else if !m.tasks.ToggleTask(id) {
			m.errMsg = "Tarefa não encontrada."
		} else {
			idStr := lipgloss.NewStyle().Foreground(accentNeon).Render(fmt.Sprintf("%d", id))
			m.output = []string{fmt.Sprintf("Tarefa %s marcada como feita/não feita.", idStr)}
		}

	case "remove", "rm", "delete":
		id, err := strconv.Atoi(param)
		if err != nil {
			m.errMsg = "Use: remove <id> (ex: remove 1)"
		} else if !m.tasks.RemoveTask(id) {
			m.errMsg = fmt.Sprintf("Tarefa #%d não encontrada", id)
		} else {
			idStr := lipgloss.NewStyle().Foreground(accentNeon).Render(fmt.Sprintf("%d", id))
			m.output = []string{fmt.Sprintf("🗑️  Tarefa %s removida", idStr)}
		}

	case "edit", "e", "update":
		parts := strings.SplitN(param, " ", 2)
		if len(parts) < 2 {
			m.errMsg = "Use: edit <id> <novo título>"
		} else {
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				m.errMsg = "ID deve ser um número"
			} else if !m.tasks.EditTask(id, parts[1]) {
				m.errMsg = fmt.Sprintf("Tarefa #%d não encontrada", id)
			} else {
				idStr := lipgloss.NewStyle().Foreground(accentNeon).Render(fmt.Sprintf("%d", id))
				m.output = []string{fmt.Sprintf("✏️  Tarefa %s atualizada", idStr)}
			}
		}

	case "quit", "q", "exit":
		fmt.Println("👋 Até logo! O VEX estará aqui quando precisar.")
		os.Exit(0)

	case "clear", "cls":
		m.output = []string{}

	default:
		m.errMsg = fmt.Sprintf("Comando '%s' não reconhecido. Digite 'help' para ajuda.", command)
	}

	return m
}

func Run() error {
	p := bubbletea.NewProgram(initialModel(), bubbletea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
