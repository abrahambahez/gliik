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
	metaFile := filepath.Join(instructionDir, "meta.yaml")

	metaData, err := os.ReadFile(metaFile)
	if err != nil {
		return "", fmt.Errorf("failed to read meta.yaml: %w", err)
	}

	var meta Meta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return "", fmt.Errorf("failed to parse meta.yaml: %w", err)
	}

	return meta.Version, nil
}

func BumpVersion(name, description string) (string, string, error) {
	if err := ValidateName(name); err != nil {
		return "", "", err
	}

	instructionDir := filepath.Join(config.GetInstructionsDir(), name)
	metaFile := filepath.Join(instructionDir, "meta.yaml")

	metaData, err := os.ReadFile(metaFile)
	if err != nil {
		return "", "", fmt.Errorf("failed to read meta.yaml: %w", err)
	}

	var meta Meta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return "", "", fmt.Errorf("failed to parse meta.yaml: %w", err)
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
		return "", "", fmt.Errorf("failed to marshal meta.yaml: %w", err)
	}

	if err := os.WriteFile(metaFile, newMetaData, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write meta.yaml: %w", err)
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
	metaFile := filepath.Join(instructionDir, "meta.yaml")

	metaData, err := os.ReadFile(metaFile)
	if err != nil {
		return "", fmt.Errorf("failed to read meta.yaml: %w", err)
	}

	var meta Meta
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return "", fmt.Errorf("failed to parse meta.yaml: %w", err)
	}

	oldVersion := meta.Version
	meta.Version = version
	if description != "" {
		meta.Description = description
	}

	newMetaData, err := yaml.Marshal(&meta)
	if err != nil {
		return "", fmt.Errorf("failed to marshal meta.yaml: %w", err)
	}

	if err := os.WriteFile(metaFile, newMetaData, 0644); err != nil {
		return "", fmt.Errorf("failed to write meta.yaml: %w", err)
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
