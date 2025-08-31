package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

// Game represents the chess game TUI
type Game struct {
	chessGame   *chess.Game
	input       textinput.Model
	status      string
	err         string
	selected    string
	validMoves  []chess.Move
	gameMode    GameMode
	aiClient    *AIClient
	gameHistory []string
	isAITurn    bool
}

// NewGame creates a new chess game
func NewGame() *Game {
	return NewGameWithMode(ModeHumanVsHuman)
}

// NewGameWithMode creates a new chess game with a specific mode
func NewGameWithMode(mode GameMode) *Game {
	input := textinput.New()
	input.Placeholder = "e2e4"
	input.Focus()
	input.CharLimit = 10
	input.Width = 20

	game := &Game{
		chessGame:   chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{})),
		input:       input,
		status:      "White's turn",
		validMoves:  []chess.Move{},
		gameMode:    mode,
		gameHistory: []string{},
		isAITurn:    false,
	}

	// Initialize AI client if playing against AI
	if mode == ModeHumanVsAI {
		game.aiClient = NewAIClient("")
	}

	return game
}

// Init initializes the game
func (g *Game) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		g.input.Cursor.BlinkCmd(),
	)
}

// Update handles game updates
func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keyboard shortcuts
		switch msg.String() {
		case "q", "ctrl+c":
			return g, tea.Quit
		case "r":
			return g, g.resetGame()
		case "h":
			return g, g.showHelp()
		case "enter":
			// Only handle enter if we have input to process and it's not AI's turn
			if g.input.Value() != "" && !g.isAITurn {
				return g, g.makeMove(g.input.Value())
			}
		}
	}

	// Only update text input if it's not AI's turn
	var cmd tea.Cmd
	if !g.isAITurn {
		g.input, cmd = g.input.Update(msg)
	}
	return g, cmd
}

// View renders the game
func (g *Game) View() string {
	var sb strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		Render("♔ Chess TUI ♛")
	sb.WriteString(title + "\n\n")

	// Board
	sb.WriteString(g.renderBoard())
	sb.WriteString("\n\n")

	// Game mode
	modeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00AAFF"))
	var modeText string
	switch g.gameMode {
	case ModeHumanVsHuman:
		modeText = "Human vs Human"
	case ModeHumanVsAI:
		modeText = "Human vs AI"
	}
	sb.WriteString(modeStyle.Render("Mode: "+modeText) + "\n")

	// Status
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	sb.WriteString(statusStyle.Render(g.status) + "\n")

	// Error message
	if g.err != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		sb.WriteString(errStyle.Render("Error: "+g.err) + "\n")
	}

	// Input
	if g.isAITurn {
		sb.WriteString("\n🤖 AI is thinking...")
	} else {
		sb.WriteString("\nEnter move (e.g., e2e4): ")
		sb.WriteString(g.input.View())
	}

	// Help
	sb.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	sb.WriteString(helpStyle.Render("Commands: [q]uit, [r]eset, [h]elp"))

	return sb.String()
}

// renderBoard renders the chess board
func (g *Game) renderBoard() string {
	board := g.chessGame.Position().Board()
	var sb strings.Builder

	// File labels (a-h)
	sb.WriteString("   ")
	for file := 0; file < 8; file++ {
		sb.WriteString(fmt.Sprintf(" %c ", 'a'+file))
	}
	sb.WriteString("\n")

	// Board squares
	for rank := 7; rank >= 0; rank-- {
		// Rank label (1-8)
		sb.WriteString(fmt.Sprintf(" %d ", rank+1))

		for file := 0; file < 8; file++ {
			square := chess.Square(rank*8 + file)
			piece := board.Piece(square)

			// Determine square color
			isLight := (rank+file)%2 == 0
			var bgColor string
			if isLight {
				bgColor = "#F0D9B5" // Light square
			} else {
				bgColor = "#B58863" // Dark square
			}

			// Determine piece color
			var fgColor string
			if piece.Color() == chess.White {
				fgColor = "#FFFFFF"
			} else {
				fgColor = "#000000"
			}

			// Get piece symbol
			symbol := g.getPieceSymbol(piece)

			// Style the square
			style := lipgloss.NewStyle().
				Background(lipgloss.Color(bgColor)).
				Foreground(lipgloss.Color(fgColor)).
				Bold(true).
				Width(3).
				Align(lipgloss.Center)

			sb.WriteString(style.Render(symbol))
		}

		// Rank label (1-8)
		sb.WriteString(fmt.Sprintf(" %d ", rank+1))
		sb.WriteString("\n")
	}

	// File labels (a-h)
	sb.WriteString("   ")
	for file := 0; file < 8; file++ {
		sb.WriteString(fmt.Sprintf(" %c ", 'a'+file))
	}

	return sb.String()
}

// getPieceSymbol returns the Unicode symbol for a chess piece
func (g *Game) getPieceSymbol(piece chess.Piece) string {
	if piece == chess.NoPiece {
		return " "
	}

	symbols := map[chess.Piece]string{
		chess.WhitePawn:   "♙",
		chess.WhiteRook:   "♖",
		chess.WhiteKnight: "♘",
		chess.WhiteBishop: "♗",
		chess.WhiteQueen:  "♕",
		chess.WhiteKing:   "♔",
		chess.BlackPawn:   "♟",
		chess.BlackRook:   "♜",
		chess.BlackKnight: "♞",
		chess.BlackBishop: "♝",
		chess.BlackQueen:  "♛",
		chess.BlackKing:   "♚",
	}

	if symbol, ok := symbols[piece]; ok {
		return symbol
	}
	return "?"
}

// makeMove attempts to make a move
func (g *Game) makeMove(moveStr string) tea.Cmd {
	return func() tea.Msg {
		// Clear previous error
		g.err = ""

		// Try to make the move
		err := g.chessGame.MoveStr(moveStr)
		if err != nil {
			g.err = err.Error()
			return nil
		}

		// Add move to history
		g.gameHistory = append(g.gameHistory, moveStr)

		// Update status
		g.updateStatus()

		// Clear input
		g.input.SetValue("")

		// If playing against AI and it's now AI's turn, get AI move
		if g.gameMode == ModeHumanVsAI && g.chessGame.Position().Turn() == chess.Black {
			g.isAITurn = true
			return g.getAIMove()
		}

		return nil
	}
}

// resetGame resets the game to starting position
func (g *Game) resetGame() tea.Cmd {
	return func() tea.Msg {
		g.chessGame = chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))
		g.status = "White's turn"
		g.err = ""
		g.input.SetValue("")
		g.gameHistory = []string{}
		g.isAITurn = false
		return nil
	}
}

// showHelp shows help information
func (g *Game) showHelp() tea.Cmd {
	return func() tea.Msg {
		g.status = "Help: Use algebraic notation (e.g., e2e4, Nf3, O-O)"
		return nil
	}
}

// updateStatus updates the game status
func (g *Game) updateStatus() {
	if g.chessGame.Outcome() != chess.NoOutcome {
		switch g.chessGame.Outcome() {
		case chess.WhiteWon:
			g.status = "White wins!"
		case chess.BlackWon:
			g.status = "Black wins!"
		case chess.Draw:
			g.status = "Draw!"
		}
	} else {
		if g.chessGame.Position().Turn() == chess.White {
			g.status = "White's turn"
		} else {
			g.status = "Black's turn"
		}
	}
}

// getAIMove gets a move from the AI
func (g *Game) getAIMove() tea.Cmd {
	return func() tea.Msg {
		if g.aiClient == nil {
			g.err = "AI client not initialized"
			return nil
		}

		// Get current board state
		boardState := g.getBoardState()

		// Get AI move
		aiMove, err := g.aiClient.GetAIMove(boardState, g.gameHistory)
		if err != nil {
			g.err = "AI error: " + err.Error()
			return nil
		}

		// Apply AI move
		err = g.chessGame.MoveStr(aiMove)
		if err != nil {
			g.err = "Invalid AI move: " + err.Error()
			return nil
		}

		// Add AI move to history
		g.gameHistory = append(g.gameHistory, aiMove)

		// Update status and clear AI turn flag
		g.updateStatus()
		g.isAITurn = false

		return nil
	}
}

// getBoardState returns the current board state as a string
func (g *Game) getBoardState() string {
	return g.chessGame.Position().Board().Draw()
}
