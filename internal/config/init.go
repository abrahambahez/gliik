package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultModel string `yaml:"default_model"`
	Editor       string `yaml:"editor"`
}

func Initialize() error {
	gliikHome := GetGliikHome()
	instructionsDir := GetInstructionsDir()
	configFile := GetConfigFile()

	if _, err := os.Stat(gliikHome); err == nil {
		return fmt.Errorf("Gliik is already initialized at %s", gliikHome)
	}

	if err := os.MkdirAll(instructionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create instructions directory: %w", err)
	}

	defaultConfig := Config{
		DefaultModel: "claude-sonnet-4-20250514",
		Editor:       "vim",
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
