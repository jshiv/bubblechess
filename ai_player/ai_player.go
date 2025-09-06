package ai_player

import (
	"bufio"
	"bytes"
	"context"
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
	Logger    *ColoredLogger
}

// NewAIPlayer creates a new AI player
func NewAIPlayer(ollamaURL, model, color string, logger *ColoredLogger) *AIPlayer {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	if model == "" {
		model = "gemma3n:latest" // Default model, adjust as needed
	}
	if logger == nil {
		logger = NewAIPlayerLogger()
	}

	return &AIPlayer{
		OllamaURL: ollamaURL,
		Model:     model,
		Client: &http.Client{
			Timeout: 60 * time.Second, // Reduced timeout to 1 minute for faster responses
		},
		Color:  color,
		Logger: logger,
	}
}

// GetMove gets the next move from the AI player
func (ai *AIPlayer) GetMove(boardState string, gameHistory []string) (*ChessMove, error) {
	ai.Logger.Debug("🎯 %sAI GetMove called - Color: %s, Board: %d chars, History: %d moves%s",
		ColorBlue, ai.Color, len(boardState), len(gameHistory), ColorReset)

	prompt := ai.buildPrompt(boardState, gameHistory)
	ai.Logger.Debug("📝 %sGenerated prompt: %d chars%s", ColorCyan, len(prompt), ColorReset)

	request := OllamaRequest{
		Model:  ai.Model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature":    0.3, // Slightly higher for faster decisions
			"top_p":          0.8, // Lower for more focused responses
			"top_k":          20,  // Limit vocabulary for faster generation
			"repeat_penalty": 1.1, // Prevent repetitive thinking
		},
	}

	ai.Logger.Debug("🚀 %sCalling Ollama API - Model: %s%s", ColorGreen, ai.Model, ColorReset)

	response, err := ai.callOllama(request)
	if err != nil {
		ai.Logger.Error("❌ %sOllama API call failed: %v%s", ColorRed, err, ColorReset)
		return nil, fmt.Errorf("failed to call Ollama: %w", err)
	}

	ai.Logger.Debug("✅ %sOllama API call successful - Response: %d chars%s", ColorGreen, len(response.Response), ColorReset)

	move, err := ai.parseMove(response.Response)
	if err != nil {
		ai.Logger.Error("❌ %sFailed to parse AI response: %v - Raw: %s%s", ColorRed, err, response.Response, ColorReset)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	ai.Logger.Debug("🎉 %sSuccessfully parsed AI move: %s%s", ColorGreen, move.Notation, ColorReset)
	return move, nil
}

// buildPrompt creates a prompt for the AI to generate a chess move
func (ai *AIPlayer) buildPrompt(boardState string, gameHistory []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are a chess AI playing as ")
	prompt.WriteString(ai.Color)
	prompt.WriteString(". Make a quick, solid move.\n\n")

	prompt.WriteString("Current board position:\n")
	prompt.WriteString(boardState)
	prompt.WriteString("\n\n")

	if len(gameHistory) > 0 {
		prompt.WriteString("Game history (last 3 moves):\n")
		start := len(gameHistory) - 3
		if start < 0 {
			start = 0
		}
		for i, move := range gameHistory[start:] {
			prompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, move))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("SPEED INSTRUCTIONS:\n")
	prompt.WriteString("1. Think FAST - spend no more than 10-15 seconds analyzing\n")
	prompt.WriteString("2. Look for obvious tactics (checks, captures, threats) first\n")
	prompt.WriteString("3. If no tactics, make a developing move (develop pieces, control center)\n")
	prompt.WriteString("4. Avoid overthinking - pick a reasonable move quickly\n")
	prompt.WriteString("5. DO NOT spend time on deep positional analysis\n\n")

	prompt.WriteString("CRITICAL FORMAT:\n")
	prompt.WriteString("1. You MUST respond with ONLY the move in SHORT ALGEBRAIC NOTATION\n")
	prompt.WriteString("2. Use SHORT notation format: e4, e5, Nf3, Nc6, Bb5, etc.\n")
	prompt.WriteString("3. For castling, use O-O (kingside) or O-O-O (queenside)\n")
	prompt.WriteString("4. For captures, use exd5 (pawn captures) or Nxe5 (piece captures)\n")
	prompt.WriteString("5. DO NOT include any explanations, analysis, or additional text\n")
	prompt.WriteString("6. DO NOT use long notation like e2e4, g1f3\n")
	prompt.WriteString("7. Your response must be exactly one move in short algebraic notation\n\n")

	prompt.WriteString("Your move (short algebraic notation only): ")

	finalPrompt := prompt.String()
	ai.Logger.Debug("📝 %sPrompt construction complete - Length: %d chars, Speed: fast_thinking%s",
		ColorCyan, len(finalPrompt), ColorReset)

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

	ai.Logger.Info("🚀 %sStarting Ollama API call - Model: %s, Prompt: %d chars%s",
		ColorGreen, request.Model, len(request.Prompt), ColorReset)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Reduced timeout to 1 minute for faster responses
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
	var thinkingBuffer strings.Builder
	var lastProgressTime time.Time
	startTime := time.Now()
	lineCount := 0

	ai.Logger.Info("📖 %sStarting to read streaming response%s", ColorBlue, ColorReset)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if line == "" {
			continue
		}

		// Parse streaming response - handle both "thinking" and "response" fields
		var streamResp struct {
			Response string `json:"response"`
			Thinking string `json:"thinking"`
			Done     bool   `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			ai.Logger.Debug("⚠️ %sFailed to parse streaming response line: %s - Error: %v%s",
				ColorYellow, line, err, ColorReset)
			continue
		}

		// Capture thinking content (this is where Ollama shows its analysis)
		if streamResp.Thinking != "" {
			thinkingBuffer.WriteString(streamResp.Thinking)

			// Log thinking progress every 15 seconds
			if time.Since(lastProgressTime) > 15*time.Second {
				elapsed := time.Since(startTime)
				currentThinking := thinkingBuffer.String()
				// Show last 100 characters of thinking to avoid log spam
				if len(currentThinking) > 100 {
					currentThinking = "..." + currentThinking[len(currentThinking)-100:]
				}
				ai.Logger.Info("🧠 %sOllama thinking progress - Elapsed: %v, Length: %d chars, Current: %s%s",
					ColorPurple, elapsed.Round(time.Second), thinkingBuffer.Len(), currentThinking, ColorReset)
				lastProgressTime = time.Now()
			}
		}

		// Add to full response (this is the actual move when done)
		if streamResp.Response != "" {
			fullResponse.WriteString(streamResp.Response)
			ai.Logger.Info("📝 %sResponse content received: %s%s", ColorCyan, streamResp.Response, ColorReset)
		}

		// Check if done
		if streamResp.Done {
			elapsed := time.Since(startTime)
			ai.Logger.Info("✅ %sOllama response completed - Time: %v, Response: %d chars, Thinking: %d chars, Lines: %d%s",
				ColorGreen, elapsed.Round(100*time.Millisecond), fullResponse.Len(), thinkingBuffer.Len(), lineCount, ColorReset)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		ai.Logger.Error("❌ %sScanner error: %v - Lines processed: %d%s", ColorRed, err, lineCount, ColorReset)
		return nil, fmt.Errorf("failed to read streaming response: %w", err)
	}

	// Log final response details
	ai.Logger.Info("📊 %sStreaming response summary - Lines: %d, Response: %d chars, Thinking: %d chars, Final: %s%s",
		ColorBlue, lineCount, fullResponse.Len(), thinkingBuffer.Len(), fullResponse.String(), ColorReset)

	// Create final response
	response := &OllamaResponse{
		Response: fullResponse.String(),
	}

	return response, nil
}

// parseMove parses the AI's response and extracts the chess move
func (ai *AIPlayer) parseMove(response string) (*ChessMove, error) {
	ai.Logger.Debug("🔍 %sParsing AI response - Raw: %s, Length: %d chars%s",
		ColorBlue, response, len(response), ColorReset)

	// Clean up the response
	response = strings.TrimSpace(response)
	response = strings.Split(response, "\n")[0] // Take only the first line
	ai.Logger.Debug("🧹 %sCleaned response: %s%s", ColorCyan, response, ColorReset)

	// Remove common prefixes/suffixes that AI might add
	originalResponse := response
	response = strings.TrimPrefix(response, "Move: ")
	response = strings.TrimPrefix(response, "The best move is ")
	response = strings.TrimPrefix(response, "I suggest ")
	response = strings.TrimSuffix(response, ".")
	response = strings.TrimSuffix(response, "!")
	response = strings.TrimSuffix(response, "?")

	if originalResponse != response {
		ai.Logger.Debug("✂️ %sRemoved prefixes/suffixes - Original: %s, Cleaned: %s%s",
			ColorYellow, originalResponse, response, ColorReset)
	}

	// Validate that it looks like a chess move
	if !ai.isValidMoveNotation(response) {
		ai.Logger.Error("❌ %sInvalid move notation - Cleaned: %s, Original: %s%s",
			ColorRed, response, originalResponse, ColorReset)
		return nil, fmt.Errorf("invalid move notation: %s", response)
	}

	ai.Logger.Debug("✅ %sMove notation validated: %s%s", ColorGreen, response, ColorReset)

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
	ai.Logger.Info("🔍 %sTesting Ollama connection - URL: %s%s", ColorBlue, ai.OllamaURL, ColorReset)

	// Test basic connectivity
	resp, err := ai.Client.Get(ai.OllamaURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	ai.Logger.Info("✅ %sOllama connection test successful%s", ColorGreen, ColorReset)
	return nil
}

// TestModelResponse tests if the specific model can respond
func (ai *AIPlayer) TestModelResponse() error {
	ai.Logger.Info("🧪 %sTesting model response - Model: %s%s", ColorPurple, ai.Model, ColorReset)

	// Create a simple test request
	testRequest := OllamaRequest{
		Model:  ai.Model,
		Prompt: "Say 'hello' in one word.",
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
		},
	}

	jsonData, err := json.Marshal(testRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal test request: %w", err)
	}

	// Create context with shorter timeout for test
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", ai.OllamaURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	startTime := time.Now()
	resp, err := ai.Client.Do(req)
	if err != nil {
		return fmt.Errorf("test request failed: %w", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(startTime)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("test request returned status %d: %s", resp.StatusCode, string(body))
	}

	var testResponse OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&testResponse); err != nil {
		return fmt.Errorf("failed to decode test response: %w", err)
	}

	ai.Logger.Info("✅ %sModel test successful - Model: %s, Time: %v, Response: %s%s",
		ColorGreen, ai.Model, elapsed.Round(100*time.Millisecond), testResponse.Response, ColorReset)

	return nil
}
