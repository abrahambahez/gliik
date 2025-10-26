package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "gliik",
	Short:   "A CLI tool for managing and executing AI prompts",
	Version: version,
	Long: `Gliik is a CLI tool for managing and executing AI prompts (called "instructions") following UNIX philosophy: composability, minimalism, and clear separation of concerns.

Instructions are stored in directories managed by Gliik (default: ~/.gliik/instructions/) and can contain variables that are resolved from stdin or CLI flags.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
