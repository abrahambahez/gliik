# Configuration Migration Specification

## Overview
Migrate Gliik from `~/.gliik/` to follow XDG Base Directory specification using `~/.config/gliik/` with configurable instruction directory paths.

## Goals
1. Follow standard configuration conventions (`~/.config/gliik/`)
2. Allow users to specify custom instruction directories
3. Maintain simplicity and minimize code changes
4. Clean break from old structure (breaking change acceptable)

---

## Implementation Changes

### 1. Update Path Functions
**File**: `internal/config/paths.go`

```go
package config

import (
	"os"
	"path/filepath"
)

func GetGliikHome() string {
	// Check for XDG_CONFIG_HOME first (standard)
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "gliik")
	}
	
	// Fall back to ~/.config/gliik
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "gliik")
}

func GetInstructionsDir() string {
	// Try to load from config first
	cfg, err := Load()
	if err == nil && cfg.InstructionsDir != "" {
		return expandPath(cfg.InstructionsDir)
	}
	
	// Default to config directory
	return filepath.Join(GetGliikHome(), "instructions")
}

func GetConfigFile() string {
	return filepath.Join(GetGliikHome(), "config.yaml")
}

func expandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
```

### 2. Update Config Structure
**File**: `internal/config/init.go`

```go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultModel    string `yaml:"default_model"`
	Editor          string `yaml:"editor"`
	InstructionsDir string `yaml:"instructions_dir,omitempty"`
}

func Initialize(instructionsDir string) error {
	gliikHome := GetGliikHome()
	configFile := GetConfigFile()

	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("Gliik is already initialized at %s", gliikHome)
	}

	if err := os.MkdirAll(gliikHome, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	defaultConfig := Config{
		DefaultModel:    "claude-sonnet-4-20250514",
		Editor:          "vim",
		InstructionsDir: instructionsDir,
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	actualDir := GetInstructionsDir()
	if err := os.MkdirAll(actualDir, 0755); err != nil {
		return fmt.Errorf("failed to create instructions directory: %w", err)
	}

	return nil
}
```

### 3. Add Config Loader
**File**: `internal/config/config.go` (new file)

```go
package config

import (
	"fmt"
	"os"
	
	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	configFile := GetConfigFile()
	
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	return &cfg, nil
}
```

### 4. Add --dir Flag to Init Command
**File**: `cmd/init.go`

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Gliik configuration",
	Long:  `Creates the ~/.config/gliik directory structure and default configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		
		if err := config.Initialize(dir); err != nil {
			return err
		}
		
		fmt.Printf("✓ Initialized Gliik at %s\n", config.GetGliikHome())
		if dir != "" {
			fmt.Printf("✓ Instructions directory: %s\n", dir)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("dir", "d", "", "Custom instructions directory path")
	rootCmd.AddCommand(initCmd)
}
```

### 5. Update GetInstructionsDir to Handle Path Expansion
**File**: `internal/config/paths.go`

Add helper at the end of the file:

```go
func expandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
```

---

## Documentation Updates

### README.md

**Add to Setup Section:**

```markdown
### Setup

1. Initialize Gliik:
```bash
# Default location (~/.config/gliik/instructions)
gliik init

# Or specify custom instructions directory
gliik init --dir ~/my-instructions
gliik init -d /absolute/path/to/instructions
```

2. Set your Anthropic API key:
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

Add this to your `~/.bashrc` or `~/.zshrc` to persist.
```

**Update Commands Section:**

```markdown
## Commands

### `gliik init [--dir <path>]`
Initialize Gliik configuration at `~/.config/gliik`

**Flags:**
- `--dir, -d` - Custom instructions directory (default: `~/.config/gliik/instructions`)

**Examples:**
```bash
# Default location
gliik init

# Custom directory (relative to home)
gliik init --dir ~/my-instructions

# Custom directory (absolute path)
gliik init -d /path/to/instructions
```
```

**Update File Structure Section:**

```markdown
## File Structure

```
~/.config/gliik/
├── config.yaml          # Configuration
└── instructions/        # Default instruction location
    └── <name>/
        ├── system.txt   # Prompt template
        └── meta.yaml    # Version and description
```

**Update Configuration Section:**

```markdown
## Configuration

Located at `~/.config/gliik/config.yaml`:

```yaml
default_model: claude-sonnet-4-20250514
editor: vim
instructions_dir: ""  # Set via `gliik init --dir` or edit manually
```

### Custom Instructions Directory

Set during initialization:
```bash
gliik init --dir ~/my-instructions
```

Or edit `~/.config/gliik/config.yaml` manually:
```yaml
# Relative to home directory
instructions_dir: "~/my-instructions"

# Or absolute path
instructions_dir: "/absolute/path/to/instructions"
```

**Update Environment Variables Section:**

```markdown
## Environment Variables

- `ANTHROPIC_API_KEY` - Your Anthropic API key (required)
- `EDITOR` - Text editor for editing instructions (default: vim)
- `XDG_CONFIG_HOME` - Custom config directory (default: ~/.config)
```

### docs/PRD/prd-mvp.md

No updates required. PRD remains unchanged.

---

## Testing Checklist

- [ ] `gliik init` creates `~/.config/gliik/`
- [ ] `gliik init` creates `~/.config/gliik/instructions/`
- [ ] Config file created at `~/.config/gliik/config.yaml` with correct structure
- [ ] `XDG_CONFIG_HOME` environment variable respected when set
- [ ] Custom `instructions_dir` works with relative paths (`~/custom`)
- [ ] Custom `instructions_dir` works with absolute paths (`/absolute/path`)
- [ ] Empty/omitted `instructions_dir` defaults to `~/.config/gliik/instructions`
- [ ] `gliik add` creates instructions in configured directory
- [ ] `gliik list` reads from configured directory
- [ ] `gliik edit` opens files from configured directory
- [ ] `gliik run` executes instructions from configured directory
- [ ] `gliik remove` deletes from configured directory
- [ ] All commands work when `instructions_dir` is set to custom path

---

## Files Modified

### Code Changes
- `internal/config/paths.go` - Updated all path functions for new location
- `internal/config/init.go` - Added `InstructionsDir` field and `InitializeWithDir()` function
- `internal/config/config.go` - New file with `Load()` function
- `cmd/init.go` - Added `--dir` flag support

### Documentation Changes
- `README.md` - Updated paths, added `--dir` flag documentation, custom directory instructions

### Files Unchanged
- All other `cmd/*.go` files - Use config path functions (no direct changes needed)
- All `internal/instruction/*.go` files - Use config functions
- `internal/ai/client.go` - No changes needed
- `main.go` - No changes needed
- `docs/PRD/prd-mvp.md` - No updates required

---

## Migration Notes

### Breaking Change
This is a breaking change. Users must:
1. Re-run `gliik init` to create new structure
2. Manually move existing instructions from `~/.gliik/instructions/` to `~/.config/gliik/instructions/`

### Migration Commands
```bash
# Backup old data
cp -r ~/.gliik/instructions ~/gliik-backup

# Initialize new structure (default location)
gliik init

# Or initialize with custom location
gliik init --dir ~/my-instructions

# Copy instructions to new location
cp -r ~/gliik-backup/* ~/.config/gliik/instructions/
# Or if using custom directory:
cp -r ~/gliik-backup/* ~/my-instructions/

# Remove old directory
rm -rf ~/.gliik
```

---

## Implementation Order

1. Create `internal/config/config.go` with `Load()` function
2. Update `Config` struct in `internal/config/init.go` to add `InstructionsDir` field
3. Add `InitializeWithDir(customDir string)` function in `internal/config/init.go`
4. Update `cmd/init.go` to add `--dir` flag and call `InitializeWithDir()`
5. Update `GetGliikHome()` in `internal/config/paths.go` to use `~/.config/gliik`
6. Add `expandPath()` helper function in `internal/config/paths.go`
7. Update `GetInstructionsDir()` to load from config if present
8. Test `gliik init` with default config
9. Test `gliik init --dir ~/custom` with relative path
10. Test `gliik init --dir /absolute/path` with absolute path
11. Test all commands (`add`, `list`, `edit`, `run`, `remove`) with default location
12. Test all commands with custom `instructions_dir` from `--dir` flag
13. Test with `XDG_CONFIG_HOME` environment variable set
14. Update `README.md` with new paths and `--dir` flag documentation

---

## Key Design Decisions

### Why XDG Base Directory?
- Industry standard for Linux/Unix configuration
- Keeps `$HOME` clean
- Separates config from data
- Respects `XDG_CONFIG_HOME` for flexibility

### Why Allow Custom Instructions Directory?
- Users may want instructions in cloud-synced folders (Dropbox, iCloud)
- Teams may want shared instruction repositories
- Power users may organize differently
- Minimal code complexity for significant flexibility
- `--dir` flag provides convenience during initialization

### Why Empty String Means Default?
- YAML clarity: omitted field vs explicit override
- No need for special sentinel values
- Natural behavior: no config = use default

### Why Add --dir Flag to Init?
- User convenience: set custom directory at initialization time
- Reduces manual config file editing
- Clearer intent: custom path is a first-class option
- Still allows manual editing of config.yaml if needed later
