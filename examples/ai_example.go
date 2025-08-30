package main

import (
	"fmt"
	"log"
	"time"

	"chess-tui/ai_player"
)

func main() {
	fmt.Println("🤖 BubbleChess AI Player Example")
	fmt.Println("=================================")

	// Load or create configuration
	config, err := ai_player.LoadConfig("ai_config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Using Ollama URL: %s\n", config.OllamaURL)
	fmt.Printf("Using Model: %s\n", config.Model)

	// Create AI game in Human vs AI mode
	game := ai_player.NewAIGame(ai_player.ModeHumanVsAI, config)

	// Test AI connection
	fmt.Println("\n🔌 Testing AI connection...")
	if err := game.TestAIConnection(); err != nil {
		log.Fatalf("AI connection test failed: %v", err)
	}
	fmt.Println("✅ AI connection successful!")

	// Example board state (starting position)
	boardState := `  a b c d e f g h
8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 8
7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 7
6 . . . . . . . . 6
5 . . . . . . . . 5
4 . . . . . . . . 4
3 . . . . . . . . 3
2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ 2
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 1
  a b c d e f g h`

	fmt.Println("\n📋 Current board state:")
	fmt.Println(boardState)

	// Simulate a few moves
	moves := []string{"e2e4", "e7e5", "Nf3"}
	for _, move := range moves {
		game.AddMove(move)
		game.SwitchTurn()
	}

	fmt.Printf("\n📚 Game history: %v\n", game.MoveHistory)
	fmt.Printf("🔄 Current turn: %s\n", game.CurrentTurn)

	// Get AI move
	fmt.Println("\n🤖 Getting AI move...")
	start := time.Now()

	aiMove, err := game.GetAIMove(boardState)
	if err != nil {
		log.Fatalf("Failed to get AI move: %v", err)
	}

	duration := time.Since(start)
	fmt.Printf("⏱️  AI response time: %v\n", duration)
	fmt.Printf("🎯 AI suggests: %s\n", aiMove.Notation)

	// Show game status
	fmt.Println("\n📊 Game Status:")
	fmt.Println(game.GetGameStatus())

	fmt.Println("\n✨ Example completed successfully!")
}
