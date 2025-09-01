package main

import (
	"fmt"
	"log"

	"chess-tui/game"
)

func main() {
	fmt.Println("🤖 Testing AI Client Communication")
	fmt.Println("==================================")
	fmt.Println()

	// Test 1: Create AI client
	fmt.Println("🔌 Test 1: Creating AI client...")
	aiClient := game.NewAIClient("")
	fmt.Println("✅ AI client created successfully")
	fmt.Println()

	// Test 2: Test connection to a2a server
	fmt.Println("🔌 Test 2: Testing connection to a2a server...")
	if err := aiClient.TestConnection(); err != nil {
		log.Fatalf("❌ Connection failed: %v", err)
	}
	fmt.Println("✅ Connection to a2a server successful!")
	fmt.Println()

	// Test 3: Test AI move request
	fmt.Println("🎯 Test 3: Testing AI move request...")

	// Sample board state (starting position)
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

	gameHistory := []string{"e2e4"}

	fmt.Printf("   Board state:\n%s\n", boardState)
	fmt.Printf("   Game history: %v\n", gameHistory)
	fmt.Println("   Requesting AI move...")

	aiMove, err := aiClient.GetAIMove(boardState, gameHistory, "white")
	if err != nil {
		log.Fatalf("❌ AI move request failed: %v", err)
	}

	fmt.Printf("   ✅ AI responded with move: %s\n", aiMove)
	fmt.Println()

	// Test 4: Test with different board state
	fmt.Println("🎯 Test 4: Testing AI move with different board state...")

	// Board after e4 e5
	boardState2 := `  a b c d e f g h
8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 8
7 ♟ ♟ ♟ ♟ . ♟ ♟ ♟ 7
6 . . . . . . . . 6
5 . . . . ♟ . . . 5
4 . . . . ♙ . . . 4
3 . . . . . . . . 3
2 ♙ ♙ ♙ ♙ . ♙ ♙ ♙ 2
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 1
  a b c d e f g h`

	gameHistory2 := []string{"e2e4", "e7e5"}

	fmt.Printf("   Board state:\n%s\n", boardState2)
	fmt.Printf("   Game history: %v\n", gameHistory2)
	fmt.Println("   Requesting AI move...")

	aiMove2, err := aiClient.GetAIMove(boardState2, gameHistory2, "black")
	if err != nil {
		log.Fatalf("❌ AI move request failed: %v", err)
	}

	fmt.Printf("   ✅ AI responded with move: %s\n", aiMove2)
	fmt.Println()

	fmt.Println("✨ All AI client tests completed successfully!")
	fmt.Println("🎯 The AI integration should now work in the TUI game.")
}
