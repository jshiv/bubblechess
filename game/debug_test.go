package game

import (
	"testing"

	"github.com/notnil/chess"
)

func TestMoveObjectDebugging(t *testing.T) {
	// Test with LongAlgebraicNotation to see the actual move objects
	game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	// Make first two moves
	err := game.MoveStr("e2e4")
	if err != nil {
		t.Fatalf("Failed to make first move 'e2e4': %v", err)
	}

	err = game.MoveStr("e7e5")
	if err != nil {
		t.Fatalf("Failed to make second move 'e7e5': %v", err)
	}

	// Get valid moves and show their details
	validMoves := game.ValidMoves()
	t.Logf("Valid moves at position: %s", game.Position().String())

	for i, move := range validMoves {
		// Show the move in different formats
		moveStr := move.String()
		t.Logf("Move %d: %s (from %s to %s)",
			i, moveStr, move.S1(), move.S2())

		// Try to make this move using MoveStr
		err := game.MoveStr(moveStr)
		if err != nil {
			t.Logf("  MoveStr('%s') failed: %v", moveStr, err)
		} else {
			t.Logf("  MoveStr('%s') succeeded", moveStr)
			// Can't undo, so just log success
		}
	}
}
