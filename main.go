package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	bar "github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var IncreasePercent float64
var ShowProgress = false

type Model struct {
	choices  []string       // items on the to-do list
	cursor   int            // which to-do list item our cursor is pointing at
	selected map[int]string // which to-do items are selected
	progress bar.Model
}

type tickMsg time.Time

func main() {
	model := InitialModel()

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// for {

	// 	Welcome()
	// 	scanner := bufio.NewScanner(os.Stdin)

	// 	for scanner.Scan() {
	// 		Choices()
	// 		text := strings.ToUpper(scanner.Text())
	// 		switch text {
	// 		case "1":
	// 			fmt.Println("Starting 30 min ...")
	// 			notification, err := Notify("30 min pomodoro", "Starting 30 min !! 25 minutes focused + 5 minutes break")
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			fmt.Println(notification)
	// 		case "2":
	// 			fmt.Println("Starting 45 min ...")
	// 		case "3":
	// 			fmt.Println("Starting 60 min ...")
	// 		case "Q":
	// 			fmt.Println("Quitting zeero pomodoro ...")
	// 			time.Sleep(time.Second)
	// 			os.Exit(0)

	// 		}
	// 	}
	// }
}

func Welcome() {
	fmt.Println(`
============================================
	WELCOME TO ZEERO POMODORO !!
============================================

Press enter to start
	`)
}

func Choices() []string {
	return []string{
		"(1) - Start 30 min pomodoro",
		"(2) - Start 45 min pomodoro",
		"(3) - Start 60 min pomodoro",
		"(q/Q) - Quit",
	}
}

func Notify(summary, body string) (string, error) {
	cmd := exec.Command("notify-send", "-t", "3000", summary, body)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func InitialModel() Model {
	return Model{
		choices:  Choices(),
		selected: make(map[int]string),
		progress: bar.New(bar.WithDefaultGradient()),
	}
}

func (m Model) Init() tea.Cmd {
	return TickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected[0] = "aeeee"
			m.selected[1] = "aeeee123123123"
			m.selected[2] = "ZZZZZZZZZZZZZZZZZZZ"
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
				IncreasePercent = 0.1
				Notify("Focus time Complete", "Take a break of X minutes !!")
				ShowProgress = true
			} else {
				m.selected[m.cursor] = " "
			}
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(IncreasePercent)
		return m, tea.Batch(TickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header
	s := "Please choose one of the following:\n\n"
	// Iterate over our choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	if ShowProgress {
		s += m.progress.View() + "\n\n"
	}
	// The footer
	s += helpStyle("\nPress q to quit.\n")
	// Send the UI for rendering
	return s
}

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
