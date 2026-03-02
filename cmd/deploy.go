package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var deployCmd = &cobra.Command{
	Use:   "deploy <name>",
	Short: "Deploy an instruction artifact to the current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		as, _ := cmd.Flags().GetString("as")
		force, _ := cmd.Flags().GetBool("force")

		inst, err := instruction.Load(name)
		if err != nil {
			return err
		}

		source, err := filepath.Abs(filepath.Join(inst.Path, "instruction.md"))
		if err != nil {
			return fmt.Errorf("failed to resolve instruction path: %w", err)
		}

		switch inst.Meta.Type {
		case "skill":
			return deploySkill(name, source, force)
		case "context":
			target := "CLAUDE.md"
			if as != "" {
				target = as
			}
			return deployContext(source, target)
		default:
			target := name + ".md"
			if as != "" {
				target = as
			}
			return deploySymlink(source, target, force)
		}
	},
}

func deploySkill(name, source string, force bool) error {
	skillDir := filepath.Join(".claude", "skills", name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	target := filepath.Join(skillDir, "SKILL.md")
	return deploySymlink(source, target, force)
}

func deploySymlink(source, target string, force bool) error {
	if _, err := os.Lstat(target); err == nil {
		if !force {
			return fmt.Errorf("'%s' already exists (use --force to overwrite)", target)
		}
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("failed to remove existing file: %w", err)
		}
	}

	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Printf("Deployed: %s -> %s\n", target, source)
	return nil
}

func deployContext(source, target string) error {
	reference := "@" + source

	data, err := os.ReadFile(target)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", target, err)
	}

	if strings.Contains(string(data), reference) {
		fmt.Printf("Already deployed to %s\n", target)
		return nil
	}

	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", target, err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\n%s\n", reference); err != nil {
		return fmt.Errorf("failed to write to %s: %w", target, err)
	}

	fmt.Printf("Deployed: appended %s to %s\n", reference, target)
	return nil
}

func init() {
	deployCmd.Flags().String("as", "", "Override target filename")
	deployCmd.Flags().Bool("force", false, "Overwrite existing symlink")
	rootCmd.AddCommand(deployCmd)
}
