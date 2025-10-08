package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/instruction"
)

var versionCmd = &cobra.Command{
	Use:   "version <instruction> [subcommand]",
	Short: "Manage instruction versions",
	Long:  `Show, bump, or set the version of an instruction.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version, err := instruction.GetVersion(name)
		if err != nil {
			return err
		}
		fmt.Printf("%s v%s\n", name, version)
		return nil
	},
}

var versionBumpCmd = &cobra.Command{
	Use:   "bump <instruction> [description]",
	Short: "Bump the patch version of an instruction",
	Long:  `Increments the patch version number (e.g., 0.1.0 → 0.1.1)`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description := ""
		if len(args) > 1 {
			description = args[1]
		}

		oldVersion, newVersion, err := instruction.BumpVersion(name, description)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Version bumped: %s → %s\n", oldVersion, newVersion)
		return nil
	},
}

var versionSetCmd = &cobra.Command{
	Use:   "set <instruction> <version> [description]",
	Short: "Set a specific version for an instruction",
	Long:  `Sets the version to a specific value (must be valid semver: X.Y.Z)`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := args[1]
		description := ""
		if len(args) > 2 {
			description = args[2]
		}

		oldVersion, err := instruction.SetVersion(name, version, description)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Version set: %s → %s\n", oldVersion, version)
		return nil
	},
}

func init() {
	versionCmd.AddCommand(versionBumpCmd)
	versionCmd.AddCommand(versionSetCmd)
	rootCmd.AddCommand(versionCmd)
}
