# BubbleChess

A terminal-based chess game built with Go and Bubble Tea.

## Features

- **Interactive TUI Mode**: Full-featured terminal user interface with colored board, piece selection, and move input
- **Text Mode**: Simple text-based output for environments that don't support TUI or for automation
- **AI Player**: Play against local AI powered by Ollama models
- **Multiple Game Modes**: Human vs AI, AI vs AI, and Human vs Human
- **Chess Notation Support**: Supports both long algebraic notation (e2e4) and short algebraic notation (Nc6, Kxe5)
- **Castling**: Full castling support with proper move validation
- **Move Validation**: Basic chess move validation for all piece types
- **New Game Package**: Clean TUI implementation using the notnil/chess library for accurate chess logic

## Installation

```bash
go install github.com/charmbracelet/bubbles@latest
go install github.com/charmbracelet/bubbletea@latest
go install github.com/charmbracelet/lipgloss@latest
go build -o bubblechess .
```

### AI Player Setup

To use the AI player, you'll need to install and run Ollama:

```bash
# Install Ollama (macOS)
curl -fsSL https://ollama.ai/install.sh | sh

# Start Ollama service
ollama serve

# Pull a chess-capable model
ollama pull llama3.2:3b
```

## Usage

### TUI Mode (Default)

Run the game normally to get the interactive TUI:

```bash
./bubblechess
```

You'll see a menu to choose between TUI and Text modes. Use arrow keys to navigate and Enter to select.

### Text Mode

#### Option 1: Environment Variable

Set the `BUBBLECHESS_MODE` environment variable to skip the menu and go directly to text mode:

```bash
BUBBLECHESS_MODE=text ./bubblechess
```

#### Option 2: Menu Selection

Run the game normally and select "Text Mode (Simple Output)" from the menu.

### Controls

#### TUI Mode
- **Arrow Keys**: Navigate the board or menu
- **Enter**: Select menu option or submit move
- **Ctrl+C**: Exit the game

#### Text Mode
- **Type moves**: Enter chess notation (e.g., e2e4, Nc6, Kxe5)
- **Enter**: Submit move
- **Ctrl+C**: Exit the game

## Chess Notation

The game supports both long and short algebraic notation:

### Long Algebraic Notation
- `e2e4` - Move from e2 to e4
- `g8f6` - Move from g8 to f6

### Short Algebraic Notation
- `Nc6` - Knight to c6
- `Kxe5` - King captures on e5
- `O-O` or `0-0` - Kingside castling
- `O-O-O` or `0-0-0` - Queenside castling

## Examples

### Running in Text Mode
```bash
# Set environment variable for text mode
export BUBBLECHESS_MODE=text
./bubblechess

# Or run directly
BUBBLECHESS_MODE=text ./bubblechess
```

### New Game Package

The new `game` package provides a clean, modern TUI implementation that directly uses the `notnil/chess` library:

```bash
# Run the new game TUI
go run cmd/game/main.go

# Or run the example
go run examples/game_example.go
```

The new game package features:
- Beautiful Unicode chess piece rendering
- Direct integration with the chess library for accurate rules
- Clean separation of concerns (TUI vs chess logic)
- Modern Bubble Tea TUI framework
- Proper error handling and user feedback

### Sample Game Session (Text Mode)
```
=== BubbleChess (Text Mode) ===
Status: White's turn
Current Player: White

  a b c d e f g h
8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ 8
7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ 7
6 . . . . . . . . 6
5 . . . . . . . . 5
4 . . . . . . . . 4
3 . . . . . . . . 3
2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ 2
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ 1
  a b c d e f g h

Enter move (e.g. e2e4, Nc6, Kxe5): e2e4
```

## Development

The game is built with:
- **Go 1.24+**: Core language
- **Bubble Tea**: TUI framework
- **Bubbles**: UI components
- **Lipgloss**: Terminal styling

### Project Structure

```
bubblechess/
├── main.go              # Main chess game
├── ai_player/           # AI player package
│   ├── ai_player.go     # Core AI player implementation
│   ├── config.go        # Configuration management
│   ├── game_mode.go     # Game mode definitions
│   └── README.md        # AI player documentation
├── examples/            # Example programs
│   └── ai_example.go    # AI player usage example
├── ai_config.json       # AI player configuration
└── README.md            # This file
```

## AI Player

The AI player allows you to play chess against local AI models powered by Ollama. See the [AI Player README](ai_player/README.md) for detailed documentation.

### Quick AI Example

```bash
# Test the AI player
./ai_example

# Or integrate with the main game (coming soon)
```

## License

MIT License
