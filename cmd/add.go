package cmd

import (
	"fmt"
	"strings"

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
		tagsStr, _ := cmd.Flags().GetString("tags")
		lang, _ := cmd.Flags().GetString("lang")

		if description == "" {
			return fmt.Errorf("missing required flag: --description")
		}

		if tagsStr == "" {
			return fmt.Errorf("missing required flag: --tags")
		}

		if lang == "" {
			return fmt.Errorf("missing required flag: --lang")
		}

		tags := strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}

		if err := instruction.Create(name, description, tags, lang); err != nil {
			return err
		}

		fmt.Printf("Created instruction: %s\n", name)
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("description", "d", "", "Description of the instruction (required)")
	addCmd.Flags().StringP("tags", "t", "", "Comma-separated tags (required)")
	addCmd.Flags().StringP("lang", "l", "", "Language ISO code (required)")
	rootCmd.AddCommand(addCmd)
}
