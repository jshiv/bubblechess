package ai_player

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the AI player
type Config struct {
	OllamaURL     string            `json:"ollama_url"`
	Model         string            `json:"model"`
	Timeout       int               `json:"timeout_seconds"`
	Temperature   float64           `json:"temperature"`
	TopP          float64           `json:"top_p"`
	MaxRetries    int               `json:"max_retries"`
	RetryDelay    int               `json:"retry_delay_seconds"`
	MoveHistory   int               `json:"move_history_length"`
	CustomPrompts map[string]string `json:"custom_prompts,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		OllamaURL:     "http://localhost:11434",
		Model:         "llama3.2:3b",
		Timeout:       30,
		Temperature:   0.1,
		TopP:          0.9,
		MaxRetries:    3,
		RetryDelay:    2,
		MoveHistory:   5,
		CustomPrompts: make(map[string]string),
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "ai_config.json"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		config := DefaultConfig()
		if err := SaveConfig(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Load existing config
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := DefaultConfig()
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = "ai_config.json"
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// ValidateConfig validates the configuration
func (c *Config) ValidateConfig() error {
	if c.OllamaURL == "" {
		return fmt.Errorf("ollama_url cannot be empty")
	}

	if c.Model == "" {
		return fmt.Errorf("model cannot be empty")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.Temperature < 0 || c.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if c.TopP < 0 || c.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1")
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries cannot be negative")
	}

	if c.RetryDelay < 0 {
		return fmt.Errorf("retry_delay cannot be negative")
	}

	if c.MoveHistory < 0 {
		return fmt.Errorf("move_history_length cannot be negative")
	}

	return nil
}
