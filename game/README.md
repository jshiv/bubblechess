# Chess Game Package

This package provides a Terminal User Interface (TUI) for playing chess using the `github.com/notnil/chess` library.

## Features

- Beautiful chess board rendering with Unicode piece symbols
- Input handling for chess moves in algebraic notation
- Game state management (turns, game over, etc.)
- Error handling for invalid moves
- Game reset functionality
- Help system
- Game mode selection menu (Human vs Human, Human vs AI)
- AI integration via a2a JSON-RPC server (Human vs AI mode)

## Usage

The game can be run using the main executable in `cmd/game/`:

```bash
go run cmd/game/main.go
```

## Controls

### Menu Navigation
- **Up/Down arrows** or **j/k**: Navigate between menu options
- **Enter**: Select the highlighted option
- **q** or **Ctrl+C**: Quit the application

### Game Controls
- **Move input**: Type chess moves in algebraic notation (e.g., `e2e4`, `Nf3`, `O-O`)
- **Reset game**: Press `r` to reset the game to starting position
- **Help**: Press `h` to show help information
- **Quit**: Press `q` or `Ctrl+C` to exit

### AI Mode
When playing in **Human vs AI** mode:
- You play as White (first move)
- AI plays as Black (responds to your moves)
- AI moves are automatically requested from the a2a server
- The input is disabled during AI thinking time
- AI moves are validated and applied to the board

## Dependencies

- `github.com/notnil/chess` - Chess game logic and rules
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles/textinput` - Text input component
- `github.com/charmbracelet/lipgloss` - Styling and colors

## Architecture

The package is designed to be a thin wrapper around the chess library, handling only the TUI aspects while delegating all chess logic to the underlying library. This ensures:

- Correct chess rules enforcement
- Proper move validation
- Accurate game state tracking
- No duplication of chess logic
