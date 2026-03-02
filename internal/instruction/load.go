package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/gliik/internal/config"
)

var validTypes = map[string]bool{"prompt": true, "skill": true, "context": true}

func Load(name string) (*Instruction, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)

	if _, err := os.Stat(instructionDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("instruction '%s' not found", name)
	}

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read instruction.md: %w", err)
	}

	meta, systemText, err := ParseFrontmatter(string(instructionData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse instruction.md: %w", err)
	}

	if meta.Type == "" {
		meta.Type = "prompt"
	}

	if !validTypes[meta.Type] {
		return nil, fmt.Errorf("instruction '%s' has invalid type '%s': must be prompt, skill, or context", name, meta.Type)
	}

	if meta.Type == "skill" && meta.Name == "" {
		return nil, fmt.Errorf("instruction '%s' of type 'skill' requires a 'name' field in frontmatter", name)
	}

	return &Instruction{
		Name:       name,
		Path:       instructionDir,
		SystemText: systemText,
		Meta:       meta,
	}, nil
}
