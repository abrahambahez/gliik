# Gliik MVP - Product Requirements Document

## Overview

**Gliik** is a CLI tool for managing and executing AI prompts (called "instructions") following UNIX philosophy: composability, minimalism, and clear separation of concerns.

**Version:** MVP 0.1.0  
**Language:** Go  
**Philosophy:** Content-first, human-readable, pipeline-friendly

---

## Core Concepts

### Instruction
A reusable AI prompt that performs a specific task. Each instruction:
- Lives in its own directory under `~/.gliik/instructions/`
- Contains pure prompt content (`system.txt`)
- Has metadata for versioning and description (`meta.yaml`)
- Supports variable substitution via `{{variable_name}}`
- Supports OR logic for flexible input via `{{var1|var2|var3}}`
- Can be piped, chained, and composed UNIX-style

### Directory Structure
```
~/.gliik/
├── instructions/              # User's instruction repository
│   ├── enhance_resume/
│   │   ├── system.txt        # Prompt content
│   │   └── meta.yaml         # Metadata
│   └── assess_paper/
│       ├── system.txt
│       └── meta.yaml
└── config.yaml               # Global configuration
```

---

## Metadata Format

### `meta.yaml` Schema
```yaml
version: "1.0.0"              # Semver format
description: "Brief description of what this instruction does"
```

**Notes:**
- `version`: Semantic versioning (MAJOR.MINOR.PATCH)
- `description`: Single-line human-readable description

### `config.yaml` Schema
```yaml
default_model: "claude-sonnet-4-20250514"
editor: "vim"                 # Override $EDITOR if needed
# API keys use environment variables: ANTHROPIC_API_KEY
```

---

## MVP Commands

### 1. `gliik init`
Initialize Gliik in user's home directory.

**Behavior:**
- Create `~/.gliik/instructions/` directory
- Create default `~/.gliik/config.yaml`
- Output success message with next steps

**Example:**
```bash
$ gliik init
✓ Initialized Gliik at ~/.gliik
✓ Created instructions directory
✓ Created default configuration

Next steps:
  gliik add <name>     Create your first instruction
  gliik list           List all instructions
```

---

### 2. `gliik add <name>`
Create a new instruction.

**Arguments:**
- `<name>`: Instruction name (alphanumeric + underscores)

**Flags:**
- `--description, -d`: Description of the instruction

**Behavior:**
- Create `~/.gliik/instructions/<name>/` directory
- Create `system.txt` with template content
- Create `meta.yaml` with version "0.1.0" and provided description
- Open `system.txt` in `$EDITOR` (or config.yaml editor)

**Example:**
```bash
$ gliik add enhance_resume -d "Enhance resumes for job applications"
✓ Created instruction: enhance_resume
✓ Opening system.txt in editor...
```

**Default `system.txt` template:**
```
You are an AI assistant helping with {{task_name}}.

Input: {{input|text}}

Instructions:
[Edit this file to define your instruction]

Output:
[Define expected output format]
```

---

### 3. `gliik list`
List all available instructions.

**Flags:**
- None (use pipes for filtering: `gliik list | grep resume`)

**Output format:**
```
NAME              VERSION    DESCRIPTION
enhance_resume    1.0.0      Enhance resumes for job applications
assess_paper      0.2.1      Evaluate research papers for quality
summarize_text    1.1.0      Create concise summaries of long texts
```

**Behavior:**
- Read all subdirectories in `~/.gliik/instructions/`
- Parse each `meta.yaml`
- Display in table format
- Sort alphabetically by name

---

### 4. `gliik edit <name>`
Edit an instruction's content.

**Arguments:**
- `<name>`: Instruction name

**Behavior:**
- Open `~/.gliik/instructions/<name>/system.txt` in `$EDITOR`
- Validate instruction exists, error if not

**Example:**
```bash
$ gliik edit enhance_resume
# Opens system.txt in configured editor
```

---

### 5. `gliik <instruction_name> [input_options]`
Execute an instruction with AI model.

**Input Methods (priority order):**
1. **Stdin pipe**: `cat input.txt | gliik instruction_name` (requires `{{input}}` in system.txt)
2. **CLI flags**: `--variable_name value` or `--variable_name file.txt`

**Special Variable:**
- `{{input}}`: Reserved variable that ONLY accepts stdin
- Cannot be provided via CLI flags
- If stdin provided but no `{{input}}` in system.txt → error
- If `{{input}}` in system.txt but no stdin → error

**Common Flags:**
- `--output, -o <file>`: Save output to file (also prints to stdout)
- `--model, -m <model>`: Override default model
- `--<variable_name> <value>`: Variable substitution

**Variable Resolution:**
1. Parse `system.txt` for `{{variable_name}}` patterns
2. If stdin provided:
   - Check if `{{input}}` exists in system.txt
   - If yes: replace `{{input}}` with stdin content
   - If no: error (instruction doesn't accept stdin)
3. For each remaining variable:
   - Map CLI flags to variables (e.g., `--job_description job.txt`)
   - If file path provided, read file contents
   - If value provided directly, use as-is
4. Error if any variables remain undefined

**Example Usage:**
```bash
# Method 1: CLI flags with files
$ gliik enhance_resume --job_description job.txt --base_resume resume.md

# Method 2: CLI flags with inline values
$ gliik enhance_resume --job_description "Senior Developer role" --base_resume resume.md

# Method 3: Stdin pipe (requires {{input}} or {{input|...}} in system.txt)
$ cat paper.txt | gliik assess_paper

# Method 4: Combination of stdin + CLI flags
$ cat resume.md | gliik enhance_resume --job_description job.txt
# system.txt must have {{input}} for stdin and {{job_description}} for the flag

# Method 5: OR variable flexibility (system.txt has {{input|text}})
$ cat paper.txt | gliik proofread          # Uses stdin
$ gliik proofread --text paper.txt         # Uses --text flag
$ gliik proofread --text "Quick text"      # Uses inline text

# Method 6: Save output
$ gliik enhance_resume --job job.txt --resume resume.md -o enhanced_resume.md

# Method 7: Pipeline composition
$ gliik enhance_resume --job job.txt --resume resume.md | gliik proofread
# proofread must have {{input}} or {{input|...}} in system.txt

# Method 8: Multiple OR groups (system.txt: "Compare {{input|text1}} with {{reference|text2}}")
$ cat file1.txt | gliik compare --reference file2.txt
$ gliik compare --text1 file1.txt --text2 file2.txt
```

**Behavior:**
- Validate instruction exists
- Read `system.txt` and `meta.yaml`
- Perform variable substitution
- Send to AI model (Anthropic Claude API via `ANTHROPIC_API_KEY` env var)
- Stream output to stdout
- If `--output` specified, also write to file
- Error handling:
  - Missing required variables → show error with expected variables
  - Missing API key → show error with setup instructions
  - API errors → display error message

**Error Example:**
```bash
$ gliik enhance_resume
Error: Missing required variables for 'enhance_resume'

Expected variables:
  --job_description    Job posting or description
  --base_resume        Current resume content

Usage:
  gliik enhance_resume --job_description <file|text> --base_resume <file|text>

$ cat resume.md | gliik enhance_resume
Error: Instruction 'enhance_resume' does not accept stdin input

This instruction expects variables via CLI flags, not stdin.
To use stdin, add {{input}} or {{input|var}} to the instruction's system.txt

Expected variables:
  --job_description
  --base_resume

$ gliik proofread
Error: Missing required variable for 'proofread'

Expected one of:
  • stdin (use: cat file | gliik proofread)
  • --text (use: gliik proofread --text <file|value>)

$ gliik compare --text1 file1.txt
Error: Missing required variable for 'compare'

Variable '{{reference|text2}}' needs one of:
  • --reference
  • --text2
```

---

### 6. `gliik remove <name>`
Delete an instruction.

**Aliases:** `rm`

**Arguments:**
- `<name>`: Instruction name

**Flags:**
- `--force, -f`: Skip confirmation prompt

**Behavior:**
- Prompt for confirmation (unless `--force`)
- Delete `~/.gliik/instructions/<name>/` directory and contents
- Display success message

**Example:**
```bash
$ gliik remove old_instruction
Delete instruction 'old_instruction'? [y/N]: y
✓ Removed instruction: old_instruction

$ gliik rm old_instruction -f
✓ Removed instruction: old_instruction
```

---

### 7. `gliik version <name> [subcommand]`
Manage instruction versioning.

**Subcommands:**

#### `gliik version <name>`
Show current version of instruction.

**Output:**
```bash
$ gliik version enhance_resume
enhance_resume v1.0.0
```

#### `gliik version <name> bump [description]`
Increment patch version (1.0.0 → 1.0.1).

**Example:**
```bash
$ gliik version enhance_resume bump "Fixed formatting issue"
✓ Version bumped: 1.0.0 → 1.0.1
```

#### `gliik version <name> set <version> [description]`
Set specific version.

**Arguments:**
- `<version>`: Semver string (e.g., "2.0.0")
- `[description]`: Optional version description

**Example:**
```bash
$ gliik version enhance_resume set 2.0.0 "Major rewrite with new approach"
✓ Version set: 1.0.1 → 2.0.0
```

**Behavior:**
- Validate semver format
- Update `meta.yaml`
- If description provided, update description field

---

## Technical Implementation Notes

### Go Packages
- CLI framework: `github.com/spf13/cobra`
- YAML parsing: `gopkg.in/yaml.v3`
- AI API: `github.com/anthropics/anthropic-sdk-go` (or HTTP client)
- Template parsing: `text/template` or regex for `{{variable}}`

### Variable Substitution Logic
```go
// Pseudocode
1. Parse system.txt for {{...}} patterns using regex
2. For each pattern:
   - If contains "|": create VariableGroup with options ["var1", "var2", ...]
   - If simple: create VariableGroup with single option ["var"]
3. Check if stdin is provided (read from os.Stdin)
4. For each VariableGroup:
   - resolved = false
   - For each option in group (left to right):
     - If option == "input":
       - If stdin available: use stdin, resolved = true, break
       - Else: continue to next option
     - Else:
       - If CLI flag --option exists:
         - If value is file path: read file contents
         - Else: use value as-is
         - resolved = true, break
   - If not resolved: add VariableGroup to unresolved list
5. If stdin provided but no "input" in any VariableGroup:
   - ERROR "Instruction does not accept stdin input"
6. If unresolved list not empty:
   - ERROR with details of missing variables (show all options for OR groups)
7. Replace all {{...}} patterns in system.txt with resolved values
```

### AI Model Integration
- MVP: Anthropic Claude API only
- Use `ANTHROPIC_API_KEY` environment variable
- Default model: `claude-sonnet-4-20250514`
- Stream responses to stdout in real-time
- Handle API errors gracefully

### File System Operations
- Atomic writes for `meta.yaml` updates
- Validate instruction names (alphanumeric + underscore only)
- Handle missing directories gracefully
- Use `$EDITOR` environment variable (fallback: vim)

---

## Error Handling

### Common Errors
1. **Instruction not found**
   ```
   Error: Instruction 'xyz' not found
   
   Run 'gliik list' to see available instructions
   ```

2. **Missing API key**
   ```
   Error: ANTHROPIC_API_KEY environment variable not set
   
   Set your API key:
     export ANTHROPIC_API_KEY="your-key-here"
   ```

3. **Invalid version format**
   ```
   Error: Invalid version format '1.0.a'
   
   Use semantic versioning (e.g., 1.0.0)
   ```

4. **Missing required variables**
   ```
   Error: Missing required variables for 'enhance_resume'
   
   Expected: --job_description, --base_resume
   ```

5. **Stdin provided to non-stdin instruction**
   ```
   Error: Instruction 'enhance_resume' does not accept stdin input
   
   This instruction expects variables via CLI flags, not stdin.
   To use stdin, add {{input}} or {{input|var}} to system.txt
   
   Expected variables: --job_description, --base_resume
   ```

6. **Stdin or variable expected but not provided (OR group)**
   ```
   Error: Missing required variable for 'proofread'
   
   Variable '{{input|text}}' needs one of:
     • stdin (use: cat file | gliik proofread)
     • --text (use: gliik proofread --text <file|value>)
   ```

7. **Simple variable missing**
   ```
   Error: Missing required variable for 'enhance_resume'
   
   Variable '{{job_description}}' is required
   
   Usage:
     gliik enhance_resume --job_description <file|value>
   ```

---

## Out of Scope (Post-MVP)

- Multiple instruction repositories
- Instruction dependencies and chaining
- Built-in testing framework
- Semantic search
- SQLite caching for fast lookups
- Remote instruction sharing/marketplace
- Version history tracking (beyond current version)
- Instruction templates library
- Multi-model support (OpenAI, local models)
- Interactive variable prompting

---

## Success Metrics (MVP)

- User can create, edit, list, and delete instructions via CLI
- User can execute instructions with variable substitution
- Instructions work with pipes and CLI flags
- Version management works correctly
- Clean, readable code following Go best practices
- Complete documentation and examples

---

## Example Workflow

```bash
# Setup
$ gliik init
$ export ANTHROPIC_API_KEY="sk-ant-..."

# Create instruction
$ gliik add enhance_resume -d "Enhance resumes for job applications"
# Edit system.txt to add prompt logic with {{job_description}} and {{base_resume}}

# Create instruction with flexible input
$ gliik add proofread -d "Proofread and correct text"
# Edit system.txt with: "Proofread this text: {{input|text}}"

# Use instruction with CLI flags
$ gliik enhance_resume --job_description job.txt --base_resume resume.md -o enhanced.md

# Use instruction with stdin
$ cat paper.txt | gliik proofread

# Use instruction with CLI flag (same instruction)
$ gliik proofread --text paper.txt

# Update and version
$ gliik edit enhance_resume
# Make improvements...
$ gliik version enhance_resume bump "Improved formatting section"

# List all instructions
$ gliik list | grep resume

# Pipeline usage with OR variables
$ gliik enhance_resume --job job.txt --resume resume.md | gliik proofread
# Both methods work:
# - enhance_resume outputs to stdout
# - proofread accepts via stdin ({{input}} in {{input|text}})
```

---

## Next Steps

1. Set up Go project structure
2. Implement core file system operations (init, add, list, edit, remove)
3. Implement variable substitution engine
4. Integrate Anthropic Claude API
5. Implement version management
6. Add comprehensive error handling
7. Write tests for core functionality
8. Create documentation and examples
