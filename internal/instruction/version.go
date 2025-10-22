package instruction

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yourusername/gliik/internal/config"
	"gopkg.in/yaml.v3"
)

var semverRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func GetVersion(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)
	instructionFile := filepath.Join(instructionDir, "instruction.md")

	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return "", fmt.Errorf("failed to read instruction.md: %w", err)
	}

	meta, _, err := ParseFrontmatter(string(instructionData))
	if err != nil {
		return "", fmt.Errorf("failed to parse instruction.md frontmatter: %w", err)
	}

	return meta.Version, nil
}

func BumpVersion(name, description string) (string, string, error) {
	if err := ValidateName(name); err != nil {
		return "", "", err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)
	instructionFile := filepath.Join(instructionDir, "instruction.md")

	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read instruction.md: %w", err)
	}

	meta, body, err := ParseFrontmatter(string(instructionData))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse instruction.md frontmatter: %w", err)
	}

	oldVersion := meta.Version
	newVersion, err := bumpPatch(oldVersion)
	if err != nil {
		return "", "", err
	}

	meta.Version = newVersion
	if description != "" {
		meta.Description = description
	}

	newMetaData, err := yaml.Marshal(&meta)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	instructionContent := fmt.Sprintf("---\n%s---\n%s", string(newMetaData), body)

	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write instruction.md: %w", err)
	}

	return oldVersion, newVersion, nil
}

func SetVersion(name, version, description string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	if !semverRegex.MatchString(version) {
		return "", fmt.Errorf("invalid version format: must be X.Y.Z (e.g., 1.0.0)")
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)
	instructionFile := filepath.Join(instructionDir, "instruction.md")

	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return "", fmt.Errorf("failed to read instruction.md: %w", err)
	}

	meta, body, err := ParseFrontmatter(string(instructionData))
	if err != nil {
		return "", fmt.Errorf("failed to parse instruction.md frontmatter: %w", err)
	}

	oldVersion := meta.Version
	meta.Version = version
	if description != "" {
		meta.Description = description
	}

	newMetaData, err := yaml.Marshal(&meta)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	instructionContent := fmt.Sprintf("---\n%s---\n%s", string(newMetaData), body)

	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write instruction.md: %w", err)
	}

	return oldVersion, nil
}

func bumpPatch(version string) (string, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid version format: %s", version)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch+1), nil
}
