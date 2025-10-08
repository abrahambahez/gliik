package instruction

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolver_Stdin(t *testing.T) {
	variables := []Variable{
		{Raw: "{{input}}", Options: []string{"input"}},
	}

	resolver := Resolver{
		Variables: variables,
		Stdin:     "test input from stdin",
		Flags:     map[string]string{},
	}

	resolved, err := resolver.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved["{{input}}"] != "test input from stdin" {
		t.Errorf("expected 'test input from stdin', got '%s'", resolved["{{input}}"])
	}
}

func TestResolver_CLIFlag(t *testing.T) {
	variables := []Variable{
		{Raw: "{{text}}", Options: []string{"text"}},
	}

	resolver := Resolver{
		Variables: variables,
		Stdin:     "",
		Flags: map[string]string{
			"text": "value from flag",
		},
	}

	resolved, err := resolver.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved["{{text}}"] != "value from flag" {
		t.Errorf("expected 'value from flag', got '%s'", resolved["{{text}}"])
	}
}

func TestResolver_ORPriority(t *testing.T) {
	variables := []Variable{
		{Raw: "{{input|text}}", Options: []string{"input", "text"}},
	}

	t.Run("stdin wins when both available", func(t *testing.T) {
		resolver := Resolver{
			Variables: variables,
			Stdin:     "from stdin",
			Flags: map[string]string{
				"text": "from flag",
			},
		}

		resolved, err := resolver.Resolve()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resolved["{{input|text}}"] != "from stdin" {
			t.Errorf("expected 'from stdin', got '%s'", resolved["{{input|text}}"])
		}
	})

	t.Run("flag wins when stdin not available", func(t *testing.T) {
		resolver := Resolver{
			Variables: variables,
			Stdin:     "",
			Flags: map[string]string{
				"text": "from flag",
			},
		}

		resolved, err := resolver.Resolve()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resolved["{{input|text}}"] != "from flag" {
			t.Errorf("expected 'from flag', got '%s'", resolved["{{input|text}}"])
		}
	})
}

func TestResolver_MissingVariable(t *testing.T) {
	t.Run("simple variable", func(t *testing.T) {
		variables := []Variable{
			{Raw: "{{job_description}}", Options: []string{"job_description"}},
		}

		resolver := Resolver{
			Variables: variables,
			Stdin:     "",
			Flags:     map[string]string{},
		}

		_, err := resolver.Resolve()
		if err == nil {
			t.Fatal("expected error for missing variable")
		}

		if !strings.Contains(err.Error(), "missing required variable") {
			t.Errorf("expected 'missing required variable' in error, got: %v", err)
		}
	})

	t.Run("OR variable", func(t *testing.T) {
		variables := []Variable{
			{Raw: "{{input|text}}", Options: []string{"input", "text"}},
		}

		resolver := Resolver{
			Variables: variables,
			Stdin:     "",
			Flags:     map[string]string{},
		}

		_, err := resolver.Resolve()
		if err == nil {
			t.Fatal("expected error for missing variable")
		}

		if !strings.Contains(err.Error(), "needs one of") {
			t.Errorf("expected 'needs one of' in error, got: %v", err)
		}
	})
}

func TestResolver_StdinRejection(t *testing.T) {
	variables := []Variable{
		{Raw: "{{text}}", Options: []string{"text"}},
	}

	resolver := Resolver{
		Variables: variables,
		Stdin:     "stdin provided",
		Flags:     map[string]string{},
	}

	_, err := resolver.Resolve()
	if err == nil {
		t.Fatal("expected error for stdin to non-stdin instruction")
	}

	if !strings.Contains(err.Error(), "does not accept stdin") {
		t.Errorf("expected 'does not accept stdin' in error, got: %v", err)
	}
}

func TestResolver_FileResolution(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("content from file"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	variables := []Variable{
		{Raw: "{{text}}", Options: []string{"text"}},
	}

	resolver := Resolver{
		Variables: variables,
		Stdin:     "",
		Flags: map[string]string{
			"text": testFile,
		},
	}

	resolved, err := resolver.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolved["{{text}}"] != "content from file" {
		t.Errorf("expected 'content from file', got '%s'", resolved["{{text}}"])
	}
}
