package llms

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/scottraio/plum/memory"
)

type LLM interface {
	Client() LLM
	Run(memory *memory.Memory) string
}

// InjectInputsToPrompt injects the agent's input and memory into the prompt.
func InjectObjectToPrompt(obj interface{}, prompt string) string {
	tmpl, err := template.New("").Parse(prompt)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, obj); err != nil {
		return ""
	}

	raw := buf.String()
	var cleaned strings.Builder
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)

		cleaned.WriteString(trimmed)
		cleaned.WriteString("\n")

	}

	return cleaned.String()
}
