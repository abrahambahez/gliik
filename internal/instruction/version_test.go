package instruction

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBumpVersion_UpdatesFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_bump"
	description := "Original description"
	tags := []string{"test"}
	lang := "en"

	err := createVersionTestInstruction(instructionName, description, tags, lang, "1.0.0", tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	oldVersion, newVersion, err := bumpVersionWithCustomDir(instructionName, "", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if oldVersion != "1.0.0" {
		t.Errorf("expected old version '1.0.0', got '%s'", oldVersion)
	}

	if newVersion != "1.0.1" {
		t.Errorf("expected new version '1.0.1', got '%s'", newVersion)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	meta, _, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if meta.Version != "1.0.1" {
		t.Errorf("expected version '1.0.1' in frontmatter, got '%s'", meta.Version)
	}
}

func TestBumpVersion_UpdatesDescription(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_bump_desc"
	originalDesc := "Original description"
	tags := []string{"test"}
	lang := "en"

	err := createVersionTestInstruction(instructionName, originalDesc, tags, lang, "1.0.0", tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	newDesc := "Updated description"
	_, _, err = bumpVersionWithCustomDir(instructionName, newDesc, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	meta, _, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if meta.Description != newDesc {
		t.Errorf("expected description '%s', got '%s'", newDesc, meta.Description)
	}
}

func TestSetVersion_UpdatesFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_set"
	description := "Test description"
	tags := []string{"test"}
	lang := "en"

	err := createVersionTestInstruction(instructionName, description, tags, lang, "1.0.0", tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	oldVersion, err := setVersionWithCustomDir(instructionName, "2.5.3", "", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if oldVersion != "1.0.0" {
		t.Errorf("expected old version '1.0.0', got '%s'", oldVersion)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	meta, _, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if meta.Version != "2.5.3" {
		t.Errorf("expected version '2.5.3' in frontmatter, got '%s'", meta.Version)
	}
}

func TestBumpVersion_PreservesMarkdownBody(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_preserve_bump"
	description := "Test description"
	tags := []string{"test"}
	lang := "en"
	originalBody := "# Test Instruction\n\nThis is the original body with {{variable}}.\n\n- Item 1\n- Item 2"

	err := createVersionTestInstructionWithBody(instructionName, description, tags, lang, "1.0.0", originalBody, tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	_, _, err = bumpVersionWithCustomDir(instructionName, "", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	_, body, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if body != originalBody {
		t.Errorf("markdown body was not preserved.\nExpected:\n%s\n\nGot:\n%s", originalBody, body)
	}
}

func TestSetVersion_PreservesMarkdownBody(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_preserve_set"
	description := "Test description"
	tags := []string{"test"}
	lang := "en"
	originalBody := "# Test Instruction\n\nOriginal content with variables {{input|text}}.\n\n**Bold text** and *italic*."

	err := createVersionTestInstructionWithBody(instructionName, description, tags, lang, "1.0.0", originalBody, tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	_, err = setVersionWithCustomDir(instructionName, "3.2.1", "", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	_, body, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if body != originalBody {
		t.Errorf("markdown body was not preserved.\nExpected:\n%s\n\nGot:\n%s", originalBody, body)
	}
}

func TestGetVersion_ReadsFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_get"
	description := "Test description"
	tags := []string{"test"}
	lang := "en"

	err := createVersionTestInstruction(instructionName, description, tags, lang, "5.4.3", tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	version, err := getVersionWithCustomDir(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if version != "5.4.3" {
		t.Errorf("expected version '5.4.3', got '%s'", version)
	}
}

func TestBumpVersion_MultipleBumps(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_multiple_bumps"
	description := "Test description"
	tags := []string{"test"}
	lang := "en"

	err := createVersionTestInstruction(instructionName, description, tags, lang, "1.0.0", tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	_, v1, err := bumpVersionWithCustomDir(instructionName, "", tmpDir)
	if err != nil {
		t.Fatalf("first bump failed: %v", err)
	}
	if v1 != "1.0.1" {
		t.Errorf("expected version '1.0.1' after first bump, got '%s'", v1)
	}

	_, v2, err := bumpVersionWithCustomDir(instructionName, "", tmpDir)
	if err != nil {
		t.Fatalf("second bump failed: %v", err)
	}
	if v2 != "1.0.2" {
		t.Errorf("expected version '1.0.2' after second bump, got '%s'", v2)
	}

	_, v3, err := bumpVersionWithCustomDir(instructionName, "", tmpDir)
	if err != nil {
		t.Fatalf("third bump failed: %v", err)
	}
	if v3 != "1.0.3" {
		t.Errorf("expected version '1.0.3' after third bump, got '%s'", v3)
	}
}

func createVersionTestInstruction(name, description string, tags []string, lang, version string, baseDir string) error {
	return createVersionTestInstructionWithBody(name, description, tags, lang, version, "# "+name+"\n\n", baseDir)
}

func createVersionTestInstructionWithBody(name, description string, tags []string, lang, version, body string, baseDir string) error {
	instructionDir := filepath.Join(baseDir, name)

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		return err
	}

	meta := Meta{
		Version:     version,
		Description: description,
		Tags:        tags,
		Lang:        lang,
	}

	metaData, err := yaml.Marshal(&meta)
	if err != nil {
		return err
	}

	instructionContent := fmt.Sprintf("---\n%s---\n%s", string(metaData), body)

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	return os.WriteFile(instructionFile, []byte(instructionContent), 0644)
}

func getVersionWithCustomDir(name string, baseDir string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	instructionDir := filepath.Join(baseDir, name)
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

func bumpVersionWithCustomDir(name, description string, baseDir string) (string, string, error) {
	if err := ValidateName(name); err != nil {
		return "", "", err
	}

	instructionDir := filepath.Join(baseDir, name)
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

func setVersionWithCustomDir(name, version, description string, baseDir string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}

	if !semverRegex.MatchString(version) {
		return "", fmt.Errorf("invalid version format: must be X.Y.Z (e.g., 1.0.0)")
	}

	instructionDir := filepath.Join(baseDir, name)
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
