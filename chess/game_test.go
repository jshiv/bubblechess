package chess

import (
	"testing"
)

func TestNewChessGame(t *testing.T) {
	game := NewChessGame()
	if game == nil {
		t.Fatal("NewChessGame returned nil")
	}
	if game.board == nil {
		t.Fatal("Game board is nil")
	}
	if !game.currentPlayer {
		t.Fatal("Game should start with white's turn")
	}
}

func TestBoardSetup(t *testing.T) {
	board := NewBoard()

	// Test pawns
	for i := 0; i < 8; i++ {
		if board.Squares[1][i] == nil || board.Squares[1][i].Type != Pawn || board.Squares[1][i].White {
			t.Errorf("Black pawn not properly set at position [1][%d]", i)
		}
		if board.Squares[6][i] == nil || board.Squares[6][i].Type != Pawn || !board.Squares[6][i].White {
			t.Errorf("White pawn not properly set at position [6][%d]", i)
		}
	}

	// Test other pieces
	expectedPieces := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, expectedType := range expectedPieces {
		if board.Squares[0][i] == nil || board.Squares[0][i].Type != expectedType || board.Squares[0][i].White {
			t.Errorf("Black %v not properly set at position [0][%d]", expectedType, i)
		}
		if board.Squares[7][i] == nil || board.Squares[7][i].Type != expectedType || !board.Squares[7][i].White {
			t.Errorf("White %v not properly set at position [7][%d]", expectedType, i)
		}
	}
}

func TestLongAlgebraicNotation(t *testing.T) {
	game := NewChessGame()

	// Test valid pawn moves
	testCases := []struct {
		move     string
		expected bool
		desc     string
	}{
		{"e2e4", true, "White pawn e2 to e4"},
		{"d7d5", false, "Black pawn d7 to d5 (not black's turn)"},
		{"e7e5", false, "Black pawn e7 to e5 (not black's turn)"},
		{"e2e5", false, "White pawn e2 to e5 (invalid move)"},
		{"e2d3", false, "White pawn e2 to d3 (invalid diagonal without capture)"},
	}

	for _, tc := range testCases {
		result := game.isValidMove(tc.move)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.desc, tc.expected, result)
		}
	}
}

func TestShortAlgebraicNotation(t *testing.T) {
	game := NewChessGame()

	// Test knight moves
	testCases := []struct {
		move     string
		expected bool
		desc     string
	}{
		{"Ng1f3", false, "Invalid format - should be Nf3"},
		{"Nf3", true, "Knight to f3"},
		{"Nc3", true, "Knight to c3"},
		{"Nf6", false, "Black knight to f6 (not black's turn)"},
		{"Nxe5", false, "Knight capture on e5 (no piece to capture)"},
	}

	for _, tc := range testCases {
		result := game.isValidMove(tc.move)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.desc, tc.expected, result)
		}
	}
}

func TestCastlingNotation(t *testing.T) {
	game := NewChessGame()

	// Test castling moves - should be invalid initially because pieces haven't moved yet
	testCases := []struct {
		move     string
		expected bool
		desc     string
	}{
		{"O-O", false, "Kingside castling (should be invalid initially)"},
		{"0-0", false, "Kingside castling (alternative notation, should be invalid initially)"},
		{"O-O-O", false, "Queenside castling (should be invalid initially)"},
		{"0-0-0", false, "Queenside castling (alternative notation, should be invalid initially)"},
	}

	for _, tc := range testCases {
		result := game.isValidMove(tc.move)
		if result != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.desc, tc.expected, result)
		}
	}

	// Now test that castling becomes valid after moving pieces
	game.executeMove("e2e4")
	game.executeMove("e7e5")
	game.executeMove("Nf3")
	game.executeMove("Nf6")
	game.executeMove("Bc4")
	game.executeMove("Bc5")

	// Castling should now be valid
	if !game.isValidMove("O-O") {
		t.Error("Castling should be valid after moving pieces")
	}
}

func TestMoveExecution(t *testing.T) {
	game := NewChessGame()

	// Execute a move
	initialPiece := game.board.Squares[6][4] // White pawn at e2
	if initialPiece == nil || initialPiece.Type != Pawn {
		t.Fatal("Expected white pawn at e2")
	}

	// Move e2e4
	game.executeMove("e2e4")

	// Check that piece moved
	if game.board.Squares[6][4] != nil {
		t.Error("Piece should have moved from e2")
	}
	if game.board.Squares[4][4] == nil || game.board.Squares[4][4].Type != Pawn {
		t.Error("Piece should be at e4")
	}

	// Check that turn switched
	if game.currentPlayer {
		t.Error("Turn should have switched to black")
	}
}

func TestShortNotationMoveExecution(t *testing.T) {
	game := NewChessGame()

	// Execute a knight move using short notation
	initialPiece := game.board.Squares[7][1] // White knight at b1
	if initialPiece == nil || initialPiece.Type != Knight {
		t.Fatal("Expected white knight at b1")
	}

	// Move Nc3
	game.executeMove("Nc3")

	// Check that piece moved
	if game.board.Squares[7][1] != nil {
		t.Error("Knight should have moved from b1")
	}
	if game.board.Squares[5][2] == nil || game.board.Squares[5][2].Type != Knight {
		t.Error("Knight should be at c3")
	}

	// Check that turn switched
	if game.currentPlayer {
		t.Error("Turn should have switched to black")
	}
}

func TestCastlingExecution(t *testing.T) {
	game := NewChessGame()

	// Move pieces to allow castling
	game.executeMove("e2e4")
	game.executeMove("e7e5")
	game.executeMove("Nf3")
	game.executeMove("Nf6")
	game.executeMove("Bc4")
	game.executeMove("Bc5")

	// Now try kingside castling
	initialKingPos := game.board.Squares[7][4] // White king at e1
	if initialKingPos == nil || initialKingPos.Type != King {
		t.Fatal("Expected white king at e1")
	}

	// Check that castling is valid
	if !game.isValidMove("O-O") {
		t.Error("Castling should be valid after moving pieces")
	}

	game.executeMove("O-O")

	// Check that king moved
	if game.board.Squares[7][4] != nil {
		t.Error("King should have moved from e1")
	}
	if game.board.Squares[7][6] == nil || game.board.Squares[7][6].Type != King {
		t.Error("King should be at g1")
	}

	// Check that rook also moved
	if game.board.Squares[7][7] != nil {
		t.Error("Rook should have moved from h1")
	}
	if game.board.Squares[7][5] == nil || game.board.Squares[7][5].Type != Rook {
		t.Error("Rook should be at f1")
	}
}

func TestDisambiguation(t *testing.T) {
	game := NewChessGame()

	// Move pieces to create a position where disambiguation is needed
	game.executeMove("e2e4")
	game.executeMove("e7e5")
	game.executeMove("Nf3")
	game.executeMove("Nf6")
	game.executeMove("Nc3")
	game.executeMove("Nc6")

	// Now we have knights at both b1 and c3, so Nbd3 should work
	game.executeMove("Nbd3")

	// Check that the knight from b1 moved to d3
	if game.board.Squares[7][1] != nil {
		t.Error("Knight should have moved from b1")
	}
	if game.board.Squares[5][3] == nil || game.board.Squares[5][3].Type != Knight {
		t.Error("Knight should be at d3")
	}
}

func TestInvalidMoves(t *testing.T) {
	game := NewChessGame()

	// Test various invalid moves
	invalidMoves := []string{
		"",        // Empty move
		"abc",     // Too short
		"abcdef",  // Too long
		"x2e4",    // Invalid file
		"e9e4",    // Invalid rank
		"e2x4",    // Invalid destination
		"N",       // Incomplete move
		"Nx",      // Incomplete capture
		"Kx",      // Incomplete capture
		"O-O-O-O", // Invalid castling
		"0-0-0-0", // Invalid castling
	}

	for _, move := range invalidMoves {
		if game.isValidMove(move) {
			t.Errorf("Move '%s' should be invalid", move)
		}
	}
}

func TestPieceFinding(t *testing.T) {
	game := NewChessGame()

	// Test finding pieces
	testCases := []struct {
		pieceType PieceType
		isWhite   bool
		expected  bool
		desc      string
	}{
		{King, true, true, "White king should be found"},
		{Queen, true, true, "White queen should be found"},
		{Rook, true, true, "White rook should be found"},
		{Bishop, true, true, "White bishop should be found"},
		{Knight, true, true, "White knight should be found"},
		{Pawn, true, true, "White pawn should be found"},
		{King, false, true, "Black king should be found"},
		{Queen, false, true, "Black queen should be found"},
	}

	for _, tc := range testCases {
		row, col, found := game.findPiece(tc.pieceType, tc.isWhite, "")
		if found != tc.expected {
			t.Errorf("%s: expected %v, got %v", tc.desc, tc.expected, found)
		}
		if found && (row < 0 || row > 7 || col < 0 || col > 7) {
			t.Errorf("%s: invalid coordinates [%d][%d]", tc.desc, row, col)
		}
	}
}

func TestGameFlow(t *testing.T) {
	game := NewChessGame()

	// Play a few moves to test game flow
	moves := []string{
		"e2e4", // White pawn e2e4
		"e7e5", // Black pawn e7e5
		"Nf3",  // White knight Nf3
		"Nc6",  // Black knight Nc6
		"Bc4",  // White bishop Bc4
		"Bc5",  // Black bishop Bc5
	}

	for i, move := range moves {
		if !game.isValidMove(move) {
			t.Errorf("Move %d '%s' should be valid", i+1, move)
			continue
		}

		game.executeMove(move)

		// Check that turn alternates
		expectedWhite := (i%2 == 0)
		if game.currentPlayer != expectedWhite {
			var expectedPlayer string
			if expectedWhite {
				expectedPlayer = "white"
			} else {
				expectedPlayer = "black"
			}
			var currentPlayer string
			if game.currentPlayer {
				currentPlayer = "white"
			} else {
				currentPlayer = "black"
			}
			t.Errorf("After move %d, expected %s's turn, got %s's turn",
				i+1, expectedPlayer, currentPlayer)
		}
	}
}

func TestMixedNotation(t *testing.T) {
	game := NewChessGame()

	// Test mixing long and short notation in the same game
	moves := []string{
		"e2e4", // Long notation
		"e7e5", // Long notation
		"Nf3",  // Short notation
		"Nc6",  // Short notation
		"d2d4", // Long notation
		"exd4", // Short notation (capture)
	}

	for i, move := range moves {
		if !game.isValidMove(move) {
			t.Errorf("Move %d '%s' should be valid", i+1, move)
			continue
		}

		game.executeMove(move)
	}

	// Verify the final position
	if game.board.Squares[4][3] == nil || game.board.Squares[4][3].Type != Pawn {
		t.Error("Black pawn should be at d4 after exd4")
	}
}
