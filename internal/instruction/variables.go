package instruction

import (
	"fmt"
	"regexp"
	"strings"
)

type Variable struct {
	Raw     string
	Options []string
}

var variableRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func ParseVariables(systemText string) ([]Variable, error) {
	matches := variableRegex.FindAllStringSubmatch(systemText, -1)
	var variables []Variable
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		raw := match[0]
		content := match[1]

		if seen[raw] {
			return nil, fmt.Errorf("duplicate variable in instruction: %s\n\nEach variable must appear only once in system.txt\nPlease remove duplicate occurrences", raw)
		}
		seen[raw] = true

		options := strings.Split(content, "|")
		for i := range options {
			options[i] = strings.TrimSpace(options[i])
		}

		variables = append(variables, Variable{
			Raw:     raw,
			Options: options,
		})
	}

	return variables, nil
}
