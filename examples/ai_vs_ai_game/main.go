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

	fmt.Println("🤖🤖 BubbleChess AI vs AI Full Game Test")
	fmt.Println("=========================================")

	// Create a new game in Human vs AI mode (we'll manually control both sides)
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

	// Play a full game between two AIs
	totalTurns := 10
	fmt.Printf("\n🎮 Starting AI vs AI game with %d turns total...\n", totalTurns)
	fmt.Println("⚪ White AI will make moves on odd turns")
	fmt.Println("⚫ Black AI will make moves on even turns")

	var totalGameTime time.Duration
	var successfulMoves int
	var failedMoves int

	for turn := 1; turn <= totalTurns; turn++ {
		fmt.Print("\n" + strings.Repeat("=", 60))
		fmt.Printf("\n🎯 TURN %d", turn)
		fmt.Print("\n" + strings.Repeat("=", 60))

		// Determine which AI is playing
		var playerColor string
		var aiName string
		if turn%2 == 1 {
			playerColor = "white"
			aiName = "⚪ WHITE AI"
		} else {
			playerColor = "black"
			aiName = "⚫ BLACK AI"
		}

		// Get current state
		currentBoard := chessGame.GetBoardState()
		gameHistory := chessGame.GetGameHistory()
		currentTurn := chessGame.GetCurrentTurn()

		fmt.Printf("\n📋 Board state before turn %d:\n%s", turn, currentBoard)
		fmt.Printf("📚 Game history: %v\n", gameHistory)
		fmt.Printf("🔄 Current turn: %s\n", currentTurn)
		fmt.Printf("🎨 %s playing as: %s\n", aiName, playerColor)

		// Time the AI move request
		fmt.Printf("\n⏱️  Requesting move from %s...\n", aiName)
		start := time.Now()

		aiMove, err := chessGame.AIClient().GetAIMove(currentBoard, gameHistory, playerColor)
		if err != nil {
			log.Printf("❌ %s move failed: %v", aiName, err)
			failedMoves++
			continue
		}

		aiResponseTime := time.Since(start)
		fmt.Printf("✅ %s move received: %s (response time: %v)\n", aiName, aiMove, aiResponseTime)

		// Apply the AI move
		fmt.Printf("🔄 Applying %s move...\n", aiName)
		applyStart := time.Now()

		if err := chessGame.MakeMove(aiMove); err != nil {
			log.Printf("❌ Failed to apply %s move: %v", aiName, err)
			failedMoves++
			continue
		}

		applyTime := time.Since(applyStart)
		fmt.Printf("✅ %s move applied successfully (apply time: %v)\n", aiName, applyTime)

		// Show updated state
		fmt.Printf("\n📋 Board state after turn %d:\n%s", turn, chessGame.GetBoardState())
		fmt.Printf("📚 Updated game history: %v\n", chessGame.GetGameHistory())
		fmt.Printf("🔄 Current turn: %s\n", chessGame.GetCurrentTurn())

		// Summary for this turn
		totalTurnTime := aiResponseTime + applyTime
		totalGameTime += totalTurnTime
		successfulMoves++

		fmt.Printf("\n📊 Turn %d Summary:", turn)
		fmt.Printf("\n   🤖 AI Response Time: %v", aiResponseTime)
		fmt.Printf("\n   🔄 Apply Time: %v", applyTime)
		fmt.Printf("\n   ⏱️  Total Turn Time: %v", totalTurnTime)
		fmt.Printf("\n   🎯 Move: %s", aiMove)

		// Check for game end conditions
		if chessGame.GetCurrentTurn() == "Game Over" {
			fmt.Printf("\n🏁 Game ended after %d turns!", turn)
			break
		}

		// Add a small delay between moves for readability
		time.Sleep(2 * time.Second)
	}

	// Final game summary
	fmt.Print("\n" + strings.Repeat("=", 80))
	fmt.Println("\n🎉 AI vs AI Game Completed!")
	fmt.Printf("📊 Final Game Statistics:")
	fmt.Printf("\n   🎯 Total Turns: %d", totalTurns)
	fmt.Printf("\n   ✅ Successful Moves: %d", successfulMoves)
	fmt.Printf("\n   ❌ Failed Moves: %d", failedMoves)
	fmt.Printf("\n   ⏱️  Total Game Time: %v", totalGameTime)
	if successfulMoves > 0 {
		fmt.Printf("\n   📈 Average Move Time: %v", totalGameTime/time.Duration(successfulMoves))
	}
	fmt.Printf("\n📋 Final board state:\n%s", chessGame.GetBoardState())
	fmt.Printf("📚 Final game history: %v\n", chessGame.GetGameHistory())
	fmt.Println("✨ Full game test completed successfully!")
}
