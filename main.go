package main

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
}

func NewChessGame() *ChessGame {
	return &ChessGame{
		board:         NewBoard(),
		currentPlayer: true, // White starts
		moveInput:     textinput.NewModel(),
		status:        "White's turn",
		gameState:     gameStatePlaying,
	}
}

func (g *ChessGame) Init() tea.Cmd {
	g.moveInput.Placeholder = "Enter move (e.g. e2e4)"
	g.moveInput.Focus()
	return textinput.Blink
}

func (g *ChessGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if g.moveInput.Value() != "" {
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
			if g.selectedSquare[0] < 7 {
				g.selectedSquare[0]++
			}
		case tea.KeyDown:
			if g.selectedSquare[0] > 0 {
				g.selectedSquare[0]--
			}
		case tea.KeyRight:
			if g.selectedSquare[1] < 7 {
				g.selectedSquare[1]++
			}
		case tea.KeyLeft:
			if g.selectedSquare[1] > 0 {
				g.selectedSquare[1]--
			}
		}
	}

	var cmd tea.Cmd
	g.moveInput, cmd = g.moveInput.Update(msg)
	return g, cmd
}

func (g *ChessGame) isValidMove(move string) bool {
	// Basic validation - check if move is in format like "e2e4"
	if len(move) != 4 {
		return false
	}

	// Check if coordinates are valid
	fromCol := int(move[0] - 'a')
	fromRow := int(move[1] - '1')
	toCol := int(move[2] - 'a')
	toRow := int(move[3] - '1')

	if fromCol < 0 || fromCol > 7 || fromRow < 0 || fromRow > 7 ||
		toCol < 0 || toCol > 7 || toRow < 0 || toRow > 7 {
		return false
	}

	// Convert display row to array row (display row 1 = array row 0, display row 8 = array row 7)
	fromArrayRow := fromRow
	toArrayRow := toRow

	// Check if there's a piece at the source square
	if g.board.Squares[fromArrayRow][fromCol] == nil {
		return false
	}

	// Check if it's the right player's turn
	if g.board.Squares[fromArrayRow][fromCol].White != g.currentPlayer {
		return false
	}

	// Basic move validation (simplified)
	// In a full implementation, this would check piece-specific movement rules
	piece := g.board.Squares[fromArrayRow][fromCol]

	switch piece.Type {
	case Pawn:
		// Pawn movement - forward or diagonal capture
		if fromCol == toCol {
			// Forward movement - must be to empty square
			if g.board.Squares[toArrayRow][toCol] == nil {
				// Check if path is clear for two-square moves
				pathClear := true

				// White pawns can move two squares from starting position (row 6)
				if g.currentPlayer && fromArrayRow == 6 && toArrayRow == 4 {
					// Check if intermediate square is empty
					if g.board.Squares[5][fromCol] != nil {
						pathClear = false
					}
				}
				// Black pawns can move two squares from starting position (row 1)
				if !g.currentPlayer && fromArrayRow == 1 && toArrayRow == 3 {
					// Check if intermediate square is empty
					if g.board.Squares[2][fromCol] != nil {
						pathClear = false
					}
				}

				if pathClear {
					// Move forward one or two squares
					if (g.currentPlayer && (toArrayRow == fromArrayRow-1 || toArrayRow == fromArrayRow-2)) ||
						(!g.currentPlayer && (toArrayRow == fromArrayRow+1 || toArrayRow == fromArrayRow+2)) {
						return true
					}
				}
			}
		} else if abs(fromCol-toCol) == 1 {
			// Diagonal movement - must be capture
			if g.board.Squares[toArrayRow][toCol] != nil {
				// Must be capturing opponent's piece
				if g.board.Squares[toArrayRow][toCol].White != g.currentPlayer {
					// Move forward one square diagonally
					if (g.currentPlayer && toArrayRow == fromArrayRow-1) || (!g.currentPlayer && toArrayRow == fromArrayRow+1) {
						return true
					}
				}
			}
		}
	case King:
		// Simplified king movement (one square in any direction)
		if abs(fromArrayRow-toArrayRow) <= 1 && abs(fromCol-toCol) <= 1 {
			return true
		}
	default:
		// For other pieces, assume valid for simplicity
		return true
	}

	return false
}

func (g *ChessGame) executeMove(move string) {
	fromCol := int(move[0] - 'a')
	fromRow := int(move[1] - '1')
	toCol := int(move[2] - 'a')
	toRow := int(move[3] - '1')

	// Convert display row to array row (display row 1 = array row 0, display row 8 = array row 7)
	fromArrayRow := fromRow
	toArrayRow := toRow

	// Move piece
	g.board.Squares[toArrayRow][toCol] = g.board.Squares[fromArrayRow][fromCol]
	g.board.Squares[fromArrayRow][fromCol] = nil

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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	p := tea.NewProgram(NewChessGame())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
