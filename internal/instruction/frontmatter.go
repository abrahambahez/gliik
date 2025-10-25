package instruction

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFrontmatter extracts YAML frontmatter and markdown body from instruction.md content.
// Returns the parsed metadata, markdown body content, and any error encountered.
// Frontmatter must be delimited by "---" at the start and end.
func ParseFrontmatter(content string) (Meta, string, error) {
	delimiter := "---"

	// Find first delimiter
	firstDelimiter := strings.Index(content, delimiter)
	if firstDelimiter == -1 {
		return Meta{}, "", fmt.Errorf("missing frontmatter start delimiter")
	}

	// Find second delimiter (after the first one)
	searchStart := firstDelimiter + len(delimiter)
	secondDelimiter := strings.Index(content[searchStart:], delimiter)
	if secondDelimiter == -1 {
		return Meta{}, "", fmt.Errorf("missing frontmatter end delimiter")
	}

	// Adjust second delimiter position to absolute position
	secondDelimiter += searchStart

	// Extract frontmatter (between the delimiters)
	frontmatterStart := firstDelimiter + len(delimiter)
	frontmatterYAML := content[frontmatterStart:secondDelimiter]

	// Extract body (everything after second delimiter)
	bodyStart := secondDelimiter + len(delimiter)
	body := ""
	if bodyStart < len(content) {
		body = strings.TrimLeft(content[bodyStart:], "\n")
	}

	// Parse YAML into Meta struct
	var meta Meta
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &meta); err != nil {
		return Meta{}, "", fmt.Errorf("failed to parse frontmatter YAML: %w", err)
	}

	return meta, body, nil
}
