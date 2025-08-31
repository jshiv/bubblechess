# Chess CLI Application

A comprehensive chess application with a modern command-line interface built using Cobra.

## Features

- **TUI Chess Game**: Interactive terminal-based chess game
- **A2A Server**: JSON-RPC A2A protocol server for AI-powered chess moves
- **Unified CLI**: Single binary with multiple commands
- **Configurable**: Customizable server settings

## Installation

### Prerequisites

- Go 1.24.5 or later
- Ollama (for AI-powered moves)

### Build

```bash
# From the project root
go build ./cmd/chess
```

## Usage

### Basic Commands

```bash
# Start the TUI chess game (default)
./chess

# Show help
./chess --help

# Show help for a specific command
./chess server --help
```

### TUI Chess Game

The root command starts the interactive TUI chess game:

```bash
./chess
```

This launches the terminal-based chess interface where you can:
- Play chess interactively
- Use mouse and keyboard controls
- View the board with proper chess notation
- Make moves using standard chess notation

### A2A Server

Start the JSON-RPC A2A chess server:

```bash
./chess server [flags]
```

#### Server Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--ollama-url` | `-u` | `http://localhost:11434` | Ollama server URL |
| `--model` | `-m` | `gpt-oss:20b` | Ollama model to use |
| `--port` | `-p` | `8080` | Port to listen on |

#### Server Examples

```bash
# Start with default settings
./chess server

# Use a specific model and port
./chess server -m llama3.2:3b -p 8081

# Custom Ollama URL
./chess server -u http://192.168.1.100:11434 -m gpt-oss:20b

# Full custom configuration
./chess server --ollama-url http://localhost:11434 --model llama3.2:3b --port 8081
```

## Architecture

### Commands

- **Root Command** (`./chess`): Starts the TUI chess game
- **Server Command** (`./chess server`): Starts the A2A protocol server

### Integration Points

- **TUI Game**: Integrates with the existing chess game logic
- **A2A Server**: Integrates with the JSON-RPC A2A server implementation
- **AI Player**: Connects to Ollama for AI-powered move generation

## Development

### Project Structure

```
cmd/chess/
├── main.go          # Main CLI application
└── README.md        # This documentation
```

### Adding New Commands

To add new commands to the CLI:

1. **Create the command**:
```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description of new command",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}
```

2. **Add to root command**:
```go
func init() {
    rootCmd.AddCommand(newCmd)
}
```

3. **Add flags if needed**:
```go
newCmd.Flags().StringP("flag", "f", "default", "description")
```

### Building and Testing

```bash
# Build the CLI
go build ./cmd/chess

# Test the CLI
./chess --help
./chess server --help

# Run tests
go test ./cmd/chess
```

## Configuration

### Environment Variables

The CLI respects the following environment variables:

- `BUBBLECHESS_MODE`: Display mode for the TUI game
- `OLLAMA_URL`: Default Ollama server URL
- `OLLAMA_MODEL`: Default Ollama model

### Configuration Files

Configuration can be provided via:
- Command-line flags (highest priority)
- Environment variables
- Default values (lowest priority)

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Use a different port
   ./chess server -p 8081
   ```

2. **Ollama Connection Issues**
   ```bash
   # Check if Ollama is running
   curl http://localhost:11434/api/tags
   
   # Use custom Ollama URL
   ./chess server -u http://your-ollama-host:11434
   ```

3. **Model Not Found**
   ```bash
   # List available models
   ollama list
   
   # Pull a specific model
   ollama pull gpt-oss:20b
   ```

### Debug Mode

For debugging, you can run with verbose output:

```bash
# Enable debug logging
./chess server --debug
```

## Future Enhancements

- [ ] **Game Modes**: Human vs Human, Human vs AI, AI vs AI
- [ ] **Network Play**: Multiplayer over network
- [ ] **Game History**: Save and load games
- [ ] **Analysis Mode**: AI analysis of positions
- [ ] **Tournament Mode**: Multiple game support
- [ ] **Configuration File**: YAML/JSON config support

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of the BubbleChess application.
