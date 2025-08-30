# BubbleChess

A terminal-based chess game built with Go and Bubble Tea.

## Features

- **Interactive TUI Mode**: Full-featured terminal user interface with colored board, piece selection, and move input
- **Text Mode**: Simple text-based output for environments that don't support TUI or for automation
- **Chess Notation Support**: Supports both long algebraic notation (e2e4) and short algebraic notation (Nc6, Kxe5)
- **Castling**: Full castling support with proper move validation
- **Move Validation**: Basic chess move validation for all piece types

## Installation

```bash
go install github.com/charmbracelet/bubbles@latest
go install github.com/charmbracelet/bubbletea@latest
go install github.com/charmbracelet/lipgloss@latest
go build -o bubblechess .
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

## License

MIT License
