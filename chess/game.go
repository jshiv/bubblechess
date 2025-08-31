package chess

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Game state
type gameState int

const (
	gameStatePlaying gameState = iota
	gameStateCheckmate
	gameStateStalemate
	gameStateMenu
)

// Display mode
type displayMode int

const (
	displayModeTUI displayMode = iota
	displayModeText
)

// Piece represents a chess piece
type Piece struct {
	White bool
	Type  PieceType
}

type PieceType int

const (
	Pawn PieceType = iota
	Rook
	Knight
	Bishop
	Queen
	King
)

func (p Piece) String() string {
	if p.White {
		switch p.Type {
		case Pawn:
			return "♙"
		case Rook:
			return "♖"
		case Knight:
			return "♘"
		case Bishop:
			return "♗"
		case Queen:
			return "♕"
		case King:
			return "♔"
		}
	} else {
		switch p.Type {
		case Pawn:
			return "♟"
		case Rook:
			return "♜"
		case Knight:
			return "♞"
		case Bishop:
			return "♝"
		case Queen:
			return "♛"
		case King:
			return "♚"
		}
	}
	return "?"
}

// Board represents the chess board
type Board struct {
	Squares [8][8]*Piece
	// Track if pieces have moved for castling
	WhiteKingMoved          bool
	WhiteRookKingsideMoved  bool
	WhiteRookQueensideMoved bool
	BlackKingMoved          bool
	BlackRookKingsideMoved  bool
	BlackRookQueensideMoved bool
}

func NewBoard() *Board {
	board := &Board{}
	board.setupPieces()
	return board
}

func (b *Board) setupPieces() {
	// Set up pawns
	for i := 0; i < 8; i++ {
		b.Squares[1][i] = &Piece{White: false, Type: Pawn}
		b.Squares[6][i] = &Piece{White: true, Type: Pawn}
	}

	// Set up other pieces
	pieces := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for i, pieceType := range pieces {
		b.Squares[0][i] = &Piece{White: false, Type: pieceType}
		b.Squares[7][i] = &Piece{White: true, Type: pieceType}
	}
}

func (b *Board) String() string {
	var sb strings.Builder
	sb.WriteString("  a b c d e f g h\n")
	for i := 7; i >= 0; i-- {
		sb.WriteString(fmt.Sprintf("%d ", i+1))
		for j := 0; j < 8; j++ {
			if b.Squares[i][j] == nil {
				sb.WriteString(" . ")
			} else {
				sb.WriteString(fmt.Sprintf(" %s ", b.Squares[i][j]))
			}
		}
		sb.WriteString(fmt.Sprintf(" %d\n", i+1))
	}
	sb.WriteString("  a b c d e f g h\n")
	return sb.String()
}

// ChessGame represents the game state
type ChessGame struct {
	board          *Board
	currentPlayer  bool // true for white, false for black
	selectedSquare [2]int
	moveInput      textinput.Model
	status         string
	gameState      gameState
	displayMode    int
	menuSelection  int
}

func NewChessGame() *ChessGame {
	// Check environment variable for display mode
	mode := displayModeTUI
	if os.Getenv("BUBBLECHESS_MODE") == "text" {
		mode = displayModeText
	}

	// If environment variable is set, skip menu and go directly to game
	startState := gameStateMenu
	if mode == displayModeText {
		startState = gameStatePlaying
	}

	return &ChessGame{
		board:         NewBoard(),
		currentPlayer: true, // White starts
		moveInput:     textinput.NewModel(),
		status:        "White's turn",
		gameState:     startState,
		displayMode:   int(mode),
		menuSelection: 0,
	}
}

func (g *ChessGame) Init() tea.Cmd {
	g.moveInput.Placeholder = "Enter move (e.g. e2e4, Nc6, Kxe5)"
	g.moveInput.Focus()
	return textinput.Blink
}

func (g *ChessGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if g.gameState == gameStateMenu {
				if g.menuSelection == 0 {
					g.displayMode = int(displayModeTUI)
				} else {
					g.displayMode = int(displayModeText)
				}
				g.gameState = gameStatePlaying
				g.moveInput.SetValue("")
				g.updateStatus()
			} else if g.moveInput.Value() != "" {
				move := g.moveInput.Value()
				if g.isValidMove(move) {
					g.executeMove(move)
					g.moveInput.SetValue("")
					g.updateStatus()
				} else {
					g.status = "Invalid move"
				}
			}
		case tea.KeyCtrlC:
			return g, tea.Quit
		case tea.KeyUp:
			if g.gameState == gameStateMenu {
				if g.menuSelection > 0 {
					g.menuSelection--
				}
			} else if g.selectedSquare[0] < 7 {
				g.selectedSquare[0]++
			}
		case tea.KeyDown:
			if g.gameState == gameStateMenu {
				if g.menuSelection < 1 {
					g.menuSelection++
				}
			} else if g.selectedSquare[0] > 0 {
				g.selectedSquare[0]--
			}
		case tea.KeyRight:
			if g.gameState == gameStateMenu {
				if g.menuSelection < 1 {
					g.menuSelection++
				}
			} else if g.selectedSquare[1] < 7 {
				g.selectedSquare[1]++
			}
		case tea.KeyLeft:
			if g.gameState == gameStateMenu {
				if g.menuSelection > 0 {
					g.menuSelection--
				}
			} else if g.selectedSquare[1] > 0 {
				g.selectedSquare[1]--
			}
		}
	}

	var cmd tea.Cmd
	g.moveInput, cmd = g.moveInput.Update(msg)
	return g, cmd
}

// parseMove parses both long algebraic notation (e2e4) and short algebraic notation (Nc6, Kxe5)
func (g *ChessGame) parseMove(move string) (fromRow, fromCol, toRow, toCol int, err error) {
	// Try to parse as long algebraic notation first (e2e4)
	if len(move) == 4 {
		fromCol := int(move[0] - 'a')
		fromRow := int(move[1] - '1')
		toCol := int(move[2] - 'a')
		toRow := int(move[3] - '1')

		if fromCol >= 0 && fromCol <= 7 && fromRow >= 0 && fromRow <= 7 &&
			toCol >= 0 && toCol <= 7 && toRow >= 0 && toRow <= 7 {
			return fromRow, fromCol, toRow, toCol, nil
		}
	}

	// Parse as short algebraic notation
	return g.parseShortNotation(move)
}

// parseShortNotation parses short algebraic notation like "Nc6", "Kxe5", "O-O"
func (g *ChessGame) parseShortNotation(move string) (fromRow, fromCol, toRow, toCol int, err error) {
	// Handle empty or invalid moves
	if len(move) < 2 {
		return 0, 0, 0, 0, fmt.Errorf("move too short")
	}

	// Handle castling
	if move == "O-O" || move == "0-0" {
		return g.parseCastling(move, true) // Kingside
	}
	if move == "O-O-O" || move == "0-0-0" {
		return g.parseCastling(move, false) // Queenside
	}

	// Parse piece moves like "Nc6", "Kxe5", "Nbd7"
	// Format: [Piece][File/Rank][x][Destination][=][Promotion][+][#]

	// Extract piece type
	pieceType := Pawn
	pieceChar := move[0]

	switch pieceChar {
	case 'K':
		pieceType = King
	case 'Q':
		pieceType = Queen
	case 'R':
		pieceType = Rook
	case 'B':
		pieceType = Bishop
	case 'N':
		pieceType = Knight
	default:
		// If no piece specified, it's a pawn move
		pieceType = Pawn
	}

	// Find the piece on the board
	var found bool
	fromRow, fromCol, found = g.findPiece(pieceType, g.currentPlayer, move)
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("piece not found")
	}

	// Parse destination
	destStart := len(move) - 1
	for i := len(move) - 1; i >= 0; i-- {
		if move[i] == '+' || move[i] == '#' {
			destStart = i - 1
		}
	}

	// Handle promotion
	if destStart >= 2 && move[destStart-1] == '=' {
		destStart -= 2
	}

	// Extract destination coordinates
	if destStart < 2 {
		return 0, 0, 0, 0, fmt.Errorf("invalid destination")
	}

	dest := move[destStart-1 : destStart+1]
	if len(dest) != 2 {
		return 0, 0, 0, 0, fmt.Errorf("invalid destination format")
	}

	toCol = int(dest[0] - 'a')
	toRow = int(dest[1] - '1')

	if toCol < 0 || toCol > 7 || toRow < 0 || toRow > 7 {
		return 0, 0, 0, 0, fmt.Errorf("invalid destination coordinates")
	}

	// Convert chess notation coordinates to array coordinates
	// Chess notation: a1=bottom-left, h8=top-right
	// Array coordinates: [0][0]=top-left, [7][7]=bottom-right
	toArrayRow := 7 - toRow
	toArrayCol := toCol

	return fromRow, fromCol, toArrayRow, toArrayCol, nil
}

// parseCastling handles castling moves
func (g *ChessGame) parseCastling(move string, kingside bool) (fromRow, fromCol, toRow, toCol int, err error) {
	// Find the king
	var found bool
	fromRow, fromCol, found = g.findPiece(King, g.currentPlayer, "")
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("king not found")
	}

	if kingside {
		// Kingside castling: king moves 2 squares to the right
		toRow = fromRow
		toCol = fromCol + 2
	} else {
		// Queenside castling: king moves 2 squares to the left
		toRow = fromRow
		toCol = fromCol - 2
	}

	return fromRow, fromCol, toRow, toCol, nil
}

// findPiece finds a piece of the given type and color, optionally disambiguating by file/rank
func (g *ChessGame) findPiece(pieceType PieceType, isWhite bool, move string) (row, col int, found bool) {
	var candidates [][2]int

	// Find all pieces of the given type and color
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if g.board.Squares[i][j] != nil &&
				g.board.Squares[i][j].Type == pieceType &&
				g.board.Squares[i][j].White == isWhite {
				candidates = append(candidates, [2]int{i, j})
			}
		}
	}

	if len(candidates) == 0 {
		return 0, 0, false
	}

	if len(candidates) == 1 {
		return candidates[0][0], candidates[0][1], true
	}

	// Disambiguate if multiple pieces of the same type
	// Look for file/rank disambiguation in the move string
	for _, candidate := range candidates {
		if g.canDisambiguate(candidate[0], candidate[1], move) {
			return candidate[0], candidate[1], true
		}
	}

	// If still ambiguous, return the first one (simplified)
	return candidates[0][0], candidates[0][1], true
}

// canDisambiguate checks if a piece can be disambiguated by the move notation
func (g *ChessGame) canDisambiguate(row, col int, move string) bool {
	// Check if the move string contains file disambiguation (e.g., "Nbd7" means knight on b-file)
	if len(move) >= 2 {
		file := move[1]
		if file >= 'a' && file <= 'h' {
			return int(file-'a') == col
		}
	}

	// Check if the move string contains rank disambiguation (e.g., "N1d3" means knight on rank 1)
	if len(move) >= 2 {
		rank := move[1]
		if rank >= '1' && rank <= '8' {
			return int(rank-'1') == row
		}
	}

	return true
}

func (g *ChessGame) isValidMove(move string) bool {
	// Handle castling moves first
	if move == "O-O" || move == "0-0" || move == "O-O-O" || move == "0-0-0" {
		return g.isValidCastling(move)
	}

	// Parse the move
	fromRow, fromCol, toRow, toCol, err := g.parseMove(move)
	if err != nil {
		return false
	}

	// Check if there's a piece at the source square
	if g.board.Squares[fromRow][fromCol] == nil {
		return false
	}

	// Check if it's the right player's turn
	if g.board.Squares[fromRow][fromCol].White != g.currentPlayer {
		return false
	}

	// Check if destination square is occupied by own piece
	if g.board.Squares[toRow][toCol] != nil && g.board.Squares[toRow][toCol].White == g.currentPlayer {
		return false
	}

	// Basic move validation (simplified)
	// In a full implementation, this would check piece-specific movement rules
	piece := g.board.Squares[fromRow][fromCol]

	switch piece.Type {
	case Pawn:
		// Pawn movement - forward or diagonal capture
		if fromCol == toCol {
			// Forward movement - must be to empty square
			if g.board.Squares[toRow][toCol] == nil {
				// Check if path is clear for two-square moves
				pathClear := true

				// White pawns can move two squares from starting position (row 6)
				if g.currentPlayer && fromRow == 6 && toRow == 4 {
					// Check if intermediate square is empty
					if g.board.Squares[5][fromCol] != nil {
						pathClear = false
					}
				}
				// Black pawns can move two squares from starting position (row 1)
				if !g.currentPlayer && fromRow == 1 && toRow == 3 {
					// Check if intermediate square is empty
					if g.board.Squares[2][fromCol] != nil {
						pathClear = false
					}
				}

				if pathClear {
					// Move forward one or two squares
					if (g.currentPlayer && (toRow == fromRow-1 || toRow == fromRow-2)) ||
						(!g.currentPlayer && (toRow == fromRow+1 || toRow == fromRow+2)) {
						return true
					}
				}
			}
		} else if abs(fromCol-toCol) == 1 {
			// Diagonal movement - must be capture
			if g.board.Squares[toRow][toCol] != nil {
				// Must be capturing opponent's piece
				if g.board.Squares[toRow][toCol].White != g.currentPlayer {
					// Move forward one square diagonally
					if (g.currentPlayer && toRow == fromRow-1) || (!g.currentPlayer && toRow == fromRow+1) {
						return true
					}
				}
			}
		}
	case King:
		// Simplified king movement (one square in any direction)
		if abs(fromRow-toRow) <= 1 && abs(fromCol-toCol) <= 1 {
			return true
		}
	case Knight:
		// Knight moves in L-shape: 2 squares in one direction, 1 square perpendicular
		rowDiff := abs(fromRow - toRow)
		colDiff := abs(fromCol - toCol)
		return (rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)
	case Bishop:
		// Bishop moves diagonally
		if abs(fromRow-toRow) == abs(fromCol-toCol) {
			// For now, allow diagonal moves without path checking (simplified)
			return true
		}
	case Rook:
		// Rook moves horizontally or vertically
		if fromRow == toRow || fromCol == toCol {
			// For now, allow horizontal/vertical moves without path checking (simplified)
			return true
		}
	case Queen:
		// Queen combines bishop and rook movements
		// For now, allow queen moves without path checking (simplified)
		return true
	default:
		return false
	}

	return false
}

// isValidCastling checks if castling is legal according to chess rules
func (g *ChessGame) isValidCastling(move string) bool {
	isKingside := (move == "O-O" || move == "0-0")

	if g.currentPlayer { // White's turn
		if isKingside {
			// Kingside castling: King e1->g1, Rook h1->f1
			if g.board.WhiteKingMoved || g.board.WhiteRookKingsideMoved {
				return false
			}
			// Check if squares are empty
			if g.board.Squares[7][5] != nil || g.board.Squares[7][6] != nil {
				return false
			}
			return true
		} else {
			// Queenside castling: King e1->c1, Rook a1->d1
			if g.board.WhiteKingMoved || g.board.WhiteRookQueensideMoved {
				return false
			}
			// Check if squares are empty
			if g.board.Squares[7][1] != nil || g.board.Squares[7][2] != nil || g.board.Squares[7][3] != nil {
				return false
			}
			return true
		}
	} else { // Black's turn
		if isKingside {
			// Kingside castling: King e8->g8, Rook h8->f8
			if g.board.BlackKingMoved || g.board.BlackRookKingsideMoved {
				return false
			}
			// Check if squares are empty
			if g.board.Squares[0][5] != nil || g.board.Squares[0][6] != nil {
				return false
			}
			return true
		} else {
			// Queenside castling: King e8->c8, Rook a8->d8
			if g.board.BlackKingMoved || g.board.BlackRookQueensideMoved {
				return false
			}
			// Check if squares are empty
			if g.board.Squares[0][1] != nil || g.board.Squares[0][2] != nil || g.board.Squares[0][3] != nil {
				return false
			}
			return true
		}
	}
}

// trackPieceMovement tracks when pieces move for castling purposes
func (g *ChessGame) trackPieceMovement(fromRow, fromCol, toRow, toCol int) {
	piece := g.board.Squares[toRow][toCol]
	if piece == nil {
		return
	}

	// Track king movements
	if piece.Type == King {
		if piece.White {
			g.board.WhiteKingMoved = true
		} else {
			g.board.BlackKingMoved = true
		}
	}

	// Track rook movements
	if piece.Type == Rook {
		if piece.White {
			if fromCol == 0 { // Queenside rook (a1)
				g.board.WhiteRookQueensideMoved = true
			} else if fromCol == 7 { // Kingside rook (h1)
				g.board.WhiteRookKingsideMoved = true
			}
		} else {
			if fromCol == 0 { // Queenside rook (a8)
				g.board.BlackRookQueensideMoved = true
			} else if fromCol == 7 { // Kingside rook (h8)
				g.board.BlackRookKingsideMoved = true
			}
		}
	}
}

// executeCastling handles castling moves
func (g *ChessGame) executeCastling(move string) {
	isKingside := (move == "O-O" || move == "0-0")

	if g.currentPlayer { // White's turn
		if isKingside {
			// Kingside castling: King e1->g1, Rook h1->f1
			if g.board.WhiteKingMoved || g.board.WhiteRookKingsideMoved {
				g.status = "Cannot castle: King or rook has moved"
				return
			}
			// Check if squares are empty
			if g.board.Squares[7][5] != nil || g.board.Squares[7][6] != nil {
				g.status = "Cannot castle: Squares between king and rook are occupied"
				return
			}

			// Move king
			g.board.Squares[7][6] = g.board.Squares[7][4] // e1->g1
			g.board.Squares[7][4] = nil

			// Move rook
			g.board.Squares[7][5] = g.board.Squares[7][7] // h1->f1
			g.board.Squares[7][7] = nil

			// Track the moves for castling purposes
			g.board.WhiteKingMoved = true
			g.board.WhiteRookKingsideMoved = true
		} else {
			// Queenside castling: King e1->c1, Rook a1->d1
			if g.board.WhiteKingMoved || g.board.WhiteRookQueensideMoved {
				g.status = "Cannot castle: King or rook has moved"
				return
			}
			// Check if squares are empty
			if g.board.Squares[7][1] != nil || g.board.Squares[7][2] != nil || g.board.Squares[7][3] != nil {
				g.status = "Cannot castle: Squares between king and rook are occupied"
				return
			}

			// Move king
			g.board.Squares[7][2] = g.board.Squares[7][4] // e1->c1
			g.board.Squares[7][4] = nil

			// Move rook
			g.board.Squares[7][3] = g.board.Squares[7][0] // a1->d1
			g.board.Squares[7][0] = nil

			g.board.WhiteKingMoved = true
			g.board.WhiteRookQueensideMoved = true
		}
	} else { // Black's turn
		if isKingside {
			// Kingside castling: King e8->g8, Rook h8->f8
			if g.board.BlackKingMoved || g.board.BlackRookKingsideMoved {
				g.status = "Cannot castle: King or rook has moved"
				return
			}
			// Check if squares are empty
			if g.board.Squares[0][5] != nil || g.board.Squares[0][6] != nil {
				g.status = "Cannot castle: Squares between king and rook are occupied"
				return
			}

			// Move king
			g.board.Squares[0][6] = g.board.Squares[0][4] // e8->g8
			g.board.Squares[0][4] = nil

			// Move rook
			g.board.Squares[0][5] = g.board.Squares[0][7] // h8->f8
			g.board.Squares[0][7] = nil

			g.board.BlackKingMoved = true
			g.board.BlackRookKingsideMoved = true
		} else {
			// Queenside castling: King e8->c8, Rook a8->d8
			if g.board.BlackKingMoved || g.board.BlackRookQueensideMoved {
				g.status = "Cannot castle: King or rook has moved"
				return
			}
			// Check if squares are empty
			if g.board.Squares[0][1] != nil || g.board.Squares[0][2] != nil || g.board.Squares[0][3] != nil {
				g.status = "Cannot castle: Squares between king and rook are occupied"
				return
			}

			// Move king
			g.board.Squares[0][2] = g.board.Squares[0][4] // e8->c8
			g.board.Squares[0][4] = nil

			// Move rook
			g.board.Squares[0][3] = g.board.Squares[0][0] // a8->d8
			g.board.Squares[0][0] = nil

			g.board.BlackKingMoved = true
			g.board.BlackRookQueensideMoved = true
		}
	}

	// Switch player after successful castling
	g.currentPlayer = !g.currentPlayer
}

func (g *ChessGame) executeMove(move string) {
	// Handle castling moves specially
	if move == "O-O" || move == "0-0" || move == "O-O-O" || move == "0-0-0" {
		g.executeCastling(move)
		return
	}

	fromRow, fromCol, toRow, toCol, err := g.parseMove(move)
	if err != nil {
		g.status = "Error parsing move"
		return
	}

	// Track piece movement for castling
	g.trackPieceMovement(fromRow, fromCol, toRow, toCol)

	// Move piece
	g.board.Squares[toRow][toCol] = g.board.Squares[fromRow][fromCol]
	g.board.Squares[fromRow][fromCol] = nil

	// Switch player
	g.currentPlayer = !g.currentPlayer
}

func (g *ChessGame) updateStatus() {
	if g.gameState == gameStateCheckmate {
		g.status = "Checkmate! Game over"
	} else if g.gameState == gameStateStalemate {
		g.status = "Stalemate! Game over"
	} else {
		if g.currentPlayer {
			g.status = "White's turn"
		} else {
			g.status = "Black's turn"
		}
	}
}

func (g *ChessGame) View() string {
	if g.gameState == gameStateMenu {
		return g.renderMenu()
	}

	if g.displayMode == int(displayModeText) {
		return g.renderTextMode()
	}

	return g.renderTUIMode()
}

func (g *ChessGame) renderMenu() string {
	var b strings.Builder
	b.WriteString("\n🎯 BubbleChess - Choose Display Mode\n")
	b.WriteString("=====================================\n\n")

	options := []string{"TUI Mode (Interactive)", "Text Mode (Simple Output)"}
	for i, option := range options {
		if i == g.menuSelection {
			b.WriteString("▶ ")
		} else {
			b.WriteString("  ")
		}
		b.WriteString(option)
		b.WriteString("\n")
	}

	b.WriteString("\nUse ↑/↓ to navigate, Enter to select\n")
	b.WriteString("Press Ctrl+C to exit\n")
	return b.String()
}

func (g *ChessGame) renderTextMode() string {
	var b strings.Builder

	// Game status
	b.WriteString(fmt.Sprintf("=== BubbleChess (Text Mode) ===\n"))
	b.WriteString(fmt.Sprintf("Status: %s\n", g.status))
	b.WriteString(fmt.Sprintf("Current Player: %s\n\n", g.getPlayerName()))

	// Board representation
	b.WriteString(g.board.String())
	b.WriteString("\n")

	// Move input prompt
	b.WriteString("Enter move (e.g. e2e4, Nc6, Kxe5): ")
	b.WriteString(g.moveInput.Value())

	return b.String()
}

func (g *ChessGame) renderTUIMode() string {
	var b strings.Builder

	// Board view
	b.WriteString("\n")
	for i := 7; i >= 0; i-- {
		b.WriteString(fmt.Sprintf("%d ", i+1))
		for j := 0; j < 8; j++ {
			piece := g.board.Squares[i][j]
			var squareContent string

			if piece != nil {
				squareContent = fmt.Sprintf(" %s ", piece.String())
			} else {
				squareContent = "   "
			}

			if i == g.selectedSquare[0] && j == g.selectedSquare[1] {
				b.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("200")).Render(squareContent))
			} else if (i+j)%2 == 0 {
				b.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("7")).Render(squareContent))
			} else {
				b.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("15")).Render(squareContent))
			}
		}
		b.WriteString(fmt.Sprintf(" %d\n", i+1))
	}

	// Column labels with proper spacing
	b.WriteString("   ")
	for j := 0; j < 8; j++ {
		b.WriteString(fmt.Sprintf(" %c ", 'a'+j))
	}
	b.WriteString("\n")

	// Status
	b.WriteString(fmt.Sprintf("\n%s\n", g.status))

	// Input
	b.WriteString("\n")
	b.WriteString(g.moveInput.View())

	return b.String()
}

func (g *ChessGame) getPlayerName() string {
	if g.currentPlayer {
		return "White"
	}
	return "Black"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func StartGame() error {
// 	p := tea.NewProgram(NewChessGame())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Error running program: %v\n", err)
// 		os.Exit(1)
// 	}
// }

// StartGame starts the chess game
func StartGame() error {
	p := tea.NewProgram(NewChessGame())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running chess game: %w", err)
	}
	return nil
}
