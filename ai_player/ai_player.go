package ai_player

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaRequest represents the request sent to Ollama
type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// OllamaResponse represents the response from Ollama
type OllamaResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

// ChessMove represents a chess move in standard notation
type ChessMove struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Piece     string `json:"piece,omitempty"`
	Capture   bool   `json:"capture,omitempty"`
	Check     bool   `json:"check,omitempty"`
	Checkmate bool   `json:"checkmate,omitempty"`
	Notation  string `json:"notation"`
}

// AIPlayer represents an AI chess player
type AIPlayer struct {
	OllamaURL string
	Model     string
	Client    *http.Client
	Color     string // "white" or "black"
}

// NewAIPlayer creates a new AI player
func NewAIPlayer(ollamaURL, model, color string) *AIPlayer {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	if model == "" {
		model = "gemma3n:latest" // Default model, adjust as needed
	}

	return &AIPlayer{
		OllamaURL: ollamaURL,
		Model:     model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Color: color,
	}
}

// GetMove gets the next move from the AI player
func (ai *AIPlayer) GetMove(boardState string, gameHistory []string) (*ChessMove, error) {
	prompt := ai.buildPrompt(boardState, gameHistory)

	request := OllamaRequest{
		Model:  ai.Model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1, // Low temperature for more consistent moves
			"top_p":       0.9,
		},
	}

	response, err := ai.callOllama(request)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama: %w", err)
	}

	move, err := ai.parseMove(response.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return move, nil
}

// buildPrompt creates a prompt for the AI to generate a chess move
func (ai *AIPlayer) buildPrompt(boardState string, gameHistory []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a chess AI playing as ")
	prompt.WriteString(ai.Color)
	prompt.WriteString(". Analyze the current board position and suggest the best move.\n\n")

	prompt.WriteString("Current board position:\n")
	prompt.WriteString(boardState)
	prompt.WriteString("\n\n")

	if len(gameHistory) > 0 {
		prompt.WriteString("Game history (last 5 moves):\n")
		start := len(gameHistory) - 5
		if start < 0 {
			start = 0
		}
		for i, move := range gameHistory[start:] {
			prompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, move))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("Instructions:\n")
	prompt.WriteString("1. Analyze the position carefully\n")
	prompt.WriteString("2. Consider tactics, strategy, and piece safety\n")
	prompt.WriteString("3. Respond with ONLY the move in standard algebraic notation\n")
	prompt.WriteString("4. Use long notation (e2e4) or short notation (Nc6, Kxe5)\n")
	prompt.WriteString("5. For castling, use O-O or O-O-O\n")
	prompt.WriteString("6. Do not include any explanations or additional text\n\n")

	prompt.WriteString("Your move: ")

	return prompt.String()
}

// callOllama makes an HTTP request to the Ollama API
func (ai *AIPlayer) callOllama(request OllamaRequest) (*OllamaResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := ai.Client.Post(
		ai.OllamaURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// parseMove parses the AI's response and extracts the chess move
func (ai *AIPlayer) parseMove(response string) (*ChessMove, error) {
	// Clean up the response
	response = strings.TrimSpace(response)
	response = strings.Split(response, "\n")[0] // Take only the first line

	// Remove common prefixes/suffixes that AI might add
	response = strings.TrimPrefix(response, "Move: ")
	response = strings.TrimPrefix(response, "The best move is ")
	response = strings.TrimPrefix(response, "I suggest ")
	response = strings.TrimSuffix(response, ".")
	response = strings.TrimSuffix(response, "!")
	response = strings.TrimSuffix(response, "?")

	// Validate that it looks like a chess move
	if !ai.isValidMoveNotation(response) {
		return nil, fmt.Errorf("invalid move notation: %s", response)
	}

	return &ChessMove{
		Notation: response,
	}, nil
}

// isValidMoveNotation checks if the move notation looks valid
func (ai *AIPlayer) isValidMoveNotation(move string) bool {
	if move == "" {
		return false
	}

	// Check for castling
	if move == "O-O" || move == "0-0" || move == "O-O-O" || move == "0-0-0" {
		return true
	}

	// Check for long algebraic notation (e2e4)
	if len(move) == 4 {
		if (move[0] >= 'a' && move[0] <= 'h') &&
			(move[1] >= '1' && move[1] <= '8') &&
			(move[2] >= 'a' && move[2] <= 'h') &&
			(move[3] >= '1' && move[3] <= '8') {
			return true
		}
	}

	// Check for short algebraic notation (Nc6, Kxe5, etc.)
	if len(move) >= 2 {
		// First character should be a piece or file
		if (move[0] >= 'A' && move[0] <= 'Z') || (move[0] >= 'a' && move[0] <= 'h') {
			// Last two characters should be coordinates
			if len(move) >= 2 {
				lastTwo := move[len(move)-2:]
				if (lastTwo[0] >= 'a' && lastTwo[0] <= 'h') &&
					(lastTwo[1] >= '1' && lastTwo[1] <= '8') {
					return true
				}
			}
		}
	}

	return false
}

// TestConnection tests the connection to Ollama
func (ai *AIPlayer) TestConnection() error {
	request := OllamaRequest{
		Model:  ai.Model,
		Prompt: "Hello",
		Stream: false,
	}

	_, err := ai.callOllama(request)
	return err
}
