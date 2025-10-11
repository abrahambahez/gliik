package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultModel    string `yaml:"default_model"`
	Editor          string `yaml:"editor"`
	InstructionsDir string `yaml:"instructions_dir,omitempty"`
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
