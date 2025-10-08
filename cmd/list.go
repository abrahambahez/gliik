package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")

		for _, inst := range instructions {
			fmt.Fprintf(w, "%s\t%s\t%s\n", inst.Name, inst.Meta.Version, inst.Meta.Description)
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
