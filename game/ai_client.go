package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (ac *AIClient) GetAIMove(boardState string, gameHistory []string) (string, error) {
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
						Text: fmt.Sprintf(`{"board_state":"%s","player_color":"%s","game_history":%v}`, boardState, "black", gameHistory),
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
	fmt.Printf("DEBUG: Making request to %s/a2a\n", ac.serverURL)
	fmt.Printf("DEBUG: Request data: %s\n", string(jsonData))

	// Make request to the a2a endpoint
	resp, err := ac.client.Post(ac.serverURL+"/a2a", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("DEBUG: Request failed: %v\n", err)
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
	fmt.Printf("DEBUG: Response status: %d\n", resp.StatusCode)
	fmt.Printf("DEBUG: Response body: %s\n", string(bodyBytes))

	// Parse the JSON-RPC response
	var jsonrpcResponse JSONRPCResponse
	if err := json.Unmarshal(bodyBytes, &jsonrpcResponse); err != nil {
		return "", fmt.Errorf("failed to decode JSON-RPC response: %w", err)
	}

	// Debug output removed for production

	// Check for JSON-RPC errors
	if jsonrpcResponse.Error != nil {
		errorBytes, _ := json.Marshal(jsonrpcResponse.Error)
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

	// The text should contain "Generated move: <move>"
	// Extract the move from the text
	if len(text) > 16 && text[:16] == "Generated move: " {
		move := text[16:]
		return move, nil
	}

	return "", fmt.Errorf("unexpected text format: %s", text)
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
