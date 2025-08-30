package main

import (
	"testing"
)

// TestNewBoard tests that a new board is set up correctly
func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// Test that pawns are in correct positions
	for i := 0; i < 8; i++ {
		// Black pawns on row 1 (array index 1)
		if board.Squares[1][i] == nil {
			t.Errorf("Expected black pawn at [1][%d], got nil", i)
		}
		if board.Squares[1][i].Type != Pawn {
			t.Errorf("Expected pawn at [1][%d], got %v", i, board.Squares[1][i].Type)
		}
		if board.Squares[1][i].White {
			t.Errorf("Expected black pawn at [1][%d], got white", i)
		}

		// White pawns on row 6 (array index 6)
		if board.Squares[6][i] == nil {
			t.Errorf("Expected white pawn at [6][%d], got nil", i)
		}
		if board.Squares[6][i].Type != Pawn {
			t.Errorf("Expected pawn at [6][%d], got %v", i, board.Squares[6][i].Type)
		}
		if !board.Squares[6][i].White {
			t.Errorf("Expected white pawn at [6][%d], got black", i)
		}
	}

	// Test that pieces are in correct positions
	pieces := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, pieceType := range pieces {
		// Black pieces on row 0 (array index 0)
		if board.Squares[0][i] == nil {
			t.Errorf("Expected black %v at [0][%d], got nil", pieceType, i)
		}
		if board.Squares[0][i].Type != pieceType {
			t.Errorf("Expected %v at [0][%d], got %v", pieceType, i, board.Squares[0][i].Type)
		}
		if board.Squares[0][i].White {
			t.Errorf("Expected black %v at [0][%d], got white", pieceType, i)
		}

		// White pieces on row 7 (array index 7)
		if board.Squares[7][i] == nil {
			t.Errorf("Expected white %v at [7][%d], got nil", pieceType, i)
		}
		if board.Squares[7][i].Type != pieceType {
			t.Errorf("Expected %v at [7][%d], got %v", pieceType, i, board.Squares[7][i].Type)
		}
		if !board.Squares[7][i].White {
			t.Errorf("Expected white %v at [7][%d], got black", pieceType, i)
		}
	}
}

// TestNewChessGame tests that a new game is initialized correctly
func TestNewChessGame(t *testing.T) {
	game := NewChessGame()

	if game.board == nil {
		t.Error("Expected board to be initialized")
	}
	if game.currentPlayer != true {
		t.Error("Expected white to start first")
	}
	if game.status != "White's turn" {
		t.Errorf("Expected status 'White's turn', got '%s'", game.status)
	}
	if game.gameState != gameStatePlaying {
		t.Error("Expected game state to be playing")
	}
}

// TestPawnMovement tests basic pawn movement rules
func TestPawnMovement(t *testing.T) {
	game := NewChessGame()

	// Test white pawn movement (e7e6)
	if !game.isValidMove("e7e6") {
		t.Error("Expected e7e6 to be valid for white")
	}

	// Test black pawn movement (e2e4) - should fail on white's turn
	if game.isValidMove("e2e4") {
		t.Error("Expected e2e4 to be invalid on white's turn")
	}

	// Execute white's move
	game.executeMove("e7e6")

	// Now it should be black's turn
	if game.currentPlayer != false {
		t.Error("Expected current player to be black after white's move")
	}

	// Test black pawn movement (e2e4) - should work now
	if !game.isValidMove("e2e4") {
		t.Error("Expected e2e4 to be valid for black")
	}

	// Test invalid pawn movement (e7e8) - can't move 2 squares from current position
	if game.isValidMove("e7e8") {
		t.Error("Expected e7e8 to be invalid (can't move 2 squares from current position)")
	}

	// Test invalid pawn movement (e7f6) - can't move diagonally without capture
	if game.isValidMove("e7f6") {
		t.Error("Expected e7f6 to be invalid (can't move diagonally without capture)")
	}
}

// TestKingMovement tests basic king movement rules
func TestKingMovement(t *testing.T) {
	game := NewChessGame()

	// Test valid king movement (e8e7)
	if !game.isValidMove("e8e7") {
		t.Error("Expected e8e7 to be valid for white king")
	}

	// Test invalid king movement (e8e6) - can't move 2 squares
	if game.isValidMove("e8e6") {
		t.Error("Expected e8e6 to be invalid (king can't move 2 squares)")
	}

	// Test diagonal king movement (e8d7)
	if !game.isValidMove("e8d7") {
		t.Error("Expected e8d7 to be valid diagonal movement for king")
	}
}

// TestTurnOrder tests that turns alternate correctly
func TestTurnOrder(t *testing.T) {
	game := NewChessGame()

	// White starts
	if game.currentPlayer != true {
		t.Error("Expected white to start")
	}

	// White moves
	game.executeMove("e7e6")
	if game.currentPlayer != false {
		t.Error("Expected black's turn after white moves")
	}

	// Black moves
	game.executeMove("e2e4")
	if game.currentPlayer != true {
		t.Error("Expected white's turn after black moves")
	}

	// White moves again
	game.executeMove("d7d6")
	if game.currentPlayer != false {
		t.Error("Expected black's turn after white moves again")
	}
}

// TestInvalidMoves tests various invalid move scenarios
func TestInvalidMoves(t *testing.T) {
	game := NewChessGame()

	// Test empty square movement
	if game.isValidMove("e5e6") {
		t.Error("Expected moving from empty square to be invalid")
	}

	// Test wrong player's piece
	if game.isValidMove("e2e3") {
		t.Error("Expected moving black piece on white's turn to be invalid")
	}

	// Test invalid coordinates
	if game.isValidMove("i9j0") {
		t.Error("Expected invalid coordinates to be invalid")
	}

	// Test wrong move format
	if game.isValidMove("e7") {
		t.Error("Expected wrong move format to be invalid")
	}
	if game.isValidMove("e7e6e5") {
		t.Error("Expected wrong move format to be invalid")
	}
}

// TestGameState tests game state transitions
func TestGameState(t *testing.T) {
	game := NewChessGame()

	// Game should start in playing state
	if game.gameState != gameStatePlaying {
		t.Error("Expected game to start in playing state")
	}

	// Status should update correctly
	if game.status != "White's turn" {
		t.Errorf("Expected status 'White's turn', got '%s'", game.status)
	}

	// After a move, status should update
	game.executeMove("e7e6")
	game.updateStatus()
	if game.status != "Black's turn" {
		t.Errorf("Expected status 'Black's turn', got '%s'", game.status)
	}
}

// TestRealisticGame tests a realistic sequence of chess moves
func TestRealisticGame(t *testing.T) {
	game := NewChessGame()

	// Test a realistic opening sequence
	moves := []string{
		"e7e6", // White: e6
		"e2e4", // Black: e4
		"d7d6", // White: d6
		"d2d4", // Black: d4
		"c7c6", // White: c6
		"c2c4", // Black: c4
	}

	for i, move := range moves {
		if !game.isValidMove(move) {
			t.Errorf("Move %d '%s' should be valid", i+1, move)
		}

		game.executeMove(move)
		game.updateStatus()

		// Verify turn alternates - after each move, the current player should be the opposite
		// White starts (true), so after move 1 (white), current player should be black (false)
		// After move 2 (black), current player should be white (true), etc.
		expectedPlayer := (i+1)%2 == 0 // false (black) after odd moves, true (white) after even moves
		if game.currentPlayer != expectedPlayer {
			t.Errorf("After move %d '%s', expected player %v, got %v",
				i+1, move, expectedPlayer, game.currentPlayer)
		}
	}
}

// TestPieceCapture tests basic capture mechanics
func TestPieceCapture(t *testing.T) {
	game := NewChessGame()

	// Set up a scenario where white can capture black pawn
	// Move white pawn to e6
	game.executeMove("e7e6")
	// Move black pawn to e4
	game.executeMove("e2e4")
	// Move white pawn to e5
	game.executeMove("e6e5")
	// Move black pawn to d4
	game.executeMove("d2d4")

	// Now white can capture black pawn at d4
	if !game.isValidMove("e5d4") {
		t.Error("Expected e5d4 to be valid capture")
	}

	// Execute the capture
	game.executeMove("e5d4")

	// Verify the piece was captured (square d4 should now have white pawn)
	if game.board.Squares[3][3] == nil {
		t.Error("Expected white pawn at d4 after capture")
	}
	if !game.board.Squares[3][3].White {
		t.Error("Expected white pawn at d4 after capture")
	}
}

// TestBoardString tests the board string representation
func TestBoardString(t *testing.T) {
	board := NewBoard()
	boardStr := board.String()

	// Should contain row and column labels
	if len(boardStr) == 0 {
		t.Error("Board string should not be empty")
	}

	// Should contain piece symbols
	if len(boardStr) < 100 {
		t.Error("Board string should be reasonably long")
	}
}

// TestPieceString tests piece string representation
func TestPieceString(t *testing.T) {
	// Test white pieces
	whitePawn := &Piece{White: true, Type: Pawn}
	if whitePawn.String() != "♙" {
		t.Errorf("Expected white pawn to render as ♙, got %s", whitePawn.String())
	}

	whiteKing := &Piece{White: true, Type: King}
	if whiteKing.String() != "♔" {
		t.Errorf("Expected white king to render as ♔, got %s", whiteKing.String())
	}

	// Test black pieces
	blackPawn := &Piece{White: false, Type: Pawn}
	if blackPawn.String() != "♟" {
		t.Errorf("Expected black pawn to render as ♟, got %s", blackPawn.String())
	}

	blackKing := &Piece{White: false, Type: King}
	if blackKing.String() != "♚" {
		t.Errorf("Expected black king to render as ♚, got %s", blackKing.String())
	}
}
