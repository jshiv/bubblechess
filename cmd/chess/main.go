package main

import (
	"fmt"
	"os"

	"chess-tui/ai_player"
	"chess-tui/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chess",
	Short: "A chess game with TUI and A2A server capabilities",
	Long: `Chess is a comprehensive chess application that provides:

- TUI (Terminal User Interface) for playing chess interactively
- A2A (Agent-to-Agent) server for AI-powered chess moves
- Support for both human vs human and human vs AI gameplay

The root command starts the TUI version of the game.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start the TUI chess game
		fmt.Println("Starting TUI Chess Game...")
		if err := startTUIGame(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting TUI game: %v\n", err)
			os.Exit(1)
		}
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the A2A chess server",
	Long: `Start the JSON-RPC A2A chess server that provides AI-powered chess moves.

The server listens on port 8080 and implements the A2A protocol for
agent-to-agent communication. It integrates with Ollama for AI move generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start the A2A server
		fmt.Println("Starting A2A Chess Server...")
		if err := startA2AServer(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting A2A server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add server command to root
	rootCmd.AddCommand(serverCmd)

	// Add flags for server command
	serverCmd.Flags().StringP("ollama-url", "u", "http://localhost:11434", "Ollama server URL")
	serverCmd.Flags().StringP("model", "m", "gpt-oss:20b", "Ollama model to use")
	serverCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
}

func startTUIGame() error {
	// Start the TUI chess game
	fmt.Println("Starting TUI Chess Game...")

	p := tea.NewProgram(game.NewMenu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func startA2AServer(cmd *cobra.Command) error {
	// Get flags from the command that was executed
	ollamaURL, _ := cmd.Flags().GetString("ollama-url")
	model, _ := cmd.Flags().GetString("model")
	port, _ := cmd.Flags().GetInt("port")

	fmt.Printf("Starting A2A server with:\n")
	fmt.Printf("  Ollama URL: %s\n", ollamaURL)
	fmt.Printf("  Model: %s\n", model)
	fmt.Printf("  Port: %d\n", port)

	// Start the actual A2A server
	fmt.Println("Starting A2A server...")

	// Start the JSON-RPC A2A server
	// This will block and keep the server running
	if err := ai_player.StartJSONRPCA2AServer(ollamaURL, model, port); err != nil {
		return fmt.Errorf("failed to start A2A server: %w", err)
	}

	// Note: The server is now running and this function won't return
	// until the server is stopped (e.g., by Ctrl+C)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
