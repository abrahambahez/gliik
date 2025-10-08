package instruction

import (
	"fmt"
	"os"
)

type Resolver struct {
	Variables []Variable
	Stdin     string
	Flags     map[string]string
}

func (r *Resolver) Resolve() (map[string]string, error) {
	hasInputOption := false
	for _, v := range r.Variables {
		for _, opt := range v.Options {
			if opt == "input" {
				hasInputOption = true
				break
			}
		}
		if hasInputOption {
			break
		}
	}

	if r.Stdin != "" && !hasInputOption {
		return nil, fmt.Errorf("instruction does not accept stdin input\n\nThis instruction expects variables via CLI flags, not stdin.\nTo use stdin, add {{input}} or {{input|var}} to system.txt")
	}

	resolved := make(map[string]string)

	for _, variable := range r.Variables {
		var value string
		var resolvedOption string

		for _, option := range variable.Options {
			if option == "input" && r.Stdin != "" {
				value = r.Stdin
				resolvedOption = option
				break
			}

			if flagValue, exists := r.Flags[option]; exists {
				if isFile(flagValue) {
					content, err := readFile(flagValue)
					if err != nil {
						return nil, fmt.Errorf("failed to read file '%s': %w", flagValue, err)
					}
					value = content
				} else {
					value = flagValue
				}
				resolvedOption = option
				break
			}
		}

		if resolvedOption == "" {
			if len(variable.Options) == 1 {
				return nil, fmt.Errorf("missing required variable\n\nVariable '%s' is required\n\nUsage:\n  gliik <name> --%s <file|value>", variable.Raw, variable.Options[0])
			} else {
				var optionsHelp string
				for i, opt := range variable.Options {
					if opt == "input" {
						optionsHelp += fmt.Sprintf("  • stdin (use: cat file | gliik <name>)\n")
					} else {
						optionsHelp += fmt.Sprintf("  • --%s (use: gliik <name> --%s <file|value>)\n", opt, opt)
					}
					if i < len(variable.Options)-1 {
						optionsHelp = optionsHelp[:len(optionsHelp)-1]
					}
				}
				return nil, fmt.Errorf("missing required variable\n\nVariable '%s' needs one of:\n%s", variable.Raw, optionsHelp)
			}
		}

		resolved[variable.Raw] = value
	}

	return resolved, nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
