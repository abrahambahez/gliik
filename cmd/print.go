package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var printCmd = &cobra.Command{
	Use:   "print <name>",
	Short: "Print instruction prompt body to stdout",
	Long:  "Outputs the instruction's prompt body (markdown content) to stdout without variable substitution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		inst, err := instruction.Load(name)
		if err != nil {
			return err
		}

		fmt.Print(inst.SystemText)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
