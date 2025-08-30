# AI Player for BubbleChess

This package provides an AI chess player that can play against humans or other AI players using local Ollama models.

## Features

- **Local AI**: Uses Ollama to run AI models locally on your machine
- **Multiple Game Modes**: Human vs AI, AI vs AI, and Human vs Human
- **Configurable**: Easy to configure AI behavior and connection settings
- **Retry Logic**: Built-in retry mechanism for reliable AI responses
- **Move Validation**: Validates AI responses to ensure they're valid chess moves

## Prerequisites

1. **Ollama**: Install and run Ollama on your machine
   ```bash
   # Install Ollama (macOS)
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Start Ollama service
   ollama serve
   ```

2. **Chess Model**: Pull a chess-capable model
   ```bash
   # Pull a model (adjust based on your hardware)
   ollama pull llama3.2:3b
   
   # Or use a larger model for better chess play
   ollama pull llama3.2:8b
   ```

## Quick Start

### 1. Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "chess-tui/ai_player"
)

func main() {
    // Load configuration
    config, err := ai_player.LoadConfig("ai_config.json")
    if err != nil {
        log.Fatal(err)
    }

    // Create AI game
    game := ai_player.NewAIGame(ai_player.ModeHumanVsAI, config)

    // Test connection
    if err := game.TestAIConnection(); err != nil {
        log.Fatal("AI connection failed:", err)
    }

    // Get AI move
    boardState := "your chess board representation here"
    move, err := game.GetAIMove(boardState)
    if err != nil {
        log.Fatal("Failed to get AI move:", err)
    }

    fmt.Printf("AI suggests: %s\n", move.Notation)
}
```

### 2. Game Modes

```go
// Human vs AI (AI plays black)
game := ai_player.NewAIGame(ai_player.ModeHumanVsAI, config)

// AI vs AI (both players are AI)
game := ai_player.NewAIGame(ai_player.ModeAIvsAI, config)

// Human vs Human (no AI)
game := ai_player.NewAIGame(ai_player.ModeHumanVsHuman, config)
```

## Configuration

The AI player uses a JSON configuration file (`ai_config.json`) with the following options:

```json
{
  "ollama_url": "http://localhost:11434",
  "model": "llama3.2:3b",
  "timeout_seconds": 30,
  "temperature": 0.1,
  "top_p": 0.9,
  "max_retries": 3,
  "retry_delay_seconds": 2,
  "move_history_length": 5
}
```

### Configuration Options

- **ollama_url**: URL where Ollama is running (default: localhost:11434)
- **model**: Ollama model name to use for chess moves
- **timeout_seconds**: HTTP timeout for Ollama requests
- **temperature**: AI creativity (0.0 = deterministic, 2.0 = very creative)
- **top_p**: Nucleus sampling parameter for response quality
- **max_retries**: Number of retry attempts if AI fails
- **retry_delay_seconds**: Delay between retry attempts
- **move_history_length**: Number of recent moves to include in AI prompts

## AI Prompt Engineering

The AI player sends carefully crafted prompts to Ollama:

```
You are a chess AI playing as black. Analyze the current board position and suggest the best move.

Current board position:
[board representation]

Game history (last 5 moves):
1. e2e4
2. e7e5
3. Nf3

Instructions:
1. Analyze the position carefully
2. Consider tactics, strategy, and piece safety
3. Respond with ONLY the move in standard algebraic notation
4. Use long notation (e2e4) or short notation (Nc6, Kxe5)
5. For castling, use O-O or O-O-O
6. Do not include any explanations or additional text

Your move: [AI response]
```

## Supported Move Notations

The AI player accepts and generates moves in standard chess notation:

### Long Algebraic Notation
- `e2e4` - Move from e2 to e4
- `g8f6` - Move from g8 to f6

### Short Algebraic Notation
- `Nc6` - Knight to c6
- `Kxe5` - King captures on e5
- `O-O` or `0-0` - Kingside castling
- `O-O-O` or `0-0-0` - Queenside castling

## Error Handling

The AI player includes robust error handling:

- **Connection Errors**: Automatic retry with configurable attempts
- **Invalid Moves**: Response validation and parsing
- **Timeout Handling**: Configurable HTTP timeouts
- **Model Errors**: Graceful fallback and error reporting

## Performance Considerations

- **Model Size**: Larger models (8B+ parameters) provide better chess play but are slower
- **Hardware**: GPU acceleration significantly improves response times
- **Temperature**: Lower temperature (0.1) provides more consistent moves
- **Prompt Length**: Shorter prompts are processed faster

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure Ollama is running: `ollama serve`
   - Check Ollama URL in configuration

2. **Model Not Found**
   - Pull the required model: `ollama pull [model_name]`
   - Verify model name in configuration

3. **Slow Responses**
   - Use smaller models for faster responses
   - Enable GPU acceleration if available
   - Reduce temperature for more deterministic play

4. **Invalid Moves**
   - Check that the AI model understands chess notation
   - Verify board state representation is clear
   - Consider using a chess-specialized model

### Debug Mode

Enable debug logging by setting environment variables:
```bash
export DEBUG=1
export OLLAMA_DEBUG=1
```

## Examples

See the `examples/` directory for complete working examples:

- `ai_example.go` - Basic AI player usage
- Integration examples with the main chess game

## Contributing

When contributing to the AI player:

1. Follow Go coding standards
2. Add tests for new functionality
3. Update documentation for new features
4. Test with different Ollama models
5. Validate chess move parsing thoroughly

## License

MIT License - see main project license
