package ai_player

import (
	"fmt"
	"strings"
	"time"
)

// GameMode represents different ways to play against the AI
type GameMode int

const (
	ModeHumanVsAI GameMode = iota
	ModeAIvsAI
	ModeHumanVsHuman
)

// AIGame represents a chess game with AI integration
type AIGame struct {
	GameMode    GameMode
	AIWhite     *AIPlayer
	AIBlack     *AIPlayer
	MoveHistory []string
	Config      *Config
	GameState   string
	CurrentTurn string // "white" or "black"
}

// NewAIGame creates a new AI-enabled chess game
func NewAIGame(mode GameMode, config *Config) *AIGame {
	game := &AIGame{
		GameMode:    mode,
		MoveHistory: make([]string, 0),
		Config:      config,
		GameState:   "initial",
		CurrentTurn: "white",
	}

	// Initialize AI players based on game mode
	switch mode {
	case ModeHumanVsAI:
		game.AIBlack = NewAIPlayer(config.OllamaURL, config.Model, "black")
	case ModeAIvsAI:
		game.AIWhite = NewAIPlayer(config.OllamaURL, config.Model, "white")
		game.AIBlack = NewAIPlayer(config.OllamaURL, config.Model, "black")
	case ModeHumanVsHuman:
		// No AI players needed
	}

	return game
}

// GetAIMove gets the next move from the appropriate AI player
func (g *AIGame) GetAIMove(boardState string) (*ChessMove, error) {
	var aiPlayer *AIPlayer

	switch g.CurrentTurn {
	case "white":
		aiPlayer = g.AIWhite
	case "black":
		aiPlayer = g.AIBlack
	default:
		return nil, fmt.Errorf("invalid current turn: %s", g.CurrentTurn)
	}

	if aiPlayer == nil {
		return nil, fmt.Errorf("no AI player for %s", g.CurrentTurn)
	}

	// Get move from AI with retry logic
	var move *ChessMove
	var err error

	for attempt := 1; attempt <= g.Config.MaxRetries; attempt++ {
		move, err = aiPlayer.GetMove(boardState, g.MoveHistory)
		if err == nil {
			break
		}

		if attempt < g.Config.MaxRetries {
			time.Sleep(time.Duration(g.Config.RetryDelay) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("AI failed to generate move after %d attempts: %w", g.Config.MaxRetries, err)
	}

	return move, nil
}

// IsAITurn checks if it's currently the AI's turn
func (g *AIGame) IsAITurn() bool {
	switch g.GameMode {
	case ModeHumanVsAI:
		return g.CurrentTurn == "black" // AI plays black
	case ModeAIvsAI:
		return true // Both players are AI
	case ModeHumanVsHuman:
		return false // No AI players
	default:
		return false
	}
}

// GetCurrentAIPlayer returns the current AI player if it's an AI turn
func (g *AIGame) GetCurrentAIPlayer() *AIPlayer {
	if !g.IsAITurn() {
		return nil
	}

	switch g.CurrentTurn {
	case "white":
		return g.AIWhite
	case "black":
		return g.AIBlack
	default:
		return nil
	}
}

// SwitchTurn switches the current turn
func (g *AIGame) SwitchTurn() {
	if g.CurrentTurn == "white" {
		g.CurrentTurn = "black"
	} else {
		g.CurrentTurn = "white"
	}
}

// AddMove adds a move to the game history
func (g *AIGame) AddMove(move string) {
	g.MoveHistory = append(g.MoveHistory, move)

	// Keep only the last N moves as specified in config
	if len(g.MoveHistory) > g.Config.MoveHistory {
		g.MoveHistory = g.MoveHistory[len(g.MoveHistory)-g.Config.MoveHistory:]
	}
}

// GetGameStatus returns a string representation of the current game status
func (g *AIGame) GetGameStatus() string {
	var status strings.Builder

	status.WriteString(fmt.Sprintf("Game Mode: %s\n", g.getGameModeString()))
	status.WriteString(fmt.Sprintf("Current Turn: %s\n", strings.Title(g.CurrentTurn)))
	status.WriteString(fmt.Sprintf("Moves Made: %d\n", len(g.MoveHistory)))

	if g.IsAITurn() {
		status.WriteString("Next Move: AI\n")
	} else {
		status.WriteString("Next Move: Human\n")
	}

	return status.String()
}

// getGameModeString returns a human-readable string for the game mode
func (g *AIGame) getGameModeString() string {
	switch g.GameMode {
	case ModeHumanVsAI:
		return "Human vs AI"
	case ModeAIvsAI:
		return "AI vs AI"
	case ModeHumanVsHuman:
		return "Human vs Human"
	default:
		return "Unknown"
	}
}

// TestAIConnection tests the connection to Ollama for all AI players
func (g *AIGame) TestAIConnection() error {
	var errors []string

	if g.AIWhite != nil {
		if err := g.AIWhite.TestConnection(); err != nil {
			errors = append(errors, fmt.Sprintf("White AI: %v", err))
		}
	}

	if g.AIBlack != nil {
		if err := g.AIBlack.TestConnection(); err != nil {
			errors = append(errors, fmt.Sprintf("Black AI: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("AI connection test failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetAIConfig returns the current AI configuration
func (g *AIGame) GetAIConfig() *Config {
	return g.Config
}

// UpdateAIConfig updates the AI configuration
func (g *AIGame) UpdateAIConfig(newConfig *Config) error {
	if err := newConfig.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	g.Config = newConfig

	// Update AI players with new configuration
	if g.AIWhite != nil {
		g.AIWhite.OllamaURL = newConfig.OllamaURL
		g.AIWhite.Model = newConfig.Model
	}

	if g.AIBlack != nil {
		g.AIBlack.OllamaURL = newConfig.OllamaURL
		g.AIBlack.Model = newConfig.Model
	}

	return nil
}
