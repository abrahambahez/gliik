package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/gliik/internal/config"
	"gopkg.in/yaml.v3"
)

func Load(name string) (*Instruction, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)

	if _, err := os.Stat(instructionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("instruction '%s' not found", name)
	}

	metaFile := filepath.Join(instructionDir, "meta.yaml")
	metaData, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta.yaml: %w", err)
	}

	var meta Meta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse meta.yaml: %w", err)
	}

	systemFile := filepath.Join(instructionDir, "system.txt")
	systemText, err := os.ReadFile(systemFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read system.txt: %w", err)
	}

	return &Instruction{
		Name:       name,
		Path:       instructionDir,
		SystemText: string(systemText),
		Meta:       meta,
	}, nil
}
