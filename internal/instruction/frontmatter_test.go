package instruction

import (
	"strings"
	"testing"
)

func TestParseFrontmatter_ValidComplete(t *testing.T) {
	content := `---
version: "1.0.0"
description: "Test instruction"
tags:
  - tag1
  - tag2
lang: "en"
---
# Test Instruction

This is the markdown body with {{variable}}.`

	meta, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}

	if meta.Description != "Test instruction" {
		t.Errorf("expected description 'Test instruction', got '%s'", meta.Description)
	}

	if len(meta.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(meta.Tags))
	}

	if meta.Tags[0] != "tag1" || meta.Tags[1] != "tag2" {
		t.Errorf("expected tags 'tag1' and 'tag2', got '%s' and '%s'", meta.Tags[0], meta.Tags[1])
	}

	if meta.Lang != "en" {
		t.Errorf("expected lang 'en', got '%s'", meta.Lang)
	}

	expectedBody := "# Test Instruction\n\nThis is the markdown body with {{variable}}."
	if body != expectedBody {
		t.Errorf("expected body '%s', got '%s'", expectedBody, body)
	}
}

func TestParseFrontmatter_ValidMinimal(t *testing.T) {
	content := `---
version: "1.0.0"
description: "Minimal test"
---
Body content here.`

	meta, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}

	if meta.Description != "Minimal test" {
		t.Errorf("expected description 'Minimal test', got '%s'", meta.Description)
	}

	if len(meta.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(meta.Tags))
	}

	if meta.Lang != "" {
		t.Errorf("expected empty lang, got '%s'", meta.Lang)
	}

	if body != "Body content here." {
		t.Errorf("expected body 'Body content here.', got '%s'", body)
	}
}

func TestParseFrontmatter_MissingStartDelimiter(t *testing.T) {
	content := `version: "1.0.0"
description: "No delimiters"
Body content here`

	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for missing start delimiter, got nil")
	}

	if !strings.Contains(err.Error(), "missing frontmatter start delimiter") {
		t.Errorf("expected error about missing start delimiter, got: %v", err)
	}
}

func TestParseFrontmatter_MissingEndDelimiter(t *testing.T) {
	content := `---
version: "1.0.0"
description: "No end delimiter"
Body without delimiter`

	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for missing end delimiter, got nil")
	}

	if !strings.Contains(err.Error(), "missing frontmatter end delimiter") {
		t.Errorf("expected error about missing end delimiter, got: %v", err)
	}
}

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	content := `---
version: "1.0.0"
description: invalid yaml content here: : :
tags:
  - malformed
    - nested wrong
---
Body`

	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse frontmatter YAML") {
		t.Errorf("expected error about failed YAML parsing, got: %v", err)
	}
}

func TestParseFrontmatter_MarkdownBodyWithVariables(t *testing.T) {
	content := `---
version: "1.0.0"
description: "Test with variables"
---
# Header with {{title}}

Process this {{input|text}} content.

- List item with {{variable}}
- Another item

**Bold {{text}}** and *italic {{name}}*.

` + "```" + `
code block with {{code_var}}
` + "```"

	meta, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}

	if !strings.Contains(body, "{{title}}") {
		t.Error("expected body to contain {{title}}")
	}

	if !strings.Contains(body, "{{input|text}}") {
		t.Error("expected body to contain {{input|text}}")
	}

	if !strings.Contains(body, "{{variable}}") {
		t.Error("expected body to contain {{variable}}")
	}

	if !strings.Contains(body, "{{text}}") {
		t.Error("expected body to contain {{text}}")
	}

	if !strings.Contains(body, "{{name}}") {
		t.Error("expected body to contain {{name}}")
	}

	if !strings.Contains(body, "{{code_var}}") {
		t.Error("expected body to contain {{code_var}}")
	}

	if !strings.Contains(body, "# Header") {
		t.Error("expected body to contain markdown header")
	}

	if !strings.Contains(body, "**Bold") {
		t.Error("expected body to contain markdown bold")
	}
}

func TestParseFrontmatter_EmptyBody(t *testing.T) {
	content := `---
version: "1.0.0"
description: "Empty body test"
---`

	meta, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}

	if body != "" {
		t.Errorf("expected empty body, got '%s'", body)
	}
}

func TestParseFrontmatter_BodyWithLeadingNewlines(t *testing.T) {
	content := `---
version: "1.0.0"
description: "Test"
---


Body after newlines`

	_, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if body != "Body after newlines" {
		t.Errorf("expected body 'Body after newlines', got '%s'", body)
	}
}
