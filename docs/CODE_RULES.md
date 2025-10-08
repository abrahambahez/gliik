# Gliik Code Rules

## Core Philosophy

### Lean Code & MVP Mindset
- Less code is always better
- For each feature, evaluate the simplest, leanest implementation
- Avoid over-engineering - build what's needed, not what might be needed
- Prefer simple solutions over clever ones
- Question every line: does this truly need to exist?

### Minimalism in Practice
- Use standard library before adding dependencies
- One clear way to do things, not multiple options
- Remove code as enthusiastically as you add it
- If a feature doesn't serve the MVP, defer it

## Code Style

### Comments Policy
- **No random inline comments**
- Code must be self-documenting through clear naming
- Use longer, descriptive function and variable names instead of explanatory comments
- Comments are ONLY for Godoc documentation

**Good:**
```go
// ParseTemplateVariables extracts variable patterns from instruction content
// and returns a list of VariableGroups for resolution.
func ParseTemplateVariables(content string) ([]VariableGroup, error) {
    variablePattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)
    // ...
}
```

**Bad:**
```go
// Parse variables
func Parse(c string) ([]VG, error) {
    re := regexp.MustCompile(`\{\{([^}]+)\}\}`) // regex for vars
    // loop through matches
    // ...
}
```

### Naming Conventions
- Long, descriptive names over short cryptic ones
- Function names should describe what they do: `ResolveVariableFromStdinOrFlags` not `Resolve`
- Variable names should describe what they hold: `unresolvedVariableGroups` not `unresolved`
- If you need a comment to explain it, the name isn't good enough

### Output & Messages
- No emojis in code, CLI output, or error messages
- Clear, direct error messages
- Follow standard Go formatting (gofmt, goimports)

### Debugging
- Keep debugging process clean
- Remove all debug output after debugging sessions
- No leftover print statements, temporary logging, or debug UI feedback
- Use proper logging only where necessary for production diagnostics

## Go-Specific Guidelines

### Language Idioms
- Follow Go proverbs: simple, clear, idiomatic
- Errors are values - handle them explicitly
- Accept interfaces, return structs
- Use `internal/` for private packages
- Prefer composition over inheritance

### Dependencies
- Prefer standard library over external packages
- Question every dependency: is this truly necessary?
- For MVP, only essential dependencies:
  - `github.com/spf13/cobra` (CLI framework)
  - `gopkg.in/yaml.v3` (YAML parsing)
  - `github.com/anthropics/anthropic-sdk-go` (AI API)

### Error Handling
- Explicit error handling, no silent failures
- Clear error messages that guide users to solutions
- Error context should flow up the stack

### Testing
- Test what matters for MVP functionality
- Focus on core logic (template engine, variable resolution)
- Don't test trivial code
- Tests should be simple and clear

## Project-Specific

### File Operations
- Atomic writes for critical files (meta.yaml, config.yaml)
- Clean error messages for file system issues
- Validate paths before operations

### Template Engine
- Keep parsing logic simple and testable
- Clear separation: parse → resolve → substitute
- Explicit handling of stdin vs. flags vs. OR groups

### CLI Design
- Follow UNIX philosophy: do one thing well
- Composable through pipes
- Clear, consistent command structure
- Minimal flags, maximum clarity
