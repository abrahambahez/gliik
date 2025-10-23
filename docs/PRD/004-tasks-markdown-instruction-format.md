# Tasks: Markdown Instruction Format Refactoring

Based on PRD: `004-markdown-instruction-format.md`

## Relevant Files

- `internal/instruction/instruction.go` - Contains Meta struct and Instruction struct definitions
- `internal/instruction/frontmatter.go` - ✓ CREATED: Frontmatter parsing utilities for instruction.md files
- `internal/instruction/frontmatter_test.go` - ✓ CREATED: Unit tests for frontmatter parsing (all tests passing)
- `internal/instruction/load.go` - ✓ UPDATED: Loads instructions from disk using instruction.md with frontmatter parsing
- `internal/instruction/load_test.go` - ✓ CREATED: Tests for Load function with instruction.md format
- `internal/instruction/create.go` - ✓ UPDATED: Creates new instructions using instruction.md with frontmatter template
- `internal/instruction/create_test.go` - ✓ CREATED: Tests for Create function with instruction.md format (all tests passing)
- `cmd/edit.go` - ✓ UPDATED: Edit command opens instruction.md file in editor (removed --meta flag)
- `cmd/edit_test.go` - ✓ CREATED: Tests for edit command verifying instruction.md path and persistence (all tests passing)
- `internal/instruction/version.go` - ✓ UPDATED: Version management reads/writes frontmatter and preserves markdown body
- `internal/instruction/version_test.go` - ✓ CREATED: Tests for version management (all tests passing)
- `internal/instruction/list.go` - ✓ UPDATED: ListAll() reads from instruction.md instead of meta.yaml (breaking change)
- `internal/instruction/variables.go` - Variable parsing logic (works unchanged with markdown)
- `internal/instruction/variables_test.go` - ✓ UPDATED: Added comprehensive markdown scenario tests (all passing)
- `internal/instruction/load_test.go` - ✓ UPDATED: Added test to verify frontmatter exclusion from SystemText

### Notes

- The variable substitution logic in `variables.go` should work unchanged since it operates on the SystemText string regardless of source format
- Tests should verify that markdown formatting (headers, lists, etc.) doesn't break variable substitution
- Frontmatter parsing can use simple string splitting on `---` delimiters followed by `gopkg.in/yaml.v3` parsing

## Edge Cases & Key Findings

### Frontmatter Security
- **VERIFIED**: SystemText (sent to LLM) contains ONLY markdown body, never frontmatter metadata
- Test `TestLoad_SystemTextExcludesFrontmatter` ensures no metadata leakage (version, description, tags, lang)
- Frontmatter delimiters (`---`) are completely stripped from LLM input

### Variable Substitution in Markdown
- **Works seamlessly** across all markdown contexts:
  - Headers (`# {{title}}`, `## {{subtitle}}`)
  - Lists (`- Item {{var}}`, `1. {{numbered}}`)
  - Bold/Italic text (`**{{bold}}**`, `*{{italic}}*`)
  - Inline code (`` `{{code}}` ``)
  - Code blocks (``` {{var}} ```)
  - Complex nested formatting
- Variables are format-agnostic - regex pattern `{{...}}` works regardless of surrounding markdown

### Breaking Changes
- **ListAll()** now only displays instruction.md format instructions
- Old format instructions (system.txt + meta.yaml) will not appear in `gliik list`
- Old format instructions will fail to load/edit/version with updated functions
- Users must migrate existing instructions to new format

### Migration Requirements
- Convert `system.txt` → markdown body in `instruction.md`
- Convert `meta.yaml` → YAML frontmatter in `instruction.md`
- Format: `---\n<yaml frontmatter>\n---\n<markdown body>`
- No code changes needed for variable substitution logic

## Tasks

- [x] 1.0 Create frontmatter parsing utilities
  - [x] 1.1 Create new file `internal/instruction/frontmatter.go` with parsing functions
  - [x] 1.2 Implement `ParseFrontmatter(content string) (Meta, string, error)` function that splits frontmatter from body
  - [x] 1.3 Add validation for frontmatter delimiters (must start and end with `---`)
  - [x] 1.4 Implement frontmatter extraction using string splitting on `---` markers
  - [x] 1.5 Parse YAML frontmatter into Meta struct using `gopkg.in/yaml.v3`
  - [x] 1.6 Extract markdown body (everything after second `---`)
  - [x] 1.7 Add error handling for missing/malformed frontmatter with clear error messages
  - [x] 1.8 Create `internal/instruction/frontmatter_test.go` with comprehensive test cases
  - [x] 1.9 Test valid frontmatter parsing with all metadata fields
  - [x] 1.10 Test error cases: missing delimiters, invalid YAML, missing required fields
  - [x] 1.11 Test frontmatter with markdown body containing variable patterns

- [x] 2.0 Refactor instruction loading to use instruction.md
  - [x] 2.1 [depends on: 1.0] Update `Load()` function in `internal/instruction/load.go` to use instruction.md
  - [x] 2.2 [depends on: 1.0] Change file path from `system.txt` and `meta.yaml` to `instruction.md`
  - [x] 2.3 [depends on: 1.0] Call `ParseFrontmatter()` to extract metadata and content
  - [x] 2.4 [depends on: 1.0] Keep existing validation for tags and lang fields
  - [x] 2.5 Update error message when instruction.md doesn't exist
  - [x] 2.6 Remove code that reads system.txt and meta.yaml
  - [x] 2.7 Test loading instructions with valid instruction.md files
  - [x] 2.8 Test error handling for missing instruction.md
  - [x] 2.9 Test error handling for malformed frontmatter

- [x] 3.0 Refactor instruction creation to use instruction.md with frontmatter
  - [x] 3.1 Update `Create()` function in `internal/instruction/create.go`
  - [x] 3.2 Remove code that creates system.txt and meta.yaml files
  - [x] 3.3 Implement template generation for instruction.md with frontmatter structure
  - [x] 3.4 Format frontmatter with provided metadata (version, description, tags, lang)
  - [x] 3.5 Initialize markdown body with empty content or example template
  - [x] 3.6 Write instruction.md to disk with proper formatting
  - [x] 3.7 Update editor to open instruction.md instead of system.txt
  - [x] 3.8 Test creating new instruction generates valid instruction.md
  - [x] 3.9 Verify created instruction.md has correct frontmatter format
  - [x] 3.10 Verify created instruction can be loaded successfully with Load()

- [x] 4.0 Update edit command to work with instruction.md
  - [x] 4.1 [depends on: 2.0] Update `editCmd` in `cmd/edit.go` to reference instruction.md
  - [x] 4.2 Remove the `--meta` flag and related logic
  - [x] 4.3 Change file path construction to use instruction.md
  - [x] 4.4 Update command description to reflect editing single file
  - [x] 4.5 Test edit command opens instruction.md correctly
  - [x] 4.6 Verify changes to instruction.md are persisted

- [x] 5.0 Update version management for frontmatter format
  - [x] 5.1 [depends on: 1.0, 2.0] Review `internal/instruction/version.go` implementation
  - [x] 5.2 Update version reading to parse instruction.md frontmatter
  - [x] 5.3 Update version writing to modify frontmatter and preserve markdown body
  - [x] 5.4 Implement logic to read current instruction.md content
  - [x] 5.5 Parse frontmatter, update version field in Meta struct
  - [x] 5.6 Regenerate frontmatter with updated version
  - [x] 5.7 Write updated instruction.md preserving markdown body
  - [x] 5.8 Test version bump updates frontmatter correctly
  - [x] 5.9 Test version set updates frontmatter correctly
  - [x] 5.10 Verify markdown body is preserved during version updates

- [x] 6.0 Verify variable substitution works with markdown content
  - [x] 6.1 [depends on: 2.0, 3.0] Review `internal/instruction/variables.go` to understand current implementation
  - [x] 6.2 Create test instruction.md with markdown formatting and variables
  - [x] 6.3 Test simple variable substitution `{{variable}}` in markdown content
  - [x] 6.4 Test OR variable substitution `{{var1|var2}}` in markdown headers and lists
  - [x] 6.5 Test `{{input}}` handling with markdown content
  - [x] 6.6 Test variables within markdown code blocks
  - [x] 6.7 Test variables within markdown bold/italic text
  - [x] 6.8 Add test cases to `internal/instruction/variables_test.go` for markdown scenarios
  - [x] 6.9 Verify all existing variable tests still pass
  - [x] 6.10 Document any edge cases discovered during testing
