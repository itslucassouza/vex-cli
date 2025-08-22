//v1

// package main

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	bubbletea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"

// 	"github.com/seunome/bcli/internal/tasks"
// )

// type model struct {
// 	tasks       *tasks.TaskService
// 	input       string
// 	output      []string
// 	errMsg      string
// 	showWelcome bool
// }

// var (
// 	titleStyle = lipgloss.NewStyle().
// 			Bold(true).
// 			Foreground(lipgloss.Color("#FF79C6")).
// 			Align(lipgloss.Center).
// 			Padding(1, 4).
// 			Border(lipgloss.DoubleBorder()).
// 			BorderForeground(lipgloss.Color("#BD93F9"))

// 	welcomeStyle = lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("#8BE9FD")).
// 			Align(lipgloss.Center).
// 			Padding(1, 4).
// 			Border(lipgloss.NormalBorder()).
// 			BorderForeground(lipgloss.Color("#50FA7B"))
// 	promptStyle = lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("#50FA7B")).
// 			Bold(true)

// 	outputStyle = lipgloss.NewStyle().
// 			Foreground(lipgloss.Color("#F8F8F2")).
// 			Padding(1, 2).
// 			Border(lipgloss.NormalBorder(), true).
// 			Margin(1, 2)
// )

// func initialModel() model {
// 	svc := tasks.NewTaskService()

// 	return model{
// 		tasks:       svc,
// 		input:       "",
// 		output:      []string{},
// 		errMsg:      "",
// 		showWelcome: true,
// 	}
// }

// func (m model) Init() bubbletea.Cmd {
// 	return bubbletea.Tick(time.Second*2, func(t time.Time) bubbletea.Msg {
// 		return hideWelcomeMsg{}
// 	})
// }

// type hideWelcomeMsg struct{}

// func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
// 	switch msg := msg.(type) {

// 	case hideWelcomeMsg:
// 		return model{
// 			tasks:       m.tasks,
// 			input:       m.input,
// 			output:      []string{"Digite 'help' para ver comandos."},
// 			errMsg:      "",
// 			showWelcome: false,
// 		}, nil

// 	case bubbletea.KeyMsg:
// 		if m.showWelcome {
// 			// qualquer tecla já fecha o welcome screen
// 			m.showWelcome = false
// 			m.output = []string{"Digite 'help' para ver comandos."}
// 			return m, nil
// 		}

// 		switch msg.String() {

// 		case "ctrl+c", "esc", "q":
// 			return m, bubbletea.Quit

// 		case "backspace":
// 			if len(m.input) > 0 {
// 				m.input = m.input[:len(m.input)-1]
// 			}

// 		case "enter":
// 			m = m.processCommand(m.input)
// 			m.input = ""

// 		default:
// 			if len(msg.String()) == 1 {
// 				m.input += msg.String()
// 			}
// 		}
// 	}

// 	return m, nil
// }

// func (m model) View() string {
// 	if m.showWelcome {
// 		return welcomeScreenView()
// 	}

// 	s := titleStyle.Render("✨ bcli - Todo List Interativo ✨") + "\n\n"
// 	s += outputStyle.Render(strings.Join(m.output, "\n")) + "\n\n"
// 	s += promptStyle.Render("bcli> ") + m.input

// 	if m.errMsg != "" {
// 		s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render(m.errMsg)
// 	}

// 	return s
// }

// func welcomeScreenView() string {
// 	welcomeText := `
// ██████╗  ██████╗  ██████╗ ██╗
// ██╔══██╗██╔═══██╗██╔═══██╗██║
// ██████╔╝██║   ██║██║   ██║██║
// ██╔═══╝ ██║   ██║██║   ██║██║
// ██║     ╚██████╔╝╚██████╔╝███████╗
// ╚═╝      ╚═════╝  ╚═════╝ ╚══════╝

// Bem-vindo ao
// bcli - Seu gestor de tarefas via terminal!

// Aguarde 3 segundos ou pressione qualquer tecla para começar...
// `
// 	return welcomeStyle.Render(welcomeText)
// }

// func (m model) processCommand(cmd string) model {
// 	m.errMsg = ""
// 	cmd = strings.TrimSpace(cmd)
// 	if cmd == "" {
// 		return m
// 	}

// 	args := strings.SplitN(cmd, " ", 2)
// 	command := strings.ToLower(args[0])
// 	param := ""
// 	if len(args) > 1 {
// 		param = args[1]
// 	}

// 	switch command {
// 	case "help":
// 		m.output = []string{
// 			"Comandos:",
// 			"  list                 - listar tarefas",
// 			"  add <título>         - adicionar tarefa",
// 			"  done <id>            - marcar tarefa como feita/não feita",
// 			"  remove <id>          - remover tarefa",
// 			"  edit <id> <título>   - editar tarefa",
// 			"  quit, q, esc         - sair",
// 		}

// 	case "list":
// 		tasks := m.tasks.ListTasks()
// 		if len(tasks) == 0 {
// 			m.output = []string{"Nenhuma tarefa encontrada."}
// 		} else {
// 			lines := []string{}
// 			for _, t := range tasks {
// 				check := "[ ]"
// 				if t.Done { // aqui no seu struct TaskService é `Done` e não `Status`
// 					check = "[x]"
// 				}
// 				lines = append(lines, fmt.Sprintf("%s %d %s", check, t.ID, t.Title))
// 			}
// 			m.output = lines
// 		}

// 	case "add":
// 		if param == "" {
// 			m.errMsg = "Use: add <título>"
// 		} else {
// 			task := m.tasks.AddTask(param)

// 			check := "[ ]"
// 			if task.Done {
// 				check = "[x]"
// 			}

// 			m.output = []string{
// 				fmt.Sprintf("Tarefa adicionada: %s %d %s", check, task.ID, task.Title),
// 			}
// 		}

// 	case "done":
// 		id, err := strconv.Atoi(param)
// 		if err != nil {
// 			m.errMsg = "Use: done <id>"
// 		} else if !m.tasks.ToggleTask(id) {
// 			m.errMsg = "Tarefa não encontrada."
// 		} else {
// 			m.output = []string{fmt.Sprintf("Tarefa %d marcada como feita/não feita.", id)}
// 		}

// 	case "remove":
// 		id, err := strconv.Atoi(param)
// 		if err != nil {
// 			m.errMsg = "Use: remove <id>"
// 		} else if !m.tasks.RemoveTask(id) {
// 			m.errMsg = "Tarefa não encontrada."
// 		} else {
// 			m.output = []string{fmt.Sprintf("Tarefa %d removida.", id)}
// 		}

// 	case "edit":
// 		parts := strings.SplitN(param, " ", 2)
// 		if len(parts) < 2 {
// 			m.errMsg = "Use: edit <id> <novo título>"
// 		} else {
// 			id, err := strconv.Atoi(parts[0])
// 			if err != nil {
// 				m.errMsg = "Use: edit <id> <novo título>"
// 			} else if !m.tasks.EditTask(id, parts[1]) {
// 				m.errMsg = "Tarefa não encontrada."
// 			} else {
// 				m.output = []string{fmt.Sprintf("Tarefa %d atualizada.", id)}
// 			}
// 		}

// 	case "quit", "q", "exit":
// 		fmt.Println("Até logo!")
// 		return m

// 	default:
// 		m.errMsg = "Comando não reconhecido. Digite 'help' para ajuda."
// 	}

// 	return m
// }

// func main() {
// 	p := bubbletea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Erro:", err)
// 		os.Exit(1)
// 	}
// }