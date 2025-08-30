package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chess-tui/ai_player"
)

func main() {
	// Parse command line flags
	ollamaURL := flag.String("ollama-url", "http://localhost:11434", "Ollama server URL")
	model := flag.String("model", "llama3.2:3b", "Ollama model to use")
	flag.Parse()

	// Start the JSON-RPC A2A server
	log.Printf("Starting Chess JSON-RPC A2A Server...")
	log.Printf("Ollama URL: %s", *ollamaURL)
	log.Printf("Model: %s", *model)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping server...")
		os.Exit(0)
	}()

	// Start the server
	if err := ai_player.StartJSONRPCA2AServer(*ollamaURL, *model); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
