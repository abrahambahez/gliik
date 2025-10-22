package instruction

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCreate_GeneratesValidInstructionMD(t *testing.T) {
	tmpDir := t.TempDir()

	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	os.Setenv("EDITOR", "true")

	instructionName := "test_instruction"
	description := "Test instruction for creation"
	tags := []string{"test", "example"}
	lang := "en"

	err := createWithCustomDir(instructionName, description, tags, lang, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	contentStr := string(content)

	if !strings.HasPrefix(contentStr, "---\n") {
		t.Error("instruction.md should start with ---")
	}

	if !strings.Contains(contentStr, "version:") || !strings.Contains(contentStr, "0.1.0") {
		t.Error("instruction.md should contain version: 0.1.0")
	}

	if !strings.Contains(contentStr, "description: "+description) {
		t.Errorf("instruction.md should contain description: %s", description)
	}

	if !strings.Contains(contentStr, "- test") {
		t.Error("instruction.md should contain tag 'test'")
	}

	if !strings.Contains(contentStr, "- example") {
		t.Error("instruction.md should contain tag 'example'")
	}

	if !strings.Contains(contentStr, "lang: "+lang) {
		t.Errorf("instruction.md should contain lang: %s", lang)
	}

	if !strings.Contains(contentStr, "# "+instructionName) {
		t.Errorf("instruction.md should contain markdown header: # %s", instructionName)
	}

	delimiterCount := strings.Count(contentStr, "---")
	if delimiterCount < 2 {
		t.Errorf("instruction.md should have at least 2 '---' delimiters, found %d", delimiterCount)
	}
}

func TestCreate_FrontmatterFormat(t *testing.T) {
	tmpDir := t.TempDir()

	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	os.Setenv("EDITOR", "true")

	instructionName := "format_test"
	description := "Test frontmatter format"
	tags := []string{"format"}
	lang := "es"

	err := createWithCustomDir(instructionName, description, tags, lang, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instructionFile := filepath.Join(tmpDir, instructionName, "instruction.md")
	content, err := os.ReadFile(instructionFile)
	if err != nil {
		t.Fatalf("failed to read instruction.md: %v", err)
	}

	meta, body, err := ParseFrontmatter(string(content))
	if err != nil {
		t.Fatalf("failed to parse frontmatter: %v", err)
	}

	if meta.Version != "0.1.0" {
		t.Errorf("expected version '0.1.0', got '%s'", meta.Version)
	}

	if meta.Description != description {
		t.Errorf("expected description '%s', got '%s'", description, meta.Description)
	}

	if len(meta.Tags) != 1 || meta.Tags[0] != "format" {
		t.Errorf("expected tags ['format'], got %v", meta.Tags)
	}

	if meta.Lang != lang {
		t.Errorf("expected lang '%s', got '%s'", lang, meta.Lang)
	}

	expectedBody := "# " + instructionName
	if !strings.HasPrefix(body, expectedBody) {
		t.Errorf("expected body to start with '%s', got '%s'", expectedBody, body)
	}
}

func TestCreate_LoadableInstruction(t *testing.T) {
	tmpDir := t.TempDir()

	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	os.Setenv("EDITOR", "true")

	instructionName := "loadable_test"
	description := "Test that created instruction is loadable"
	tags := []string{"test", "load"}
	lang := "en"

	err := createWithCustomDir(instructionName, description, tags, lang, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	instr, err := loadWithCustomDir(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("failed to load created instruction: %v", err)
	}

	if instr.Name != instructionName {
		t.Errorf("expected name '%s', got '%s'", instructionName, instr.Name)
	}

	if instr.Meta.Version != "0.1.0" {
		t.Errorf("expected version '0.1.0', got '%s'", instr.Meta.Version)
	}

	if instr.Meta.Description != description {
		t.Errorf("expected description '%s', got '%s'", description, instr.Meta.Description)
	}

	if len(instr.Meta.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(instr.Meta.Tags))
	}

	if instr.Meta.Lang != lang {
		t.Errorf("expected lang '%s', got '%s'", lang, instr.Meta.Lang)
	}
}

func createWithCustomDir(name, description string, tags []string, lang string, baseDir string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	if err := ValidateTags(tags); err != nil {
		return err
	}

	if err := ValidateLanguageCode(lang); err != nil {
		return err
	}

	instructionDir := filepath.Join(baseDir, name)

	if _, err := os.Stat(instructionDir); err == nil {
		return fmt.Errorf("instruction already exists")
	}

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		return err
	}

	meta := Meta{
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
	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		return err
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
		return err
	}

	return nil
}
