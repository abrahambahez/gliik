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

func ValidateLanguageCode(lang string) error {
	validLangRegex := regexp.MustCompile(`^[a-z]{2}$`)
	if !validLangRegex.MatchString(lang) {
		return fmt.Errorf("invalid language code '%s': must be ISO 639-1 two-letter lowercase code (e.g., 'en', 'es', 'fr')", lang)
	}
	return nil
}

func ValidateTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	validTagRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	for _, tag := range tags {
		if !validTagRegex.MatchString(tag) {
			return fmt.Errorf("invalid tag '%s': must be lowercase alphanumeric with hyphens only", tag)
		}
	}
	return nil
}

func Create(name, description string, tags []string, lang string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	if err := ValidateTags(tags); err != nil {
		return err
	}

	if err := ValidateLanguageCode(lang); err != nil {
		return err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)

	if _, err := os.Stat(instructionDir); err == nil {
		return fmt.Errorf("instruction '%s' already exists", name)
	}

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		return fmt.Errorf("failed to create instruction directory: %w", err)
	}

	meta := Meta{
		Version:     "0.1.0",
		Description: description,
		Tags:        tags,
		Lang:        lang,
	}

	metaData, err := yaml.Marshal(&meta)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	instructionContent := fmt.Sprintf("---\n%s---\n# %s\n\n", string(metaData), name)

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		return fmt.Errorf("failed to write instruction.md: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, instructionFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	return nil
}
