package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"chess-tui/game"
)

func main() {
	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Reduced logging for cleaner output
	}))
	slog.SetDefault(logger)

	fmt.Println("⚡ BubbleChess AI Speed Test")
	fmt.Println("============================")

	// Create a new game in Human vs AI mode
	chessGame := game.NewGameWithMode(game.ModeHumanVsAI)

	// Test AI connection first
	fmt.Println("\n🔌 Testing AI connection...")
	if err := chessGame.AIClient().TestConnection(); err != nil {
		log.Fatalf("AI connection test failed: %v", err)
	}
	fmt.Println("✅ AI connection successful!")

	// Test 3 quick AI moves to measure speed improvement
	fmt.Println("\n⚡ Testing AI speed with 3 moves...")

	// Set up a simple position
	chessGame.MakeMove("e4")
	chessGame.MakeMove("e5")

	fmt.Printf("📋 Board after setup moves:\n%s\n", chessGame.GetBoardState())

	var totalTime time.Duration
	var successfulMoves int

	for i := 1; i <= 3; i++ {
		fmt.Printf("\n🎯 AI Move %d:", i)

		// Get current state
		currentBoard := chessGame.GetBoardState()
		gameHistory := chessGame.GetGameHistory()
		playerColor := "white" // AI is playing white

		// Time the AI move request
		start := time.Now()

		aiMove, err := chessGame.AIClient().GetAIMove(currentBoard, gameHistory, playerColor)
		if err != nil {
			log.Printf("❌ AI move %d failed: %v", i, err)
			continue
		}

		aiResponseTime := time.Since(start)
		fmt.Printf(" %s (%.1fs)", aiMove, aiResponseTime.Seconds())

		// Apply the AI move
		if err := chessGame.MakeMove(aiMove); err != nil {
			log.Printf("❌ Failed to apply AI move %d: %v", i, err)
			continue
		}

		totalTime += aiResponseTime
		successfulMoves++
	}

	// Results
	fmt.Printf("\n\n📊 Speed Test Results:")
	fmt.Printf("\n   ✅ Successful Moves: %d/3", successfulMoves)
	fmt.Printf("\n   ⏱️  Total Time: %.1fs", totalTime.Seconds())
	if successfulMoves > 0 {
		fmt.Printf("\n   📈 Average Time: %.1fs per move", totalTime.Seconds()/float64(successfulMoves))
	}

	if totalTime.Seconds() < 30 {
		fmt.Printf("\n   🚀 SPEED IMPROVEMENT: AI is now thinking much faster!")
	} else if totalTime.Seconds() < 60 {
		fmt.Printf("\n   ⚡ GOOD SPEED: AI is thinking reasonably fast")
	} else {
		fmt.Printf("\n   🐌 STILL SLOW: AI needs more optimization")
	}

	fmt.Printf("\n📋 Final board state:\n%s", chessGame.GetBoardState())
	fmt.Println("\n✨ Speed test completed!")
}
