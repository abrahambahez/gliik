package config

import (
	"os"
	"path/filepath"
)

func GetGliikHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".gliik")
}

func GetInstructionsDir() string {
	return filepath.Join(GetGliikHome(), "instructions")
}

func GetConfigFile() string {
	return filepath.Join(GetGliikHome(), "config.yaml")
}
