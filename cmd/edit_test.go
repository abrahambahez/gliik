package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/gliik/internal/instruction"
	"gopkg.in/yaml.v3"
)

func TestEditCommand_InstructionMDPath(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_edit"
	description := "Test instruction for editing"
	tags := []string{"test"}
	lang := "en"

	err := createTestInstruction(instructionName, description, tags, lang, tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	inst, err := loadTestInstruction(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("failed to load instruction: %v", err)
	}

	instructionFile := filepath.Join(inst.Path, "instruction.md")

	if _, err := os.Stat(instructionFile); os.IsNotExist(err) {
		t.Errorf("instruction.md file does not exist at expected path: %s", instructionFile)
	}

	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	if !strings.Contains(string(content), "---") {
		t.Error("instruction.md should contain frontmatter delimiters")
	}

	if !strings.Contains(string(content), description) {
		t.Error("instruction.md should contain the description")
	}
}

func TestEditCommand_ChangesPersist(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_persist"
	description := "Original description"
	tags := []string{"test"}
	lang := "en"

	err := createTestInstruction(instructionName, description, tags, lang, tmpDir)
	if err != nil {
		t.Fatalf("failed to create test instruction: %v", err)
	}

	inst, err := loadTestInstruction(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("failed to load instruction: %v", err)
	}

	instructionFile := filepath.Join(inst.Path, "instruction.md")

	originalContent, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	modifiedContent := strings.ReplaceAll(string(originalContent), "# "+instructionName, "# Modified Header\n\nNew content added.")

	err = os.WriteFile(instructionFile, []byte(modifiedContent), 0644)
	if err != nil {
		t.Fatalf("failed to write modified content: %v", err)
	}

	persistedContent, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md after modification: %v", err)
	}

	if !strings.Contains(string(persistedContent), "Modified Header") {
		t.Error("changes to instruction.md should persist")
	}

	if !strings.Contains(string(persistedContent), "New content added") {
		t.Error("new content should be persisted in instruction.md")
	}

	reloaded, err := loadTestInstruction(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("failed to reload instruction after modification: %v", err)
	}

	if !strings.Contains(reloaded.SystemText, "Modified Header") {
		t.Error("modified content should be loaded by Load() function")
	}
}

func createTestInstruction(name, description string, tags []string, lang string, baseDir string) error {
	instructionDir := filepath.Join(baseDir, name)

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		return err
	}

	meta := instruction.Meta{
		Version:     "0.1.0",
		Description: description,
		Tags:        tags,
		Lang:        lang,
	}

	metaData, err := yaml.Marshal(&meta)
	if err != nil {
		return err
	}

	instructionContent := fmt.Sprintf("---\n%s---\n# %s\n\n", string(metaData), name)

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	return os.WriteFile(instructionFile, []byte(instructionContent), 0644)
}

func loadTestInstruction(name string, baseDir string) (*instruction.Instruction, error) {
	instructionDir := filepath.Join(baseDir, name)

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return nil, err
	}

	meta, systemText, err := instruction.ParseFrontmatter(string(instructionData))
	if err != nil {
		return nil, err
	}

	return &instruction.Instruction{
		Name:       name,
		Path:       instructionDir,
		SystemText: systemText,
		Meta:       meta,
	}, nil
}
