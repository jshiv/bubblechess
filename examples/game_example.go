package main

import (
	"fmt"
	"log"

	"chess-tui/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Starting Chess TUI...")
	fmt.Println("Use 'q' to quit, 'r' to reset, 'h' for help")
	fmt.Println("Enter moves in long algebraic notation (e.g., e2e4)")
	fmt.Println()

	// Create and run the game
	p := tea.NewProgram(game.NewGame())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running game: %v", err)
	}
}
