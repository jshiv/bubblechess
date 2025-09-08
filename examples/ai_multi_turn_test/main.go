package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"chess-tui/game"
)

func main() {
	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	fmt.Println("🤖 BubbleChess AI Multi-Turn Test")
	fmt.Println("=================================")

	// Create a new game in Human vs AI mode
	chessGame := game.NewGameWithMode(game.ModeHumanVsAI)

	// Test AI connection first
	fmt.Println("\n🔌 Testing AI connection...")
	if err := chessGame.AIClient().TestConnection(); err != nil {
		log.Fatalf("AI connection test failed: %v", err)
	}
	fmt.Println("✅ AI connection successful!")

	// Get initial board state
	boardState := chessGame.GetBoardState()
	fmt.Println("\n📋 Initial board state:")
	fmt.Println(boardState)

	// Test multiple AI moves with detailed timing
	fmt.Println("\n🎮 Starting multi-turn AI test...")

	// Simulate a few human moves first to set up an interesting position
	humanMoves := []string{"e4", "Nf3", "Bc4"}
	for i, move := range humanMoves {
		fmt.Printf("\n👤 Human move %d: %s", i+1, move)
		start := time.Now()

		if err := chessGame.MakeMove(move); err != nil {
			log.Printf("❌ Human move failed: %v", err)
			continue
		}

		duration := time.Since(start)
		fmt.Printf(" (took %v)", duration)
		fmt.Printf("\n📋 Board after move %d:\n%s", i+1, chessGame.GetBoardState())
	}

	// Now test AI moves for several turns
	aiMovesToTest := 5
	fmt.Printf("\n🤖 Testing %d AI moves with detailed timing...\n", aiMovesToTest)

	for turn := 1; turn <= aiMovesToTest; turn++ {
		fmt.Print("\n" + strings.Repeat("=", 50))
		fmt.Printf("\n🎯 AI TURN %d", turn)
		fmt.Print("\n" + strings.Repeat("=", 50))

		// Get current state
		currentBoard := chessGame.GetBoardState()
		gameHistory := chessGame.GetGameHistory()
		playerColor := "black" // AI is playing black in this test

		fmt.Printf("\n📋 Board state before AI move %d:\n%s", turn, currentBoard)
		fmt.Printf("📚 Game history: %v\n", gameHistory)
		fmt.Printf("🎨 AI playing as: %s\n", playerColor)

		// Time the AI move request
		fmt.Printf("\n⏱️  Requesting AI move %d...\n", turn)
		start := time.Now()

		aiMove, err := chessGame.AIClient().GetAIMove(currentBoard, gameHistory, playerColor)
		if err != nil {
			log.Printf("❌ AI move %d failed: %v", turn, err)
			continue
		}

		aiResponseTime := time.Since(start)
		fmt.Printf("✅ AI move %d received: %s (response time: %v)\n", turn, aiMove, aiResponseTime)

		// Apply the AI move
		fmt.Printf("🔄 Applying AI move %d...\n", turn)
		applyStart := time.Now()

		if err := chessGame.MakeMove(aiMove); err != nil {
			log.Printf("❌ Failed to apply AI move %d: %v", turn, err)
			continue
		}

		applyTime := time.Since(applyStart)
		fmt.Printf("✅ AI move %d applied successfully (apply time: %v)\n", turn, applyTime)

		// Show updated state
		fmt.Printf("\n📋 Board state after AI move %d:\n%s", turn, chessGame.GetBoardState())
		fmt.Printf("📚 Updated game history: %v\n", chessGame.GetGameHistory())
		fmt.Printf("🔄 Current turn: %s\n", chessGame.GetCurrentTurn())

		// Summary for this turn
		totalTime := aiResponseTime + applyTime
		fmt.Printf("\n📊 Turn %d Summary:", turn)
		fmt.Printf("\n   🤖 AI Response Time: %v", aiResponseTime)
		fmt.Printf("\n   🔄 Apply Time: %v", applyTime)
		fmt.Printf("\n   ⏱️  Total Time: %v", totalTime)

		// Add a small delay between moves for readability
		time.Sleep(1 * time.Second)
	}

	fmt.Print("\n" + strings.Repeat("=", 60))
	fmt.Println("\n🎉 Multi-turn AI test completed!")
	fmt.Printf("📊 Final board state:\n%s", chessGame.GetBoardState())
	fmt.Printf("📚 Final game history: %v\n", chessGame.GetGameHistory())
	fmt.Println("✨ Test completed successfully!")
}
