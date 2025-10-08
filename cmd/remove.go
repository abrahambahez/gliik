package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an instruction",
	Long:    `Deletes an instruction and all its files. Requires confirmation unless --force is used.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		inst, err := instruction.Load(name)
		if err != nil {
			return err
		}

		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("Delete instruction '%s'? [y/N]: ", name)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Operation cancelled")
				return nil
			}
		}

		if err := os.RemoveAll(inst.Path); err != nil {
			return fmt.Errorf("failed to remove instruction: %w", err)
		}

		fmt.Printf("✓ Removed instruction: %s\n", name)
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}
