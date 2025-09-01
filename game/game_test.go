package game

import (
	"testing"

	"github.com/notnil/chess"
)

func TestNewGame(t *testing.T) {
	g := NewGame()
	if g.chessGame == nil {
		t.Error("Expected chess game to be initialized")
	}
	if g.status != "White's turn" {
		t.Errorf("Expected status 'White's turn', got '%s'", g.status)
	}
}

func TestGetPieceSymbol(t *testing.T) {
	g := NewGame()

	// Test white pieces
	if symbol := g.getPieceSymbol(chess.WhitePawn); symbol != "♙" {
		t.Errorf("Expected white pawn symbol '♙', got '%s'", symbol)
	}
	if symbol := g.getPieceSymbol(chess.WhiteKing); symbol != "♔" {
		t.Errorf("Expected white king symbol '♔', got '%s'", symbol)
	}

	// Test black pieces
	if symbol := g.getPieceSymbol(chess.BlackPawn); symbol != "♟" {
		t.Errorf("Expected black pawn symbol '♟', got '%s'", symbol)
	}
	if symbol := g.getPieceSymbol(chess.BlackQueen); symbol != "♛" {
		t.Errorf("Expected black queen symbol '♛', got '%s'", symbol)
	}

	// Test no piece
	if symbol := g.getPieceSymbol(chess.NoPiece); symbol != " " {
		t.Errorf("Expected no piece symbol ' ', got '%s'", symbol)
	}
}

func TestUpdateStatus(t *testing.T) {
	g := NewGame()

	// Test initial status
	g.updateStatus()
	if g.status != "White's turn" {
		t.Errorf("Expected status 'White's turn', got '%s'", g.status)
	}

	// Test after a move
	err := g.chessGame.MoveStr("e4")
	if err != nil {
		t.Fatalf("Failed to make move: %v", err)
	}
	g.updateStatus()
	if g.status != "Black's turn" {
		t.Errorf("Expected status 'Black's turn', got '%s'", g.status)
	}
}

func TestMoveNotationHandling(t *testing.T) {
	g := NewGame()

	// Test short algebraic notation (e4)
	err := g.chessGame.MoveStr("e4")
	if err != nil {
		t.Errorf("Failed to make short algebraic move 'e4': %v", err)
	}

	// Test that short algebraic notation (e5) is accepted
	// After e4, black can play e5 (which means e7e5)
	err = g.chessGame.MoveStr("e5")
	if err != nil {
		t.Errorf("Failed to make short algebraic move 'e5': %v", err)
	}

	// Verify the position is correct after both moves
	expectedFEN := "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
	if g.chessGame.Position().String() != expectedFEN {
		t.Errorf("Expected position %s, got %s", expectedFEN, g.chessGame.Position().String())
	}
}

func TestShortAlgebraicNotation(t *testing.T) {
	g := NewGame()

	// Make some opening moves using short algebraic notation
	// Each move must be valid for the current position
	moves := []string{"e4", "e5", "Nf3", "Nc6", "Bb5"}

	for i, move := range moves {
		err := g.chessGame.MoveStr(move)
		if err != nil {
			t.Errorf("Failed to make move %d '%s': %v", i+1, move, err)
			// Print current position for debugging
			t.Logf("Current position after move %d: %s", i, g.chessGame.Position().String())
		}
	}

	// Verify the game is still ongoing
	if g.chessGame.Outcome() != chess.NoOutcome {
		t.Errorf("Expected game to be ongoing, got outcome: %v", g.chessGame.Outcome())
	}
}

// TestLongAlgebraicNotation removed - game now uses AlgebraicNotation

func TestNotationRequirements(t *testing.T) {
	g := NewGame()

	// Test that the game now accepts short algebraic notation
	// These short notation moves should all work
	shortMoves := []string{"e4", "e5", "Nf3", "Nc6", "Bb5"}

	for i, move := range shortMoves {
		err := g.chessGame.MoveStr(move)
		if err != nil {
			t.Errorf("Failed to make short algebraic move %d '%s': %v", i+1, move, err)
		}
	}

	// Verify the final position
	expectedFEN := "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"
	if g.chessGame.Position().String() != expectedFEN {
		t.Errorf("Expected position %s, got %s", expectedFEN, g.chessGame.Position().String())
	}
}
