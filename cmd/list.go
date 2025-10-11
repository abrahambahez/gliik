package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instructions",
	Long:  `Displays a table of all available instructions with their versions and descriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instructions, err := instruction.ListAll()
		if err != nil {
			return err
		}

		if len(instructions) == 0 {
			fmt.Println("No instructions found. Use 'gliik add <name>' to create one.")
			return nil
		}

		for _, inst := range instructions {
			tags := strings.Join(inst.Meta.Tags, ", ")
			fmt.Printf("%s v%s [%s] - %s\n\n",
				inst.Name,
				inst.Meta.Version,
				tags,
				inst.Meta.Description)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
