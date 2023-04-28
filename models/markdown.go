package models

import (
	"strings"
)

type Markdown struct {
}

func (m *Markdown) Split(text string) []string {
	// Define the separators to split the Markdown file
	separators := []string{
		"\n## ",
		"\n### ",
		"\n#### ",
		"\n##### ",
		"\n###### ",
		"```\n\n",
		"\n\n***\n\n",
		"\n\n---\n\n",
		"\n\n___\n\n",
		"\n\n",
		"\n",
	}

	// Split the Markdown text into an array of strings
	// using each separator in turn
	parts := []string{text}
	for _, sep := range separators {
		newParts := make([]string, 0)
		for _, part := range parts {
			subParts := strings.Split(part, sep)
			newParts = append(newParts, subParts...)
		}
		parts = newParts
	}
	return parts
}
