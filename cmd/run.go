package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gliik/internal/ai"
	"github.com/yourusername/gliik/internal/instruction"
)

var outputFile string

var runCmd = &cobra.Command{
	Use:                "run <instruction>",
	Short:              "Execute an instruction",
	Long:               `Executes an instruction with the AI, resolving variables from stdin or CLI flags.`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("instruction name is required")
		}

		instructionName := args[0]

		inst, err := instruction.Load(instructionName)
		if err != nil {
			return err
		}

		variables, err := instruction.ParseVariables(inst.SystemText)
		if err != nil {
			return err
		}

		tempCmd := &cobra.Command{}
		tempCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file")

		for _, v := range variables {
			for _, opt := range v.Options {
				if opt != "input" {
					tempCmd.Flags().String(opt, "", fmt.Sprintf("Value for %s", opt))
				}
			}
		}

		if err := tempCmd.ParseFlags(args[1:]); err != nil {
			return err
		}

		return executeInstruction(instructionName, tempCmd)
	},
}

func executeInstruction(name string, cmd *cobra.Command) error {
	inst, err := instruction.Load(name)
	if err != nil {
		return err
	}

	variables, err := instruction.ParseVariables(inst.SystemText)
	if err != nil {
		return err
	}

	stdin := ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		stdin = string(stdinBytes)
	}

	flags := make(map[string]string)
	for _, v := range variables {
		for _, opt := range v.Options {
			if opt == "input" {
				continue
			}
			if cmd.Flags().Changed(opt) {
				value, _ := cmd.Flags().GetString(opt)
				flags[opt] = value
			}
		}
	}

	resolver := instruction.Resolver{
		Variables: variables,
		Stdin:     stdin,
		Flags:     flags,
	}

	resolved, err := resolver.Resolve()
	if err != nil {
		return err
	}

	finalPrompt := inst.SystemText
	for varRaw, value := range resolved {
		finalPrompt = strings.ReplaceAll(finalPrompt, varRaw, value)
	}

	aiClient, err := ai.NewClient()
	if err != nil {
		return err
	}

	response, err := aiClient.Complete(finalPrompt)
	if err != nil {
		return err
	}

	fmt.Print(response)

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(response), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
