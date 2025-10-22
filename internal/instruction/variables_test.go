package instruction

import (
	"strings"
	"testing"
)

func TestParseVariables_Simple(t *testing.T) {
	text := "Process this {{text}}"
	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 1 {
		t.Fatalf("expected 1 variable, got %d", len(vars))
	}

	if vars[0].Raw != "{{text}}" {
		t.Errorf("expected raw '{{text}}', got '%s'", vars[0].Raw)
	}

	if len(vars[0].Options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(vars[0].Options))
	}

	if vars[0].Options[0] != "text" {
		t.Errorf("expected option 'text', got '%s'", vars[0].Options[0])
	}
}

func TestParseVariables_OR(t *testing.T) {
	text := "Process {{input|text}}"
	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 1 {
		t.Fatalf("expected 1 variable, got %d", len(vars))
	}

	if vars[0].Raw != "{{input|text}}" {
		t.Errorf("expected raw '{{input|text}}', got '%s'", vars[0].Raw)
	}

	if len(vars[0].Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(vars[0].Options))
	}

	if vars[0].Options[0] != "input" {
		t.Errorf("expected first option 'input', got '%s'", vars[0].Options[0])
	}

	if vars[0].Options[1] != "text" {
		t.Errorf("expected second option 'text', got '%s'", vars[0].Options[1])
	}
}

func TestParseVariables_Multiple(t *testing.T) {
	text := "Process {{job}} and {{input|resume}}"
	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 2 {
		t.Fatalf("expected 2 variables, got %d", len(vars))
	}

	if vars[0].Raw != "{{job}}" {
		t.Errorf("expected first raw '{{job}}', got '%s'", vars[0].Raw)
	}

	if vars[1].Raw != "{{input|resume}}" {
		t.Errorf("expected second raw '{{input|resume}}', got '%s'", vars[1].Raw)
	}

	if len(vars[0].Options) != 1 {
		t.Errorf("expected 1 option for first variable, got %d", len(vars[0].Options))
	}

	if len(vars[1].Options) != 2 {
		t.Errorf("expected 2 options for second variable, got %d", len(vars[1].Options))
	}
}

func TestParseVariables_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected int
	}{
		{"empty text", "", 0},
		{"no variables", "just plain text", 0},
		{"spaces in options", "{{ input | text }}", 1},
		{"multiple spaces", "{{  input  |  text  }}", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vars, err := ParseVariables(tc.text)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(vars) != tc.expected {
				t.Errorf("expected %d variables, got %d", tc.expected, len(vars))
			}

			if tc.name == "spaces in options" && len(vars) > 0 {
				if len(vars[0].Options) != 2 {
					t.Errorf("expected 2 options, got %d", len(vars[0].Options))
				}
				if vars[0].Options[0] != "input" || vars[0].Options[1] != "text" {
					t.Errorf("expected trimmed options 'input' and 'text', got '%s' and '%s'",
						vars[0].Options[0], vars[0].Options[1])
				}
			}
		})
	}
}

func TestParseVariables_Duplicate(t *testing.T) {
	text := "Process {{text}} and then {{text}} again"
	_, err := ParseVariables(text)

	if err == nil {
		t.Fatal("expected error for duplicate variable, got nil")
	}

	expectedMsg := "duplicate variable in instruction: {{text}}"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestParseVariables_MarkdownHeaders(t *testing.T) {
	text := `# Header with {{title}}
## Subheader with {{input|context}}
### Process {{variable}}`

	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(vars))
	}

	if vars[0].Raw != "{{title}}" {
		t.Errorf("expected first variable '{{title}}', got '%s'", vars[0].Raw)
	}

	if vars[1].Raw != "{{input|context}}" {
		t.Errorf("expected second variable '{{input|context}}', got '%s'", vars[1].Raw)
	}

	if vars[2].Raw != "{{variable}}" {
		t.Errorf("expected third variable '{{variable}}', got '%s'", vars[2].Raw)
	}
}

func TestParseVariables_MarkdownLists(t *testing.T) {
	text := `- Item with {{var1}}
- Another item with {{var2|option}}
1. Numbered item {{var3}}
2. Second numbered {{input}}`

	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 4 {
		t.Fatalf("expected 4 variables, got %d", len(vars))
	}
}

func TestParseVariables_MarkdownFormatting(t *testing.T) {
	text := `**Bold text with {{bold_var}}**
*Italic with {{italic_var}}*
***Bold italic {{combined}}***
` + "`code with {{code_var}}`"

	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 4 {
		t.Fatalf("expected 4 variables, got %d", len(vars))
	}

	expectedVars := []string{"{{bold_var}}", "{{italic_var}}", "{{combined}}", "{{code_var}}"}
	for i, expected := range expectedVars {
		if vars[i].Raw != expected {
			t.Errorf("expected variable %d to be '%s', got '%s'", i, expected, vars[i].Raw)
		}
	}
}

func TestParseVariables_MarkdownCodeBlocks(t *testing.T) {
	text := "```\ncode block with {{var1}}\n```\n\nRegular text {{var2}}\n\n```python\nfunction({{param}})\n```"

	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(vars))
	}

	if vars[0].Raw != "{{var1}}" {
		t.Errorf("expected first variable '{{var1}}', got '%s'", vars[0].Raw)
	}

	if vars[1].Raw != "{{var2}}" {
		t.Errorf("expected second variable '{{var2}}', got '%s'", vars[1].Raw)
	}

	if vars[2].Raw != "{{param}}" {
		t.Errorf("expected third variable '{{param}}', got '%s'", vars[2].Raw)
	}
}

func TestParseVariables_ComplexMarkdown(t *testing.T) {
	text := `## PURPOSE

You are an assistant for {{role}}.

## PROCESS

Follow these steps:

1. Analyze the {{input|request}}
2. Generate output for **{{output_type}}**
3. Use ` + "`{{format}}`" + ` formatting

### Notes

- Important: *{{note}}*
- See {{reference|doc}}`

	vars, err := ParseVariables(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vars) != 6 {
		t.Fatalf("expected 6 variables, got %d", len(vars))
	}

	expectedVars := []string{"{{role}}", "{{input|request}}", "{{output_type}}", "{{format}}", "{{note}}", "{{reference|doc}}"}
	for i, expected := range expectedVars {
		if vars[i].Raw != expected {
			t.Errorf("expected variable %d to be '%s', got '%s'", i, expected, vars[i].Raw)
		}
	}
}
