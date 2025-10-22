# PRD: Markdown Instruction Format

## Introduction/Overview

Currently, Gliik instructions store their prompt content in `system.txt` and metadata in `meta.yaml` as separate files within each instruction directory. This PRD outlines the refactoring to consolidate both prompt content and metadata into a single Markdown file with YAML frontmatter (`instruction.md`).

**Problem:** The current two-file approach (system.txt + meta.yaml) creates unnecessary file management overhead and separates logically connected data (the prompt and its metadata).

**Goal:** Simplify the instruction file structure by using a single Markdown file with YAML frontmatter that contains all metadata and prompt content, while enabling richer prompt formatting capabilities.

## Goals

1. Replace the dual-file structure (system.txt + meta.yaml) with a single `instruction.md` file per instruction
2. Preserve all existing metadata fields (version, description, tags, lang)
3. Enable markdown formatting within prompt content for better structure and readability
4. Support multiple logical sections within prompts (e.g., Context, Task, Examples)
5. Maintain all existing variable substitution functionality (`{{variable}}`, `{{var1|var2}}`, `{{input}}`)
6. Ensure the refactor is a clean breaking change without backward compatibility

## User Stories

1. **As a Gliik user**, I want to view and edit all instruction information in a single file, so that I can understand and modify prompts more efficiently.

2. **As a prompt author**, I want to use markdown formatting (bold, italic, lists, code blocks) in my prompts, so that I can create well-structured, readable AI instructions.

3. **As a prompt author**, I want to organize my prompts into sections (e.g., ## Context, ## Task, ## Examples), so that complex prompts are easier to understand and maintain.

4. **As a developer**, I want the instruction loading logic to parse frontmatter metadata, so that all instruction data comes from a single canonical source.

5. **As a Gliik user**, I want `gliik edit <name>` to open the markdown file, so that I can edit everything in one editor session.

## Functional Requirements

### FR1: Instruction File Format
The instruction file must be named `instruction.md` and located at `~/.gliik/instructions/<name>/instruction.md`.

### FR2: YAML Frontmatter Structure
The frontmatter must be delimited by `---` markers and contain all current metadata fields:
```markdown
---
version: "1.0.0"
description: "Brief description of the instruction"
tags: ["tag1", "tag2"]
lang: "en"
---

Prompt content goes here...
```

### FR3: Metadata Fields
The frontmatter must support the following fields with the same validation rules as current implementation:
- `version` (string): Semantic version (e.g., "1.0.0")
- `description` (string): Brief description of the instruction
- `tags` (array of strings): At least one tag required, lowercase alphanumeric with hyphens only
- `lang` (string): ISO 639-1 two-letter lowercase language code (e.g., "en", "es")

### FR4: Markdown Body Support
The markdown body must support:
- Standard markdown formatting (bold, italic, lists, code blocks, etc.)
- Multiple sections using markdown headers (##, ###, etc.)
- All existing variable substitution patterns (`{{variable}}`, `{{input}}`, `{{var1|var2}}`)

### FR5: Frontmatter Parsing
The system must:
- Parse YAML frontmatter using `---` delimiters
- Extract metadata into the existing `Meta` struct
- Extract the markdown body as the prompt content (SystemText)
- Return clear errors if frontmatter is malformed or missing required fields

### FR6: Create Command Modification
`gliik add <name>` must:
- Create `instruction.md` (not system.txt + meta.yaml)
- Initialize the file with frontmatter containing provided metadata
- Open the markdown file in the editor for prompt content entry

### FR7: Edit Command Modification
`gliik edit <name>` must:
- Open `instruction.md` for editing (removing the need for --meta flag)
- No longer support editing separate system.txt or meta.yaml files

### FR8: Load Function Refactoring
The `instruction.Load()` function must:
- Read `instruction.md` instead of system.txt and meta.yaml
- Parse frontmatter to extract metadata
- Extract markdown body as prompt content
- Validate all metadata fields with existing validation rules
- Return appropriate errors for missing or malformed files

### FR9: Variable Substitution Preservation
All existing variable substitution logic must work identically on markdown content:
- `{{variable}}` patterns must be detected and resolved
- `{{var1|var2|var3}}` OR logic must function as before
- `{{input}}` stdin handling must remain unchanged

### FR10: Validation
The system must validate:
- Frontmatter is present and valid YAML
- All required metadata fields exist (tags, lang must not be empty)
- Language code is valid ISO 639-1 format
- Tags match the pattern `^[a-z0-9-]+$`
- Instruction name matches `^[a-zA-Z0-9_]+$`

## Non-Goals (Out of Scope)

1. **Backward compatibility**: No support for reading old system.txt + meta.yaml format
2. **Migration tool**: No automated command to convert existing instructions to new format
3. **Markdown comment stripping**: Comments in markdown (e.g., `<!-- -->`) will be sent to the AI as-is
4. **Frontmatter-only editing**: No special mode to edit just metadata without opening the full file
5. **Multiple instruction files**: Each instruction still has exactly one instruction.md file

## Design Considerations

### Example Instruction File Structure

```markdown
---
version: "1.0.0"
description: "Generate a professional cover letter from resume and job description"
tags: ["writing", "job-search"]
lang: "en"
---

## Context

You are an expert career coach helping job seekers create compelling cover letters.

## Task

Generate a professional cover letter based on the provided resume and job description.

The cover letter should:
- Be concise (3-4 paragraphs)
- Highlight relevant experience from the resume
- Address key requirements from the job description
- Maintain a professional yet personable tone

## Input

**Resume:**
{{resume}}

**Job Description:**
{{job_description|input}}
```

### File Structure After Refactor

```
~/.gliik/
├── instructions/
│   └── <name>/
│       └── instruction.md        (replaces system.txt + meta.yaml)
└── config.yaml
```

## Technical Considerations

### Dependencies
- Continue using `gopkg.in/yaml.v3` for YAML parsing
- Consider using `gopkg.in/yaml.v3` frontmatter extraction or implement simple string splitting on `---` delimiters

### Implementation Approach
1. Update `instruction.Instruction` struct if needed (likely no changes needed to struct itself)
2. Refactor `instruction.Load()` to read and parse instruction.md
3. Refactor `instruction.Create()` to generate instruction.md with frontmatter
4. Update `cmd/edit.go` to remove --meta flag and only edit instruction.md
5. Update all file path references from system.txt/meta.yaml to instruction.md
6. Update validation to ensure frontmatter presence and correctness

### Affected Files
Based on current codebase analysis:
- `/internal/instruction/instruction.go` - May need parsing helper functions
- `/internal/instruction/load.go` - Complete refactor of Load() function
- `/internal/instruction/create.go` - Refactor Create() to generate markdown with frontmatter
- `/cmd/edit.go` - Simplify to only edit instruction.md
- `/internal/instruction/list.go` - Update file references
- `/internal/instruction/version.go` - Update to read/write frontmatter
- `/internal/instruction/variables.go` - Likely no changes needed
- `/internal/instruction/resolver.go` - Likely no changes needed

### Error Handling
Must provide clear error messages for:
- Missing instruction.md file
- Malformed frontmatter (missing --- delimiters)
- Invalid YAML in frontmatter
- Missing required metadata fields
- Invalid metadata values (bad language code, invalid tags)

## Success Metrics

1. All existing instructions can be manually recreated in the new format without data loss
2. All variable substitution patterns continue to work identically
3. Markdown formatting in prompts is preserved when sent to AI
4. Edit workflow is simplified (single file vs. two files)
5. No regression in existing functionality (list, version, remove, execute commands)

## Open Questions

1. Should frontmatter validation be strict (fail on unknown fields) or permissive (ignore unknown fields)?
2. Should we provide example template when creating new instructions?
3. How should we handle instructions that have system.txt/meta.yaml? Display error message with migration guidance?
4. Should markdown content be processed/rendered before being sent to AI, or sent as raw markdown?
