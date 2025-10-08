package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Gliik configuration",
	Long:  `Creates the ~/.gliik directory structure and default configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Initialize(); err != nil {
			return err
		}
		fmt.Printf("✓ Initialized Gliik at %s\n", config.GetGliikHome())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
