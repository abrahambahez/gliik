package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	configFile := GetConfigFile()

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
