package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/gliik/internal/config"
	"gopkg.in/yaml.v3"
)

func ListAll() ([]Instruction, error) {
	instructionsDir := config.GetInstructionsDir()

	entries, err := os.ReadDir(instructionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read instructions directory: %w", err)
	}

	var instructions []Instruction

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		instructionPath := filepath.Join(instructionsDir, name)
		metaFile := filepath.Join(instructionPath, "meta.yaml")

		metaData, err := os.ReadFile(metaFile)
		if err != nil {
			continue
		}

		var meta Meta
		if err := yaml.Unmarshal(metaData, &meta); err != nil {
			continue
		}

		instructions = append(instructions, Instruction{
			Name: name,
			Path: instructionPath,
			Meta: meta,
		})
	}

	return instructions, nil
}
