package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new instruction",
	Long:  `Creates a new instruction with the specified name and opens it in your editor.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description, _ := cmd.Flags().GetString("description")

		if err := instruction.Create(name, description); err != nil {
			return err
		}

		fmt.Printf("✓ Created instruction: %s\n", name)
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("description", "d", "", "Description of the instruction")
	rootCmd.AddCommand(addCmd)
}
