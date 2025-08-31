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
	err := g.chessGame.MoveStr("e2e4")
	if err != nil {
		t.Fatalf("Failed to make move: %v", err)
	}
	g.updateStatus()
	if g.status != "Black's turn" {
		t.Errorf("Expected status 'Black's turn', got '%s'", g.status)
	}
}
