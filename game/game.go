package game

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

// Game represents the chess game TUI
type Game struct {
	chessGame     *chess.Game
	input         textinput.Model
	status        string
	err           string
	selected      string
	validMoves    []chess.Move
	gameMode      GameMode
	aiClient      *AIClient
	gameHistory   []string
	isAITurn      bool
	aiMovePending bool
}

// aiMoveRequestedMsg is a message that signals the AI move should be requested
type aiMoveRequestedMsg struct{}

// aiMoveCompletedMsg is a message that signals the AI move has been completed
type aiMoveCompletedMsg struct{}

// NewGame creates a new chess game
func NewGame() *Game {
	return NewGameWithMode(ModeHumanVsHuman)
}

// NewGameWithMode creates a new chess game with a specific mode
func NewGameWithMode(mode GameMode) *Game {
	input := textinput.New()
	input.Placeholder = "e4"
	input.Focus()
	input.CharLimit = 10
	input.Width = 20

	game := &Game{
		chessGame:     chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{})),
		input:         input,
		status:        "White's turn",
		validMoves:    []chess.Move{},
		gameMode:      mode,
		gameHistory:   []string{},
		isAITurn:      false,
		aiMovePending: false,
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
				slog.Debug("Enter pressed", "input_value", g.input.Value())
				return g, g.makeMove(g.input.Value())
			}
		}
	case aiMoveRequestedMsg:
		// AI move was requested, execute it
		slog.Debug("Received aiMoveRequestedMsg, executing getAIMove")
		return g, g.getAIMove()
	case aiMoveCompletedMsg:
		// AI move completed, refresh the TUI
		slog.Debug("Received aiMoveCompletedMsg, refreshing TUI")
		return g, nil
	default:
		// Check if AI move is pending
		if g.aiMovePending {
			slog.Debug("AI move pending, executing getAIMove")
			g.aiMovePending = false
			return g, g.getAIMove()
		}
	}

	// Only update text input if it's not AI's turn
	var cmd tea.Cmd
	if !g.isAITurn {
		slog.Debug("Updating text input", "isAITurn", g.isAITurn)
		g.input, cmd = g.input.Update(msg)
	} else {
		slog.Debug("Skipping text input update", "isAITurn", g.isAITurn)
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

	// Debug info
	slog.Debug("Game state", "gameMode", g.gameMode, "isAITurn", g.isAITurn, "turn", g.chessGame.Position().Turn())
	sb.WriteString(fmt.Sprintf("DEBUG: gameMode=%d, isAITurn=%t, turn=%s\n",
		g.gameMode, g.isAITurn, g.chessGame.Position().Turn()))

	// Additional debug info
	slog.Debug("View function state", "status", g.status, "err", g.err, "input_focused", !g.isAITurn)
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
		sb.WriteString("\nEnter move (e.g., e4): ")
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

// convertLongToShortNotation converts long algebraic notation to short algebraic notation
func (g *Game) convertLongToShortNotation(moveStr string) string {
	// If it's already short notation (less than 4 characters), return as is
	if len(moveStr) < 4 {
		return moveStr
	}

	// For pawn moves like "e2e4" -> "e4"
	if len(moveStr) == 4 && moveStr[0] >= 'a' && moveStr[0] <= 'h' &&
		moveStr[2] >= 'a' && moveStr[2] <= 'h' &&
		moveStr[1] >= '2' && moveStr[1] <= '7' &&
		moveStr[3] >= '2' && moveStr[3] <= '8' {
		return string(moveStr[2:4]) // Return destination square
	}

	// For other moves, return as is for now
	// TODO: Add more conversion logic for pieces, captures, etc.
	return moveStr
}

// makeMove attempts to make a move
func (g *Game) makeMove(moveStr string) tea.Cmd {
	return func() tea.Msg {
		slog.Debug("makeMove function started", "move", moveStr)

		// Clear previous error
		g.err = ""

		// Try to make the move
		err := g.chessGame.MoveStr(moveStr)
		if err != nil {
			slog.Debug("Move failed", "error", err)
			g.err = err.Error()
			return nil
		}
		slog.Debug("Move successful", "current_turn", g.chessGame.Position().Turn())

		// Add move to history
		g.gameHistory = append(g.gameHistory, moveStr)
		slog.Debug("Move added to history", "history_length", len(g.gameHistory))

		// Update status
		g.updateStatus()
		slog.Debug("Status updated", "new_status", g.status)

		// Clear input
		g.input.SetValue("")

		// If playing against AI and it's now AI's turn, get AI move
		slog.Debug("Checking AI turn", "gameMode", g.gameMode, "turn", g.chessGame.Position().Turn())
		if g.gameMode == ModeHumanVsAI {
			// In Human vs AI mode, after the human makes a move, it's the AI's turn to respond
			// The AI will play as the opposite color of the current turn
			slog.Debug("AI turn detected, setting aiMovePending flag")
			g.isAITurn = true
			g.aiMovePending = true
			g.status = "🤖 AI is thinking..."
			slog.Debug("aiMovePending set to true")
		} else {
			slog.Debug("Not AI turn", "gameMode", g.gameMode, "turn", g.chessGame.Position().Turn())
		}

		slog.Debug("makeMove returning nil")
		return nil
	}
}

// resetGame resets the game to starting position
func (g *Game) resetGame() tea.Cmd {
	return func() tea.Msg {
		g.chessGame = chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))
		g.status = "White's turn"
		g.err = ""
		g.input.SetValue("")
		g.gameHistory = []string{}
		g.isAITurn = false
		g.aiMovePending = false
		return nil
	}
}

// showHelp shows help information
func (g *Game) showHelp() tea.Cmd {
	return func() tea.Msg {
		g.status = "Help: Use algebraic notation (e.g., e4, Nf3, O-O)"
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
		slog.Debug("getAIMove function called")

		if g.aiClient == nil {
			slog.Debug("AI client is nil")
			g.err = "AI client not initialized"
			return nil
		}

		slog.Debug("AI client found, getting board state")
		// Get current board state
		boardState := g.getBoardState()

		slog.Debug("Board state", "board", boardState)
		slog.Debug("Game history", "history", g.gameHistory)
		slog.Debug("Calling AI client GetAIMove")

		// Get AI move (this will block, but it's the only way to make it work)
		playerColor := "white"
		if g.chessGame.Position().Turn() == chess.Black {
			playerColor = "black"
		}
		aiMove, err := g.aiClient.GetAIMove(boardState, g.gameHistory, playerColor)
		if err != nil {
			slog.Debug("AI error", "error", err)
			g.err = "AI error: " + err.Error()
			return nil
		}

		slog.Debug("AI move received", "move", aiMove)

		// Convert AI move from long to short notation if needed
		convertedMove := g.convertLongToShortNotation(aiMove)
		slog.Debug("Converted AI move", "original", aiMove, "converted", convertedMove)

		// Apply AI move
		err = g.chessGame.MoveStr(convertedMove)
		if err != nil {
			slog.Debug("Invalid AI move error", "error", err)
			g.err = "Invalid AI move: " + err.Error()

			// Send error back to AI server and request a new move
			slog.Debug("Sending error to AI server and requesting new move")
			newMove, retryErr := g.retryAIMoveWithError(boardState, g.gameHistory, err.Error(), playerColor)
			if retryErr != nil {
				slog.Debug("Retry failed", "error", retryErr)
				return nil
			}

			// Convert the retry move as well
			convertedRetryMove := g.convertLongToShortNotation(newMove)
			slog.Debug("Converted retry move", "original", newMove, "converted", convertedRetryMove)

			// Try to apply the new move
			err = g.chessGame.MoveStr(convertedRetryMove)
			if err != nil {
				slog.Debug("Second AI move also failed", "error", err)
				g.err = "AI failed to make valid move after retry"
				return nil
			}

			aiMove = newMove // Use the successful move
		} else {
			slog.Debug("✅ AI move applied successfully", "move", convertedMove, "position_after", g.chessGame.Position().String())
		}

		// Add AI move to history
		g.gameHistory = append(g.gameHistory, aiMove)
		slog.Debug("📝 AI move added to history", "history_length", len(g.gameHistory), "full_history", g.gameHistory)

		// Update status and clear AI turn flags
		g.updateStatus()
		g.isAITurn = false
		g.aiMovePending = false // Reset the pending flag

		slog.Debug("🎉 AI move completed successfully",
			"new_turn", g.chessGame.Position().Turn(),
			"isAITurn", g.isAITurn,
			"aiMovePending", g.aiMovePending,
			"status", g.status,
			"position_after", g.chessGame.Position().String())
		return aiMoveCompletedMsg{}
	}
}

// getBoardState returns the current board state as a string
func (g *Game) getBoardState() string {
	// Return FEN notation which is better for AI understanding
	return g.chessGame.Position().String()
}

// retryAIMoveWithError sends the error back to the AI server and requests a new move
func (g *Game) retryAIMoveWithError(boardState string, gameHistory []string, errorMsg string, playerColor string) (string, error) {
	slog.Debug("Retrying AI move with error", "error", errorMsg)

	// Use the AI client to make the retry request
	return g.aiClient.GetAIMoveWithError(boardState, gameHistory, errorMsg, playerColor)
}
