# Gliik MVP - Lean Implementation Plan

## Philosophy
Build vertically, not horizontally. Implement one complete feature slice at a time, testing as you go. Ship the smallest working version, then iterate.

---

## Phase 0: Project Setup (30 min)

### Goal
Get a working Go CLI that does nothing but prints "Hello Gliik"

### Tasks
1. Initialize Go module: `go mod init github.com/yourusername/gliik`
2. Install Cobra: `go get -u github.com/spf13/cobra@latest`
3. Create basic project structure:
```
gliik/
├── main.go
├── cmd/
│   └── root.go
├── internal/
│   ├── config/
│   ├── instruction/
│   └── ai/
└── go.mod
```
4. Create root command with Cobra
5. Build and run: `go build && ./gliik`

**Test:** Run `./gliik` and see help output

---

## Phase 1: Init Command (1 hour)

### Goal
`gliik init` creates directory structure

### Tasks
1. Implement `cmd/init.go` with Cobra
2. Create `internal/config/paths.go`:
   - `GetGliikHome()` → `~/.gliik`
   - `GetInstructionsDir()` → `~/.gliik/instructions`
3. Create `internal/config/init.go`:
   - `Initialize()` function:
     - Check if already initialized
     - Create directories
     - Create default `config.yaml`
4. Wire up init command

**Test:** 
```bash
$ ./gliik init
✓ Initialized Gliik at ~/.gliik
$ ls ~/.gliik
config.yaml  instructions/
```

**Deliverable:** Working `init` command

---

## Phase 2: Add Command (2 hours)

### Goal
`gliik add <name>` creates instruction scaffold

### Tasks
1. Implement `cmd/add.go`
2. Create `internal/instruction/instruction.go`:
   - `Instruction` struct:
     ```go
     type Instruction struct {
         Name        string
         Path        string
         SystemText  string
         Meta        Meta
     }
     
     type Meta struct {
         Version     string
         Description string
     }
     ```
3. Create `internal/instruction/create.go`:
   - `Create(name, description string)` function:
     - Validate name (alphanumeric + underscore)
     - Create directory
     - Write default `system.txt`
     - Write `meta.yaml` with version "0.1.0"
4. Use `gopkg.in/yaml.v3` for YAML handling
5. Open `system.txt` in `$EDITOR` after creation

**Test:**
```bash
$ ./gliik add test_instruction -d "Test instruction"
✓ Created instruction: test_instruction
$ ls ~/.gliik/instructions/test_instruction
meta.yaml  system.txt
$ cat ~/.gliik/instructions/test_instruction/meta.yaml
version: "0.1.0"
description: "Test instruction"
```

**Deliverable:** Working `add` command

---

## Phase 3: List & Edit Commands (1 hour)

### Goal
List instructions and edit them

### Tasks
1. Implement `cmd/list.go`:
   - `internal/instruction/list.go`:
     - `ListAll()` → scan instructions directory
     - Read each `meta.yaml`
     - Return slice of Instructions
   - Format as table (use `text/tabwriter`)
2. Implement `cmd/edit.go`:
   - `internal/instruction/load.go`:
     - `Load(name string)` → read instruction
   - Open system.txt in `$EDITOR`

**Test:**
```bash
$ ./gliik list
NAME              VERSION    DESCRIPTION
test_instruction  0.1.0      Test instruction

$ ./gliik edit test_instruction
# Opens in editor
```

**Deliverable:** Working `list` and `edit` commands

---

## Phase 4: Variable Parser (2 hours)

### Goal
Parse `{{variable}}` and `{{var1|var2}}` from system.txt

### Tasks
1. Create `internal/instruction/variables.go`:
   ```go
   type Variable struct {
       Raw     string   // "{{input|text}}"
       Options []string // ["input", "text"]
   }
   
   func ParseVariables(systemText string) []Variable
   ```
2. Use regex to find `{{...}}` patterns:
   - Pattern: `\{\{([^}]+)\}\}`
   - Split on `|` to get options
3. Write comprehensive tests:
   - Simple variable: `{{text}}`
   - OR variable: `{{input|text}}`
   - Multiple: `{{job}} and {{input|resume}}`
   - Edge cases: spaces, special chars

**Test:**
```go
// Create test file
func TestParseVariables(t *testing.T) {
    text := "Process {{input|text}} with {{config}}"
    vars := ParseVariables(text)
    // Assert: 2 variables, first has 2 options, second has 1
}
```

**Deliverable:** Variable parser with tests

---

## Phase 5: Variable Resolution (3 hours)

### Goal
Resolve variables from stdin and CLI flags

### Tasks
1. Create `internal/instruction/resolver.go`:
   ```go
   type Resolver struct {
       Variables []Variable
       Stdin     string
       Flags     map[string]string
   }
   
   func (r *Resolver) Resolve() (map[string]string, error)
   ```
2. Implement resolution logic:
   - Check if stdin available
   - For each Variable:
     - Iterate through Options (left to right)
     - If option is "input" and stdin available: resolve
     - If option in Flags: resolve (check if file, read if needed)
     - If resolved: break, next variable
   - Return unresolved errors with helpful messages
3. Implement file detection and reading:
   - `isFile(path string)` → check if file exists
   - `readFile(path string)` → read contents
4. Handle stdin check separately:
   - If stdin provided but no "input" in variables: error

**Test:**
```go
func TestResolver(t *testing.T) {
    // Test stdin resolution
    // Test CLI flag resolution
    // Test OR priority
    // Test missing variables
    // Test stdin rejection
}
```

**Deliverable:** Variable resolver with tests

---

## Phase 6: AI Integration (2 hours)

### Goal
Send resolved prompt to Claude API

### Tasks
1. Create `internal/ai/client.go`:
   ```go
   type Client struct {
       APIKey string
       Model  string
   }
   
   func (c *Client) Complete(prompt string) (string, error)
   ```
2. Use Anthropic API:
   - Endpoint: `https://api.anthropic.com/v1/messages`
   - Headers: `x-api-key`, `anthropic-version: 2023-06-01`
   - Body: `{model, max_tokens, messages: [{role: "user", content: prompt}]}`
3. Stream response to stdout in real-time:
   - Use `bufio.Scanner` to read line by line
   - Print each line as received
4. Handle errors:
   - Missing API key
   - API errors
   - Network errors

**Test:**
```bash
$ export ANTHROPIC_API_KEY="sk-ant-..."
$ echo "What is 2+2?" | go run main.go test_simple
# Should get Claude response
```

**Deliverable:** Working AI client

---

## Phase 7: Run Command (3 hours)

### Goal
Execute instructions end-to-end

### Tasks
1. Implement `cmd/run.go` or make it implicit:
   - Cobra command: `gliik <instruction_name>`
   - Use `cobra.MinimumNArgs(1)` validation
2. Wire everything together:
   ```go
   func ExecuteInstruction(name string, args []string) error {
       // 1. Load instruction
       inst := instruction.Load(name)
       
       // 2. Parse variables
       vars := instruction.ParseVariables(inst.SystemText)
       
       // 3. Check stdin
       stdin := readStdin()
       
       // 4. Parse CLI flags dynamically
       flags := parseDynamicFlags(args, vars)
       
       // 5. Resolve variables
       resolver := instruction.Resolver{vars, stdin, flags}
       resolved, err := resolver.Resolve()
       
       // 6. Replace variables in system.txt
       finalPrompt := replaceVariables(inst.SystemText, resolved)
       
       // 7. Call AI
       client := ai.NewClient()
       response, err := client.Complete(finalPrompt)
       
       // 8. Output
       fmt.Print(response)
       if outputFlag {
           writeFile(outputPath, response)
       }
   }
   ```
3. Implement dynamic flag parsing:
   - Extract variable names from Variables
   - Create Cobra flags on-the-fly
   - Parse flag values
4. Implement stdin reading:
   - Check if stdin is pipe: `stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0`
   - Read all: `ioutil.ReadAll(os.Stdin)`
5. Add `--output` flag handling

**Test:**
```bash
# Create test instruction
$ cat > ~/.gliik/instructions/echo_test/system.txt << EOF
Echo this: {{input|text}}
EOF

# Test stdin
$ echo "Hello" | ./gliik echo_test

# Test flag
$ ./gliik echo_test --text "Hello"

# Test output
$ ./gliik echo_test --text "Hello" -o output.txt
$ cat output.txt
```

**Deliverable:** Working instruction execution

---

## Phase 8: Remove Command (30 min)

### Goal
Delete instructions

### Tasks
1. Implement `cmd/remove.go` (alias: `rm`)
2. Add confirmation prompt (unless `--force`)
3. Use `os.RemoveAll()` to delete directory

**Test:**
```bash
$ ./gliik remove test_instruction
Delete instruction 'test_instruction'? [y/N]: y
✓ Removed instruction: test_instruction
```

**Deliverable:** Working `remove` command

---

## Phase 9: Version Command (1.5 hours)

### Goal
Manage instruction versions

### Tasks
1. Implement `cmd/version.go` with subcommands:
   - `gliik version <name>` - show current
   - `gliik version <name> bump [desc]` - increment patch
   - `gliik version <name> set <version> [desc]` - set specific
2. Create `internal/instruction/version.go`:
   - `GetVersion(name)` - read meta.yaml
   - `BumpVersion(name, desc)` - parse semver, increment patch, write
   - `SetVersion(name, version, desc)` - validate semver, write
3. Use simple string manipulation for semver:
   ```go
   func BumpPatch(version string) string {
       parts := strings.Split(version, ".")
       patch, _ := strconv.Atoi(parts[2])
       return fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch+1)
   }
   ```
4. Validate semver format with regex: `^\d+\.\d+\.\d+$`

**Test:**
```bash
$ ./gliik version test_instruction
test_instruction v0.1.0

$ ./gliik version test_instruction bump "Fixed bug"
✓ Version bumped: 0.1.0 → 0.1.1

$ ./gliik version test_instruction set 1.0.0 "First stable"
✓ Version set: 0.1.1 → 1.0.0
```

**Deliverable:** Working `version` command

---

## Phase 10: Polish & Error Handling (2 hours)

### Goal
Better UX and error messages

### Tasks
1. Improve all error messages:
   - Add suggestions and examples
   - Use colors (optional): `github.com/fatih/color`
2. Add validation everywhere:
   - Instruction name format
   - File paths
   - API key presence
3. Add help text to all commands
4. Create README.md with:
   - Installation instructions
   - Quick start guide
   - Examples
5. Test edge cases:
   - Missing directories
   - Corrupted YAML
   - Empty instructions
   - Network failures

**Test:** Manual testing of all error scenarios

**Deliverable:** Production-ready CLI

---

## Phase 11: Build & Release (1 hour)

### Goal
Distribute the binary

### Tasks
1. Add build script:
   ```bash
   #!/bin/bash
   GOOS=darwin GOARCH=amd64 go build -o gliik-darwin-amd64
   GOOS=linux GOARCH=amd64 go build -o gliik-linux-amd64
   GOOS=windows GOARCH=amd64 go build -o gliik-windows-amd64.exe
   ```
2. Create installation instructions in README
3. Optional: Setup GitHub Actions for releases
4. Tag version: `git tag v0.1.0`

**Deliverable:** Distributable binaries

---

## Total Time Estimate: ~18 hours

Broken down:
- **Core functionality**: 12 hours (Phases 1-7)
- **Additional features**: 3 hours (Phases 8-9)
- **Polish**: 3 hours (Phases 10-11)

---

## Development Tips

### 1. Test as you build
```bash
# Keep a test instruction handy
$ cat > ~/.gliik/instructions/test/system.txt << EOF
You are a helpful assistant.
Input: {{input|text}}
Respond concisely.
EOF
```

### 2. Use main.go as integration test
```go
// main.go - simple and clean
package main

import "github.com/yourusername/gliik/cmd"

func main() {
    cmd.Execute()
}
```

### 3. Build often
```bash
# During development
$ go run main.go <command>

# For testing full binary
$ go build && ./gliik <command>
```

### 4. Defer optimization
- Don't cache metadata initially
- Don't add SQLite in MVP
- Don't optimize file reads
- Focus on correctness first

### 5. Use Go's standard library
- `os` for file operations
- `text/template` or `strings.Replace` for variable substitution
- `regexp` for variable parsing
- `net/http` for API calls
- No need for heavy frameworks

---

## Post-MVP Improvements (Future)

Not in scope for MVP, but good to keep in mind:
1. SQLite cache for fast `list` operations
2. Instruction templates/marketplace
3. Multi-model support (OpenAI, local models)
4. Version history (not just current version)
5. Instruction dependencies and chaining syntax
6. Config for multiple API keys
7. Interactive mode for missing variables
8. Better streaming (word-by-word vs line-by-line)
9. Progress indicators for long operations
10. Shell completions (bash, zsh, fish)

---

## Critical Path

The absolute minimum to have a working MVP:
1. Init (Phase 1)
2. Add (Phase 2)
3. Variable Parser (Phase 4)
4. Variable Resolver (Phase 5)
5. AI Integration (Phase 6)
6. Run (Phase 7)

**Minimum viable timeline: 12 hours**

Everything else (list, edit, remove, version) can be added after you have working instruction execution.

---

## Getting Started Checklist

- [ ] Phase 0: Project setup
- [ ] Phase 1: Init command
- [ ] Phase 2: Add command
- [ ] Phase 3: List & edit
- [ ] Phase 4: Variable parser
- [ ] Phase 5: Variable resolver
- [ ] Phase 6: AI integration
- [ ] Phase 7: Run command
- [ ] Phase 8: Remove command
- [ ] Phase 9: Version command
- [ ] Phase 10: Polish
- [ ] Phase 11: Build & release

Start with Phase 0 and work sequentially. Each phase should be fully working before moving to the next.
