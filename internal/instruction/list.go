package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourusername/gliik/internal/config"
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
		instructionFile := filepath.Join(instructionPath, "instruction.md")

		instructionData, err := os.ReadFile(instructionFile)
		if err != nil {
			continue
		}

		meta, _, err := ParseFrontmatter(string(instructionData))
		if err != nil {
			continue
		}

		if len(meta.Tags) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: instruction '%s' missing required field 'tags' in frontmatter\n", name)
		}

		if meta.Lang == "" {
			fmt.Fprintf(os.Stderr, "Warning: instruction '%s' missing required field 'lang' in frontmatter\n", name)
		}

		instructions = append(instructions, Instruction{
			Name: name,
			Path: instructionPath,
			Meta: meta,
		})
	}

	return instructions, nil
}
