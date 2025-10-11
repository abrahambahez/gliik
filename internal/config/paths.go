package config

import (
	"os"
	"path/filepath"
)

func GetGliikHome() string {
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "gliik")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "gliik")
}

func GetInstructionsDir() string {
	cfg, err := Load()
	if err == nil && cfg.InstructionsDir != "" {
		return expandPath(cfg.InstructionsDir)
	}

	return filepath.Join(GetGliikHome(), "instructions")
}

func GetConfigFile() string {
	return filepath.Join(GetGliikHome(), "config.yaml")
}

func expandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
