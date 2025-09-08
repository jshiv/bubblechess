package ai_player

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

// LogLevel represents different log levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ColoredLogger provides a custom logger with colors and shorter timestamps
type ColoredLogger struct {
	level  LogLevel
	logger *log.Logger
}

// NewColoredLogger creates a new colored logger
func NewColoredLogger(level LogLevel) *ColoredLogger {
	return &ColoredLogger{
		level:  level,
		logger: log.New(os.Stdout, "", 0), // No prefix, we'll handle it ourselves
	}
}

// formatTime returns a shorter timestamp format
func (cl *ColoredLogger) formatTime() string {
	return time.Now().Format("15:04:05")
}

// getColor returns the appropriate color for the log level
func (cl *ColoredLogger) getColor(level LogLevel) string {
	switch level {
	case LevelDebug:
		return ColorGray
	case LevelInfo:
		return ColorCyan
	case LevelWarn:
		return ColorYellow
	case LevelError:
		return ColorRed
	default:
		return ColorWhite
	}
}

// getLevelString returns the string representation of the log level
func (cl *ColoredLogger) getLevelString(level LogLevel) string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO "
	case LevelWarn:
		return "WARN "
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// log prints a message with the specified level
func (cl *ColoredLogger) log(level LogLevel, format string, args ...interface{}) {
	if level < cl.level {
		return
	}

	timestamp := cl.formatTime()
	color := cl.getColor(level)
	levelStr := cl.getLevelString(level)
	reset := ColorReset

	// Format the message
	message := fmt.Sprintf(format, args...)

	// Create the final log line with color
	logLine := fmt.Sprintf("%s%s %s%s %s%s",
		ColorGray, timestamp,
		color, levelStr,
		reset, message)

	cl.logger.Println(logLine)
}

// Debug logs a debug message
func (cl *ColoredLogger) Debug(format string, args ...interface{}) {
	cl.log(LevelDebug, format, args...)
}

// Info logs an info message
func (cl *ColoredLogger) Info(format string, args ...interface{}) {
	cl.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func (cl *ColoredLogger) Warn(format string, args ...interface{}) {
	cl.log(LevelWarn, format, args...)
}

// Error logs an error message
func (cl *ColoredLogger) Error(format string, args ...interface{}) {
	cl.log(LevelError, format, args...)
}

// Printf logs a message at info level (for compatibility with log.Logger)
func (cl *ColoredLogger) Printf(format string, args ...interface{}) {
	cl.Info(format, args...)
}

// Print logs a message at info level (for compatibility with log.Logger)
func (cl *ColoredLogger) Print(v ...interface{}) {
	cl.Info(fmt.Sprint(v...))
}

// SetOutput sets the output destination (for compatibility with log.Logger)
func (cl *ColoredLogger) SetOutput(w io.Writer) {
	cl.logger = log.New(w, "", 0)
}

// SetFlags sets the logger flags (for compatibility with log.Logger)
func (cl *ColoredLogger) SetFlags(flag int) {
	// We handle our own formatting, so this is a no-op
}

// SetPrefix sets the logger prefix (for compatibility with log.Logger)
func (cl *ColoredLogger) SetPrefix(prefix string) {
	// We handle our own formatting, so this is a no-op
}

// Writer returns the underlying writer (for compatibility with log.Logger)
func (cl *ColoredLogger) Writer() io.Writer {
	return os.Stdout
}

// NewA2ALogger creates a logger specifically for A2A server with a nice prefix
func NewA2ALogger() *ColoredLogger {
	logger := NewColoredLogger(LevelInfo)
	return logger
}

// NewAIPlayerLogger creates a logger specifically for AI player with a nice prefix
func NewAIPlayerLogger() *ColoredLogger {
	logger := NewColoredLogger(LevelInfo)
	return logger
}
