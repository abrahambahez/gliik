package instruction

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/yourusername/gliik/internal/config"
	"gopkg.in/yaml.v3"
)

var validNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("instruction name cannot be empty")
	}
	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("instruction name must contain only alphanumeric characters and underscores")
	}
	return nil
}

func Create(name, description string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)

	if _, err := os.Stat(instructionDir); err == nil {
		return fmt.Errorf("instruction '%s' already exists", name)
	}

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		return fmt.Errorf("failed to create instruction directory: %w", err)
	}

	systemFile := filepath.Join(instructionDir, "system.txt")
	if err := os.WriteFile(systemFile, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create system.txt: %w", err)
	}

	meta := Meta{
		Version:     "0.1.0",
		Description: description,
	}

	metaData, err := yaml.Marshal(&meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta.yaml: %w", err)
	}

	metaFile := filepath.Join(instructionDir, "meta.yaml")
	if err := os.WriteFile(metaFile, metaData, 0644); err != nil {
		return fmt.Errorf("failed to write meta.yaml: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, systemFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	return nil
}
