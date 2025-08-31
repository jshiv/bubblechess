package game

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GameMode represents the type of game
type GameMode int

const (
	ModeMenu GameMode = iota
	ModeHumanVsHuman
	ModeHumanVsAI
)

// Menu represents the game mode selection menu
type Menu struct {
	cursor int
	modes  []string
}

// NewMenu creates a new menu
func NewMenu() *Menu {
	return &Menu{
		cursor: 0,
		modes: []string{
			"Human vs Human",
			"Human vs AI",
		},
	}
}

// Init initializes the menu
func (m *Menu) Init() tea.Cmd {
	return nil
}

// Update handles menu updates
func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.modes)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				return NewGameWithMode(ModeHumanVsHuman), nil
			case 1:
				return NewGameWithMode(ModeHumanVsAI), nil
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the menu
func (m *Menu) View() string {
	var sb strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Render("♔ Chess TUI ♛")
	sb.WriteString(title + "\n\n")

	// Subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("Select Game Mode")
	sb.WriteString(subtitle + "\n\n")

	// Menu options
	for i, mode := range m.modes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle()
		if m.cursor == i {
			style = style.Foreground(lipgloss.Color("#00FF00")).Bold(true)
		} else {
			style = style.Foreground(lipgloss.Color("#888888"))
		}

		sb.WriteString(style.Render(cursor+" "+mode) + "\n")
	}

	// Instructions
	sb.WriteString("\n")
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Render("Use ↑/↓ or j/k to navigate, Enter to select, q to quit")
	sb.WriteString(instructions)

	return sb.String()
}
