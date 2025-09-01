package ai_player

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ChessRequest represents a chess move request from the A2A client
type ChessRequest struct {
	BoardState  string   `json:"board_state,omitempty"`
	PlayerColor string   `json:"player_color,omitempty"`
	GameHistory []string `json:"game_history,omitempty"`
}

// ChessResponse represents a chess move response from the AI
type ChessResponse struct {
	Move string `json:"move"`
}

// JSONRPCA2AServer represents an A2A server using the generated JSON-RPC spec
type JSONRPCA2AServer struct {
	aiPlayer *AIPlayer
	server   *http.Server
	logger   *log.Logger
}

// NewJSONRPCA2AServer creates a new A2A server using the generated JSON-RPC spec
func NewJSONRPCA2AServer(ollamaURL, model string, port int, logger *log.Logger) (*JSONRPCA2AServer, error) {
	// Create AI player
	aiPlayer := NewAIPlayer(ollamaURL, model, "black")

	// Test connection to Ollama
	if err := aiPlayer.TestConnection(); err != nil {
		return nil, fmt.Errorf("failed to test Ollama connection: %w", err)
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// Add A2A endpoints
	mux.HandleFunc("/", handleJSONRPCRoot)
	mux.HandleFunc("/.well-known/agent.json", handleJSONRPCAgentCard)
	mux.HandleFunc("/a2a", handleJSONRPCEndpoint(aiPlayer, logger))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &JSONRPCA2AServer{
		aiPlayer: aiPlayer,
		server:   httpServer,
		logger:   logger,
	}, nil
}

// Start starts the JSON-RPC A2A server
func (s *JSONRPCA2AServer) Start() error {
	s.logger.Printf("Starting JSON-RPC A2A Chess Server on :8080")
	s.logger.Printf("AI Model: %s", s.aiPlayer.Model)
	s.logger.Printf("Ollama URL: %s", s.aiPlayer.OllamaURL)

	return s.server.ListenAndServe()
}

// Stop stops the JSON-RPC A2A server gracefully
func (s *JSONRPCA2AServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleJSONRPCRoot handles the root endpoint
func handleJSONRPCRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"service":  "Chess JSON-RPC A2A Server",
		"version":  "1.0.0",
		"protocol": "A2A (Agent-to-Agent) with JSON-RPC 2.0",
		"endpoints": map[string]string{
			"agent_card": "/.well-known/agent.json",
			"a2a":        "/a2a",
		},
		"description": "A2A protocol server for chess AI moves using Ollama and generated JSON-RPC spec",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleJSONRPCAgentCard handles the A2A agent discovery endpoint
func handleJSONRPCAgentCard(w http.ResponseWriter, r *http.Request) {
	agentCard := AgentCard{
		Name:               "Chess AI Player",
		Description:        "An AI chess player that generates moves using Ollama models",
		Url:                "http://localhost:8080",
		Version:            "1.0.0",
		ProtocolVersion:    "1.0.0",
		PreferredTransport: "JSONRPC",
		Capabilities: AgentCapabilities{
			Streaming:         &[]bool{false}[0],
			PushNotifications: &[]bool{false}[0],
		},
		DefaultInputModes:  []string{"text/plain", "application/json"},
		DefaultOutputModes: []string{"text/plain", "application/json"},
		Skills: []AgentSkill{
			{
				Name:        "chess_move_generation",
				Description: "Generate chess moves using AI analysis",
				InputModes:  []string{"text/plain", "application/json"},
				OutputModes: []string{"text/plain", "application/json"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentCard)
}

// handleJSONRPCEndpoint handles A2A JSON-RPC protocol requests
func handleJSONRPCEndpoint(aiPlayer *AIPlayer, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendJSONRPCError(w, -32600, "Method Not Allowed", "Only POST method is supported", nil)
			return
		}

		// Parse the request body to determine the method
		var rawRequest map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&rawRequest); err != nil {
			sendJSONRPCError(w, -32700, "Parse error", err.Error(), nil)
			return
		}

		// Extract method and ID for error handling
		method, _ := rawRequest["method"].(string)
		requestID := rawRequest["id"]

		// Handle different A2A methods
		switch method {
		case "message/send":
			handleJSONRPCMessageSend(w, r, rawRequest, aiPlayer, logger)
		case "tasks/send":
			handleJSONRPCTasksSend(w, r, rawRequest, aiPlayer, logger)
		default:
			sendJSONRPCError(w, -32601, "Method not found", fmt.Sprintf("Method '%s' not found", method), requestID)
		}
	}
}

// handleJSONRPCMessageSend handles the message/send method for JSON-RPC
func handleJSONRPCMessageSend(w http.ResponseWriter, r *http.Request, request map[string]interface{}, aiPlayer *AIPlayer, logger *log.Logger) {
	logger.Printf("Received A2A message/send request")
	logger.Printf("Raw request: %+v", request)

	// Extract ID for error handling
	requestID := request["id"]

	// Parse the request using the generated spec
	var requestSendMessage SendMessageRequest
	requestBytes, _ := json.Marshal(request)
	logger.Printf("Request bytes: %s", string(requestBytes))
	if err := json.Unmarshal(requestBytes, &requestSendMessage); err != nil {
		logger.Printf("Failed to parse SendMessageRequest: %v", err)
		sendJSONRPCError(w, -32602, "Invalid params", fmt.Sprintf("Failed to parse request: %v", err), requestID)
		return
	}
	logger.Printf("Parsed request: %+v", requestSendMessage)

	// Parse chess request from message
	var chessReq ChessRequest
	if err := parseChessRequestFromJSONRPCMessage(requestSendMessage.Params.Message, &chessReq); err != nil {
		sendJSONRPCError(w, -32602, "Invalid params", fmt.Sprintf("Failed to parse chess request: %v", err), requestID)
		return
	}

	// Process chess request
	result, err := processChessRequest(chessReq, aiPlayer, logger)
	if err != nil {
		sendJSONRPCError(w, -32603, "Internal error", fmt.Sprintf("Chess processing failed: %v", err), requestID)
		return
	}

	// Create A2A message response
	responseMessage := Message{
		Kind:      "message",
		MessageId: fmt.Sprintf("msg_%d", time.Now().Unix()),
		Role:      MessageRoleAgent,
		Parts: []MessagePartsElem{
			TextPart{
				Kind: "text",
				Text: fmt.Sprintf("Generated move: %s", result.Move),
			},
		},
	}

	// Create A2A success response
	response := SendMessageSuccessResponse{
		Jsonrpc: "2.0",
		Id:      requestID,
		Result: SendMessageSuccessResponseResult{
			Kind:      "message",
			MessageId: responseMessage.MessageId,
			Role:      responseMessage.Role,
			Parts:     responseMessage.Parts,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleJSONRPCTasksSend handles the A2A tasks/send method
func handleJSONRPCTasksSend(w http.ResponseWriter, r *http.Request, rawRequest map[string]interface{}, aiPlayer *AIPlayer, logger *log.Logger) {
	logger.Printf("Received A2A tasks/send request")

	// For now, we'll handle this the same as message/send
	// In a full implementation, this would create a task and return task status
	handleJSONRPCMessageSend(w, r, rawRequest, aiPlayer, logger)
}

// parseChessRequestFromJSONRPCMessage parses chess request from JSON-RPC A2A message
func parseChessRequestFromJSONRPCMessage(message Message, req *ChessRequest) error {
	for _, part := range message.Parts {
		// Try to convert to TextPart
		partBytes, _ := json.Marshal(part)
		var textPart TextPart
		if err := json.Unmarshal(partBytes, &textPart); err == nil && textPart.Kind == "text" {
			// Try to parse as JSON first
			if err := json.Unmarshal([]byte(textPart.Text), req); err == nil {
				return nil
			}

			// If not JSON, try to parse as simple board state
			req.BoardState = strings.TrimSpace(textPart.Text)
			return nil
		}
	}

	return fmt.Errorf("no text part found in message")
}

// sendJSONRPCError sends a JSON-RPC error response
func sendJSONRPCError(w http.ResponseWriter, code int, message, data string, id interface{}) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// processChessRequest processes a chess request and returns a move
func processChessRequest(req ChessRequest, aiPlayer *AIPlayer, logger *log.Logger) (*ChessResponse, error) {
	logger.Printf("🎮 [JSONRPCA2A] Processing chess request - Player: %s, Board state length: %d, History: %v",
		req.PlayerColor, len(req.BoardState), req.GameHistory)

	// Set AI player color based on request
	aiPlayer.Color = req.PlayerColor
	logger.Printf("🎨 [JSONRPCA2A] AI player color set to: %s", aiPlayer.Color)

	// Log board state for debugging
	logger.Printf("📊 [JSONRPCA2A] Board state: %s", req.BoardState)
	if len(req.GameHistory) > 0 {
		logger.Printf("📜 [JSONRPCA2A] Game history: %v", req.GameHistory)
	}

	// Get AI move
	logger.Printf("🤖 [JSONRPCA2A] Requesting AI move...")
	startTime := time.Now()

	// Start a goroutine to log progress
	progressCtx, cancelProgress := context.WithCancel(context.Background())
	defer cancelProgress()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-progressCtx.Done():
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				logger.Printf("⏱️ [JSONRPCA2A] Still thinking... (elapsed: %v)", elapsed.Round(time.Second))
			}
		}
	}()

	aiMove, err := aiPlayer.GetMove(req.BoardState, req.GameHistory)
	cancelProgress() // Stop progress logging

	elapsed := time.Since(startTime)

	if err != nil {
		logger.Printf("❌ [JSONRPCA2A] AI move generation failed after %v: %v", elapsed, err)
		return nil, fmt.Errorf("AI move generation failed: %w", err)
	}

	logger.Printf("✅ [JSONRPCA2A] AI move generated successfully in %v: %s", elapsed, aiMove.Notation)

	return &ChessResponse{
		Move: aiMove.Notation,
	}, nil
}

// StartJSONRPCA2AServer starts the JSON-RPC A2A server
func StartJSONRPCA2AServer(ollamaURL, model string, port int) error {
	logger := log.New(log.Writer(), "[JSONRPCA2A] ", log.LstdFlags)

	server, err := NewJSONRPCA2AServer(ollamaURL, model, port, logger)
	if err != nil {
		return fmt.Errorf("failed to create JSON-RPC A2A server: %w", err)
	}

	return server.Start()
}
