# Gliik

A CLI tool for managing and executing AI prompts (called "instructions") following UNIX philosophy: composability, minimalism, and clear separation of concerns.

>[!note]
> This is an early version. Supports both Anthropic Claude (cloud) and Ollama (local) providers.
>

## Features

- **Multiple AI providers**: Anthropic Claude (cloud) or Ollama (local)
- Store reusable AI prompts as instructions
- **Markdown-based format** with YAML frontmatter for metadata
- Variable substitution with `{{variable}}` syntax
- OR logic for flexible input: `{{input|text}}`
- Pipe stdin or use CLI flags
- File path support for variables
- Version management for instructions
- Simple, composable commands
- Rich formatting support (headers, lists, code blocks, etc.)
- Streaming responses for real-time output

## Installation

### Prerequisites

- Go 1.24 or later
- **For Anthropic provider**: Anthropic API key (set as `ANTHROPIC_API_KEY` env variable)
- **For Ollama provider**: Ollama installed and running locally ([ollama.com](https://ollama.com))

### Build from Source

```bash
git clone https://github.com/yourusername/gliik.git
cd gliik
go build
```

This creates a `gliik` binary in the current directory.

To build with version from git tags:
```bash
go build -ldflags "-X 'github.com/yourusername/gliik/cmd.version=$(git describe --tags --always --dirty)'"
```

### Install to PATH

```bash
# Move binary to a directory in your PATH
sudo mv gliik /usr/local/bin/

# Or add to your home bin directory
mv gliik ~/bin/
```

### Setup

1. Initialize Gliik:
```bash
gliik init
```

2. **Choose your AI provider** by editing `~/.config/gliik/config.yaml`:

   **Option A: Anthropic Claude (default)**
   ```yaml
   provider: anthropic
   anthropic:
     model: claude-sonnet-4-20250514
   ```
   Then set your API key:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```

   **Option B: Ollama (local)**
   ```yaml
   provider: ollama
   ollama:
     endpoint: http://localhost:11434
     model: llama3.2
   ```
   Make sure Ollama is running:
   ```bash
   ollama serve
   ollama pull llama3.2
   ```

Add the `export` to your `~/.bashrc` or `~/.zshrc` to persist the API key.

## Quick Start

### Create an Instruction

```bash
gliik add summarize -d "Summarize text" -t summary -t text -l en
# Opens editor to edit instruction.md
```

In the editor, you'll see a template with frontmatter and markdown body:
```markdown
---
version: "0.1.0"
description: "Summarize text"
tags:
  - summary
  - text
lang: "en"
---
# summarize

Please summarize the following text in one sentence:

{{input|text}}
```

### Execute with Stdin

```bash
cat article.txt | gliik run summarize
```

### Execute with CLI Flag

```bash
gliik run summarize --text "Long text to summarize..."
```

### Execute with File Path

```bash
gliik run summarize --text article.txt
```

## Commands

### `gliik --version`
Show Gliik tool version

### `gliik init`
Initialize Gliik configuration at `~/.gliik`

### `gliik add <name> -d "description"`
Create a new instruction and open in editor

### `gliik list`
List all instructions with versions and descriptions

### `gliik edit <name>`
Edit an instruction's instruction.md file

### `gliik run <name> [flags]`
Execute an instruction with AI

### `gliik remove <name> [-f]`
Delete an instruction (with optional force flag)

### `gliik version <name>`
Show instruction version

### `gliik version bump <name> [description]`
Increment patch version (0.1.0 → 0.1.1)

### `gliik version set <name> <version> [description]`
Set specific version (must be X.Y.Z format)

## Variable Syntax

### Simple Variable
```
Hello {{name}}!
```
Usage: `gliik run greet --name "Alice"`

### OR Variable (Multiple Options)
```
Process this: {{input|text}}
```
Usage (either):
- `cat file.txt | gliik run process`
- `gliik run process --text "content"`

### Reserved: `{{input}}`
The `input` option is reserved for stdin only:
```
{{input}}       # Accepts only stdin
{{input|text}}  # Accepts stdin OR --text flag
```

## File Structure

```
~/.gliik/
├── config.yaml          # Configuration
└── instructions/
    └── <name>/
        └── instruction.md   # Single file with YAML frontmatter + markdown body
```

Each `instruction.md` follows this format:
```markdown
---
version: "1.0.0"
description: "Brief description"
tags:
  - tag1
  - tag2
lang: "en"
---
# Your prompt content here

With {{variable}} substitution support.
```

## Configuration

Located at `~/.config/gliik/config.yaml`:

```yaml
default_model: claude-sonnet-4-20250514
editor: vim
instructions_dir: ~/.gliik/instructions  # Optional: custom instructions directory
provider: anthropic  # or "ollama"

# Provider-specific configuration
anthropic:
  model: claude-sonnet-4-20250514

ollama:
  endpoint: http://localhost:11434
  model: llama3.2
```

**Configuration options:**
- `provider`: Choose between `"anthropic"` (cloud) or `"ollama"` (local)
- `anthropic.model`: Which Claude model to use
- `ollama.endpoint`: Ollama server URL (default: `http://localhost:11434`)
- `ollama.model`: Which Ollama model to use (run `ollama list` to see available models)

## Environment Variables

- `ANTHROPIC_API_KEY` - Your Anthropic API key (required only when using `provider: anthropic`)
- `EDITOR` - Text editor for editing instructions (default: vim)

## Examples

### Code Review Instruction

Create:
```bash
gliik add review -d "Review code for issues" -t code-review -t development -l en
```

instruction.md:
```markdown
---
version: "0.1.0"
description: "Review code for issues"
tags:
  - code-review
  - development
lang: "en"
---
# Code Review

You are a code reviewer. Review the following code and provide feedback:

{{input|code}}

## Focus Areas

- Bugs and errors
- Performance issues
- Best practices
```

Use:
```bash
cat main.go | gliik run review
# or
gliik run review --code main.go
```

### Resume Tailoring

Create:
```bash
gliik add tailor_resume -d "Tailor resume to job" -t career -t job-seeking -l en
```

instruction.md:
```markdown
---
version: "0.1.0"
description: "Tailor resume to job"
tags:
  - career
  - job-seeking
lang: "en"
---
# Resume Tailoring

Tailor this resume to match the job description:

## Resume
{{resume}}

## Job Description
{{job}}

Provide suggestions for improvements.
```

Use:
```bash
gliik run tailor_resume --resume resume.pdf --job job_desc.txt
```

## Tips

- Use descriptive instruction names with underscores: `review_code`, `summarize_text`
- Keep prompts focused and single-purpose
- Use OR variables for flexibility: `{{input|text}}`
- Version your instructions as you improve them
- File paths are automatically detected and read
- Leverage markdown formatting for better prompt organization (headers, lists, emphasis)
- Use tags to categorize instructions for easier discovery
- Frontmatter metadata is never sent to the LLM - only the markdown body

## License

MIT
