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

		fmt.Printf("Initialized Gliik at %s\n", config.GetGliikHome())
		if dir != "" {
			fmt.Printf("Instructions directory: %s\n", config.GetInstructionsDir())
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("dir", "d", "", "Custom instructions directory path")
	rootCmd.AddCommand(initCmd)
}
