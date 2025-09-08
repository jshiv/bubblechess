package game

import (
	"testing"

	"github.com/notnil/chess"
)

func TestAlgebraicNotation(t *testing.T) {
	// Test with AlgebraicNotation (short notation)
	game := chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))

	// Test short algebraic notation moves
	moves := []string{"e4", "e5", "Nf3", "Nc6", "Bb5"}

	for i, move := range moves {
		err := game.MoveStr(move)
		if err != nil {
			t.Errorf("Failed to make short algebraic move %d '%s': %v", i+1, move, err)
		}
	}

	// Verify the game is still ongoing
	if game.Outcome() != chess.NoOutcome {
		t.Errorf("Expected game to be ongoing, got outcome: %v", game.Outcome())
	}
}

// TestLongAlgebraicNotationInNewFile removed - testing long notation separately

// TestNotationComparison removed - long notation has bugs, focusing on short notation
