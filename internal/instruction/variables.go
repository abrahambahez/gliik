package instruction

import (
	"regexp"
	"strings"
)

type Variable struct {
	Raw     string
	Options []string
}

var variableRegex = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func ParseVariables(systemText string) []Variable {
	matches := variableRegex.FindAllStringSubmatch(systemText, -1)
	var variables []Variable

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		raw := match[0]
		content := match[1]

		options := strings.Split(content, "|")
		for i := range options {
			options[i] = strings.TrimSpace(options[i])
		}

		variables = append(variables, Variable{
			Raw:     raw,
			Options: options,
		})
	}

	return variables
}
