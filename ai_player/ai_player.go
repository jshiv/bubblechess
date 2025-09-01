package ai_player

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
	slog.Debug("🎯 AI GetMove called", "color", ai.Color, "board_state_length", len(boardState), "history_length", len(gameHistory))

	prompt := ai.buildPrompt(boardState, gameHistory)
	slog.Debug("📝 Generated prompt", "prompt_length", len(prompt))

	request := OllamaRequest{
		Model:  ai.Model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1, // Low temperature for more consistent moves
			"top_p":       0.9,
		},
	}

	slog.Debug("🚀 Calling Ollama API", "model", ai.Model)

	response, err := ai.callOllama(request)
	if err != nil {
		slog.Error("❌ Ollama API call failed", "error", err)
		return nil, fmt.Errorf("failed to call Ollama: %w", err)
	}

	slog.Debug("✅ Ollama API call successful", "response_length", len(response.Response))

	move, err := ai.parseMove(response.Response)
	if err != nil {
		slog.Error("❌ Failed to parse AI response", "error", err, "raw_response", response.Response)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	slog.Debug("🎉 Successfully parsed AI move", "move", move.Notation)
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

	prompt.WriteString("CRITICAL INSTRUCTIONS:\n")
	prompt.WriteString("1. Analyze the position carefully for tactics, strategy, and piece safety\n")
	prompt.WriteString("2. You MUST respond with ONLY the move in LONG ALGEBRAIC NOTATION\n")
	prompt.WriteString("3. Use LONG notation format: e2e4, e7e5, g1f3, d7d5, etc.\n")
	prompt.WriteString("4. For castling, use e1g1 (white kingside) or e8c8 (black queenside)\n")
	prompt.WriteString("5. DO NOT include any explanations, analysis, or additional text\n")
	prompt.WriteString("6. DO NOT use short notation like e4, Nf3, O-O\n")
	prompt.WriteString("7. Your response must be exactly one move in long algebraic notation\n\n")

	prompt.WriteString("Your move (long algebraic notation only): ")

	finalPrompt := prompt.String()
	slog.Debug("📝 Prompt construction complete", "prompt_length", len(finalPrompt), "notation_requirement", "long_algebraic_only")

	return finalPrompt
}

// callOllama makes an HTTP request to the Ollama API with streaming support
func (ai *AIPlayer) callOllama(request OllamaRequest) (*OllamaResponse, error) {
	// Enable streaming for better progress tracking
	request.Stream = true

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	slog.Info("🚀 Starting Ollama API call", "model", request.Model, "prompt_length", len(request.Prompt))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Increased timeout to 60 seconds
	defer cancel()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", ai.OllamaURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := ai.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Handle streaming response
	var fullResponse strings.Builder
	var lastResponseTime time.Time
	startTime := time.Now()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse streaming response
		var streamResp struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			slog.Debug("Failed to parse streaming response line", "line", line, "error", err)
			continue
		}

		// Add to full response
		if streamResp.Response != "" {
			fullResponse.WriteString(streamResp.Response)

			// Log progress every 500ms
			if time.Since(lastResponseTime) > 500*time.Millisecond {
				elapsed := time.Since(startTime)
				slog.Info("💭 Ollama thinking...",
					"elapsed", elapsed.Round(100*time.Millisecond),
					"response_length", fullResponse.Len(),
					"current_response", streamResp.Response)
				lastResponseTime = time.Now()
			}
		}

		// Check if done
		if streamResp.Done {
			elapsed := time.Since(startTime)
			slog.Info("✅ Ollama response completed",
				"total_time", elapsed.Round(100*time.Millisecond),
				"total_response_length", fullResponse.Len())
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read streaming response: %w", err)
	}

	// Create final response
	response := &OllamaResponse{
		Response: fullResponse.String(),
	}

	return response, nil
}

// parseMove parses the AI's response and extracts the chess move
func (ai *AIPlayer) parseMove(response string) (*ChessMove, error) {
	slog.Debug("🔍 Parsing AI response", "raw_response", response, "response_length", len(response))

	// Clean up the response
	response = strings.TrimSpace(response)
	response = strings.Split(response, "\n")[0] // Take only the first line
	slog.Debug("🧹 Cleaned response", "cleaned_response", response)

	// Remove common prefixes/suffixes that AI might add
	originalResponse := response
	response = strings.TrimPrefix(response, "Move: ")
	response = strings.TrimPrefix(response, "The best move is ")
	response = strings.TrimPrefix(response, "I suggest ")
	response = strings.TrimSuffix(response, ".")
	response = strings.TrimSuffix(response, "!")
	response = strings.TrimSuffix(response, "?")

	if originalResponse != response {
		slog.Debug("✂️ Removed prefixes/suffixes", "original", originalResponse, "cleaned", response)
	}

	// Validate that it looks like a chess move
	if !ai.isValidMoveNotation(response) {
		slog.Error("❌ Invalid move notation", "cleaned_response", response, "original_response", originalResponse)
		return nil, fmt.Errorf("invalid move notation: %s", response)
	}

	slog.Debug("✅ Move notation validated", "final_move", response)

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
