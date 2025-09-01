package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// AIClient represents a client for communicating with the a2a server
type AIClient struct {
	serverURL string
	client    *http.Client
}

// NewAIClient creates a new AI client
func NewAIClient(serverURL string) *AIClient {
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	return &AIClient{
		serverURL: serverURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// MessageSendParams represents the parameters for message/send method
type MessageSendParams struct {
	Message Message `json:"message"`
}

// Message represents an A2A message
type Message struct {
	Kind      string             `json:"kind"`
	MessageID string             `json:"messageId"`
	Role      string             `json:"role"`
	Parts     []MessagePartsElem `json:"parts"`
}

// MessagePartsElem represents message parts (interface type)
type MessagePartsElem interface{}

// TextPart represents a text part
type TextPart struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
}

// ChessRequest represents a chess move request to the AI
type ChessRequest struct {
	BoardState  string   `json:"board_state"`
	PlayerColor string   `json:"player_color"`
	GameHistory []string `json:"game_history"`
}

// ChessResponse represents a chess move response from the AI
type ChessResponse struct {
	Move string `json:"move"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// GetAIMove requests a move from the AI via the a2a server
func (ac *AIClient) GetAIMove(boardState string, gameHistory []string, playerColor string) (string, error) {
	return ac.getAIMoveInternal(boardState, gameHistory, "", playerColor)
}

// GetAIMoveWithError requests a move from the AI with error information from the previous attempt
func (ac *AIClient) GetAIMoveWithError(boardState string, gameHistory []string, errorMsg string, playerColor string) (string, error) {
	return ac.getAIMoveInternal(boardState, gameHistory, errorMsg, playerColor)
}

// getAIMoveInternal is the internal implementation for getting AI moves
func (ac *AIClient) getAIMoveInternal(boardState string, gameHistory []string, errorMsg string, playerColor string) (string, error) {
	// Create the JSON-RPC request
	jsonrpcRequest := JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  "message/send",
		ID:      1,
		Params: MessageSendParams{
			Message: Message{
				Kind:      "message",
				MessageID: fmt.Sprintf("msg_%d", time.Now().Unix()),
				Role:      "user",
				Parts: []MessagePartsElem{
					TextPart{
						Kind: "text",
						Text: ac.buildRequestText(boardState, gameHistory, errorMsg, playerColor),
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(jsonrpcRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON-RPC request: %w", err)
	}

	// Debug output
	slog.Debug("Making request to AI server", "url", ac.serverURL+"/a2a")
	slog.Debug("Request data", "data", string(jsonData))

	// Make request to the a2a endpoint
	resp, err := ac.client.Post(ac.serverURL+"/a2a", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Debug("Request failed", "error", err)
		return "", fmt.Errorf("failed to make request to a2a server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("a2a server returned status: %d", resp.StatusCode)
	}

	// Read the full response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug output
	slog.Debug("Response received", "status", resp.StatusCode)
	slog.Debug("Response body", "body", string(bodyBytes))

	// Parse the JSON-RPC response
	var jsonrpcResponse JSONRPCResponse
	if err := json.Unmarshal(bodyBytes, &jsonrpcResponse); err != nil {
		return "", fmt.Errorf("failed to decode JSON-RPC response: %w", err)
	}

	// Debug output removed for production

	// Enhanced debug output for response parsing
	slog.Debug("Parsing AI response", "has_result", jsonrpcResponse.Result != nil, "has_error", jsonrpcResponse.Error != nil)

	// Check for JSON-RPC errors
	if jsonrpcResponse.Error != nil {
		errorBytes, _ := json.Marshal(jsonrpcResponse.Error)
		slog.Debug("JSON-RPC error received", "error", string(errorBytes))
		return "", fmt.Errorf("JSON-RPC error: %s", string(errorBytes))
	}

	// Extract the result from the JSON-RPC response
	// The result contains a message with parts
	resultMap, ok := jsonrpcResponse.Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("result is not a map")
	}

	// Extract the parts from the result
	parts, ok := resultMap["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("no parts found in result")
	}

	// Get the first part (should be text)
	firstPart, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("first part is not a map")
	}

	// Extract the text from the part
	text, ok := firstPart["text"].(string)
	if !ok {
		return "", fmt.Errorf("text field not found in part")
	}

	slog.Debug("📝 AI response text received", "text", text, "text_length", len(text))

	// Try to extract the move from various possible response formats
	var move string

	// Format 1: "Generated move: <move>"
	if len(text) > 16 && text[:16] == "Generated move: " {
		move = text[16:]
		slog.Debug("✅ Extracted move using 'Generated move:' format", "move", move)
	} else if len(text) > 7 && text[:7] == "Move: " {
		// Format 2: "Move: <move>"
		move = text[7:]
		slog.Debug("✅ Extracted move using 'Move:' format", "move", move)
	} else if len(text) > 0 {
		// Format 3: Just the move itself (clean response)
		// Check if it looks like a valid chess move
		cleanedText := strings.TrimSpace(text)
		if len(cleanedText) >= 2 && len(cleanedText) <= 5 {
			// Basic validation - should be 2-5 characters for chess moves
			move = cleanedText
			slog.Debug("✅ Extracted move as direct response", "move", move)
		} else {
			slog.Debug("❌ Response doesn't match expected move format", "text", text)
			return "", fmt.Errorf("unexpected text format: %s", text)
		}
	} else {
		slog.Debug("❌ Empty or invalid response text", "text", text)
		return "", fmt.Errorf("empty or invalid response text")
	}

	// Validate that we extracted a move
	if move == "" {
		slog.Debug("❌ No move extracted from response", "text", text)
		return "", fmt.Errorf("no move extracted from response: %s", text)
	}

	slog.Debug("🎯 Successfully extracted AI move", "move", move, "original_text", text)
	return move, nil
}

// buildRequestText builds the request text for the AI
func (ac *AIClient) buildRequestText(boardState string, gameHistory []string, errorMsg string, playerColor string) string {
	// Convert game history to proper JSON array format
	historyJSON, _ := json.Marshal(gameHistory)

	if errorMsg == "" {
		return fmt.Sprintf(`{"board_state":"%s","player_color":"%s","game_history":%s}`, boardState, playerColor, string(historyJSON))
	}
	return fmt.Sprintf(`{"board_state":"%s","player_color":"%s","game_history":%s,"last_move_error":"%s"}`, boardState, playerColor, string(historyJSON), errorMsg)
}

// TestConnection tests the connection to the a2a server
func (ac *AIClient) TestConnection() error {
	resp, err := ac.client.Get(ac.serverURL)
	if err != nil {
		return fmt.Errorf("failed to connect to a2a server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("a2a server returned status: %d", resp.StatusCode)
	}

	return nil
}
