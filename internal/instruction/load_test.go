package instruction

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_ValidInstructionMD(t *testing.T) {
	tmpDir := t.TempDir()

	originalGetInstructionsDir := getInstructionsDirForTest
	defer func() { getInstructionsDirForTest = originalGetInstructionsDir }()

	getInstructionsDirForTest = func() string {
		return tmpDir
	}

	instructionName := "test_instruction"
	instructionDir := filepath.Join(tmpDir, instructionName)

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		t.Fatalf("failed to create instruction directory: %v", err)
	}

	instructionContent := `---
version: "1.0.0"
description: "Test instruction for loading"
tags:
  - test
  - example
lang: "en"
---
# Test Instruction

This is a test instruction with {{variable}} and {{input|text}}.

Process this content.`

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		t.Fatalf("failed to write instruction.md: %v", err)
	}

	instr, err := loadWithCustomDir(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instr.Name != instructionName {
		t.Errorf("expected name '%s', got '%s'", instructionName, instr.Name)
	}

	if instr.Meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", instr.Meta.Version)
	}

	if instr.Meta.Description != "Test instruction for loading" {
		t.Errorf("expected description 'Test instruction for loading', got '%s'", instr.Meta.Description)
	}

	if len(instr.Meta.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(instr.Meta.Tags))
	}

	if instr.Meta.Tags[0] != "test" || instr.Meta.Tags[1] != "example" {
		t.Errorf("expected tags 'test' and 'example', got '%s' and '%s'", instr.Meta.Tags[0], instr.Meta.Tags[1])
	}

	if instr.Meta.Lang != "en" {
		t.Errorf("expected lang 'en', got '%s'", instr.Meta.Lang)
	}

	expectedSystemText := `# Test Instruction

This is a test instruction with {{variable}} and {{input|text}}.

Process this content.`

	if instr.SystemText != expectedSystemText {
		t.Errorf("expected system text:\n%s\n\ngot:\n%s", expectedSystemText, instr.SystemText)
	}
}

func TestLoad_SystemTextExcludesFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_no_frontmatter"
	instructionDir := filepath.Join(tmpDir, instructionName)

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		t.Fatalf("failed to create instruction directory: %v", err)
	}

	instructionContent := `---
version: "2.0.0"
description: "Verify frontmatter exclusion"
tags:
  - security
  - privacy
lang: "es"
---
## CRITICAL: SystemText must NOT contain frontmatter

The LLM should only receive this markdown body.

Metadata like version, description, tags, and lang should be in Meta, not SystemText.`

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	if err := os.WriteFile(instructionFile, []byte(instructionContent), 0644); err != nil {
		t.Fatalf("failed to write instruction.md: %v", err)
	}

	instr, err := loadWithCustomDir(instructionName, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(instr.SystemText, "---") {
		t.Error("SystemText contains frontmatter delimiters '---', it should only contain the body")
	}

	if strings.Contains(instr.SystemText, "version:") {
		t.Error("SystemText contains 'version:' metadata, frontmatter leaked into body")
	}

	if strings.Contains(instr.SystemText, "description:") {
		t.Error("SystemText contains 'description:' metadata, frontmatter leaked into body")
	}

	if strings.Contains(instr.SystemText, "tags:") {
		t.Error("SystemText contains 'tags:' metadata, frontmatter leaked into body")
	}

	if strings.Contains(instr.SystemText, "lang:") {
		t.Error("SystemText contains 'lang:' metadata, frontmatter leaked into body")
	}

	if !strings.HasPrefix(instr.SystemText, "## CRITICAL") {
		t.Errorf("SystemText should start with markdown body, got: %s", instr.SystemText[:50])
	}

	if !strings.Contains(instr.SystemText, "The LLM should only receive this markdown body") {
		t.Error("SystemText missing expected body content")
	}
}

var getInstructionsDirForTest = func() string {
	return ""
}

func TestLoad_MissingInstructionMD(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "test_instruction"
	instructionDir := filepath.Join(tmpDir, instructionName)

	if err := os.MkdirAll(instructionDir, 0755); err != nil {
		t.Fatalf("failed to create instruction directory: %v", err)
	}

	_, err := loadWithCustomDir(instructionName, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing instruction.md, got nil")
	}

	if !strings.Contains(err.Error(), "instruction.md") {
		t.Errorf("expected error message to mention 'instruction.md', got: %v", err)
	}
}

func TestLoad_MissingInstructionDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	instructionName := "nonexistent_instruction"

	_, err := loadWithCustomDir(instructionName, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing instruction directory, got nil")
	}
}

func TestLoad_MalformedFrontmatter(t *testing.T) {
	t.Run("invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		instructionName := "test_invalid_yaml"
		instructionDir := filepath.Join(tmpDir, instructionName)

		if err := os.MkdirAll(instructionDir, 0755); err != nil {
			t.Fatalf("failed to create instruction directory: %v", err)
		}

		invalidContent := `---
version: "1.0.0"
description: invalid yaml: : :
tags:
  - bad
    - nested wrong
---
Body content`

		instructionFile := filepath.Join(instructionDir, "instruction.md")
		if err := os.WriteFile(instructionFile, []byte(invalidContent), 0644); err != nil {
			t.Fatalf("failed to write instruction.md: %v", err)
		}

		_, err := loadWithCustomDir(instructionName, tmpDir)
		if err == nil {
			t.Fatal("expected error for invalid YAML, got nil")
		}

		if !strings.Contains(err.Error(), "failed to parse") {
			t.Errorf("expected error about failed parsing, got: %v", err)
		}
	})

	t.Run("missing start delimiter", func(t *testing.T) {
		tmpDir := t.TempDir()
		instructionName := "test_no_start_delimiter"
		instructionDir := filepath.Join(tmpDir, instructionName)

		if err := os.MkdirAll(instructionDir, 0755); err != nil {
			t.Fatalf("failed to create instruction directory: %v", err)
		}

		invalidContent := `version: "1.0.0"
description: "No start delimiter"
tags:
  - test
Body content`

		instructionFile := filepath.Join(instructionDir, "instruction.md")
		if err := os.WriteFile(instructionFile, []byte(invalidContent), 0644); err != nil {
			t.Fatalf("failed to write instruction.md: %v", err)
		}

		_, err := loadWithCustomDir(instructionName, tmpDir)
		if err == nil {
			t.Fatal("expected error for missing start delimiter, got nil")
		}

		if !strings.Contains(err.Error(), "frontmatter") {
			t.Errorf("expected error about frontmatter, got: %v", err)
		}
	})

	t.Run("missing end delimiter", func(t *testing.T) {
		tmpDir := t.TempDir()
		instructionName := "test_no_end_delimiter"
		instructionDir := filepath.Join(tmpDir, instructionName)

		if err := os.MkdirAll(instructionDir, 0755); err != nil {
			t.Fatalf("failed to create instruction directory: %v", err)
		}

		invalidContent := `---
version: "1.0.0"
description: "No end delimiter"
tags:
  - test
Body content without delimiter`

		instructionFile := filepath.Join(instructionDir, "instruction.md")
		if err := os.WriteFile(instructionFile, []byte(invalidContent), 0644); err != nil {
			t.Fatalf("failed to write instruction.md: %v", err)
		}

		_, err := loadWithCustomDir(instructionName, tmpDir)
		if err == nil {
			t.Fatal("expected error for missing end delimiter, got nil")
		}

		if !strings.Contains(err.Error(), "frontmatter") {
			t.Errorf("expected error about frontmatter, got: %v", err)
		}
	})
}

func loadWithCustomDir(name string, baseDir string) (*Instruction, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	instructionDir := filepath.Join(baseDir, name)

	if _, err := os.Stat(instructionDir); os.IsNotExist(err) {
		return nil, err
	}

	instructionFile := filepath.Join(instructionDir, "instruction.md")
	instructionData, err := os.ReadFile(instructionFile)
	if err != nil {
		return nil, err
	}

	meta, systemText, err := ParseFrontmatter(string(instructionData))
	if err != nil {
		return nil, err
	}

	return &Instruction{
		Name:       name,
		Path:       instructionDir,
		SystemText: systemText,
		Meta:       meta,
	}, nil
}
