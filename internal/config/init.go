package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AnthropicConfig holds configuration for the Anthropic provider.
type AnthropicConfig struct {
	Model string `yaml:"model"`
}

// OllamaConfig holds configuration for the Ollama provider.
type OllamaConfig struct {
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

// OpenAIConfig holds configuration for the OpenAI provider.
type OpenAIConfig struct {
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

// Config represents the Gliik configuration file structure.
type Config struct {
	DefaultModel    string `yaml:"default_model"`
	Editor          string `yaml:"editor"`
	InstructionsDir string `yaml:"instructions_dir,omitempty"`
	// Provider specifies the LLM provider to use for instruction execution.
	// Valid values: "anthropic" (default) or "ollama".
	Provider  string          `yaml:"provider"`
	Anthropic AnthropicConfig `yaml:"anthropic"`
	Ollama    OllamaConfig    `yaml:"ollama"`
	OpenAI    OpenAIConfig    `yaml:"openai"`
}

// ValidateProvider checks if the provider value is either "anthropic", "ollama", or "openai".
func (c *Config) ValidateProvider() error {
	if c.Provider != "anthropic" && c.Provider != "ollama" && c.Provider != "openai" {
		return fmt.Errorf("invalid provider '%s': must be 'anthropic', 'ollama', or 'openai'", c.Provider)
	}
	return nil
}

func Initialize(instructionsDir string) error {
	gliikHome := GetGliikHome()
	configFile := GetConfigFile()

	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("Gliik is already initialized at %s", gliikHome)
	}

	if err := os.MkdirAll(gliikHome, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	defaultConfig := Config{
		DefaultModel:    "claude-sonnet-4-20250514",
		Editor:          "vim",
		InstructionsDir: instructionsDir,
		Provider:        "anthropic",
		Anthropic: AnthropicConfig{
			Model: "claude-sonnet-4-20250514",
		},
		Ollama: OllamaConfig{
			Endpoint: "http://localhost:11434",
			Model:    "llama3.2",
		},
		OpenAI: OpenAIConfig{
			Endpoint: "https://api.openai.com/v1",
			Model:    "gpt-4o-mini",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	actualDir := GetInstructionsDir()
	if err := os.MkdirAll(actualDir, 0755); err != nil {
		return fmt.Errorf("failed to create instructions directory: %w", err)
	}

	return nil
}
