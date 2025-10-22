package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an instruction",
	Long:  `Opens the instruction's instruction.md file in your editor.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		inst, err := instruction.Load(name)
		if err != nil {
			return err
		}

		instructionFile := filepath.Join(inst.Path, "instruction.md")

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		editorCmd := exec.Command(editor, instructionFile)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
