package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/gliik/internal/config"
)

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

	if len(meta.Tags) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: instruction '%s' missing required field 'tags' in frontmatter\n", name)
	}

	if meta.Lang == "" {
		fmt.Fprintf(os.Stderr, "Warning: instruction '%s' missing required field 'lang' in frontmatter\n", name)
	}

	return &Instruction{
		Name:       name,
		Path:       instructionDir,
		SystemText: systemText,
		Meta:       meta,
	}, nil
}
