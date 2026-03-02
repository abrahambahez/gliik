package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/config"
)

var syncCmd = &cobra.Command{
	Use:   "sync [pull|push]",
	Short: "Sync instructions via git",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subcommand := args[0]
		if subcommand != "pull" && subcommand != "push" {
			return fmt.Errorf("invalid subcommand '%s': must be pull or push", subcommand)
		}

		gliikHome := config.GetGliikHome()
		gitDir := filepath.Join(gliikHome, ".git")

		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			return fmt.Errorf("'%s' is not a git repository\n\nTo set up sync, run:\n  cd %s && git init && git remote add origin <url>", gliikHome, gliikHome)
		}

		out, err := exec.Command("git", "-C", gliikHome, subcommand).CombinedOutput()
		fmt.Print(string(out))
		if err != nil {
			return fmt.Errorf("git %s failed: %w", subcommand, err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
